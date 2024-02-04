alias b := build
alias r := run

bindir:
  #!/usr/bin/env bash
  [ -e ./bin ] || mkdir bin
build: bindir
  #!/usr/bin/env bash
  set -e
  go mod download
  go build -o ./bin/main
run: build
  ./bin/main

