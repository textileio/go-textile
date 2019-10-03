package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log"
	logger "github.com/whyrusleeping/go-logging"
)

type SubsystemInfo map[string]string

// logsCall godoc
// @Summary Access subsystem logs
// @Description List or change the verbosity of one or all subsystems log output. Textile logs
// @Description piggyback on the IPFS event logs
// @Tags utils
// @Produce application/json
// @Param subsystem path string false "subsystem logging identifier (omit for all)"
// @Param X-Textile-Opts header string false "level: Log-level (one of: debug, info, warning, error, critical, or "" to get current), tex-only: Whether to list/change only Textile subsystems, or all available subsystems" default(level=,tex-only="false")
// @Success 200 {object} core.SubsystemInfo "subsystems"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /logs/{subsystem} [post]
func (a *Api) logsCall(g *gin.Context) {

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
