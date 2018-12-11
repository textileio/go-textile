package core

import (
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"
	logger "gx/ipfs/QmcaSwFc5RBg8yCq54QURwEU4nwjfCpjbpmaAm4VbdGLKv/go-logging"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *api) logsCall(g *gin.Context) {

	subsystem := g.Param("subsystem")
	var subsystems []string
	if subsystem == "" {
		subsystems = logging.GetSubsystems()
	} else {
		subsystems = []string{subsystem}
	}

	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}
	level := strings.ToUpper(opts["level"])
	texOnly := opts["tex-only"] == "true"

	result := make(map[string]string)
	for _, system := range subsystems {
		if texOnly && !strings.HasPrefix(system, "tex") {
			continue
		}
		var llevel logger.Level
		if level != "" && g.Request.Method == "POST" {
			// validate log level
			llevel, err = logger.LogLevel(level)
			if err != nil {
				g.String(http.StatusBadRequest, err.Error())
				return
			}
		} else {
			llevel = logger.GetLevel(system)
		}
		// validate subsystem + set log level
		err = logging.SetLogLevel(system, llevel.String())
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		result[system] = llevel.String()
	}
	g.JSON(http.StatusOK, result)
}
