build:
	go get github.com/kardianos/govendor
	govendor init
	govendor add +external
	docker-compose build
	rm -rf vendor