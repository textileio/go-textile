# textile-go

![banner](https://s3.amazonaws.com/textile.public/Textile_Logo_Horizontal.png)

---

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/textile-go)](https://goreportcard.com/report/github.com/textileio/textile-go) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![CircleCI](https://circleci.com/gh/textileio/textile-go/tree/master.svg?style=shield)](https://circleci.com/gh/textileio/textile-go/tree/master)

## Status

[![Throughput Graph](https://graphs.waffle.io/textileio/textile-go/throughput.svg)](https://waffle.io/textileio/textile-go/metrics/throughput)

## What is Textile?

[Textile](https://www.textile.io) provides encrypted, recoverable, schema-based, and cross-application data storage built on [IPFS](https://github.com/ipfs) and [libp2p](https://github.com/libp2p). We like to think of it as a decentralized Firebase with built-in protocols for sharing and recovery.

This repository contains the core Textile node and daemon, a CLI client, and a mobile client for building an iOS/Android application.

See [textile-mobile](https://github.com/textileio/textile-mobile/) for the [Textile Photos](https://www.textile.photos) iOS/Android app.

## Install

Download the [latest release](https://github.com/textileio/textile-go/releases/latest) for your OS.

## Usage

```
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
```

Textile uses an HD Wallet as an account key manager. You may use the name derived account seed on multiple devices to sync wallet data. To get started, run:

```
$ textile wallet init
```

This will generate a recovery phrase for _all your accounts_. You may specify a word count and password as well (run with `--help` for usage).

Next, use an account from you wallet to initialize a node. First time users should just use _Account 0_, which is printed out by the `wallet init` subcommand. Use the `accounts` subcommand to access deeper derived wallet accounts.

```
$ textile init -s <account_seed>
```

Finally, start the daemon:

```
$ textile daemon
```

TODO: Run through creating a thread, adding images, comments, etc.

## Contributing

```
$ go get github.com/textileio/textile-go
```

You'll need a few different tools here to get setup...

#### Install `dep`

Golang package manager:

```
$ brew install dep
```

#### Install `gx`

IPFS package manager:

```
$ go get -u github.com/whyrusleeping/gx
$ go get -u github.com/whyrusleeping/gx-go
```

#### Install dependencies

Finally, download dependencies managed by `gx` and `dep`:

```
$ make setup
```

#### Run tests

```
$ make test_compile
```

## Building

There are various things to build:

#### CLI/daemon

```
$ make build
```

#### iOS Framework

```
$ go get golang.org/x/mobile/cmd/gomobile
$ gomobile init
$ make ios_framework
```

#### Android Framework

```
$ go get golang.org/x/mobile/cmd/gomobile
$ gomobile init
$ make android_framework
```

## Acknowledgments

While almost entirely different now, this project was jumpstarted from OpenBazaar. Thanks to @cpacia, @drwasho and the rest of the contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go).

And of course, thank you Protocal Labs for the incredible FOSS effort and constant inspiration.

## License

MIT
