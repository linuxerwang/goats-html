package processors

import (
	"fmt"
	"io"

	"github.com/linuxerwang/goats-html/util"
)

type GoIfProcessor struct {
	BaseProcessor
	conditional string
}

func (i *GoIfProcessor) Process(writer io.Writer, ctx *TagContext) {
	expr, err := ctx.RewriteExpression(i.conditional)
	if err != nil {
		panic(err)
	}

	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("if %s {\n", expr))
	case "closure":
		io.WriteString(writer, fmt.Sprintf("if (%s) {\n", expr))
	}

	if i.next != nil {
		i.next.Process(writer, ctx)
	}

	io.WriteString(writer, "}\n")
}

func NewIfProcessor(conditional string) *GoIfProcessor {
	processor := &GoIfProcessor{
		conditional: util.TrimWhiteSpaces(conditional),
	}
	return processor
}
