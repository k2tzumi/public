test:
	go test -v ./pkg/...

linux:
	docker run -ti --rm -v $(GOPATH):/go/ \
		-e CC=gcc \
		-w /go/src/cirello.io/cci golang \
		/bin/bash -c 'go build -o cci.linux ./cmd/cci'

local:
	CC=gcc go build -o cci.darwin ./cmd/cci
