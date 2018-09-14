#!/bin/bash

rm -rf $1
git add $1
git commit -m "$1: re-import $1"
./scripts/add-repo.sh $1
