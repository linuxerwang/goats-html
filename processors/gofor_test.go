package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/symbolmgr"
)

func TestNewForProcessorOneVarGo(t *testing.T) {
	processor := NewForProcessor(" \tbook \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	ctx.symMgr.Peek()["bookshelf"] = &symbolmgr.Symbol{
		Name:    "bookshelf",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != `{
__loopItems := bookshelf.Books
for _, book := range __loopItems {
DUMMY}
}
` {
		t.Errorf("Expected if block was not found. ", result.String())
	}
}

func TestNewForProcessorOneVarClosure(t *testing.T) {
	processor := NewForProcessor(" \tbook \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	ctx.symMgr.Peek()["bookshelf"] = &symbolmgr.Symbol{
		Name:    "bookshelf",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != `{
var __loopItems = bookshelf.getBooks();
if (__loopItems) {
goog.array.forEach(__loopItems, function(book, __index) {
DUMMY}, this);
}
}
` {
		t.Errorf("Expected if block was not found. ", result.String())
	}
}

func TestNewForProcessorTwoVarGo(t *testing.T) {
	processor := NewForProcessor(" \tidx, book \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	ctx.symMgr.Peek()["bookshelf"] = &symbolmgr.Symbol{
		Name:    "bookshelf",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != `{
__loopItems := bookshelf.Books
for idx, book := range __loopItems {
DUMMY}
}
` {
		t.Errorf("Expected for block was not found. ", result.String())
	}
}

func TestNewForProcessorTwoVarClosure(t *testing.T) {
	processor := NewForProcessor(" \tidx, book \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	ctx.symMgr.Peek()["bookshelf"] = &symbolmgr.Symbol{
		Name:    "bookshelf",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != `{
var __loopItems = bookshelf.getBooks();
if (__loopItems) {
goog.array.forEach(__loopItems, function(book, idx) {
DUMMY}, this);
}
}
` {
		t.Errorf("Expected for block was not found. ", result.String())
	}
}

func TestNewForProcessorThreeVarGo(t *testing.T) {
	processor := NewForProcessor(" \t@loopVar , idx , book \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	ctx.symMgr.Peek()["bookshelf"] = &symbolmgr.Symbol{
		Name:    "bookshelf",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != `{
__loopItems := bookshelf.Books
__loopTotal := len(__loopItems)
loopVar := &runtime.LoopVar{
Total: __loopTotal,
Counter0: -1,
RevCounter: __loopTotal + 1,
RevCounter0: __loopTotal,
First: __loopTotal > 1,
}
for idx, book := range __loopItems {
loopVar.Counter++
loopVar.Counter0++
loopVar.RevCounter--
loopVar.RevCounter0--
loopVar.Last = (loopVar.Total == loopVar.Counter)
DUMMY
loopVar.First = false
}
}
` {
		t.Errorf("Expected for block was not found. ", result.String())
	}
}

func TestNewForProcessorThreeVarClosure(t *testing.T) {
	processor := NewForProcessor(" \t@loopVar , idx , book \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	ctx.symMgr.Peek()["bookshelf"] = &symbolmgr.Symbol{
		Name:    "bookshelf",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != `{
var __loopItems = bookshelf.getBooks();
if (__loopItems) {
var __loopTotal = __loopItems.length;
var loopVar = {
total: __loopTotal,
counter0: -1,
revCounter: __loopTotal + 1,
revCounter0: __loopTotal,
first: __loopTotal > 1
};
goog.array.forEach(__loopItems, function(book, idx) {
loopVar.counter++;
loopVar.counter0++;
loopVar.revCounter--;
loopVar.revCounter0--;
loopVar.last = (loopVar.total == loopVar.counter);
DUMMY
loopVar.first = false;
}, this);
}
}
` {
		t.Errorf("Expected for block was not found. ", result.String())
	}
}

func TestForProcessorImportsOneVar(t *testing.T) {
	processor := NewForProcessor("\t userType: accounts.User.Types\t ")
	var result bytes.Buffer
	referer := pkgmgr.NewDummyAliasReferer()
	ctx := NewTagContext(pkgmgr.New("dummy"), referer, "go")
	ctx.symMgr.Peek()["accounts"] = &symbolmgr.Symbol{
		Name:    "accounts",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)
	switch len(referer.Aliases()) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.Aliases()["accounts"]; !ok {
			t.Errorf("Expected import \"accounts\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.Aliases())
	}
}

func TestForProcessorImportsTwoVar(t *testing.T) {
	processor := NewForProcessor("\t idx, userType: accounts.User.Types\t ")
	var result bytes.Buffer
	referer := pkgmgr.NewDummyAliasReferer()
	ctx := NewTagContext(pkgmgr.New("dummy"), referer, "go")
	ctx.symMgr.Peek()["accounts"] = &symbolmgr.Symbol{
		Name:    "accounts",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)
	switch len(referer.Aliases()) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.Aliases()["accounts"]; !ok {
			t.Errorf("Expected import \"accounts\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.Aliases())
	}
}

func TestForProcessorImportsThreeVar(t *testing.T) {
	processor := NewForProcessor("\t @loopVar, _, userType: accounts.User.Types\t ")
	var result bytes.Buffer
	referer := pkgmgr.NewDummyAliasReferer()
	ctx := NewTagContext(pkgmgr.New("dummy"), referer, "go")
	ctx.symMgr.Peek()["accounts"] = &symbolmgr.Symbol{
		Name:    "accounts",
		Type:    symbolmgr.TypeFor,
		IsPb:    true,
		PkgImpt: nil,
	}
	processor.Process(&result, ctx)
	switch len(referer.Aliases()) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.Aliases()["accounts"]; !ok {
			t.Errorf("Expected import \"accounts\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.Aliases())
	}
}
