all: darwin

darwin:
	vgo build -o bookmarkd ./cmd/bookmarkd

linux:
	docker run -ti --rm -e CC=gcc -v $(GOPATH):/go/ \
		-w /go/src/cirello.io/cci golang \
		/bin/bash -c 'go get -u golang.org/x/vgo && vgo build -o sdci.linux ./cmd/sdci'

test:
	go get -u golang.org/x/vgo
	vgo test -v ./...
