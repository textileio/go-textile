#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
IOS_DIR="${DIR}/../dist/ios"
ANDROID_DIR="${DIR}/../dist/android"
VER=$1

# fetch iOS framework
mkdir -p ${IOS_DIR}
curl -L -o textile-go_v${VER}_ios-framework.tar.gz https://github.com/textileio/textile-go/releases/download/v${VER}/textile-go_v${VER}_ios-framework.tar.gz
tar xvfz textile-go_v${VER}_ios-framework.tar.gz
rm textile-go_v${VER}_ios-framework.tar.gz
mv Mobile.framework ${IOS_DIR}
mv protos ${IOS_DIR}

# fetch Android framework
mkdir -p ${ANDROID_DIR}
curl -L -o textile-go_v${VER}_android-aar.tar.gz https://github.com/textileio/textile-go/releases/download/v${VER}/textile-go_v${VER}_android-aar.tar.gz
tar xvfz textile-go_v${VER}_android-aar.tar.gz
rm textile-go_v${VER}_android-aar.tar.gz
mv mobile.aar ${ANDROID_DIR}
mv protos ${ANDROID_DIR}
