package goats

import (
	"bytes"
	"testing"
)

func TestCaseProcessor(t *testing.T) {
	processor := NewCaseProcessor("\t proto1.book.Category_HISTORY\t ")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not called.")
	}
	if result.String() != "case proto1.Book.Category_HISTORY:\nDUMMY" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestCaseProcessorImports(t *testing.T) {
	processor := NewCaseProcessor("\t proto1.book.Category_HISTORY\t ")
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)
	switch len(ctx.GetImports()) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		imports := ctx.GetImports()
		if _, ok := imports["proto1"]; !ok {
			t.Errorf("Expected import \"proto1\" but was not found:", imports)
		}
	default:
		t.Errorf("Expect one import but found different: ", ctx.GetImports())
	}
}
