package cmd

import (
	"bytes"
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
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
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
		"thread": x.Thread,
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
	var files []repo.File
	if err := millNode(reader, writer.FormDataContentType(), info.Schema, files, func(out string) {
		output(out, nil)
	}); err != nil {
		return err
	}

	// post to thread

	return nil
}

func millNode(reader *bytes.Reader, ctype string, node *schema.Node, files []repo.File, out func(string)) error {
	if node.Mill != "" {
		file := repo.File{}
		res, err := executeJsonCmd(POST, "mills"+node.Mill, params{
			opts:    node.Opts,
			payload: reader,
			ctype:   ctype,
		}, &file)
		if err != nil {
			return err
		}
		files = append(files, file)
		out(res)
	}

	for _, n := range node.Nodes {
		reader.Seek(0, 0)
		if err := millNode(reader, ctype, n, files, out); err != nil {
			return err
		}
	}
	return nil
}
