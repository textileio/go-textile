package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	_ "expvar"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	logging "github.com/ipfs/go-log"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdsMap map[string]func() error

func threadBlocksCommand(cmds cmdsMap, parent *kingpin.CmdClause, names []string) *kingpin.CmdClause {
	cmd := parent.Command(names[0], "Paginates blocks in a thread")
	for _, name := range names[1:] {
		cmd = cmd.Alias(name)
	}

	blockListThreadID := cmd.Arg("thread", "Thread ID").Required().String()
	blockListOffset := cmd.Flag("offset", "Offset ID to start listing from").Short('o').String()
	blockListLimit := cmd.Flag("limit", "List page size").Short('l').Default("5").Int()
	blockListDots := cmd.Flag("dots", "Return GraphViz dots instead of JSON").Short('d').Bool()
	cmds[cmd.FullCommand()] = func() error {
		return BlockList(*blockListThreadID, *blockListOffset, *blockListLimit, *blockListDots)
	}

	return cmd
}

func threadFilesCommand(cmds cmdsMap, parent *kingpin.CmdClause, names []string) *kingpin.CmdClause {
	cmd := parent.Command(names[0], "Paginates the files of a thread, or of all threads")
	for _, name := range names[1:] {
		cmd = cmd.Alias(name)
	}

	threadID := cmd.Arg("thread", "Thread ID, omit for all").String()
	offset := cmd.Flag("offset", "Offset ID to start listing from").Short('o').String()
	limit := cmd.Flag("limit", "List page size").Short('l').Default("5").Int()
	cmds[cmd.FullCommand()] = func() error {
		return FileListThread(*threadID, *offset, *limit)
	}

	return cmd
}

func blockFilesCommand(cmds cmdsMap, parent *kingpin.CmdClause, names []string) *kingpin.CmdClause {
	cmd := parent.Command(names[0], "Commands to interact with File Blocks")
	for _, name := range names[1:] {
		cmd = cmd.Alias(name)
	}

	listCmd := cmd.Command("list", "List all files within a File Block").Alias("ls").Default()
	listBlockID := listCmd.Arg("block", "File Block ID").Required().String()
	cmds[listCmd.FullCommand()] = func() error {
		return FileListBlock(*listBlockID)
	}

	getCmd := cmd.Command("get", "Get a specific file within the File Block")
	getBlockID := getCmd.Arg("block", "File Block ID").Required().String()
	getIndex := getCmd.Flag("index", "The index of the file to fetch").Default("0").Int()
	getPath := getCmd.Flag("path", "The link path of the file to fetch").Default(".").String()
	getContent := getCmd.Flag("content", "If provided, the decrypted content of the file is retrieved").Bool()
	cmds[getCmd.FullCommand()] = func() error {
		return FileGetBlock(*getBlockID, *getIndex, *getPath, *getContent)
	}

	return cmd
}

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
	log  = logging.Logger("tex-main")

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
	appCmd      = kingpin.New("textile", "Textile is a set of tools and trust-less infrastructure for building censorship resistant and privacy preserving applications")
	apiAddr     = appCmd.Flag("api", "API Address to use").Envar("API").Default("http://127.0.0.1:40600").String()
	apiVersion  = appCmd.Flag("api-version", "API version to use").Envar("API_VERSION").Default("v0").String()
	logDebug    = appCmd.Flag("debug", "Set the logging level to debug").Bool()
	appUsername = appCmd.Flag("username", "Specify the username (address) if required for Basic Auth").Envar("TEXTILE_USERNAME").String()
	appPassword = appCmd.Flag("password", "Specify the password (pincode) used for datastore encryption and Basic Auth (omit if no auth/encryption is used)").Envar("TEXTILE_PASSWORD").String()
)

func Run() error {
	cmds := make(cmdsMap)

	// configure
	appCmd.UsageTemplate(kingpin.CompactUsageTemplate)

	// ================================

	// account
	accountCmd := appCmd.Command("account", "Manage the account that the initialised textile repository is associated with")

	// account get
	accountGetCmd := accountCmd.Command("get", "Shows the local peer's account info as a contact").Default()
	cmds[accountGetCmd.FullCommand()] = AccountGet

	// account seed
	accountSeedCmd := accountCmd.Command("seed", "Shows the local peer's account seed, treat this as top secret")
	cmds[accountSeedCmd.FullCommand()] = AccountSeed

	// account address
	accountAddressCmd := accountCmd.Command("address", "Shows the local peer's account address")
	cmds[accountAddressCmd.FullCommand()] = AccountAddress

	// account sync
	accountSyncCmd := accountCmd.Command("sync", "Syncs the local account peer with other peers found on the network")
	accountSyncWait := accountSyncCmd.Flag("wait", "Stops searching after 'wait' seconds have elapsed (max 30s)").Default("2").Int()
	cmds[accountSyncCmd.FullCommand()] = func() error {
		return AccountSync(*accountSyncWait)
	}

	// ================================

	// block
	blockCmd := appCmd.Command("block", "Threads are composed of an append-only log of blocks, use these commands to manage them").Alias("blocks")

	// block list
	threadBlocksCommand(cmds, blockCmd, []string{"list", "ls"})

	// block meta
	blockMetaCmd := blockCmd.Command("meta", "Get the metadata for a block").Alias("get")
	blockMetaBlockID := blockMetaCmd.Arg("block", "Block ID").Required().String()
	cmds[blockMetaCmd.FullCommand()] = func() error {
		return BlockMeta(*blockMetaBlockID)
	}

	// block ignore
	blockIgnoreCmd := blockCmd.Command("ignore", "Remove a block by marking it to be ignored").Alias("remove").Alias("rm")
	blockIgnoreBlockID := blockIgnoreCmd.Arg("block", "Block ID").Required().String()
	cmds[blockIgnoreCmd.FullCommand()] = func() error {
		return BlockIgnore(*blockIgnoreBlockID)
	}

	// block file alias
	blockFilesCommand(cmds, blockCmd, []string{"files", "file"})

	// ================================

	// bots
	botsCmd := appCmd.Command("bots", "Commands to manage bots").Alias("bot")

	// bots list
	botsListCmd := botsCmd.Command("list", "List info about all active bots").Alias("ls").Default()
	cmds[botsListCmd.FullCommand()] = BotsList

	// bots disable
	botsDisableCmd := botsCmd.Command("disable", "Disable a bot")
	botsDisableID := botsDisableCmd.Arg("id", "ID of the bot").Required().String()
	cmds[botsDisableCmd.FullCommand()] = func() error {
		return BotsDisable(*botsDisableID)
	}

	// bots enable
	botsEnableCmd := botsCmd.Command("enable", "Enable a bot")
	botsEnableID := botsEnableCmd.Arg("id", "ID of the bot").Required().String()
	botsEnableCafe := botsEnableCmd.Flag("cafe-api", "Whether to serve bot on the Cafe API (public)").Short('c').Bool()

	cmds[botsEnableCmd.FullCommand()] = func() error {
		return BotsEnable(*botsEnableID, *botsEnableCafe)
	}

	// bots create
	botsNewCmd := botsCmd.Command("create", "Initialize a new bot for development")
	botsNewName := botsNewCmd.Arg("name", "Name of the bot").Required().String()
	cmds[botsNewCmd.FullCommand()] = func() error {
		return BotsCreate(*botsNewName)
	}

	// ================================

	// cafe
	cafeCmd := appCmd.Command("cafe", "Commands to manage cafes").Alias("cafes")

	// cafe add
	cafeAddCmd := cafeCmd.Command("add", `Registers with a cafe and saves an expiring service session token.
An access token is required to register, and should be obtained separately from the target cafe.`)
	cafeAddPeerID := cafeAddCmd.Arg("peer", "The host cafe's IPFS peer ID").Required().String()
	cafeAddToken := cafeAddCmd.Flag("token", "An access token supplied by the cafe").Short('t').Required().String()
	// @todo is this consistent with the rest?
	cmds[cafeAddCmd.FullCommand()] = func() error {
		return CafeAdd(*cafeAddPeerID, *cafeAddToken)
	}

	// cafe list
	cafeListCmd := cafeCmd.Command("list", "List info about all active cafe sessions").Alias("ls").Default()
	cmds[cafeListCmd.FullCommand()] = CafeList

	// cafe get
	cafeGetCmd := cafeCmd.Command("get", "Gets and displays info about a cafe session")
	cafeGetCafeID := cafeGetCmd.Arg("cafe", "Cafe ID").Required().String()
	cmds[cafeGetCmd.FullCommand()] = func() error {
		return CafeGet(*cafeGetCafeID)
	}

	// cafe delete
	cafeDeleteCmd := cafeCmd.Command("delete", "Deregisters a cafe (content will expire based on the cafe's service rules)").Alias("del").Alias("remove").Alias("rm")
	cafeDeleteCafeID := cafeDeleteCmd.Arg("cafe", "Cafe ID").Required().String()
	cmds[cafeDeleteCmd.FullCommand()] = func() error {
		return CafeDelete(*cafeDeleteCafeID)
	}

	// cafe messages
	cafeMessagesCmd := cafeCmd.Command("messages", "Check for messages at all cafes. New messages are downloaded and processed opportunistically.")
	cmds[cafeMessagesCmd.FullCommand()] = CafeMessages

	// ================================

	// chat
	chatCmd := appCmd.Command("chat", `Starts an interactive chat session in a thread`)
	chatThreadID := chatCmd.Arg("thread", "Thread ID").Required().String()
	cmds[chatCmd.FullCommand()] = func() error {
		return Chat(*chatThreadID)
	}

	// ================================

	// comment
	commentCmd := appCmd.Command("comment", "Comments are added as blocks in a thread, which target another block, usually a file(s)").Alias("comments")

	// comment add
	commentAddCmd := commentCmd.Command("add", "Attach a comment to a block")
	commentAddBlockID := commentAddCmd.Arg("block", "The Block ID to attach the comment to").Required().String()
	commentAddBody := commentAddCmd.Arg("body", "Text to use as the comment").Required().String()
	cmds[commentAddCmd.FullCommand()] = func() error {
		return CommentAdd(*commentAddBlockID, *commentAddBody)
	}

	// comment list
	commentListCmd := commentCmd.Command("list", "Get the comments that are attached to a block").Alias("ls").Default()
	commentListBlockID := commentListCmd.Arg("block", "The Block ID which the comments attached to").Required().String()
	cmds[commentListCmd.FullCommand()] = func() error {
		return CommentList(*commentListBlockID)
	}

	// comment get
	commentGetCmd := commentCmd.Command("get", "Get a comment by its own Block ID")
	commentGetBlockID := commentGetCmd.Arg("comment-block", "Comment Block ID").Required().String()
	cmds[commentGetCmd.FullCommand()] = func() error {
		return CommentGet(*commentGetBlockID)
	}

	// comment ignore
	commentIgnoreCmd := commentCmd.Command("ignore", "Ignore a comment by its own Block ID").Alias("remove").Alias("rm")
	commentIgnoreBlockID := commentIgnoreCmd.Arg("comment-block", "Comment Block ID").Required().String()
	cmds[commentIgnoreCmd.FullCommand()] = func() error {
		return CommentIgnore(*commentIgnoreBlockID)
	}

	// ================================

	// config
	configCmd := appCmd.Command("config", "Get or set configuration variables").Alias("conf")
	configName := configCmd.Arg("name", "If provided, will restrict the operation to this specific configuration variable, e.g. 'Addresses.API'").String()
	configValue := configCmd.Arg("value", `If provided, will set the specific configuration variable to this JSON escaped value, e.g. '"127.0.0.1:40600"'`).String()
	cmds[configCmd.FullCommand()] = func() error {
		return Config(*configName, *configValue)
	}

	// ================================

	// contact
	contactCmd := appCmd.Command("contact", "Manage local contacts and find other contacts on the network").Alias("contacts")

	// contact add
	contactAddCmd := contactCmd.Command("add", "Adds a contact by display name or account address to known contacts")
	contactAddName := contactAddCmd.Flag("name", "Add by display name").Short('n').String()
	contactAddAddress := contactAddCmd.Flag("address", "Add by account address").Short('a').String()
	contactAddWait := contactAddCmd.Flag("wait", "Stops searching after [wait] seconds have elapsed").Int()
	cmds[contactAddCmd.FullCommand()] = func() error {
		return ContactAdd(*contactAddName, *contactAddAddress, *contactAddWait)
	}

	// contact list
	contactListCmd := contactCmd.Command("list", "List known contacts").Alias("ls").Default()
	cmds[contactListCmd.FullCommand()] = ContactList

	// contact get
	contactGetCmd := contactCmd.Command("get", "Gets a known local contact")
	contactGetAddress := contactGetCmd.Arg("address", "Account Address").Required().String()
	cmds[contactGetCmd.FullCommand()] = func() error {
		return ContactGet(*contactGetAddress)
	}

	// contact delete
	contactDeleteCmd := contactCmd.Command("delete", "Deletes a known contact").Alias("del").Alias("remove").Alias("rm")
	contactDeleteAddress := contactDeleteCmd.Arg("address", "Account Address").Required().String()
	cmds[contactDeleteCmd.FullCommand()] = func() error {
		return ContactDelete(*contactDeleteAddress)
	}

	// contact search
	contactSearchCmd := contactCmd.Command("search", "Searches locally and on the network for contacts").Alias("find")
	contactSearchName := contactSearchCmd.Flag("name", "Search by display name").Short('n').String()
	contactSearchAddress := contactSearchCmd.Flag("address", "Search by account address").Short('a').String()
	contactSearchLocal := contactSearchCmd.Flag("only-local", "Only search local contacts").Bool()
	contactSearchRemote := contactSearchCmd.Flag("only-remote", "Only search remote contacts").Bool()
	contactSearchLimit := contactSearchCmd.Flag("limit", "Stops searching after [limit] results are found").Default("5").Int()
	contactSearchWait := contactSearchCmd.Flag("wait", "Stops searching after [wait] seconds have elapsed (max 30s)").Default("2").Int()
	cmds[contactSearchCmd.FullCommand()] = func() error {
		return ContactSearch(*contactSearchName, *contactSearchAddress, *contactSearchLocal, *contactSearchRemote, *contactSearchLimit, *contactSearchWait)
	}
	// @todo why not make this part of `textile contact list`?
	// @todo perhaps our `list` and `get` commands can be merged

	// ================================

	// daemon
	daemonCmd := appCmd.Command("daemon", "Start a node daemon session")
	daemonRepo := daemonCmd.Flag("repo", "Specify a custom path to the repo directory").Short('r').String()
	daemonBaseRepo := daemonCmd.Flag("base-repo", "Specify a custom path to the base repo directory").Short('b').String()
	daemonAccountAddress := daemonCmd.Flag("account-address", "Specify an existing account address").Short('a').String()
	daemonPin := daemonCmd.Flag("pin", "Specify the pin code for datastore encryption (omit no pin code was used during init)").Short('p').String()
	daemonDocs := daemonCmd.Flag("serve-docs", "Whether to serve the local REST API docs").Short('s').Bool()
	cmds[daemonCmd.FullCommand()] = func() error {
		repo, err := getRepo(*daemonRepo, *daemonBaseRepo, *daemonAccountAddress)
		if err != nil {
			return err
		}
		return Daemon(repo, *daemonPin, *daemonDocs, *logDebug)
	}

	// ================================

	// docs
	docsCmd := appCmd.Command("docs", "Prints the CLI help as HTML")
	cmds[docsCmd.FullCommand()] = Docs

	// ================================

	// feed
	feedCmd := appCmd.Command("feed", `Paginates post (join|leave|files|message) and annotation (comment|like) block types as a consumable feed.

The --mode option dictates how the feed is displayed:

-  "chrono": All feed block types are shown. Annotations always nest their target post, i.e., the post a comment is about.
-  "annotated": Annotations are nested under post targets, but are not shown in the top-level feed.
-  "stacks": Related blocks are chronologically grouped into "stacks". A new stack is started if an unrelated block
   breaks continuity. This mode is used by Textile Photos.

Stacks may include:

- The initial post with some nested annotations. Newer annotations may have already been listed.
- One or more annotations about a post. The newest annotation assumes the "top" position in the stack. Additional
 annotations are nested under the target. Newer annotations may have already been listed in the case as well.`)
	feedThreadID := feedCmd.Arg("thread", "Thread ID, omit for all").String()
	feedOffset := feedCmd.Flag("offset", "Offset ID to start listening from").Short('o').String()
	feedLimit := feedCmd.Flag("limit", "List page size").Short('l').Default("3").Int()
	feedMode := feedCmd.Flag("mode", "Feed mode, one of: chrono, annotated, stacks").Short('m').Default("chrono").String()
	// ^ when kingpin v2 lands with enumerables, we could move the usage docs to the enum docs
	cmds[feedCmd.FullCommand()] = func() error {
		return Feed(*feedThreadID, *feedOffset, *feedLimit, *feedMode)
	}

	// ================================

	// file
	fileCmd := appCmd.Command("file", "Manage Textile Files Blocks").Alias("files").Alias("data")
	// @todo rename this to Textile Data Blocks: https://github.com/textileio/meta/issues/31

	// file list
	fileListCmd := fileCmd.Command("list", `Get all the files, or just the files for a specific thread or block`).Alias("ls").Default()

	// file list thread
	fileListThreadCmd := threadFilesCommand(cmds, fileListCmd, []string{"thread"})
	fileListThreadCmd.Default()

	// file list block
	blockFilesCommand(cmds, fileListCmd, []string{"block"})

	// file keys
	fileKeysCmd := fileCmd.Command("keys", "Shows the encryption keys for each content/meta pair for the given block DAG").Alias("key")
	fileKeysDataID := fileKeysCmd.Arg("block-data", "Block Data ID").Required().String()
	cmds[fileKeysCmd.FullCommand()] = func() error {
		return FileKeys(*fileKeysDataID)
	}

	// file add
	fileAddCmd := fileCmd.Command("add", `Adds a file, directory, or hash to a thread. Files not supported by the thread schema are ignored`)
	fileAddThreadID := fileAddCmd.Arg("thread", "Thread ID").Required().String()
	fileAddPath := fileAddCmd.Arg("path", "The path to the file or directory to add, can also be an existing hash. If omitted, you must provide a stdin blob input.").String()
	fileAddCaption := fileAddCmd.Flag("caption", "File(s) caption").Short('c').String()
	fileAddGroup := fileAddCmd.Flag("group", "If provided, group a directory's files together into a single object, includes nested directories").Short('g').Bool()
	fileAddVerbose := fileAddCmd.Flag("verbose", "Prints files as they are milled").Short('v').Bool()
	cmds[fileAddCmd.FullCommand()] = func() error {
		return FileAdd(*fileAddPath, *fileAddThreadID, *fileAddCaption, *fileAddGroup, *fileAddVerbose)
	}

	// file ignore
	fileIgnoreCmd := fileCmd.Command("ignore", `Ignores a thread file by its own block ID`).Alias("remove").Alias("rm")
	fileIgnoreBlockID := fileIgnoreCmd.Arg("files-block", "Files Block ID").Required().String()
	cmds[fileIgnoreCmd.FullCommand()] = func() error {
		return FileIgnore(*fileIgnoreBlockID)
	}

	// file get
	fileGetCmd := fileCmd.Command("get", "Get the metadata or content of a specific file")
	fileGetHash := fileGetCmd.Arg("hash", "File Hash").Required().String()
	fileGetContent := fileGetCmd.Flag("content", "If provided, the decrypted content of the file is retrieved").Bool()
	cmds[fileGetCmd.FullCommand()] = func() error {
		return FileGet(*fileGetHash, *fileGetContent)
	}

	// ================================

	// init
	initCmd := appCmd.Command("init", "Configure textile to use the account by creating a local repository to house its data")
	initRepo := initCmd.Flag("repo", "Specify a custom path to the repo directory").Short('r').String()
	initBaseRepo := initCmd.Flag("base-repo", "Specify a custom path to the base repo directory").Short('b').String()
	initAccountSeed := initCmd.Arg("account-seed", "The account seed to use, if you do not have one, refer to: textile wallet --help").Required().String()
	initPin := initCmd.Flag("pin", "Specify a pin for datastore encryption").Short('p').String()
	initIpfsServerMode := initCmd.Flag("server", "Apply IPFS server profile").Bool()
	initIpfsSwarmPorts := initCmd.Flag("swarm-ports", "Set the swarm ports (TCP,WS). A random TCP port is chosen by default").String()
	initLogFiles := initCmd.Flag("log-files", "If true, writes logs to rolling files, if false, writes logs to stdout").Default("false").Bool()
	initApiBindAddr := initCmd.Flag("api-bind-addr", "Set the local API address").Default("127.0.0.1:40600").String()
	initCafeApiBindAddr := initCmd.Flag("cafe-bind-addr", "Set the cafe REST API address").Default("0.0.0.0:40601").String()
	initGatewayBindAddr := initCmd.Flag("gateway-bind-addr", "Set the IPFS gateway address").Default("127.0.0.1:5050").String()
	initProfilingBindAddr := initCmd.Flag("profile-bind-addr", "Set the profiling address").Default("127.0.0.1:6060").String()
	initCafe := initCmd.Flag("cafe", "Open the p2p cafe service for other peers").Bool()
	initCafeOpen := initCmd.Flag("cafe-open", "Open the p2p cafe service for other peers").Hidden().Bool() // hidden alias
	initCafeURL := initCmd.Flag("cafe-url", "Specify a custom URL of this cafe, e.g., https://mycafe.com").Envar("CAFE_HOST_URL").String()
	initCafeNeighborURL := initCmd.Flag("cafe-neighbor-url", "Specify the URL of a secondary cafe. Must return cafe info, e.g., via a Gateway: https://my-gateway.yolo.com/cafe, or a cafe API: https://my-cafe.yolo.com").Envar("CAFE_HOST_NEIGHBOR_URL").String()
	cmds[initCmd.FullCommand()] = func() error {
		kp, err := keypair.Parse(*initAccountSeed)
		if err != nil {
			return fmt.Errorf("parse account seed failed: %s", err)
		}

		account, ok := kp.(*keypair.Full)
		if !ok {
			return keypair.ErrInvalidKey
		}

		var repo = *initRepo
		var baseRepo = *initBaseRepo

		// default if neither repo or base-repo is specified, set a value for repo
		if len(repo) == 0 && len(baseRepo) == 0 {
			repo, err = getDefaultRepo()
			if err != nil {
				return err
			}
		}

		config := core.InitConfig{
			Account:         account,
			PinCode:         *initPin, // @todo rename to pin
			RepoPath:        repo,
			BaseRepoPath:    baseRepo,
			SwarmPorts:      *initIpfsSwarmPorts,
			ApiAddr:         *initApiBindAddr,
			CafeApiAddr:     *initCafeApiBindAddr,
			GatewayAddr:     *initGatewayBindAddr,
			ProfilingAddr:   *initProfilingBindAddr,
			IsMobile:        false,
			IsServer:        *initIpfsServerMode,
			LogToDisk:       *initLogFiles,
			Debug:           *logDebug,
			CafeOpen:        *initCafe || *initCafeOpen,
			CafeURL:         *initCafeURL,
			CafeNeighborURL: *initCafeNeighborURL,
		}

		return InitCommand(config)
	}

	// ================================
	// Invites are blocks
	// but they do not stay on the update chain / graph
	// and inbound invites get indexed into the invites table

	// invite
	inviteCmd := appCmd.Command("invite", `Invites allow other users to join threads.

There are two types of invites, direct account-to-account and external:

- Account-to-account invites are encrypted with the invitee's account address (public key).
- External invites are encrypted with a single-use key and are useful for onboarding new users.`).Alias("invites")

	// invite create
	inviteCreateCmd := inviteCmd.Command("create", "Creates a direct account-to-account or external invite to a thread")
	inviteCreateThreadID := inviteCreateCmd.Arg("thread", "Thread ID").Required().String()
	inviteCreateAddress := inviteCreateCmd.Flag("address", "Account Address, omit to create an external invite").Short('a').String()
	inviteCreateWait := inviteCreateCmd.Flag("wait", "Stops searching after [wait] seconds have elapsed (max 30s)").Default("2").Int()
	cmds[inviteCreateCmd.FullCommand()] = func() error {
		return InviteCreate(*inviteCreateThreadID, *inviteCreateAddress, *inviteCreateWait)
	}

	// invite list
	inviteListCmd := inviteCmd.Command("list", "Lists all pending thread invites").Alias("ls").Default()
	cmds[inviteListCmd.FullCommand()] = InviteList

	// invite accept
	inviteAcceptCmd := inviteCmd.Command("accept", "Accepts a direct account-to-account or external invite to a thread")
	inviteAcceptID := inviteAcceptCmd.Arg("id", "Invite ID that you have received").Required().String()
	inviteAcceptKey := inviteAcceptCmd.Flag("key", "Key for an external invite").Short('k').String()
	cmds[inviteAcceptCmd.FullCommand()] = func() error {
		return InviteAccept(*inviteAcceptID, *inviteAcceptKey)
	}

	// invite ignore
	inviteIgnoreCmd := inviteCmd.Command("ignore", "Ignores a direct account-to-account invite to a thread").Alias("remove").Alias("rm")
	inviteIgnoreID := inviteIgnoreCmd.Arg("id", "Invite ID that you wish to ignore").Required().String()
	cmds[inviteIgnoreCmd.FullCommand()] = func() error {
		return InviteIgnore(*inviteIgnoreID)
	}

	// ================================

	// ipfs
	ipfsCmd := appCmd.Command("ipfs", "Provides access to some IPFS commands")

	// ipfs peer
	ipfsPeerCmd := ipfsCmd.Command("peer", "Shows the local node's IPFS peer ID").Alias("id").Default()
	cmds[ipfsPeerCmd.FullCommand()] = func() error {
		return IpfsPeer()
	}

	// ipfs swarm
	ipfsSwarmCmd := ipfsCmd.Command("swarm", "Provides access to a limited set of IPFS swarm commands")

	// ipfs swarm connect
	ipfsSwarmConnectCmd := ipfsSwarmCmd.Command("connect", `Opens a new direct connection to a peer address`)
	ipfsSwarmConnectAddress := ipfsSwarmConnectCmd.Arg("address", `An IPFS multiaddr, such as: /ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ`).String()
	cmds[ipfsSwarmConnectCmd.FullCommand()] = func() error {
		return IpfsSwarmConnect(*ipfsSwarmConnectAddress)
	}

	// ipfs swarm peers
	ipfsSwarmPeersCmd := ipfsSwarmCmd.Command("peers", "Lists the set of peers this node is connected to")
	ipfsSwarmPeersVerbose := ipfsSwarmPeersCmd.Flag("verbose", "Display all extra information").Short('v').Bool()
	ipfsSwarmPeersStreams := ipfsSwarmPeersCmd.Flag("streams", "Also list information about open streams for search peer").Short('s').Bool()
	ipfsSwarmPeersLatency := ipfsSwarmPeersCmd.Flag("latency", "Also list information about the latency to each peer").Short('l').Bool()
	ipfsSwarmPeersDirection := ipfsSwarmPeersCmd.Flag("direction", "Also list information about the direction of connection").Short('d').Bool()
	cmds[ipfsSwarmPeersCmd.FullCommand()] = func() error {
		return IpfsSwarmPeers(*ipfsSwarmPeersVerbose, *ipfsSwarmPeersStreams, *ipfsSwarmPeersLatency, *ipfsSwarmPeersDirection)
	}

	// ipfs cat
	ipfsCatCmd := ipfsCmd.Command("cat", "Displays the data behind an IPFS CID (hash)")
	ipfsCatHash := ipfsCatCmd.Arg("hash", "IPFS CID").Required().String()
	ipfsCatKey := ipfsCatCmd.Flag("key", "Encryption key").Short('k').String()
	cmds[ipfsCatCmd.FullCommand()] = func() error {
		return IpfsCat(*ipfsCatHash, *ipfsCatKey)
	}

	// ================================

	// like
	likeCmd := appCmd.Command("like", `Likes are added as blocks in a thread, which target another block`).Alias("likes")

	// like add
	likeAddCmd := likeCmd.Command("add", "Attach a like to a block")
	likeAddBlockID := likeAddCmd.Arg("block", "Block ID to like, usually a file's block").Required().String()
	cmds[likeAddCmd.FullCommand()] = func() error {
		return LikeAdd(*likeAddBlockID)
	}

	// like list
	likeListCmd := likeCmd.Command("list", "Get likes that are attached to a block").Alias("ls").Default()
	likeListBlockID := likeListCmd.Arg("block", "Block ID to like, usually a file's block").Required().String()
	cmds[likeListCmd.FullCommand()] = func() error {
		return LikeList(*likeListBlockID)
	}

	// like get
	likeGetCmd := likeCmd.Command("get", "Get a like by its own Block ID")
	likeGetLikeID := likeGetCmd.Arg("like-block", "Like Block ID").Required().String()
	cmds[likeGetCmd.FullCommand()] = func() error {
		return LikeGet(*likeGetLikeID)
	}

	// like ignore
	likeIgnoreCmd := likeCmd.Command("ignore", "Ignore a like by its own Block ID").Alias("remove").Alias("rm")
	likeIgnoreLikeID := likeIgnoreCmd.Arg("like-block", "Like Block ID").Required().String()
	cmds[likeIgnoreCmd.FullCommand()] = func() error {
		return LikeIgnore(*likeIgnoreLikeID)
	}

	// ================================

	// log
	logCmd := appCmd.Command("log", `List or change the verbosity of one or all subsystems log output. Textile logs piggyback on the IPFS event logs.`).Alias("logs")
	logSubsystem := logCmd.Flag("subsystem", "The subsystem logging identifier, omit for all").Short('s').String()
	logLevel := logCmd.Flag("level", "One of: debug, info, warning, error, critical. Omit to get current level.").Short('l').String()
	logTextileOnly := logCmd.Flag("textile-only", "Whether to list/change only Textile subsystems, or all available subsystems").Short('t').Bool()
	cmds[logCmd.FullCommand()] = func() error {
		return Logs(*logSubsystem, *logLevel, *logTextileOnly)
	}

	// ================================

	// message
	messageCmd := appCmd.Command("message", "Manage Textile Messages").Alias("messages")

	// message add
	messageAddCmd := messageCmd.Command("add", "Adds a message to a thread")
	messageAddThreadID := messageAddCmd.Arg("thread", "Thread ID").Required().String()
	messageAddBody := messageAddCmd.Arg("body", "The message to add the thread").Required().String()
	cmds[messageAddCmd.FullCommand()] = func() error {
		return MessageAdd(*messageAddThreadID, *messageAddBody)
	}

	// message list
	messageListCmd := messageCmd.Command("list", "Paginates thread messages").Alias("ls").Default()
	messageListThreadID := messageListCmd.Arg("thread", "Thread ID, omit to paginate all messages").String()
	messageListOffset := messageListCmd.Flag("offset", "Offset ID to start the listing from").Short('o').String()
	messageListLimit := messageListCmd.Flag("limit", "List page size").Default("10").Short('l').Int()
	cmds[messageListCmd.FullCommand()] = func() error {
		return MessageList(*messageListThreadID, *messageListOffset, *messageListLimit)
	}

	// message get
	messageGetCmd := messageCmd.Command("get", "Gets a message by its own Block ID")
	messageGetBlockID := messageGetCmd.Arg("message-block", "Message Block ID").String()
	cmds[messageGetCmd.FullCommand()] = func() error {
		return MessageGet(*messageGetBlockID)
	}

	// message ignore
	messageIgnoreCmd := messageCmd.Command("ignore", "Ignores a message by its own Block ID").Alias("remove").Alias("rm")
	messageIgnoreBlockID := messageIgnoreCmd.Arg("message-block", "Message Block ID").String()
	cmds[messageIgnoreCmd.FullCommand()] = func() error {
		return MessageIgnore(*messageIgnoreBlockID)
	}

	// ================================

	// migrate
	migrateCmd := appCmd.Command("migrate", "Migrate the node repository and exit")
	migrateRepo := migrateCmd.Flag("repo", "Specify a custom path to the repo directory").Short('r').String()
	migrateBaseRepo := migrateCmd.Flag("base-repo", "Specify a custom path to the base repo directory").Short('b').String()
	migrateAccountAddress := migrateCmd.Flag("account-address", "Specify an existing account address").Short('a').String()
	cmds[migrateCmd.FullCommand()] = func() error {
		repo, err := getRepo(*migrateRepo, *migrateBaseRepo, *migrateAccountAddress)
		if err != nil {
			return err
		}
		return Migrate(repo, *appPassword)
	}

	// ================================
	// Notifications are local-only, and most block updates generate them
	// E.g. https://github.com/textileio/go-textile/blob/72a910879b5b8135d3cf65c5348beeb5aa4226a0/core/threads_service.go#L395

	// notification
	notificationCmd := appCmd.Command("notification", "Manage notifications that have been generated by thread and account activity").Alias("notifications")

	// notification list
	notificationListCmd := notificationCmd.Command("list", "Lists all notifications").Alias("ls").Default()
	cmds[notificationListCmd.FullCommand()] = func() error {
		return NotificationList()
	}

	// notification read
	notificationReadCmd := notificationCmd.Command("read", "Marks a notification as read")
	notificationReadID := notificationReadCmd.Arg("id", "Notification ID, set to [all] to mark all notifications as read").Required().String()
	cmds[notificationReadCmd.FullCommand()] = func() error {
		return NotificationRead(*notificationReadID)
	}

	// delete
	// @todo add delete notification command
	// https://github.com/textileio/go-textile/issues/823

	// ================================

	// ping
	pingCmd := appCmd.Command("ping", "Pings another peer on the network, returning [online] or [offline]")
	pingAddress := pingCmd.Arg("address", "The address of the other peer on the network").Required().String()
	cmds[pingCmd.FullCommand()] = func() error {
		return Ping(*pingAddress)
	}

	// publish
	publishCmd := appCmd.Command("publish", "Publishes stdin to a topic on the network")
	publishTopic := publishCmd.Arg("topic", "The topic to publish to").Required().String()
	cmds[publishCmd.FullCommand()] = func() error {
		return Publish(*publishTopic)
	}

	// ================================

	// profile
	profileCmd := appCmd.Command("profile", `Manage the profile for your Textile Account, each peer will have its own profile`)

	// profile get
	profileGetCmd := profileCmd.Command("get", "Gets the local peer profile").Default()
	cmds[profileGetCmd.FullCommand()] = func() error {
		return ProfileGet()
	}

	// profile set
	profileSetCmd := profileCmd.Command("set", "Sets the profile name and avatar of the peer")
	profileSetName := profileSetCmd.Flag("name", "Set the peer's display name").Short('n').String()
	profileSetAvatar := profileSetCmd.Flag("avatar", "Set the peer's avatar from an image path (JPEG, PNG, or GIF)").Short('a').String()
	cmds[profileSetCmd.FullCommand()] = func() error {
		return ProfileSet(*profileSetName, *profileSetAvatar)
	}

	// ================================

	// observe
	observeCmd := appCmd.Command("observe", "Observe updates in a thread or all threads. An update is generated when a new block is added to a thread.").Alias("subscribe").Alias("listen").Alias("stream")
	observeThreadID := observeCmd.Arg("thread", "Thread ID, omit for all").String()
	observeType := observeCmd.Flag("type", "Only be alerted to specific type of updates, possible values: merge, ignore, flag, join, announce, leave, text, files comment, like. Can be used multiple times, e.g., --type files --type comment").Short('k').Strings()
	cmds[observeCmd.FullCommand()] = func() error {
		return ObserveCommand(*observeThreadID, *observeType)
	}

	// ================================

	// summary
	summaryCmd := appCmd.Command("summary", "Get a summary of the local node's data")
	cmds[summaryCmd.FullCommand()] = func() error {
		return Summary()
	}

	// ================================
	// @todo this documentation should be moved to docs.textile.io

	// thread
	threadCmd := appCmd.Command("thread", `Threads are distributed sets of encrypted files, often shared between peers, governed by schemas.
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

	// thread add
	threadAddCmd := threadCmd.Command("add", "Adds and joins a new thread")
	threadAddName := threadAddCmd.Arg("name", "The name to use for the new thread").Required().String()
	threadAddKey := threadAddCmd.Flag("key", "A locally unique key used by an app to identify this thread on recovery").Short('k').String()
	threadAddType := threadAddCmd.Flag("type", "Set the thread type to one of: private, read_only, public, open").Short('t').Default("private").String()
	threadAddSharing := threadAddCmd.Flag("sharing", "Set the thread sharing style to one of: not_shared, invite_only, shared").Short('s').Default("not_shared").String()
	threadAddWhitelist := threadAddCmd.Flag("whitelist", "A contact address. When supplied, the thread will not allow additional peers, useful for 1-1 chat/file sharing. Can be used multiple times to include multiple contacts").Short('w').Strings()
	threadAddSchema := threadAddCmd.Flag("schema", "Thread schema ID. Supersedes schema filename").String()
	threadAddSchemaFile := threadAddCmd.Flag("schema-file", "Thread schema filename, supersedes the built-in schema flags").String() // @note could be swapped to .File() perhaps
	threadAddBlob := threadAddCmd.Flag("blob", "Use the built-in blob schema for generic data").Bool()
	threadAddCameraRoll := threadAddCmd.Flag("camera-roll", "Use the built-in camera roll schema").Bool()
	threadAddMedia := threadAddCmd.Flag("media", "Use the built-in media schema").Bool()
	cmds[threadAddCmd.FullCommand()] = func() error {
		return ThreadAdd(*threadAddName, *threadAddKey, *threadAddType, *threadAddSharing, *threadAddWhitelist, *threadAddSchema, *threadAddSchemaFile, *threadAddBlob, *threadAddCameraRoll, *threadAddMedia)
	}

	// thread list
	threadListCmd := threadCmd.Command("list", "Lists info on all threads").Alias("ls").Default()
	cmds[threadListCmd.FullCommand()] = func() error {
		return ThreadList()
	}

	// thread get
	threadGetCmd := threadCmd.Command("get", "Gets and displays info about a thread")
	threadGetThreadID := threadGetCmd.Arg("thread", "Thread ID").Required().String()
	cmds[threadGetCmd.FullCommand()] = func() error {
		return ThreadGet(*threadGetThreadID)
	}

	// thread peer
	threadPeerCmd := threadCmd.Command("peer", "Lists all peers in a thread").Alias("peers")
	threadPeerThreadID := threadPeerCmd.Arg("thread", "Thread ID").Required().String()
	cmds[threadPeerCmd.FullCommand()] = func() error {
		return ThreadPeer(*threadPeerThreadID)
	}

	// thread rename
	threadRenameCmd := threadCmd.Command("rename", "Renames a thread. Only the initiator of a thread can rename it.").Alias("mv")
	threadRenameThreadID := threadRenameCmd.Arg("thread", "Thread ID").Required().String()
	threadRenameName := threadRenameCmd.Arg("name", "The name to rename the thread to").Required().String()
	cmds[threadRenameCmd.FullCommand()] = func() error {
		return ThreadRename(*threadRenameName, *threadRenameThreadID)
	}

	// thread abandon
	threadAbandonCmd := threadCmd.Command("abandon", "Abandon a thread. If no one is else remains participating, the thread dissipates.").Alias("unsubscribe").Alias("leave").Alias("remove").Alias("rm")
	threadAbandonThreadID := threadAbandonCmd.Arg("thread", "Thread ID").Required().String()
	cmds[threadAbandonCmd.FullCommand()] = func() error {
		return ThreadAbandon(*threadAbandonThreadID)
	}

	// thread snapshot
	// A snapshot is an encrypted object containing thread metadata and the latest block hash, which is enough to recover the thread.
	threadSnapshotCmd := threadCmd.Command("snapshot", "Manage thread snapshots").Alias("snapshots")

	// thread snapshot create
	threadSnapshotCreateCmd := threadSnapshotCmd.Command("create", "Snapshots all threads and pushes to registered cafes").Alias("make")
	cmds[threadSnapshotCreateCmd.FullCommand()] = func() error {
		return ThreadSnapshotCreate()
	}

	// thread snapshot search
	threadSnapshotSearchCmd := threadSnapshotCmd.Command("search", "Searches the network for thread snapshots").Alias("find").Default()
	threadSnapshotSearchWait := threadSnapshotSearchCmd.Flag("wait", "Stops searching after [wait] seconds have elapse (max 30s)").Short('w').Default("2").Int()
	cmds[threadSnapshotSearchCmd.FullCommand()] = func() error {
		return ThreadSnapshotSearch(*threadSnapshotSearchWait)
	}

	// thread snapshot apply
	threadSnapshotApplyCmd := threadSnapshotCmd.Command("apply", "Applies a single thread snapshot")
	threadSnapshotApplyID := threadSnapshotApplyCmd.Arg("snapshot", "The ID of the snapshot to apply").Required().String()
	threadSnapshotApplyWait := threadSnapshotApplyCmd.Flag("wait", "Stops searching after [wait] seconds have elapse (max 30s)").Short('w').Default("2").Int()
	cmds[threadSnapshotApplyCmd.FullCommand()] = func() error {
		return ThreadSnapshotApply(*threadSnapshotApplyID, *threadSnapshotApplyWait)
	}

	// thread file
	threadFilesCommand(cmds, threadCmd, []string{"files", "file"})

	// thread block
	threadBlocksCommand(cmds, threadCmd, []string{"blocks", "block"})

	// ================================
	// Tokens are local-only,
	// essentially passwords to get a session token (JWT) to a cafe.

	// token
	tokenCmd := appCmd.Command("token", "Tokens allow other peers to register with a cafe peer").Alias("tokens")

	// token create
	tokenCreateCmd := tokenCmd.Command("add", `Generates an access token (44 random bytes) and saves a bcrypt hashed version for future lookup.
The response contains a base58 encoded version of the random bytes token.`).Alias("create").Alias("generate").Alias("init")
	tokenCreateNoStore := tokenCreateCmd.Flag("no-store", "If used instead of token, the token is generated but not stored in the local cafe database").Short('n').Bool()
	// @todo at some point remove --no-store, if no one is using it
	tokenCreateToken := tokenCreateCmd.Flag("token", "If used instead of no-store, use this existing token rather than creating a new one").Short('t').String()
	cmds[tokenCreateCmd.FullCommand()] = func() error {
		return TokenCreate(*tokenCreateToken, *tokenCreateNoStore)
	}

	// token list
	tokenListCmd := tokenCmd.Command("list", "List info about all stored cafe tokens").Alias("ls").Default()
	cmds[tokenListCmd.FullCommand()] = func() error {
		return TokenList()
	}

	// token validate
	tokenValidateCmd := tokenCmd.Command("validate", "Check validity of existing cafe access token").Alias("valid")
	tokenValidateToken := tokenValidateCmd.Arg("token", "The token to validate").Required().String()
	cmds[tokenValidateCmd.FullCommand()] = func() error {
		return TokenValidate(*tokenValidateToken)
	}

	// token delete
	tokenDeleteCmd := tokenCmd.Command("delete", "Removes an existing cafe token").Alias("del").Alias("remove").Alias("rm")
	tokenDeleteToken := tokenDeleteCmd.Arg("token", "The token to delete").Required().String()
	cmds[tokenDeleteCmd.FullCommand()] = func() error {
		return TokenRemove(*tokenDeleteToken)
	}

	// ================================

	// version
	versionCmd := appCmd.Command("version", "Print the current version and exit")
	versionGit := versionCmd.Flag("git", "Show full git version summary").Short('g').Bool()
	cmds[versionCmd.FullCommand()] = func() error {
		return Version(*versionGit)
	}

	// ================================

	// wallet
	walletCmd := appCmd.Command("wallet", "Create a new wallet, or list its available accounts").Alias("wallets")

	// wallet init
	walletInitCmd := walletCmd.Command("create", "Generate a hierarchical deterministic wallet and output the first child account. A wallet is a seed that deterministically generates child accounts. Child accounts are used to interact with textile. Formula: Autogenerated Mnemonic + Optionally Specified Passphrase = Generated Seed. The seed, mnemonic, and passphrase must be kept top secret. The mnemonic and passphrase must be remembered by you.").Alias("init").Alias("generate")
	walletInitPassphrase := walletInitCmd.Arg("passphrase", "If provided, the resultant wallet seed will be salted with this passphrase, resulting in a different (and more unique) wallet seed than if just the mnemonic was used.").String()
	walletInitWords := walletInitCmd.Flag("words", "How many words to use for the autogenerated mnemonic? 12, 15, 18, 21, 24").Short('w').Default("12").Int()
	cmds[walletInitCmd.FullCommand()] = func() error {
		return WalletInit(*walletInitWords, *walletInitPassphrase)
	}

	// wallet accounts
	walletAccountsCmd := walletCmd.Command("accounts", "List the available accounts (within a specific range) within the wallet's deterministic bounds. Formula: Account = Account Index + Parent Private Key from Parent Seed. Parent Seed = Wallet Mnemonic + Passphrase.").Alias("account")
	walletAccountsMnemonic := walletAccountsCmd.Arg("mnemonic", "The autogenerated mnemonic of the wallet").Required().String()
	walletAccountsPassphrase := walletAccountsCmd.Arg("passphrase", "If the wallet was generated with a passphrase, specify it here to ensure the accounts you are listing are for the same wallet").String()
	walletAccountsDepth := walletAccountsCmd.Flag("depth", "Number of accounts to show").Short('d').Default("1").Int()
	walletAccountsOffset := walletAccountsCmd.Flag("offset", "Account depth to start from").Short('o').Default("0").Int()
	cmds[walletAccountsCmd.FullCommand()] = func() error {
		return WalletAccounts(*walletAccountsMnemonic, *walletAccountsPassphrase, *walletAccountsDepth, *walletAccountsOffset)
	}

	// ================================

	hideGlobalsFlagsFor(
		daemonCmd,
		initCmd,
		walletCmd,
	)

	// commands
	cmd := kingpin.MustParse(appCmd.Parse(os.Args[1:]))
	for key, value := range cmds {
		if key == cmd {
			return value()
		}
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

	req.SetBasicAuth(*appUsername, *appPassword)

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	cancel := func() {
		tr.CancelRequest(req)
	}

	if res != nil && res.StatusCode == 401 {
		err = fmt.Errorf("error: unauthorized")
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

func nextPage() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press ctrl+c to exit, press enter for next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
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

func getDefaultRepo() (string, error) {
	// get homedir
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("get homedir failed: %s", err))
	}
	return filepath.Join(home, ".textile", "repo"), nil
}

// Get the full repo path for the user, will create it if missing
func getRepo(repo string, baseRepo string, accountAddress string) (string, error) {
	if len(repo) > 0 {
		return repo, nil
	} else if len(baseRepo) > 0 && len(accountAddress) > 0 {
		return path.Join(baseRepo, accountAddress), nil
	} else if len(baseRepo) == 0 && len(accountAddress) == 0 {
		return getDefaultRepo()
	} else {
		return "", fmt.Errorf("you must specify --base-repo and --account-address flags")
	}
}

func hideGlobalsFlagsFor(cmds ...*kingpin.CmdClause) {

	m := map[string]bool{}
	for _, c := range cmds {
		m[c.Model().Name] = true
	}

	appCmd.PreAction(func(ctx *kingpin.ParseContext) error {
		if ctx.SelectedCommand == nil {
			return nil
		}
		if m[ctx.String()] {
			for _, r := range appCmd.Model().Flags {
				if r.Name != "help" {
					appCmd.GetFlag(r.Name).Hidden()
				}
			}
		}
		return nil
	})
}
