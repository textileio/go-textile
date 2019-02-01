package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	iface "gx/ipfs/QmX9YciaxRii8TARoEbmavzaeTUAe7BozeAgydsThNcTpy/go-ipfs/core/coreapi/interface"

	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
)

var errMissingFilePath = errors.New("missing file path")
var errMissingFileId = errors.New("missing file block ID")
var errNothingToAdd = errors.New("nothing to add")
var errMissingTarget = errors.New("missing file(s) target")

func init() {
	register(&filesCmd{})
}

type filesCmd struct {
	Add    addFilesCmd `command:"add" description:"Add file(s) to a thread"`
	List   lsFilesCmd  `command:"ls" description:"Paginate thread files"`
	Get    getFilesCmd `command:"get" description:"Get a thread file by ID"`
	Ignore rmFilesCmd  `command:"ignore" description:"Ignore thread files"`
	Keys   keysCmd     `command:"keys" description:"Show file keys"`
}

func (x *filesCmd) Name() string {
	return "files"
}

func (x *filesCmd) Short() string {
	return "Manage thread files"
}

func (x *filesCmd) Long() string {
	return `
Files are added as blocks in a thread.
Use this command to add, list, get, and ignore files.
The 'key' command provides access to file encryption keys. 
`
}

const batchSize = 10

type millOpts struct {
	val map[string]string
}

func newMillOpts(ext map[string]string) millOpts {
	c := make(map[string]string)
	for k, v := range ext {
		c[k] = v
	}
	return millOpts{val: c}
}

func (m millOpts) setPlaintext(v bool) {
	m.val["plaintext"] = strconv.FormatBool(v)
}

func (m millOpts) setUse(v string) {
	m.val["use"] = v
}

type addFilesCmd struct {
	Client  ClientOptions `group:"Client Options"`
	Thread  string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Caption string        `short:"c" long:"caption" description:"File(s) caption."`
	Group   bool          `short:"g" long:"group" description:"Group directory files."`
	Verbose bool          `short:"v" long:"verbose" description:"Prints files as they are milled."`
}

func (x *addFilesCmd) Usage() string {
	return `

Adds a file or directory of files to a thread. Files not supported 
by the thread schema are ignored. Nested directories are included.
An existing file hash may also be used as input.
Use the --group option to add directory files as a single object.  
Omit the --thread option to use the default thread (if selected).
`
}

func (x *addFilesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread":  x.Thread,
		"caption": x.Caption,
		"group":   strconv.FormatBool(x.Group),
		"verbose": strconv.FormatBool(x.Verbose),
	}
	return callAddFiles(args, opts)
}

func callAddFiles(args []string, opts map[string]string) error {
	var pth string
	var fi os.FileInfo

	var err error
	fi, err = os.Stdin.Stat()
	if err != nil {
		return err
	}
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		if len(args) == 0 {
			return errMissingFilePath
		}

		// check if path references a cid
		ipth, err := iface.ParsePath(args[0])
		if err == nil {
			pth = ipth.String()
		} else {
			pth, err = homedir.Expand(args[0])
			if err != nil {
				pth = args[0]
			}

			fi, err = os.Stat(pth)
			if err != nil {
				return err
			}
		}
	}

	group := opts["group"] == "true"
	verbose := opts["verbose"] == "true"

	// fetch schema
	threadId := opts["thread"]
	if threadId == "" {
		threadId = "default"
	}
	var thrd *core.ThreadInfo
	if _, err := executeJsonCmd(GET, "threads/"+threadId, params{}, &thrd); err != nil {
		return err
	}

	if thrd.Schema == nil {
		return core.ErrThreadSchemaRequired
	}

	var pths []string
	var dirs []core.Directory
	var count int

	start := time.Now()

	if fi != nil && fi.IsDir() {
		err := filepath.Walk(pth, func(pth string, fi os.FileInfo, err error) error {
			if fi.IsDir() || fi.Name() == ".DS_Store" {
				return nil
			}
			pths = append(pths, pth)
			return nil
		})
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Found %d file", len(pths))
		if len(pths) != 1 {
			msg += "s"
		}
		output(msg)

		batches := batchPaths(pths, batchSize)
		for i, batch := range batches {

			ready := make(chan core.Directory, batchSize)
			go millBatch(batch, thrd.Schema, ready, verbose)

			var cerr error
		loop:
			for {
				select {
				case dir, ok := <-ready:
					if !ok {
						break loop
					}

					if !group {
						caption := strings.TrimSpace(fmt.Sprintf("%s (%d)", opts["caption"], count+1))
						block, err := add([]core.Directory{dir}, threadId, caption, verbose)
						if err != nil {
							cerr = err
							break loop
						}

						output(fmt.Sprintf("File %d target: %s", count+1, block.Target))
					} else {
						dirs = append(dirs, dir)
					}

					count++
				}
			}
			if cerr != nil {
				return cerr
			}

			output(fmt.Sprintf("Milled batch %d/%d", i+1, len(batches)))
		}

	} else {
		dir, err := mill(pth, thrd.Schema, verbose)
		if err != nil {
			return err
		}

		block, err := add([]core.Directory{dir}, threadId, opts["caption"], verbose)
		if err != nil {
			return err
		}
		output(fmt.Sprintf("File target: %s", block.Target))

		count++
	}

	if group && len(dirs) > 0 {
		block, err := add(dirs, threadId, opts["caption"], verbose)
		if err != nil {
			return err
		}
		output(fmt.Sprintf("Group target: %s", block.Target))
	}

	dur := time.Now().Sub(start)

	if count == 0 {
		return errNothingToAdd
	}

	msg := fmt.Sprintf("Added %d file", count)
	if count != 1 {
		msg += "s"
	}
	output(fmt.Sprintf("%s in %s", msg, dur.String()))

	return nil
}

func add(dirs []core.Directory, threadId string, caption string, verbose bool) (*core.BlockInfo, error) {
	data, err := json.Marshal(&dirs)
	if err != nil {
		return nil, err
	}

	var block *core.BlockInfo
	res, err := executeJsonCmd(POST, "threads/"+threadId+"/files", params{
		opts:    map[string]string{"caption": caption},
		payload: bytes.NewReader(data),
		ctype:   "application/json",
	}, &block)
	if err != nil {
		return nil, err
	}

	if verbose {
		output(res)
	}
	return block, nil
}

func mill(pth string, node *schema.Node, verbose bool) (core.Directory, error) {
	ref, err := iface.ParsePath(pth)
	if err == nil {
		parts := strings.Split(ref.String(), "/")
		pth = parts[len(parts)-1]
	}

	var f *os.File
	if ref == nil {
		if pth == "" {
			f = os.Stdin
		} else {
			var err error
			f, err = os.Open(pth)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				return nil, err
			}

			if fi.IsDir() {
				return nil, nil
			}
		}
	}

	var reader io.ReadSeeker
	var ctype string

	dir := make(core.Directory)

	// traverse the schema and collect generated files
	if node.Mill != "" {
		var res string
		file := &repo.File{}

		mopts := newMillOpts(node.Opts)
		mopts.setPlaintext(node.Plaintext)

		if node.Mill == "/json" {
			reader = f
			ctype = "application/json"
		} else if ref != nil {
			mopts.setUse(pth)
		} else {
			r, ct, err := multipartReader(f)
			if err != nil {
				return nil, err
			}
			reader = r
			ctype = ct
		}

		var err error
		res, file, err = handleStep(node.Mill, reader, mopts, ctype)
		if err != nil {
			return nil, err
		}

		if verbose {
			output(res)
		}

		dir[schema.SingleFileTag] = *file

	} else if len(node.Links) > 0 {

		// determine order
		steps, err := schema.Steps(node.Links)
		if err != nil {
			return nil, err
		}

		// send each link
		for _, step := range steps {
			var res string
			file := &repo.File{}

			mopts := newMillOpts(step.Link.Opts)
			mopts.setPlaintext(step.Link.Plaintext)

			if step.Link.Use == schema.FileTag {
				if reader != nil {
					reader.Seek(0, 0)
				} else {
					if step.Link.Mill == "/json" {
						reader = f
						ctype = "application/json"
					} else if ref != nil {
						mopts.setUse(pth)
					} else {
						r, ct, err := multipartReader(f)
						if err != nil {
							return nil, err
						}
						reader = r
						ctype = ct
					}
				}

				res, file, err = handleStep(step.Link.Mill, reader, mopts, ctype)
				if err != nil {
					return nil, err
				}

			} else {
				if dir[step.Link.Use].Hash == "" {
					return nil, errors.New(step.Link.Use + " not found")
				}
				mopts.setUse(dir[step.Link.Use].Hash)

				res, err = executeJsonCmd(POST, "mills"+step.Link.Mill, params{
					opts: mopts.val,
				}, &file)
				if err != nil {
					return nil, err
				}
			}

			if verbose {
				output(res)
			}

			dir[step.Name] = *file
		}
	} else {
		return nil, schema.ErrEmptySchema
	}

	return dir, nil
}

func millBatch(pths []string, node *schema.Node, ready chan core.Directory, verbose bool) {
	wg := sync.WaitGroup{}

	for _, pth := range pths {
		wg.Add(1)

		go func(p string) {
			dir, err := mill(p, node, verbose)
			if err != nil {
				output("mill error: " + err.Error())
			} else {
				ready <- dir
			}
			wg.Done()
		}(pth)

	}

	wg.Wait()
	close(ready)
}

func batchPaths(pths []string, size int) [][]string {
	var batches [][]string

	for i := 0; i < len(pths); i += size {
		end := i + size

		if end > len(pths) {
			end = len(pths)
		}

		batches = append(batches, pths[i:end])
	}

	return batches
}

func handleStep(mil string, reader io.Reader, opts millOpts, ctype string) (string, *repo.File, error) {
	var file *repo.File

	res, err := executeJsonCmd(POST, "mills"+mil, params{
		opts:    opts.val,
		payload: reader,
		ctype:   ctype,
	}, &file)
	if err != nil {
		return "", nil, err
	}

	return res, file, nil
}

func multipartReader(f *os.File) (io.ReadSeeker, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(f.Name()))
	if err != nil {
		return nil, "", err
	}
	if _, err = io.Copy(part, f); err != nil {
		return nil, "", err
	}
	writer.Close()
	return bytes.NewReader(body.Bytes()), writer.FormDataContentType(), nil
}

type lsFilesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  int           `short:"l" long:"limit" description:"List page size." default:"5"`
}

func (x *lsFilesCmd) Usage() string {
	return `

Paginates thread files.
Omit the --thread option to paginate all files.
Specify "default" to use the default thread (if selected).
`
}

func (x *lsFilesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"offset": x.Offset,
		"limit":  strconv.Itoa(x.Limit),
	}
	return callLsFiles(opts)
}

func callLsFiles(opts map[string]string) error {
	var list []core.ThreadFilesInfo
	res, err := executeJsonCmd(GET, "files", params{opts: opts}, &list)
	if err != nil {
		return err
	}

	output(res)

	limit, err := strconv.Atoi(opts["limit"])
	if err != nil {
		return err
	}
	if len(list) < limit {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return callLsFiles(map[string]string{
		"thread": opts["thread"],
		"offset": list[len(list)-1].Block,
		"limit":  opts["limit"],
	})
}

type getFilesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getFilesCmd) Usage() string {
	return `

Gets a thread file by block ID.
`
}

func (x *getFilesCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingFileId
	}

	var info core.ThreadFilesInfo
	res, err := executeJsonCmd(GET, "files/"+args[0], params{}, &info)
	if err != nil {
		return err
	}

	output(res)
	return nil
}

type rmFilesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmFilesCmd) Usage() string {
	return `

Ignores a thread file by its block ID.
This adds an "ignore" thread block targeted at the file.
Ignored blocks are by default not returned when listing. 
`
}

func (x *rmFilesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmBlocks(args)
}

type keysCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *keysCmd) Usage() string {
	return `

Shows file keys under the given target from an add.
`
}

func (x *keysCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingTarget
	}

	var keys core.Keys
	res, err := executeJsonCmd(GET, "keys/"+args[0], params{}, &keys)
	if err != nil {
		return err
	}

	output(res)
	return nil
}
