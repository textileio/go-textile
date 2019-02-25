# textile-go

![banner](https://s3.amazonaws.com/textile.public/Textile_Logo_Horizontal.png)

---

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/textile-go)](https://goreportcard.com/report/github.com/textileio/textile-go) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![CircleCI](https://circleci.com/gh/textileio/textile-go/tree/master.svg?style=shield)](https://circleci.com/gh/textileio/textile-go/tree/master)

## Status

[![Throughput Graph](https://graphs.waffle.io/textileio/textile-go/throughput.svg)](https://waffle.io/textileio/textile-go/metrics/throughput)

This repository contains the core Textile node and daemon, a command-line client, and a mobile client for building an iOS/Android application.

See [textile-mobile](https://github.com/textileio/textile-mobile/) for the [Textile Photos](https://www.textile.photos) iOS/Android app.

## What is Textile?

[Textile](https://www.textile.io) provides encrypted, recoverable, schema-based, and cross-application data storage built on [IPFS](https://github.com/ipfs) and [libp2p](https://github.com/libp2p). We like to think of it as a decentralized data wallet with built-in protocols for sharing and recovery, or more simply, **an open and programmable iCloud**.

**Please see the [Wiki](https://github.com/textileio/textile-go/wiki) for more**.

## Install

Download the [latest release](https://github.com/textileio/textile-go/releases/latest) for your OS.

## Usage

    ~ $ textile --help
    Usage:
      textile [OPTIONS] <command>
    
    Help Options:
      -h, --help  Show this help message
    
    Available commands:
      add            Add file(s) to a thread
      address        Show wallet address
      blocks         View thread blocks
      cafes          Manage cafes
      chat           Start a thread chat
      comments       Manage thread comments
      daemon         Start the daemon
      get            Get a thread file by ID
      ignore         Ignore a thread file
      init           Init the node repo and exit
      invites        Manage thread invites
      keys           Show file keys
      likes          Manage thread likes
      ls             Paginate thread files
      messages       Manage thread messages
      migrate        Migrate the node repo and exit
      notifications  Manage notifications
      peer           Show peer ID
      ping           Ping another peer
      profile        Manage public profile
      sub            Subscribe to thread updates
      swarm          Access IPFS swarm commands
      threads        Manage threads
      version        Print version and exit
      wallet         Manage or create an account wallet

## Quick-start

#### Initialize a new wallet.

    $ textile wallet init

This will generate a mnemonic phrase for accessing / recovering derived accounts. You may specify a word count and password as well (run with `--help` for usage).

#### Initialize a peer with an account.

Next, use an account seed from your wallet to initialize a new peer. First time users should just use the first account’s (Account 0) seed, which is printed out by the `wallet init` sub-command. The private seed begins with “S”. The public address begins with “P”. Use the `accounts` sub-command to access deeper derived wallet accounts.

    $ textile init -s <account_seed>

#### Start the daemon.

    $ textile daemon

You can now use the command-line client to interact with your running peer.

## Adding Files

Files are tracked by threads. So, let’s start there.

#### Create a new thread.

    $ textile threads add "hello world" --photos

This will create and join a thread backed by the built-in photos schema. Use the `--help` flag on any sub-command for more options and info.

#### Add a file to the thread.

    $ textile add <image path> --caption "beautiful"

The thread schema encodes the image at various width and extracts exif data. The resulting files are added to the thread under one directory. You also add an entire directory.

    $ textile add <dir path> --caption "more beauty"

#### Browse files.

The command-line client is not really meant to provide a great UX for browsing account files. However, you can easily paginate through them with `ls`.
Note: A file’s ID is just its block (update) ID.

    $ textile ls --thread <thread ID>

#### Comment on a file.

    $ textile comments add "good eye" --block <block ID>

#### Like a file.

    $ textile likes add --block <block ID>

## Sharing files / chatting

In order to start sharing or chatting with someone else, you’ll first need an open thread. Open threads allow invites to other peers.

    $ textile threads add "dog photos" --photos --open

Again, we used the built-in photos schema, but this time we’ve opened the thread to invites. Invites allow other peers to join threads. There are two types of invites: direct peer-to-peer and external.

- Peer-to-peer invites are encrypted with the invitee's public key.
- External invites are encrypted with a single-use key and are useful for on-boarding new users. Once an external invite and its key are shared, you should considered it public, since any number of peers can use it to join.

#### Create a direct peer-to-peer thread invite.

    $ textile invites create --thread <thread ID> --peer <peer ID>

The receiving peer will be notified of the invite. They can list all pending direct invites.

    $ textile invites ls

The result is something like:

    [
        {
            "id": "QmUv8783yptknBHCSSnscWNLZdz5K8uhpHZYaWnPkMxu4i",
            "name": "dog photos",
            "inviter": "fido",
            "date": "2018-12-07T13:02:57-08:00"
        }
    ]

#### Accept a direct peer-to-peer invite.

    $ textile invites accept QmUv8783yptknBHCSSnscWNLZdz5K8uhpHZYaWnPkMxu4i

#### Create an “external” thread invite.

This is done by simply omitting the `--peer` flag with the `invites create` command.

    $ textile invites create --thread <thread ID>

The result is something like:

    {
        "invite": "QmcDmpmBr6qB5QGvsUaTZZtwpGpevGgiSEa7C3AJE9EZiU",
        "key": "aKrQmYCMiCQvkyjnm4sFhxdZaFH8g9h7EaLxdBGsZCVjsoyMPzQJQUyPrn7G"
    }

Your friend can use the resulting address and key to accept the invite and join the thread.

    $ textile invites accept QmcDmpmBr6qB5QGvsUaTZZtwpGpevGgiSEa7C3AJE9EZiU --key aKrQmYCMiCQvkyjnm4sFhxdZaFH8g9h7EaLxdBGsZCVjsoyMPzQJQUyPrn7G

At this point, both of you can add and receive files via this thread. You can also exchange plain text messages.

#### Add a text message to a thread.

    $ textile messages add "nice photos" --thread <thread ID>

#### Start a chat in a thread.

    $ textile chat --thread <thread ID>

This will start an interactive chat session with other thread peers.

## Building File Schemas
To-do.

## Using a Cafe
To-do.

## Hosting a Cafe
To-do.

## Contributing

#### Go get the source code.

    $ go get github.com/textileio/textile-go

You can ignore the `gx` package errors. You'll need two package managers to get setup…

#### Install the golang package manager, `dep`.

    $ brew install dep

#### Install the IPFS package manager, `gx`.

    $ go get -u github.com/whyrusleeping/gx
    $ go get -u github.com/whyrusleeping/gx-go

#### Install the dependencies managed by `dep` and `gx`.

    $ go get github.com/ahmetb/govvv
    $ cd $GOPATH/src/github.com/textileio/textile-go
    $ make setup

#### Run the tests.

    $ make test_compile

## Building

There are various things to build…

#### CLI/daemon

    $ make build

#### iOS Framework

    $ go get golang.org/x/mobile/cmd/gomobile
    $ gomobile init
    $ make ios_framework

#### Android Framework

    $ go get golang.org/x/mobile/cmd/gomobile
    $ gomobile init
    $ make android_framework

#### Docs

    $ go get -u github.com/swaggo/swag/cmd/swag
    $ make build_docs

## Acknowledgments

While now almost entirely different, this project was jump-started from [OpenBazaar](https://openbazaar.org/). Thanks to @cpacia, @drwasho and the rest of the contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go).
And of course, thank you, [Protocal Labs](https://protocol.ai/), for the incredible FOSS effort and constant inspiration.

## License

MIT
