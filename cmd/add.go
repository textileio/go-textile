package cmd

import (
	"bytes"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/repo"
	"gopkg.in/abiosoft/ishell.v2"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

//var errMissingThreadId = errors.New("missing thread id")

func init() {
	register(&addCmd{})
}

type addCmd struct {
	Files addFilesCmd `command:"files"`
}

func (x *addCmd) Name() string {
	return "add"
}

func (x *addCmd) Short() string {
	return "Add file(s), comments, and likes to a thread"
}

func (x *addCmd) Long() string {
	return `
Adds file(s), comments, and likes to a thread.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *addCmd) Shell() *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
	}
	cmd.AddCmd((&addFilesCmd{}).Shell())
	return cmd
}

type addFilesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID."`
}

func (x *addFilesCmd) Name() string {
	return "files"
}

func (x *addFilesCmd) Short() string {
	return "Add file(s) to a thread"
}

func (x *addFilesCmd) Long() string {
	return "Adds a file or a directory of files to a thread."
}

func (x *addFilesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
	}
	return callAddThreads(args, opts, nil)
}

func (x *addFilesCmd) Shell() *ishell.Cmd {
	return nil
}

func callAddFiles(args []string, opts map[string]string, ctx *ishell.Context) error {
	if len(args) == 0 {
		return errMissingImagePath
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
	var list *[]repo.File

	// add raw file
	res, err := executeJsonCmd(POST, "files", params{
		args:    args,
		payload: &body,
		ctype:   writer.FormDataContentType(),
	}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)

	// encode each
	res, err := executeJsonCmd(POST, "threads/"+opts["thread"], params{
		args:    args,
		payload: &body,
		ctype:   writer.FormDataContentType(),
	}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)

	return nil
}
