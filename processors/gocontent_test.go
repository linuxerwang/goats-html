package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/symbolmgr"
)

func TestNewContentProcessor(t *testing.T) {
	processor := NewContentProcessor(" \t book.Name\t ")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	ctx.symMgr.Push(createSymbolMapForGoContent())
	processor.Process(&result, ctx)

	if dummy.Called {
		t.Errorf("Expect not calling the dummy processor but was called.")
	}
	if result.String() != "__impl.WriteString(runtime.EscapeContent(book.Name))\n" {
		t.Errorf("Expected block was not found. ", result.String())
	}
}

func TestContentProcessorImports(t *testing.T) {
	processor := NewContentProcessor("\t 1 + proto1.Book.Category_History")
	var result bytes.Buffer
	referer := pkgmgr.NewDummyAliasReferer()
	ctx := NewTagContext(pkgmgr.New("dummy"), referer, "go")
	ctx.symMgr.Push(createSymbolMapForGoContent())
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

func createSymbolMapForGoContent() map[string]*symbolmgr.Symbol {
	sm := make(map[string]*symbolmgr.Symbol)
	sm["book"] = &symbolmgr.Symbol{
		Name:    "book",
		Type:    symbolmgr.TypeArg,
		IsPb:    true,
		PkgImpt: nil,
	}

	pi := &pkgmgr.PkgImport{}
	pi.SetName("proto1")
	pi.SetAlias("proto1")
	pi.SetPath("a/b/proto1")
	pi.SetPbPkg("a.b.proto1")
	sm["proto1"] = &symbolmgr.Symbol{
		Name:    "proto1",
		Type:    symbolmgr.TypeImport,
		IsPb:    true,
		PkgImpt: pi,
	}

	return sm
}
