package processors

import (
	"fmt"
	"io"
)

type GoCommentProcessor struct {
	BaseProcessor
	comment string
}

func (c *GoCommentProcessor) Process(writer io.Writer, ctx *TagContext) {
	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("__impl.WriteString(`%s`)\n", "<!--"+c.comment+"-->"))
	}

	if c.next != nil {
		c.next.Process(writer, ctx)
	}
}

func NewCommentProcessor(comment string) *GoCommentProcessor {
	processor := &GoCommentProcessor{
		comment: comment,
	}
	return processor
}
