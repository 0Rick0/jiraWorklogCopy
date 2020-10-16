#!/bin/bash

OSES=(
  'linux:amd64:'
  'darwin:amd64:'
  'windows:amd64:.exe'
)

if [ ! -d target ]; then
    mkdir "target"
fi

for INFO in "${OSES[@]}" ; do
  IFS=: read -r OS ARCH SUFFIX <<< "${INFO}"
  [ ! -d "target/${OS}" ] && mkdir "target/${OS}"
  GOOS=$OS GOARCH=$ARCH go build -o "target/${OS}/jiraWorklogCopy${SUFFIX}" cmd/jiraWorklogCopy/main.go
done