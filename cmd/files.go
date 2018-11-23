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

	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingFilePath = errors.New("missing file path")
var errMissingFileBlockId = errors.New("missing file block id")

func init() {
	register(&addCmd{})
	register(&lsCmd{})
	register(&getCmd{})
}

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
}

func (x *addCmd) Name() string {
	return "add"
}

func (x *addCmd) Short() string {
	return "Add file(s) to a thread"
}

func (x *addCmd) Long() string {
	return `
Adds a file or directory to a thread.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *addCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread":  x.Thread,
		"caption": x.Caption,
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

	// first, ensure schema is present
	threadId := opts["thread"]
	if threadId == "" {
		threadId = "default"
	}
	var info *core.ThreadInfo
	if _, err := executeJsonCmd(GET, "threads/"+threadId, params{}, &info); err != nil {
		return err
	}

	if info.Schema == nil {
		return core.ErrThreadSchemaRequired
	}

	pth, err := homedir.Expand(args[0])
	if err != nil {
		pth = args[0]
	}

	f, err := os.Open(pth)
	if err != nil {
		return err
	}
	defer f.Close()

	var reader io.ReadSeeker
	var ctype string

	dir := make(core.Directory)

	// traverse the schema and collect generated files
	if info.Schema.Mill != "" {
		var res string
		file := &repo.File{}

		mopts := newMillOpts(info.Schema.Opts)
		mopts.setPlaintext(info.Schema.Plaintext)

		if info.Schema.Mill == "/json" {
			reader = f
			ctype = "application/json"
		} else {
			r, ct, err := multipartReader(f)
			if err != nil {
				return err
			}
			reader = r
			ctype = ct
		}

		res, file, err = handleStep(info.Schema.Mill, reader, mopts, ctype)
		if err != nil {
			return err
		}
		output(res, nil)

		dir[schema.SingleFileTag] = *file

	} else if len(info.Schema.Links) > 0 {

		// determine order
		steps, err := schema.Steps(info.Schema.Links)
		if err != nil {
			return err
		}

		// send each link
		for _, step := range steps {
			var res string
			file := &repo.File{}
			output("\""+step.Name+"\":", nil)

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
							return err
						}
						reader = r
						ctype = ct
					}
				}

				res, file, err = handleStep(step.Link.Mill, reader, mopts, ctype)
				if err != nil {
					return err
				}

			} else {
				if dir[step.Link.Use].Hash == "" {
					return errors.New(step.Link.Use + " not found")
				}
				mopts.setUse(dir[step.Link.Use].Hash)

				res, err = executeJsonCmd(POST, "mills"+step.Link.Mill, params{
					opts: mopts.val,
				}, &file)
				if err != nil {
					return err
				}
			}
			output(res, nil)

			dir[step.Name] = *file
		}
	} else {
		return schema.ErrEmptySchema
	}

	data, err := json.Marshal(&dir)
	if err != nil {
		return err
	}

	var block *core.BlockInfo
	res, err := executeJsonCmd(POST, "threads/"+threadId+"/files", params{
		opts:    map[string]string{"caption": opts["caption"]},
		payload: bytes.NewReader(data),
		ctype:   "application/json",
	}, &block)
	if err != nil {
		return err
	}

	output("\"block\":", nil)
	output(res, nil)

	return nil
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

type lsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for all."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  string        `short:"l" long:"limit" description:"List page size." default:"25"`
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
