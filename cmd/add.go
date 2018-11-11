package cmd

import (
	"bytes"
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/repo"
	"gopkg.in/abiosoft/ishell.v2"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var errMissingImagePath = errors.New("missing image path")

func init() {
	register(&addCmd{})
}

type addCmd struct {
	Images addImagesCmd `command:"images"`
}

func (x *addCmd) Name() string {
	return "add"
}

func (x *addCmd) Short() string {
	return "Add files and data to a thread"
}

func (x *addCmd) Long() string {
	return `
Adds files and data to a thread.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *addCmd) Shell() *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
	}
	cmd.AddCmd((&addImagesCmd{}).Shell())
	return cmd
}

type addImagesCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
}

func (x *addImagesCmd) Name() string {
	return "images"
}

func (x *addImagesCmd) Short() string {
	return "Add images(s) to a thread"
}

func (x *addImagesCmd) Long() string {
	return "Adds an image or directory of images to a thread."
}

func (x *addImagesCmd) Execute(args []string) error {
	setApi(x.Client)
	opts := map[string]string{
		"thread": x.Thread,
	}
	return callAddImages(args, opts, nil)
}

func (x *addImagesCmd) Shell() *ishell.Cmd {
	return nil
}

func callAddImages(args []string, opts map[string]string, ctx *ishell.Context) error {
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
	reader := bytes.NewReader(body.Bytes())

	// gather up all the files
	var list *[]repo.File

	// get the thread schema
	// for each

	// encode each
	res, err := executeJsonCmd(POST, "images", params{
		args: args,
		opts: map[string]string{
			"width":   "1600",
			"quality": "75",
		},
		payload: reader,
		ctype:   writer.FormDataContentType(),
	}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)

	reader.Seek(0, 0)
	res, err = executeJsonCmd(POST, "images", params{
		args: args,
		opts: map[string]string{
			"width":   "800",
			"quality": "75",
		},
		payload: reader,
		ctype:   writer.FormDataContentType(),
	}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)

	reader.Seek(0, 0)
	res, err = executeJsonCmd(POST, "images", params{
		args: args,
		opts: map[string]string{
			"width":   "320",
			"quality": "75",
		},
		payload: reader,
		ctype:   writer.FormDataContentType(),
	}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)

	reader.Seek(0, 0)
	res, err = executeJsonCmd(POST, "images", params{
		args: args,
		opts: map[string]string{
			"width":   "100",
			"quality": "75",
		},
		payload: reader,
		ctype:   writer.FormDataContentType(),
	}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)

	// add raw file
	reader.Seek(0, 0)
	res, err = executeJsonCmd(POST, "blobs", params{
		args:    args,
		payload: reader,
		ctype:   writer.FormDataContentType(),
	}, &list)
	if err != nil {
		return err
	}
	output(res, ctx)

	return nil
}
