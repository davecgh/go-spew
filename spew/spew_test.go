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
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"
	"testing"
)

var scsDefault = spew.NewDefaultConfig()

// spewFunc is used to identify which public function of the spew package or
// ConfigState a test applies to.
type spewFunc int

const (
	fCSFdump spewFunc = iota
	fCSFprint
	fCSFprintf
	fCSFprintln
	fCSPrint
	fCSPrintln
	fCSErrorf
	fCSNewFormatter
	fErrorf
	fFprint
	fFprintln
	fPrint
	fPrintln
)

// Map of spewFunc values to names for pretty printing.
var spewFuncStrings = map[spewFunc]string{
	fCSFdump:        "ConfigState.Fdump",
	fCSFprint:       "ConfigState.Fprint",
	fCSFprintf:      "ConfigState.Fprintf",
	fCSFprintln:     "ConfigState.Fprintln",
	fCSPrint:        "ConfigState.Print",
	fCSPrintln:      "ConfigState.Println",
	fCSErrorf:       "ConfigState.Errorf",
	fCSNewFormatter: "ConfigState.NewFormatter",
	fErrorf:         "spew.Errorf",
	fFprint:         "spew.Fprint",
	fFprintln:       "spew.Fprintln",
	fPrint:          "spew.Print",
	fPrintln:        "spew.Println",
}

func (f spewFunc) String() string {
	if s, ok := spewFuncStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown spewFunc (%d)", int(f))
}

// spewTest is used to describe a test to be performed against the public
// functions of the spew package or ConfigState.
type spewTest struct {
	cs     *spew.ConfigState
	f      spewFunc
	format string
	in     interface{}
	want   string
}

// spewTests houses the tests to be performed against the public functions of
// the spew package and ConfigState.
//
// These tests are only intended to ensure the public functions are exercised
// and are intentionally not exhaustive of types.  The exhaustive type
// tests are handled in the dump and format tests.
var spewTests = []spewTest{
	{scsDefault, fCSFdump, "", int8(127), "(int8) 127\n"},
	{scsDefault, fCSFprint, "", int16(32767), "32767"},
	{scsDefault, fCSFprintf, "%v", int32(2147483647), "2147483647"},
	{scsDefault, fCSFprintln, "", int(2147483647), "2147483647\n"},
	{scsDefault, fCSPrint, "", int64(9223372036854775807), "9223372036854775807"},
	{scsDefault, fCSPrintln, "", uint8(255), "255\n"},
	{scsDefault, fCSErrorf, "%#v", uint16(65535), "(uint16)65535"},
	{scsDefault, fCSNewFormatter, "%v", uint32(4294967295), "4294967295"},
	{scsDefault, fErrorf, "%v", uint64(18446744073709551615), "18446744073709551615"},
	{scsDefault, fFprint, "", float32(3.14), "3.14"},
	{scsDefault, fFprintln, "", float64(6.28), "6.28\n"},
	{scsDefault, fPrint, "", true, "true"},
	{scsDefault, fPrintln, "", false, "false\n"},
}

// redirStdout is a helper function to return the standard output from f as a
// byte slice.
func redirStdout(f func()) ([]byte, error) {
	tempFile, err := ioutil.TempFile("", "ss-test")
	if err != nil {
		return nil, err
	}
	fileName := tempFile.Name()
	defer os.Remove(fileName) // Ignore error

	origStdout := os.Stdout
	os.Stdout = tempFile
	f()
	os.Stdout = origStdout
	tempFile.Close()

	return ioutil.ReadFile(fileName)
}

// TestSpew executes all of the tests described by spewTests.
func TestSpew(t *testing.T) {
	t.Logf("Running %d tests", len(spewTests))
	for i, test := range spewTests {
		buf := new(bytes.Buffer)
		switch test.f {
		case fCSFdump:
			test.cs.Fdump(buf, test.in)

		case fCSFprint:
			test.cs.Fprint(buf, test.in)

		case fCSFprintf:
			test.cs.Fprintf(buf, test.format, test.in)

		case fCSFprintln:
			test.cs.Fprintln(buf, test.in)

		case fCSPrint:
			b, err := redirStdout(func() { test.cs.Print(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fCSPrintln:
			b, err := redirStdout(func() { test.cs.Println(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fCSErrorf:
			err := test.cs.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fCSNewFormatter:
			fmt.Fprintf(buf, test.format, test.cs.NewFormatter(test.in))

		case fErrorf:
			err := spew.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fFprint:
			spew.Fprint(buf, test.in)

		case fFprintln:
			spew.Fprintln(buf, test.in)

		case fPrint:
			b, err := redirStdout(func() { spew.Print(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fPrintln:
			b, err := redirStdout(func() { spew.Println(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		default:
			t.Errorf("%v #%d unrecognized function", test.f, i)
			continue
		}
		s := buf.String()
		if test.want != s {
			t.Errorf("ConfigState #%d\n got: %s want: %s", i, s, test.want)
			continue
		}
	}
}
