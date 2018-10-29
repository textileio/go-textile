# textile-go

![banner](https://s3.amazonaws.com/textile.public/Textile_Logo_Horizontal.png)

---

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/textile-go)](https://goreportcard.com/report/github.com/textileio/textile-go) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![CircleCI](https://circleci.com/gh/textileio/textile-go/tree/master.svg?style=shield)](https://circleci.com/gh/textileio/textile-go/tree/master)

## Status

[![Throughput Graph](https://graphs.waffle.io/textileio/textile-go/throughput.svg)](https://waffle.io/textileio/textile-go/metrics/throughput)

## What is Textile?

Riding on [IPFS](https://github.com/ipfs) and [libp2p](https://github.com/libp2p), [Textile](https://www.textile.io) aims to provide a set of straightforward primitives for building decentralized mobile applications.

This repository currently contains a CLI/daemon, a desktop application, and iOS/Android mobile bindings for running a Textile Photos node. See [textile-mobile](https://github.com/textileio/textile-mobile/) for the [Textile Photos](https://www.textile.photos) iOS/Android app.

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
  daemon   Start a node daemon
  init     Init the node repo and exit
  migrate  Migrate the node repo and exit
  shell    Start a node shell
  version  Print version and exit
  wallet   Manage a wallet of accounts
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

Finally, start the daemon or interactive shell:

```
$ textile daemon|shell
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

#### Install `node`

NodeJS is used for git hooks and some build tooling:

```
$ brew install node
```

#### Install dependencies

Finally, download deps managed by `gx` and `dep`:

```
$ npm run setup
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

#### Desktop Application

WARNING: Desktop is an unmaintained experiment. Security issues may exist.

The build is made by a vendored version of `go-astilectron-bundler`. Due to Go's painful package management, you'll want to delete any `go-astilectron`-related binaries and source code you have installed from `github.com/asticode` in your `$GOPATH`. Then you can install the vendored `go-astilectron-bundler`:

```
$ go install ./vendor/github.com/asticode/go-astilectron-bundler/astilectron-bundler
```

Pick your OS: Linux, Darwin, or Windows:

```
$ cd desktop
$ astilectron-bundler -v -c bundler.linux.json
$ astilectron-bundler -v -c bundler.darwin.json
$ astilectron-bundler -v -c bundler.windows.json
```

Double-click the built app in `desktop/output`, or run it directly:

```
$ cd desktop && go run *.go
```

See [go-astilectron-bundler](https://github.com/asticode/go-astilectron-bundler) for more build configurations.

Note: Because `cgo` is required, you'll need to setup a proper C toolchain for cross-OS-compilation.

## Commitizen

The easiest way to write a valid commit message is to use the `npm` script:

```
$ npm run cm
```

This will start the interactive commit prompt.

## Acknowledgments

While almost entirely different now, this project was jumpstarted from OpenBazaar. Thanks to @cpacia, @drwasho and the rest of the contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go).

## License

MIT
