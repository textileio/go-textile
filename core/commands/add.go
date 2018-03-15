package commands

import (
	"strconv"

	"github.com/textileio/textile-go/repo"

	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
	"gx/ipfs/QmabLouZTZwhfALuBcssPvkzhbYGMb4394huT7HY4LQ6d3/go-ipfs-cmds"
	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/coreunix"
	"fmt"
	"os"
)

var walletAddPhotoCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline: "Add a photo to the wallet on ipfs.",
		ShortDescription: `
Adds contents of a photo <path> to the wallet on ipfs.
`,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.FileArg("path", true, true, "The path to the photo to be added to wallet.").EnableStdin(),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) {
		repoDir, _ := req.Options[repoDirKwd].(string)

		r, err := fsrepo.Open(repoDir)
		if err != nil { // NB: repo is owned by the node
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		n, err := core.NewNode(req.Context, &core.BuildCfg{Repo: r})
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		// TODO: handle directories
		// this check doesn't seem to work as it will always say directory
		//if req.Files.IsDirectory() {
		//	res.SetError(errors.New("directories not yet supported"), cmdkit.ErrNormal)
		//}

		// just get the first file
		file, err := req.Files.NextFile()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		dir, err := repo.PinPhoto(file, file.FileName(), n)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		size, err := dir.Size()
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}
		obj := &coreunix.AddedObject{
			Hash: dir.Cid().Hash().B58String(),
			Name: file.FileName(),
			Size: strconv.FormatUint(size, 10),
		}

		fmt.Fprintf(os.Stdout, "added %s %s\n", obj.Hash, obj.Name)
	},
	Type: coreunix.AddedObject{},
}
