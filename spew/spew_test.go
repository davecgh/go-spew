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

// stringizeWants converts a slice of wanted test output into a format suitable
// for an test error message.
func stringizeWants(wants []string) string {
	s := ""
	for i, want := range wants {
		if i > 0 {
			s += fmt.Sprintf("want%d: %s", i+1, want)
		} else {
			s += "want: " + want
		}
	}
	return s
}

// testFailed returns whether or not a test failed by checking if the result
// of the test is in the slice of wanted strings.
func testFailed(result string, wants []string) bool {
	for _, want := range wants {
		if result == want {
			return false
		}
	}
	return true
}

// panicer is used to intentionally cause a panic for testing spew properly
// handles them
type panicer int

func (p panicer) String() string {
	panic("test panic")
}

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
	{fCSFdump, "", int8(127), "(int8) 127\n"},
	{fCSFprint, "", int16(32767), "32767"},
	{fCSFprintf, "%v", int32(2147483647), "2147483647"},
	{fCSFprintln, "", int(2147483647), "2147483647\n"},
	{fCSPrint, "", int64(9223372036854775807), "9223372036854775807"},
	{fCSPrintln, "", uint8(255), "255\n"},
	{fCSErrorf, "%#v", uint16(65535), "(uint16)65535"},
	{fCSNewFormatter, "%v", uint32(4294967295), "4294967295"},
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
	scs := spew.NewDefaultConfig()

	t.Logf("Running %d tests", len(spewTests))
	for i, test := range spewTests {
		buf := new(bytes.Buffer)
		switch test.f {
		case fCSFdump:
			scs.Fdump(buf, test.in)

		case fCSFprint:
			scs.Fprint(buf, test.in)

		case fCSFprintf:
			scs.Fprintf(buf, test.format, test.in)

		case fCSFprintln:
			scs.Fprintln(buf, test.in)

		case fCSPrint:
			b, err := redirStdout(func() { scs.Print(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fCSPrintln:
			b, err := redirStdout(func() { scs.Println(test.in) })
			if err != nil {
				t.Errorf("%v #%d %v", test.f, i, err)
				continue
			}
			buf.Write(b)

		case fCSErrorf:
			err := scs.Errorf(test.format, test.in)
			buf.WriteString(err.Error())

		case fCSNewFormatter:
			fmt.Fprintf(buf, test.format, scs.NewFormatter(test.in))

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
