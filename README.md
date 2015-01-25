# FireGo
[FirePHP](https://github.com/firephp/firephp-core) ported to Go.

It partially implements the [FirePHP Protocol](http://www.firephp.org/Wiki/Reference/Protocol), supporting:

1. Log
2. Info
3. Warn
4. Error

The `TRACE`, `EXCEPTION`, `TABLE` and `GROUP` are not implemented - I still need to understand whether it is desirable and possible to port these message types.

Also, it does not analyse the backtrace to feed the information with extra information such filename and line. http://golang.org/pkg/runtime/#Stack should do the trick.

[Check the example to see it working](examples/example.go):

`# go run examples/example.go`

