#!/bin/bash

git remote | grep -v -E "origin|_" | while read line; do ./scripts/push-repo.sh $line; done;

