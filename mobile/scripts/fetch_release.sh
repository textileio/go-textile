#!/bin/bash

DIR="$(dirname "$(stat -f "$0")")"
IOS_DIR="${DIR}/../ios"
ANDROID_DIR="${DIR}/../android"
VER=$1

# fetch iOS framework
mkdir -p ${IOS_DIR}
curl -L -o textile-go_v${VER}_ios-framework.tar.gz https://github.com/textileio/textile-go/releases/download/v${VER}/textile-go_v${VER}_ios-framework.tar.gz
tar xvfz textile-go_v${VER}_ios-framework.tar.gz
rm textile-go_v${VER}_ios-framework.tar.gz
mv Mobile.framework ${IOS_DIR}

# fetch Android framework
mkdir -p ${ANDROID_DIR}
curl -L -o textile-go_v${VER}_android-aar.tar.gz https://github.com/textileio/textile-go/releases/download/v${VER}/textile-go_v${VER}_android-aar.tar.gz
tar xvfz textile-go_v${VER}_android-aar.tar.gz
rm textile-go_v${VER}_android-aar.tar.gz
mv textilego.aar ${ANDROID_DIR}
