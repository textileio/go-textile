package mobile

import (
	"os"
	"testing"

	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
)

var cafePath = "testdata/.textile3"
var cafe *core.Textile
var cafeApi = "http://127.0.0.1:5000"
var mobilePath = "testdata/.textile4"
var mobile *Mobile

func TestMobile_SetupCafes(t *testing.T) {
	var err error
	mobile, err = createAndStartMobile(mobilePath, false)
	if err != nil {
		t.Fatal(err)
	}

	// start a cafe
	_ = os.RemoveAll(cafePath)
	err = core.InitRepo(core.InitConfig{
		Account:     keypair.Random(),
		RepoPath:    cafePath,
		CafeApiAddr: "127.0.0.1:5000",
		CafeOpen:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	cafe, err = core.NewTextile(core.RunConfig{
		RepoPath: cafePath,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = cafe.Start()
	if err != nil {
		t.Fatal(err)
	}

	<-mobile.OnlineCh()
	<-cafe.OnlineCh()
}

func TestMobile_RegisterCafe(t *testing.T) {
	// create a token
	token, err := cafe.CreateCafeToken("", true)
	if err != nil {
		t.Fatal(err)
	}

	// register with cafe
	err = mobile.RegisterCafe(cafeApi, token)
	if err != nil {
		t.Fatal(err)
	}

	// add some data
}
