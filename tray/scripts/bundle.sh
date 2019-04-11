#!/usr/bin/env bash

set -e

OS="all"

while getopts o: option
do
case "${option}"
in
o) OS=${OPTARG};;
esac
done

go install ../vendor/github.com/asticode/go-astilectron-bundler/astilectron-bundler
astilectron-bundler -v -c bundler.$OS.json