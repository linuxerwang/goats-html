package goats

import (
	"bytes"
	"testing"

	"golang.org/x/net/html"
)

func TestNewDocTypeProcessor(t *testing.T) {
	processor := NewDocTypeProcessor("html", []html.Attribute{
		html.Attribute{
			Key: "public",
			Val: "-//W3C//DTD XHTML 1.0 Transitional//EN",
		},
		html.Attribute{
			Key: "system",
			Val: "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd",
		}})
	dummy := NewDummyProcessor()
	processor.SetNext(dummy)
	var result bytes.Buffer
	ctx := NewTagContext()
	processor.Process(&result, ctx)

	if !dummy.Called {
		t.Errorf("Expect calling the dummy processor but was not.")
	}
	if result.String() != "if !__impl.GetSettings().OmitDocType {\n"+
		"__impl.WriteString(\"<!DOCTYPE html PUBLIC \\\"-//W3C//DTD XHTML 1.0 Transitional//EN\\\" "+
		"\\\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\\\">\\n\")\n}\nDUMMY" {
		t.Errorf("Expected doctype was not found. ", result.String())
	}
}
