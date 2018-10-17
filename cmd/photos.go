package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/thread"
	"gopkg.in/abiosoft/ishell.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func addPhoto(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo path"))
		return
	}
	if len(c.Args) == 1 {
		c.Err(errors.New("missing thread id"))
		return
	}
	threadId := c.Args[1]

	// try to get path with home dir tilda
	path, err := homedir.Expand(c.Args[0])
	if err != nil {
		path = c.Args[0]
	}

	// open the file
	f, err := os.Open(path)
	if err != nil {
		c.Err(err)
		return
	}
	defer f.Close()

	c.Print("caption (optional): ")
	caption := c.ReadLine()

	// do the add
	f.Seek(0, 0)
	added, err := core.Node.AddPhoto(path)
	if err != nil {
		c.Err(err)
		return
	}

	// add to thread
	_, thrd := core.Node.GetThread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", threadId)))
		return
	}
	if _, err := thrd.AddPhoto(added.Id, caption, []byte(added.Key)); err != nil {
		c.Err(err)
		return
	}
}

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
	block, fromThread, err := getBlockAndThreadForDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	// lookup destination thread
	_, toThread := core.Node.GetThread(threadId)
	if toThread == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", threadId)))
		return
	}

	// get the file key from the original block
	key, err := fromThread.Decrypt(block.DataKeyCipher)
	if err != nil {
		c.Err(err)
		return
	}

	// TODO: owner challenge

	// finally, add to destination
	if _, err := toThread.AddPhoto(id, caption, key); err != nil {
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

	_, thrd := core.Node.GetThread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	btype := repo.PhotoBlock
	blocks := thrd.Blocks("", -1, &btype, nil)
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

	block, thrd, err := getBlockAndThreadForDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	data, err := thrd.GetBlockData(fmt.Sprintf("%s/photo", id), block)
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

	block, thrd, err := getBlockAndThreadForDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	meta, err := thrd.GetPhotoMetaData(id, block)
	if err != nil {
		c.Err(err)
		return
	}
	jsonb, err := json.MarshalIndent(meta, "", "    ")
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

	block, thrd, err := getBlockAndThreadForDataId(id)
	if err != nil {
		c.Err(err)
		return
	}

	key, err := thrd.GetBlockDataKey(block)
	if err != nil {
		c.Err(err)
		return
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	c.Println(blue(string(key)))
}

func addPhotoComment(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]
	c.Print("comment: ")
	body := c.ReadLine()

	block, err := core.Node.GetBlock(id)
	if err != nil {
		c.Err(err)
		return
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
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

	block, err := core.Node.GetBlock(id)
	if err != nil {
		c.Err(err)
		return
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
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

	block, err := core.Node.GetBlock(id)
	if err != nil {
		c.Err(err)
		return
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	btype := repo.CommentBlock
	blocks := thrd.Blocks("", -1, &btype, &block.Id)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no comments found on: %s", block.Id))
	} else {
		c.Println(fmt.Sprintf("%v comments:", len(blocks)))
	}

	cyan := color.New(color.FgHiCyan).SprintFunc()
	for _, b := range blocks {
		body := "nil"
		var authorUn string
		if b.DataCaptionCipher != nil {
			bodyb, err := thrd.Decrypt(b.DataCaptionCipher)
			if err != nil {
				c.Err(err)
				return
			}
			body = string(bodyb)
		}
		if b.AuthorUsernameCipher != nil {
			authorUnb, err := thrd.Decrypt(b.AuthorUsernameCipher)
			if err != nil {
				c.Err(err)
				return
			}
			authorUn = string(authorUnb)
		} else {
			authorUn = b.AuthorId[:8]
		}
		c.Println(cyan(fmt.Sprintf("%s: %s: %s", b.Id, authorUn, body)))
	}
}

func listPhotoLikes(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing block id"))
		return
	}
	id := c.Args[0]

	block, err := core.Node.GetBlock(id)
	if err != nil {
		c.Err(err)
		return
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
		return
	}

	btype := repo.LikeBlock
	blocks := thrd.Blocks("", -1, &btype, &block.Id)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no likes found on: %s", block.Id))
	} else {
		c.Println(fmt.Sprintf("%v likes:", len(blocks)))
	}

	cyan := color.New(color.FgHiCyan).SprintFunc()
	for _, b := range blocks {
		var authorUn string
		if b.AuthorUsernameCipher != nil {
			authorUnb, err := thrd.Decrypt(b.AuthorUsernameCipher)
			if err != nil {
				c.Err(err)
				return
			}
			authorUn = string(authorUnb)
		} else {
			authorUn = b.AuthorId[:8]
		}
		c.Println(cyan(fmt.Sprintf("%s: %s", b.Id, authorUn)))
	}
}

func getBlockAndThreadForDataId(dataId string) (*repo.Block, *thread.Thread, error) {
	block, err := core.Node.GetBlockByDataId(dataId)
	if err != nil {
		return nil, nil, err
	}
	if block.Type != repo.PhotoBlock {
		return nil, nil, errors.New("not a photo block, aborting")
	}
	_, thrd := core.Node.GetThread(block.ThreadId)
	if thrd == nil {
		return nil, nil, errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}
	return block, thrd, nil
}
