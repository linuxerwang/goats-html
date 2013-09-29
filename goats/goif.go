package goats

import (
	"fmt"
	"io"
)

type GoIfProcessor struct {
	BaseProcessor
	conditional string
}

func (i *GoIfProcessor) Process(writer io.Writer, context *TagContext) {
	expr := context.RewriteExpression(i.conditional)
	io.WriteString(writer, fmt.Sprintf("if %s {\n", expr))

	if i.next != nil {
		i.next.Process(writer, context)
	}

	io.WriteString(writer, "}\n")
}

func NewIfProcessor(conditional string) *GoIfProcessor {
	processor := &GoIfProcessor{
		conditional: TrimWhiteSpaces(conditional),
	}
	return processor
}
