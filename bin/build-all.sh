#!/usr/bin/env bash
#!/bin/bash

set -e

BASE=$(dirname $0)
OUTDIR=${BASE}/../out
BINARYNAME=stufy
CWD=$(pwd)
version=$(echo "${TRAVIS_BRANCH:-dev}" | sed 's/^v//')

sha256sumPath="sha256sum"

if [[ "$OSTYPE" == "darwin"* ]]; then
sha256sumPath="shasum -a 256"
fi

function build {
  local arch=$1; shift
  local os=$1; shift
  local armv=$1
  local ext=""
  local arch_name="${arch}"

  if [ "${os}" == "windows" ]; then
      ext=".exe"
  fi
  if [ "${arch_name}" == "arm" ]; then
      arch_name="${arch_name}v${armv}"
  fi
  cd ${CWD}
  echo "building ${BINARYNAME} (${os} ${arch_name})..."
  GOARCH=${arch} GOOS=${os} GOARM=${armv} go build -ldflags="-s -w -X main.Version=${version}" -o $OUTDIR/${BINARYNAME}_${os}_${arch_name}${ext} ./cli || {
    echo >&2 "error: while building ${BINARYNAME} (${os} ${arch_name})"
    return 1
  }

  echo "zipping ${BINARYNAME} (${os} ${arch_name})..."
  cd $OUTDIR
  cp ${BINARYNAME}_${os}_${arch_name}${ext} ${BINARYNAME}${ext}
  zip "${BINARYNAME}_${os}_${arch_name}.zip" "${BINARYNAME}${ext}" || {
    echo >&2 "error: cannot zip file ${BINARYNAME}_${os}_${arch_name}${ext}"
    return 1
  }
  echo "${BINARYNAME}_${os}_${arch_name}${ext} - $(${sha256sumPath} ${BINARYNAME}_${os}_${arch_name}${ext} | awk '{print $1}')" >> sha256.txt
  echo "${BINARYNAME}_${os}_${arch_name}.zip - $(${sha256sumPath} ${BINARYNAME}_${os}_${arch_name}.zip | awk '{print $1}')" >> sha256.txt
  cd ${CWD}
}

build amd64 windows
build amd64 linux
build amd64 darwin
