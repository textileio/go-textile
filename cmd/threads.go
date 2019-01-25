package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"

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
	Peers      peersThreadsCmd      `command:"peers" description:"List thread peers"`
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

Thread type controls read (R), annotate (A), and write (W) access:

private   --> initiator: RAW, members:
read-only --> initiator: RAW, members: R
public    --> initiator: RAW, members: RA
open      --> initiator: RAW, members: RAW

Thread sharing style controls if (Y/N) a thread can be shared:

not-shared  --> initiator: N, members: N
invite-only --> initiator: Y, members: N
shared      --> initiator: Y, members: Y
`
}

type addThreadsCmd struct {
	Client     ClientOptions  `group:"Client Options"`
	Key        string         `short:"k" long:"key" description:"A locally unique key used by an app to identify this thread on recovery."`
	Type       string         `short:"t" long:"type" description:"Set the thread type to one of 'private', 'read_only', 'public', or 'open'." default:"private"`
	Sharing    string         `short:"s" long:"sharing" description:"Set the thread sharing style to one of 'not_shared', 'invite_only', or 'shared'." default:"not_shared"`
	Member     []string       `short:"m" long:"member" description:"A contact address. When supplied, the thread will not allow additional peers, useful for 1-1 chat/file sharing. Can be used multiple times to include multiple contacts.'"`
	Schema     flags.Filename `long:"schema" description:"Thread Schema filename. Supersedes the built-in schema flags."`
	Media      bool           `long:"media" description:"Use the built-in media Schema."`
	CameraRoll bool           `long:"camera-roll" description:"Use the built-in camera roll Schema."`
}

func (x *addThreadsCmd) Usage() string {
	return `

Adds and joins a new thread.`
}

func (x *addThreadsCmd) Execute(args []string) error {
	setApi(x.Client)

	var sch string
	switch x.Schema {
	case "":
		if x.Media {
			sch = "media"
			break
		}
		if x.CameraRoll {
			sch = "camera_roll"
			break
		}
	default:
		sch = string(x.Schema)
	}

	opts := map[string]string{
		"key":     x.Key,
		"type":    x.Type,
		"sharing": x.Sharing,
		"members": strings.Join(x.Member, ","),
		"schema":  sch,
	}
	return callAddThreads(args, opts)
}

func callAddThreads(args []string, opts map[string]string) error {
	var body []byte

	sch := opts["schema"]
	switch sch {
	case "media":
		body = []byte(textile.Media)
	case "camera_roll":
		body = []byte(textile.CameraRoll)
	default:
		if sch != "" {
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
		}
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
	output(res)
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
	output(res)
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
	output(res)
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
	output(res)
	return nil
}

type peersThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
}

func (x *peersThreadsCmd) Usage() string {
	return `

Lists all peers in a thread.
Omit the --thread option to use the default thread (if selected).
`
}

func (x *peersThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Thread == "" {
		x.Thread = "default"
	}
	var result []core.ContactInfo
	res, err := executeJsonCmd(GET, "threads/"+x.Thread+"/peers", params{}, &result)
	if err != nil {
		return err
	}
	output(res)
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
	output(res)
	return nil
}
