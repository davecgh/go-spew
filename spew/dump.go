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
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

// dumpState contains information about the state of a dump operation.
type dumpState struct {
	w              io.Writer
	depth          int
	pointers       map[uintptr]int
	ignoreNextType bool
	ignoreNextPad  bool
}

// pad performs indentation according to the depth level and Config.Indent
// option.
func (d *dumpState) pad() {
	if d.ignoreNextPad {
		d.ignoreNextPad = false
		return
	}
	d.w.Write(bytes.Repeat([]byte(Config.Indent), d.depth))
}

// dumpPtr handles formatting of pointers by indirecting them as necessary.
func (d *dumpState) dumpPtr(v reflect.Value) {
	// Remove pointers at or below the current depth from map used to detect
	// circular refs.
	for k, depth := range d.pointers {
		if depth >= d.depth {
			delete(d.pointers, k)
		}
	}

	// Keep list of all dereferenced pointers to show later.
	pointerChain := make([]uintptr, 0)

	// Figure out how many levels of indirection there are by derferencing
	// pointers and unpacking interfaces down the chain while detecting circular
	// references.
	nilFound := false
	cycleFound := false
	indirects := 0
	ve := v
	for ve.Kind() == reflect.Ptr {
		indirects++
		if ve.IsNil() {
			nilFound = true
			break
		}
		addr := ve.Pointer()
		pointerChain = append(pointerChain, addr)
		if pd, ok := d.pointers[addr]; ok && pd < d.depth {
			cycleFound = true
			indirects--
			break
		}
		d.pointers[addr] = d.depth

		ve = ve.Elem()
		if ve.Kind() == reflect.Interface {
			if ve.IsNil() {
				nilFound = true
				break
			}
			ve = ve.Elem()
		}
	}

	// Display type information.
	d.w.Write(openParenBytes)
	d.w.Write(bytes.Repeat(asteriskBytes, indirects))
	d.w.Write([]byte(ve.Type().String()))
	d.w.Write(closeParenBytes)

	// Display pointer information.
	d.w.Write(openParenBytes)
	for i, addr := range pointerChain {
		if i > 0 {
			d.w.Write(pointerChainBytes)
		}
		printHexPtr(d.w, addr)
	}
	d.w.Write(closeParenBytes)

	// Display dereferenced value.
	d.w.Write(openParenBytes)
	switch {
	case nilFound == true:
		d.w.Write(nilAngleBytes)

	case cycleFound == true:
		d.w.Write(circularBytes)

	default:
		d.ignoreNextType = true
		d.dump(ve)
	}
	d.w.Write(closeParenBytes)
}

// dump is the main workhorse for dumping a value.  It uses the passed reflect
// value to figure out what kind of object we are dealing with and formats it
// appropriately.  It is a recursive function, however circular data structures
// are detected and handled properly.
func (d *dumpState) dump(v reflect.Value) {
	// Handle pointers specially.
	kind := v.Kind()
	if kind == reflect.Ptr {
		d.pad()
		d.dumpPtr(v)
		return
	}

	// Print type information unless already handled elsewhere.
	if !d.ignoreNextType {
		d.pad()
		d.w.Write(openParenBytes)
		d.w.Write([]byte(v.Type().String()))
		d.w.Write(closeParenBytes)
		d.w.Write(spaceBytes)
	}
	d.ignoreNextType = false

	// Call error/Stringer interfaces if they exist and the handle methods flag
	// is enabled
	if !Config.DisableMethods {
		if (kind != reflect.Invalid) && (kind != reflect.Interface) {
			if handled := handleMethods(d.w, v); handled {
				return
			}
		}
	}

	switch kind {
	case reflect.Invalid:
		d.w.Write(invalidAngleBytes)

	case reflect.Bool:
		printBool(d.w, v.Bool())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		printInt(d.w, v.Int())

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		printUint(d.w, v.Uint())

	case reflect.Float32:
		printFloat(d.w, v.Float(), 32)

	case reflect.Float64:
		printFloat(d.w, v.Float(), 64)

	case reflect.Complex64:
		printComplex(d.w, v.Complex(), 32)

	case reflect.Complex128:
		printComplex(d.w, v.Complex(), 64)

	case reflect.Array, reflect.Slice:
		d.w.Write(openBraceNewlineBytes)
		d.depth++
		if (Config.MaxDepth != 0) && (d.depth > Config.MaxDepth) {
			d.pad()
			d.w.Write(maxNewlineBytes)
		} else {
			numEntries := v.Len()
			for i := 0; i < numEntries; i++ {
				d.dump(unpackValue(v.Index(i)))
				if i < (numEntries - 1) {
					d.w.Write(commaNewlineBytes)
				} else {
					d.w.Write(newlineBytes)
				}
			}
		}
		d.depth--
		d.pad()
		d.w.Write(closeBraceBytes)

	case reflect.String:
		d.w.Write([]byte(strconv.Quote(v.String())))

	case reflect.Interface:
		// Do nothing.  We should never get here due to unpackValue calls.

	case reflect.Ptr:
		// Do nothing.  We should never get here since pointer have already
		// been handled above.

	case reflect.Map:
		d.w.Write(openBraceNewlineBytes)
		d.depth++
		if (Config.MaxDepth != 0) && (d.depth > Config.MaxDepth) {
			d.pad()
			d.w.Write(maxNewlineBytes)
		} else {
			numEntries := v.Len()
			keys := v.MapKeys()
			for i, key := range keys {
				d.dump(unpackValue(key))
				d.w.Write(colonSpaceBytes)
				d.ignoreNextPad = true
				d.dump(unpackValue(v.MapIndex(key)))
				if i < (numEntries - 1) {
					d.w.Write(commaNewlineBytes)
				} else {
					d.w.Write(newlineBytes)
				}
			}
		}
		d.depth--
		d.pad()
		d.w.Write(closeBraceBytes)

	case reflect.Struct:
		d.w.Write(openBraceNewlineBytes)
		d.depth++
		if (Config.MaxDepth != 0) && (d.depth > Config.MaxDepth) {
			d.pad()
			d.w.Write(maxNewlineBytes)
		} else {
			vt := v.Type()
			numFields := v.NumField()
			for i := 0; i < numFields; i++ {
				d.pad()
				vtf := vt.Field(i)
				d.w.Write([]byte(vtf.Name))
				d.w.Write(colonSpaceBytes)
				d.ignoreNextPad = true
				d.dump(unpackValue(v.Field(i)))
				if i < (numFields - 1) {
					d.w.Write(commaNewlineBytes)
				} else {
					d.w.Write(newlineBytes)
				}
			}
		}
		d.depth--
		d.pad()
		d.w.Write(closeBraceBytes)

	case reflect.Uintptr:
		printHexPtr(d.w, uintptr(v.Uint()))

	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		printHexPtr(d.w, v.Pointer())

	// There were not any other types at the time this code was written, but
	// fall back to letting the default fmt package handle it in case any new
	// types are added.
	default:
		if v.CanInterface() {
			fmt.Fprintf(d.w, "%v", v.Interface())
		} else {
			fmt.Fprintf(d.w, "%v", v.String())
		}
	}
}

// Fdump formats and displays the passed arguments to io.Writer w.  It formats
// exactly the same as Dump.
func Fdump(w io.Writer, a ...interface{}) {
	for _, arg := range a {
		if arg == nil {
			w.Write(interfaceBytes)
			w.Write(nilAngleBytes)
			w.Write(newlineBytes)
			continue
		}

		d := dumpState{w: w}
		d.pointers = make(map[uintptr]int)
		d.dump(reflect.ValueOf(arg))
		d.w.Write(newlineBytes)
	}
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

The configuration options are controlled by an exported package global,
spew.Config.  See ConfigState for options documentation.

See Fdump if you would prefer dump to an arbitrary io.Writer.
*/
func Dump(a ...interface{}) {
	Fdump(os.Stdout, a...)
}
