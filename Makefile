build:
	go build -i -o textile textile.go

build_ios_framework:
	CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios/arm64 github.com/textileio/textile-go/mobile

build_android_framework:
	gomobile bind -target=android -o textilego.aar github.com/textileio/textile-go/mobile

build_cafe:
	go get github.com/kardianos/govendor
	govendor init && govendor add +external
	docker-compose build
	rm -rf vendor/gx && rm vendor/vendor.json

P_TIMESTAMP=Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp
P_ANY=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any
PKGMAP=$(P_TIMESTAMP),$(P_ANY)
protos:
	cd pb/protos && PATH=$(PATH):$(GOPATH)/bin protoc --go_out=$(PKGMAP):.. *.proto
