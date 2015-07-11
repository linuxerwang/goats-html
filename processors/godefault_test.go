package processors

import (
	"bytes"
	"testing"
)

func TestDefaultProcessor(t *testing.T) {
	processor := NewDefaultProcessor()
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(NewDummyAliasReferer())
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not called.")
	}
	if result.String() != "default:\nDUMMY" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}
