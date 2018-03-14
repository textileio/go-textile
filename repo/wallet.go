package repo

import (
	"time"
	"bytes"
	"encoding/json"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi/interface"
)

type Photo map[string]string

type WalletData struct {
	Photos []Photo `json:"photos"`
}

type Wallet struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Data WalletData `json:"data"`
	LastHash string `json:"last_hash"`
}

func NewWallet(node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)

	wallet := &Wallet{
		Created: time.Now(),
		Updated: time.Now(),
		Data: WalletData{
			Photos: make([]Photo, 0),
		},
	}

	wb, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	p, err := api.Unixfs().Add(node.Context(), bytes.NewReader(wb))
	if err != nil {
		return err
	}

	// done automatically?
	if err := api.Pin().Add(node.Context(), p); err != nil {
		return err
	}

	err = wallet.publish(p, node)
	if err != nil {
		return err
	}

	return nil
}

func (w *Wallet) publish(path iface.Path, node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)
	_, err := api.Name().Publish(node.Context(), path)
	return err
}

/*
package repo

import (
	"time"
	"context"
	"encoding/json"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/pin"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/path"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/namesys"
	dag "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/merkledag"
	ci "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

type Photo map[string]string

type WalletData struct {
	Photos []Photo `json:"photos"`
}

type Wallet struct {
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Data WalletData `json:"data"`
	LastHash string `json:"last_hash"`
}

func NewWallet(ctx context.Context, pub namesys.Publisher, pins pin.Pinner, key ci.PrivKey) error {

	// create an empty wallet
	wallet := &Wallet{
		Created: time.Now(),
		Updated: time.Now(),
		Data: WalletData{
			Photos: make([]Photo, 0),
		},
	}

	wb, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	// create a dag node from the empty wallet data
	wn := dag.NodeWithData(wb)

	// pin it
	err = pins.Pin(ctx, wn, true)
	if err != nil {
		return err
	}

	err = pins.Flush()
	if err != nil {
		return err
	}

	return pub.Publish(ctx, key, path.FromCid(wn.Cid()))
}

func (w *Wallet) publish(path iface.Path, node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)
	_, err := api.Name().Publish(node.Context(), path)
	return err
}

 */