setup:
	go mod download
	go get github.com/ahmetb/govvv
	npm install

test:
	./test_compile.sh

fmt:
	echo 'Formatting with prettier...'
	npx prettier --write "./**" 2> /dev/null || true
	echo 'Formatting with goimports...'
	goimports -w -l `find . -type f -name '*.go' -not -path './vendor/*'`

lint:
	echo 'Linting with prettier...'
	npx prettier --check "./**" 2> /dev/null || true
	echo 'Linting with golint...'
	golint `go list ./... | grep -v /vendor/`

build:
	$(eval FLAGS := $$(shell govvv -flags | sed 's/main/github.com\/textileio\/go-textile\/common/g'))
	go build -ldflags "-w $(FLAGS)" -i -o textile textile.go
	mkdir -p dist
	mv textile dist/

install:
	mv dist/textile $$GOPATH/bin

ios:
	$(eval FLAGS := $$(shell govvv -flags | sed 's/main/github.com\/textileio\/go-textile\/common/g'))
	env go111module=off gomobile bind -ldflags "-w $(FLAGS)" -v -target=ios github.com/textileio/go-textile/mobile github.com/textileio/go-textile/core
	mkdir -p mobile/dist/ios/ && cp -r Mobile.framework mobile/dist/ios/
	rm -rf Mobile.framework

android:
	$(eval FLAGS := $$(shell govvv -flags | sed 's/main/github.com\/textileio\/go-textile\/common/g'))
	env go111module=off gomobile bind -ldflags "-w $(FLAGS)" -v -target=android -o mobile.aar github.com/textileio/go-textile/mobile github.com/textileio/go-textile/core
	mkdir -p mobile/dist/android/ && mv mobile.aar mobile/dist/android/

protos:
	$(eval P_TIMESTAMP := Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp)
	$(eval P_ANY := Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any)
	$(eval PKGMAP := $$(P_TIMESTAMP),$$(P_ANY))
	cd pb/protos; protoc --go_out=$(PKGMAP):.. *.proto

protos_js:
	rm -rf mobile/dist
	mkdir mobile/dist
	cd mobile; rm -rf node_modules && npm i @textile/protobufjs --no-save
	cd mobile; ./node_modules/.bin/pbjs -t static-module -o dist/index.js ../pb/protos/*
	cd mobile; ./node_modules/.bin/pbts -o dist/index.d.ts dist/index.js

.PHONY: docs
docs:
	go get github.com/swaggo/swag/cmd/swag
	swag init -g core/api.go
	npm i -g swagger-markdown
	swagger-markdown -i docs/swagger.yaml -o docs/swagger.md

# Additional dependencies needed below:
# $ brew install jq
# $ brew install grep

.PHONY: publish
publish:
	make protos_js
	$(eval VERSION := $$(shell ggrep -oP 'const Version = "\K[^"]+' common/version.go))
	cd mobile; jq '.version = "$(VERSION)"' package.json > package.json.tmp && mv package.json.tmp package.json
	cd mobile; npm publish
	cd mobile; jq '.version = "0.0.0"' package.json > package.json.tmp && mv package.json.tmp package.json

docker:
	$(eval VERSION := $$(shell ggrep -oP 'const Version = "\K[^"]+' common/version.go))
	docker build -t go-textile:$(VERSION) .

docker_cafe:
	$(eval VERSION := $$(shell ggrep -oP 'const Version = "\K[^"]+' common/version.go))
	docker build -t go-textile:$(VERSION)-cafe -f Dockerfile.cafe .
