package goats

import (
	"bytes"
	"testing"
)

func TestNewReplaceableProcessor(t *testing.T) {
	processor := NewReplaceableProcessor("Tmpl", " \tSlot1 \t", []*Argument{})
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "if __impl.slot1 == nil {\nDUMMY} else {\n"+
		"  args := &TmplSlot1ReplArgs{\n  }\n  __impl.slot1(args)\n}\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}
