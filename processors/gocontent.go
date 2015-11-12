package processors

import (
	"fmt"
	"io"
)

type GoContentProcessor struct {
	BaseProcessor
	expression string
}

func (c *GoContentProcessor) Process(writer io.Writer, ctx *TagContext) {
	expr, err := ctx.RewriteExpression(c.expression)
	if err != nil {
		panic(err)
	}

	switch ctx.OutputFormat {
	case "go":
		if ctx.AutoEscape {
			io.WriteString(writer, fmt.Sprintf("__impl.WriteString(runtime.EscapeContent(%s))\n", expr))
		} else {
			io.WriteString(writer, fmt.Sprintf("__impl.WriteString(runtime.IndirectString(%s))\n", expr))
		}
	case "closure":
		io.WriteString(writer, fmt.Sprintf("var __text_node = goog.dom.createTextNode(%s);\n", expr))
		io.WriteString(writer, fmt.Sprintf("goog.dom.appendChild(__tag_stack[__tag_stack.length-1], __text_node);\n"))
	}

	// go:content is a terminal processor.
}

func NewContentProcessor(expression string) *GoContentProcessor {
	processor := &GoContentProcessor{
		expression: expression,
	}
	return processor
}
