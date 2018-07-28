test:
	go get -u golang.org/x/vgo
	vgo test -v cirello.io/cci/...

linux:
	docker run -ti --rm -v $(PWD)/../:/go/src/cirello.io/ \
		-w /go/src/cirello.io/cci golang \
		/bin/bash -c 'go get -u golang.org/x/vgo && vgo build -o cci.linux ./cmd/cci'

