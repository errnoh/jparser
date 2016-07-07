package jparser

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse(strings.NewReader(jsonStream))
	}
}

func TestParse(t *testing.T) {
	expected := "(object)[{ Bacon:string } { Bacon2:number } { Bacon3:(string)[string string string] } { Bacon4:null }] <nil>"
	got := fmt.Sprint(Parse(strings.NewReader(jsonStream4)))
	if got != expected {
		t.Fatalf("Expected \n'%s'\n, got \n'%s'", expected, got)
	}
}

func TestStringValue(t *testing.T) {
	expected := "(object)[{ Bacon:foo } { Bacon2:4 } { Bacon3:(string)[1 2 3 ] } { Bacon4:<nil> } ]"
	p, _ := Parse(strings.NewReader(jsonStream4))
	got := p.StringValue()
	if got != expected {
		t.Fatalf("Expected \n'%s'\n, got \n'%s'", expected, got)
	}
}
