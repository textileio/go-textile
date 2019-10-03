package api

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// publish godoc
// @Summary Publish payload to topic
// @Description Publishes payload bytes to a topic on the network.
// @Tags utils
// @Accept application/octet-stream
// @Produce text/plain
// @Param X-Textile-Args header string true "topic"
// @Success 204 {string} string "ok"
// @Failure 500 {string} string "Internal Server Error"
// @Router /publish [post]
func (a *Api) publish(g *gin.Context) {
	args, err := a.readArgs(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	if len(args) == 0 {
		g.String(http.StatusBadRequest, "missing topic")
		return
	}

	payload, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	err = a.Node.Publish(payload, args[0])
	if err != nil {
		a.abort500(g, err)
		return
	}

	g.Status(http.StatusNoContent)
}
