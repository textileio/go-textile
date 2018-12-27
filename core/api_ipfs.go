package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
)

func (a *api) swarmConnect(g *gin.Context) {
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

func (a *api) swarmPeers(g *gin.Context) {
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
