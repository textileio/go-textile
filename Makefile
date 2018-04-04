build:
	./build.sh

linux_binary:
	./build.sh linux/amd64

ios_framework:
	CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios/arm64 github.com/textileio/textile-go/mobile

android_framework:
	gomobile bind -target=android -o textilego.aar github.com/textileio/textile-go/mobile

clean_build:
	rm -rf dist && rm -f Mobile.framework

clean: clean_build
