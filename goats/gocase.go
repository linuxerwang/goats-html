package goats

import (
	"fmt"
	"io"
)

type GoCaseProcessor struct {
	BaseProcessor
	expression string
}

func (s *GoCaseProcessor) Process(writer io.Writer, context *TagContext) {
	context.MaybeAddImports(s.expression)
	io.WriteString(writer, fmt.Sprintf("case %s:\n", ToGoString(ToCamelExpr(s.expression))))

	if s.next != nil {
		s.next.Process(writer, context)
	}
}

func NewCaseProcessor(expression string) *GoCaseProcessor {
	processor := &GoCaseProcessor{
		expression: TrimWhiteSpaces(expression),
	}
	return processor
}
