package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema/textile"
	"github.com/textileio/go-textile/util"
)

var errMissingThreadId = errors.New("missing thread id")
var errMissingThreadName = errors.New("missing thread name")

func init() {
	register(&threadsCmd{})
}

type threadsCmd struct {
	Add        addThreadsCmd        `command:"add" description:"Add a new thread"`
	List       lsThreadsCmd         `command:"ls" description:"List threads"`
	Get        getThreadsCmd        `command:"get" description:"Get a thread"`
	GetDefault getDefaultThreadsCmd `command:"default" description:"Get default thread"`
	Peers      peersThreadsCmd      `command:"peers" description:"List thread peers"`
	Rename     renameThreadsCmd     `command:"rename" description:"Rename thread"`
	Remove     rmThreadsCmd         `command:"rm" description:"Remove a thread"`
	Snapshots  snapshotsThreadsCmd  `command:"snapshots" description:"Manage thread snapshots"`
}

func (x *threadsCmd) Name() string {
	return "threads"
}

func (x *threadsCmd) Short() string {
	return "Manage threads"
}

func (x *threadsCmd) Long() string {
	return `
Threads are distributed sets of encrypted files, often shared between peers,
governed by schemas.
Use this command to add, list, get, and remove threads. See below for
additional commands.

Thread type controls read (R), annotate (A), and write (W) access:

private   --> initiator: RAW, members:
read_only --> initiator: RAW, members: R
public    --> initiator: RAW, members: RA
open      --> initiator: RAW, members: RAW

Thread sharing style controls if (Y/N) a thread can be shared:

not_shared  --> initiator: N, members: N
invite_only --> initiator: Y, members: N
shared      --> initiator: Y, members: Y`
}

type addThreadsCmd struct {
	Client     ClientOptions  `group:"Client Options"`
	Key        string         `short:"k" long:"key" description:"A locally unique key used by an app to identify this thread on recovery."`
	Type       string         `short:"t" long:"type" description:"Set the thread type to one of 'private', 'read_only', 'public', or 'open'." default:"private"`
	Sharing    string         `short:"s" long:"sharing" description:"Set the thread sharing style to one of 'not_shared', 'invite_only', or 'shared'." default:"not_shared"`
	Member     []string       `short:"m" long:"member" description:"A contact address. When supplied, the thread will not allow additional peers, useful for 1-1 chat/file sharing. Can be used multiple times to include multiple contacts.'"`
	Schema     string         `long:"schema" description:"Thread schema ID. Supersedes schema filename."`
	SchemaFile flags.Filename `long:"schema-file" description:"Thread schema filename. Supersedes the built-in schema flags."`
	Blob       bool           `long:"blob" description:"Use the built-in blob schema for generic data."`
	CameraRoll bool           `long:"camera-roll" description:"Use the built-in camera roll schema."`
	Media      bool           `long:"media" description:"Use the built-in media schema."`
}

func (x *addThreadsCmd) Usage() string {
	return `

Adds and joins a new thread. See 'textile threads --help' for more about threads.`
}

func (x *addThreadsCmd) Execute(args []string) error {
	setApi(x.Client)

	var body []byte
	if x.Schema == "" {
		if x.SchemaFile != "" {
			path, err := homedir.Expand(string(x.SchemaFile))
			if err != nil {
				return err
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
		} else if x.Blob {
			body = []byte(textile.Blob)
		} else if x.CameraRoll {
			body = []byte(textile.CameraRoll)
		} else if x.Media {
			body = []byte(textile.Media)
		}
	}

	if body != nil {
		var schemaf pb.FileIndex
		if _, err := executeJsonPbCmd(POST, "mills/schema", params{
			payload: bytes.NewReader(body),
			ctype:   "application/json",
		}, &schemaf); err != nil {
			return err
		}
		x.Schema = schemaf.Hash
	}

	res, err := executeJsonCmd(POST, "threads", params{
		args: args,
		opts: map[string]string{
			"key":     x.Key,
			"type":    x.Type,
			"sharing": x.Sharing,
			"members": strings.Join(x.Member, ","),
			"schema":  x.Schema,
		},
	}, nil)
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

	res, err := executeJsonCmd(GET, "threads", params{}, nil)
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

	res, err := executeJsonCmd(GET, "threads/"+util.TrimQuotes(args[0]), params{}, nil)
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

	res, err := executeJsonCmd(GET, "threads/default", params{}, nil)
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
Omit the --thread option to use the default thread (if selected).`
}

func (x *peersThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	if x.Thread == "" {
		x.Thread = "default"
	}

	res, err := executeJsonCmd(GET, "threads/"+x.Thread+"/peers", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type renameThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Thread string        `short:"t" long:"thread" description:"Thread ID. Omit for default."`
}

func (x *renameThreadsCmd) Usage() string {
	return `

Renames a thread. Only the initiator can rename a thread.
Omit the --thread option to use the default thread (if selected).`
}

func (x *renameThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingThreadName
	}
	if x.Thread == "" {
		x.Thread = "default"
	}

	res, err := executeStringCmd(PUT, "threads/"+x.Thread+"/name", params{args: args})
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

	res, err := executeStringCmd(DEL, "threads/"+util.TrimQuotes(args[0]), params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type snapshotsThreadsCmd struct {
	Create createSnapshotsThreadsCmd `command:"create" description:"Create thread snapshots"`
	Search searchSnapshotsThreadsCmd `command:"search" description:"Search for thread snapshots"`
	Apply  applySnapshotsThreadsCmd  `command:"apply" description:"Apply a single thread snapshot"`
}

func (x *snapshotsThreadsCmd) Usage() string {
	return `

Use this command to create, search, and apply thread snapshots.`
}

type createSnapshotsThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *createSnapshotsThreadsCmd) Usage() string {
	return `

Snapshots all threads and pushes to registered cafes.`
}

func (x *createSnapshotsThreadsCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := callCreateSnapshotsThreads()
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func callCreateSnapshotsThreads() (string, error) {
	return executeStringCmd(POST, "snapshots", params{})
}

type searchSnapshotsThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 30s)." default:"2"`
}

func (x *searchSnapshotsThreadsCmd) Usage() string {
	return `

Searches the network for thread snapshots.`
}

func (x *searchSnapshotsThreadsCmd) Execute(args []string) error {
	setApi(x.Client)

	handleSearchStream("snapshots/search", params{
		opts: map[string]string{
			"wait": strconv.Itoa(x.Wait),
		},
	})
	return nil
}

type applySnapshotsThreadsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Wait   int           `long:"wait" description:"Stops searching after 'wait' seconds have elapsed (max 30s)." default:"2"`
}

func (x *applySnapshotsThreadsCmd) Usage() string {
	return `

Applies a single thread snapshot.`
}

func (x *applySnapshotsThreadsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingSnapshotId
	}
	id := args[0]

	results := handleSearchStream("snapshots/search", params{
		opts: map[string]string{
			"wait": strconv.Itoa(x.Wait),
		},
	})

	var result *pb.QueryResult
	for _, r := range results {
		if r.Id == id {
			result = &r
		}
	}

	if result == nil {
		output("Could not find snapshot with ID: " + id)
		return nil
	}

	if err := applySnapshot(result); err != nil {
		return err
	}

	return nil
}

func applySnapshot(result *pb.QueryResult) error {
	snap := new(pb.Thread)
	if err := ptypes.UnmarshalAny(result.Value, snap); err != nil {
		return err
	}
	data, err := pbMarshaler.MarshalToString(result.Value)
	if err != nil {
		return err
	}

	res, err := executeStringCmd(PUT, "threads/"+snap.Id, params{
		payload: strings.NewReader(data),
		ctype:   "application/json",
	})
	if err != nil {
		return err
	}
	if res == "" {
		output("applied " + result.Id)
	} else {
		output("error applying " + result.Id + ": " + res)
	}
	return nil
}
