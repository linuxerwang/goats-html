package goats

import (
	"bytes"
	"testing"
)

func TestSwitchTextProcessor(t *testing.T) {
	processor := NewSwitchProcessor(" \t book.Id \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not called.")
	}
	if result.String() != "switch book.Id {\nDUMMY}\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}
