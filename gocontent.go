package goats

import (
	"fmt"
	"io"
)

type GoContentProcessor struct {
	BaseProcessor
	expression string
}

func (c *GoContentProcessor) Process(writer io.Writer, context *TagContext) {
	expr := context.RewriteExpression(c.expression)

	if context.AutoEscape {
		io.WriteString(writer, fmt.Sprintf("__impl.WriteString(runtime.EscapeContent(%s))\n", expr))
	} else {
		io.WriteString(writer, fmt.Sprintf("__impl.WriteString(runtime.IndirectString(%s))\n", expr))
	}

	// go:content is a terminal processor.
}

func NewContentProcessor(expression string) *GoContentProcessor {
	processor := &GoContentProcessor{
		expression: expression,
	}
	return processor
}
