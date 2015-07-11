package processors

import (
	"fmt"
	"io"
)

type GoSettingsProcessor struct {
	BaseProcessor
	Name string
}

func (s *GoSettingsProcessor) Process(writer io.Writer, context *TagContext) {
	// Start of a local scope
	io.WriteString(writer, "{\n")

	io.WriteString(writer, fmt.Sprintf("%s := %s", s.Name, "__impl.GetSettings()"))
	io.WriteString(writer, "\n")

	if s.next != nil {
		s.next.Process(writer, context)
	}

	// End of a local scope
	io.WriteString(writer, "}\n")
}

func NewSettingsProcessor(name string) *GoSettingsProcessor {
	processor := &GoSettingsProcessor{
		Name: name,
	}
	return processor
}
