package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"time"

	njwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/pin"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/jwt"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

// GET /sessions/challenge/?account_addr=<address> (header=>token)
func (c *cafeApi) getSessionChallenge(g *gin.Context) {
	addr := g.Query("account_addr")
	accnt, err := keypair.Parse(addr)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we don't want to handle account seeds, just addresses
		log.Warning(errInvalidAddress)
		c.abort(g, http.StatusBadRequest, fmt.Errorf(errInvalidAddress))
		return
	}

	// generate a new random nonce
	nonce := &pb.CafeClientNonce{
		Value:   ksuid.New().String(),
		Address: addr,
		Date:    ptypes.TimestampNow(),
	}
	err = c.node.datastore.CafeClientNonces().Add(nonce)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	pbJSON(g, http.StatusOK, nonce)
}

// POST /sessions/:pid/?account_addr=<address>&challenge=<challenge>&n=<nonce> (header=>token, body=sig)
func (c *cafeApi) createSession(g *gin.Context) {
	// are we open?
	if !c.node.cafe.open {
		log.Warning("cafe is not open")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	addr := g.Query("account_addr")
	accnt, err := keypair.Parse(addr)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	pid, err := peer.IDB58Decode(g.Param("pid"))
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	// check nonce
	snonce := c.node.datastore.CafeClientNonces().Get(g.Query("challenge"))
	if snonce == nil {
		log.Warning("challenge not found")
		c.abort(g, http.StatusForbidden, nil)
		return
	}
	if snonce.Address != accnt.Address() {
		log.Warning("invalid address")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()

	buf.Grow(bytes.MinRead)
	_, err = buf.ReadFrom(g.Request.Body)
	if err != nil && err != io.EOF {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	sig := buf.Bytes()

	payload := []byte(g.Query("challenge") + g.Query("nonce"))
	err = accnt.Verify(payload, sig)
	if err != nil {
		log.Warning("verification failed")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	now := ptypes.TimestampNow()
	client := &pb.CafeClient{
		Id:      pid.Pretty(),
		Address: accnt.Address(),
		Created: now,
		Seen:    now,
		Token:   g.GetString("token"),
	}
	err = c.node.datastore.CafeClients().Add(client)
	if err != nil {
		// check if already exists
		client = c.node.datastore.CafeClients().Get(pid.Pretty())
		if client == nil {
			err = fmt.Errorf("get or create client failed")
			log.Warning(err)
			c.abort(g, http.StatusInternalServerError, err)
			return
		}
	}

	session, err := jwt.NewSession(
		c.node.node.PrivateKey,
		pid,
		c.node.cafe.Protocol(),
		defaultSessionDuration,
		c.node.cafe.info,
	)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	err = c.node.datastore.CafeClientNonces().Delete(snonce.Value)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	pbJSON(g, http.StatusCreated, session)
}

// POST /sessions/refresh/:pid (header=>refresh, body=access)
func (c *cafeApi) refreshSession(g *gin.Context) {
	// are we _still_ open?
	if !c.node.cafe.open {
		log.Warning("cafe is not open")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()

	buf.Grow(bytes.MinRead)
	_, err := buf.ReadFrom(g.Request.Body)
	if err != nil && err != io.EOF {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	// ensure access and refresh are a valid pair
	access, _ := njwt.Parse(string(buf.Bytes()), c.verifyKeyFunc)
	if access == nil {
		log.Warning("error parsing access token")
		c.abort(g, http.StatusForbidden, nil)
		return
	}
	refresh, _ := njwt.Parse(g.GetString("token"), c.verifyKeyFunc)
	if refresh == nil {
		log.Warning("error parsing refresh token")
		c.abort(g, http.StatusForbidden, nil)
		return
	}
	accessClaims, err := jwt.ParseClaims(access.Claims)
	if err != nil {
		log.Warning("error parsing access claims")
		c.abort(g, http.StatusForbidden, nil)
		return
	}
	refreshClaims, err := jwt.ParseClaims(refresh.Claims)
	if err != nil {
		log.Warning("error parsing refresh claims")
		c.abort(g, http.StatusForbidden, nil)
		return
	}
	if refreshClaims.Id[1:] != accessClaims.Id {
		log.Warning("token id mismatch")
		c.abort(g, http.StatusForbidden, nil)
		return
	}
	if refreshClaims.Subject != accessClaims.Subject {
		log.Warning("token subject mismatch ")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	// get a new session
	spid, err := peer.IDB58Decode(accessClaims.Subject)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}
	session, err := jwt.NewSession(
		c.node.node.PrivateKey,
		spid,
		c.node.cafe.Protocol(),
		defaultSessionDuration,
		c.node.cafe.info,
	)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	pbJSON(g, http.StatusCreated, session)
}

// DELETE /sessions/:pid (header=>access)
func (c *cafeApi) deleteSession(g *gin.Context) {
	pid := g.GetString("from")
	err := c.node.datastore.CafeClientThreads().DeleteByClient(pid)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	err = c.node.datastore.CafeClientMessages().DeleteByClient(pid, -1)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	err = c.node.datastore.CafeClients().Delete(pid)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) store(g *gin.Context) {
	var err error
	var aid *cid.Cid

	form, err := g.MultipartForm()
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	files := form.File["file"]

	var f multipart.File
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	for _, file := range files {
		f, err = file.Open()
		if err != nil {
			log.Warning(err)
			c.abort(g, http.StatusBadRequest, err)
			return
		}

		aid, err = ipfs.AddObject(c.node.Ipfs(), f, true)
		if err != nil {
			_, _ = f.Seek(0, 0)
			aid, err = ipfs.AddData(c.node.Ipfs(), f, true, false)
		}
		if err != nil {
			log.Warning(err)
			c.abort(g, http.StatusBadRequest, err)
			return
		}

		log.Debugf("stored %s", aid.Hash().B58String())

		f.Close()
		f = nil
	}

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) unstore(g *gin.Context) {
	id, err := cid.Decode(g.Param("cid"))
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	pinned, err := c.node.Ipfs().Pinning.CheckIfPinned(id)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	for _, p := range pinned {
		if p.Mode != pin.NotPinned {
			err = ipfs.UnpinCid(c.node.Ipfs(), p.Key, true)
			if err != nil {
				log.Warning(err)
				c.abort(g, http.StatusBadRequest, err)
				return
			}

			log.Debugf("unstored %s", p.Key.Hash().B58String())
		}
	}

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) storeThread(g *gin.Context) {
	from := g.GetString("from")
	id := g.Param("id")

	client := c.node.datastore.CafeClients().Get(from)
	if client == nil {
		log.Warning("client not found")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()

	buf.Grow(bytes.MinRead)
	_, err := buf.ReadFrom(g.Request.Body)
	if err != nil && err != io.EOF {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	err = c.node.datastore.CafeClientThreads().AddOrUpdate(&pb.CafeClientThread{
		Id:         id,
		Client:     client.Id,
		Ciphertext: buf.Bytes(),
	})
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	log.Debugf("stored thread %s", id)

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) unstoreThread(g *gin.Context) {
	from := g.GetString("from")
	id := g.Param("id")

	client := c.node.datastore.CafeClients().Get(from)
	if client == nil {
		log.Warning("client not found")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	err := c.node.datastore.CafeClientThreads().Delete(id, client.Id)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	log.Debugf("unstored thread %s", id)

	g.Status(http.StatusNoContent)
}

func (c *cafeApi) checkMessages(g *gin.Context) {
	client := c.node.datastore.CafeClients().Get(g.GetString("from"))
	if client == nil {
		log.Warning("client not found")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	err := c.node.datastore.CafeClients().UpdateLastSeen(client.Id, time.Now())
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	res := &pb.CafeMessages{
		Messages: make([]*pb.CafeMessage, 0),
	}
	msgs := c.node.datastore.CafeClientMessages().ListByClient(client.Id, inboxMessagePageSize)
	for _, msg := range msgs {
		res.Messages = append(res.Messages, &pb.CafeMessage{
			Id:   msg.Id,
			Peer: msg.Peer,
			Date: msg.Date,
		})
	}

	pbJSON(g, http.StatusOK, res)
}

func (c *cafeApi) deleteMessages(g *gin.Context) {
	client := c.node.datastore.CafeClients().Get(g.GetString("from"))
	if client == nil {
		log.Warning("client not found")
		c.abort(g, http.StatusForbidden, nil)
		return
	}

	// delete the most recent page
	err := c.node.datastore.CafeClientMessages().DeleteByClient(client.Id, inboxMessagePageSize)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	// check for more
	remaining := c.node.datastore.CafeClientMessages().CountByClient(client.Id)

	res := &pb.CafeDeleteMessagesAck{More: remaining > 0}
	pbJSON(g, http.StatusOK, res)
}

func (c *cafeApi) deliverMessage(g *gin.Context) {
	from := g.Param("from")
	clientId := g.Param("to")

	client := c.node.datastore.CafeClients().Get(clientId)
	if client == nil {
		log.Warningf("received message for unknown client %s", clientId)
		g.Status(http.StatusOK)
		return
	}

	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()

	buf.Grow(bytes.MinRead)
	_, err := buf.ReadFrom(g.Request.Body)
	if err != nil && err != io.EOF {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	body := buf.Bytes()

	// pin inner node
	nenv := new(pb.Envelope)
	err = proto.Unmarshal(body, nenv)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	tenv := new(pb.ThreadEnvelope)
	err = ptypes.UnmarshalAny(nenv.Message.Payload, tenv)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	oid, err := ipfs.AddObject(c.node.Ipfs(), bytes.NewReader(tenv.Node), true)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	node, err := ipfs.NodeAtCid(c.node.Ipfs(), *oid)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}
	if tenv.Block != nil {
		_, err = ipfs.AddData(c.node.Ipfs(), bytes.NewReader(tenv.Block), true, false)
		if err != nil {
			log.Warning(err)
			c.abort(g, http.StatusBadRequest, err)
			return
		}
	}
	_, err = extractNode(c.node.Ipfs(), node, tenv.Block == nil)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	// pin envelope
	id, err := ipfs.AddData(c.node.Ipfs(), bytes.NewReader(body), true, false)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	msgId := id.Hash().B58String()
	err = c.node.datastore.CafeClientMessages().AddOrUpdate(&pb.CafeClientMessage{
		Id:     msgId,
		Peer:   from,
		Client: client.Id,
		Date:   ptypes.TimestampNow(),
	})
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusInternalServerError, err)
		return
	}

	go func() {
		err = c.node.cafe.notifyClient(client.Id)
		if err != nil {
			log.Debugf("unable to notify client: %s", client.Id)
		}
	}()

	log.Debugf("delivered message %s", msgId)

	g.Status(http.StatusOK)
}

func (c *cafeApi) search(g *gin.Context) {
	from := g.GetString("from")

	pid, err := peer.IDB58Decode(from)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()

	buf.Grow(bytes.MinRead)
	_, err = buf.ReadFrom(g.Request.Body)
	if err != nil && err != io.EOF {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	// parse body as a service envelope
	pmes := new(pb.Envelope)
	err = proto.Unmarshal(buf.Bytes(), pmes)
	if err != nil {
		log.Warning(err)
		c.abort(g, http.StatusBadRequest, err)
		return
	}

	// handle the message as a JSON stream
	rpmesCh, errCh, cancel := c.node.cafe.HandleStream(pmes, pid)
	g.Stream(func(w io.Writer) bool {
		select {
		case <-g.Request.Context().Done():
			log.Debug("closing request stream")
			close(cancel)

		case err := <-errCh:
			log.Warning(err)
			c.abort(g, http.StatusBadRequest, err)
			return false

		case rpmes, ok := <-rpmesCh:
			if !ok {
				g.Status(http.StatusOK)
				return false
			}
			log.Debugf("responding with %s", rpmes.Message.Type.String())

			payload, err := proto.Marshal(rpmes)
			if err != nil {
				log.Warning(err)
				c.abort(g, http.StatusInternalServerError, err)
				return false
			}

			size := make([]byte, 2)
			binary.LittleEndian.PutUint16(size, uint16(len(payload)))

			payload = append(size, payload...)
			g.Data(http.StatusOK, "application/octet-stream", payload)
		}
		return true
	})
}

// ReverseProxyBotAPI generates a function for per-method reverse proxy
func (c *cafeApi) reverseProxyBotAPI(method string) gin.HandlerFunc {
	conf := c.node.Config()
	return func(g *gin.Context) {
		id := g.Param("id")
		enabled := false
		for _, b := range conf.Bots {
			if b.ID == id && b.CafeAPI == true {
				enabled = true
				break
			}
		}
		if enabled {
			s := fmt.Sprintf("/api/v0/bots/id/%s", id)
			director := func(req *http.Request) {
				req.URL.Path = s
				req.Method = method
				req.URL.Scheme = "http"
				req.URL.Host = c.node.config.Addresses.API
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(g.Writer, g.Request)
		} else {
			g.String(404, "")
		}
	}
}

// sendError sends the error to the gin context
func sendError(g *gin.Context, err error, statusCode int) {
	g.String(statusCode, err.Error())
}

// pbJSON responds with a JSON rendered protobuf message
func pbJSON(g *gin.Context, status int, msg proto.Message) {
	str, err := pbMarshaler.MarshalToString(msg)
	if err != nil {
		sendError(g, err, http.StatusBadRequest)
		return
	}
	g.Data(status, "application/json", []byte(str))
}
