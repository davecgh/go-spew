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

package spew_test

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

type Flag int

const (
	flagOne Flag = iota
	flagTwo
)

var flagStrings = map[Flag]string{
	flagOne: "flagOne",
	flagTwo: "flagTwo",
}

func (f Flag) String() string {
	if s, ok := flagStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown flag (%d)", int(f))
}

type Bar struct {
	flag Flag
	data uintptr
}

type Foo struct {
	unexportedField Bar
	ExportedField   map[interface{}]interface{}
}

// This example demonstrates how to use Dump to dump variables to stdout.
func ExampleDump() {
	// The following package level declarations are assumed for this example:
	/*
		type Flag int

		const (
			flagOne Flag = iota
			flagTwo
		)

		var flagStrings = map[Flag]string{
			flagOne: "flagOne",
			flagTwo: "flagTwo",
		}

		func (f Flag) String() string {
			if s, ok := flagStrings[f]; ok {
				return s
			}
			return fmt.Sprintf("Unknown flag (%d)", int(f))
		}

		type Bar struct {
			flag Flag
			data uintptr
		}

		type Foo struct {
			unexportedField Bar
			ExportedField   map[interface{}]interface{}
		}
	*/

	// Setup some sample data structures for the example.
	bar := Bar{Flag(flagTwo), uintptr(0)}
	s1 := Foo{bar, map[interface{}]interface{}{"one": true}}
	f := Flag(5)

	// Dump!
	spew.Dump(s1, f)

	// Output:
	// (spew_test.Foo) {
	//  unexportedField: (spew_test.Bar) {
	//   flag: (spew_test.Flag) flagTwo,
	//   data: (uintptr) <nil>
	//  },
	//  ExportedField: (map[interface {}]interface {}) {
	//   (string) "one": (bool) true
	//  }
	// }
	// (spew_test.Flag) Unknown flag (5)
	//
}

// This example demonstrates how to use Printf to display a variable with a
// format string and inline formatting.
func ExamplePrintf() {
	// Create a double pointer to a uint 8.
	ui8 := uint8(5)
	pui8 := &ui8
	ppui8 := &pui8

	// Create a circular data type.
	type circular struct {
		ui8 uint8
		c   *circular
	}
	c := circular{ui8: 1}
	c.c = &c

	// Print!
	spew.Printf("ppui8: %v\n", ppui8)
	spew.Printf("circular: %v\n", c)

	// Output:
	// ppui8: <**>5
	// circular: {1 <*>{1 <*><shown>}}
}
