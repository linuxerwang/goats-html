package goats

import (
	"bytes"
	"testing"
)

func TestNewCommentProcessor(t *testing.T) {
	processor := NewCommentProcessor("Comment with \"quote\"")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "__impl.WriteString(`<!--Comment with \"quote\"-->`)\nDUMMY" {
		t.Errorf("Expected comment block was not found. ", result.String())
	}
}
