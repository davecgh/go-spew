package spew_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type NestedErr struct {
	Err error
}

func (e *NestedErr) Error() string {
	return e.Err.Error()
}

type ErrWrap struct {
	NestedErr *NestedErr
}

func (e ErrWrap) Error() string {
	return e.NestedErr.Error()
}

func TestDumpCustomErrorWrapDoesNotPanic(t *testing.T) {
	err := ErrWrap{
		NestedErr: &NestedErr{
			Err: fmt.Errorf("this is an error"),
		},
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("spew.Sdump panicked: %v", r)
		}
	}()
	_ = spew.Sdump(err)
}

func TestDumpNestedCustomErrorWrapDoesNotPanic(t *testing.T) {
	type errWrap struct {
		NestedErr *NestedErr
	}

	err := errWrap{
		NestedErr: &NestedErr{
			Err: fmt.Errorf("this is an error"),
		},
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("spew.Sdump panicked: %v", r)
		}
	}()
	_ = spew.Sdump(err)
}

type mapStringer map[int]int

func (x mapStringer) String() string {
	return fmt.Sprintf("{M1:%d}", len(x))
}

func TestDumpMapValueStringerDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("spew.Sdump panicked: %v", r)
		}
	}()
	out := spew.Sdump(mapStringer{1: 2})
	if strings.Contains(out, "PANIC=") {
		t.Fatalf("spew.Sdump output contains panic: %s", out)
	}
}
