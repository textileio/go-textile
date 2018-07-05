build:
	./build.sh

linux_binary:
	./build.sh linux/amd64

build_desktop:
	$(MAKE) -C ./desktop build

ios_framework:
	CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios/arm64 github.com/textileio/textile-go/mobile github.com/textileio/textile-go/net
	# cp -r Mobile.framework ~/github/textileio/textile-mobile/ios/

android_framework:
	gomobile bind -target=android -o textilego.aar github.com/textileio/textile-go/mobile github.com/textileio/textile-go/net
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

build_swarm_services:
	go get github.com/kardianos/govendor
	cd central && govendor init && govendor add +external
	cd relay && govendor init && govendor add +external
	docker-compose -f docker-compose.swarm.yml build
	rm -rf central/vendor && rm -rf relay/vendor

build_local_services:
	go get github.com/kardianos/govendor
	cd central && govendor init && govendor add +external
	cd relay && govendor init && govendor add +external
	docker-compose build
	rm -rf central/vendor && rm -rf relay/vendor
