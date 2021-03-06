package goats

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"path/filepath"
	"strings"
	txttpl "text/template"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/processors"
	"github.com/linuxerwang/goats-html/util"
	"golang.org/x/net/html"
)

const (
	// Enum for tag processing types.
	TagProcessingGoRegular = iota
	TagProcessingGoSwitch
)

const tmplIfaceBody = (`type {{.Name}}TemplateArgs struct {
	{{range .Args}} {{.Declare}}
{{end}}
}

{{$tmplName := .Name}}
{{range .Replaceables}}
type {{$tmplName}}{{.Name}}ReplArgs struct {
	{{range .Args}} {{.Declare}}
{{end}}
}

type {{$tmplName}}{{.Name}}ReplFunc func(*{{$tmplName}}{{.Name}}ReplArgs)

{{end}}

type {{.Name}}Template interface {
	runtime.Template
	Render(*{{.Name}}TemplateArgs) error
{{range .Replaceables}}
	Replace{{.Name}}({{$tmplName}}{{.Name}}ReplFunc)
{{end}}
}
`)

const tmplImplBody = (`{{$tmplName := .Name}}
type {{$tmplName}}TemplateImpl struct {
	*runtime.BaseTemplate
	builtinFilter *runtime.BuiltinFilter
{{range .Replaceables}}
	{{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc
{{end}}
}

func (__impl *{{$tmplName}}TemplateImpl) Render(__args *{{.Name}}TemplateArgs) error {
	@@RENDER@@
	return nil
}

{{range .Replaceables}}
func (__impl *{{$tmplName}}TemplateImpl) Replace{{.Name}}({{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc) {
	__impl.{{.HiddenName}} = {{.HiddenName}}
}
{{end}}

func New{{.Name}}Template(writer io.Writer, settings *runtime.TemplateSettings) {{.Name}}Template {
	template := &{{.Name}}TemplateImpl{}
	template.BaseTemplate = runtime.NewBaseTemplate(writer, settings)
	template.builtinFilter = runtime.NewBuiltinFilter()
	return template
}
`)

const tmplClosureBody = (`{{$tmplName := .Name}}


/**
 * The {{$tmplName}} template for packge {{.ClosurePkgName}}.
 *
 * @constructor
{{if .OutputExport}} * @export
{{end}} */
{{.ClosurePkgName}}.{{$tmplName}}Template = function() {
	/**
	 * The caller attrs function.
	 * @private {Function}
	 */
	this.callerAttrsFunc_ = null;
};

/**
 * The render method renders the template.
 *
 * @param {Element} __parent The parent element.
 * @param {Object} __args The template arguments.
{{if .OutputExport}} * @export
{{end}} */
{{.ClosurePkgName}}.{{$tmplName}}Template.prototype.render = function(__parent, __args) {
	var __tag_stack = [];
	if (__parent) {
		__tag_stack.push(__parent);
	}

	var __self = this;

@@RENDER@@
};

/**
 * Sets the caller attrs function.
 *
 * @param {Function} callerAttrsFunc The function to pass the caller's attrs.
 */
{{.ClosurePkgName}}.{{$tmplName}}Template.prototype.setCallerAttrsFunc = function(callerAttrsFunc) {
	this.callerAttrsFunc_ = callerAttrsFunc;
};

{{$ClosurePkgName := .ClosurePkgName}}
{{range .Replaceables}}/**
 * The replacement function for replaceable {{.HiddenName}}.
 *
 * @param {string} {{.HiddenName}} The name of the replaceable.
 */
{{$ClosurePkgName}}.{{$tmplName}}Template.prototype.replace{{.Name}} = function({{.HiddenName}}) {
	/*
	 * @private {Function}
	 */
	this.{{.HiddenName}}_ = {{.HiddenName}};
};
{{end}}
`)

const tmplProxyBody = (`{{$tmplName := .Name}}
type {{.HiddenName}}TemplateProxy struct {
	*runtime.BaseTemplate
{{range .Replaceables}}  {{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc
{{end}}
}

func (__proxy *{{.HiddenName}}TemplateProxy) Render(args *{{.Name}}TemplateArgs) error {
	err := runtime.CallRpc("{{.Pkg}}",
		"{{.Name}}",
		__proxy.GetSettings(),
		args,
		__proxy.GetWriter())
	return err
}

{{$name := .HiddenName}}

{{range .Replaceables}}  func (__impl *{{$name}}TemplateProxy) Replace{{.Name}}({{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc) {
	__impl.{{.HiddenName}} = {{.HiddenName}}
}
{{end}}

func New{{.Name}}Template(writer io.Writer, settings *runtime.TemplateSettings) {{.Name}}Template {
	template := &{{.HiddenName}}TemplateProxy{}
	template.BaseTemplate = runtime.NewBaseTemplate(writer, settings)
	return template
}
`)

const tmplMainFile = (`package main
import(
	"bytes"
	"{{.Pkg}}"
	"github.com/linuxerwang/goats-html/runtime"
	"os"
)

func main() {
	settings := runtime.TemplateSettings{}
	var buffer bytes.Buffer
	switch os.Args[1] {
{{range .Templates}}  case "{{.Name}}":
		args := {{.PkgName}}.{{.Name}}TemplateArgs{}
		runtime.DecodeRpcRequestOrFail(os.Stdin, &settings, &args)
		template := {{.PkgName}}.New{{.Name}}Template(&buffer, &settings)
		template.Render(&args)
{{end}}
	default:
		panic("Unknown template name: " + os.Args[1])
	}
	os.Stdout.Write(buffer.Bytes())
}
`)

const (
	ImplFileGoSuffix      = "_impl.go"
	ImplFileClosureSuffix = ".closure.js"
	ProxyFileSuffix       = "_proxy.go"
)

var (
	MainFileName = filepath.Join("cmd", "main.go")
)

// List of void elements. Void elements are those that can't have any contents.
var voidElements = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

var multipleAttrs = map[string]bool{
	"go:arg":    true,
	"go:attr":   true,
	"go:import": true,
	"go:var":    true,
}

func formatSource(unformated string) string {
	formated, err := format.Source([]byte(unformated))
	if err != nil {
		log.Println("Failed to format the output template, ", err)
		return unformated
	}
	return string(formated)
}

type GoatsReplaceable struct {
	Name       string
	HiddenName string
	Args       []*processors.Argument
}

type GoatsReplace struct {
	Name       string
	HiddenName string
	Args       []*processors.Argument
}

type GoatsTemplate struct {
	Parser          *GoatsParser
	OutputPath      string
	OutputIfaceFile string
	OutputImplFile  string
	OutputProxyFile string
	OutputExport    bool
	Pkg             string
	PkgName         string
	ClosurePkgName  string
	Name            string
	HiddenName      string
	Args            []*processors.Argument
	RootNode        *html.Node
	NeedsDocType    bool
	Replaceables    []*GoatsReplaceable
	Replaces        []*GoatsReplace
	pkgRefs         *pkgmgr.PkgRefs
}

func NewGoatsTemplate(parser *GoatsParser, tmplName string, args []*processors.Argument,
	rootNode *html.Node, needsDocType bool, pkgRefs *pkgmgr.PkgRefs) *GoatsTemplate {
	prefix := util.ToSnakeCase(tmplName)
	suffix := ImplFileGoSuffix
	if parser.Settings.OutputFormat == "closure" {
		suffix = ImplFileClosureSuffix
	}
	return &GoatsTemplate{
		Parser:          parser,
		OutputPath:      parser.OutputPath,
		OutputIfaceFile: fmt.Sprintf("%s.go", prefix),
		OutputImplFile:  fmt.Sprintf("%s%s", prefix, suffix),
		OutputProxyFile: fmt.Sprintf("%s%s", prefix, ProxyFileSuffix),
		OutputExport:    parser.Settings.OutputExport,
		Pkg:             parser.Pkg,
		PkgName:         filepath.Base(parser.Pkg),
		Name:            tmplName,
		HiddenName:      util.ToHiddenName(tmplName),
		Args:            args,
		RootNode:        rootNode,
		NeedsDocType:    needsDocType,
		pkgRefs:         pkgRefs,
	}
}

func (t *GoatsTemplate) IsDirty() bool {
	if t.Parser.IsFileOld(t.OutputIfaceFile) {
		return true
	}
	if t.Parser.IsFileOld(t.OutputImplFile) {
		return true
	}
	if t.Parser.IsFileOld(t.OutputProxyFile) {
		return true
	}
	return false
}

func (t *GoatsTemplate) genIfacePkgDecl(output io.Writer) {
	io.WriteString(output, fmt.Sprintf("package %s\n\n", t.PkgName))
}

func (t *GoatsTemplate) genIfaceImports(output io.Writer) {
	io.WriteString(output, "\"github.com/linuxerwang/goats-html/runtime\"\n")
	t.pkgRefs.GenerateImports(output, pkgmgr.GenInterfaceImports)
}

func (t *GoatsTemplate) genIfaceBody(output io.Writer) {
	tmpl, err := txttpl.New("iface-body").Parse(tmplIfaceBody)
	if err != nil {
		log.Fatal("Failed to generate interface body,", err)
	}

	err = tmpl.Execute(output, t)
	if err != nil {
		log.Fatal("Failed to generate interface body,", err)
	}
}

func (t *GoatsTemplate) genImplPkgDecl(output io.Writer) {
	// Note that to make build tags to work there must be an empty line between the
	// build tags line and the package line.
	io.WriteString(output, "// +build !goats_devmod\n\n")
	io.WriteString(output, fmt.Sprintf("package %s\n\n", t.PkgName))
}

func (t *GoatsTemplate) genImplImports(output io.Writer) {
	io.WriteString(output, "\"github.com/linuxerwang/goats-html/runtime\"\n")
	io.WriteString(output, "\"io\"\n")
	t.pkgRefs.GenerateImports(output, pkgmgr.GenImplImports)
}

func (t *GoatsTemplate) genImplBody(output io.Writer) {
	tmpl, err := txttpl.New("impl-body").Parse(tmplImplBody)
	if err != nil {
		log.Fatal("Failed to generate implementation body,", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, t)
	if err != nil {
		log.Fatal("Failed to generate implementation body,", err)
	}

	io.WriteString(output, strings.Replace(buffer.String(), "@@RENDER@@", t.genRenderContent(output), 1))
}

func (t *GoatsTemplate) genProxyPkgDecl(output io.Writer) {
	// Note that to make build tags to work there must be an empty line between the
	// build tags line and the package line.
	io.WriteString(output, "// +build goats_devmod\n\n")
	io.WriteString(output, fmt.Sprintf("package %s\n\n", t.PkgName))
}

func (t *GoatsTemplate) genProxyImports(output io.Writer) {
	io.WriteString(output, "\"github.com/linuxerwang/goats-html/runtime\"\n")
	io.WriteString(output, "\"io\"\n")
}

func (t *GoatsTemplate) genProxyBody(output io.Writer) {
	tmpl, err := txttpl.New("proxy-body").Parse(tmplProxyBody)
	if err != nil {
		log.Fatal("Failed to generate proxy body,", err)
	}

	err = tmpl.Execute(output, t)
	if err != nil {
		log.Fatal("Failed to generate proxy body,", err)
	}
}

func (t *GoatsTemplate) genClosurePkgDoc(output io.Writer) {
	io.WriteString(output, fmt.Sprintf("/**\n * @fileoverview Template.\n */\n\n"))
}

func (t *GoatsTemplate) genClosureProvides(output io.Writer) {
	io.WriteString(output, fmt.Sprintf("goog.provide('%s.%sTemplate');\n", t.ClosurePkgName, t.Name))
}

func (t *GoatsTemplate) genClosureCommonRequires(output io.Writer) {
	io.WriteString(output, "goog.require('goats.runtime.TagAttrs');\n")
	io.WriteString(output, "goog.require('goog.dom');\n")
}

func (t *GoatsTemplate) dumpClosureCommonRequires(requires map[string]bool) {
	requires["goats.runtime.TagAttrs"] = true
	requires["goog.dom"] = true
}

func (t *GoatsTemplate) genClosureRequires(output io.Writer) {
	t.pkgRefs.GenerateRequires(output)
}

func (t *GoatsTemplate) dumpClosureRequires(requires map[string]bool) {
	t.pkgRefs.DumpRequires(requires)
}

func (t *GoatsTemplate) genClosureBody(output io.Writer) {
	tmpl, err := txttpl.New("closure-body").Parse(tmplClosureBody)
	if err != nil {
		log.Fatal("Failed to generate closure body,", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, t)
	if err != nil {
		log.Fatal("Failed to generate closure body,", err)
	}

	io.WriteString(output, strings.Replace(buffer.String(), "@@RENDER@@", t.genRenderContent(output), 1))
}

// genRenderContent traverses the HTML node and render the template.
func (t *GoatsTemplate) genRenderContent(output io.Writer) string {
	var headProcessor processors.Processor = processors.NewHeadProcessor()

	var argProcessor processors.Processor = processors.NewArgProcessor(t.Args)
	headProcessor.SetNext(argProcessor)

	t.buildProcessorChain(argProcessor, t.RootNode)

	ctx := processors.NewTagContext(t.Parser.PkgMgr, t.pkgRefs, t.Parser.Settings.OutputFormat)
	if t.NeedsDocType {
		docTypeProcessor := processors.NewDocTypeProcessor(t.Parser.DocTypeTag, t.Parser.DocTypeAttrs)
		docTypeProcessor.SetNext(headProcessor)
		headProcessor = docTypeProcessor
	}

	var renderBuffer bytes.Buffer
	headProcessor.Process(&renderBuffer, ctx)
	return renderBuffer.String()
}

func (t *GoatsTemplate) buildProcessorChain(preProcessor processors.Processor, node *html.Node) {
	if node.Type == html.CommentNode {
		if t.Parser.Settings.KeepComments {
			processor := processors.NewCommentProcessor(node.Data)
			preProcessor.SetNext(processor)
			preProcessor = processor
		}
		return
	}

	if node.Type == html.TextNode {
		processor := processors.NewTextProcessor(node.Data)
		preProcessor.SetNext(processor)
		preProcessor = processor
		return
	}

	if node.Type != html.ElementNode {
		panic(fmt.Sprintf("Expect element node but got node type %d", node.Type))
		return
	}

	goAttrs := t.getAttrMap(node)

	if val, ok := goAttrs["go:settings"]; ok {
		settingsProcessor := processors.NewSettingsProcessor(val)
		preProcessor.SetNext(settingsProcessor)
		preProcessor = settingsProcessor
	}

	if val, ok := goAttrs["go:var"]; ok {
		varProcessor := processors.NewVarsProcessor(val)
		preProcessor.SetNext(varProcessor)
		preProcessor = varProcessor
	}

	if val, ok := goAttrs["go:if"]; ok {
		ifProcessor := processors.NewIfProcessor(val)
		preProcessor.SetNext(ifProcessor)
		preProcessor = ifProcessor
	} else if val, ok := goAttrs["go:case"]; ok {
		caseProcessor := processors.NewCaseProcessor(val)
		preProcessor.SetNext(caseProcessor)
		preProcessor = caseProcessor
	} else if _, ok := goAttrs["go:default"]; ok {
		defaultProcessor := processors.NewDefaultProcessor()
		preProcessor.SetNext(defaultProcessor)
		preProcessor = defaultProcessor
	}

	if val, ok := goAttrs["go:for"]; ok {
		forProcessor := processors.NewForProcessor(val)
		preProcessor.SetNext(forProcessor)
		preProcessor = forProcessor
	} else if _, ok := goAttrs["go:switch"]; ok {
		// Tag with go:switch can only contain sub tags with go:case and/or go:default.
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				attrs := t.getAttrMap(c)
				_, hasCase := attrs["go:case"]
				_, hasDefault := attrs["go:default"]
				if !hasCase && !hasDefault {
					log.Fatal(
						"Tag with go:switch can only contain sub tags with go:case and/or go:default.")
				}

				t.handleTag(preProcessor, node, TagProcessingGoSwitch)
			}
		}
		return
	}

	if val, ok := goAttrs["go:template"]; ok && node != t.RootNode {
		// Convert to an in-package template call.
		callProcessor := processors.NewCallProcessor("", "", t.Parser.Settings.OutputPkgPrefix, t.ClosurePkgName, val, processors.ParseArgDefs(goAttrs["go:arg"]), nil, node.Attr)
		preProcessor.SetNext(callProcessor)
		preProcessor = callProcessor
		return
	}

	if val, ok := goAttrs["go:replaceable"]; ok && node != t.RootNode {
		replaceableProcessor := processors.NewReplaceableProcessor(t.Name, val, processors.ParseArgDefs(goAttrs["go:arg"]))
		preProcessor.SetNext(replaceableProcessor)
		preProcessor = replaceableProcessor
	}

	if val, ok := goAttrs["go:call"]; ok {
		pkgPath, relPkgPath, callName := t.pkgRefs.ParseTmplCall(val)

		var replacements []*processors.Replacement
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				if len(util.TrimWhiteSpaces(c.Data)) != 0 {
					log.Fatal("Node with go:call can only contain nodes with go:replace or spaces.")
				}
				continue
			}

			var found bool
			replacement := &processors.Replacement{
				Args: []*processors.Argument{},
			}
			for _, attr := range c.Attr {
				if attr.Key == "go:replace" {
					found = true
					head := processors.NewHeadProcessor()
					t.buildProcessorChain(head, c)
					replacement.Name = attr.Val
					replacement.Head = head
					replacements = append(replacements, replacement)
				} else if attr.Key == "go:arg" {
					replacement.Args = append(replacement.Args, processors.NewArgDef(attr.Val))
				}
			}
			if !found {
				log.Fatal("Node with go:call can only contain nodes with go:replace.")
			}
		}

		callProcessor := processors.NewCallProcessor(
			pkgPath, relPkgPath, t.Parser.Settings.OutputPkgPrefix, t.ClosurePkgName, callName, processors.ParseArgCalls(goAttrs["go:arg"]), replacements, node.Attr)
		preProcessor.SetNext(callProcessor)
		preProcessor = callProcessor

		return
	}
	t.handleTag(preProcessor, node, TagProcessingGoRegular)
}

func (t *GoatsTemplate) handleTag(preProcessor processors.Processor, node *html.Node, tagProcessingType int) {
	// Static tag attributes.
	var nonGoAttrs []html.Attribute
	for _, attr := range node.Attr {
		if !strings.HasPrefix(attr.Key, "go:") {
			nonGoAttrs = append(nonGoAttrs, attr)
		}
	}

	goAttrs := t.getAttrMap(node)
	omitTag := ""
	if val, ok := goAttrs["go:omit-tag"]; ok {
		omitTag = val
	}
	var firstTag bool
	if _, ok := goAttrs["go:template"]; ok {
		firstTag = true
	}

	tagProcessor := processors.NewTagProcessor(node.Data, omitTag, firstTag, !voidElements[node.Data], node.Attr)
	preProcessor.SetNext(tagProcessor)
	preProcessor = tagProcessor

	if tagProcessingType == TagProcessingGoSwitch {
		switchProcessor := processors.NewSwitchProcessor(goAttrs["go:switch"])
		preProcessor.SetNext(switchProcessor)
		preProcessor = switchProcessor
	}

	if val, ok := goAttrs["go:content"]; ok {
		contentProcessor := processors.NewContentProcessor(val)
		preProcessor.SetNext(contentProcessor)
		preProcessor = contentProcessor
	} else {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			currentNode := c

			if c.Type == html.TextNode {
				var buffer bytes.Buffer
				io.WriteString(&buffer, c.Data)

				var n *html.Node
				for n = c.NextSibling; n != nil && n.Type == html.TextNode; n = n.NextSibling {
					// Merge consecutive text nodes.
					io.WriteString(&buffer, n.Data)
				}

				currentNode = &html.Node{
					Type:        html.TextNode,
					Data:        buffer.String(),
					PrevSibling: c.PrevSibling,
					NextSibling: n,
				}
			}

			head := processors.NewHeadProcessor()
			t.buildProcessorChain(head, currentNode)
			tagProcessor.AddChild(head)
		}
	}
}

func (t *GoatsTemplate) getAttrMap(node *html.Node) map[string]string {
	attrs := map[string]string{}
	for _, attr := range node.Attr {
		if _, ok := multipleAttrs[attr.Key]; ok {
			if existing, ok := attrs[attr.Key]; ok {
				if existing != "" {
					existing += ";"
				}
				existing += attr.Val
				attrs[attr.Key] = existing
			} else {
				attrs[attr.Key] = attr.Val
			}
		} else {
			attrs[attr.Key] = attr.Val
		}
	}
	return attrs
}
