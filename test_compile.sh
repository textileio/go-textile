#!/bin/bash

set -e
pwd
go test -coverprofile=central.dao.cover.out ./central/dao
go test -coverprofile=central.controllers.cover.out ./central/controllers
go test -coverprofile=core.cover.out ./core
go test -coverprofile=mobile.cover.out ./mobile
go test -coverprofile=repo.db.cover.out ./repo/db
go test -coverprofile=repo.photos.cover.out ./repo/photos
go test -coverprofile=repo.schema.cover.out ./repo/schema
go test -coverprofile=net.cover.out ./net
go test -coverprofile=crypto.cover.out ./crypto
echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | \
awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
rm -rf *.cover.out
rm -rf coverage.out
