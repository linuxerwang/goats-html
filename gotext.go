package goats

import (
	"fmt"
	"html"
	"io"
)

type GoTextProcessor struct {
	BaseProcessor
	text string
}

func (i *GoTextProcessor) Process(writer io.Writer, context *TagContext) {
	// i.text is the (possiblly) merged text, its white space handling logic:
	//     - If before normalizing it has leading white space,
	//       it should have one leading space after normalizing.
	//     - If before normalizing it has trailing white space,
	//       it should have one trailing space after normalizing.
	//     - For all white spaces in the middle, reduces consecutive white spaces to one space.
	//
	hasLeadingSpace := HasLeadingSpace(i.text)
	hasTrailingSpace := HasTrailingSpace(i.text)
	text := NormalizeText(i.text)
	if text != "" {
		if hasLeadingSpace {
			text = " " + text
		}
		if hasTrailingSpace {
			text = text + " "
		}
		if context.AutoEscape {
			text = html.EscapeString(text)
		}
		io.WriteString(writer, fmt.Sprintf("__impl.WriteString(\"%s\")", text))
		io.WriteString(writer, "\n")
	}
	// go text is a terminal processor.
}

func NewTextProcessor(text string) *GoTextProcessor {
	return &GoTextProcessor{
		text: text,
	}
}
