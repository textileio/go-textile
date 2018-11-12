package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema"
	"gopkg.in/abiosoft/ishell.v2"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var errMissingFilePath = errors.New("missing file path")

func init() {
	register(&addCmd{})
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

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}
	writer.Close()
	reader := bytes.NewReader(body.Bytes())

	// traverse the schema and collect generated files
	dir := make(map[string]*repo.File)
	f, err := millNode(reader, writer.FormDataContentType(), info.Schema, dir, func(out string) {
		output(out, nil)
	})
	if err != nil {
		return err
	}

	// the schemas could have generated a directory or a single file
	var api string
	var body2 interface{}
	if len(dir) != 0 {
		api = "files"
		body2 = &dir
	} else if f != nil {
		api = "file"
		body2 = &f
	} else {
		return nil
	}

	data, err := json.Marshal(body2)
	if err != nil {
		return err
	}

	var block *core.BlockInfo
	res, err := executeJsonCmd(POST, "threads/"+threadId+"/"+api, params{
		opts:    map[string]string{"caption": opts["caption"]},
		payload: bytes.NewReader(data),
		ctype:   "application/json",
	}, &block)
	if err != nil {
		return err
	}

	output(res, nil)
	return nil
}

func millNode(reader *bytes.Reader, ctype string, node *schema.Node, dir map[string]*repo.File, out func(string)) (*repo.File, error) {
	file := &repo.File{}

	if node.Mill != "" {
		res, err := executeJsonCmd(POST, "mills"+node.Mill, params{
			opts:    node.Opts,
			payload: reader,
			ctype:   ctype,
		}, &file)
		if err != nil {
			return nil, err
		}
		out(res)
	}

	for l, n := range node.Nodes {
		reader.Seek(0, 0)
		f, err := millNode(reader, ctype, n, dir, out)
		if err != nil {
			return nil, err
		}
		dir[l] = f
	}

	return file, nil
}
