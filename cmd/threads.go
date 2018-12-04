package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema/textile"
)

var errMissingThreadId = errors.New("missing thread id")

func init() {
	register(&threadsCmd{})
}

type threadsCmd struct {
	Add        addThreadsCmd        `command:"add" description:"Add a new thread"`
	List       lsThreadsCmd         `command:"ls" description:"List threads"`
	Get        getThreadsCmd        `command:"get" description:"Get a thread"`
	GetDefault getDefaultThreadsCmd `command:"default" description:"Get default thread"`
	Remove     rmThreadsCmd         `command:"rm" description:"Remove a thread"`
}

func (x *threadsCmd) Name() string {
	return "threads"
}

func (x *threadsCmd) Short() string {
	return "Manage threads"
}

func (x *threadsCmd) Long() string {
	return `
Threads are distributed sets of encrypted files between peers,
governed by build-in or custom Schemas.
Use this command to add, list, get, join, invite, and remove threads.

Open threads are the most common thread type. Open threads allow 
any member to invite new members.

Private threads are primarily used internally for backup/recovery 
purposes and 1-to-1 communication channels.
`
}

type addThreadsCmd struct {
	Client ClientOptions  `group:"Client Options"`
	Key    string         `short:"k" long:"key" description:"A locally unique key used by an app to identify this thread on recovery."`
	Open   bool           `short:"o" long:"open" description:"Set the thread type to open (default private)."`
	Schema flags.Filename `short:"s" long:"schema" description:"Thread Schema filename. Superseded by built-in schema flags."`
	Photos bool           `long:"photos" description:"Use the built-in photo Schema."`
}

func (x *addThreadsCmd) Usage() string {
	return `

Adds and joins a new thread.`
}

func (x *addThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	ttype := "private"
	if x.Open {
		ttype = "open"
	}

	var sch string
	if x.Schema == "" && x.Photos {
		sch = "photos"
	} else {
		sch = string(x.Schema)
	}

	opts := map[string]string{
		"key":    x.Key,
		"type":   ttype,
		"schema": sch,
	}
	return callAddThreads(args, opts)
}

func callAddThreads(args []string, opts map[string]string) error {
	var body []byte

	sch := opts["schema"]
	if sch != "" && sch != "photos" {
		path, err := homedir.Expand(sch)
		if err != nil {
			path = sch
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		body, err = ioutil.ReadAll(file)
		if err != nil {
			return err
		}
	} else if sch == "photos" {
		body = []byte(textile.Photos)
	}

	if body != nil {
		var schemaf *repo.File
		if _, err := executeJsonCmd(POST, "mills/schema", params{
			payload: bytes.NewReader(body),
			ctype:   "application/json",
		}, &schemaf); err != nil {
			return err
		}

		opts["schema"] = schemaf.Hash
	}

	var info *core.ThreadInfo
	res, err := executeJsonCmd(POST, "threads", params{args: args, opts: opts}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type lsThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *lsThreadsCmd) Usage() string {
	return `

Lists info on all threads.`
}

func (x *lsThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	var list []core.ThreadInfo
	res, err := executeJsonCmd(GET, "threads", params{}, &list)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type getThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getThreadsCmd) Usage() string {
	return `

Gets and displays info about a thread.`
}

func (x *getThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingThreadId
	}
	var info *core.ThreadInfo
	res, err := executeJsonCmd(GET, "threads/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type getDefaultThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getDefaultThreadsCmd) Usage() string {
	return `

Gets and displays info about the default thread (if selected).`
}

func (x *getDefaultThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	var info *core.ThreadInfo
	res, err := executeJsonCmd(GET, "threads/default", params{}, &info)
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}

type rmThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmThreadsCmd) Usage() string {
	return `

Leaves and removes a thread.`
}

func (x *rmThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingThreadId
	}
	res, err := executeStringCmd(DEL, "threads/"+args[0], params{})
	if err != nil {
		return err
	}
	output(res, nil)
	return nil
}
