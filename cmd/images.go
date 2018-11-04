package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"gopkg.in/abiosoft/ishell.v2"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

var errMissingImagePath = errors.New("missing image path")

func init() {
	register(&imagesCmd{})
}

type imagesCmd struct {
	Add addImagesCmd `command:"add"`
	//List   lsImagesCmd  `command:"ls"`
	//Get    getImagesCmd `command:"get"`
	//Delete delImagesCmd `command:"del"`
}

func (x *imagesCmd) Name() string {
	return "images"
}

func (x *imagesCmd) Short() string {
	return "Manage images"
}

func (x *imagesCmd) Long() string {
	return "Add, ls, get, and del images."
}

func (x *imagesCmd) Shell() *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
	}
	cmd.AddCmd((&addImagesCmd{}).Shell())
	//cmd.AddCmd((&lsImagesCmd{}).Shell())
	//cmd.AddCmd((&getImagesCmd{}).Shell())
	//cmd.AddCmd((&delImagesCmd{}).Shell())
	return cmd
}

type addImagesCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *addImagesCmd) Name() string {
	return "add"
}

func (x *addImagesCmd) Short() string {
	return "Add an image"
}

func (x *addImagesCmd) Long() string {
	return "Encodes, encrypts, and adds an image to the wallet account."
}

func (x *addImagesCmd) Execute(args []string) error {
	setApi(x.Client)
	return callAddImages(args, nil)
}

func (x *addImagesCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if err := callAddImages(c.Args, c); err != nil {
				c.Err(err)
			}
		},
	}
}

func callAddImages(args []string, ctx *ishell.Context) error {
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
	var list *[]core.AddDataResult
	res, err := executeJsonCmd(POST, "images", params{
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

////////////////////////////////////////////////////////////////////////

func sharePhoto(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo id"))
		return
	}
	if len(c.Args) == 1 {
		c.Err(errors.New("missing destination thread id"))
		return
	}
	id := c.Args[0]
	threadId := c.Args[1]

	c.Print("caption (optional): ")
	caption := c.ReadLine()

	// get the original block
	block, err := getPhotoBlockByDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	// lookup destination thread
	toThread := core.Node.Thread(threadId)
	if toThread == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", threadId)))
		return
	}

	// TODO: owner challenge
	// finally, add to destination
	if _, err := toThread.AddPhoto(id, caption, block.DataKey); err != nil {
		c.Err(err)
		return
	}
}

func listPhotos(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	threadId := c.Args[0]

	thrd := core.Node.Thread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.FileBlock)
	blocks := core.Node.Blocks("", -1, query)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no photos found in: %s", thrd.Id))
	} else {
		c.Println(fmt.Sprintf("%v photos:", len(blocks)))
	}

	magenta := color.New(color.FgHiMagenta).SprintFunc()
	for _, block := range blocks {
		c.Println(magenta(fmt.Sprintf("id: %s, block: %s", block.DataId, block.Id)))
	}
}

func getPhoto(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo id"))
		return
	}
	if len(c.Args) == 1 {
		c.Err(errors.New("missing out directory"))
		return
	}
	id := c.Args[0]

	// try to get path with home dir tilda
	dest, err := homedir.Expand(c.Args[1])
	if err != nil {
		dest = c.Args[1]
	}

	block, err := getPhotoBlockByDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	data, err := core.Node.BlockData(fmt.Sprintf("%s/photo", id), block)
	if err != nil {
		c.Err(err)
		return
	}

	path := filepath.Join(dest, id)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		c.Err(err)
		return
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	c.Println(blue("saved to " + path))
}

func getPhotoMetadata(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo id"))
		return
	}
	id := c.Args[0]

	block, err := getPhotoBlockByDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	jsonb, err := json.MarshalIndent(block.DataMetadata, "", "    ")
	if err != nil {
		c.Err(err)
		return
	}

	black := color.New(color.FgHiBlack).SprintFunc()
	c.Println(black(string(jsonb)))
}

func getPhotoKey(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo id"))
		return
	}
	id := c.Args[0]

	block, err := getPhotoBlockByDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	c.Println(blue(base58.FastBase58Encoding(block.DataKey)))
}

func addPhotoComment(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]
	c.Print("comment: ")
	body := c.ReadLine()

	block, err := core.Node.Block(id)
	if err != nil {
		c.Err(err)
		return
	}
	thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	if _, err := thrd.AddComment(block.Id, body); err != nil {
		c.Err(err)
		return
	}
}

func addPhotoLike(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]

	block, err := core.Node.Block(id)
	if err != nil {
		c.Err(err)
		return
	}
	thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	if _, err := thrd.AddLike(block.Id); err != nil {
		c.Err(err)
		return
	}
}

func listPhotoComments(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]

	block, err := core.Node.Block(id)
	if err != nil {
		c.Err(err)
		return
	}
	thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.CommentBlock)
	blocks := core.Node.Blocks("", -1, query)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no comments found on: %s", block.Id))
	} else {
		c.Println(fmt.Sprintf("%v comments:", len(blocks)))
	}

	cyan := color.New(color.FgHiCyan).SprintFunc()
	for _, b := range blocks {
		username := core.Node.ContactUsername(b.AuthorId)
		c.Println(cyan(fmt.Sprintf("%s: %s: %s", b.Id, username, b.DataCaption)))
	}
}

func listPhotoLikes(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]

	block, err := core.Node.Block(id)
	if err != nil {
		c.Err(err)
		return
	}
	thrd := core.Node.Thread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.LikeBlock)
	blocks := core.Node.Blocks("", -1, query)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no likes found on: %s", block.Id))
	} else {
		c.Println(fmt.Sprintf("%v likes:", len(blocks)))
	}

	cyan := color.New(color.FgHiCyan).SprintFunc()
	for _, b := range blocks {
		username := core.Node.ContactUsername(b.AuthorId)
		c.Println(cyan(fmt.Sprintf("%s: %s", b.Id, username)))
	}
}

func getPhotoBlockByDataId(dataId string) (*repo.Block, error) {
	block, err := core.Node.BlockByDataId(dataId)
	if err != nil {
		return nil, err
	}
	if block.Type != repo.FileBlock {
		return nil, errors.New("not a photo block, aborting")
	}
	return block, nil
}
