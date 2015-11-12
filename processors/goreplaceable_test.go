package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
)

func TestNewReplaceableProcessorGo(t *testing.T) {
	processor := NewReplaceableProcessor("Tmpl", " \tSlot1 \t", []*Argument{})
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "if __impl.slot1 == nil {\nDUMMY} else {\n"+
		"  args := &TmplSlot1ReplArgs{\n  }\n  __impl.slot1(args)\n}\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestNewReplaceableProcessorClosure(t *testing.T) {
	processor := NewReplaceableProcessor("Tmpl", " \tSlot1 \t", []*Argument{})
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "if (this.slot1_ == null) {\nDUMMY} else {\n"+
		"  var __args = {};\n  this.slot1_(__element, __args);\n}\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}
