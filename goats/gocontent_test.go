package goats

import (
	"bytes"
	"testing"
)

func TestNewContentProcessor(t *testing.T) {
	processor := NewContentProcessor(" \t book.Name\t ")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if dummy.Called {
		t.Errorf("Expect not calling the dummy processor but was called.")
	}
	if result.String() != "__impl.WriteString(runtime.EscapeContent(book.Name))\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestContentProcessorImports(t *testing.T) {
	processor := NewContentProcessor("\t 1 + proto1.Book.Category_HISTORY\t ")
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
