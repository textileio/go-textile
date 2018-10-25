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

## Usage
```
~ $ go run textile.go --help
Usage:
  textile [OPTIONS]

Application Options:
  -v, --version            print the version number and exit
  -r, --repo-dir=          specify a custom repository path
  -l, --log-level=         set the logging level [debug, info, notice, warning, error, critical] (default: debug)
  -n, --no-log-files       do not save logs on disk
  -d, --daemon             start in a non-interactive daemon mode
  -s, --server             start in server mode
  -g, --gateway-bind-addr= set the gateway address (default: 127.0.0.1:random)
      --swarm-ports=       set the swarm ports (tcp,ws) (default: random)
  -c, --cafe=              cafe host address
      --cafe-bind-addr=    set the cafe address
      --cafe-db-hosts=     set the cafe mongo db hosts uri
      --cafe-db-name=      set the cafe mongo db name
      --cafe-db-user=      set the cafe mongo db user
      --cafe-db-password=  set the cafe mongo db user password
      --cafe-db-tls        use TLS for the cafe mongo db connection
      --cafe-token-secret= set the cafe token secret
      --cafe-referral-key= set the cafe referral key

Help Options:
  -h, --help               Show this help message
```

## CLI Usage
```
~ $ go run textile.go
Textile
version: 0.1.2
repo: /Users/sander/.textile/repo
gateway: 127.0.0.1:5446
type 'help' for available commands
>>> help

Commands:
  cafe                manage cafe session
  clear               clear the screen
  device              manage connected devices
  exit                exit the program
  fetch-messages      fetch offline messages from the DHT
  help                display help
  id                  show node id
  notification        manage notifications
  photo               manage photos
  ping                ping another peer
  profile             manage cafe profiles
  start               start the node
  stop                stop the node
  swarm               same as ipfs swarm
  thread              manage threads
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

## Webhook Integration

After building the CLI/daemon, run either `shell` or `daemon` with the optional `--slack-hook` (or `-s`) and/or `--welcome-image` flags. Adding a welcome image will automatically post a welcome message and image to all Threads that the textile daemon has joined. While this is setup to support Slack webhooks, it should work for any webhook that can mimic [Slack's API](https://api.slack.com/incoming-webhooks).

```
textile daemon -s https://hooks.slack.com/services/blah/blah --welcome-image ./welcome.png
```


## Contributing

The easiest way to write a valid commit message is to use the `npm` script:

```
$ npm run cm
```

This will start the interactive commit prompt.

## Acknowledgments

Thanks to @cpacia, @drwasho and the rest of the OpenBazaar contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go). 

## License

MIT
