package processors

import (
	"fmt"
	"io"
)

type GoDefaultProcessor struct {
	BaseProcessor
}

func (s *GoDefaultProcessor) Process(writer io.Writer, ctx *TagContext) {
	io.WriteString(writer, fmt.Sprintf("default:\n"))

	if s.next != nil {
		s.next.Process(writer, ctx)
	}
}

func NewDefaultProcessor() *GoDefaultProcessor {
	processor := &GoDefaultProcessor{}
	return processor
}
