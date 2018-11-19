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

	path, err := homedir.Expand(args[0])
	if err != nil {
		path = args[0]
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, f); err != nil {
		return err
	}
	writer.Close()
	reader := bytes.NewReader(body.Bytes())

	// traverse the schema and collect generated files
	var payload interface{}
	if info.Schema.Mill != "" {
		file := &repo.File{}

		res, err := executeJsonCmd(POST, "mills"+info.Schema.Mill, params{
			opts:    info.Schema.Opts,
			payload: reader,
			ctype:   writer.FormDataContentType(),
		}, &file)
		if err != nil {
			return err
		}
		output(res, nil)
		payload = &file

	} else if len(info.Schema.Links) > 0 {
		dir := make(map[string]*repo.File)

		// determine order
		steps, err := schema.Steps(info.Schema.Links)
		if err != nil {
			return err
		}

		// send each link
		for _, step := range steps {
			file := &repo.File{}
			output("\""+step.Name+"\":", nil)

			if step.Link.Use == schema.FileTag {
				reader.Seek(0, 0)
				res, err := executeJsonCmd(POST, "mills"+step.Link.Mill, params{
					opts:    step.Link.Opts,
					payload: reader,
					ctype:   writer.FormDataContentType(),
				}, &file)
				if err != nil {
					return err
				}
				dir[step.Name] = file
				output(res, nil)

			} else {
				if dir[step.Link.Use] == nil {
					return errors.New(step.Link.Use + " not found")
				}
				if len(step.Link.Opts) == 0 {
					step.Link.Opts = make(map[string]string)
				}
				step.Link.Opts["use"] = dir[step.Link.Use].Hash
				res, err := executeJsonCmd(POST, "mills"+step.Link.Mill, params{
					opts: step.Link.Opts,
				}, &file)
				if err != nil {
					return err
				}
				dir[step.Name] = file
				output(res, nil)
			}
		}
		payload = &dir

	} else {
		return schema.ErrEmptySchema
	}

	data, err := json.Marshal(payload)
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

type lsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Offset string        `short:"o" long:"offset" description:"Offset ID to start listing from."`
	Limit  string        `short:"l" long:"limit" description:"List page size." default:"25"`
}

func (x *lsCmd) Name() string {
	return "ls"
}

func (x *lsCmd) Short() string {
	return "Paginate files in a thread"
}

func (x *lsCmd) Long() string {
	return `
Paginates files in a thread.
Omit the --thread option to use the default thread (if selected).
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
	threadId := opts["thread"]
	if threadId == "" {
		threadId = "default"
	}

	var list []core.ThreadFilesInfo
	res, err := executeJsonCmd(GET, "threads/"+threadId+"/files", params{opts: opts}, &list)
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
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
	Block  string        `short:"b" long:"block" description:"File Block ID."`
}

func (x *getCmd) Name() string {
	return "get"
}

func (x *getCmd) Short() string {
	return "Get a file in a thread"
}

func (x *getCmd) Long() string {
	return `
Gets a file in a thread.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *getCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
		"block":  x.Block,
	}
	return callGet(args, opts)
}

func (x *getCmd) Shell() *ishell.Cmd {
	return nil
}

func callGet(args []string, opts map[string]string) error {
	if len(args) == 0 {
		return errMissingFileBlockId
	}

	threadId := opts["thread"]
	if threadId == "" {
		threadId = "default"
	}

	var info core.ThreadFilesInfo
	res, err := executeJsonCmd(GET, "threads/"+threadId+"/files/"+args[0], params{opts: opts}, &info)
	if err != nil {
		return err
	}

	output(res, nil)
	return nil
}
