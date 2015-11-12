package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/symbolmgr"
)

func TestSwitchTextProcessorGo(t *testing.T) {
	processor := NewSwitchProcessor(" \t book.id \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	ctx.symMgr.Peek()["book"] = &symbolmgr.Symbol{
		Name:    "book",
		Type:    symbolmgr.TypeVar,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not called.")
	}
	if result.String() != "switch book.Id {\nDUMMY}\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestSwitchTextProcessorClosure(t *testing.T) {
	processor := NewSwitchProcessor(" \t book.id \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	ctx.symMgr.Peek()["book"] = &symbolmgr.Symbol{
		Name:    "book",
		Type:    symbolmgr.TypeVar,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not called.")
	}
	if result.String() != "switch (book.getId()) {\nDUMMY}\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}
