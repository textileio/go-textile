# go-textile

[![Made by Textile](https://img.shields.io/badge/made%20by-Textile-informational.svg?style=popout-square)](https://textile.io)
[![Chat on Slack](https://img.shields.io/badge/slack-slack.textile.io-informational.svg?style=popout-square)](https://slack.textile.io)
[![GitHub license](https://img.shields.io/github/license/textileio/go-textile.svg?style=popout-square)](./LICENSE)
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
-   [Contributing](#contributing)
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

    env GO111MODULE=on go get github.com/textileio/go-textile
    env GO111MODULE=on go install github.com/textileio/go-textile/cmd/textile

[Installation instructions](https://docs.textile.io/install/the-daemon/) for pre-built binaries are in [the docs](https://docs.textile.io).

## Usage

Go to https://godoc.org/github.com/textileio/go-textile.

The [Tour of Textile](https://docs.textile.io/a-tour-of-textile/) goes through many examples and use cases. `textile --help` provides a quick look at the available APIs. For a full overview of every CLI command available, refer to our [Command Line Documentation](https://docs.textile.io/develop/clients/command-line/).

### Requirements

-   go >= 1.12

Extra setup steps are needed to build the bindings for iOS or Android, as `gomobile` does not yet support [go modules](https://github.com/golang/go/wiki/Modules). You'll need to **move the go-textile source** into your `GOPATH` (like pre-go1.11 development), before installing and initializing the `gomobile` tools:

    go get golang.org/x/mobile/cmd/gomobile
    gomobile init

Now you can execute the iOS and Android build tasks below. For the other build tasks, the source must _not_ be under `GOPATH`. Go 1.13 is supposed to bring module support to `gomobile`, at which point we can remove this madness!

### Install dependencies:

    make setup

### Build `textile`:

    make textile

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
