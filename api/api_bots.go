package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/repo/config"
)

// botsList lists all running bots
func (a *Api) botsList(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.Bots.List())
}

func (a *Api) botsDisable(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	botID := opts["id"]

	configPath := path.Join(a.RepoPath, "textile")

	original, err := ioutil.ReadFile(configPath)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// load config
	conf := config.Config{}
	err = json.Unmarshal(original, &conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	var targetDel int = -1
	for i, b := range conf.Bots {
		if b.ID == botID {
			targetDel = i
		}
	}

	if targetDel == -1 {
		g.String(http.StatusBadRequest, "Bot not enabled")
		return
	}

	// Remove the bot from the array
	conf.Bots[targetDel] = conf.Bots[len(conf.Bots)-1]
	conf.Bots[len(conf.Bots)-1] = config.EnabledBot{}
	conf.Bots = conf.Bots[:len(conf.Bots)-1]

	jsn, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := ioutil.WriteFile(configPath, jsn, 0666); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.Writer.WriteHeader(http.StatusNoContent)
}

func (a *Api) botsEnable(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	botID := opts["id"]
	cafeApi, err := strconv.ParseBool(opts["cafe"])
	if err != nil {
		cafeApi = false
	}

	botPath := path.Join(a.RepoPath, "bots", botID)
	// Check that the provided dir exists
	if _, err := os.Stat(botPath); os.IsNotExist(err) {
		g.String(http.StatusBadRequest, "Bot not known")
		return
	}

	configPath := path.Join(a.RepoPath, "textile")

	original, err := ioutil.ReadFile(configPath)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// load config
	conf := config.Config{}
	err = json.Unmarshal(original, &conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	for _, b := range conf.Bots {
		if b.ID == botID {
			g.String(http.StatusOK, "Bot already enabled")
			return
		}
	}

	conf.Bots = append(conf.Bots, config.EnabledBot{botID, cafeApi})

	jsn, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := ioutil.WriteFile(configPath, jsn, 0666); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.Writer.WriteHeader(http.StatusNoContent)
}

// botsGet is the GET endpoint for all bots
func (a *Api) botsGet(c *gin.Context) {
	botID := c.Param("id")

	if botID == "" {
		c.String(http.StatusBadRequest, "bot id required")
		return
	}

	if !a.Bots.Exists(botID) {
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
	botID := c.Param("id")
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
	botID := c.Param("id")
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
	botID := c.Param("id")

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
