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
	"gopkg.in/abiosoft/ishell.v2"
)

var errMissingThreadId = errors.New("missing thread id")

func init() {
	register(&threadsCmd{})
}

type threadsCmd struct {
	Add        addThreadsCmd        `command:"add"`
	List       lsThreadsCmd         `command:"ls"`
	Get        getThreadsCmd        `command:"get"`
	GetDefault getDefaultThreadsCmd `command:"default"`
	Remove     rmThreadsCmd         `command:"rm"`
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

func (x *threadsCmd) Shell() *ishell.Cmd {
	return nil
}

type addThreadsCmd struct {
	Client ClientOptions  `group:"Client Options"`
	Key    string         `short:"k" long:"key" description:"A locally unique key used by an app to identify this thread on recovery."`
	Open   bool           `short:"o" long:"open" description:"Set the thread type to open (default private)."`
	Schema flags.Filename `short:"s" long:"schema" description:"Thread Schema filename. Superseded by built-in schema flags."`
	Photos bool           `long:"photos" description:"Use the built-in photo Schema."`
}

func (x *addThreadsCmd) Name() string {
	return "add"
}

func (x *addThreadsCmd) Short() string {
	return "Add a new thread"
}

func (x *addThreadsCmd) Long() string {
	return "Adds and joins a new thread."
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

func (x *addThreadsCmd) Shell() *ishell.Cmd {
	return nil
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

func (x *lsThreadsCmd) Name() string {
	return "ls"
}

func (x *lsThreadsCmd) Short() string {
	return "List threads"
}

func (x *lsThreadsCmd) Long() string {
	return "Lists info on all threads."
}

func (x *lsThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callLsThreads()
}

func (x *lsThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callLsThreads() error {
	var list *[]core.ThreadInfo
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

func (x *getThreadsCmd) Name() string {
	return "get"
}

func (x *getThreadsCmd) Short() string {
	return "Get a thread"
}

func (x *getThreadsCmd) Long() string {
	return "Gets and displays info about a thread."
}

func (x *getThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetThreads(args, nil)
}

func (x *getThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callGetThreads(args []string, ctx *ishell.Context) error {
	if len(args) == 0 {
		return errMissingThreadId
	}
	var info *core.ThreadInfo
	res, err := executeJsonCmd(GET, "threads/"+args[0], params{}, &info)
	if err != nil {
		return err
	}
	output(res, ctx)
	return nil
}

type getDefaultThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getDefaultThreadsCmd) Name() string {
	return "default"
}

func (x *getDefaultThreadsCmd) Short() string {
	return "Get default thread"
}

func (x *getDefaultThreadsCmd) Long() string {
	return "Gets and displays info about the default thread (if selected)."
}

func (x *getDefaultThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callGetDefaultThreads()
}

func (x *getDefaultThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callGetDefaultThreads() error {
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

func (x *rmThreadsCmd) Name() string {
	return "rm"
}

func (x *rmThreadsCmd) Short() string {
	return "Remove a thread"
}

func (x *rmThreadsCmd) Long() string {
	return "Leaves and removes a thread."
}

func (x *rmThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmThreads(args)
}

func (x *rmThreadsCmd) Shell() *ishell.Cmd {
	return nil
}

func callRmThreads(args []string) error {
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
