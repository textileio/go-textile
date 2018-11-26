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

	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingFilePath = errors.New("missing file path")
var errMissingFileBlockId = errors.New("missing file block id")
var errNothingToAdd = errors.New("nothing to add")

func init() {
	register(&addCmd{})
	register(&lsCmd{})
	register(&getCmd{})
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

type addCmd struct {
	Client  ClientOptions `group:"Client Options"`
	Thread  string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Caption string        `short:"c" long:"caption" description:"File(s) caption."`
	Group   bool          `short:"g" long:"group" description:"Group directory files."`
	Verbose bool          `short:"v" long:"verbose" description:"Prints files as they are processed."`
}

func (x *addCmd) Name() string {
	return "add"
}

func (x *addCmd) Short() string {
	return "Add file(s) to a thread"
}

func (x *addCmd) Long() string {
	return `
Adds a file or directory of files to a thread. Files not supported 
by the thread schema are ignored. Nested directories are included.
Use the --group option to add directory files as a single object.  
Omit the --thread option to use the default thread (if selected).
`
}

func (x *addCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread":  x.Thread,
		"caption": x.Caption,
		"group":   strconv.FormatBool(x.Group),
		"verbose": strconv.FormatBool(x.Verbose),
	}
	return callAdd(args, opts)
}

func (x *addCmd) Shell() *ishell.Cmd {
	return nil
}

func callAdd(args []string, opts map[string]string) error {
	if len(args) == 0 {
		return errMissingFilePath
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

	pth, err := homedir.Expand(args[0])
	if err != nil {
		pth = args[0]
	}

	fi, err := os.Stat(pth)
	if err != nil {
		return err
	}

	var pths []string
	var dirs []core.Directory
	var count int

	start := time.Now()

	if fi.IsDir() {
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
		output(msg, nil)

		batches := batchPaths(pths, batchSize)
		for i, batch := range batches {
			res, err := millBatch(batch, thrd.Schema, verbose)
			if err != nil {
				return err
			}

			output(fmt.Sprintf("Milled batch %d/%d", i+1, len(batches)), nil)

			if group {
				for _, dir := range res {
					if dir != nil {
						dirs = append(dirs, dir)
						count++
					}
				}
			} else {
				for _, dir := range res {
					if dir != nil {
						caption := strings.TrimSpace(fmt.Sprintf("%s (%d)", opts["caption"], count+1))
						block, err := add([]core.Directory{dir}, threadId, caption, verbose)
						if err != nil {
							return err
						}

						output(fmt.Sprintf("File %d target: %s", count+1, block.Target), nil)

						count++
					}
				}
			}
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
		output(fmt.Sprintf("File target: %s", block.Target), nil)

		count++
	}

	if group && len(dirs) > 0 {
		block, err := add(dirs, threadId, opts["caption"], verbose)
		if err != nil {
			return err
		}
		output(fmt.Sprintf("Group target: %s", block.Target), nil)
	}

	dur := time.Now().Sub(start)

	if count == 0 {
		return errNothingToAdd
	}

	msg := fmt.Sprintf("Added %d file", count)
	if count != 1 {
		msg += "s"
	}
	output(fmt.Sprintf("%s in %s", msg, dur.String()), nil)

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
		output(res, nil)
	}
	return block, nil
}

func mill(pth string, node *schema.Node, verbose bool) (core.Directory, error) {
	f, err := os.Open(pth)
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
		} else {
			r, ct, err := multipartReader(f)
			if err != nil {
				return nil, err
			}
			reader = r
			ctype = ct
		}

		res, file, err = handleStep(node.Mill, reader, mopts, ctype)
		if err != nil {
			return nil, err
		}

		if verbose {
			output(res, nil)
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
				output(res, nil)
			}

			dir[step.Name] = *file
		}
	} else {
		return nil, schema.ErrEmptySchema
	}

	return dir, nil
}

func millBatch(pths []string, node *schema.Node, verbose bool) ([]core.Directory, error) {
	tmp := make([]core.Directory, len(pths))

	wg := sync.WaitGroup{}
	for i, pth := range pths {
		wg.Add(1)

		go func(i int, p string) {
			dir, err := mill(p, node, verbose)
			if err != nil {
				output("mill error: "+err.Error(), nil)
			} else {
				tmp[i] = dir
			}
			wg.Done()
		}(i, pth)

	}
	wg.Wait()

	return tmp, nil
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

type lsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  string        `short:"l" long:"limit" description:"List page size." default:"5"`
}

func (x *lsCmd) Name() string {
	return "ls"
}

func (x *lsCmd) Short() string {
	return "Paginate thread files"
}

func (x *lsCmd) Long() string {
	return `
Paginates thread files.
Omit the --thread option to paginate all files.
Specify "default" to use the default thread (if selected).
`
}

func (x *lsCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"offset": x.Offset,
		"limit":  x.Limit,
	}
	return callLs(opts)
}

func (x *lsCmd) Shell() *ishell.Cmd {
	return nil
}

func callLs(opts map[string]string) error {
	var list []core.ThreadFilesInfo
	res, err := executeJsonCmd(GET, "files", params{opts: opts}, &list)
	if err != nil {
		return err
	}

	output(res, nil)

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

	return callLs(map[string]string{
		"thread": opts["thread"],
		"offset": list[len(list)-1].Block,
		"limit":  opts["limit"],
	})
}

type getCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getCmd) Name() string {
	return "get"
}

func (x *getCmd) Short() string {
	return "Get a thread file"
}

func (x *getCmd) Long() string {
	return `
Gets a thread file by specifying a Thread Block ID.
`
}

func (x *getCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGet(args)
}

func (x *getCmd) Shell() *ishell.Cmd {
	return nil
}

func callGet(args []string) error {
	if len(args) == 0 {
		return errMissingFileBlockId
	}

	var info core.ThreadFilesInfo
	res, err := executeJsonCmd(GET, "files/"+args[0], params{}, &info)
	if err != nil {
		return err
	}

	output(res, nil)
	return nil
}
