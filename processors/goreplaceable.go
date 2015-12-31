package processors

import (
	"fmt"
	"io"

	"github.com/linuxerwang/goats-html/util"
)

type GoReplaceableProcessor struct {
	BaseProcessor
	tmplName   string
	slotName   string
	hiddenName string
	args       []*Argument
}

func (r *GoReplaceableProcessor) Process(writer io.Writer, ctx *TagContext) {
	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("if __impl.%s == nil {\n", r.hiddenName))
		for _, arg := range r.args {
			io.WriteString(writer, fmt.Sprintf("  %s := %s\n", arg.Name, arg.Val))
		}

		if r.next != nil {
			r.next.Process(writer, ctx)
		}

		io.WriteString(writer, "} else {\n")
		io.WriteString(writer, fmt.Sprintf("  args := &%s%sReplArgs{\n", r.tmplName, r.slotName))
		for _, arg := range r.args {
			io.WriteString(writer, fmt.Sprintf("    %s: %s,\n", util.ToPublicName(arg.Name), arg.Val))
		}
		io.WriteString(writer, fmt.Sprintf("  }\n"))
		io.WriteString(writer, fmt.Sprintf("  __impl.%s(args)\n", r.hiddenName))
		io.WriteString(writer, "}\n")
	case "closure":
		io.WriteString(writer, fmt.Sprintf("if (__self.%s_ == null) {\n", r.hiddenName))
		for _, arg := range r.args {
			io.WriteString(writer, fmt.Sprintf("  var %s = %s;\n", arg.Name, arg.Val))
		}

		if r.next != nil {
			r.next.Process(writer, ctx)
		}

		io.WriteString(writer, "} else {\n")
		io.WriteString(writer, fmt.Sprintf("  var __args = {};\n"))
		for _, arg := range r.args {
			io.WriteString(writer, fmt.Sprintf("__args.%s = %s;\n", util.ToPublicName(arg.Name), arg.Val))
		}
		io.WriteString(writer, fmt.Sprintf("  __self.%s_(__element, __args);\n", r.hiddenName))
		io.WriteString(writer, "}\n")
	}
}

func NewReplaceableProcessor(tmplName string, slotName string, args []*Argument) *GoReplaceableProcessor {
	slotName = util.TrimWhiteSpaces(slotName)
	processor := &GoReplaceableProcessor{
		tmplName:   tmplName,
		slotName:   slotName,
		hiddenName: util.ToHiddenName(slotName),
		args:       args,
	}
	return processor
}
