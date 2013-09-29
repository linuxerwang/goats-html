package goats

import (
	"bytes"
	"testing"
)

func TestNewIfProcessor(t *testing.T) {
	processor := NewIfProcessor("\t a.science_books[0].name == a.science_books[1].name \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "if a.ScienceBooks[0].Name == a.ScienceBooks[1].Name {\nDUMMY}\n" {
		t.Errorf("Expected if block was not found. ", result.String())
	}
}

func TestIfProcessorImports(t *testing.T) {
	processor := NewIfProcessor("\t a.science_books[0].Type == shelf.Book.Type_HISTORY \t")
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)
	switch len(ctx.GetImports()) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 2:
		imports := ctx.GetImports()
		if _, ok := imports["a"]; !ok {
			t.Errorf("Expected import \"a\" but was not found:", imports)
		}
		if _, ok := imports["shelf"]; !ok {
			t.Errorf("Expected import \"shelf\" but was not found:", imports)
		}
	default:
		t.Errorf("Expect one import but found different: ", ctx.GetImports())
	}
}
