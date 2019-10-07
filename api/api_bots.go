package api

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// botsList lists all running bots
func (a *Api) botsList(c *gin.Context) {
	bots := a.Bots.List() // TODO: this should be a pb {items:[]}
	c.JSON(200, bots)
}

// botsGet is the GET endpoint for all bots
func (a *Api) botsGet(c *gin.Context) {
	botID := c.Param("root")
	if !a.Bots.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := a.Bots.Get(botID, qbytes)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	statusInt := int(botResponse.Status)

	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

// botsPost is the POST endpoint for all bots
func (a *Api) botsPost(c *gin.Context) {
	botID := c.Param("root")
	log.Errorf("botID: %s", botID)
	if !a.Bots.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
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

	botResponse, err := a.Bots.Post(botID, qbytes, body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

func (a *Api) botsDelete(c *gin.Context) {
	botID := c.Param("root")
	if !a.Bots.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := a.Bots.Delete(botID, qbytes)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

func (a *Api) botsPut(c *gin.Context) {
	botID := c.Param("root")

	if !a.Bots.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
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

	botResponse, err := a.Bots.Put(botID, qbytes, body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}
