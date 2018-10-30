package cmd

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"gopkg.in/abiosoft/ishell.v2"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

func init() {
	register(&addCmd{})
}

type addCmd struct {
	Image  addImageCmd  `command:"image"`
	Thread addThreadCmd `command:"thread"`
}

func (x *addCmd) Name() string {
	return "add"
}

func (x *addCmd) Short() string {
	return "Add images and threads"
}

func (x *addCmd) Long() string {
	return "Add is a subcommand for adding images and threads to the wallet account."
}

func (x *addCmd) Shell() *ishell.Cmd {
	cmd := &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
	}
	cmd.AddCmd((&addImageCmd{}).Shell())
	return cmd
}

type addImageCmd struct{}

func (x *addImageCmd) Name() string {
	return "image"
}

func (x *addImageCmd) Short() string {
	return "Add an image"
}

func (x *addImageCmd) Long() string {
	return "Encodes, encrypts, and adds an image to the wallet account."
}

func (x *addImageCmd) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("missing image path")
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

	var added *struct {
		Items []core.AddDataResult `json:"items"`
	}
	if err := executeJsonCmd(POST, "add/"+x.Name(), params{
		args:    args,
		payload: &body,
		ctype:   writer.FormDataContentType(),
	}, &added); err != nil {
		return err
	}

	for _, item := range added.Items {
		fmt.Println("id:  " + item.Id)
		fmt.Println("key: " + item.Key)
	}
	return nil
}

func (x *addImageCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing image path"))
				return
			}

			path, err := homedir.Expand(c.Args[0])
			if err != nil {
				path = c.Args[0]
			}

			added, err := core.Node.AddImageByPath(path)
			if err != nil {
				c.Err(err)
				return
			}

			c.Println(Grey("id:  ") + Green(added.Id))
			c.Println(Grey("key: ") + Green(added.Key))
		},
	}
}

type addThreadCmd struct{}

func (x *addThreadCmd) Name() string {
	return "thread"
}

func (x *addThreadCmd) Short() string {
	return "Add a new thread"
}

func (x *addThreadCmd) Long() string {
	return "Adds a new thread for tracking a set of files between peers."
}

func (x *addThreadCmd) Execute(args []string) error {
	res, err := executeStringCmd(POST, "add/"+x.Name(), params{args: args})
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (x *addThreadCmd) Shell() *ishell.Cmd {
	return &ishell.Cmd{
		Name:     x.Name(),
		Help:     x.Short(),
		LongHelp: x.Long(),
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing thread name"))
				return
			}
			name := c.Args[0]

			sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
			if err != nil {
				c.Err(err)
				return
			}

			thrd, err := core.Node.AddThread(name, sk, true)
			if err != nil {
				c.Err(err)
				return
			}

			c.Println(Grey("id:  ") + Cyan(thrd.Id))
		},
	}
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
	_, toThread := core.Node.Thread(threadId)
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

	_, thrd := core.Node.Thread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.PhotoBlock)
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
	_, thrd := core.Node.Thread(block.ThreadId)
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
	_, thrd := core.Node.Thread(block.ThreadId)
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
	_, thrd := core.Node.Thread(block.ThreadId)
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
		c.Println(cyan(fmt.Sprintf("%s: %s: %s", b.Id, getUsername(b.AuthorId), b.DataCaption)))
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
	_, thrd := core.Node.Thread(block.ThreadId)
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
		c.Println(cyan(fmt.Sprintf("%s: %s", b.Id, getUsername(b.AuthorId))))
	}
}

func getPhotoBlockByDataId(dataId string) (*repo.Block, error) {
	block, err := core.Node.BlockByDataId(dataId)
	if err != nil {
		return nil, err
	}
	if block.Type != repo.PhotoBlock {
		return nil, errors.New("not a photo block, aborting")
	}
	return block, nil
}

func getUsername(peerId string) string {
	contact := core.Node.Contact(peerId)
	if contact != nil {
		return contact.Username
	}
	return peerId[:8]
}
