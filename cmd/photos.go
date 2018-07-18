package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/thread"
	"gopkg.in/abiosoft/ishell.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func AddPhoto(c *ishell.Context) {
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
	added, err := core.Node.Wallet.AddPhoto(path)
	if err != nil {
		c.Err(err)
		return
	}

	// clean up
	if err = os.Remove(added.PinRequest.PayloadPath); err != nil {
		c.Err(err)
		return
	}

	// add to thread
	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", threadId)))
		return
	}
	addr, err := thrd.AddPhoto(added.Id, caption, []byte(added.Key))
	if err != nil {
		c.Err(err)
		return
	}

	// show user root id
	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan(fmt.Sprintf("added photo %s to %s. added block %s.", added.Id, thrd.Id, addr.B58String())))
}

func SharePhoto(c *ishell.Context) {
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
	_, toThread := core.Node.Wallet.GetThread(threadId)
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
	addr, err := toThread.AddPhoto(id, caption, key)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green(fmt.Sprintf("shared photo %s to %s. added block %s.", id, toThread.Id, addr.B58String())))
}

func ListPhotos(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread id"))
		return
	}
	threadId := c.Args[0]

	_, thrd := core.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
		return
	}

	blocks := thrd.Blocks("", -1, repo.PhotoBlock)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no photos found in: %s", threadId))
	} else {
		c.Println(fmt.Sprintf("found %v photos in: %s", len(blocks), threadId))
	}

	magenta := color.New(color.FgHiMagenta).SprintFunc()
	for _, block := range blocks {
		c.Println(magenta(fmt.Sprintf("id: %s, block: %s", block.DataId, block.Id)))
	}
}

func GetPhoto(c *ishell.Context) {
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

func CatPhotoMetadata(c *ishell.Context) {
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

	data, err := thrd.GetBlockData(fmt.Sprintf("%s/meta", id), block)
	if err != nil {
		c.Err(err)
		return
	}
	var meta model.PhotoMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
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

func GetPhotoKey(c *ishell.Context) {
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

func getBlockAndThreadForDataId(dataId string) (*repo.Block, *thread.Thread, error) {
	block, err := core.Node.Wallet.GetBlockByDataId(dataId)
	if err != nil {
		return nil, nil, err
	}
	_, thrd := core.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		return nil, nil, errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}
	return block, thrd, nil
}
