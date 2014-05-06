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
	"reflect"
	"sort"
	"strconv"
	"unsafe"
)

// offsetPtr, offsetScalar, and offsetFlag are the offsets for the internal
// reflect.Value fields.
var offsetPtr, offsetScalar, offsetFlag uintptr

// reflectValueOld mirrors the struct layout of the reflect package Value type
// before golang commit ecccf07e7f9d.
var reflectValueOld struct {
	typ  unsafe.Pointer
	val  unsafe.Pointer
	flag uintptr
}

// reflectValueNew mirrors the struct layout of the reflect package Value type
// after golang commit ecccf07e7f9d.
var reflectValueNew struct {
	typ    unsafe.Pointer
	ptr    unsafe.Pointer
	scalar uintptr
	flag   uintptr
}

func init() {
	// Older versions of reflect.Value stored small integers directly in the
	// ptr field (which is named val in the older versions).  Newer versions
	// added a new field named scalar for this purpose which unfortuantely
	// comes before the flag field.  Further the new field is before the
	// flag field, so the offset of the flag field is different as well.
	// This code constructs a new reflect.Value from a known small integer
	// and checks if the val field within it matches.  When it matches, the
	// old style reflect.Value is being used.  Otherwise it's the new style.
	v := 0xf00
	vv := reflect.ValueOf(v)
	upv := unsafe.Pointer(uintptr(unsafe.Pointer(&vv)) +
		unsafe.Offsetof(reflectValueOld.val))

	// Assume the old style by default.
	offsetPtr = unsafe.Offsetof(reflectValueOld.val)
	offsetScalar = 0
	offsetFlag = unsafe.Offsetof(reflectValueOld.flag)

	// Use the new style offsets if the ptr field doesn't match the value
	// since it must be in the new scalar field.
	if int(*(*uintptr)(upv)) != v {
		offsetPtr = unsafe.Offsetof(reflectValueNew.ptr)
		offsetScalar = unsafe.Offsetof(reflectValueNew.scalar)
		offsetFlag = unsafe.Offsetof(reflectValueNew.flag)
	}
}

// flagIndir indicates whether the value field of a reflect.Value is the actual
// data or a pointer to the data.
const flagIndir = 1 << 1

// unsafeReflectValue converts the passed reflect.Value into a one that bypasses
// the typical safety restrictions preventing access to unaddressable and
// unexported data.  It works by digging the raw pointer to the underlying
// value out of the protected value and generating a new unprotected (unsafe)
// reflect.Value to it.
//
// This allows us to check for implementations of the Stringer and error
// interfaces to be used for pretty printing ordinarily unaddressable and
// inaccessible values such as unexported struct fields.
func unsafeReflectValue(v reflect.Value) (rv reflect.Value) {
	indirects := 1
	vt := v.Type()
	upv := unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetPtr)
	rvf := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetFlag))
	if rvf&flagIndir != 0 {
		vt = reflect.PtrTo(v.Type())
		indirects++
	} else if offsetScalar != 0 {
		upv = unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetScalar)
	}

	pv := reflect.NewAt(vt, upv)
	rv = pv
	for i := 0; i < indirects; i++ {
		rv = rv.Elem()
	}
	return rv
}

// Some constants in the form of bytes to avoid string overhead.  This mirrors
// the technique used in the fmt package.
var (
	panicBytes            = []byte("(PANIC=")
	plusBytes             = []byte("+")
	iBytes                = []byte("i")
	trueBytes             = []byte("true")
	falseBytes            = []byte("false")
	interfaceBytes        = []byte("(interface {})")
	commaNewlineBytes     = []byte(",\n")
	newlineBytes          = []byte("\n")
	openBraceBytes        = []byte("{")
	openBraceNewlineBytes = []byte("{\n")
	closeBraceBytes       = []byte("}")
	asteriskBytes         = []byte("*")
	colonBytes            = []byte(":")
	colonSpaceBytes       = []byte(": ")
	openParenBytes        = []byte("(")
	closeParenBytes       = []byte(")")
	spaceBytes            = []byte(" ")
	pointerChainBytes     = []byte("->")
	nilAngleBytes         = []byte("<nil>")
	maxNewlineBytes       = []byte("<max depth reached>\n")
	maxShortBytes         = []byte("<max>")
	circularBytes         = []byte("<already shown>")
	circularShortBytes    = []byte("<shown>")
	invalidAngleBytes     = []byte("<invalid>")
	openBracketBytes      = []byte("[")
	closeBracketBytes     = []byte("]")
	percentBytes          = []byte("%")
	precisionBytes        = []byte(".")
	openAngleBytes        = []byte("<")
	closeAngleBytes       = []byte(">")
	openMapBytes          = []byte("map[")
	closeMapBytes         = []byte("]")
	lenEqualsBytes        = []byte("len=")
	capEqualsBytes        = []byte("cap=")
)

// hexDigits is used to map a decimal value to a hex digit.
var hexDigits = "0123456789abcdef"

// catchPanic handles any panics that might occur during the handleMethods
// calls.
func catchPanic(w io.Writer, v reflect.Value) {
	if err := recover(); err != nil {
		w.Write(panicBytes)
		fmt.Fprintf(w, "%v", err)
		w.Write(closeParenBytes)
	}
}

// handleMethods attempts to call the Error and String methods on the underlying
// type the passed reflect.Value represents and outputes the result to Writer w.
//
// It handles panics in any called methods by catching and displaying the error
// as the formatted value.
func handleMethods(cs *ConfigState, w io.Writer, v reflect.Value) (handled bool) {
	// We need an interface to check if the type implements the error or
	// Stringer interface.  However, the reflect package won't give us an
	// interface on certain things like unexported struct fields in order
	// to enforce visibility rules.  We use unsafe to bypass these restrictions
	// since this package does not mutate the values.
	if !v.CanInterface() {
		v = unsafeReflectValue(v)
	}

	// Choose whether or not to do error and Stringer interface lookups against
	// the base type or a pointer to the base type depending on settings.
	// Technically calling one of these methods with a pointer receiver can
	// mutate the value, however, types which choose to satisify an error or
	// Stringer interface with a pointer receiver should not be mutating their
	// state inside these interface methods.
	var viface interface{}
	if !cs.DisablePointerMethods {
		if !v.CanAddr() {
			v = unsafeReflectValue(v)
		}
		viface = v.Addr().Interface()
	} else {
		if v.CanAddr() {
			v = v.Addr()
		}
		viface = v.Interface()
	}

	// Is it an error or Stringer?
	switch iface := viface.(type) {
	case error:
		defer catchPanic(w, v)
		if cs.ContinueOnMethod {
			w.Write(openParenBytes)
			w.Write([]byte(iface.Error()))
			w.Write(closeParenBytes)
			w.Write(spaceBytes)
			return false
		}

		w.Write([]byte(iface.Error()))
		return true

	case fmt.Stringer:
		defer catchPanic(w, v)
		if cs.ContinueOnMethod {
			w.Write(openParenBytes)
			w.Write([]byte(iface.String()))
			w.Write(closeParenBytes)
			w.Write(spaceBytes)
			return false
		}
		w.Write([]byte(iface.String()))
		return true
	}
	return false
}

// printBool outputs a boolean value as true or false to Writer w.
func printBool(w io.Writer, val bool) {
	if val {
		w.Write(trueBytes)
	} else {
		w.Write(falseBytes)
	}
}

// printInt outputs a signed integer value to Writer w.
func printInt(w io.Writer, val int64, base int) {
	w.Write([]byte(strconv.FormatInt(val, base)))
}

// printUint outputs an unsigned integer value to Writer w.
func printUint(w io.Writer, val uint64, base int) {
	w.Write([]byte(strconv.FormatUint(val, base)))
}

// printFloat outputs a floating point value using the specified precision,
// which is expected to be 32 or 64bit, to Writer w.
func printFloat(w io.Writer, val float64, precision int) {
	w.Write([]byte(strconv.FormatFloat(val, 'g', -1, precision)))
}

// printComplex outputs a complex value using the specified float precision
// for the real and imaginary parts to Writer w.
func printComplex(w io.Writer, c complex128, floatPrecision int) {
	r := real(c)
	w.Write(openParenBytes)
	w.Write([]byte(strconv.FormatFloat(r, 'g', -1, floatPrecision)))
	i := imag(c)
	if i >= 0 {
		w.Write(plusBytes)
	}
	w.Write([]byte(strconv.FormatFloat(i, 'g', -1, floatPrecision)))
	w.Write(iBytes)
	w.Write(closeParenBytes)
}

// printHexPtr outputs a uintptr formatted as hexidecimal with a leading '0x'
// prefix to Writer w.
func printHexPtr(w io.Writer, p uintptr) {
	// Null pointer.
	num := uint64(p)
	if num == 0 {
		w.Write(nilAngleBytes)
		return
	}

	// Max uint64 is 16 bytes in hex + 2 bytes for '0x' prefix
	buf := make([]byte, 18)

	// It's simpler to construct the hex string right to left.
	base := uint64(16)
	i := len(buf) - 1
	for num >= base {
		buf[i] = hexDigits[num%base]
		num /= base
		i--
	}
	buf[i] = hexDigits[num]

	// Add '0x' prefix.
	i--
	buf[i] = 'x'
	i--
	buf[i] = '0'

	// Strip unused leading bytes.
	buf = buf[i:]
	w.Write(buf)
}

// valuesSorter implements sort.Interface to allow a slice of reflect.Value
// elements to be sorted.
type valuesSorter struct {
	values []reflect.Value
}

// Len returns the number of values in the slice.  It is part of the
// sort.Interface implementation.
func (s *valuesSorter) Len() int {
	return len(s.values)
}

// Swap swaps the values at the passed indices.  It is part of the
// sort.Interface implementation.
func (s *valuesSorter) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

// Less returns whether the value at index i should sort before the
// value at index j.  It is part of the sort.Interface implementation.
func (s *valuesSorter) Less(i, j int) bool {
	switch s.values[i].Kind() {
	case reflect.Bool:
		return !s.values[i].Bool() && s.values[j].Bool()
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return s.values[i].Int() < s.values[j].Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return s.values[i].Uint() < s.values[j].Uint()
	case reflect.Float32, reflect.Float64:
		return s.values[i].Float() < s.values[j].Float()
	case reflect.String:
		return s.values[i].String() < s.values[j].String()
	case reflect.Uintptr:
		return s.values[i].Uint() < s.values[j].Uint()
	}
	return s.values[i].String() < s.values[j].String()
}

// sortValues is a generic sort function for native types: int, uint, bool,
// string and uintptr.  Other inputs are sorted according to their
// Value.String() value to ensure display stability.
func sortValues(values []reflect.Value) {
	if len(values) == 0 {
		return
	}
	sort.Sort(&valuesSorter{values})
}
