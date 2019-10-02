package bots

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/core"
)

type api struct {
	node    *core.Textile
	service Service
}

// botsGet is the GET endpoint for all bots
func (a *api) botsGet(c *gin.Context) {
	botID := c.Param("root")
	if !a.service.Exists(botID) { // bot doesn't exist yet
		// log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := a.service.Get(botID, qbytes)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

// botsPost is the POST endpoint for all bots
func (a *api) botsPost(c *gin.Context) {
	botID := c.Param("root")

	if !a.service.Exists(botID) { // bot doesn't exist yet
		// log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	botResponse, err := a.service.Post(botID, qbytes, body)
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

func (a *api) botsDelete(c *gin.Context) {
	botID := c.Param("root")
	if !a.service.Exists(botID) { // bot doesn't exist yet
		// log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := a.service.Delete(botID, qbytes)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

func (a *api) botsPut(c *gin.Context) {
	botID := c.Param("root")
	if !a.service.Exists(botID) { // bot doesn't exist yet
		// log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	botResponse, err := a.service.Put(botID, qbytes, body)
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}
