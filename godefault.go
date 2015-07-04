package goats

import (
	"fmt"
	"io"
)

type GoDefaultProcessor struct {
	BaseProcessor
}

func (s *GoDefaultProcessor) Process(writer io.Writer, context *TagContext) {
	io.WriteString(writer, fmt.Sprintf("default:\n"))

	if s.next != nil {
		s.next.Process(writer, context)
	}
}

func NewDefaultProcessor() *GoDefaultProcessor {
	processor := &GoDefaultProcessor{}
	return processor
}
