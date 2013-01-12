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

// ConfigState is used to describe configuration options used by spew to format
// and display values.  There is currently only a single global instance, Config,
// that is used to control all Formatter and Dump functionality.  This state
// is designed so that it would be fairly simple to add the ability to have
// unique config per Formatter or dumpState instance if there is demand for
// such a feature.
type ConfigState struct {
	// MaxDepth controls the maximum number of levels to descend into nested
	// data structures.  The default, 0, means there is no limit.
	//
	// NOTE: Circular data structures are properly detected, so it is not
	// necessary to set this value unless you specifically want to limit deeply
	// nested data structures.
	MaxDepth int

	// Indent specifies the string to use for each indentation level.  It is
	// a single space by default.  If you would like more indentation, you might
	// set this to a tab with "\t" or perhaps two spaces with "  ".
	Indent string

	// DisableMethods specifies whether or not error and Stringer interfaces are
	// invoked for types that implement them.
	DisableMethods bool

	// DisablePointerMethods specifies whether or not to check for and invoke
	// error and Stringer interfaces on types which only accept a pointer
	// receiver when the current type is not a pointer.
	//
	// NOTE: This might be an unsafe action since calling one of these methods
	// with a pointer receiver could technically mutate the value, however,
	// in practice, types which choose to satisify an error or Stringer
	// interface with a pointer receiver should not be mutating their state
	// inside these interface methods.
	DisablePointerMethods bool
}

// Config is the active configuration in use by spew.  The configuration
// can be changed by modifying the contents of spew.Config.
var Config ConfigState = ConfigState{Indent: " "}

var defaultConfig = ConfigState{Indent: " "}