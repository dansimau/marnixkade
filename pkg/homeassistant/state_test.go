package homeassistant_test

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type testStruct struct {
	Foo string
	Bar string
	Baz string
}

type fooStruct struct {
	Foo string
	Bar string
	Wut int

	testStruct
}

func TestUnmarshalJSON(t *testing.T) {
	var s fooStruct
	if err := json.Unmarshal([]byte(`{"foo": "foo", "bar": "bar", "baz": "baz", "wut": "123"}`), &s); err != nil {
		t.Fatal(err)
	}

	spew.Dump(s)
}
