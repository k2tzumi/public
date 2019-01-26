# HTTP log monitoring console program

Install with:

```sh
go get -u cirello.io/logmon/...
```

Assuming your `$GOPATH` is part of your `$PATH`, you can use `loggen` to
generate logs that can be parsed with `logmon`

```sh
loggen &
logmon access.log
```

Refer to this [document](DESIGN.md) to check some design considerations.