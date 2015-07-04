package goats

import (
	"bytes"
	"testing"
)

func TestNewVarDef(t *testing.T) {
	varDef := newVarDef("  book \t: \tbookshelf.books[0] ")
	if varDef.Name != "book" {
		t.Error("Expect variable name \"book\" but was not.")
	}
	if varDef.Val != "bookshelf.books[0]" {
		t.Error("Expect variable expression \"bookshelf.books[0]\" but was not.")
	}
}

func TestNewVarProcessor(t *testing.T) {
	varProcessor := NewVarProcessor("  book \t: \tbookshelf.books[0] ")
	dummy := NewDummyProcessor()
	varProcessor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	varProcessor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "{\nbook := bookshelf.Books[0]\n\nDUMMY}\n" {
		t.Errorf("Expected var block was not found. ", result.String())
	}
}

func TestNewVarsProcessor(t *testing.T) {
	varProcessor := NewVarsProcessor(
		"  book \t: \tbookshelf.books[0] ; \t author\t: book.authors[0]\t")
	dummy := NewDummyProcessor()
	varProcessor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	varProcessor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "{\nbook := bookshelf.Books[0]\n\nauthor := book.Authors[0]\n\nDUMMY}\n" {
		t.Errorf("Expected vars block was not found. ", result.String())
	}
}
