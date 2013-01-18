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

// spewFunc is used to identify which public function of the spew package or
// SpewState a test applies to.
type spewFunc int

const (
	fSSFdump spewFunc = iota
	fSSFprint
	fSSFprintf
	fSSFprintln
	fSSPrint
	fSSPrintln
	fSSErrorf
	fSSNewFormatter
	fErrorf
	fFprint
	fFprintln
	fPrint
	fPrintln
)

// Map of spewFunc values to names for pretty printing.
var spewFuncStrings = map[spewFunc]string{
	fSSFdump:        "SpewState.Fdump",
	fSSFprint:       "SpewState.Fprint",
	fSSFprintf:      "SpewState.Fprintf",
	fSSFprintln:     "SpewState.Fprintln",
	fSSPrint:        "SpewState.Print",
	fSSPrintln:      "SpewState.Println",
	fSSErrorf:       "SpewState.Errorf",
	fSSNewFormatter: "SpewState.NewFormatter",
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
// functions of the spew package or SpewState.
type spewTest struct {
	f      spewFunc
	format string
	in     interface{}
	want   string
}

// spewTests houses the tests to be performed against the public functions of
// the spew package and SpewState.
//
// These tests are only intended to ensure the public functions are exercised
// and are intentionally not exhaustive of types.  The exhaustive type
// tests are handled in the dump and format tests.
var spewTests = []spewTest{
	{fSSFdump, "", int8(127), "(int8) 127\n"},
	{fSSFprint, "", int16(32767), "32767"},
	{fSSFprintf, "%v", int32(2147483647), "2147483647"},
	{fSSFprintln, "", int(2147483647), "2147483647\n"},
	{fSSPrint, "", int64(9223372036854775807), "9223372036854775807"},
	{fSSPrintln, "", uint8(255), "255\n"},
	{fSSErrorf, "%#v", uint16(65535), "(uint16)65535"},
	{fSSNewFormatter, "%v", uint32(4294967295), "4294967295"},
	{fErrorf, "%v", uint64(18446744073709551615), "18446744073709551615"},
	{fFprint, "", float32(3.14), "3.14"},
	{fFprintln, "", float64(6.28), "6.28\n"},
	{fPrint, "", true, "true"},
	{fPrintln, "", false, "false\n"},
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
	ss := new(spew.SpewState)

	t.Logf("Running %d tests", len(spewTests))
	for i, test := range spewTests {
		buf := new(bytes.Buffer)
		switch test.f {
		case fSSFdump:
			ss.Fdump(buf, test.in)

		case fSSFprint:
			ss.Fprint(buf, test.in)

		case fSSFprintf:
			ss.Fprintf(buf, test.format, test.in)

		case fSSFprintln:
			ss.Fprintln(buf, test.in)

		case fSSPrint:
			b, err := redirStdout(func() { ss.Print(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fSSPrintln:
			b, err := redirStdout(func() { ss.Println(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fSSErrorf:
			err := ss.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fSSNewFormatter:
			fmt.Fprintf(buf, test.format, ss.NewFormatter(test.in))

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
			t.Errorf("SpewState #%d\n got: %s want: %s", i, s, test.want)
			continue
		}
	}
}
