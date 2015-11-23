package processors

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/util"
	"golang.org/x/net/html"
)

type Replacement struct {
	Name string
	Head Processor
	Args []*Argument
}

type GoCallProcessor struct {
	BaseProcessor
	pkgPath        string
	templateName   string
	closurePkgName string
	args           []*Argument
	replacements   []*Replacement
	callerAttrs    []html.Attribute
}

func (c *GoCallProcessor) Process(writer io.Writer, ctx *TagContext) {
	var argType string
	var newTemplateName string
	var pi pkgmgr.AliasGetter
	if c.pkgPath == "" {
		// In-package template call.
		switch ctx.OutputFormat {
		case "go":
			argType = fmt.Sprintf("%sTemplateArgs", c.templateName)
			newTemplateName = fmt.Sprintf("New%sTemplate", c.templateName)
		case "closure":
			require := fmt.Sprintf("%s.%sTemplate", c.closurePkgName, c.templateName)
			ctx.pkgRefs.RefClosureRequire(require)
			newTemplateName = require
		}
	} else {
		switch ctx.OutputFormat {
		case "go":
			pi = ctx.pkgRefs.RefByPath(c.pkgPath, false)
			argType = fmt.Sprintf("%s.%sTemplateArgs", pi.Alias(), c.templateName)
			newTemplateName = fmt.Sprintf("%s.New%sTemplate", pi.Alias(), c.templateName)
		case "closure":
			require := fmt.Sprintf("%s.%sTemplate", c.pkgPath, c.templateName)
			ctx.pkgRefs.RefClosureRequire(require)
			newTemplateName = require
		}
	}

	// Start of local scope
	io.WriteString(writer, "{\n")

	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("__args := &%s {}\n", argType))
	case "closure":
		io.WriteString(writer, fmt.Sprintf("var __args = {};\n"))
	}
	for _, argDef := range c.args {
		ctx.ExprParser.Evaluate(argDef.Val, writer, func(expr string) {
			switch ctx.OutputFormat {
			case "go":
				io.WriteString(writer, fmt.Sprintf("__args.%s = %s\n", util.ToPublicName(argDef.Name), expr))
			case "closure":
				if argDef.IsPb {
					io.WriteString(writer, fmt.Sprintf("__args[%s] = %s.getJsonData();\n", strconv.Quote(argDef.Name), expr))
				} else {
					io.WriteString(writer, fmt.Sprintf("__args[%s] = %s;\n", strconv.Quote(argDef.Name), expr))
				}
			}
		})
	}

	// Call template.
	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("__tplt := %s(__impl.GetWriter(), __impl.GetSettings())\n", newTemplateName))
	case "closure":
		io.WriteString(writer, fmt.Sprintf("__tplt = new %s();\n", newTemplateName))
	}
	// Caller Attributes.
	if c.callerAttrs != nil {
		switch ctx.OutputFormat {
		case "go":
			io.WriteString(writer, "__tplt.SetCallerAttrsFunc(func() (runtime.TagAttrs, bool, bool) {\n")
			io.WriteString(writer, "__callerAttrs := runtime.TagAttrs{}\n")
			io.WriteString(writer, "var __hasOmitTag bool\n")
			io.WriteString(writer, "var __omitTag bool\n")
			for _, attr := range c.callerAttrs {
				if attr.Key == "go:omit-tag" {
					io.WriteString(writer, "__hasOmitTag = true\n")
					v, err := ctx.RewriteExpression(attr.Val)
					if err != nil {
						panic(err)
					}
					io.WriteString(writer, fmt.Sprintf("__omitTag = %s\n", v))
				} else if attr.Key == "go:attr" {
					varName, varVal := util.SplitVarDef(attr.Val)
					ctx.ExprParser.Evaluate(varVal, writer, func(expr string) {
						io.WriteString(writer, fmt.Sprintf("__callerAttrs.AddAttr(\"%s\", %s)\n", varName, expr))
					})
				} else if !strings.HasPrefix(attr.Key, "go:") {
					// Static attributes
					io.WriteString(writer,
						fmt.Sprintf("__callerAttrs.AddAttr(\"%s\", \"%s\")\n", attr.Key, attr.Val))
				}
			}
			io.WriteString(writer, "return __callerAttrs, __hasOmitTag, __omitTag\n")
			io.WriteString(writer, "})\n")
		case "closure":
			io.WriteString(writer, "__tplt.setCallerAttrsFunc(function() {\n")
			io.WriteString(writer, "var __callerAttrs = {};\n")
			io.WriteString(writer, "var __hasOmitTag = false;")
			io.WriteString(writer, "var __omitTag = false;\n")
			for _, attr := range c.callerAttrs {
				if attr.Key == "go:omit-tag" {
					io.WriteString(writer, "__hasOmitTag = true;\n")
					v, err := ctx.RewriteExpression(attr.Val)
					if err != nil {
						panic(err)
					}
					io.WriteString(writer, fmt.Sprintf("__omitTag = %s;\n", v))
				} else if attr.Key == "go:attr" {
					varName, varVal := util.SplitVarDef(attr.Val)
					ctx.ExprParser.Evaluate(varVal, writer, func(expr string) {
						io.WriteString(writer, fmt.Sprintf("__callerAttrs[\"%s\"] = %s;\n", varName, expr))
					})
				} else if !strings.HasPrefix(attr.Key, "go:") {
					// Static attributes
					io.WriteString(writer, fmt.Sprintf("__callerAttrs[\"%s\"] = %s;\n", attr.Key, attr.Val))
				}
			}
			io.WriteString(writer, "return {\nattrs: __callerAttrs, hasOmitTag: __hasOmitTag, omitTag: __omitTag};\n")
			io.WriteString(writer, "});\n")
		}
	}
	// Replacements.
	for _, replacement := range c.replacements {
		switch ctx.OutputFormat {
		case "go":
			argType := fmt.Sprintf("%s%sReplArgs", c.templateName, replacement.Name)
			if c.pkgPath == "" {
				io.WriteString(writer,
					fmt.Sprintf("  __tplt.Replace%s(func(__args *%s) {\n", replacement.Name, argType))
			} else {
				io.WriteString(writer,
					fmt.Sprintf("  __tplt.Replace%s(func(__args *%s.%s) {\n",
						replacement.Name, pi.Alias(), argType))
			}

			for _, arg := range replacement.Args {
				io.WriteString(writer, fmt.Sprintf("  %s := __args.%s\n", arg.Name, util.ToPublicName(arg.Name)))
			}
			replacement.Head.Process(writer, ctx)

			io.WriteString(writer, "})\n")
		case "closure":
			if c.pkgPath == "" {
				io.WriteString(writer, fmt.Sprintf("  __tplt.replace%s(func(__args) {\n", replacement.Name))
			} else {
				io.WriteString(writer, fmt.Sprintf("  __tplt.replace%s(func(__args) {\n", replacement.Name))
			}

			for _, arg := range replacement.Args {
				io.WriteString(writer, fmt.Sprintf("  %s := __args[\"%s\"];\n", arg.Name, util.ToPublicName(arg.Name)))
			}
			replacement.Head.Process(writer, ctx)

			io.WriteString(writer, "});\n")
		}
	}

	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, "__tplt.Render(__element, __args);\n")
	case "closure":
		io.WriteString(writer, "__tplt.render(__element, __args);\n")
	}

	// Start of local scope.
	io.WriteString(writer, "}\n")

	// go:call is a terminal processor.
}

func NewCallProcessor(pkgPath, closurePkgName, templateName string, args []*Argument,
	replacements []*Replacement, callerAttrs []html.Attribute) *GoCallProcessor {
	processor := &GoCallProcessor{
		pkgPath:        pkgPath,
		closurePkgName: closurePkgName,
		templateName:   templateName,
		args:           args,
		replacements:   replacements,
		callerAttrs:    callerAttrs,
	}
	return processor
}
