#!/bin/bash

set -e
pwd
go test -coverprofile=cafe.dao.cover.out ./cafe/dao
go test -coverprofile=cafe.cover.out ./cafe
go test -coverprofile=core.cover.out ./core
go test -coverprofile=util.cover.out ./util
go test -coverprofile=wallet.cover.out ./wallet
go test -coverprofile=mobile.cover.out ./mobile
go test -coverprofile=repo.db.cover.out ./repo/db
go test -coverprofile=repo.schema.cover.out ./repo/schema
go test -coverprofile=crypto.cover.out ./crypto
echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | \
awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
rm -rf *.cover.out
rm -rf coverage.out
