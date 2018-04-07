build:
	$(MAKE) -C ./textile build

linux_binary:
	$(MAKE) -C ./textile linux_binary

ios_framework:
	CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios/arm64 github.com/textileio/textile-go/mobile github.com/textileio/textile-go/net

android_framework:
	gomobile bind -target=android -o textilego.aar github.com/textileio/textile-go/mobile github.com/textileio/textile-go/net

clean_build:
	rm -rf dist && rm -f Mobile.framework

clean: clean_build

docker_build:
	docker build -t circleci:1.10 .
