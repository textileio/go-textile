package commands

import (
	"fmt"
	"io"

	"github.com/textileio/textile-go/repo"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/commands/e"
	"gx/ipfs/QmabLouZTZwhfALuBcssPvkzhbYGMb4394huT7HY4LQ6d3/go-ipfs-cmds"
	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
)

var WalletCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Interact with the wallet.",
		ShortDescription: `
'textile wallet' is a tool to manipulate the data wallet
`,
	},

	Options: []cmdkit.Option{},
	Subcommands: map[string]*cmds.Command{
		"cat": walletCatCmd,
		"add": walletAddPhotoCmd,
	},
}

var walletCatCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Print decrypted wallet.",
		ShortDescription: `
'textile wallet cat' decrypts and prints the wallet
`,
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		err := cmds.EmitOnce(res, &repo.Wallet{})
		if err != nil {
			log.Error(err)
		}
	},
	Type: repo.Wallet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeEncoder(func(req *cmds.Request, w io.Writer, v interface{}) error {
			bs, ok := v.(*repo.Wallet)
			if !ok {
				return e.TypeErr(bs, v)
			}
			_, err := fmt.Fprintf(w, "%s", bs)
			return err
		}),
	},
}
