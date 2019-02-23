P_TIMESTAMP=Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp
P_ANY=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any
PKGMAP=$(P_TIMESTAMP),$(P_ANY)

clean:
	rm -rf vendor

setup:
	go get github.com/ahmetb/govvv
	dep ensure
	gx install

test_compile:
	./test_compile.sh

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

cross_build_linux:
	$(eval FLAGS := $$(shell govvv -flags -pkg github.com/textileio/textile-go/common))
	export CGO_ENABLED=1
	docker pull karalabe/xgo-latest
	go get github.com/karalabe/xgo
	xgo -go 1.11.1 -ldflags "-w $(FLAGS)" --targets=linux/amd64 .
	chmod +x textile-go-linux-amd64
	mkdir -p dist
	mv textile-go-linux-amd64 dist/

build_ios_framework:
	$(eval FLAGS := $$(shell govvv -flags -pkg github.com/textileio/textile-go/common))
	gomobile bind -ldflags "-w $(FLAGS)" -target=ios github.com/textileio/textile-go/mobile
	cp -r Mobile.framework mobile/dist/
	rm -rf Mobile.framework

build_android_aar:
	$(eval FLAGS := $$(shell govvv -flags -pkg github.com/textileio/textile-go/common))
	gomobile bind -ldflags "-w $(FLAGS)" -target=android -o mobile.aar github.com/textileio/textile-go/mobile
	mv mobile.aar mobile/dist/

build_mobile:
	make build_ios_framework
	build_android_aar
	make protos_js

# Additional dependencies needed:
# $ brew install jq
# $ brew install grep
publish_mobile:
	#$(eval VERSION := $$(shell ggrep -oP 'const Version = "\K[^"]+' common/version.go))
	#cd mobile; jq '.version = "$(VERSION)"' package.json > package.json.tmp && mv package.json.tmp package.json
	#cd mobile; npm publish
	#cd mobile; jq '.version = "0.0.0"' package.json > package.json.tmp && mv package.json.tmp package.json

protos:
	cd pb/protos && PATH=$(PATH):$(GOPATH)/bin protoc --go_out=$(PKGMAP):.. *.proto

protos_js:
	cd mobile; yarn install --ignore-scripts
	cd mobile; node node_modules/@textile/protobufjs/cli/bin/pbjs -t static-module -w es6 -o dist/index.js ../pb/protos/*
	cd mobile; node node_modules/@textile/protobufjs/cli/bin/pbts -o dist/index.d.ts dist/index.js

build_docs:
	go get github.com/swaggo/swag/cmd/swag
	swag init -g core/api.go
