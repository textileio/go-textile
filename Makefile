build:
	gox -osarch="linux/amd64 linux/386 linux/arm darwin/amd64" -output="textile-go-{{.OS}}-{{.Arch}}"
	CC="x86_64-w64-mingw32-gcc" CXX="x86_64-w64-mingw32-g++" gox -cgo -osarch="windows/386 windows/amd64" -output="textile-go-{{.OS}}-{{.Arch}}"
	mv textile-go-* dist

build_desktop:
	$(MAKE) -C ./desktop build

build_ios_framework:
	CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios/arm64 github.com/textileio/textile-go/mobile
	# cp -r Mobile.framework ~/github/textileio/textile-mobile/ios/

build_android_framework:
	gomobile bind -target=android -o textilego.aar github.com/textileio/textile-go/mobile
	# cp -r textilego.aar ~/github/textileio/textile-mobile/android/textilego/

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
