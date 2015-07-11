package processors

import (
	"fmt"
	"io"

	"github.com/linuxerwang/goats-html/util"
)

type GoCaseProcessor struct {
	BaseProcessor
	expression string
}

func (s *GoCaseProcessor) Process(writer io.Writer, context *TagContext) {
	context.MaybeAddImports(s.expression)
	io.WriteString(writer, fmt.Sprintf("case %s:\n", util.ToGoString(util.ToCamelExpr(s.expression))))

	if s.next != nil {
		s.next.Process(writer, context)
	}
}

func NewCaseProcessor(expression string) *GoCaseProcessor {
	processor := &GoCaseProcessor{
		expression: util.TrimWhiteSpaces(expression),
	}
	return processor
}
