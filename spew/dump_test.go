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
Test Summary:
NOTE: For each test, a nil pointer, a single pointer and double pointer to the
base test element are also tested to ensure proper indirection across all types.

- Max int8, int16, int32, int64, int
- Max uint8, uint16, uint32, uint64, uint
- Boolean true and false
- Standard complex64 and complex128
- Array containing standard ints
- Array containing type with custom formatter on pointer receiver only
- Slice containing standard float32 values
- Slice containing type with custom formatter on pointer receiver only
- Standard string
- Nil interface
- Map with string keys and int vals
- Map with custom formatter type on pointer receiver only keys and vals
- Map with interface keys and values
- Struct with primitives
- Struct that contains another struct
- Struct that contains custom type with Stringer pointer interface via both
  exported and unexported fields
- Uintptr to 0 (null pointer)
- Uintptr address of real variable
- Unsafe.Pointer to 0 (null pointer)
- Unsafe.Pointer to address of real variable
- Nil channel
- Standard int channel
- Function with no params and no returns
- Function with param and no returns
- Function with multiple params and multiple returns
- Struct that is circular through self referencing
- Structs that are circular through cross referencing
- Structs that are indirectly circular
*/

package spew_test

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
	"unsafe"
)

// custom type to test Stinger interface on pointer receiver.
type pstringer string

// String implements the Stringer interface for testing invocation of custom
// stringers on types with only pointer receivers.
func (s *pstringer) String() string {
	return "stringer " + string(*s)
}

// xref1 and xref2 are cross referencing structs for testing circular reference
//  detection.
type xref1 struct {
	ps2 *xref2
}
type xref2 struct {
	ps1 *xref1
}

// indirCir1, indirCir2, and indirCir3 are used to generate an indirect circular
// reference for testing detection.
type indirCir1 struct {
	ps2 *indirCir2
}
type indirCir2 struct {
	ps3 *indirCir3
}
type indirCir3 struct {
	ps1 *indirCir1
}

// dumpTest is used to describe a test to be perfomed against the Dump method.
type dumpTest struct {
	in   interface{}
	want string
}

// dumpTests houses all of the tests to be performed against the Dump method.
var dumpTests = make([]dumpTest, 0)

// addDumpTest is a helper method to append the passed input and desired result
// to dumpTests
func addDumpTest(in interface{}, want string) {
	test := dumpTest{in, want}
	dumpTests = append(dumpTests, test)
}

func addIntDumpTests() {
	// Max int8.
	v := int8(127)
	nv := (*int8)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "int8"
	vs := "127"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Max int16.
	v2 := int16(32767)
	nv2 := (*int16)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "int16"
	v2s := "32767"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")

	// Max int32.
	v3 := int32(2147483647)
	nv3 := (*int32)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "int32"
	v3s := "2147483647"
	addDumpTest(v3, "("+v3t+") "+v3s+"\n")
	addDumpTest(pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")\n")
	addDumpTest(&pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(<nil>)\n")

	// Max int64.
	v4 := int64(9223372036854775807)
	nv4 := (*int64)(nil)
	pv4 := &v4
	v4Addr := fmt.Sprintf("%p", pv4)
	pv4Addr := fmt.Sprintf("%p", &pv4)
	v4t := "int64"
	v4s := "9223372036854775807"
	addDumpTest(v4, "("+v4t+") "+v4s+"\n")
	addDumpTest(pv4, "(*"+v4t+")("+v4Addr+")("+v4s+")\n")
	addDumpTest(&pv4, "(**"+v4t+")("+pv4Addr+"->"+v4Addr+")("+v4s+")\n")
	addDumpTest(nv4, "(*"+v4t+")(<nil>)\n")

	// Max int.
	v5 := int(2147483647)
	nv5 := (*int)(nil)
	pv5 := &v5
	v5Addr := fmt.Sprintf("%p", pv5)
	pv5Addr := fmt.Sprintf("%p", &pv5)
	v5t := "int"
	v5s := "2147483647"
	addDumpTest(v5, "("+v5t+") "+v5s+"\n")
	addDumpTest(pv5, "(*"+v5t+")("+v5Addr+")("+v5s+")\n")
	addDumpTest(&pv5, "(**"+v5t+")("+pv5Addr+"->"+v5Addr+")("+v5s+")\n")
	addDumpTest(nv5, "(*"+v5t+")(<nil>)\n")
}

func addUintDumpTests() {
	// Max uint8.
	v := uint8(255)
	nv := (*uint8)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "uint8"
	vs := "255"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Max uint16.
	v2 := uint16(65535)
	nv2 := (*uint16)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "uint16"
	v2s := "65535"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")

	// Max uint32.
	v3 := uint32(4294967295)
	nv3 := (*uint32)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "uint32"
	v3s := "4294967295"
	addDumpTest(v3, "("+v3t+") "+v3s+"\n")
	addDumpTest(pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")\n")
	addDumpTest(&pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(<nil>)\n")

	// Max uint64.
	v4 := uint64(18446744073709551615)
	nv4 := (*uint64)(nil)
	pv4 := &v4
	v4Addr := fmt.Sprintf("%p", pv4)
	pv4Addr := fmt.Sprintf("%p", &pv4)
	v4t := "uint64"
	v4s := "18446744073709551615"
	addDumpTest(v4, "("+v4t+") "+v4s+"\n")
	addDumpTest(pv4, "(*"+v4t+")("+v4Addr+")("+v4s+")\n")
	addDumpTest(&pv4, "(**"+v4t+")("+pv4Addr+"->"+v4Addr+")("+v4s+")\n")
	addDumpTest(nv4, "(*"+v4t+")(<nil>)\n")

	// Max uint.
	v5 := uint(4294967295)
	nv5 := (*uint)(nil)
	pv5 := &v5
	v5Addr := fmt.Sprintf("%p", pv5)
	pv5Addr := fmt.Sprintf("%p", &pv5)
	v5t := "uint"
	v5s := "4294967295"
	addDumpTest(v5, "("+v5t+") "+v5s+"\n")
	addDumpTest(pv5, "(*"+v5t+")("+v5Addr+")("+v5s+")\n")
	addDumpTest(&pv5, "(**"+v5t+")("+pv5Addr+"->"+v5Addr+")("+v5s+")\n")
	addDumpTest(nv5, "(*"+v5t+")(<nil>)\n")
}

func addBoolDumpTests() {
	// Boolean true.
	v := bool(true)
	nv := (*bool)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "bool"
	vs := "true"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Boolean false.
	v2 := bool(false)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "bool"
	v2s := "false"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
}

func addFloatDumpTests() {
	// Standard float32.
	v := float32(3.1415)
	nv := (*float32)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "float32"
	vs := "3.1415"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Standard float64.
	v2 := float64(3.1415926)
	nv2 := (*float64)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "float64"
	v2s := "3.1415926"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")
}

func addComplexDumpTests() {
	// Standard complex64.
	v := complex(float32(6), -2)
	nv := (*complex64)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "complex64"
	vs := "(6-2i)"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Standard complex128.
	v2 := complex(float64(-6), 2)
	nv2 := (*complex128)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "complex128"
	v2s := "(-6+2i)"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")
}

func addArrayDumpTests() {
	// Array containing standard ints.
	v := [3]int{1, 2, 3}
	nv := (*[3]int)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "int"
	vs := "{\n (" + vt + ") 1,\n (" + vt + ") 2,\n (" + vt + ") 3\n}"
	addDumpTest(v, "([3]"+vt+") "+vs+"\n")
	addDumpTest(pv, "(*[3]"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**[3]"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*[3]"+vt+")(<nil>)\n")

	// Array containing type with custom formatter on pointer receiver only.
	v2 := [3]pstringer{"1", "2", "3"}
	nv2 := (*[3]pstringer)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.pstringer"
	v2s := "{\n (" + v2t + ") stringer 1,\n (" + v2t + ") stringer 2,\n (" +
		v2t + ") stringer 3\n}"
	addDumpTest(v2, "([3]"+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*[3]"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**[3]"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*[3]"+v2t+")(<nil>)\n")
}

func addSliceDumpTests() {
	// Slice containing standard float32 values.
	v := []float32{3.14, 6.28, 12.56}
	nv := (*[]float32)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "float32"
	vs := "{\n (" + vt + ") 3.14,\n (" + vt + ") 6.28,\n (" + vt + ") 12.56\n}"
	addDumpTest(v, "([]"+vt+") "+vs+"\n")
	addDumpTest(pv, "(*[]"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**[]"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*[]"+vt+")(<nil>)\n")

	// Slice containing type with custom formatter on pointer receiver only.
	v2 := []pstringer{"1", "2", "3"}
	nv2 := (*[]pstringer)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.pstringer"
	v2s := "{\n (" + v2t + ") stringer 1,\n (" + v2t + ") stringer 2,\n (" +
		v2t + ") stringer 3\n}"
	addDumpTest(v2, "([]"+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*[]"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**[]"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*[]"+v2t+")(<nil>)\n")
}

func addStringDumpTests() {
	// Standard string.
	v := "test"
	nv := (*string)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "string"
	vs := "\"test\""
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")
}

func addNilInterfaceDumpTests() {
	// Nil interface.
	var v interface{}
	nv := (*interface{})(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "interface {}"
	vs := "<nil>"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")
}

func addMapDumpTests() {
	// Map with string keys and int vals.
	v := map[string]int{"one": 1}
	nv := (*map[string]int)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "map[string]int"
	vt1 := "string"
	vt2 := "int"
	vs := "{\n (" + vt1 + ") \"one\": (" + vt2 + ") 1\n}"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Map with custom formatter type on pointer receiver only keys and vals.
	v2 := map[pstringer]pstringer{"one": "1"}
	nv2 := (*map[pstringer]pstringer)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "map[spew_test.pstringer]spew_test.pstringer"
	v2t1 := "spew_test.pstringer"
	v2t2 := "spew_test.pstringer"
	v2s := "{\n (" + v2t1 + ") stringer one: (" + v2t2 + ") stringer 1\n}"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")

	// Map with interface keys and values.
	v3 := map[interface{}]interface{}{"one": 1}
	nv3 := (*map[interface{}]interface{})(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "map[interface {}]interface {}"
	v3t1 := "string"
	v3t2 := "int"
	v3s := "{\n (" + v3t1 + ") \"one\": (" + v3t2 + ") 1\n}"
	addDumpTest(v3, "("+v3t+") "+v3s+"\n")
	addDumpTest(pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")\n")
	addDumpTest(&pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(<nil>)\n")
}

func addStructDumpTests() {
	// Struct with primitives.
	type s1 struct {
		a int8
		b uint8
	}
	v := s1{127, 255}
	nv := (*s1)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "spew_test.s1"
	vt2 := "int8"
	vt3 := "uint8"
	vs := "{\n a: (" + vt2 + ") 127,\n b: (" + vt3 + ") 255\n}"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Struct that contains another struct.
	type s2 struct {
		s1 s1
		b  bool
	}
	v2 := s2{s1{127, 255}, true}
	nv2 := (*s2)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.s2"
	v2t2 := "spew_test.s1"
	v2t3 := "int8"
	v2t4 := "uint8"
	v2t5 := "bool"
	v2s := "{\n s1: (" + v2t2 + ") {\n  a: (" + v2t3 + ") 127,\n  b: (" +
		v2t4 + ") 255\n },\n b: (" + v2t5 + ") true\n}"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")

	// Struct that contains custom type with Stringer pointer interface via both
	// exported and unexported fields.
	type s3 struct {
		s pstringer
		S pstringer
	}
	v3 := s3{"test", "test2"}
	nv3 := (*s3)(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "spew_test.s3"
	v3t2 := "spew_test.pstringer"
	v3s := "{\n s: (" + v3t2 + ") stringer test,\n S: (" + v3t2 +
		") stringer test2\n}"
	addDumpTest(v3, "("+v3t+") "+v3s+"\n")
	addDumpTest(pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")\n")
	addDumpTest(&pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(<nil>)\n")
}

func addUintptrDumpTests() {
	// Null pointer.
	v := uintptr(0)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "uintptr"
	vs := "<nil>"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")

	// Address of real variable.
	i := 1
	v2 := uintptr(unsafe.Pointer(&i))
	nv2 := (*uintptr)(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "uintptr"
	v2s := fmt.Sprintf("%p", &i)
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")
}

func addUnsafePointerDumpTests() {
	// Null pointer.
	v := unsafe.Pointer(uintptr(0))
	nv := (*unsafe.Pointer)(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "unsafe.Pointer"
	vs := "<nil>"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Address of real variable.
	i := 1
	v2 := unsafe.Pointer(&i)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "unsafe.Pointer"
	v2s := fmt.Sprintf("%p", &i)
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")
}

func addChanDumpTests() {
	// Nil channel.
	var v chan int
	pv := &v
	nv := (*chan int)(nil)
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "chan int"
	vs := "<nil>"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Real channel.
	v2 := make(chan int)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "chan int"
	v2s := fmt.Sprintf("%p", v2)
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
}

func addFuncDumpTests() {
	// Function with no params and no returns.
	v := addIntDumpTests
	nv := (*func())(nil)
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "func()"
	vs := fmt.Sprintf("%p", v)
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(<nil>)\n")

	// Function with param and no returns.
	v2 := TestDump
	nv2 := (*func(*testing.T))(nil)
	pv2 := &v2
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "func(*testing.T)"
	v2s := fmt.Sprintf("%p", v2)
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(<nil>)\n")

	// Function with multiple params and multiple returns.
	var v3 = func(i int, s string) (b bool, err error) {
		return true, nil
	}
	nv3 := (*func(int, string) (bool, error))(nil)
	pv3 := &v3
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "func(int, string) (bool, error)"
	v3s := fmt.Sprintf("%p", v3)
	addDumpTest(v3, "("+v3t+") "+v3s+"\n")
	addDumpTest(pv3, "(*"+v3t+")("+v3Addr+")("+v3s+")\n")
	addDumpTest(&pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(<nil>)\n")
}

func addCircularDumpTests() {
	// Struct that is circular through self referencing.
	type circular struct {
		c *circular
	}
	v := circular{nil}
	v.c = &v
	pv := &v
	vAddr := fmt.Sprintf("%p", pv)
	pvAddr := fmt.Sprintf("%p", &pv)
	vt := "spew_test.circular"
	vs := "{\n c: (*" + vt + ")(" + vAddr + ")({\n  c: (*" + vt + ")(" +
		vAddr + ")(<already shown>)\n })\n}"
	vs2 := "{\n c: (*" + vt + ")(" + vAddr + ")(<already shown>)\n}"
	addDumpTest(v, "("+vt+") "+vs+"\n")
	addDumpTest(pv, "(*"+vt+")("+vAddr+")("+vs2+")\n")
	addDumpTest(&pv, "(**"+vt+")("+pvAddr+"->"+vAddr+")("+vs2+")\n")

	// Structs that are circular through cross referencing.
	v2 := xref1{nil}
	ts2 := xref2{&v2}
	v2.ps2 = &ts2
	pv2 := &v2
	ts2Addr := fmt.Sprintf("%p", &ts2)
	v2Addr := fmt.Sprintf("%p", pv2)
	pv2Addr := fmt.Sprintf("%p", &pv2)
	v2t := "spew_test.xref1"
	v2t2 := "spew_test.xref2"
	v2s := "{\n ps2: (*" + v2t2 + ")(" + ts2Addr + ")({\n  ps1: (*" + v2t +
		")(" + v2Addr + ")({\n   ps2: (*" + v2t2 + ")(" + ts2Addr +
		")(<already shown>)\n  })\n })\n}"
	v2s2 := "{\n ps2: (*" + v2t2 + ")(" + ts2Addr + ")({\n  ps1: (*" + v2t +
		")(" + v2Addr + ")(<already shown>)\n })\n}"
	addDumpTest(v2, "("+v2t+") "+v2s+"\n")
	addDumpTest(pv2, "(*"+v2t+")("+v2Addr+")("+v2s2+")\n")
	addDumpTest(&pv2, "(**"+v2t+")("+pv2Addr+"->"+v2Addr+")("+v2s2+")\n")

	// Structs that are indirectly circular.
	v3 := indirCir1{nil}
	tic2 := indirCir2{nil}
	tic3 := indirCir3{&v3}
	tic2.ps3 = &tic3
	v3.ps2 = &tic2
	pv3 := &v3
	tic2Addr := fmt.Sprintf("%p", &tic2)
	tic3Addr := fmt.Sprintf("%p", &tic3)
	v3Addr := fmt.Sprintf("%p", pv3)
	pv3Addr := fmt.Sprintf("%p", &pv3)
	v3t := "spew_test.indirCir1"
	v3t2 := "spew_test.indirCir2"
	v3t3 := "spew_test.indirCir3"
	v3s := "{\n ps2: (*" + v3t2 + ")(" + tic2Addr + ")({\n  ps3: (*" + v3t3 +
		")(" + tic3Addr + ")({\n   ps1: (*" + v3t + ")(" + v3Addr +
		")({\n    ps2: (*" + v3t2 + ")(" + tic2Addr +
		")(<already shown>)\n   })\n  })\n })\n}"
	v3s2 := "{\n ps2: (*" + v3t2 + ")(" + tic2Addr + ")({\n  ps3: (*" + v3t3 +
		")(" + tic3Addr + ")({\n   ps1: (*" + v3t + ")(" + v3Addr +
		")(<already shown>)\n  })\n })\n}"
	addDumpTest(v3, "("+v3t+") "+v3s+"\n")
	addDumpTest(pv3, "(*"+v3t+")("+v3Addr+")("+v3s2+")\n")
	addDumpTest(&pv3, "(**"+v3t+")("+pv3Addr+"->"+v3Addr+")("+v3s2+")\n")
}

// TestDump executes all of the tests described by dumpTests.
func TestDump(t *testing.T) {
	// Setup tests.
	addIntDumpTests()
	addUintDumpTests()
	addBoolDumpTests()
	addFloatDumpTests()
	addComplexDumpTests()
	addArrayDumpTests()
	addSliceDumpTests()
	addStringDumpTests()
	addNilInterfaceDumpTests()
	addMapDumpTests()
	addStructDumpTests()
	addUintptrDumpTests()
	addUnsafePointerDumpTests()
	addChanDumpTests()
	addFuncDumpTests()
	addCircularDumpTests()

	t.Logf("Running %d tests", len(dumpTests))
	for i, test := range dumpTests {
		buf := new(bytes.Buffer)
		spew.Fdump(buf, test.in)
		s := buf.String()
		if test.want != s {
			t.Errorf("Dump #%d\n got: %s want: %s", i, s, test.want)
			continue
		}
	}
}
