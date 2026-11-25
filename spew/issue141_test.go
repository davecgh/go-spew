package spew_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type namedMapWithStringer map[int]int

func (m namedMapWithStringer) String() string {
	return fmt.Sprintf("named map stringer with %d entries", len(m))
}

func TestNamedMapStringer(t *testing.T) {
	out := spew.Sdump(namedMapWithStringer{1: 2})
	if strings.Contains(out, "PANIC") {
		t.Fatalf("unexpected panic while dumping a named map with a String method: %s", out)
	}
	if !strings.Contains(out, "named map stringer with 1 entries") {
		t.Fatalf("String method was not called: %s", out)
	}
}
