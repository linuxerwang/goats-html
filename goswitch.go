package goats

import (
	"fmt"
	"io"
)

type GoSwitchProcessor struct {
	BaseProcessor
	expression string
}

func (s *GoSwitchProcessor) Process(writer io.Writer, context *TagContext) {
	expr := context.RewriteExpression(s.expression)
	io.WriteString(writer, fmt.Sprintf("switch %s {\n", expr))

	if s.next != nil {
		s.next.Process(writer, context)
	}

	io.WriteString(writer, "}\n")
}

func NewSwitchProcessor(expression string) *GoSwitchProcessor {
	processor := &GoSwitchProcessor{
		expression: TrimWhiteSpaces(expression),
	}
	return processor
}
