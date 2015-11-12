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

func (gcp *GoCaseProcessor) Process(writer io.Writer, ctx *TagContext) {
	ctx.MaybeAddImports(gcp.expression)

	s, err := ctx.RewriteExpression(gcp.expression)
	if err != nil {
		panic("Error, " + err.Error())
	}

	io.WriteString(writer, fmt.Sprintf("case %s:\n", s))

	if gcp.next != nil {
		gcp.next.Process(writer, ctx)
	}

	switch ctx.OutputFormat {
	case "closure":
		io.WriteString(writer, "\nbreak;\n")
	}
}

func NewCaseProcessor(expression string) *GoCaseProcessor {
	processor := &GoCaseProcessor{
		expression: util.TrimWhiteSpaces(expression),
	}
	return processor
}
