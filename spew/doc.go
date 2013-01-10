/*
 * Copyright (c) 2013 Dave Collins <dave@davec.name>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
Package spew implements a deep pretty printer for Go data structures to aid in
debugging.

A quick overview of the additional features spew provides over the built-in
printing facilities for Go data types are as follows:

	* Pointers are dereferenced and followed
	* Circular data structures are detected and handled properly
	* Custom error/Stringer interfaces are optionally invoked, including
	  on unexported types
	* Custom types which only implement the error/Stringer interfaces via
	  a pointer receiver are optionally invoked when passing non-pointer
	  variables

There are two different approaches spew allows for dumping Go data structures:

	* Dump style which prints with newlines, customizable indentation,
	  and additional debug information such as types and all pointer addresses
	  used to indirect to the final value
	* A custom Formatter interface that integrates cleanly with the standard fmt
	  package and replaces %v and %+v to provide inline printing similar
	  to the default %v while providing the additional functionality outlined
	  above and passing unsupported format verb/flag combinations such a %x,
	  %q, and %#v along to fmt

Quick Start

This section demonstrates how to quickly get started with spew.  See the
sections below for further details on formatting and configuration options.

To dump a variable with full newlines, indentation, type, and pointer
information use Dump or Fdump:
	spew.Dump(myVar1, myVar2, ...)
	spew.Fdump(someWriter, myVar1, myVar2, ...)

Alternatively, if you would prefer to use format strings with a compacted inline
printing style, use the convenience wrappers Printf, Fprintf, etc with either
%v (most compact) or %+v (adds pointer addresses):
	spew.Printf("myVar1: %v -- myVar2: %+v", myVar1, myVar2)
	spew.Fprintf(someWriter, "myVar1: %v -- myVar2: %+v", myVar1, myVar2)

Configuration Options

The following configuration options are available:
	spew.Config.MaxDepth
		Maximum number of levels to descend into nested data structures.
		There is no limit by default.

	spew.Config.Indent
		String to use for each indentation level for Dump functions.
		It is a single space by default.  A popular alternative is "\t".

	spew.Config.DisableMethods
		Disables invocation of error and Stringer interface methods.
		Method invocation is enabled by default.

	spew.Config.DisablePointerMethods
		Disables invocation of error and Stringer interface methods on types
		which only accept pointer receivers from non-pointer variables.
		Pointer method invocation is enabled by default.

Dump Usage

Simply call spew.Dump with a list of variables you want to dump:

	spew.Dump(myVar1, myVar2, ...)

You may also call spew.Fdump if you would prefer to output to an arbitrary
io.Writer.  For example, to dump to standard error:

	spew.Fdump(os.Stderr, myVar1, myVar2, ...)

Sample Dump Output

See the Dump example for details on the setup of the types and variables being
shown here.

	(main.Foo) {
	 unexportedField: (*main.Bar)(0xf84002e210)({
	  flag: (main.Flag) flagTwo,
	  data: (uintptr) <nil>
	 }),
	 ExportedField: (map[interface {}]interface {}) {
	  (string) "one": (bool) true
	 }
	}

Custom Formatter

Spew provides a custom formatter the implements the fmt.Formatter interface
so that it integrates cleanly with standard fmt package printing functions. The
formatter is useful for inline printing of smaller data types similar to the
standard %v format specifier.

The spew formatter only responds to the %v and %+v verb combinations.  Any other
variations such as %x, %q, and %#v will be sent to the the standard fmt package
for formatting.  In addition, the spew formatter ignores the width and precision
arguments (however they will still work on the format specifiers spew does not
handle).

Custom Formatter Usage

The simplest way to make use of the spew custom formatter is to call one of the
convenience functions such as spew.Printf, spew.Println, or spew.Printf.  The
functions have the exact same syntax you are most likely already familiar with:

	spew.Printf("myVar1: %v -- myVar2: %+v", myVar1, myVar2)
	spew.Println(myVar, myVar2)
	spew.Fprintf(os.Stderr, "myVar1: %v -- myVar2: %+v", myVar1, myVar2)

See the Index for the full list convenience functions.

Sample Formatter Output

Double pointer to a uint8 via %v:
	<**>5

Circular struct with a uint8 field and a pointer to itself via %+v:
	{ui8:1 c:<*>(0xf84002d200){ui8:1 c:<*>(0xf84002d200)<shown>}}

See the Printf example for details on the setup of variables being shown
here.

Errors

Since it is possible for custom error/Stringer interfaces to panic, spew
detects them and handles them internally by printing the panic information
inline with the output.  Since spew is intended to provide deep pretty printing
capabilities on structures, it intentionally does not return any errors.
*/
package spew
