# textile-go

Textile CLI, desktop app, mobile bindings, and REST API. See [textile-mobile](https://github.com/textileio/textile-mobile/) for iOS and Android apps.

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/textile-go)](https://goreportcard.com/report/github.com/textileio/textile-go) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![CircleCI](https://circleci.com/gh/textileio/textile-go/tree/master.svg?style=shield)](https://circleci.com/gh/textileio/textile-go/tree/master)

![textile](https://s3.amazonaws.com/textile.public/cli_3.png)

This repository contains a cross platform cli, desktop application, and iOS/Android mobile bindings for running a Textile node. See [Textile Photos](https://www.textile.photos) for more info. 

Until [Textile Photos](https://www.textile.photos) is ready for public release, this library will be rapidly evolving. 

## Install

```
go get github.com/textileio/textile-go
```

You'll need a few different tools here to get setup...

#### Install `dep`

Golang package manager:

```
brew install dep
```

#### Install `gx`

IPFS package manager:

```
go get -u github.com/whyrusleeping/gx
go get -u github.com/whyrusleeping/gx-go
```

#### Install `node`

NodeJS is used for git hooks and some build tooling:

```
brew install node
```

#### Install dependencies

Finally, download deps managed by `gx` and `dep`:

```
npm run setup
```

This will start the interactive commit prompt.

## Building

These instructions assume the build OS is either Darwin or Linux. As such, `mingw-w64` is needed for Windows cross-compiled builds. `brew install mingw-w64` for Darwin, `apt-get install mingw-w64` for Debian, etc.

There are various things to build:

#### The CLI:

```
make build
```

#### The iOS Framework:

```
make ios_framework
```

#### The Android Framework:

```
make android_framework
```

#### The Desktop Application

The build is made by a vendored version of `go-astilectron-bundler`. Due to Go's painful package management, you'll want to delete any `go-astilectron`-related binaries and source code you have installed from `github.com/asticode` in your `$GOPATH`. Then you can install the vendored `go-astilectron-bundler`:
```
go install ./vendor/github.com/asticode/go-astilectron-bundler/astilectron-bundler
```
Run `make` to build the app for Darwin, Linux, and Windows:
```
make build_desktop
```
Double-click the built app in `desktop/output/darwin-amd64`, or run it directly:
```
cd desktop && go run *.go
```
See [go-astilectron-bundler](https://github.com/asticode/go-astilectron-bundler) for more build configurations.

## Contributing

The easiest way to write a valid commit message is to use the `npm` script:

```
npm run cm
```

## Acknowledgments

Thanks to @cpacia, @drwasho and the rest of the OpenBazaar contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go). 

## License

MIT
