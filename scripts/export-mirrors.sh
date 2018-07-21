#!/bin/bash

export REPOS="HumorChecker dynamolock errors goherokuname runner supervisor bookmarkd snippetsd"

for repo in $REPOS
do
	./scripts/push-repo.sh $repo
done

