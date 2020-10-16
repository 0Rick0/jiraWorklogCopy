#!/bin/bash

OSES=(
  'linux:amd64:jiraWorklogCopyLinux'
  'darwin:amd64:jiraWorklogCopyDarwin'
  'windows:amd64:jiraWorklogCopy.exe'
)

if [ ! -d target ]; then
    mkdir "target"
fi

for INFO in "${OSES[@]}" ; do
  IFS=: read -r OS ARCH FILENAME <<< "${INFO}"
  [ ! -d "target/${OS}" ] && mkdir "target/${OS}"
  GOOS=$OS GOARCH=$ARCH go build -o "target/${OS}/${FILENAME}" cmd/jiraWorklogCopy/main.go
done