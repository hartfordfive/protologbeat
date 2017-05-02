#!/bin/bash

if [ "$1" == "" ]; then
  echo "Usage: ./build-bin [VERSION]"
  exit 1
fi

BEAT_NAME="protologbeat"
VERSION=$1
ARCH=amd64

for OS in linux darwin windows; do
  echo "Building ${BEAT_NAME} v${VERSION} for ${OS}/${ARCH}..."
  GOOS=$OS GOARCH=$ARCH go build -ldflags "-s -w -X main.build_date=`date +%Y-%m-%d` -X main.version=${VERSION} -X main.commit_hash=`git rev-parse --verify HEAD`" -o ${BEAT_NAME}-${VERSION}-${OS}-x86_64
  
  zip ${BEAT_NAME}-${VERSION}-${OS}-x86_64.zip ${BEAT_NAME}-${VERSION}-${OS}-x86_64 &&  mv ${BEAT_NAME}-${VERSION}-${OS}-x86_64.zip bin/
  tar -cvzf ${BEAT_NAME}-${VERSION}-${OS}-x86_64.tar.gz ${BEAT_NAME}-${VERSION}-${OS}-x86_64 && mv ${BEAT_NAME}-${VERSION}-${OS}-x86_64.tar.gz bin/ 
  rm ${BEAT_NAME}-${VERSION}-${OS}-x86_64
done

