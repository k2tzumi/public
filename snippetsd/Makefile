all: assets darwin

assets:
	cd frontend; npm run build;
	go-bindata-assetfs -o bindata_assetfs.go -pkg generated frontend/build/...
	mv bindata_assetfs.go generated

darwin:
	go build -o snippetsd ./cmd/snippetsd

linux-docker:
	docker run -ti --rm -v $(PWD)/../:/go/src/cirello.io/ \
		-w /go/src/cirello.io/snippetsd golang \
		/bin/bash -c 'go build -o snippetsd.linux ./cmd/snippetsd'

linux:
	GOOS=linux go build -o snippetsd.linux ./cmd/snippetsd

test:
	go test -v ./...
