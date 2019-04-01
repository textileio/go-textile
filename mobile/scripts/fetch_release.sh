#!/bin/bash

set -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
IOS_DIR="${DIR}/../dist/ios"
ANDROID_DIR="${DIR}/../dist/android"
VER=$1

# fetch iOS framework
rm -rf ${IOS_DIR}
mkdir ${IOS_DIR}
curl -L -o go-textile_v${VER}_ios-framework.tar.gz https://github.com/textileio/go-textile/releases/download/v${VER}/go-textile_v${VER}_ios-framework.tar.gz
tar xvfz go-textile_v${VER}_ios-framework.tar.gz
rm go-textile_v${VER}_ios-framework.tar.gz
mv Mobile.framework ${IOS_DIR}
mv protos ${IOS_DIR}

# fetch Android framework
rm -rf ${ANDROID_DIR}
mkdir ${ANDROID_DIR}
curl -L -o go-textile_v${VER}_android-aar.tar.gz https://github.com/textileio/go-textile/releases/download/v${VER}/go-textile_v${VER}_android-aar.tar.gz
tar xvfz go-textile_v${VER}_android-aar.tar.gz
rm go-textile_v${VER}_android-aar.tar.gz
mv mobile.aar ${ANDROID_DIR}
mv protos ${ANDROID_DIR}
