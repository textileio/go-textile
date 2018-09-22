#!/bin/bash

set -e
pwd

# crypto
go test -coverprofile=crypto.cover.out ./crypto

# keypair
go test -coverprofile=keypair.cover.out ./keypair

# strkey
go test -coverprofile=strkey.cover.out ./strkey

# wallet
go test -coverprofile=wallet.cover.out ./wallet

# repo
go test -coverprofile=repo.db.cover.out ./repo/db
go test -coverprofile=repo.schema.cover.out ./repo/schema
go test -coverprofile=repo.migrations.cover.out ./repo/migrations

# thread
go test -coverprofile=thread.cover.out ./thread

# photo
go test -coverprofile=photo.cover.out ./photo

# core
go test -coverprofile=core.cover.out ./core

# mobile
go test -coverprofile=mobile.cover.out ./mobile

# gateway
go test -coverprofile=gateway.cover.out ./gateway

# cafe
go test -coverprofile=cafe.dao.cover.out ./cafe/dao
go test -coverprofile=cafe.crypto.cover.out ./cafe/crypto
go test -coverprofile=cafe.cover.out ./cafe

echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | \
awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
rm -rf *.cover.out
rm -rf coverage.out
