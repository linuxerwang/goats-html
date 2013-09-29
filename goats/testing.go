package goats

import (
	"io"
)

type dummyProcessor struct {
	BaseProcessor
	Called bool
}

func (dp *dummyProcessor) Process(writer io.Writer, context *TagContext) {
	io.WriteString(writer, "DUMMY")
	dp.Called = true
}

func NewDummyProcessor() *dummyProcessor {
	return &dummyProcessor{}
}
