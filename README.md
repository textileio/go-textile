# textile-go

Textile's REST API and daemon

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/textileio/textile-go)](https://goreportcard.com/report/github.com/textileio/textile-go) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)

This repository contains Textile's API and daemon used to create a cross platform cli and mobile frameworks. The ```textile-go``` library is primarily used in the [Textile Photos](https://www.textile.photos) mobile application. 

Until [Textile Photos](https://www.textile.photos) is ready for public release, this library will be rapidly evolving.

## Contributing

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

## Building

Build the cli based daemon:

```
make build
```

Build the iOS Framework:

```
make ios_framework
``` 

## Acknowledgments

Thanks to @cpacia, @drwasho and the rest of the OpenBazaar contributors for their work on [openbazaar-go](https://github.com/OpenBazaar/openbazaar-go). 

## License

MIT
