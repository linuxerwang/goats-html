package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/symbolmgr"
)

func TestCaseProcessorGo(t *testing.T) {
	processor := NewCaseProcessor("\t proto1.book.Category_History")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)

	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	ctx.symMgr.Push(createSymbolMapForGoCase())

	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not called.")
	}
	if result.String() != "case proto1.book.Category_History:\nDUMMY" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestCaseProcessorClosure(t *testing.T) {
	processor := NewCaseProcessor("\t proto1.book.Category_History")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	ctx.symMgr.Push(createSymbolMapForGoCase())
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not called.")
	}
	if result.String() != "case proto1.book.Category_History:\nDUMMY\nbreak;\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestCaseProcessorImports(t *testing.T) {
	processor := NewCaseProcessor("\t proto1.book.Category_History\t ")
	var result bytes.Buffer
	referer := pkgmgr.NewDummyAliasReferer()
	ctx := NewTagContext(pkgmgr.New("dummy"), referer, "go")
	ctx.symMgr.Push(createSymbolMapForGoCase())
	processor.Process(&result, ctx)
	switch len(referer.Aliases()) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.Aliases()["proto1"]; !ok {
			t.Errorf("Expected import \"proto1\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.Aliases())
	}
}

func createSymbolMapForGoCase() map[string]*symbolmgr.Symbol {
	pi := &pkgmgr.PkgImport{}
	pi.SetName("proto1")
	pi.SetAlias("proto1")
	pi.SetPath("a/b/proto")
	pi.SetPbPkg("a.b.proto")
	sm := make(map[string]*symbolmgr.Symbol)
	sm["proto1"] = &symbolmgr.Symbol{
		Name:    "proto1",
		Type:    symbolmgr.TypeImport,
		IsPb:    true,
		PkgImpt: pi,
	}

	return sm
}
