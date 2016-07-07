package jparser

import (
	"strings"
	"testing"
)

func BenchmarkCmp(b *testing.B) {
	r1, _ := Parse(strings.NewReader(jsonStream))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compare(r1, r1)
	}
}

func TestCmp(t *testing.T) {
	r1, _ := Parse(strings.NewReader(jsonStream))
	r2, _ := Parse(strings.NewReader(jsonStream2))
	r3, _ := Parse(strings.NewReader(jsonStream3))
	c1 := Compare(r1, r2)
	if c1&Incompatible != Incompatible || c1.String() != "incompatible" {
		t.Fail()
	}
	c2 := Compare(r2, r2)
	if c2 != Equal || c2.String() != "equal" {
		t.Fail()
	}
	c3 := Compare(r2, r3)
	if c3 != Subset || c3.String() != "subset" {
		t.Fail()
	}
}
