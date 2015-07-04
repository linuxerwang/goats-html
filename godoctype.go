package goats

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

const (
	DOCTYPE = "<!DOCTYPE %s"
	PUBLIC  = "public"
	SYSTEM  = "system"
)

type GoDocTypeProcessor struct {
	BaseProcessor
	name  string
	attrs []html.Attribute
}

func (d *GoDocTypeProcessor) Process(writer io.Writer, context *TagContext) {
	var doctype bytes.Buffer
	doctype.WriteString(fmt.Sprintf(DOCTYPE, d.name))
	for _, attr := range d.attrs {
		if attr.Key == PUBLIC {
			doctype.WriteString(" PUBLIC")
		} else if attr.Key == SYSTEM {
			// Do nothing
		}
		doctype.WriteString(" \\\"" + attr.Val + "\\\"")
	}
	doctype.WriteString(">\\n")

	io.WriteString(writer, "if !__impl.GetSettings().OmitDocType {\n")
	io.WriteString(writer, fmt.Sprintf("__impl.WriteString(\"%s\")\n", doctype.String()))
	io.WriteString(writer, "}\n")

	if d.next != nil {
		d.next.Process(writer, context)
	}
}

func NewDocTypeProcessor(name string, attrs []html.Attribute) *GoDocTypeProcessor {
	processor := &GoDocTypeProcessor{
		name:  name,
		attrs: attrs,
	}
	return processor
}
