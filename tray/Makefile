ifeq ($(OS),Windows_NT) 
	DETECTED_OS := windows
else
	DETECTED_OS := $(shell sh -c 'uname 2>/dev/null || echo unknown' | tr A-Z a-z)
endif

build:
	go install ../vendor/github.com/asticode/go-astilectron-bundler/astilectron-bundler
	astilectron-bundler -v -c bundler.$(DETECTED_OS).json

run: build
	go run *.go
