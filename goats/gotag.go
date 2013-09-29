package goats

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"fmt"
	"goats-html/goats/runtime"
	"io"
	"strings"
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

func (t *GoTagProcessor) Process(writer io.Writer, context *TagContext) {
	// There might be imports for constants or enums
	for _, attr := range t.attrs {
		if attr.Key == "go:attr" {
			parts := strings.Split(attr.Val, ":")
			context.MaybeAddImports(ToGoString(TrimWhiteSpaces(parts[1])))
		}
	}

	// Start local scope.
	io.WriteString(writer, "{\n")
	if t.firstTag {
		t.processFirstTag(writer, context)
		t.processChildrenTags(writer, context)
		t.maybeCloseTag(writer)
	} else {
		t.processSubseqTag(writer, context)
		t.processChildrenTags(writer, context)
		t.maybeCloseTag(writer)
	}
	// End local scope.
	io.WriteString(writer, "}\n")
}

func (t *GoTagProcessor) processFirstTag(writer io.Writer, context *TagContext) {
	// Vars for attributes
	io.WriteString(writer, "var __omitTag bool = false\n")
	io.WriteString(writer, "var __attrs = &runtime.TagAttrs{}\n")
	// Local attributes
	t.genLocalAttrs(writer, context)
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
		t.genConditionalTag(writer)
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
		t.genConditionalTag(writer)
	}
	io.WriteString(writer, "}\n")
}

func (t *GoTagProcessor) processSubseqTag(writer io.Writer, context *TagContext) {
	switch t.omitTag {
	case "true":
		// Do nothing
	case "", "false":
		if t.hasOnlyStaticAttrs() {
			attrs := &runtime.TagAttrs{}
			for _, attr := range t.attrs {
				if !strings.HasPrefix(attr.Key, "go:") {
					attrs.AddAttr(attr.Key, attr.Val)
				}
			}
			var tagBuffer bytes.Buffer
			attrs.GenTagAndAttrs(&tagBuffer, t.tagName)
			io.WriteString(writer,
				fmt.Sprintf("__impl.WriteString(\"%s\")\n",
					strings.Replace(tagBuffer.String(), "\"", "\\\"", -1)))
		} else {
			t.genLocalAttrs(writer, context)
			io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
		}
	default:
		io.WriteString(writer, fmt.Sprintf("if !%s {\n", t.omitTag))
		io.WriteString(writer, "var __omitTag bool = false\n")
		io.WriteString(writer, "  var __attrs = &runtime.TagAttrs{}\n")
		io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
		t.genLocalAttrs(writer, context)
		t.genConditionalTag(writer)
		io.WriteString(writer, "} else {\n")
		io.WriteString(writer, "  __omitTag = true\n")
		// Do nothing
		io.WriteString(writer, "}\n")
	}
}

func (t *GoTagProcessor) genConditionalTag(writer io.Writer) {
	io.WriteString(writer, fmt.Sprintf("__omitTag = %s\n", t.omitTag))
	io.WriteString(writer, "if !__omitTag {\n")
	io.WriteString(writer, fmt.Sprintf("  __attrs.GenTagAndAttrs(__impl.GetWriter(), \"%s\")\n", t.tagName))
	io.WriteString(writer, "}\n")
}

func (t *GoTagProcessor) genLocalAttrs(writer io.Writer, context *TagContext) {
	for _, attr := range t.attrs {
		if attr.Key == "go:attr" {
			parts := strings.Split(attr.Val, ":")
			expr := context.RewriteExpression(parts[1])
			io.WriteString(writer,
				fmt.Sprintf("__attrs.AddAttr(\"%s\", %s)\n", TrimWhiteSpaces(parts[0]), expr))
		} else if !strings.HasPrefix(attr.Key, "go:") {
			// Static attrs
			io.WriteString(writer,
				fmt.Sprintf("__attrs.AddAttr(\"%s\", \"%s\")\n",
					TrimWhiteSpaces(attr.Key), TrimWhiteSpaces(attr.Val)))
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

func (t *GoTagProcessor) processChildrenTags(writer io.Writer, context *TagContext) {
	head := NewHeadProcessor()
	if t.next != nil {
		head.SetNext(t.next)
	}
	tail := t.getTail(head)
	callback := NewCallbackProcessor(func() {
		for _, child := range t.childProcessors {
			child.Process(writer, context)
		}
	})
	tail.SetNext(callback)
	head.Process(writer, context)
}

func (t *GoTagProcessor) maybeCloseTag(writer io.Writer) {
	if t.needsClosing {
		if t.firstTag || t.isDynamicOmitTag() {
			io.WriteString(writer, "if !__omitTag {\n")
			t.closeTag(writer)
			io.WriteString(writer, "}\n")
		} else if t.omitTag == "" || t.omitTag == "false" {
			t.closeTag(writer)
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
