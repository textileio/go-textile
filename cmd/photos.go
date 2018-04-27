package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/segmentio/ksuid"
	"gopkg.in/abiosoft/ishell.v2"

	"github.com/textileio/textile-go/core"
)

func AddPhoto(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo path"))
		return
	}

	// try to get path with home dir tilda
	pp, err := homedir.Expand(c.Args[0])
	if err != nil {
		pp = c.Args[0]
	}

	// open the file
	f, err := os.Open(pp)
	if err != nil {
		c.Err(err)
		return
	}
	defer f.Close()

	// try to create a thumbnail version
	th, _, err := image.Decode(f)
	if err != nil {
		c.Err(err)
		return
	}
	th = imaging.Resize(th, 300, 0, imaging.Lanczos)
	thb := new(bytes.Buffer)
	if err = jpeg.Encode(thb, th, nil); err != nil {
		c.Err(err)
		return
	}
	tp := filepath.Join(core.Node.RepoPath, "tmp", ksuid.New().String()+".jpg")
	if err = ioutil.WriteFile(tp, thb.Bytes(), 0644); err != nil {
		c.Err(err)
		return
	}

	// parse album
	album := "default"
	if len(c.Args) > 1 {
		album = c.Args[1]
	}

	// do the add
	f.Seek(0, 0)
	mr, err := core.Node.AddPhoto(pp, tp, album)
	if err != nil {
		c.Err(err)
		return
	}

	// clean up
	if err = os.Remove(tp); err != nil {
		c.Err(err)
		return
	}
	if err = os.Remove(mr.PayloadPath); err != nil {
		c.Err(err)
		return
	}

	// show user root cid
	cyan := color.New(color.FgCyan).SprintFunc()
	c.Println(cyan("added " + mr.Boundary + " to thread " + album))
}

func SharePhoto(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo cid"))
		return
	}
	if len(c.Args) == 1 {
		c.Err(errors.New("missing destination thread name"))
		return
	}
	cid := c.Args[0]
	dest := c.Args[1]

	mr, err := core.Node.SharePhoto(cid, dest)
	if err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("shared " + cid + " to thread " + dest + " (new cid: " + mr.Boundary + ")"))
}

func ListPhotos(c *ishell.Context) {
	album := "default"
	if len(c.Args) > 0 {
		album = c.Args[0]
	}

	a := core.Node.Datastore.Albums().GetAlbumByName(album)
	if a == nil {
		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", album)))
		return
	}

	sets := core.Node.Datastore.Photos().GetPhotos("", -1, "album='"+a.Id+"'")
	if len(sets) == 0 {
		c.Println(fmt.Sprintf("no photos found in: %s", album))
	} else {
		c.Println(fmt.Sprintf("found %v photos in: %s", len(sets), album))
	}

	magenta := color.New(color.FgHiMagenta).SprintFunc()
	for _, s := range sets {
		c.Println(magenta(fmt.Sprintf("cid: %s, name: %s%s", s.Cid, s.MetaData.Name, s.MetaData.Ext)))
	}
}

func GetPhoto(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing photo cid"))
		return
	}
	if len(c.Args) == 1 {
		c.Err(errors.New("missing out directory"))
		return
	}
	hash := c.Args[0]

	// try to get path with home dir tilda
	dest, err := homedir.Expand(c.Args[1])
	if err != nil {
		dest = c.Args[1]
	}

	set, a, err := core.Node.LoadPhotoAndAlbum(hash)
	if err != nil {
		c.Err(err)
		return
	}

	pb, err := core.Node.GetFile(fmt.Sprintf("%s/photo", hash), a)
	if err != nil {
		c.Err(err)
		return
	}

	path := filepath.Join(dest, set.MetaData.Name+set.MetaData.Ext)
	if err := ioutil.WriteFile(path, pb, 0644); err != nil {
		c.Err(err)
		return
	}

	blue := color.New(color.FgHiBlue).SprintFunc()
	c.Println(blue("wrote " + hash + " to " + path))
}
