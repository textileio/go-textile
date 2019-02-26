clean:
	rm -rf vendor

setup:
	go get github.com/ahmetb/govvv
	dep ensure
	gx install
	cd mobile; rm -rf node_modules && npm i @textile/protobufjs --no-save

test:
	./test.sh

fmt:
	goimports -w -l `find . -type f -name '*.go' -not -path './vendor/*'`

lint:
	golint `go list ./... | grep -v /vendor/`

build:
	$(eval FLAGS := $$(shell govvv -flags -pkg github.com/textileio/textile-go/common))
	go build -ldflags "-w $(FLAGS)" -i -o textile textile.go
	mv textile dist/

install:
	mv dist/textile /usr/local/bin

linux:
	$(eval FLAGS := $$(shell govvv -flags -pkg github.com/textileio/textile-go/common))
	export CGO_ENABLED=1
	docker pull karalabe/xgo-latest
	go get github.com/karalabe/xgo
	xgo -go 1.11.1 -ldflags "-w $(FLAGS)" --targets=linux/amd64 .
	chmod +x textile-go-linux-amd64
	mkdir -p dist
	mv textile-go-linux-amd64 dist/

ios:
	$(eval FLAGS := $$(shell govvv -flags -pkg github.com/textileio/textile-go/common))
	gomobile bind -ldflags "-w $(FLAGS)" -target=ios github.com/textileio/textile-go/mobile
	cp -r Mobile.framework mobile/dist/
	rm -rf Mobile.framework

android:
	$(eval FLAGS := $$(shell govvv -flags -pkg github.com/textileio/textile-go/common))
	gomobile bind -ldflags "-w $(FLAGS)" -target=android -o mobile.aar github.com/textileio/textile-go/mobile
	mv mobile.aar mobile/dist/

protos:
	$(eval P_TIMESTAMP := Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp)
	$(eval P_ANY := Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any)
	$(eval PKGMAP := $$(P_TIMESTAMP),$$(P_ANY))
	cd pb/protos; protoc --go_out=$(PKGMAP):.. *.proto
	mkdir -p mobile/dist
	cd mobile; ./node_modules/.bin/pbjs -t static-module -w es6 -o dist/index.js ../pb/protos/*
	cd mobile; ./node_modules/.bin/pbts -o dist/index.d.ts dist/index.js

.PHONY: docs
docs:
	go get github.com/swaggo/swag/cmd/swag
	swag init -g core/api.go

# Additional dependencies needed:
# $ brew install jq
# $ brew install grep
.PHONY: publish
publish:
	make protos
	$(eval VERSION := $$(shell ggrep -oP 'const Version = "\K[^"]+' common/version.go))
	cd mobile; jq '.version = "$(VERSION)"' package.json > package.json.tmp && mv package.json.tmp package.json
	cd mobile; npm publish
	cd mobile; jq '.version = "0.0.0"' package.json > package.json.tmp && mv package.json.tmp package.json
