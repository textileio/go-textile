package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
)

func (a *api) ipfsId(g *gin.Context) {
	pid, err := a.node.PeerId()
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
// @Param X-Textile-Args header string true "peer address")
// @Success 200 {array} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /swarm/connect [get]
func (a *api) ipfsSwarmConnect(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing peer multi address")
		return
	}

	res, err := ipfs.SwarmConnect(a.node.node, []string{args[0]})
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
// @Router /swarm/peers [get]
func (a *api) ipfsSwarmPeers(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	verbose := opts["verbose"] == "true"
	latency := opts["latency"] == "true"
	streams := opts["streams"] == "true"
	direction := opts["direction"] == "true"

	res, err := ipfs.SwarmPeers(a.node.node, verbose, latency, streams, direction)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.JSON(http.StatusOK, res)
}

// ipfsCat godoc
// @Summary Cat IPFS data
// @Description Displays the data behind an IPFS CID (hash)
// @Tags ipfs
// @Produce application/octet-stream
// @Param cid path string true "ipfs/ipns cid"
// @Param X-Textile-Opts header string false "key: Key to decrypt data on-the-fly" default(key=)
// @Success 200 {array} byte "data"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ipfs [get]
func (a *api) ipfsCat(g *gin.Context) {
	cid := g.Param("cid")
	if cid == "" {
		g.String(http.StatusBadRequest, "Missing IPFS CID")
	}

	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	data, err := ipfs.DataAtPath(a.node.node, cid)
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
