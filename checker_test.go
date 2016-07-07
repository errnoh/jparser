package jparser

import (
	"strings"
	"testing"
)

const (
	jsonStream  = `{"Message": "Hello", "Array": [1, "foo", 3], "Null": null, "bacon": false, "Number": 1.234}`
	jsonStream2 = `{"MessageFoo": "Hello", "Array": [1, 3, 3], "Arr2": [{"Bacon": "foo"}], "Null": 3, "Number": 1.234}`
	jsonStream3 = `{"Bacon": "watch", "MessageFoo": "Hello", "Array": [1, 3, 3], "Null": 3, "Number": 1.234, "Arr2": [{"Bacon": "foo"}]}`
	jsonStream4 = `[{"Bacon": "foo"}, {"Bacon2": 4}, {"Bacon3": ["1","2","3"]}, {"Bacon4": null}]`
)

func BenchmarkGet(b *testing.B) {
	r3, _ := Parse(strings.NewReader(jsonStream3))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r3.Get("Arr2.[0].Bacon")
	}
}

func TestGet(t *testing.T) {
	r3, _ := Parse(strings.NewReader(jsonStream3))
	v, ok := r3.Get("Bacon")
	if !ok || v != "watch" {
		t.Fatalf("Expected 'watch', got '%s'", v)
	}
	v, ok = r3.Get("Array.[0]")
	if !ok || v != "1" {
		t.Fatalf("Expected '1', got '%s'", v)
	}
	v, ok = r3.Get("Arr2.[0].Bacon")
	if !ok || v != "foo" {
		t.Fatalf("Expected 'foo', got '%s'", v)
	}
	v, ok = r3.Get("Arr2.[0.Bacon")
	if ok || v != "" {
		t.Fatalf("Got '%s %s' instead of fail", ok, v)
	}
	v, ok = r3.Get("Arr2.[0].Crispy")
	if ok || v != "" {
		t.Fatalf("Got '%s %s' instead of fail", ok, v)
	}
}

func TestGet2(t *testing.T) {
	r, _ := Parse(strings.NewReader(`7`))
	v, ok := r.Get(".")
	if !ok || v != "7" {
		t.Fatalf("Expected '7', got '%s'", v)
	}
}
