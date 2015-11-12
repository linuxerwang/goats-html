package processors

import (
	"fmt"
	"html"
	"io"

	"github.com/linuxerwang/goats-html/util"
)

type GoTextProcessor struct {
	BaseProcessor
	text string
}

func (i *GoTextProcessor) Process(writer io.Writer, ctx *TagContext) {
	// i.text is the (possiblly) merged text, its white space handling logic:
	//     - If before normalizing it has leading white space,
	//       it should have one leading space after normalizing.
	//     - If before normalizing it has trailing white space,
	//       it should have one trailing space after normalizing.
	//     - For all white spaces in the middle, reduces consecutive white spaces to one space.
	//
	hasLeadingSpace := util.HasLeadingSpace(i.text)
	hasTrailingSpace := util.HasTrailingSpace(i.text)
	text := util.NormalizeText(i.text)
	if text != "" {
		if hasLeadingSpace {
			text = " " + text
		}
		if hasTrailingSpace {
			text = text + " "
		}
		if ctx.AutoEscape {
			text = html.EscapeString(text)
		}
		switch ctx.OutputFormat {
		case "go":
			io.WriteString(writer, fmt.Sprintf("__impl.WriteString(\"%s\")\n", text))
		case "closure":
			io.WriteString(writer, fmt.Sprintf("var __text_node = goog.dom.createTextNode(\"%s\");\n", text))
			io.WriteString(writer, fmt.Sprintf("goog.dom.appendChild(__tag_stack[__tag_stack.length-1], __text_node);\n"))
		}
	}
	// go text is a terminal processor.
}

func NewTextProcessor(text string) *GoTextProcessor {
	return &GoTextProcessor{
		text: text,
	}
}
