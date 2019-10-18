package cmd

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	shared "github.com/textileio/go-textile-core/bots"
	"github.com/textileio/go-textile/ipfs"
)

// BotsList lists all enabled bots
func BotsList() error {
	res, err := executeJsonCmd(http.MethodGet, "bots/list", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

// BotsCreate writes a new bot config to the current repo
func BotsCreate(name string) error {
	// create an identity for the ipfs peer
	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return err
	}
	peerIdent, err := ipfs.IdentityConfig(sk)
	if err != nil {
		return err
	}

	conf := &shared.HostConfig{
		Name:           name,
		ID:             peerIdent.PeerID,
		ReleaseVersion: 0,
		ReleaseHash:    "",
		Params:         map[string]string{},
	}

	jsn, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("./config", jsn, 0666); err != nil {
		return err
	}

	res := fmt.Sprintf("Bot secret key: %s", peerIdent.PrivKey)
	output(res)
	return nil
}

// BotsDisable disables a bot
func BotsDisable(id string) error {
	res, err := executeJsonCmd(http.MethodPost, "bots/disable", params{
		opts: map[string]string{"id": id},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

// BotsEnable enables a known bot
func BotsEnable(id string, cafe bool) error {
	res, err := executeJsonCmd(http.MethodPost, "bots/enable", params{
		opts: map[string]string{"id": id, "cafe": strconv.FormatBool(cafe)},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
