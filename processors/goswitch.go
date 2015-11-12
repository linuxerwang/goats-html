package processors

import (
	"fmt"
	"io"

	"github.com/linuxerwang/goats-html/util"
)

type GoSwitchProcessor struct {
	BaseProcessor
	expression string
}

func (s *GoSwitchProcessor) Process(writer io.Writer, ctx *TagContext) {
	expr, err := ctx.RewriteExpression(s.expression)
	if err != nil {
		panic(err)
	}

	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("switch %s {\n", expr))
	case "closure":
		io.WriteString(writer, fmt.Sprintf("switch (%s) {\n", expr))
	}

	if s.next != nil {
		s.next.Process(writer, ctx)
	}

	io.WriteString(writer, "}\n")
}

func NewSwitchProcessor(expression string) *GoSwitchProcessor {
	processor := &GoSwitchProcessor{
		expression: util.TrimWhiteSpaces(expression),
	}
	return processor
}
