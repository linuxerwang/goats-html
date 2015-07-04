package goats

import (
	"fmt"
	"io"
)

type GoReplaceableProcessor struct {
	BaseProcessor
	tmplName   string
	slotName   string
	hiddenName string
	args       []*Argument
}

func (r *GoReplaceableProcessor) Process(writer io.Writer, context *TagContext) {
	io.WriteString(writer, fmt.Sprintf("if __impl.%s == nil {\n", r.hiddenName))
	for _, arg := range r.args {
		io.WriteString(writer, fmt.Sprintf("  %s := %s\n", arg.Name, arg.Val))
	}

	if r.next != nil {
		r.next.Process(writer, context)
	}

	io.WriteString(writer, "} else {\n")
	io.WriteString(writer, fmt.Sprintf("  args := &%s%sReplArgs{\n", r.tmplName, r.slotName))
	for _, arg := range r.args {
		io.WriteString(writer, fmt.Sprintf("    %s: %s,\n", ToPublicName(arg.Name), arg.Val))
	}
	io.WriteString(writer, fmt.Sprintf("  }\n"))
	io.WriteString(writer, fmt.Sprintf("  __impl.%s(args)\n", r.hiddenName))
	io.WriteString(writer, "}\n")
}

func NewReplaceableProcessor(tmplName string, slotName string, args []*Argument) *GoReplaceableProcessor {
	slotName = TrimWhiteSpaces(slotName)
	processor := &GoReplaceableProcessor{
		tmplName:   tmplName,
		slotName:   slotName,
		hiddenName: ToHiddenName(slotName),
		args:       args,
	}
	return processor
}
