# go-textile

![banner](https://s3.amazonaws.com/textile.public/Textile_Logo_Horizontal.png)

---

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/go-textile)](https://goreportcard.com/report/github.com/textileio/go-textile) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![CircleCI](https://circleci.com/gh/textileio/go-textile/tree/master.svg?style=shield)](https://circleci.com/gh/textileio/go-textile/tree/master)

## Status

[![Throughput Graph](https://graphs.waffle.io/textileio/go-textile/throughput.svg)](https://waffle.io/textileio/go-textile/metrics/throughput)

This repository contains the core Textile node and daemon, a command-line client, and a mobile client for building an iOS/Android application.

See [textile-mobile](https://github.com/textileio/textile-mobile/) for the [Textile Photos](https://www.textile.photos) iOS/Android app.

## What is Textile?

[Textile](https://www.textile.io) provides encrypted, recoverable, schema-based, and cross-application data storage built on [IPFS](https://github.com/ipfs) and [libp2p](https://github.com/libp2p). We like to think of it as a decentralized data wallet with built-in protocols for sharing and recovery, or more simply, **an open and programmable iCloud**.

**Please see the [Wiki](https://github.com/textileio/go-textile/wiki) for more**.

## Install

Download the [latest release](https://github.com/textileio/go-textile/releases/latest) for your OS or jump to [Docker](https://github.com/textileio/go-textile#docker). You can also install the Textile [desktop tray app](https://github.com/textileio/go-textile/releases/latest) to run local web/desktop apps that leverage Textile tools.

## Usage

    ~ $ textile --help
    Usage:
      textile [OPTIONS] <command>

    Help Options:
      -h, --help  Show this help message

    Available commands:
      account        Manage a wallet account
      blocks         View thread blocks
      cafes          Manage cafes
      chat           Start a thread chat
      commands       List available commands
      comments       Manage thread comments
      config         Get and set config values
      contacts       Manage contacts
      daemon         Start the daemon
      files          Manage thread files
      init           Init the node repo and exit
      invites        Manage thread invites
      ipfs           Access IPFS commands
      likes          Manage thread likes
      logs           List and control Textile subsystem logs.
      ls             Paginate thread content
      messages       Manage thread messages
      migrate        Migrate the node repo and exit
      notifications  Manage notifications
      ping           Ping another peer
      profile        Manage public profile
      sub            Subscribe to thread updates
      threads        Manage threads
      tokens         Manage Cafe access tokens
      version        Print version and exit
      wallet         Manage or create an account wallet

## Quick-start

#### Initialize a new wallet

    $ textile wallet init

This will generate a mnemonic phrase for accessing/recovering derived accounts. You may specify a word count and password as well (run with `--help` for usage).

#### Initialize a peer with an account

Next, use an account seed from your wallet to initialize a new peer. First time users should just use the first account’s (Account 0) seed, which is printed out by the `wallet init` sub-command. The private seed begins with “S”. The public address begins with “P”. Use the `accounts` sub-command to access deeper derived wallet accounts.

    $ textile init -s <account_seed>

#### Start the daemon

    $ textile daemon

You can now use the command-line client to interact with your running peer.

## Adding Files

Files are tracked by threads. So, let’s start there.

#### Create a new thread

    $ textile threads add "hello world" --media

This will create and join a thread backed by the built-in media schema. Use the `--help` flag on any sub-command for more options and info.

#### Add a file

    $ textile files add <image path> --caption "beautiful"

The thread schema encodes the image at various width and extracts exif data. The resulting files are added to the thread under one directory. You also add an entire directory.

    $ textile files add <dir path> --caption "more beauty"

#### Browse a thread feed

The command-line client is not really meant to provide a great UX for browsing thread content. However, you can easily paginate the feed with `ls`.

    $ textile ls --thread <thread ID>

#### Comment on a file

    $ textile comments add "good eye" --block <block ID>

#### Like a file

    $ textile likes add --block <block ID>

## Sharing files / chatting

In order to start sharing or chatting with someone else, you’ll first need an open and shared thread. An `open` threads allows other to read and write, while `shared` means anyone can join via an invite. See `textile threads --help` for much more about threads, access control types, and share settings.

    $ textile threads add "dog photos" --media --type=open --sharing=shared

There are two types of invites: direct peer-to-peer and external.

- Peer-to-peer invites are encrypted with the invitee's public key.
- External invites are encrypted with a single-use key and are useful for on-boarding new users.

#### Create a direct peer-to-peer thread invite

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

#### Accept a direct peer-to-peer invite

    $ textile invites accept QmUv8783yptknBHCSSnscWNLZdz5K8uhpHZYaWnPkMxu4i

#### Create an “external” thread invite

This is done by simply omitting the `--peer` flag with the `invites create` command.

    $ textile invites create --thread <thread ID>

The result is something like:

    {
        "invite": "QmcDmpmBr6qB5QGvsUaTZZtwpGpevGgiSEa7C3AJE9EZiU",
        "key": "aKrQmYCMiCQvkyjnm4sFhxdZaFH8g9h7EaLxdBGsZCVjsoyMPzQJQUyPrn7G"
    }

Your friend can use the resulting address and key to accept the invite and join the thread.

    $ textile invites accept QmcDmpmBr6qB5QGvsUaTZZtwpGpevGgiSEa7C3AJE9EZiU --key aKrQmYCMiCQvkyjnm4sFhxdZaFH8g9h7EaLxdBGsZCVjsoyMPzQJQUyPrn7G

At this point, both of you can add and receive files via this thread. You can also exchange text messages (chat).

#### Add a text message to a thread

    $ textile messages add "nice photos" --thread <thread ID>

#### Start a chat in a thread

    $ textile chat --thread <thread ID>

This will start an interactive chat session with other thread peers.

## Docker

See available tags [here](https://hub.docker.com/r/textile/go-textile/tags).

#### Run a Textile node

    $ docker run -it --name textile-node \
      -p 4001:4001 -p 8081:8081 -p 5050:5050 -p 127.0.0.1:40600:40600 \
      textile/go-textile:latest

#### Run a Textile node as a _cafe_

    $ docker run -it --name textile-cafe-node \
      -p 4001:4001 -p 8081:8081 -p 5050:5050 -p 127.0.0.1:40600:40600 -p 40601:40601 \
      -e CAFE_HOST_URL=<public_URL> -e CAFE_HOST_PUBLIC_IP=<public_IP> \
      textile/go-textile:latest-cafe

A cafe node can issue client sessions (JWTs) to other nodes. In order to issue valid sessions, the cafe must know its public IP address and the machine's public facing URL. The `CAFE_HOST_PUBLIC_IP` and `CAFE_HOST_URL` environment variable values are written to the textile config file. Read more about cafe host config settings [here](https://github.com/textileio/go-textile/wiki/Config-File#cafe).

## Contributing

#### Go get the source code

    $ go get github.com/textileio/go-textile

You can ignore the `gx` package errors. You'll need two package managers to get setup…

#### Install the golang package manager, `dep`

MacOS:

    $ brew install dep
    
Debian:

    $ sudo apt-get install go-dep

#### Install the IPFS package manager, `gx`

    $ go get -u github.com/whyrusleeping/gx
    $ go get -u github.com/whyrusleeping/gx-go

#### Install the dependencies managed by `dep` and `gx`

    $ cd $GOPATH/src/github.com/textileio/go-textile
    $ make setup

#### Run the tests

    $ make test

## Building

There are various things to build…

#### CLI/daemon

    $ make build

#### iOS Framework

    $ go get golang.org/x/mobile/cmd/gomobile
    $ gomobile init
    $ make ios

#### Android Framework

    $ go get golang.org/x/mobile/cmd/gomobile
    $ gomobile init
    $ make android

#### Docs

    $ make docs

#### Tray app

The build is made by a vendored version of `go-astilectron-bundler`. Due to Go's painful package management, you'll want to delete any `go-astilectron`-related binaries and source code you have installed from `github.com/asticode` in your `$GOPATH`. Then you can install the vendored `go-astilectron-bundler`:

```
go install ./vendor/github.com/asticode/go-astilectron-bundler/astilectron-bundler
```

Change into the `tray` folder and build the app:

```
cd tray
astilectron-bundler -v
```

Double-click the built app in `tray/output/{darwin,linux,windows}-amd64`, or run it directly:

```
go run *.go
```

You can also build the architecture-specific versions with:

```
astilectron-bundler -v -c bundler.{darwin,linux,windows}.json
```

##### Linux

On Linux, you also have to `apt-get install libappindicator1 xclip libgconf-2-4` due to an issue with building Electron-based apps.

## Acknowledgments

While now almost entirely different, this project was jump-started from [OpenBazaar](https://openbazaar.org/). Thanks to @cpacia, @drwasho and the rest of the contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go).
And of course, thank you, [Protocal Labs](https://protocol.ai/), for the incredible FOSS effort and constant inspiration.

## License

MIT
