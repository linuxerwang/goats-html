package processors

import (
	"bytes"
	"testing"
)

func TestNewContentProcessor(t *testing.T) {
	processor := NewContentProcessor(" \t book.Name\t ")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(NewDummyAliasReferer())
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
	referer := NewDummyAliasReferer()
	ctx := NewTagContext(referer)
	processor.Process(&result, ctx)
	switch len(referer.aliases) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.aliases["proto1"]; !ok {
			t.Errorf("Expected import \"proto1\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.aliases)
	}
}
