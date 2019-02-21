package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/repo/config"
)

func getKeyValue(path string, object interface{}) (interface{}, error) {
	keys := strings.Split(path, "/")
	v := reflect.ValueOf(object)
	for _, key := range keys {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("unable to parse struct path; ecountered %T", v)
		}

		v = v.FieldByName(key)
	}
	return v.Interface(), nil
}

// getConfig godoc
// @Summary Get active config settings
// @Description Report the currently active config settings, which may differ from the values
// @Description specifed when setting/patching values.
// @Tags config
// @Produce application/json
// @Param path path string false "config path (e.g., Addresses/API)"
// @Success 200 {object} mill.Json "new config value"
// @Failure 400 {string} string "Bad Request"
// @Router /config/{path} [get]
func (a *api) getConfig(g *gin.Context) {
	path := g.Param("path")
	conf := a.node.Config()

	if path == "" {
		g.JSON(http.StatusOK, conf)
	} else {
		value, err := getKeyValue(path[1:], conf)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		g.JSON(http.StatusOK, value)
	}
}

// patchConfig godoc
// @Summary Get active config settings
// @Description When patching config values, valid JSON types must be used. For example, a string
// @Description should be escaped or wrapped in single quotes (e.g., \"127.0.0.1:40600\") and
// @Description arrays and objects work fine (e.g. '{"API": "127.0.0.1:40600"}') but should be
// @Description wrapped in single quotes. Be sure to restart the daemon for changes to take effect.
// @Description See https://tools.ietf.org/html/rfc6902 for details on RFC6902 JSON patch format.
// @Tags config
// @Accept application/json
// @Param patch body mill.Json true "An RFC6902 JSON patch (array of ops)"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Router /config [patch]
func (a *api) patchConfig(g *gin.Context) {
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer g.Request.Body.Close()

	// decode request body into a RFC 6902 patch (array of ops)
	patch, err := jsonpatch.DecodePatch(body)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	configPath := path.Join(a.node.repoPath, "textile")

	original, err := ioutil.ReadFile(configPath)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// apply json patch to config
	modified, err := patch.Apply(original)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	// make sure our config is still valid
	conf := config.Config{}
	err = json.Unmarshal(modified, &conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	json, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	err = ioutil.WriteFile(configPath, json, 0666)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.Writer.WriteHeader(http.StatusNoContent)
}

// setConfig godoc
// @Summary Replce config settings.
// @Description Replace entire config file contents. The config command controls configuration
// @Description variables. It works much like 'git config'. The configuration values are stored
// @Description in a config file inside the Textile repository.
// @Tags config
// @Accept application/json
// @Param config body mill.Json true "JSON document"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Router /config [put]
func (a *api) setConfig(g *gin.Context) {
	configPath := path.Join(a.node.repoPath, "textile")

	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	defer g.Request.Body.Close()

	conf := config.Config{}
	err = json.Unmarshal(body, &conf)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	json, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	err = ioutil.WriteFile(configPath, json, 0666)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.Writer.WriteHeader(http.StatusNoContent)
}
