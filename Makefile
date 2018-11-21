build:
	go build -ldflags "-w" -i -o textile textile.go
	mv textile dist/

build_ios_framework:
	gomobile bind -ldflags "-w" -target=ios github.com/textileio/textile-go/mobile

build_android_framework:
	gomobile bind -ldflags "-w" -target=android -o textilego.aar github.com/textileio/textile-go/mobile

P_TIMESTAMP=Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp
P_ANY=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any
PKGMAP=$(P_TIMESTAMP),$(P_ANY)
protos:
	cd pb/protos && PATH=$(PATH):$(GOPATH)/bin protoc --go_out=$(PKGMAP):.. *.proto
