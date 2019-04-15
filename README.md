# go-textile

![banner](https://s3.amazonaws.com/textile.public/Textile_Logo_Horizontal.png)

---

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/go-textile)](https://goreportcard.com/report/github.com/textileio/go-textile) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![CircleCI](https://circleci.com/gh/textileio/go-textile/tree/master.svg?style=shield)](https://circleci.com/gh/textileio/go-textile/tree/master)

## Status

[![Throughput Graph](https://graphs.waffle.io/textileio/go-textile/throughput.svg)](https://waffle.io/textileio/go-textile/metrics/throughput)

This repository contains a the core Textile daemon + command-line client, as well as bindings for mobile (iOS/Android) applications.

See [textile-mobile](https://github.com/textileio/textile-mobile/) for the [Textile Photos](https://www.textile.photos) iOS/Android app.

## What is Textile?

[Textile](https://www.textile.io) provides encrypted, recoverable, schema-based, and cross-application data storage built on [IPFS](https://github.com/ipfs) and [libp2p](https://github.com/libp2p). We like to think of it as a decentralized data wallet with built-in protocols for sharing and recovery, or more simply, **an open and programmable iCloud**.

**Please see the [docs](https://docs.textile.io/) for more**.

## Install

Download the [latest release](https://github.com/textileio/go-textile/releases/latest) for your OS or jump to [Docker](https://github.com/textileio/go-textile#docker). You can also install the Textile [desktop tray app](https://github.com/textileio/go-textile/releases/latest) to run local web/desktop apps that leverage Textile tools.

## Usage

Check out the [docs site](https://docs.textile.io/) for more detailed usage instructions and tutorials.

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

## Contributing

**Go >= 1.12 is required.**

    git clone git@github.com:textileio/go-textile.git

#### Run the tests

    make test

## Building

There are various things to build… first off, run setup:

    make setup

If you plan on building the bindings for iOS or Android, install and init the `gomobile` tools:

    go get golang.org/x/mobile/cmd/gomobile
    gomobile init

#### CLI/daemon

    make build

#### iOS Framework

    make ios

#### Android Framework

    make android

#### Docs

    make docs

#### Desktop

Install `go-astilectron-bundler`:

    go get github.com/asticode/go-astilectron-bundler/...

Change into the `tray` folder and build the app:

    cd tray
    astilectron-bundler -v

Double-click the built app in `tray/output/{darwin,linux,windows}-amd64`, or run it directly:

    go run *.go

You can also build the architecture-specific versions with:

    astilectron-bundler -v -c bundler.{darwin,linux,windows}.json

##### Linux

On Linux, you also have to `apt-get install libappindicator1 xclip libgconf-2-4` due to an issue with building Electron-based apps.

## Acknowledgments

While now almost entirely different, this project was jump-started from [OpenBazaar](https://openbazaar.org/). Thanks to @cpacia, @drwasho and the rest of the contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go).
And of course, thank you, [Protocal Labs](https://protocol.ai/), for the incredible FOSS effort and constant inspiration.

## License

MIT
