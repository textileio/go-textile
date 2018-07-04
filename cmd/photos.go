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
		c.Err(errors.New("missing thread name"))
		return
	}
	threadName := c.Args[1]

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
	if err = os.Remove(added.RemoteRequest.PayloadPath); err != nil {
		c.Err(err)
		return
	}

	// add to thread
	_, thrd := core.Node.Wallet.GetThreadByName(threadName)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread %s", threadName)))
		return
	}
	tadded, err := thrd.AddPhoto(added.Id, caption, added.Key)
	if err != nil {
		c.Err(err)
		return
	}

	// show user root id
	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan("added " + added.Id + " to thread " + thrd.Name + " with block " + tadded.Id))
}

func SharePhoto(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo id"))
		return
	}
	if len(c.Args) == 1 {
		c.Err(errors.New("missing destination thread name"))
		return
	}
	id := c.Args[0]
	threadName := c.Args[1]

	c.Print("caption (optional): ")
	caption := c.ReadLine()

	// get the original block
	block, fromThread, err := getBlockAndThreadForTarget(id)
	if err != nil {
		c.Err(err)
		return
	}

	// lookup destination thread
	_, toThread := core.Node.Wallet.GetThreadByName(threadName)
	if toThread == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread named %s", threadName)))
		return
	}

	// get the file key from the original block
	key, err := fromThread.Decrypt(block.TargetKey)
	if err != nil {
		c.Err(err)
		return
	}

	// TODO: owner challenge

	// finally, add to destination
	shared, err := toThread.AddPhoto(id, caption, key)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("shared " + id + " to thread " + toThread.Name + " (new id: " + shared.Id + ")"))
}

func ListPhotos(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing thread name"))
		return
	}
	threadName := c.Args[0]

	_, thrd := core.Node.Wallet.GetThreadByName(threadName)
	if thrd == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadName)))
		return
	}

	blocks := thrd.Blocks("", -1, repo.PhotoBlock)
	if len(blocks) == 0 {
		c.Println(fmt.Sprintf("no photos found in: %s", threadName))
	} else {
		c.Println(fmt.Sprintf("found %v photos in: %s", len(blocks), threadName))
	}

	magenta := color.New(color.FgHiMagenta).SprintFunc()
	for _, block := range blocks {
		c.Println(magenta(fmt.Sprintf("id: %s, block: %s", block.Target, block.Id)))
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

	block, thrd, err := getBlockAndThreadForTarget(id)
	if err != nil {
		c.Err(err)
		return
	}

	file, err := thrd.GetFileData(fmt.Sprintf("%s/photo", id), block)
	if err != nil {
		c.Err(err)
		return
	}

	path := filepath.Join(dest, id)
	if err := ioutil.WriteFile(path, file, 0644); err != nil {
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

	block, thrd, err := getBlockAndThreadForTarget(id)
	if err != nil {
		c.Err(err)
		return
	}

	file, err := thrd.GetFileData(fmt.Sprintf("%s/meta", id), block)
	if err != nil {
		c.Err(err)
		return
	}
	var meta model.PhotoMetadata
	if err := json.Unmarshal(file, &meta); err != nil {
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

	block, thrd, err := getBlockAndThreadForTarget(id)
	if err != nil {
		c.Err(err)
		return
	}

	key, err := thrd.GetFileKey(block)
	if err != nil {
		c.Err(err)
		return
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	c.Println(blue(key))
}

func getBlockAndThreadForTarget(id string) (*repo.Block, *thread.Thread, error) {
	block, err := core.Node.Wallet.GetBlockByTarget(id)
	if err != nil {
		return nil, nil, err
	}
	thrd := core.Node.Wallet.GetThread(block.ThreadPubKey)
	if thrd == nil {
		return nil, nil, errors.New(fmt.Sprintf("could not find thread %s", block.ThreadPubKey))
	}
	return block, thrd, nil
}
