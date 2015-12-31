package processors

import (
	"fmt"
	"io"
)

type GoSettingsProcessor struct {
	BaseProcessor
	Name string
}

func (s *GoSettingsProcessor) Process(writer io.Writer, ctx *TagContext) {
	// Start of a local scope
	io.WriteString(writer, "{\n")

	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("%s := %s\n", s.Name, "__impl.GetSettings()"))
	case "closure":
		io.WriteString(writer, fmt.Sprintf("var %s = %s;\n", s.Name, "__self.__getSettings()"))
	}

	if s.next != nil {
		s.next.Process(writer, ctx)
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
