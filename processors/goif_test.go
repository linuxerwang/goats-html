package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/symbolmgr"
)

func TestNewIfProcessorGo(t *testing.T) {
	processor := NewIfProcessor("\t a.science_books[0].name == a.science_books[1].name \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	ctx.symMgr.Peek()["a"] = &symbolmgr.Symbol{
		Name:    "a",
		Type:    symbolmgr.TypeVar,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "if a.ScienceBooks[0].Name==a.ScienceBooks[1].Name {\nDUMMY}\n" {
		t.Errorf("Expected if block was not found. ", result.String())
	}
}

func TestNewIfProcessorClosure(t *testing.T) {
	processor := NewIfProcessor("\t a.science_books[0].name == a.science_books[1].name \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	ctx.symMgr.Peek()["a"] = &symbolmgr.Symbol{
		Name:    "a",
		Type:    symbolmgr.TypeVar,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "if (a.getScienceBooks()[0].getName()==a.getScienceBooks()[1].getName()) {\nDUMMY}\n" {
		t.Errorf("Expected if block was not found. ", result.String())
	}
}

func TestIfProcessorImports(t *testing.T) {
	processor := NewIfProcessor("\t a.science_books[0].Type == shelf.book.Type_HISTORY \t")
	var result bytes.Buffer
	referer := pkgmgr.NewDummyAliasReferer()
	ctx := NewTagContext(pkgmgr.New("dummy"), referer, "go")
	ctx.symMgr.Peek()["a"] = &symbolmgr.Symbol{
		Name:    "a",
		Type:    symbolmgr.TypeVar,
		IsPb:    true,
		PkgImpt: nil,
	}
	sym := &symbolmgr.Symbol{
		Name:    "shelf",
		Type:    symbolmgr.TypeImport,
		IsPb:    true,
		PkgImpt: &pkgmgr.PkgImport{},
	}
	sym.PkgImpt.SetName("shelf")
	sym.PkgImpt.SetAlias("shelf")
	sym.PkgImpt.SetPath("x/y/shelf")
	sym.PkgImpt.SetPbPkg("x.y.shelf")
	ctx.symMgr.Peek()["shelf"] = sym

	processor.Process(&result, ctx)
	switch len(referer.Aliases()) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 2:
		if _, ok := referer.Aliases()["a"]; !ok {
			t.Errorf("Expected import \"a\" but was not found.")
		}
		if _, ok := referer.Aliases()["shelf"]; !ok {
			t.Errorf("Expected import \"shelf\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.Aliases())
	}
}
