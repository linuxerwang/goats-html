package runtime

import (
	"fmt"
	"io"
)

type TagAttrs map[string]string

func (ta TagAttrs) AddAttr(name string, a interface{}) {
	val := EscapeAttr(fmt.Sprint(a))
	if existingVal, found := ta[name]; found && IsMergeable(name) {
		switch name {
		case "class":
			ta[name] = existingVal + " " + val
		case "style":
			ta[name] = existingVal + "; " + val
		default:
			ta[name] = existingVal + ", " + val
		}
	} else {
		ta[name] = val
	}
}

func (ta TagAttrs) MergeFrom(fromAttrs map[string]string) {
	for key, val := range fromAttrs {
		ta.AddAttr(key, val)
	}
}

func (ta TagAttrs) GenTagAndAttrs(writer io.Writer, tagName string) {
	io.WriteString(writer, "<")
	io.WriteString(writer, tagName)
	for name, val := range ta {
		io.WriteString(writer, " ")
		io.WriteString(writer, name)
		io.WriteString(writer, "=\"")
		io.WriteString(writer, EscapeContent(val))
		io.WriteString(writer, "\"")
	}
	io.WriteString(writer, ">")
}
