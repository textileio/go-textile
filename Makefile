goinstall:
	go build -i -o textile textile.go
	mv textile $(GOPATH)/bin/

build:
	./build.sh

linux_binary:
	./build.sh linux/amd64
	cd dist && tar -czvf textile-go-linux-amd64.tar.gz textile-go-linux-amd64 && cd ..

build_desktop:
	$(MAKE) -C ./desktop build

ios_framework:
	CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios/arm64 github.com/textileio/textile-go/mobile
	# cp -r Mobile.framework ~/github/textileio/textile-mobile/ios/

android_framework:
	gomobile bind -target=android -o textilego.aar github.com/textileio/textile-go/mobile
	# cp -r textilego.aar ~/github/textileio/textile-mobile/android/textilego/

P_TIMESTAMP = Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp
P_ANY = Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any
PKGMAP = $(P_TIMESTAMP),$(P_ANY)
protos:
	cd pb/protos && PATH=$(PATH):$(GOPATH)/bin protoc --go_out=$(PKGMAP):.. *.proto

clean:
	rm -rf dist && rm -f Mobile.framework && rm -rf textilego.aar && rm -rf textilego-sources.jar

build_test:
	docker build -f Dockerfile.circleci -t circleci:1.10 .

build_cafe:
	go get github.com/kardianos/govendor
	govendor init && govendor add +external
	docker-compose build
	rm -rf vendor/gx && rm vendor/vendor.json
