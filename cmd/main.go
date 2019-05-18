package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
	"io"
	"io/ioutil"
	"net/url"

	"os/signal"
	"path/filepath"

	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/fatih/color"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	_ "expvar"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	_ "net/http/pprof"
	"os"

	logging "github.com/ipfs/go-log"
)

type method string // e.g. http.MethodGet

type params struct {
	args    []string
	opts    map[string]string
	payload io.Reader
	ctype   string
}

var (
	// ================================

	// shared
	node *core.Textile
	log = logging.Logger("tex-main")

	// requests
	pbMarshaler = jsonpb.Marshaler{
		OrigName: true,
		Indent:   "    ",
	}
	pbUnmarshaler = jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	// locals
	errMissingSearchInfo = fmt.Errorf("missing search info")


	// ================================

	// colors
	Grey   = color.New(color.FgHiBlack).SprintFunc()
	Green  = color.New(color.FgHiGreen).SprintFunc()
	Cyan   = color.New(color.FgHiCyan).SprintFunc()
	Yellow = color.New(color.FgHiYellow).SprintFunc()

	// ================================

	// app
	appCmd     = kingpin.New("textile", "Textile is a set of tools and trust-less infrastructure for building censorship resistant and privacy preserving applications")
	apiAddr    = appCmd.Flag("api", "API Address to use.").Envar("API").Default("http://127.0.0.1:40600").String()
	apiVersion = appCmd.Flag("api-version", "API version to use.").Envar("API_VERSION").Default("v0").String()

	// ipfs
	ipfsServerMode = appCmd.Flag("server", "Apply IPFS server profile").Bool()
	ipfsSwarmPorts = appCmd.Flag("swarm-ports", "Set the swarm ports (TCP,WS). A random TCP port is chosen by default").String()

	// log
	logNoFiles = appCmd.Flag("no-log-files", "Write logs to stdout instead of rolling files").Bool()
	logDebug   = appCmd.Flag("debug", "Set the logging level to debug").Bool()

	// address
	apiBindAddr       = appCmd.Flag("api-bind-addr", "Set the local API address").Default("127.0.0.1:40600").String()
	cafeApiBindAddr   = appCmd.Flag("cafe-bind-addr", "Set the cafe REST API address").Default("0.0.0.0:40601").String()
	gatewayBindAddr   = appCmd.Flag("gateway-bind-addr", "Set the IPFS gateway address").Default("127.0.0.1:5050").String()
	profilingBindAddr = appCmd.Flag("profile-bind-addr", "Set the profiling address").Default("127.0.0.1:6060").String()

	// cafe
	cafeOpen        = appCmd.Flag("cafe-open", "Open the p2p Cafe Service for other peers").Bool()
	cafeURL         = appCmd.Flag("cafe-url", "Specify a custom URL of this cafe, e.g., https://mycafe.com").Envar("CAFE_HOST_URL").String()
	cafeNeighborURL = appCmd.Flag("cafe-neighbor-url", "Specify the URL of a secondary cafe. Must return cafe info, e.g., via a Gateway: https://my-gateway.yolo.com/cafe, or a Cafe API: https://my-cafe.yolo.com").Envar("CAFE_HOST_NEIGHBOR_URL").String()

	// @note removes the short names for the above, as they w ere conflicting with command ones

	// ================================

	// account
	accountCmd = appCmd.Command("account", "Manage a wallet account")

	// get
	accountGetCmd = accountCmd.Command("get", "Shows the local peer's account info as a contact")

	// seed
	accountSeedCmd = accountCmd.Command("seed", "Shows the local peer's account seed")

	// address
	accountAddressCmd = accountCmd.Command("address", "Shows the local peer's account address")

	// sync
	accountSyncCmd = accountCmd.Command("sync", "Syncs the local account peer with other peers found on the network")
	accountSyncWait = accountSyncCmd.Flag("wait", "Stops searching after 'wait' seconds have elapsed (max 30s)").Default("2").Int()

	// ================================
	// @todo the verbose docs are better suited for https://docs.textile.io

	// block
	blockCmd = appCmd.Command("block", "Blocks are the raw components in a thread. Think of them as an append-only log of thread updates where each update is hash-linked to its parent(s). New / recovering peers can sync history by simply traversing the hash tree.").Alias("blocks")

	// list
	blockListCmd = blockCmd.Command("list", "Paginates blocks in a thread").Alias("ls")
	blockListThreadID = blockListCmd.Flag("thread", "Thread ID").Default("default").Short('t').String()
	blockListOffset = blockListCmd.Flag("offset", "Offset ID to start listing from").Short('o').String()
	blockListLimit  = blockListCmd.Flag("limit", "List page size").Short('l').Default("5").Int()
	blockListDots = blockListCmd.Flag("dots", "Return GraphViz dots instead of JSON").Short('d').Bool()

	// meta
	blockMetaCmd = blockCmd.Command("meta", "Get the metadata for a block").Alias("get")
	blockMetaBlockID = blockMetaCmd.Arg("id", "Block ID").Required().String()

	// remove
	// @todo does this ignore, or delete?
	blockRemoveCmd = blockCmd.Command("remove", "Remove a block").Alias("rm")
	blockRemoveBlockID = blockRemoveCmd.Arg("id", "Block ID").Required().String()

	// files
	blockFileCmd = blockCmd.Command("file", "Get the files or a file of a block").Alias("files")
	blockFileBlockID = blockFileCmd.Arg("id", "Block ID").Required().String()
	blockFileIndex = blockFileCmd.Flag("index", "If provided, the index of a specific file to retrieve").Default("0").Int()
	blockFilePath = blockFileCmd.Flag("path", "If provided, the path of a specific file to retrieve").String()
	blockFileContent = blockFileCmd.Flag("content", "If provided alongside a path, the content of the specific file is retrieved").Bool()

	// ================================

	// cafe
	cafeCmd = appCmd.Command("cafe", "Commands to manage cafes").Alias("cafes")

	// add
	cafeAddCmd = cafeCmd.Command("add", `Registers with a cafe and saves an expiring service session token.
	An access token is required to register, and should be obtained separately from the target Cafe.`)
	cafeAddToken = cafeAddCmd.Flag("token", "An access token supplied by the Cafe.").Short('t').Required().String()

	// list
	cafeListCmd = cafeCmd.Command("list", "List info about all active cafe sessions").Alias("ls")

	// get
	cafeGetCmd = cafeCmd.Command("get", "Gets and displays info about a cafe session")
	cafeGetCafeID = cafeGetCmd.Arg("id", "Cafe ID").Required().String()

	// remove
	// @todo does this ignore, or delete?
	cafeRemoveCmd = cafeCmd.Command("remove", "Deregisters a cafe (content will expire based on the cafe's service rules)").Alias("rm")
	cafeRemoveCafeID = cafeRemoveCmd.Arg("id", "Cafe ID").Required().String()

	// messages
	cafeMessagesCmd = cafeCmd.Command("messages", "Check for messages at all cafes. New messages are downloaded and processed opportunistically.")

	// ================================

	// chat
	chatCmd = appCmd.Command("chat", `Starts an interactive chat session in a thread`)
	chatThreadID = chatCmd.Flag("thread", "Thread ID").Default("default").Short('t').String()

	// ================================

	// comment
	commentCmd = appCmd.Command("comment",  "Comments are added as blocks in a thread, which target another block, usually a file(s)").Alias("comments")

	// add
	commentAddCmd = commentCmd.Command("add", "Adds a comment to a thread block")
	commentAddBlockID = commentAddCmd.Flag("block", "Thread Block ID, usually a file(s) block").Short('b').Required().String()
	commentAddBody = commentAddCmd.Arg("body", "Comment body to add to the thread").Required().String()

	// list
	commentListCmd = commentCmd.Command("list", "Lists comments on a thread block").Alias("ls")
	commentListBlockID = commentListCmd.Arg("block", "Thread Block ID, usually a file(s) block").Required().String()

	// get
	commentGetCmd = commentCmd.Command("get", "Gets a thread comment by its comment block ID")
	commentGetCommentID = commentGetCmd.Arg("id", "Comment Block ID").Required().String()

	// ignore
	commentIgnoreCmd = commentCmd.Command("ignore", `Ignores a thread comment by its block ID.
This adds an "ignore" thread block targeted at the comment.
Ignored blocks are by default not returned when listing.`).Alias("remove").Alias("rm")
	commentIgnoreCommentID = commentIgnoreCmd.Arg("id", "Comment Block ID").Required().String()

	// ================================
	// @todo this documentation is better suited for docs.textile.io

	// config
	configCmd = appCmd.Command("config", `The config command controls configuration variables.
It works much like 'git config'. The configuration
values are stored in a config file inside your Textile
repository.

Getting config values will report the currently active
config settings. This may differ from the values specified
when setting values.

When changing values, valid JSON types must be used.
For example, a string should be escaped or wrapped in
single quotes (e.g., \"127.0.0.1:40600\") and arrays and
objects work fine (e.g. '{"API": "127.0.0.1:40600"}')
but should be wrapped in single quotes. Be sure to restart
the daemon for changes to take effect.

Examples:

Get the value of the 'Addresses.API' key:

  $ textile config Addresses.API
  $ textile config Addresses/API # Alternative syntax

Print the entire Textile config file to console:

  $ textile config

Set the value of the 'Addresses.API' key:

  $ textile config Addresses.API \"127.0.0.1:40600\"`)
	configName = configCmd.Arg("name", "The name of the configuration variable you wish to retrieve").Required().String()
	configValue = configCmd.Arg("value", "If provided, the value to set the configuration variable to").String()


	// ================================

	// contact
	contactCmd = appCmd.Command("contact", "Manage local contacts and find other contacts on the network").Alias("contacts")

	// add
	contactAddCmd =  contactCmd.Command("add", "Adds a contact by display name or account address to known contacts")
	contactAddName = contactAddCmd.Flag("name", "Add by display name").Short('n').String()
	contactAddAddress = contactAddCmd.Flag("address", "Add by account address").Short('a').String()
	contactAddWait = contactAddCmd.Flag("wait", "Stops searching after [wait] seconds have elapsed").Int()

	// ls
	contactListCmd = contactCmd.Command("list", "List known contacts").Alias("ls")

	// get
	contactGetCmd =  contactCmd.Command("get", "Gets a known contact")
	contactGetAddress = contactGetCmd.Arg("address", "Account Address").Required().String()

	// delete
	// @todo does this ignore, or delete?
	contactRemoveCmd =  contactCmd.Command("remove", "Removes a known contact").Alias("rn")
	contactRemoveAddress = contactRemoveCmd.Arg("address", "Account Address").Required().String()

	// search
	contactSearchCmd = contactCmd.Command("search", "Searches locally and on the network for contacts").Alias("find")
	contactSearchName = contactSearchCmd.Flag("name", "Search by display name").Short('n').String()
	contactSearchAddress = contactSearchCmd.Flag("address", "Search by account address").Short('a').String()
	contactSearchLocal = contactSearchCmd.Flag("only-local", "Only search local contacts").Bool()
	contactSearchRemote = contactSearchCmd.Flag("only-remote", "Only search remote contacts").Bool()
	contactSearchLimit = contactSearchCmd.Flag("limit", "Stops searching after [limit] results are found").Default("5").Int()
	contactSearchWait = contactSearchCmd.Flag("wait", "Stops searching after [wait] seconds have elapsed (max 30s)").Default("2").Int()


	// ================================

	// daemon
	daemonCmd      = appCmd.Command("daemon", "Start a node daemon session")
	daemonRepoPath = daemonCmd.Flag("repo-dir", "Specify a custom repository path").Short('r').String()
	daemonPinCode  = daemonCmd.Flag("pin-code", "Specify the pin code for datastore encryption (omit no pin code was used during init)").Short('p').String()
	daemonDocs     = daemonCmd.Flag("serve-docs", "Whether to serve the local REST API docs").Short('s').Bool()
	// @note use global debug flag, as otherwise conflict arises

	// ================================

	// docs
	// docsCmd = appCmd.Command("docs", "Prints markdown docs for the command-line client")
	// @todo add commands as alias for --help

	// ================================
	// @todo this documentation is better suited for docs.textile.io

	// feed
	feedCmd = appCmd.Command("feed", `Paginates post (join|leave|files|message) and annotation (comment|like) block types as a consumable feed.
The --mode option dictates how the feed is displayed:

-  "chrono": All feed block types are shown. Annotations always nest their target post, i.e., the post a comment is about.
-  "annotated": Annotations are nested under post targets, but are not shown in the top-level feed.
-  "stacks": Related blocks are chronologically grouped into "stacks". A new stack is started if an unrelated block
   breaks continuity. This mode is used by Textile Photos. Stacks may include:

*  The initial post with some nested annotations. Newer annotations may have already been listed.
*  One or more annotations about a post. The newest annotation assumes the "top" position in the stack. Additional
     annotations are nested under the target. Newer annotations may have already been listed in the case as well.`)
	feedThreadID = feedCmd.Flag("thread", "Thread ID, omit for all").Short('t').String()
	feedOffset = feedCmd.Flag("offset", "Offset ID to start listening from").Short('o').String()
	feedLimit = feedCmd.Flag("limit", "List page size").Short('l').Default("3").Int()
	feedMode = feedCmd.Flag("mode", "Feed mode, one of: chrono, annotated, stacks").Short('m').Default("chrono").String()

	// ================================

	// file
	fileCmd = appCmd.Command("file", "Interact with Textile Files").Alias("files")

	// list
	fileListCmd = fileCmd.Command("list", `Paginates thread files.`).Alias("ls")
	fileListThreadID = fileListCmd.Flag("thread", "Thread ID").Default("default").Short('t').String()
	fileListOffset = fileListCmd.Flag("offset", "Offset ID to start listing from").Short('o').String()
	fileListLimit  = fileListCmd.Flag("limit", "List page size").Short('l').Default("5").Int()

	// keys
	fileKeysCmd    = fileCmd.Command("keys", "Shows file keys under the given target from an add.").Alias("key")
	fileKeysFileID = fileKeysCmd.Arg("id", "File ID").Required().String()

	// add
	fileAddCmd = fileCmd.Command("add", `Adds a file, directory, or hash to a thread.
Files not supported by the thread schema are ignored.
Nested directories are included.`)
	fileAddPath    = fileAddCmd.Arg("path", "The path to the file or directory to add, can also be an existing hash").Required().String()
	fileAddThreadID  = fileAddCmd.Flag("thread", "Thread ID").Default("default").Short('t').String()
	fileAddCaption = fileAddCmd.Flag("caption", "File(s) caption.").Short('c').String()
	fileAddGroup   = fileAddCmd.Flag("group", "If provided, group a directory's files together into a single object.").Short('g').Bool()
	fileAddVerbose = fileAddCmd.Flag("verbose", "Prints files as they are milled").Short('v').Bool()

	// ignore
	// @todo is this a block id or a file id?
	fileIgnoreCmd = fileCmd.Command("ignore", `Ignores a thread file by its block ID.
This adds an "ignore" thread block targeted at the file.
Ignored blocks are by default not returned when listing.`).Alias("remove").Alias("rm")
	fileIgnoreFileID = fileIgnoreCmd.Arg("id", "File ID").Required().String()

	// get
	fileGetCmd = fileCmd.Command("get", "Get the metadata or content of a specific file")
	fileGetFileID = fileGetCmd.Arg("id", "File ID").Required().String()
	fileGetContent = fileGetCmd.Flag("content", "If provided, the decrypted content of the file is retrieved").Bool()

	// ================================

	// init
	initCmd         = appCmd.Command("init", "Initialize the node repository and exit")
	initAccountSeed = initCmd.Flag("seed", "Account seed (run 'wallet' command to generate new seeds)").Short('s').Required().String()
	initPinCode     = initCmd.Flag("pin-code", "Specify a pin code for datastore encryption").Short('p').String()
	initRepoPath    = initCmd.Flag("repo-dir", "Specify a custom repository path").Short('r').String()
	// add address flags
	// add cafe flags
	// add log flags

	// ================================

	// invite
	inviteCmd  = appCmd.Command("invite", `Invites allow other users to join threads. There are two types of
invites, direct account-to-account and external:

- Account-to-account invites are encrypted with the invitee's account address (public key).
- External invites are encrypted with a single-use key and are useful for onboarding new users.`).Alias("invites")

	// create
	inviteCreateCmd = inviteCmd.Command("create", `Creates a direct account-to-account or external invite to a thread.`)
	inviteCreateThreadID = inviteCreateCmd.Flag("thread", "Thread ID").Default("default").Short('t').String()
	inviteCreateAddress = inviteCreateCmd.Flag("address", "Account Address, omit to create an external invite").Short('a').String()
	inviteCreateWait = inviteCreateCmd.Flag("wait", "Stops searching after [wait] seconds have elapsed (max 30s)").Default("2").Int()

	// list
	inviteListCmd = inviteCmd.Command("list", "Lists all pending thread invites").Alias("ls")

	// accept
	inviteAcceptCmd = inviteCmd.Command("accept", `Accepts a direct account-to-account or external invite to a thread`)
	inviteAcceptKey = inviteAcceptCmd.Flag("key", "Key for an external invite").Short('k').String()
	inviteAcceptID = inviteAcceptCmd.Arg("id", "Invite ID that you have received").Required().String()

	// ignore
	inviteIgnoreCmd = inviteCmd.Command("ignore", `Ignores a direct account-to-account invite to a thread`).Alias("remove").Alias("rm")
	inviteIgnoreID = inviteIgnoreCmd.Arg("id", "Invite ID that you wish to ignore").Required().String()


	// ================================

	// ipfs
	ipfsCmd = appCmd.Command("ipfs", "Provides access to some IPFS commands")

	// id
	ipfsIdCmd = ipfsCmd.Command("id", "Shows the local node's IPFS peer ID")

	// swarm
	ipfsSwarmCmd = ipfsCmd.Command("swarm", "Provides access to a limited set of IPFS swarm commands")

	// swarm connect
	ipfsSwarmConnectCmd = ipfsSwarmCmd.Command("connect", `Opens a new direct connection to a peer address`)
	ipfsSwarmConnectAddress = ipfsSwarmConnectCmd.Arg("address", `An IPFS multiaddr, such as: /ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ`).String()

	// swarm peers
	ipfsSwarmPeersCmd = ipfsSwarmCmd.Command("peers", "Lists the set of peers this node is connected to")
	ipfsSwarmPeersVerbose = ipfsSwarmPeersCmd.Flag("verbose", "Display all extra information").Short('v').Bool()
	ipfsSwarmPeersStreams  = ipfsSwarmPeersCmd.Flag("streams", "Also list information about open streams for search peer").Short('s').Bool()
	ipfsSwarmPeersLatency = ipfsSwarmPeersCmd.Flag("latency", "Also list information about the latency to each peer").Short('l').Bool()
	ipfsSwarmPeersDirection = ipfsSwarmPeersCmd.Flag("direction", "Also list information about the direction of connection").Short('d').Bool()

	// cat
	ipfsCatCmd = ipfsCmd.Command("cat", "Displays the data behind an IPFS CID (hash)")
	ipfsCatKey =  ipfsCatCmd.Flag("key", "Encryption key").Short('k').String()
	ipfsCatHash = ipfsCatCmd.Arg("hash", "IPFS CID").Required().String()


	// ================================
	// note, so this was quite inconsistent before, sometimes an arg, sometimes a flag
	// also a few typos in the file
	// also, we need to get a consistent ignore/remove/rm naming convention
	// also, why sometimes a thread block, why sometimes a file block, we need to be clear on this
	// as right now, it says thread block, but then usually a file's block - it needs to be one or the other!

	// like
	likeCmd = appCmd.Command("like", `Likes are added as blocks in a thread, which target another block, usually a file(s).
Use this command to add, list, get, and ignore likes.`).Alias("likes")

	// add
	likeAddCmd = likeCmd.Command("add", "Adds a like to a thread block")
	likeAddBlockID = likeAddCmd.Arg("block", "Thread Block ID, usually a file's block").Required().String()

	// list
	likeListCmd = likeCmd.Command("list", "List likes on a thread block").Alias("ls")
	likeListBlockID = likeListCmd.Arg("block", "Thread Block ID, usually a file's block").Required().String()

	// get
	likeGetCmd = likeCmd.Command("get", "Get a like by its Block ID")
	likeGetLikeID = likeGetCmd.Arg("id", "The like's Block ID").Required().String()

	// ignore
	likeIgnoreCmd = likeCmd.Command("ignore", "Ignore a like by its Block ID").Alias("remove").Alias("rm")
	likeIgnoreLikeID = likeIgnoreCmd.Arg("id", "The like's Block ID").Required().String()


	// ================================

	// log
	logCmd = appCmd.Command("log", `List or change the verbosity of one or all subsystems log output. Textile logs piggyback on the IPFS event logs.`).Alias("logs")
	logSubsystem = logCmd.Flag("subsystem", "The subsystem logging identifier, omit for all").Short('s').String()
	logLevel = logCmd.Flag("level", "One of: debug, info, warning, error, critical. Omit to get current level.").Short('l').String()
	logTexOnly = logCmd.Flag("tex-only", "Whether to list/change only Textile subsystems, or all available subsystems").Short('t').Bool()
	// ^ @todo why is this called tex, instead of textile-only

	// ================================
	// @todo, is really describing what a message is necessary in the CLI docs
	// that is what docs.textile.io is for

	// message
	messageCmd = appCmd.Command("message", "Messages are added as blocks in a thread").Alias("messages")

	// add
	messageAddCmd = messageCmd.Command("add", "Adds a message to a thread")
	messageAddThreadID = messageAddCmd.Flag("thread", "Thread ID").Default("default").String()
	messageAddBody = messageAddCmd.Arg("body", "The message to add the thread").String()


	// list
	messageListCmd = messageCmd.Command("list", "Paginates thread messages").Alias("ls")
	messageListThreadID = messageListCmd.Flag("thread", "Thread ID, omit to paginate all messages").Short('t').String()
	messageListOffset = messageListCmd.Flag("offset", "Offset ID to start the listing from").Short('o').String()
	messageListLimit = messageListCmd.Flag("limit", "List page size").Default("10").Short('l').Int()

	// get
	messageGetCmd = messageCmd.Command("get", "Gets a thread message by its Block ID")
	messageGetMessageID = messageGetCmd.Arg("id", "The Block ID for the Message").String()

	// ignore
	messageIgnoreCmd  = messageCmd.Command("ignore", "Ignores a thread message by its Block ID").Alias("remove").Alias("rm")
	messageIgnoreMessageID = messageIgnoreCmd.Arg("id", "The Block ID for the Message").String()


	// ================================

	// migrate
	migrateCmd      = appCmd.Command("migrate", "Migrate the node repository and exit")
	migrateRepoPath = migrateCmd.Flag("repo-dir", "Specify a custom repository path").Short('r').String()
	migratePinCode  = migrateCmd.Flag("pin-code", "Specify the pin code for datastore encryption (omit of none was used during init)").Short('p').String()


	// ================================

	// notification
	notificationCmd = appCmd.Command("notification", "Manage notifications that have been generated by thread and account activity").Alias("notifications")

	// list
	notificationListCmd = notificationCmd.Command("list", "Lists notifications").Alias("ls")
	// @todo only unread?

	// read
	notificationReadCmd = notificationCmd.Command("read", "Marks a notification as read")
	notificationReadID = notificationReadCmd.Arg("id", "Notification ID, set to [all] to mark all notifications as read").Required().String()

	// ================================

	// ping
	pingCmd = appCmd.Command("ping", "Pings another peer on the network, returning [online] or [offline]")
	pingAddress = pingCmd.Arg("address", "The address of the other peer on the network").Required().String()

	// ================================

	// profile
	profileCmd = appCmd.Command("profile", `Use this command to view and update the peer profile. A Textile account will
show a profile for each of its peers, e.g., mobile, desktop, etc.`)

	// get
	profileGetCmd = profileCmd.Command("get", "Gets the local peer profile")

	// set
	profileSetCmd = profileCmd.Command("set", "Sets the profile name and avatar of the peer")
	profileSetName = profileSetCmd.Flag("name", "Set the peer's display name").Short('n').String()
	profileSetAvatar = profileSetCmd.Flag("avatar", "Set the peer's avatar from an image path (JPEG, PNG, or GIF)").Short('a').String()

	// set name
	// set avatar
	// @todo surely name and avatar should become flags rather than commands
	// @todo implemented accordingly

	// ================================

	// subscribe
	subscribeCmd = appCmd.Command("subscribe", "Subscribes to updates in a thread or all threads. An update is generated when a new block is added to a thread.").Alias("sub")
	subscribeThreadID = subscribeCmd.Flag("thread", "Thread ID, omit for all").Short('t').String()
	subscribeType = subscribeCmd.Flag("type", "Only be alerted to specific type of updates, possible values: merge, ignore, flag, join, announce, leave, text, files comment, like. Can be used multiple times, e.g., --type files --type comment").Short('k').Strings()

	// ================================

	// summary
	summaryCmd = appCmd.Command("summary", "Get a summary of the local node's data")

	// ================================
	// @todo this documentation should be moved to docs.textile.io

	// thread
	threadCmd = appCmd.Command("thread", `Threads are distributed sets of encrypted files, often shared between peers, governed by schemas.
Use this command to add, list, get, and remove threads. See below for additional commands.

Control over thread access and sharing is handled by a combination of the --type and --sharing flags.
An immutable member address "whitelist" gives the initiator fine-grained control.
The table below outlines access patterns for the thread initiator and the whitelist members.
An empty whitelist is taken to be "everyone", which is the default.

Thread type controls read (R), annotate (A), and write (W) access:

private   --> initiator: RAW, whitelist:
read_only --> initiator: RAW, whitelist: R
public    --> initiator: RAW, whitelist: RA
open      --> initiator: RAW, whitelist: RAW

Thread sharing style controls if (Y/N) a thread can be shared:

not_shared  --> initiator: N, whitelist: N
invite_only --> initiator: Y, whitelist: N
shared      --> initiator: Y, whitelist: Y`).Alias("threads")

	// add
	threadAddCmd = threadCmd.Command("add", "Adds and joins a new thread")
	threadAddKey = threadAddCmd.Flag("key", "A locally unique key used by an app to identify this thread on recovery").Short('k').String()
	threadAddType = threadAddCmd.Flag("type", "Set the thread type to one of: private, read_only, public, open").Short('t').Default("private").String()
	threadAddSharing = threadAddCmd.Flag("sharing", "Set the thread sharing style to one of: not_shared, invite_only, shared").Short('s').Default("not_shared").String()
	threadAddWhitelist = threadAddCmd.Flag("whitelist", "A contact address. When supplied, the thread will not allow additional peers, useful for 1-1 chat/file sharing. Can be used multiple times to include multiple contacts").Short('w').Strings()
	threadAddSchema = threadAddCmd.Flag("schema", "Thread schema ID. Supersedes schema filename").String()
	threadAddSchemaFile = threadAddCmd.Flag("schema-file", "Thread schema filename, supersedes the built-in schema flags").String() // @todo could swap this to file perhaps
	threadAddBlob = threadAddCmd.Flag("blob", "Use the built-in blob schema for generic data").Bool()
	threadAddCameraRoll = threadAddCmd.Flag("camera-roll", "Use the built-in camera roll schema").Bool()
	threadAddMedia = threadAddCmd.Flag("media", "Use the built-in media schema").Bool()

	// list
	threadListCmd = threadCmd.Command("list", "Lists info on all threads").Alias("ls")

	// get
	threadGetCmd = threadCmd.Command("get", "Gets and displays info about a thread")
	threadGetThreadID = threadGetCmd.Arg("id", "Thread ID").Required().String()
	// @todo remove all trimQuotes in the cmd/* folder

	// default
	threadDefaultCmd = threadCmd.Command("default", "Gets and displays info about the default thread (if selected")

	// peer
	threadPeerCmd = threadCmd.Command("peer", "Lists all peers in a thread").Alias("peers")
	threadPeerThreadID = threadPeerCmd.Flag("thread", "Thread ID").Default("default").Short('t').String()

	// rename
	threadRenameCmd = threadCmd.Command("rename", "Renames a thread. Only the initiator can rename the thread.").Alias("mv")
	threadRenameThreadID = threadRenameCmd.Flag("thread", "Thread ID").Default("default").Short('t').String()
	threadRenameName = threadRenameCmd.Arg("name", "The name to rename the thread to").Required().String()

	// remove
	threadRemoveCmd = threadCmd.Command("remove", "Leaves and removes a thread").Alias("rm")
	threadRemoveThreadID = threadRemoveCmd.Arg("id", "Thread ID").Required().String()

	// snapshot
	threadSnapshotCmd = threadCmd.Command("snapshot", "Manage thread snapshots").Alias("snapshots")

	// snapshot create
	threadSnapshotCreateCmd = threadSnapshotCmd.Command("create", "Snapshots all threads and pushes to registered cafes").Alias("make")

	// snapshot search
	threadSnapshotSearchCmd = threadSnapshotCmd.Command("search", "Searches the network for thread snapshots").Alias("find")
	threadSnapshotSearchWait = threadSnapshotSearchCmd.Flag("wait", "Stops searching after [wait] seconds have elapse (max 30s)").Short('w').Default("2").Int()

	// snapshot apply
	threadSnapshotApplyCmd = threadSnapshotCmd.Command("apply", "Applies a single thread snapshot")
	threadSnapshotApplyWait = threadSnapshotApplyCmd.Flag("wait", "Stops searching after [wait] seconds have elapse (max 30s)").Short('w').Default("2").Int()
	threadSnapshotApplyID = threadSnapshotApplyCmd.Arg("id", "The id of the snapshot to apply").Required().String()

	// ================================
	// @todo see if kingpin has automated no-* support, in which case we do "*").Default("true")
	// token
	tokenCmd = appCmd.Command("token", "Tokens allow other peers to register with a cafe peer").Alias("tokens")

	// create
	tokenCreateCmd = tokenCmd.Command("create", `Generates an access token (44 random bytes) and saves a bcrypt hashed version for future lookup.
The response contains a base58 encoded version of the random bytes token.`)
	tokenCreateNoStore = tokenCreateCmd.Flag("no-store", "If used instead of token, the token is generated but not stored in the local Cafe database.").Short('n').Bool()
	tokenCreateToken = tokenCreateCmd.Flag("token", "If used instead of no-store, use this existing token rather than creating a new one.").Short('t').String()

	// list
	tokenListCmd = tokenCmd.Command("list", "List info about all stored cafe tokens").Alias("ls")

	// validate
	tokenValidateCmd = tokenCmd.Command("validate", "Check validity of existing cafe access token").Alias("valid")
	tokenValidateToken = tokenValidateCmd.Arg("token", "The token to validate").Required().String()

	// remove
	tokenRemoveCmd = tokenCmd.Command("remove", "Removes an existing cafe token").Alias("rm")
	// @todo does this delete or ignore?
	tokenRemoveToken = tokenRemoveCmd.Arg("token", "The token to remove").Required().String()

	// ================================

	// version
	versionCmd = appCmd.Command("version", "Print the current version and exit")
	versionGit = versionCmd.Flag("git", "Show full git version summary").Short('g').Bool()

	// ================================

	// wallet
	walletCmd = appCmd.Command("wallet", "Initialize a new wallet, or view accounts from an existing wallet").Alias("wallets")

	// wallet init
	walletInitCmd       = walletCmd.Command("init", "Initializes a new account wallet backed by a mnemonic recovery phrase")
	walletInitWordCount = walletInitCmd.Flag("word-count", "Number of mnemonic recovery phrase words: 12,15,18,21,24").Short('w').Default("12").Int()
	walletInitPassword  = walletInitCmd.Flag("password", "Mnemonic recovery phrase password (omit if none)").Short('p').String()

	// wallet accounts
	walletAccountsCmd      = walletCmd.Command("accounts", "Shows the derived accounts (address/seed pairs) in a wallet").Alias("account")
	walletAccountsPassword = walletAccountsCmd.Flag("password", "Mnemonic recovery phrase password (omit if none)").Short('p').String()
	walletAccountsDepth    = walletAccountsCmd.Flag("depth", "Number of accounts to show").Short('d').Default("1").Int()
	walletAccountsOffset   = walletAccountsCmd.Flag("offset", "Account depth to start from").Short('o').Default("0").Int()

)

func Run() error {
	switch kingpin.MustParse(appCmd.Parse(os.Args[1:])) {

	// account
	case accountGetCmd.FullCommand():
		return AccountGet()

	case accountSeedCmd.FullCommand():
		return AccountSeed()

	case accountAddressCmd.FullCommand():
		return AccountAddress()

	case accountSyncCmd.FullCommand():
		return AccountSync(*accountSyncWait)

	// block
	case blockListCmd.FullCommand():
		return BlockList(*blockListThreadID, *blockListOffset, *blockListLimit, *blockListDots)

	case blockMetaCmd.FullCommand():
		return BlockMeta(*blockMetaBlockID)

	case blockRemoveCmd.FullCommand():
		return BlockRemove(*blockRemoveBlockID)

	case blockFileCmd.FullCommand():
		return BlockFile(*blockFileBlockID, *blockFileIndex, *blockFilePath,  *blockFileContent)

	// cafe
	case cafeAddCmd.FullCommand():
		return CafeAdd(*cafeAddToken)

	case cafeListCmd.FullCommand():
		return CafeList()

	case cafeGetCmd.FullCommand():
		return CafeGet(*cafeGetCafeID)

	case cafeRemoveCmd.FullCommand():
		return CafeRemove(*cafeRemoveCafeID)

	case cafeMessagesCmd.FullCommand():
		return CafeMessages()

	// chat
	case chatCmd.FullCommand():
		return Chat(*chatThreadID)

	// comments
	case commentAddCmd.FullCommand():
		return CommentAdd(*commentAddBlockID, *commentAddBody)

	case commentListCmd.FullCommand():
		return CommentList(*commentListBlockID)

	case commentGetCmd.FullCommand():
		return CommentGet(*commentGetCommentID)

	case commentIgnoreCmd.FullCommand():
		return CommentIgnore(*commentIgnoreCommentID)

	// config
	case  configCmd.FullCommand():
		return Config(*configName, *configValue)

	// contacts
	case contactAddCmd.FullCommand():
		return ContactAdd(*contactAddName, *contactAddAddress, *contactAddWait)

	case contactListCmd.FullCommand():
		return ContactList()

	case contactGetCmd.FullCommand():
		return ContactGet(*contactGetAddress)

	case contactRemoveCmd.FullCommand():
		return ContactRemove(*contactRemoveAddress)

	case contactSearchCmd.FullCommand():
		return ContactSearch(*contactSearchName, *contactSearchAddress, *contactSearchLocal, *contactSearchRemote, *contactSearchLimit, *contactSearchWait)

	// daemon
	case daemonCmd.FullCommand():
		repoPath, err := getRepoPath(*daemonRepoPath)
		if err != nil {
			return err
		}
		return Daemon(repoPath, *daemonPinCode, *daemonDocs, *logDebug)

		// feed
	case feedCmd.FullCommand():
		return Feed(*feedThreadID, *feedOffset, *feedLimit, *feedMode)

	// file
	case fileListCmd.FullCommand():
		return FileList(*fileListThreadID, *fileListOffset, *fileListLimit)

	case fileKeysCmd.FullCommand():
		return FileKeys(*fileKeysFileID)

	case fileIgnoreCmd.FullCommand():
		return FileIgnore(*fileIgnoreFileID)

	case fileGetCmd.FullCommand():
		return FileGet(*fileGetFileID, *fileGetContent)

	case fileAddCmd.FullCommand():
		return FileAdd(*fileAddPath, *fileAddThreadID, *fileAddCaption, *fileAddGroup, *fileAddVerbose)

	// init
	case initCmd.FullCommand():
		kp, err := keypair.Parse(*initAccountSeed)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("parse account seed failed: %s", err))
		}

		account, ok := kp.(*keypair.Full)
		if !ok {
			return keypair.ErrInvalidKey
		}

		repoPath, err := getRepoPath(*initRepoPath)
		if err != nil {
			return err
		}

		config := core.InitConfig{
			Account:         account,
			PinCode:         *initPinCode,
			RepoPath:        repoPath,
			SwarmPorts:      *ipfsSwarmPorts,
			ApiAddr:         *apiBindAddr,
			CafeApiAddr:     *cafeApiBindAddr,
			GatewayAddr:     *gatewayBindAddr,
			ProfilingAddr:   *profilingBindAddr,
			IsMobile:        false,
			IsServer:        *ipfsServerMode,
			LogToDisk:       *logNoFiles == false,
			Debug:           *logDebug,
			CafeOpen:        *cafeOpen,
			CafeURL:         *cafeURL,
			CafeNeighborURL: *cafeNeighborURL,
		}

		return InitCommand(config)

	// invite
	case inviteCreateCmd.FullCommand():
		return InviteCreate(*inviteCreateThreadID, *inviteCreateAddress, *inviteCreateWait)

	case inviteListCmd.FullCommand():
		return InviteList()

	case inviteAcceptCmd.FullCommand():
		return InviteAccept(*inviteAcceptID, *inviteAcceptKey)

	case inviteIgnoreCmd.FullCommand():
		return InviteIgnore(*inviteIgnoreID)

	// ipfs
	case  ipfsIdCmd.FullCommand():
		return IPFSId()

	case ipfsSwarmConnectCmd.FullCommand():
		return IPFSSwarmConnect(*ipfsSwarmConnectAddress)

	case  ipfsSwarmPeersCmd.FullCommand():
		return IPFSSwarmPeers(*ipfsSwarmPeersVerbose, *ipfsSwarmPeersStreams, *ipfsSwarmPeersLatency, *ipfsSwarmPeersDirection)

	case ipfsCatCmd.FullCommand():
		return IPFSCat(*ipfsCatHash, *ipfsCatKey)

	// like
	case likeAddCmd.FullCommand():
		return LikeAdd(*likeAddBlockID)

	case likeListCmd.FullCommand():
		return LikeAdd(*likeListBlockID)

	case likeGetCmd.FullCommand():
		return LikeAdd(*likeGetLikeID)

	case likeIgnoreCmd.FullCommand():
		return LikeAdd(*likeIgnoreLikeID)

	// log
	case logCmd.FullCommand():
		return Logs(*logSubsystem, *logLevel,  *logTexOnly)

	// message
	case messageAddCmd.FullCommand():
		return MessageAdd(*messageAddThreadID, *messageAddBody)

	case messageListCmd.FullCommand():
		return MessageList(*messageListThreadID, *messageListOffset, *messageListLimit)

	case messageGetCmd.FullCommand():
		return MessageGet(*messageGetMessageID)

	case messageIgnoreCmd.FullCommand():
		return MessageIgnore(*messageIgnoreMessageID)

	// migrate
	case migrateCmd.FullCommand():
		repoPath, err := getRepoPath(*migrateRepoPath)
		if err != nil {
			return err
		}
		return Migrate(repoPath, *migratePinCode)

	// notification
	case notificationListCmd.FullCommand():
		return NotificationList()

	case notificationReadCmd.FullCommand():
		return NotificationRead(*notificationReadID)

	// ping
	case pingCmd.FullCommand():
		return Ping(*pingAddress)

	// profile
	case profileGetCmd.FullCommand():
		return ProfileGet()

	case profileSetCmd.FullCommand():
		return ProfileSet(*profileSetName, *profileSetAvatar)

	// subscribe
	case subscribeCmd.FullCommand():
		return SubscribeCommand(*subscribeThreadID, *subscribeType)

	// summary
	case summaryCmd.FullCommand():
		return Summary()

	// thread
	case threadAddCmd.FullCommand():
		return ThreadAdd(*threadAddKey, *threadAddType, *threadAddSharing, *threadAddWhitelist, *threadAddSchema, *threadAddSchemaFile,  *threadAddBlob, *threadAddCameraRoll, *threadAddMedia)

	case threadListCmd.FullCommand():
		return ThreadList()

	case threadGetCmd.FullCommand():
		return ThreadGet(*threadGetThreadID)

	case threadDefaultCmd.FullCommand():
		return ThreadDefault()

	case threadPeerCmd.FullCommand():
		return ThreadPeer(*threadPeerThreadID)

	case threadRenameCmd.FullCommand():
		return ThreadRename(*threadRenameName, *threadRenameThreadID)

	case threadRemoveCmd.FullCommand():
		return ThreadRemove(*threadRemoveThreadID)

	case threadSnapshotCreateCmd.FullCommand():
		return ThreadSnapshotCreate()

	case threadSnapshotSearchCmd.FullCommand():
		return ThreadSnapshotSearch(*threadSnapshotSearchWait)

	case threadSnapshotApplyCmd.FullCommand():
		return ThreadSnapshotApply(*threadSnapshotApplyID, *threadSnapshotApplyWait)

	// token
	case tokenCreateCmd.FullCommand():
		return TokenCreate(*tokenCreateToken, *tokenCreateNoStore)

	case tokenListCmd.FullCommand():
		return TokenList()

	case tokenValidateCmd.FullCommand():
		return TokenValidate(*tokenValidateToken)

	case tokenRemoveCmd.FullCommand():
		return TokenRemove(*tokenRemoveToken)

	// version
	case versionCmd.FullCommand():
		return Version(*versionGit)

	// wallet
	case walletInitCmd.FullCommand():
		return WalletInit(*walletInitWordCount, *walletInitPassword)

	case walletAccountsCmd.FullCommand():
		return WalletAccounts(*walletAccountsPassword, *walletAccountsDepth, *walletAccountsOffset)

	}

	return nil
}

func executeStringCmd(meth method, pth string, pars params) (string, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := util.UnmarshalString(res.Body)
	if err != nil {
		return "", err
	}

	return body, nil
}

func executeJsonCmd(meth method, pth string, pars params, target interface{}) (string, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf(body)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", nil
	}

	if target == nil {
		target = new(interface{})
	}
	if err := json.Unmarshal(data, target); err != nil {
		return "", err
	}
	jsn, err := json.MarshalIndent(target, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsn), nil
}

func executeJsonPbCmd(meth method, pth string, pars params, target proto.Message) (string, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf(body)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := pbUnmarshaler.Unmarshal(bytes.NewReader(data), target); err != nil {
		return "", err
	}
	jsn, err := pbMarshaler.MarshalToString(target)
	if err != nil {
		return "", err
	}

	return jsn, nil
}

func executeBlobCmd(meth method, pth string, pars params) error {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(body)
	}

	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		return err
	}

	return nil
}

func request(meth method, pth string, pars params) (*http.Response, func(), error) {
	apiURL := fmt.Sprintf("%s/api/%s/%s", *apiAddr, *apiVersion, pth)
	if *logDebug {
		fmt.Println(apiURL)
	}
	req, err := http.NewRequest(string(meth), apiURL, pars.payload)
	if err != nil {
		return nil, nil, err
	}

	if len(pars.args) > 0 {
		var args []string
		for _, arg := range pars.args {
			args = append(args, url.PathEscape(arg))
		}
		req.Header.Set("X-Textile-Args", strings.Join(args, ","))
	}

	if len(pars.opts) > 0 {
		var items []string
		for k, v := range pars.opts {
			items = append(items, k+"="+url.PathEscape(v))
		}
		req.Header.Set("X-Textile-Opts", strings.Join(items, ","))
	}

	if pars.ctype != "" {
		req.Header.Set("Content-Type", pars.ctype)
	}

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	cancel := func() {
		tr.CancelRequest(req)
	}

	return res, cancel, err
}

func handleSearchStream(pth string, param params) []pb.QueryResult {
	var results []pb.QueryResult
	outputCh := make(chan interface{})

	cancel := func() {}
	done := make(chan struct{})
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		defer func() {
			cancel()
			done <- struct{}{}
		}()

		var res *http.Response
		var err error
		res, cancel, err = request(http.MethodPost, pth, param)
		if err != nil {
			outputCh <- err.Error()
			return
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			body, err := util.UnmarshalString(res.Body)
			if err != nil {
				outputCh <- err.Error()
			} else {
				outputCh <- body
			}
			return
		}

		decoder := json.NewDecoder(res.Body)
		for decoder.More() {
			var result pb.QueryResult
			if err := pbUnmarshaler.UnmarshalNext(decoder, &result); err == io.EOF {
				return
			} else if err != nil {
				outputCh <- err.Error()
				return
			}
			results = append(results, result)

			out, err := pbMarshaler.MarshalToString(&result)
			if err != nil {
				outputCh <- err.Error()
				return
			}
			outputCh <- out
		}
	}()

	for {
		select {
		case val := <-outputCh:
			output(val)

		case <-quit:
			fmt.Println("Interrupted")
			if cancel != nil {
				fmt.Printf("Canceling...")
				cancel()
			}
			fmt.Print("done\n")
			os.Exit(1)

		case <-done:
			return results
		}
	}
}

// https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func confirm(q string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", q)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func output(val interface{}) {
	if val.(string) == "" {
		val = "ok"
	}
	fmt.Println(val)
}

// Get the repo path for the user, will create it if missing
// Unless provided, it defaults to ~/.textile/repo
func getRepoPath(repoPath string) (string, error) {
	if len(repoPath) == 0 {
		// get homedir
		home, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("get homedir failed: %s", err))
		}

		// ensure app folder is created
		appDir := filepath.Join(home, ".textile")
		if err := os.MkdirAll(appDir, 0755); err != nil {
			return "", fmt.Errorf(fmt.Sprintf("create repo directory failed: %s", err))
		}
		repoPath = filepath.Join(appDir, "repo")
	}
	return repoPath, nil
}

