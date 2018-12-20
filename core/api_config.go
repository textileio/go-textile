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
