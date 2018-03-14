package repo

import (
	"time"
	"bytes"
	"image/jpeg"
	"fmt"
	"encoding/json"

	"github.com/textileio/textile-go/repo/images"

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

func (w *Wallet) PinPhoto(base64ImageData string, node *core.IpfsNode) error {
	// decode image
	im, cfg, err := images.DecodeImageData(base64ImageData)
	if err != nil {
		return err
	}
	imb := new(bytes.Buffer)
	if err = jpeg.Encode(imb, im, &jpeg.Options{ Quality: 100 }); err != nil {
		return err
	}

	// create thumbnail
	th := images.ResizeImage(im, cfg, 80, 80)
	thb := new(bytes.Buffer)
	if err = jpeg.Encode(thb, th, nil); err != nil {
		return err
	}

	// add files to ipfs
	api := coreapi.NewCoreAPI(node)
	imp, err := api.Unixfs().Add(node.Context(), bytes.NewReader(imb.Bytes()))
	if err != nil {
		return err
	}
	thp, err := api.Unixfs().Add(node.Context(), bytes.NewReader(thb.Bytes()))
	if err != nil {
		return err
	}

	fmt.Println(imp.Cid().String())
	fmt.Println(thp.Cid().String())

	// done automatically?
	//if err := api.Pin().Add(node.Context(), p); err != nil {
	//	return err
	//}

	return nil
}

func (w *Wallet) publish(path iface.Path, node *core.IpfsNode) error {
	api := coreapi.NewCoreAPI(node)
	_, err := api.Name().Publish(node.Context(), path)
	return err
}

func (w *Wallet) String() string {
	return "TODO"
}
