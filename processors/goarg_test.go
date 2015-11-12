package processors

import (
	"bytes"
	"testing"

	"github.com/linuxerwang/goats-html/pkgmgr"
)

func TestNewArgDefWithValue(t *testing.T) {
	argDef := " book \t: \tproto.Book \t = \t books[0] "
	arg := NewArgDef(argDef)
	assertBookArg(t, arg, true)
}

func TestNewArgDefWithoutValue(t *testing.T) {
	argDef := " book \t: \tproto.Book \t "
	arg := NewArgDef(argDef)
	assertBookArg(t, arg, false)
}

func TestNewArgCall(t *testing.T) {
	argCall := " book \t: \t data.comic_books[0] "
	arg := NewArgCall(argCall)
	assertBookCall(t, arg)
}

func TestParseArgDefs(t *testing.T) {
	argDefs := " book \t: \tproto.Book \t = \t books[0] ; author \t : proto.Author  = author;"
	args := ParseArgDefs(argDefs)
	if len(args) != 2 {
		t.Errorf("Expect two arguments parsed but was %d", len(args))
	}
	assertBookArg(t, args[0], true)
	assertAuthorArg(t, args[1])
}

func TestParseArgCalls(t *testing.T) {
	argDefs := " book \t: \t data.comic_books[0]  ; author \t : data.book_authors[0];"
	args := ParseArgCalls(argDefs)
	if len(args) != 2 {
		t.Errorf("Expect two arguments parsed but was %d", len(args))
	}
	assertBookCall(t, args[0])
	assertAuthorCall(t, args[1])
}

func TestProcessorGo(t *testing.T) {
	argDef := " book \t: \tproto.Book \t = \t books[0] "
	arg := NewArgDef(argDef)
	processor := NewArgProcessor([]*Argument{arg})
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "go")
	processor.Process(&result, ctx)
	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "book := __args.Book\nDUMMY" {
		t.Errorf("Expected arg block was not found. ", result.String())
	}
}

func TestProcessorClosure(t *testing.T) {
	argDef := " book \t: \tproto.Book \t = \t books[0] "
	arg := NewArgDef(argDef)
	processor := NewArgProcessor([]*Argument{arg})
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext(pkgmgr.New("dummy"), pkgmgr.NewDummyAliasReferer(), "closure")
	processor.Process(&result, ctx)
	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "var book = __args[\"book\"];\nDUMMY" {
		t.Errorf("Expected arg block was not found. ", result.String())
	}
}

//----------- Private functions -----------

func assertBookArg(t *testing.T, arg *Argument, withValue bool) {
	if arg.Name != "book" {
		t.Error("Expect arg name \"book\" but was not.")
	}
	if arg.PkgName != "proto" {
		t.Error("Expect pkg name \"proto\" but was not.")
	}
	if arg.Type != "proto.Book" {
		t.Error("Expect arg type \"proto.Book\" but was not.")
	}
	if withValue {
		if arg.Val != "books[0]" {
			t.Error("Expect arg value \"books[0]\" but was not.")
		}
	} else if arg.Val != "" {
		t.Error("Expect empty arg value but was not.")
	}
	if arg.Declare != "Book proto.Book" {
		t.Error("Expect arg declare \"Book proto.Book\" but was not.", arg.Declare)
	}
}

func assertAuthorArg(t *testing.T, arg *Argument) {
	if arg.Name != "author" {
		t.Error("Expect arg name \"author\" but was not.")
	}
	if arg.PkgName != "proto" {
		t.Error("Expect pkg name \"proto\" but was not.")
	}
	if arg.Type != "proto.Author" {
		t.Error("Expect arg type \"proto.Author\" but was not.")
	}
	if arg.Val != "author" {
		t.Error("Expect arg value \"author\" but was not.")
	}
	if arg.Declare != "Author proto.Author" {
		t.Error("Expect arg declare \"Author proto.Author\" but got.", arg.Declare)
	}
}

func assertBookCall(t *testing.T, arg *Argument) {
	if arg.Name != "book" {
		t.Error("Expect arg name \"book\" but was not.")
	}
	if arg.Val != "data.comic_books[0]" {
		t.Error("Expect arg value \"data.comic_books[0]\" but got ", arg.Val)
	}
}

func assertAuthorCall(t *testing.T, arg *Argument) {
	if arg.Name != "author" {
		t.Error("Expect arg name \"author\" but was not.")
	}
	if arg.Val != "data.book_authors[0]" {
		t.Error("Expect arg value \"data.book_authors[0]\" but was not.")
	}
}
