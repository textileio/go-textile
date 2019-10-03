package api

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
	"github.com/textileio/go-textile/repo/config"
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
		if !v.IsValid() {
			return nil, fmt.Errorf("empty struct value")
		}
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
func (a *Api) getConfig(g *gin.Context) {
	pth := g.Param("path")
	conf := a.Node.Config()

	if pth == "" {
		g.JSON(http.StatusOK, conf)
	} else {
		value, err := getKeyValue(pth[1:], conf)
		if err != nil {
			g.String(http.StatusBadRequest, err.Error())
			return
		}
		g.JSON(http.StatusOK, value)
	}
}

// patchConfig godoc
// @Summary Set/update config settings
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
func (a *Api) patchConfig(g *gin.Context) {
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

	configPath := path.Join(a.RepoPath, "textile")

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

// setConfig godoc
// @Summary Replace config settings.
// @Description Replace entire config file contents. The config command controls configuration
// @Description variables. It works much like 'git config'. The configuration values are stored
// @Description in a config file inside the Textile repository.
// @Tags config
// @Accept application/json
// @Param config body mill.Json true "JSON document"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Router /config [put]
func (a *Api) setConfig(g *gin.Context) {
	configPath := path.Join(a.RepoPath, "textile")

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
