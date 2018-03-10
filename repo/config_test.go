package repo

import (
	"reflect"
	"testing"

	//"github.com/ipfs/go-ipfs/repo/fsrepo"
	"io/ioutil"
	//"os"
	//"path/filepath"
	//"time"
)

const testConfigFolder = "testdata"
const testConfigPath = "testdata/config"

func TestGetApiConfig(t *testing.T) {
	configFile, err := ioutil.ReadFile(testConfigPath)
	if err != nil {
		t.Error(err)
	}

	config, err := GetAPIConfig(configFile)
	if config.Username != "TestUsername" {
		t.Error("Expected TestUsername, got ", config.Username)
	}
	if config.Password != "TestPassword" {
		t.Error("Expected TestPassword, got ", config.Password)
	}
	if !config.Authenticated {
		t.Error("Expected Authenticated = true")
	}
	if len(config.AllowedIPs) != 1 || config.AllowedIPs[0] != "127.0.0.1" {
		t.Error("Expected AllowedIPs = [127.0.0.1]")
	}
	if config.CORS == nil {
		t.Error("Cors is not set")
	}
	if reflect.ValueOf(config.HTTPHeaders).Kind() != reflect.Map {
		t.Error("Headers is not a map")
	}
	if config.Enabled != true {
		t.Error("Enabled is not true")
	}
	if !config.SSL {
		t.Error("Expected SSL = true")
	}
	if config.SSLCert == "" {
		t.Error("Expected test SSL cert, got ", config.SSLCert)
	}
	if config.SSLKey == "" {
		t.Error("Expected test SSL key, got ", config.SSLKey)
	}
	if err != nil {
		t.Error("GetAPIAuthentication threw an unexpected error")
	}

	_, err = GetAPIConfig([]byte{})
	if err == nil {
		t.Error("GetAPIAuthentication didn`t throw an error")
	}
}

//func TestGetResolverConfig(t *testing.T) {
//	configFile, err := ioutil.ReadFile(testConfigPath)
//	if err != nil {
//		t.Error(err)
//	}
//	resolvers, err := GetResolverConfig(configFile)
//	if err != nil {
//		t.Error("GetResolverUrl threw an unexpected error")
//	}
//	if resolvers.Id != "https://resolver.onename.com/" {
//		t.Error("resolverUrl does not equal expected value")
//	}
//}

//func TestExtendConfigFile(t *testing.T) {
//	r, err := fsrepo.Open(testConfigFolder)
//	if err != nil {
//		t.Error("fsrepo.Open threw an unexpected error", err)
//		return
//	}
//	configFile, err := ioutil.ReadFile(testConfigPath)
//	if err != nil {
//		t.Error(err)
//	}
//	config, _ := GetWalletConfig(configFile)
//	originalMaxFee := config.MaxFee
//	newMaxFee := config.MaxFee + 1
//	if err := extendConfigFile(r, "Wallet.MaxFee", newMaxFee); err != nil {
//		t.Error("extendConfigFile threw an unexpected error ", err)
//		return
//	}
//	configFile, err = ioutil.ReadFile(testConfigPath)
//	if err != nil {
//		t.Error(err)
//	}
//	config, _ = GetWalletConfig(configFile)
//	if config.MaxFee != newMaxFee {
//		t.Errorf("Expected maxFee to be %v, got %v", newMaxFee, config.MaxFee)
//		return
//	}
//	// Reset maxFee to original value
//	extendConfigFile(r, "Wallet.MaxFee", originalMaxFee)
//
//	// Teardown
//	os.RemoveAll(filepath.Join(testConfigFolder, "datastore"))
//	os.RemoveAll(filepath.Join(testConfigFolder, "repo.lock"))
//}

func TestInitConfig(t *testing.T) {
	config, err := InitConfig(testConfigFolder)
	if config == nil {
		t.Error("config empty", err)
	}
	if err != nil {
		t.Error("InitConfig threw an unexpected error")
	}
	if config.Addresses.Gateway != "/ip4/127.0.0.1/tcp/4002" {
		t.Error("config.Addresses.Gateway is not set")
	}
}
