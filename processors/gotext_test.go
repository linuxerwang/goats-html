package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
)

func TestNewTextProcessorGo(t *testing.T) {
	processor := NewTextProcessor(" \tWhatever \" ' > & <   text \n \t\t  content\n ")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	processor.Process(&result, ctx)

	if dummy.Called {
		t.Errorf("Expect not calling the dummy processor but was called.")
	}
	if result.String() != "__impl.WriteString(\" Whatever &#34; &#39; &gt; &amp; &lt; text content \")\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestNewTextProcessorClosure(t *testing.T) {
	processor := NewTextProcessor(" \tWhatever \" ' > & <   text \n \t\t  content\n ")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	processor.Process(&result, ctx)

	if dummy.Called {
		t.Errorf("Expect not calling the dummy processor but was called.")
	}
	if result.String() != "var __text_node = goog.dom.createTextNode(\" Whatever &#34; &#39; &gt; &amp; &lt; text content \");\n"+
		"goog.dom.appendChild(__tag_stack[__tag_stack.length-1], __text_node);\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}
