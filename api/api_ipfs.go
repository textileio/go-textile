package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/mr-tron/base58/base58"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

// ipfsId godoc
// @Summary Get IPFS peer ID
// @Description Displays underlying IPFS peer ID
// @Tags ipfs
// @Produce text/plain
// @Success 200 {string} string "peer id"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ipfs/id [get]
func (a *Api) ipfsId(g *gin.Context) {
	pid, err := a.Node.PeerId()
	if err != nil {
		a.abort500(g, err)
		return
	}
	g.String(http.StatusOK, pid.Pretty())
}

// ipfsSwarmConnect godoc
// @Summary Opens a new direct connection to a peer address
// @Description Opens a new direct connection to a peer using an IPFS multiaddr
// @Tags ipfs
// @Produce application/json
// @Param X-Textile-Args header string true "peer address"
// @Success 200 {array} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ipfs/swarm/connect [post]
func (a *Api) ipfsSwarmConnect(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing peer multi address")
		return
	}

	res, err := ipfs.SwarmConnect(a.Node.Ipfs(), []string{args[0]})
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, res)
}

// ipfsSwarmPeers godoc
// @Summary List swarm peers
// @Description Lists the set of peers this node is connected to
// @Tags ipfs
// @Produce application/json
// @Param X-Textile-Opts header string false "verbose: Display all extra information, latency: Also list information about latency to each peer, streams: Also list information about open streams for each peer, direction: Also list information about the direction of connection" default(verbose="false",latency="false",streams="false",direction="false")
// @Success 200 {object} ipfs.ConnInfos "connection"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ipfs/swarm/peers [get]
func (a *Api) ipfsSwarmPeers(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	verbose := opts["verbose"] == "true"
	latency := opts["latency"] == "true"
	streams := opts["streams"] == "true"
	direction := opts["direction"] == "true"

	res, err := ipfs.SwarmPeers(a.Node.Ipfs(), verbose, latency, streams, direction)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, res)
}

// ipfsCat godoc
// @Summary Cat IPFS data
// @Description Displays the data behind an IPFS CID (hash) or Path
// @Tags ipfs
// @Produce application/octet-stream
// @Param path path string true "ipfs/ipns cid"
// @Param X-Textile-Opts header string false "key: Key to decrypt data on-the-fly" default(key=)
// @Success 200 {array} byte "data"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ipfs/cat/{path} [get]
func (a *Api) ipfsCat(g *gin.Context) {
	pth := g.Param("path")
	if pth == "" {
		g.String(http.StatusBadRequest, "Missing IPFS CID")
	}

	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	data, err := ipfs.DataAtPath(a.Node.Ipfs(), pth[1:])
	if err != nil {
		g.String(http.StatusNotFound, err.Error())
		return
	}

	var plaintext []byte
	if opts["key"] != "" {
		key, err := base58.Decode(opts["key"])
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
		}
		plaintext, err = crypto.DecryptAES(data, key)
		if err != nil {
			g.String(http.StatusUnauthorized, err.Error())
		}
	} else {
		plaintext = data
	}

	g.Data(http.StatusOK, "application/octet-stream", plaintext)
}

// ipfsPubsubPub godoc
// @Summary Publish a message
// @Description Publishes a message to a given pubsub topic
// @Tags ipfs
// @Accept application/octet-stream
// @Produce text/plain
// @Success 204 {string} string "ok"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ipfs/pubsub/pub/{topic} [post]
func (a *Api) ipfsPubsubPub(g *gin.Context) {
	topic := g.Param("topic")
	if topic == "" {
		g.String(http.StatusBadRequest, "missing topic")
	}

	payload, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	err = ipfs.Publish(a.Node.Ipfs(), topic, payload)
	if err != nil {
		a.abort500(g, err)
		return
	}

	g.Status(http.StatusNoContent)
}

// ipfsPubsubSub godoc
// @Summary Subscribe messages
// @Description Subscribes to messages on a given topic
// @Tags ipfs
// @Produce text/event-stream with events, or just application/json
// @Success 200 {string} string with events, or []byte "results stream"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ipfs/pubsub/sub/{topic} [get]
func (a *Api) ipfsPubsubSub(g *gin.Context) {
	topic := g.Param("topic")
	if topic == "" {
		g.String(http.StatusBadRequest, "Missing topic")
	}

	events := g.Query("events")
	queryId := g.Query("queryId")

	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	msgs := make(chan iface.PubSubMessage, 10)
	ctx := a.Node.Ipfs().Context()
	var id string
	if queryId != "" {
		id = queryId
	} else if opts["queryId"] != "" {
		id = opts["queryId"]
	} else {
		id = ksuid.New().String()
	}

	go func() {
		if err := ipfs.Subscribe(a.Node.Ipfs(), ctx, topic, true, msgs); err != nil {
			close(msgs)
			a.abort500(g, err)
			log.Errorf("ipfs pubsub sub stopped with error: %s", err.Error())
			return
		}
	}()
	log.Infof("ipfs pubsub sub started for %s", topic)

	// EventSource of web browser will not pending on first msg from g.SSEvent,
	// initHttpStatusOK here will not pending on first msg from g.Data for consistency
	initHttpStatusOK := false

	g.Stream(func(w io.Writer) bool {
		if events != "true" && !initHttpStatusOK {
			initHttpStatusOK = true
			g.Data(http.StatusOK, "application/octet-stream", []byte{})
			return true
		}

		select {
		case <-g.Request.Context().Done():
			log.Infof("ipfs pubsub sub shutdown for %s", topic)
			return false
		case msg, ok := <-msgs:
			if !ok {
				log.Infof("ipfs pubsub sub shutdown for %s", topic)
				return false
			}

			mPeer := msg.From()
			if mPeer.Pretty() == a.Node.Ipfs().Identity.Pretty() {
				return true
			}

			value, err := proto.Marshal(&pb.Strings{
				Values: []string{string(msg.Data())},
			})
			if err != nil {
				g.String(http.StatusBadRequest, err.Error())
				break
			}

			res := &pb.QueryResult{
				Id:    fmt.Sprintf("%x", msg.Seq()),
				Value: &any.Any{
					// Can't let TypeUrl to use ipfs official "/pubsub.pb.Message"
					// from github.com/libp2p/go-libp2p-pubsub/pb/rpc.pb.go ,
					// because rpc.pb.go import github.com/gogo/protobuf/proto , so
					// proto.RegisterType((*Message)(nil), "pubsub.pb.Message") in
					// rpc.pb.go , can't let "/pubsub.pb.Message" be found in
					// pbMarshaler.MarshalToString which import github.com/golang/protobuf
					TypeUrl: "/Strings",
					Value:   value,
				},
			}
			str, err := pbMarshaler.MarshalToString(&pb.MobileQueryEvent{
				Id:   id,
				Type: pb.MobileQueryEvent_DATA,
				Data: res,
			})
			if err != nil {
				g.String(http.StatusBadRequest, err.Error())
				break
			}

			if events == "true" || opts["events"] == "true" {
				g.SSEvent("update", str)
			} else {
				g.Data(http.StatusOK, "application/json", []byte(str))
				g.Writer.Write([]byte("\n"))
			}
		}
		return true
	})
}
