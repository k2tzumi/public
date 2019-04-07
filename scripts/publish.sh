#!/bin/bash

exit 0

export REPOS="HumorChecker dynamolock errors goherokuname runner supervisor bookmarkd snippetsd cci mlflappygopher kohrah-ani oversight pglock junk svc gists exp poo alfredemoji"

for repo in $REPOS;
do
	cd $repo;

	ROOT_PACKAGE="$(go list | sed 's/cirello.io/github.com\/cirello-io/g')"
	PACKAGE_NAME="$(go list | sed 's/cirello.io/github.com\/cirello-io/g' | xargs basename)"

	echo "PACKAGE NAME $PACKAGE_NAME";
	echo "ROOT PACKAGE $ROOT_PACKAGE";

	go list ./... | while read line; do

		DIR="$(echo $line | sed 's/cirello.io\///')"
		echo $DIR;
		mkdir -p $GOPATH/src/github.com/cirello-io/export/$DIR

		cat > $GOPATH/src/github.com/cirello-io/export/$DIR/index.html <<EOT
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="cirello.io/$PACKAGE_NAME git https://$ROOT_PACKAGE">
<meta http-equiv="refresh" content="0; url=https://godoc.org/$line">
</head>
<body>
Redirecting to docs at <a href="https://godoc.org/$line">godoc.org/$line</a>...
</body>
</html>
EOT

	done;

	cd -;
done;