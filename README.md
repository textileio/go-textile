# go-textile

[![Made by Textile](https://img.shields.io/badge/made%20by-Textile-informational.svg?style=popout-square)](https://textile.io)
[![Chat on Slack](https://img.shields.io/badge/slack-slack.textile.io-informational.svg?style=popout-square)](https://slack.textile.io)
[![GitHub license](https://img.shields.io/github/license/textileio/photos-desktop.svg?style=popout-square)](./LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/textileio/go-textile?style=flat-square)](https://goreportcard.com/report/github.com/textileio/go-textile?style=flat-square)
[![CircleCI branch](https://img.shields.io/circleci/project/github/textileio/go-textile/master.svg?style=popout-square)](https://circleci.com/gh/textileio/go-textile)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=popout-square)](https://github.com/RichardLitt/standard-readme)

> Textile implementation in Go

This repository contains the core API, daemon, and command-line client, as well as bindings for mobile (iOS/Android) applications.

[Textile](https://www.textile.io) provides encrypted, recoverable, schema-based, and cross-application data storage built on [IPFS](https://github.com/ipfs) and [libp2p](https://github.com/libp2p). We like to think of it as a decentralized data wallet with built-in protocols for sharing and recovery, or more simply, **an open and programmable iCloud**.

**Please see [Textile Docs](https://docs.textile.io/) for more**.

Join us on our [public Slack channel](https://slack.textile.io/) for news, discussions, and status updates. [Check out our blog](https://medium.com/textileio) for the latest posts and announcements.

## Table of Contents

-   [Security](#security)
-   [Background](#background)
-   [Install](#install)
-   [Usage](#usage)
-   [Develop](#develop)
-   [Contribute](#contribute)
-   [Changelog](#changelog)
-   [License](#license)

## Security

Textile is still under heavy development and no part of it should be used before a thorough review of the underlying code and an understanding that APIs and protocols may change rapidly. There may be coding mistakes and the underlying protocols may contain design flaws. Please [let us know](mailto:contact@textile.io) immediately if you have discovered a security vulnerability.

Please also read the [security note](https://github.com/ipfs/go-ipfs#security-issues) for [go-ipfs](https://github.com/ipfs/go-ipfs).

## Background

Textile is a set of tools and trust-less infrastructure for building _censorship resistant_ and _privacy preserving_ applications.

While interoperable with the whole [IPFS](https://ipfs.io/) peer-to-peer network, Textile-flavored peers represent an additional layer or sub-network of **users, applications, and services**.

With good encryption defaults and anonymous, disposable application services like [cafes](https://docs.textile.io/concepts/cafes/), Textile aims to bring the decentralized internet to real products that people love.

[Continue reading](https://docs.textile.io/concepts/) about Textile...

## Install

[Installation instructions](https://docs.textile.io/install/the-daemon/) for the command-line tool and daemon are in [the docs](https://docs.textile.io).

## Usage

The [Tour of Textile](https://docs.textile.io/a-tour-of-textile/) goes through many examples and use cases. `textile --help-long` provides a quick look at the available APIs:

```
$ textile --help-long
usage: textile [<flags>] <command> [<args> ...]

Textile is a set of tools and trust-less infrastructure for building censorship resistant and privacy preserving
applications

Flags:
  --help              Show context-sensitive help (also try --help-long and --help-man).
  --api="http://127.0.0.1:40600"
                      API Address to use
  --api-version="v0"  API version to use
  --debug             Set the logging level to debug

Commands:
  help [<command>...]
    Show help.


  account get
    Shows the local peer's account info as a contact


  account seed
    Shows the local peer's account seed


  account address
    Shows the local peer's account address


  account sync [<flags>]
    Syncs the local account peer with other peers found on the network

    --wait=2  Stops searching after 'wait' seconds have elapsed (max 30s)

  block list [<flags>]
    Paginates blocks in a thread

    -t, --thread="default"  Thread ID
    -o, --offset=OFFSET     Offset ID to start listing from
    -l, --limit=5           List page size
    -d, --dots              Return GraphViz dots instead of JSON

  block meta <block>
    Get the metadata for a block


  block ignore <block>
    Remove a block by marking it to be ignored


  block file [<flags>] <files-block>
    Get the files, or a specific file, of a Files Block

    --index=0    If provided, the index of a specific file to retrieve
    --path=PATH  If provided, the path of a specific file to retrieve
    --content    If provided alongside a path, the content of the specific file is retrieved

  cafe add --token=TOKEN <peer>
    Registers with a cafe and saves an expiring service session token. An access token is required to register, and
    should be obtained separately from the target cafe.

    -t, --token=TOKEN  An access token supplied by the cafe

  cafe list
    List info about all active cafe sessions


  cafe get <cafe>
    Gets and displays info about a cafe session


  cafe delete <cafe>
    Deregisters a cafe (content will expire based on the cafe's service rules)


  cafe messages
    Check for messages at all cafes. New messages are downloaded and processed opportunistically.


  chat [<flags>]
    Starts an interactive chat session in a thread

    -t, --thread="default"  Thread ID

  comment add <block> <body>
    Attach a comment to a block


  comment list <block>
    Get the comments that are attached to a block


  comment get <comment-block>
    Get a comment by its own Block ID


  comment ignore <comment-block>
    Ignore a comment by its own Block ID


  config [<name>] [<value>]
    Get or set configuration variables


  contact add [<flags>]
    Adds a contact by display name or account address to known contacts

    -n, --name=NAME        Add by display name
    -a, --address=ADDRESS  Add by account address
        --wait=WAIT        Stops searching after [wait] seconds have elapsed

  contact list
    List known contacts


  contact get <address>
    Gets a known contact


  contact delete <address>
    Deletes a known contact


  contact search [<flags>]
    Searches locally and on the network for contacts

    -n, --name=NAME        Search by display name
    -a, --address=ADDRESS  Search by account address
        --only-local       Only search local contacts
        --only-remote      Only search remote contacts
        --limit=5          Stops searching after [limit] results are found
        --wait=2           Stops searching after [wait] seconds have elapsed (max 30s)

  daemon [<flags>]
    Start a node daemon session

    -r, --repo-dir=REPO-DIR  Specify a custom repository path
    -p, --pin-code=PIN-CODE  Specify the pin code for datastore encryption (omit no pin code was used during init)
    -s, --serve-docs         Whether to serve the local REST API docs

  docs
    Prints the CLI help as HTML


  feed [<flags>]
    Paginates post (join|leave|files|message) and annotation (comment|like) block types as a consumable feed.

    The --mode option dictates how the feed is displayed:

    - "chrono": All feed block types are shown. Annotations always nest their target post, i.e., the post a comment
    is about. - "annotated": Annotations are nested under post targets, but are not shown in the top-level feed. -
    "stacks": Related blocks are chronologically grouped into "stacks". A new stack is started if an unrelated block

      breaks continuity. This mode is used by Textile Photos.

    Stacks may include:

    - The initial post with some nested annotations. Newer annotations may have already been listed. - One or more
    annotations about a post. The newest annotation assumes the "top" position in the stack. Additional

      annotations are nested under the target. Newer annotations may have already been listed in the case as well.
    -t, --thread=THREAD  Thread ID, omit for all
    -o, --offset=OFFSET  Offset ID to start listening from
    -l, --limit=3        List page size
    -m, --mode="chrono"  Feed mode, one of: chrono, annotated, stacks

  file list [<flags>]
    Paginates thread files

    -t, --thread="default"  Thread ID
    -o, --offset=OFFSET     Offset ID to start listing from
    -l, --limit=5           List page size

  file keys <target-block>
    Shows file keys under the given target


  file add [<flags>] [<path>]
    Adds a file, directory, or hash to a thread. Files not supported by the thread schema are ignored

    -t, --thread="default"  Thread ID
    -c, --caption=CAPTION   File(s) caption
    -g, --group             If provided, group a directory's files together into a single object, includes nested
                            directories
    -v, --verbose           Prints files as they are milled

  file ignore <files-block>
    Ignores a thread file by its own block ID


  file get [<flags>] <hash>
    Get the metadata or content of a specific file

    --content  If provided, the decrypted content of the file is retrieved

  init --seed=SEED [<flags>]
    Initialize the node repository and exit

    -s, --seed=SEED                Account seed (run 'wallet' command to generate new seeds)
    -p, --pin-code=PIN-CODE        Specify a pin code for datastore encryption
    -r, --repo-dir=REPO-DIR        Specify a custom repository path
        --server                   Apply IPFS server profile
        --swarm-ports=SWARM-PORTS  Set the swarm ports (TCP,WS). A random TCP port is chosen by default
        --log-files                If true, writes logs to rolling files, if false, writes logs to stdout
        --api-bind-addr="127.0.0.1:40600"
                                   Set the local API address
        --cafe-bind-addr="0.0.0.0:40601"
                                   Set the cafe REST API address
        --gateway-bind-addr="127.0.0.1:5050"
                                   Set the IPFS gateway address
        --profile-bind-addr="127.0.0.1:6060"
                                   Set the profiling address
        --cafe-open                Open the p2p cafe service for other peers
        --cafe-url=CAFE-URL        Specify a custom URL of this cafe, e.g., https://mycafe.com
        --cafe-neighbor-url=CAFE-NEIGHBOR-URL
                                   Specify the URL of a secondary cafe. Must return cafe info, e.g., via a Gateway:
                                   https://my-gateway.yolo.com/cafe, or a cafe API: https://my-cafe.yolo.com

  invite create [<flags>]
    Creates a direct account-to-account or external invite to a thread

    -t, --thread="default"  Thread ID
    -a, --address=ADDRESS   Account Address, omit to create an external invite
        --wait=2            Stops searching after [wait] seconds have elapsed (max 30s)

  invite list
    Lists all pending thread invites


  invite accept [<flags>] <id>
    Accepts a direct account-to-account or external invite to a thread

    -k, --key=KEY  Key for an external invite

  invite ignore <id>
    Ignores a direct account-to-account invite to a thread


  ipfs peer
    Shows the local node's IPFS peer ID


  ipfs swarm connect [<address>]
    Opens a new direct connection to a peer address


  ipfs swarm peers [<flags>]
    Lists the set of peers this node is connected to

    -v, --verbose    Display all extra information
    -s, --streams    Also list information about open streams for search peer
    -l, --latency    Also list information about the latency to each peer
    -d, --direction  Also list information about the direction of connection

  ipfs cat [<flags>] <hash>
    Displays the data behind an IPFS CID (hash)

    -k, --key=KEY  Encryption key

  like add <block>
    Attach a like to a block


  like list <block>
    Get likes that are attached to a block


  like get <like-block>
    Get a like by its own Block ID


  like ignore <like-block>
    Ignore a like by its own Block ID


  log [<flags>]
    List or change the verbosity of one or all subsystems log output. Textile logs piggyback on the IPFS event logs.

    -s, --subsystem=SUBSYSTEM  The subsystem logging identifier, omit for all
    -l, --level=LEVEL          One of: debug, info, warning, error, critical. Omit to get current level.
    -t, --textile-only         Whether to list/change only Textile subsystems, or all available subsystems

  message add [<flags>] [<body>]
    Adds a message to a thread

    -t, --thread="default"  Thread ID

  message list [<flags>]
    Paginates thread messages

    -t, --thread=THREAD  Thread ID, omit to paginate all messages
    -o, --offset=OFFSET  Offset ID to start the listing from
    -l, --limit=10       List page size

  message get [<message-block>]
    Gets a message by its own Block ID


  message ignore [<message-block>]
    Ignores a message by its own Block ID


  migrate [<flags>]
    Migrate the node repository and exit

    -r, --repo-dir=REPO-DIR  Specify a custom repository path
    -p, --pin-code=PIN-CODE  Specify the pin code for datastore encryption (omit of none was used during init)

  notification list
    Lists all notifications


  notification read <id>
    Marks a notification as read


  ping <address>
    Pings another peer on the network, returning [online] or [offline]


  profile get
    Gets the local peer profile


  profile set name <value>
    Sets the profile name of the peer


  profile set avatar <value>
    Sets the profile avatar of the peer


  subscribe [<flags>]
    Subscribes to updates in a thread or all threads. An update is generated when a new block is added to a thread.

    -t, --thread=THREAD  Thread ID, omit for all
    -k, --type=TYPE ...  Only be alerted to specific type of updates, possible values: merge, ignore, flag, join,
                         announce, leave, text, files comment, like. Can be used multiple times, e.g., --type files
                         --type comment

  summary
    Get a summary of the local node's data


  thread add [<flags>] <name>
    Adds and joins a new thread

    -k, --key=KEY                  A locally unique key used by an app to identify this thread on recovery
    -t, --type="private"           Set the thread type to one of: private, read_only, public, open
    -s, --sharing="not_shared"     Set the thread sharing style to one of: not_shared, invite_only, shared
    -w, --whitelist=WHITELIST ...  A contact address. When supplied, the thread will not allow additional peers,
                                   useful for 1-1 chat/file sharing. Can be used multiple times to include multiple
                                   contacts
        --schema=SCHEMA            Thread schema ID. Supersedes schema filename
        --schema-file=SCHEMA-FILE  Thread schema filename, supersedes the built-in schema flags
        --blob                     Use the built-in blob schema for generic data
        --camera-roll              Use the built-in camera roll schema
        --media                    Use the built-in media schema

  thread list
    Lists info on all threads


  thread get <thread>
    Gets and displays info about a thread


  thread default
    Gets and displays info about the default thread (if selected


  thread peer [<flags>]
    Lists all peers in a thread

    -t, --thread="default"  Thread ID

  thread rename [<flags>] <name>
    Renames a thread. Only the initiator of a thread can rename it.

    -t, --thread="default"  Thread ID

  thread unsubscribe <thread>
    Unsubscribes from the thread, and if no one else remains subscribed, deletes it


  thread snapshot create
    Snapshots all threads and pushes to registered cafes


  thread snapshot search [<flags>]
    Searches the network for thread snapshots

    -w, --wait=2  Stops searching after [wait] seconds have elapse (max 30s)

  thread snapshot apply [<flags>] <snapshot>
    Applies a single thread snapshot

    -w, --wait=2  Stops searching after [wait] seconds have elapse (max 30s)

  token create [<flags>]
    Generates an access token (44 random bytes) and saves a bcrypt hashed version for future lookup. The response
    contains a base58 encoded version of the random bytes token.

    -n, --no-store     If used instead of token, the token is generated but not stored in the local cafe database
    -t, --token=TOKEN  If used instead of no-store, use this existing token rather than creating a new one

  token list
    List info about all stored cafe tokens


  token validate <token>
    Check validity of existing cafe access token


  token delete <token>
    Removes an existing cafe token


  version [<flags>]
    Print the current version and exit

    -g, --git  Show full git version summary

  wallet init [<flags>]
    Initializes a new account wallet backed by a mnemonic phrase

    -w, --word-count=12      Number of mnemonic phrase words: 12,15,18,21,24
    -p, --password=PASSWORD  Mnemonic password (omit if none)

  wallet accounts [<flags>] <mnemonic>
    Shows the derived accounts (address/seed pairs) in a wallet

    -p, --password=PASSWORD  Mnemonic password (omit if none)
    -d, --depth=1            Number of accounts to show
    -o, --offset=0           Account depth to start from
```

## Develop

    git clone git@github.com:textileio/go-textile.git

### Requirements

-   go >= 1.12
-   node >= 10.0

Extra setup steps are needed to build the bindings for iOS or Android, as `gomobile` does not yet support [go modules](https://github.com/golang/go/wiki/Modules). You'll need to **move the go-textile source** into your `GOPATH` (like pre-go1.11 development), before installing and initializing the `gomobile` tools:

    go get golang.org/x/mobile/cmd/gomobile
    gomobile init

Now you can execute the iOS and Android build tasks below. For the other build tasks, the source must _not_ be under `GOPATH`. Go 1.13 is supposed to bring module support to `gomobile`, at which point we can remove this madness!

### Install dependencies:

    make setup

### Build `textile`:

    make build

### Run unit tests:

    make test

### Build the iOS framework:

    make ios

### Build the Android Archive Library (aar):

    make android

### Build the swagger docs:

    make docs

## Contributing

This project is a work in progress. As such, there's a few things you can do right now to help out:

-   **Ask questions**! We'll try to help. Be sure to drop a note (on the above issue) if there is anything you'd like to work on and we'll update the issue to let others know. Also [get in touch](https://slack.textile.io) on Slack.
-   **Open issues**, [file issues](https://github.com/textileio/go-textile/issues), submit pull requests!
-   **Perform code reviews**. More eyes will help a) speed the project along b) ensure quality and c) reduce possible future bugs.
-   **Take a look at the code**. Contributions here that would be most helpful are **top-level comments** about how it should look based on your understanding. Again, the more eyes the better.
-   **Add tests**. There can never be enough tests.

Before you get started, be sure to read our [contributors guide](./CONTRIBUTING.md) and our [contributor covenant code of conduct](./CODE_OF_CONDUCT.md).

## Changelog

[Changelog is published to Releases.](https://github.com/textileio/go-textile/releases)

## License

[MIT](LICENSE)
