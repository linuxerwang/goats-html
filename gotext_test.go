package goats

import (
	"bytes"
	"testing"
)

func TestNewTextProcessor(t *testing.T) {
	processor := NewTextProcessor(" \tWhatever \" ' > & <   text \n \t\t  content\n ")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if dummy.Called {
		t.Errorf("Expect not calling the dummy processor but was called.")
	}
	if result.String() != "__impl.WriteString(\" Whatever &#34; &#39; &gt; &amp; &lt; text content \")\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}
