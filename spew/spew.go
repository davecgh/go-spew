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

package spew

import (
	"fmt"
	"io"
	"os"
)

// Errorf is a wrapper for fmt.Errorf that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the formatted string as a value that satisfies error.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Errorf(format, spew.NewFormatter(a), spew.NewFormatter(b))
func Errorf(format string, a ...interface{}) (err error) {
	return fmt.Errorf(format, convertArgs(a)...)
}

// Fprint is a wrapper for fmt.Fprint that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprint(w, spew.NewFormatter(a), spew.NewFormatter(b))
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, convertArgs(a)...)
}

// Fprintf is a wrapper for fmt.Fprintf that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprintf(w, format, spew.NewFormatter(a), spew.NewFormatter(b))
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format, convertArgs(a)...)
}

// Fprintln is a wrapper for fmt.Fprintln that treats each argument as if it
// passed with a default Formatter interface returned by NewFormatter.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprintln(w, spew.NewFormatter(a), spew.NewFormatter(b))
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, convertArgs(a)...)
}

// Print is a wrapper for fmt.Print that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Print(spew.NewFormatter(a), spew.NewFormatter(b))
func Print(a ...interface{}) (n int, err error) {
	return fmt.Print(convertArgs(a)...)
}

// Printf is a wrapper for fmt.Printf that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Printf(format, spew.NewFormatter(a), spew.NewFormatter(b))
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, convertArgs(a)...)
}

// Println is a wrapper for fmt.Println that treats each argument as if it were
// passed with a default Formatter interface returned by NewFormatter.  It
// returns the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Println(spew.NewFormatter(a), spew.NewFormatter(b))
func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(convertArgs(a)...)
}

// convertArgs accepts a slice of arguments and returns a slice of the same
// length with each argument converted to a default spew Formatter interface.
func convertArgs(args []interface{}) (formatters []interface{}) {
	formatters = make([]interface{}, len(args))
	for index, arg := range args {
		formatters[index] = NewFormatter(arg)
	}
	return formatters
}

// SpewState provides a context which can have its own configuration options.
// The configuration options can be manipulated via the Config method.  The
// methods of SpewState are equivalent to the top-level functions.
//
// A SpewState does not need any special initialization, so new(SpewState) or
// just declaring a SpewState variable, is  sufficient to initialilize a
// SpewState using the default configuration options.
type SpewState struct {
	cs     *ConfigState
}

// Config returns a pointer to the active ConfigState for the SpewState
// instance.  Set the fields of the returned structure to the desired
// configuration settings for the instance.
func (s *SpewState) Config() (cs *ConfigState) {
	if s.cs == nil {
		cs := defaultConfig
		s.cs = &cs
	}
	return s.cs
}

// Errorf is a wrapper for fmt.Errorf that treats each argument as if it were
// passed with a Formatter interface returned by s.NewFormatter.  It returns
// the formatted string as a value that satisfies error.  See NewFormatter
// for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Errorf(format, s.NewFormatter(a), s.NewFormatter(b))
func (s *SpewState) Errorf(format string, a ...interface{}) (err error) {
	return fmt.Errorf(format, s.convertArgs(a)...)
}

// Fprint is a wrapper for fmt.Fprint that treats each argument as if it were
// passed with a Formatter interface returned by s.NewFormatter.  It returns
// the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprint(w, s.NewFormatter(a), s.NewFormatter(b))
func (s *SpewState) Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, s.convertArgs(a)...)
}

// Fprintf is a wrapper for fmt.Fprintf that treats each argument as if it were
// passed with a Formatter interface returned by s.NewFormatter.  It returns
// the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprintf(w, format, s.NewFormatter(a), s.NewFormatter(b))
func (s *SpewState) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format, s.convertArgs(a)...)
}

// Fprintln is a wrapper for fmt.Fprintln that treats each argument as if it
// passed with a Formatter interface returned by s.NewFormatter.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Fprintln(w, s.NewFormatter(a), s.NewFormatter(b))
func (s *SpewState) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, s.convertArgs(a)...)
}

// Print is a wrapper for fmt.Print that treats each argument as if it were
// passed with a Formatter interface returned by s.NewFormatter.  It returns
// the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Print(s.NewFormatter(a), s.NewFormatter(b))
func (s *SpewState) Print(a ...interface{}) (n int, err error) {
	return fmt.Print(s.convertArgs(a)...)
}

// Printf is a wrapper for fmt.Printf that treats each argument as if it were
// passed with a Formatter interface returned by s.NewFormatter.  It returns
// the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Printf(format, s.NewFormatter(a), s.NewFormatter(b))
func (s *SpewState) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, s.convertArgs(a)...)
}

// Println is a wrapper for fmt.Println that treats each argument as if it were
// passed with a Formatter interface returned by s.NewFormatter.  It returns
// the number of bytes written and any write error encountered.  See
// NewFormatter for formatting details.
//
// This function is shorthand for the following syntax:
//
//	fmt.Println(s.NewFormatter(a), s.NewFormatter(b))
func (s *SpewState) Println(a ...interface{}) (n int, err error) {
	return fmt.Println(s.convertArgs(a)...)
}

/*
NewFormatter returns a custom formatter that satisfies the fmt.Formatter
interface.  As a result, it integrates cleanly with standard fmt package
printing functions.  The formatter is useful for inline printing of smaller data
types similar to the standard %v format specifier.

The custom formatter only responds to the %v and %+v verb combinations.  Any
other variations such as %x, %q, and %#v will be sent to the the standard fmt
package for formatting.  In addition, the custom formatter ignores the width and
precision arguments (however they will still work on the format specifiers not
handled by the custom formatter).

Typically this function shouldn't be called directly.  It is much easier to make
use of the custom formatter by calling one of the convenience functions such as
s.Printf, s.Println, or s.Printf.
*/
func (s *SpewState) NewFormatter(v interface{}) fmt.Formatter {
	// The Config method creates the config state if needed, so call it instead
	// of using s.cs directly to ensure the zero value SpewState is sane.
	return newFormatter(s.Config(), v)
}

// Fdump formats and displays the passed arguments to io.Writer w.  It formats
// exactly the same as Dump.
func (s *SpewState) Fdump(w io.Writer, a ...interface{}) {
	// The Config method creates the config state if needed, so call it instead
	// of using s.cs directly to ensure the zero value SpewState is sane.
	fdump(s.Config(), w, a...)
}

/*
Dump displays the passed parameters to standard out with newlines, customizable
indentation, and additional debug information such as complete types and all
pointer addresses used to indirect to the final value.  It provides the
following features over the built-in printing facilities provided by the fmt
package:

	* Pointers are dereferenced and followed
	* Circular data structures are detected and handled properly
	* Custom error/Stringer interfaces are optionally invoked, including
	  on unexported types
	* Custom types which only implement the error/Stringer interfaces via
	  a pointer receiver are optionally invoked when passing non-pointer
	  variables

The configuration options are controlled by accessing the ConfigState associated
with s via the Config method.  See ConfigState for options documentation.

See Fdump if you would prefer dump to an arbitrary io.Writer.
*/
func (s *SpewState) Dump(a ...interface{}) {
	// The Config method creates the config state if needed, so call it instead
	// of using s.cs directly to ensure the zero value SpewState is sane.
	fdump(s.Config(), os.Stdout, a...)
}

// convertArgs accepts a slice of arguments and returns a slice of the same
// length with each argument converted to a spew Formatter interface using
// the ConfigState associated with s.
func (s *SpewState) convertArgs(args []interface{}) (formatters []interface{}) {
	// The Config method creates the config state if needed, so call it instead
	// of using s.cs directly to ensure the zero value SpewState is sane.
	cs := s.Config()
	formatters = make([]interface{}, len(args))
	for index, arg := range args {
		formatters[index] = newFormatter(cs, arg)
	}
	return formatters
}
