package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/textileio/textile-go/repo/wallet"

	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit/files"
	"gx/ipfs/QmfAkMSt9Fwzk48QDJecPcwCUjnf2uG7MLnmCGTp4C6ouL/go-ipfs-cmds"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/commands"
	"gx/ipfs/QmatUACvrFK3xYg1nd2iLAKfz7Yy5YB56tnzBYHpqiUuhn/go-ipfs/core/coreunix"
)

// TODO: Add --remote flag to determine if the photo should also be added to the remote API
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
		n, err := commands.GetNode(env)
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
			return
		}

		outChan := make(chan interface{}, 8)

		// TODO: handle directories
		// this check doesn't seem to work as it will always say directory
		//if req.Files.IsDirectory() {
		//	res.SetError(errors.New("directories not yet supported"), cmdkit.ErrNormal)
		//}

		// Current hack was just to expect Full res folllowed by Thumb in list

		addAllAndPin := func(f files.File) error {
			// just get the first file
			file, err := f.NextFile()
			if err != nil {
				return err
			}
			// just get the first thumb
			thumb, err := f.NextFile()
			if err != nil {
				return err
			}

			dir, err := wallet.PinPhoto(file, file.FileName(), thumb, n, "")
			if err != nil {
				return err
			}
			size, err := dir.Size()
			if err != nil {
				return err
			}
			outChan <- &coreunix.AddedObject{
				Hash: dir.Cid().Hash().B58String(),
				Name: file.FileName(),
				Size: strconv.FormatUint(size, 10),
			}
			return nil
		}

		errCh := make(chan error)
		go func() {
			var err error
			defer func() { errCh <- err }()
			defer close(outChan)
			err = addAllAndPin(req.Files)
		}()

		defer res.Close()

		err = res.Emit(outChan)
		if err != nil {
			log.Error(err)
			return
		}
		err = <-errCh
		if err != nil {
			res.SetError(err, cmdkit.ErrNormal)
		}
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: func(req *cmds.Request, re cmds.ResponseEmitter) cmds.ResponseEmitter {
			reNext, res := cmds.NewChanResponsePair(req)
			outChan := make(chan interface{})

			progress := func(wait chan struct{}) {
				defer close(wait)

				lastHash := ""

			LOOP:
				for {
					select {
					case out, ok := <-outChan:
						if !ok {
							fmt.Fprintln(os.Stdout, lastHash)
							break LOOP
						}
						output := out.(*coreunix.AddedObject)
						if len(output.Hash) > 0 {
							lastHash = output.Hash
							fmt.Fprintf(os.Stdout, "added %s %s and thumb\n", output.Hash, output.Name)
							return
						} else {
							continue
						}
					case <-req.Context.Done():
						// don't set or print error here, that happens in the goroutine below
						return
					}
				}
			}

			go func() {
				// defer order important! First close outChan, then wait for output to finish, then close re
				defer re.Close()

				if e := res.Error(); e != nil {
					defer close(outChan)
					re.SetError(e.Message, e.Code)
					return
				}

				wait := make(chan struct{})
				go progress(wait)

				defer func() { <-wait }()
				defer close(outChan)

				for {
					v, err := res.Next()
					if !cmds.HandleError(err, res, re) {
						break
					}

					select {
					case outChan <- v:
					case <-req.Context.Done():
						re.SetError(req.Context.Err(), cmdkit.ErrNormal)
						return
					}
				}
			}()

			return reNext
		},
	},
	Type: coreunix.AddedObject{},
}
