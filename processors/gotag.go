package processors

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/linuxerwang/goats-html/runtime"
	"github.com/linuxerwang/goats-html/util"
	"golang.org/x/net/html"
)

/**
 * TODO: Fix me. Logic is out of date.
 * The logic to handle tag processing:
 * if tag is the first:
 *     if caller has omitTag:
 *         if caller omitTag evals to false:
 *             call genTagWithDynamicAttrs()
 *         else:
 *             // Do nothing
 *     else: (caller has not omitTag)
 *         if tag has omitTag:
 *             if eval(omitTag) == false:
 *                 call genTagWithDynamicAttrs()
 *         else:
 *             call genTagWithDynamicAttrs()
 * else: (tag is subsequent)
 *     if omitTag == "":
 *         omitTag <- "false"
 *     if omitTag == "true":
 *         // Do nothing
 *     else if omitTag == "false":
 *         if has only static attrs:
 *             call genTagWithStaticAttrs()
 *         else:
 *             call genTagWithDynamicAttrs()
 *     else: (has omitTag expression)
 *         if eval(omitTag) == true:
 *             call genReturn()
 *         else:
 *             call genTagWithDynamicAttrs()
 * processChildrenTags()
 *
 */

type GoTagProcessor struct {
	BaseProcessor
	tagName         string
	omitTag         string
	hasOmitTag      bool
	firstTag        bool // Whether it's the first tag of the template.
	needsClosing    bool
	childProcessors []Processor
	attrs           []html.Attribute
}

func (t *GoTagProcessor) Process(writer io.Writer, ctx *TagContext) {
	originalAutoEscape := ctx.AutoEscape
	// There might be imports for constants or enums
	for _, attr := range t.attrs {
		if attr.Key == "go:autoescape" {
			ctx.AutoEscape = (attr.Val == "true")
			break
		}
	}

	// Start local scope.
	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, "{\n")
	case "closure":
		io.WriteString(writer, "(function(){\n")
	}

	switch ctx.OutputFormat {
	case "closure":
		io.WriteString(writer, "var __element = null;\n")
		io.WriteString(writer, "var __attrs = new goats.runtime.TagAttrs();\n")
	}

	if t.firstTag {
		t.processFirstTag(writer, ctx)
	} else {
		t.processSubseqTag(writer, ctx)
	}

	switch ctx.OutputFormat {
	case "closure":
		io.WriteString(writer, "if (__element) {\n")
		io.WriteString(writer, "goog.dom.appendChild(__tag_stack[__tag_stack.length-1], __element);\n")
		io.WriteString(writer, "__tag_stack.push(__element);\n")
		io.WriteString(writer, "}\n")
	}

	t.processChildrenTags(writer, ctx)

	switch ctx.OutputFormat {
	case "closure":
		io.WriteString(writer, "if (__element) {\n")
		io.WriteString(writer, "__tag_stack.pop();\n")
		io.WriteString(writer, "}\n")
	}

	t.maybeCloseTag(writer, ctx)

	// End local scope.
	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, "}\n")
	case "closure":
		io.WriteString(writer, "})();\n")
	}

	ctx.AutoEscape = originalAutoEscape
}

func (t *GoTagProcessor) processFirstTag(writer io.Writer, ctx *TagContext) {
	switch ctx.OutputFormat {
	case "go":
		// Vars for attributes
		io.WriteString(writer, "var __omitTag bool = false\n")
		io.WriteString(writer, "var __attrs = runtime.TagAttrs{}\n")
		// Local attributes
		t.genLocalAttrs(writer, ctx)
		// Caller attributes may overwrite local attrs.
		io.WriteString(writer, "if __callerAttrsFunc := __impl.GetCallerAttrsFunc(); __callerAttrsFunc != nil {\n")
		io.WriteString(writer, "  var __callerAttrs runtime.TagAttrs\n")
		io.WriteString(writer, "  __callerAttrs, __callerHasOmitTag, __callerOmitTag := __callerAttrsFunc()\n")
		io.WriteString(writer, "  __attrs.MergeFrom(__callerAttrs)\n")
		io.WriteString(writer, "  if __callerHasOmitTag {\n")
		io.WriteString(writer, "    if __callerOmitTag {\n")
		io.WriteString(writer, "      __omitTag = true\n")
		io.WriteString(writer, "    } else {\n")
		io.WriteString(writer, fmt.Sprintf("    __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
		io.WriteString(writer, "    }\n")
		io.WriteString(writer, "  } else {\n")
		if t.hasOmitTag {
			t.genConditionalTag(writer, ctx)
		} else {
			io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
		}
		io.WriteString(writer, "  }\n")
		io.WriteString(writer, "} else {\n")
		if !t.hasOmitTag || t.omitTag == "false" {
			io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
		} else if t.omitTag == "true" {
			io.WriteString(writer, "  __omitTag = true\n")
			// Do not output tag.
		} else {
			t.genConditionalTag(writer, ctx)
		}
		io.WriteString(writer, "}\n")
	case "closure":
		// Vars for attributes
		io.WriteString(writer, "var __omitTag = false;\n")
		// Local attributes
		t.genLocalAttrs(writer, ctx)
		// Caller attributes may overwrite local attrs.
		io.WriteString(writer, "if (this.callerAttrsFunc_) {\n")
		io.WriteString(writer, "  var __callerAttrs = this.callerAttrsFunc_();\n")
		io.WriteString(writer, "  if (__callerAttrs) {\n")
		io.WriteString(writer, "    __attrs.mergeFrom(__callerAttrs.attrs);\n")
		io.WriteString(writer, "    if (__callerAttrs.hasOmitTag) {\n")
		io.WriteString(writer, "      if (__callerAttrs.omitTag) {\n")
		io.WriteString(writer, "        __omitTag = true;\n")
		io.WriteString(writer, "      }\n")
		io.WriteString(writer, "    } else {\n")
		if t.hasOmitTag {
			io.WriteString(writer, fmt.Sprintf("      __omitTag = %s;\n", t.omitTag))
		}
		io.WriteString(writer, "    }\n")
		io.WriteString(writer, "  }\n")
		io.WriteString(writer, "} else {\n")
		if t.hasOmitTag {
			io.WriteString(writer, fmt.Sprintf("  __omitTag = %s;\n", t.omitTag))
		}
		io.WriteString(writer, "}\n")
		io.WriteString(writer, "if (!__omitTag) {\n")
		io.WriteString(writer, fmt.Sprintf("  __element = goog.dom.createDom(\"%s\", __attrs.get());\n", t.tagName))
		io.WriteString(writer, "}\n")
	}
}

func (t *GoTagProcessor) processSubseqTag(writer io.Writer, ctx *TagContext) {
	switch t.omitTag {
	case "true":
		// Do nothing
	case "", "false":
		if t.hasOnlyStaticAttrs() {
			attrs := runtime.TagAttrs{}
			for _, attr := range t.attrs {
				if !strings.HasPrefix(attr.Key, "go:") {
					attrs.AddAttr(attr.Key, attr.Val)
				}
			}
			switch ctx.OutputFormat {
			case "go":
				var tagBuffer bytes.Buffer
				attrs.GenTagAndAttrs(&tagBuffer, t.tagName)
				io.WriteString(writer,
					fmt.Sprintf("__impl.WriteString(\"%s\")\n",
						strings.Replace(tagBuffer.String(), "\"", "\\\"", -1)))
			case "closure":
				io.WriteString(writer, fmt.Sprintf("__element = goog.dom.createElement(\"%s\");\n", t.tagName))
				for k, v := range map[string]string(attrs) {
					io.WriteString(writer, fmt.Sprintf("__element.setAttribute(\"%s\", \"%s\");\n", k, v))
				}
			}
		} else {
			switch ctx.OutputFormat {
			case "go":
				io.WriteString(writer, "  var __attrs = runtime.TagAttrs{}\n")
				t.genLocalAttrs(writer, ctx)
				io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
			case "closure":
				t.genLocalAttrs(writer, ctx)
				io.WriteString(writer, fmt.Sprintf("__element = goog.dom.createDom(\"%s\", __attrs.get());\n", t.tagName))
			}
		}
	default:
		switch ctx.OutputFormat {
		case "go":
			io.WriteString(writer, fmt.Sprintf("if !%s {\n", t.omitTag))
			io.WriteString(writer, "var __omitTag bool = false\n")
			io.WriteString(writer, "  var __attrs = runtime.TagAttrs{}\n")
			io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
			t.genLocalAttrs(writer, ctx)
			t.genConditionalTag(writer, ctx)
			io.WriteString(writer, "} else {\n")
			io.WriteString(writer, "  __omitTag = true\n")
			// Do nothing
			io.WriteString(writer, "}\n")
		case "closure":
			io.WriteString(writer, fmt.Sprintf("if (!%s) {\n", t.omitTag))
			io.WriteString(writer, "var __omitTag = false;\n")
			io.WriteString(writer, "  var __attrs = {};\n")
			io.WriteString(writer, fmt.Sprintf("__element = goog.dom.createDom(\"%s\", __attrs.get());\n", t.tagName))
			t.genLocalAttrs(writer, ctx)
			t.genConditionalTag(writer, ctx)
			io.WriteString(writer, "} else {\n")
			io.WriteString(writer, "  __omitTag = true;\n")
			// Do nothing
			io.WriteString(writer, "}\n")
		}
	}
}

func (t *GoTagProcessor) genConditionalTag(writer io.Writer, ctx *TagContext) {
	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, fmt.Sprintf("__omitTag = %s\n", t.omitTag))
		io.WriteString(writer, "if !__omitTag {\n")
		io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
		io.WriteString(writer, "}\n")
	case "closure":
		io.WriteString(writer, fmt.Sprintf("__omitTag = %s;\n", t.omitTag))
		io.WriteString(writer, "if (!__omitTag) {\n")
		io.WriteString(writer, fmt.Sprintf("  __element = goog.dom.createDom(\"%s\", __attrs.get());\n", t.tagName))
		io.WriteString(writer, "}\n")
	}
}

func (t *GoTagProcessor) genLocalAttrs(writer io.Writer, ctx *TagContext) {
	for _, attr := range t.attrs {
		if attr.Key == "go:attr" {
			varName, varVal := util.SplitVarDef(attr.Val)
			ctx.ExprParser.Evaluate(varVal, writer, func(expr string) {
				switch ctx.OutputFormat {
				case "go":
					io.WriteString(writer, fmt.Sprintf("__attrs.AddAttr(\"%s\", %s)\n", varName, expr))
				case "closure":
					io.WriteString(writer, fmt.Sprintf("__attrs.add(\"%s\", %s);\n", varName, expr))
				}
			})
		} else if !strings.HasPrefix(attr.Key, "go:") {
			// Static attrs
			switch ctx.OutputFormat {
			case "go":
				io.WriteString(writer,
					fmt.Sprintf("__attrs.AddAttr(\"%s\", \"%s\")\n",
						util.TrimWhiteSpaces(attr.Key), util.TrimWhiteSpaces(attr.Val)))
			case "closure":
				io.WriteString(writer, fmt.Sprintf("__attrs.add(\"%s\", \"%s\");\n", util.TrimWhiteSpaces(attr.Key), util.TrimWhiteSpaces(attr.Val)))
			}
		}
	}
}

func (t *GoTagProcessor) hasOnlyStaticAttrs() bool {
	for _, attr := range t.attrs {
		if attr.Key == "go:attr" {
			return false
		}
	}
	return true
}

func (t *GoTagProcessor) processChildrenTags(writer io.Writer, ctx *TagContext) {
	originalAutoEscape := ctx.AutoEscape
	if t.tagName == "script" || t.tagName == "style" {
		ctx.AutoEscape = false
	}

	if t.next != nil {
		t.next.Process(writer, ctx)
	}

	for _, child := range t.childProcessors {
		child.Process(writer, ctx)
	}

	ctx.AutoEscape = originalAutoEscape
}

func (t *GoTagProcessor) maybeCloseTag(writer io.Writer, ctx *TagContext) {
	switch ctx.OutputFormat {
	case "go":
		if t.needsClosing {
			if t.firstTag || t.isDynamicOmitTag() {
				io.WriteString(writer, "if (!__omitTag) {\n")
				t.closeTag(writer)
				io.WriteString(writer, "}\n")
			} else if t.omitTag == "" || t.omitTag == "false" {
				t.closeTag(writer)
			}
		}
	}
}

func (t *GoTagProcessor) closeTag(writer io.Writer) {
	io.WriteString(writer, fmt.Sprintf("__impl.WriteString(\"</%s>\")\n", t.tagName))
}

func (t *GoTagProcessor) isDynamicOmitTag() bool {
	return t.omitTag != "" && t.omitTag != "true" && t.omitTag != "false"
}

func (t *GoTagProcessor) getTail(processor Processor) Processor {
	tail := processor
	for tail.HasNext() {
		tail = tail.GetNext()
	}
	return tail
}

func (t *GoTagProcessor) AddChild(child Processor) {
	t.childProcessors = append(t.childProcessors, child)
}

func NewTagProcessor(tagName string, omitTag string, firstTag bool, needsClosing bool,
	attrs []html.Attribute) *GoTagProcessor {
	processor := &GoTagProcessor{
		tagName:      tagName,
		omitTag:      omitTag,
		hasOmitTag:   omitTag != "",
		firstTag:     firstTag,
		needsClosing: needsClosing,
		attrs:        attrs,
	}
	return processor
}
