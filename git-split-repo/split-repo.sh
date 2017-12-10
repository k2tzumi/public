#!/bin/bash
git filter-branch --prune-empty -f --tree-filter 'if [[ -e "common" ]]; then mv common runner/; fi' HEAD
git filter-branch --prune-empty -f --subdirectory-filter runner HEAD
git filter-branch --prune-empty -f --tree-filter 'rm -rf ./agent ./config ./sim ./swapi ./tasker' HEAD

gofmt -r '"golang.org/x/net/context" -> "context"' -w .
gofmt -r '"github.com/iron-io/worker/common" -> "github.com/iron-io/runner/common"' -w .
gofmt -r '"github.com/iron-io/worker/common/stats" -> "github.com/iron-io/runner/common/stats"' -w .
gofmt -r '"github.com/iron-io/worker/runner/drivers" -> "github.com/iron-io/runner/drivers"' -w .

###
# git filter-branch --prune-empty -f --tree-filter 'rm -rf ./runner ./agent/common' extract-runner
# gofmt -r '"github.com/iron-io/worker/runner" -> "github.com/iron-io/runner"' -w .
# gofmt -r '"github.com/iron-io/worker/runner" -> "github.com/iron-io/runner"' -w .
##
