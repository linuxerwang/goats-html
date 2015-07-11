package processors

import (
	"bytes"
	"testing"
)

func TestNewForProcessorOneVar(t *testing.T) {
	processor := NewForProcessor(" \tbook \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(NewDummyAliasReferer())
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

func TestNewForProcessorTwoVar(t *testing.T) {
	processor := NewForProcessor(" \tidx, book \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(NewDummyAliasReferer())
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

func TestNewForProcessorThreeVar(t *testing.T) {
	processor := NewForProcessor(" \t@loopVar , idx , book \t: bookshelf.books \t")
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(NewDummyAliasReferer())
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

func TestForProcessorImportsOneVar(t *testing.T) {
	processor := NewForProcessor("\t userType: accounts.User.Types\t ")
	var result bytes.Buffer
	referer := NewDummyAliasReferer()
	ctx := NewTagContext(referer)
	processor.Process(&result, ctx)
	switch len(referer.aliases) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.aliases["accounts"]; !ok {
			t.Errorf("Expected import \"accounts\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.aliases)
	}
}

func TestForProcessorImportsTwoVar(t *testing.T) {
	processor := NewForProcessor("\t idx, userType: accounts.User.Types\t ")
	var result bytes.Buffer
	referer := NewDummyAliasReferer()
	ctx := NewTagContext(referer)
	processor.Process(&result, ctx)
	switch len(referer.aliases) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.aliases["accounts"]; !ok {
			t.Errorf("Expected import \"accounts\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.aliases)
	}
}

func TestForProcessorImportsThreeVar(t *testing.T) {
	processor := NewForProcessor("\t @loopVar, _, userType: accounts.User.Types\t ")
	var result bytes.Buffer
	referer := NewDummyAliasReferer()
	ctx := NewTagContext(referer)
	processor.Process(&result, ctx)
	switch len(referer.aliases) {
	case 0:
		t.Errorf("Expected import not found. ")
	case 1:
		if _, ok := referer.aliases["accounts"]; !ok {
			t.Errorf("Expected import \"accounts\" but was not found.")
		}
	default:
		t.Errorf("Expect one import but found more: ", referer.aliases)
	}
}
