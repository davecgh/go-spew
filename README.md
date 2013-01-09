go-spew
=======

Go-spew implements a deep pretty printer for Go data structures to aid in
debugging.  It is still under initial development, so some of the formatting
output is likely to change, however it is already quite capable.   It is
licensed under the liberal ISC license, so it may be used in open source or
commercial projects.

## Documentation

Full `go doc` style documentation for the project can be viewed online without
installing this package by using the excellent GoDoc site here:
http://godoc.org/github.com/davecgh/go-spew/spew

You can also view the documentation locally once the package is installed with
the `godoc` tool by running `godoc -http=":6060"` and pointing your browser to
http://localhost:6060/pkg/github.com/davecgh/go-spew/spew

## Installation

```bash
$ go get github.com/davecgh/go-spew/spew
```

## Quick Start

To dump a variable with full newlines, indentation, type, and pointer
information use Dump or Fdump:

```Go
spew.Dump(myVar1, myVar2, ...)
spew.Fdump(someWriter, myVar1, myVar2, ...)
```

Alternatively, if you would prefer to use format strings with a compacted inline
printing style, use the convenience wrappers Printf, Fprintf, etc with either
%v (most compact) or %+v (adds pointer addresses):

```Go
spew.Printf("myVar1: %v -- myVar2: %+v", myVar1, myVar2)
spew.Fprintf(someWriter, "myVar1: %v -- myVar2: %+v", myVar1, myVar2)
```

## Sample Dump Output

```
(main.Foo) {
 unexportedField: (*main.Bar)(0xf84002e210)({
  flag: (main.Flag) flagTwo,
  data: (uintptr) <nil>
 }),
 ExportedField: (map[interface {}]interface {}) {
  (string) "one": (bool) true
 }
}
```

## Sample Formatter Output

Double pointer to a uint8 via %v:
```
	<**>5
```

Circular struct with a uint8 field and a pointer to itself via %+v:
```
	{ui8:1 c:<*>(0xf84002d200){ui8:1 c:<*>(0xf84002d200)<shown>}}
```

## License

Go-spew is licensed under the liberal ISC License.
