package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) botsGet(c *gin.Context) {
	botID := c.Param("root")
	botService, err := a.node.Bots()
	if err != nil {
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}
	if !botService.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := botService.Get(botID, qbytes)
	statusInt := int(botResponse.Status)

	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

func (a *api) botsPost(c *gin.Context) {
	botID := c.Param("root")
	botService, err := a.node.Bots()
	if err != nil {
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}
	if !botService.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := botService.Post(botID, qbytes)
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

func (a *api) botsDelete(c *gin.Context) {
	botID := c.Param("root")
	botService, err := a.node.Bots()
	if err != nil {
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}
	if !botService.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := botService.Delete(botID, qbytes)
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}

func (a *api) botsPut(c *gin.Context) {
	botID := c.Param("root")
	botService, err := a.node.Bots()
	if err != nil {
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}
	if !botService.Exists(botID) { // bot doesn't exist yet
		log.Errorf("error bot not found: %s", botID)
		c.String(http.StatusBadRequest, "bot not found")
		return
	}

	query := c.Request.URL.Query().Encode()
	qbytes := []byte(query)

	botResponse, err := botService.Put(botID, qbytes)
	statusInt := int(botResponse.Status)
	c.Data(statusInt, botResponse.ContentType, botResponse.Body)
}
