package goats

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	txttpl "text/template"
	"time"
)

const (
	// Enum for tag processing types.
	TagProcessingGoRegular = iota
	TagProcessingGoSwitch
)

const TemplateInterfaceFile = ("package {{.PkgName}}\n" +
	"\n" +
	"import (\n" +
	"  \"goats-html/goats/runtime\"\n" +
	"  @@IMPORT@@\n" +
	")\n" +
	"\n" +
	"type {{.Name}}TemplateArgs struct {\n" +
	"{{range .Args}} {{.Declare}}\n{{end}}" +
	"}\n" +
	"\n" +
	"{{$tmplName := .Name}}" +
	"{{range .Replaceables}}" +
	"type {{$tmplName}}{{.Name}}ReplArgs struct {\n" +
	"  {{range .Args}} {{.Declare}}\n{{end}}" +
	"}\n" +
	"type {{$tmplName}}{{.Name}}ReplFunc func(*{{$tmplName}}{{.Name}}ReplArgs)\n\n" +
	"{{end}}" +
	"type {{.Name}}Template interface {\n" +
	"  runtime.Template\n" +
	"  Render(*{{.Name}}TemplateArgs) error\n" +
	"{{range .Replaceables}}" +
	"  Replace{{.Name}}({{$tmplName}}{{.Name}}ReplFunc)\n" +
	"{{end}}" +
	"}")

const TemplateImplFile = ("// +build !goats_devmod\n" +
	"\n" +
	"package {{.PkgName}}\n" +
	"\n" +
	"import (\n" +
	"  \"goats-html/goats/runtime\"\n" +
	"  \"io\"\n" +
	"  @@IMPORT@@\n" +
	")\n" +
	"\n" +
	"{{$tmplName := .Name}}" +
	"type {{$tmplName}}TemplateImpl struct {\n" +
	"  *runtime.BaseTemplate\n" +
	"  builtinFilter *runtime.BuiltinFilter\n" +
	"{{range .Replaceables}}" +
	"  {{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc\n" +
	"{{end}}" +
	"}\n" +
	"\n" +
	"func (__impl *{{$tmplName}}TemplateImpl) Render(__args *{{.Name}}TemplateArgs) error {\n" +
	"  @@RENDER@@\n" +
	"  return nil\n" +
	"}\n" +
	"\n" +
	"{{range .Replaceables}}" +
	"  func (__impl *{{$tmplName}}TemplateImpl) Replace{{.Name}}({{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc) {\n" +
	"    __impl.{{.HiddenName}} = {{.HiddenName}}\n" +
	"}\n\n" +
	"{{end}}" +
	"\n" +
	"func New{{.Name}}Template(writer io.Writer, settings *runtime.TemplateSettings) {{.Name}}Template {\n" +
	"  template := &{{.Name}}TemplateImpl{}\n" +
	"  template.BaseTemplate = runtime.NewBaseTemplate(writer, settings)\n" +
	"  template.builtinFilter = runtime.NewBuiltinFilter()\n" +
	"  return template\n" +
	"}\n")

const TemplateProxyFile = ("// +build goats_devmod\n" +
	"\n" +
	"package {{.PkgName}}\n" +
	"\n" +
	"import (\n" +
	"  \"goats-html/goats/runtime\"\n" +
	"  \"io\"\n" +
	")\n" +
	"\n" +
	"{{$tmplName := .Name}}" +
	"type {{.HiddenName}}TemplateProxy struct {\n" +
	"  *runtime.BaseTemplate\n" +
	"{{range .Replaceables}}  {{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc\n{{end}}" +
	"}\n" +
	"\n" +
	"func (__proxy *{{.HiddenName}}TemplateProxy) Render(args *{{.Name}}TemplateArgs) error {\n" +
	"  err := runtime.CallRpc(\"{{.Pkg}}\",\n" +
	"    \"{{.Name}}\",\n" +
	"    __proxy.GetSettings(),\n" +
	"    args,\n" +
	"    __proxy.GetWriter())\n" +
	"  return err\n" +
	"}\n" +
	"\n" +
	"{{$name := .HiddenName}}" +
	"\n" +
	"{{range .Replaceables}}  func (__impl *{{$name}}TemplateProxy) Replace{{.Name}}({{.HiddenName}} {{$tmplName}}{{.Name}}ReplFunc) {\n" +
	"  __impl.{{.HiddenName}} = {{.HiddenName}}\n" +
	"}\n\n" +
	"{{end}}" +
	"\n" +
	"func New{{.Name}}Template(writer io.Writer, settings *runtime.TemplateSettings) {{.Name}}Template {\n" +
	"  template := &{{.HiddenName}}TemplateProxy{}\n" +
	"  template.BaseTemplate = runtime.NewBaseTemplate(writer, settings)\n" +
	"  return template\n" +
	"}")

const TemplateMainFile = ("package main\n\n" +
	"import(\n" +
	"  \"bytes\"\n" +
	"  \"{{.Pkg}}\"\n" +
	"  \"goats-html/goats/runtime\"\n" +
	"  \"os\"\n" +
	")\n\n" +
	"func main() {\n" +
	"  settings := runtime.TemplateSettings{}\n" +
	"  var buffer bytes.Buffer\n" +
	"  switch os.Args[1] {\n" +
	"{{range .Templates}}  case \"{{.Name}}\":" +
	"    args := {{.PkgName}}.{{.Name}}TemplateArgs{}\n" +
	"    runtime.DecodeRpcRequestOrFail(os.Stdin, &settings, &args)\n" +
	"    template := {{.PkgName}}.New{{.Name}}Template(&buffer, &settings)\n" +
	"    template.Render(&args)\n{{end}}" +
	"  default:\n" +
	"    panic(\"Unknown template name: \" + os.Args[1])\n" +
	"  }\n" +
	"  os.Stdout.Write(buffer.Bytes())\n" +
	"}\n")

const (
	ImplFileSuffix  = "_impl.go"
	ProxyFileSuffix = "_proxy.go"
)

var (
	MainFileName = path.Join("cmd", "main.go")
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

var argMatcher *regexp.Regexp = regexp.MustCompile(
	`^(?P<name>\w+)\s*:\s*(?P<type>([*]|\w|\.)+)(\s*=\s*(?P<value>(\w|\.)*))?$`)

type ParserSettings struct {
	PkgRoot      string
	TemplateDir  string
	OutputDir    string
	Clean        bool
	KeepComments bool
	SampleData   bool
}

type PkgImport struct {
	Name  string
	Alias string
	Path  string
}

func formatSource(writer io.Writer, unformated string) {
	formated, err := format.Source([]byte(unformated))
	if err != nil {
		io.WriteString(writer, unformated)
		log.Fatal("Failed to format the output template, ", err)
	}
	io.WriteString(writer, string(formated))
}

func (pi *PkgImport) GenerateImports(buffer *bytes.Buffer) {
	if pi.Alias != "" {
		buffer.WriteString(fmt.Sprintf("%s \"%s\"\n", pi.Alias, pi.Path))
	} else {
		buffer.WriteString(fmt.Sprintf("\"%s\"\n", pi.Path))
	}
}

func NewPkgImport(impt string) *PkgImport {
	var pkgAlias, pkgPath string
	if strings.Contains(impt, ":") {
		parts := strings.Split(impt, ":")
		pkgAlias = TrimWhiteSpaces(parts[0])
		pkgPath = TrimWhiteSpaces(parts[1])
	} else if strings.Contains(impt, "/") {
		pkgAlias = TrimWhiteSpaces(path.Base(impt))
		pkgPath = TrimWhiteSpaces(impt)
	} else {
		pkgAlias = TrimWhiteSpaces(impt)
		pkgPath = TrimWhiteSpaces(impt)
	}
	return &PkgImport{
		Name:  TrimWhiteSpaces(path.Base(impt)),
		Alias: pkgAlias,
		Path:  pkgPath,
	}
}

func NewPkgImportFromCall(outputDir, callStmt string) *PkgImport {
	callStmt = TrimWhiteSpaces(callStmt)
	if !strings.HasPrefix(callStmt, "#") {
		pkgPath := path.Join(
			outputDir, strings.Replace(strings.Split(callStmt, "#")[0], ".html", "_html", -1))
		pkgName := path.Base(pkgPath)
		return &PkgImport{
			Name:  pkgName,
			Alias: pkgName,
			Path:  pkgPath,
		}
	}
	return nil
}

type GoatsReplaceable struct {
	Name       string
	HiddenName string
	Args       []*Argument
}

type GoatsReplace struct {
	Name       string
	HiddenName string
	Args       []*Argument
}

type GoatsTemplate struct {
	Parser           *GoatsParser
	OutputPath       string
	OutputIfaceFile  string
	OutputImplFile   string
	OutputProxyFile  string
	Pkg              string
	PkgName          string
	Name             string
	HiddenName       string
	Args             []*Argument
	RootNode         *html.Node
	NeedsDocType     bool
	Replaceables     []*GoatsReplaceable
	Replaces         []*GoatsReplace
	ImportsInterface map[string]*PkgImport // imports for interface file
	Imports          map[string]*PkgImport // imports for non interface files
}

func NewGoatsTemplate(parser *GoatsParser, tmplName string, args []*Argument,
	rootNode *html.Node, needsDocType bool, importsIface map[string]*PkgImport) *GoatsTemplate {
	prefix := ToSnakeCase(tmplName)
	return &GoatsTemplate{
		Parser:           parser,
		OutputPath:       parser.OutputPath,
		OutputIfaceFile:  fmt.Sprintf("%s.go", prefix),
		OutputImplFile:   fmt.Sprintf("%s%s", prefix, ImplFileSuffix),
		OutputProxyFile:  fmt.Sprintf("%s%s", prefix, ProxyFileSuffix),
		Pkg:              parser.Pkg,
		PkgName:          path.Base(parser.Pkg),
		Name:             tmplName,
		HiddenName:       ToHiddenName(tmplName),
		Args:             args,
		RootNode:         rootNode,
		NeedsDocType:     needsDocType,
		ImportsInterface: importsIface,
		Imports:          map[string]*PkgImport{},
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

func (t *GoatsTemplate) generateInterface() {
	goFilePath := path.Join(t.OutputPath, t.OutputIfaceFile)
	goFile, err := os.Create(goFilePath)
	if err != nil {
		log.Fatal("Failed to create file " + goFilePath)
	}
	defer goFile.Close()

	var buffer bytes.Buffer
	tmpl, err := txttpl.New("interface").Parse(TemplateInterfaceFile)
	if err != nil {
		log.Fatal("Failed to generate file "+goFilePath, err)
	}
	err = tmpl.Execute(&buffer, t)
	if err != nil {
		log.Fatal("Failed to generate file "+goFilePath, err)
	}

	// Generate imports
	var importsBuffer bytes.Buffer
	for _, pkgImport := range t.ImportsInterface {
		pkgImport.GenerateImports(&importsBuffer)
	}
	text := strings.Replace(buffer.String(), "@@IMPORT@@", importsBuffer.String(), 1)
	formatSource(goFile, text)
}

func (t *GoatsTemplate) generateImpl() {
	goFilePath := path.Join(t.OutputPath, t.OutputImplFile)
	goFile, err := os.Create(goFilePath)
	if err != nil {
		log.Fatal("Failed to create file " + goFilePath)
	}
	defer goFile.Close()

	var buffer bytes.Buffer
	tmpl, err := txttpl.New("impl").Parse(TemplateImplFile)
	if err != nil {
		log.Fatal("Failed to generate file "+goFilePath, err)
	}
	err = tmpl.Execute(&buffer, t)
	if err != nil {
		log.Fatal("Failed to generate file "+goFilePath, err)
	}

	// Generate render content
	var headProcessor Processor = NewArgProcessor(t.Args)
	t.buildProcessorChain(headProcessor, t.RootNode)
	context := NewTagContext()
	if t.NeedsDocType {
		docTypeProcessor := NewDocTypeProcessor(t.Parser.DocTypeTag, t.Parser.DocTypeAttrs)
		docTypeProcessor.SetNext(headProcessor)
		headProcessor = docTypeProcessor
	}
	var renderBuffer bytes.Buffer
	headProcessor.Process(&renderBuffer, context)

	// manage imports
	var importsBuffer bytes.Buffer
	for _, pkgImport := range t.Imports {
		pkgImport.GenerateImports(&importsBuffer)
	}
	for impt, _ := range context.GetImports() {
		if pkgImport, ok := t.Parser.Imports[impt]; ok {
			pkgImport.GenerateImports(&importsBuffer)
		}
	}
	text := strings.Replace(buffer.String(), "@@IMPORT@@", importsBuffer.String(), 1)
	unformated := strings.Replace(text, "@@RENDER@@", renderBuffer.String(), 1)

	formatSource(goFile, unformated)
}

func (t *GoatsTemplate) generateProxy() {
	goFilePath := path.Join(t.OutputPath, t.OutputProxyFile)
	goFile, err := os.Create(goFilePath)
	if err != nil {
		log.Fatal("Failed to create file " + goFilePath)
	}
	defer goFile.Close()

	var buffer bytes.Buffer
	tmpl, err := txttpl.New("proxy").Parse(TemplateProxyFile)
	if err != nil {
		log.Fatal("Failed to generate file "+goFilePath, err)
	}
	err = tmpl.Execute(&buffer, t)
	if err != nil {
		log.Fatal("Failed to generate file ", goFilePath, err)
	}
	formatSource(goFile, buffer.String())
}

func (t *GoatsTemplate) findTemplateCall(node *html.Node) {
	// TODO: Cyclic template call detection.
	for _, attr := range node.Attr {
		if attr.Key == "go:call" {
			if pkgImport := NewPkgImportFromCall(t.Parser.Settings.OutputDir, attr.Val); pkgImport != nil {
				t.Imports[pkgImport.Alias] = pkgImport
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		t.findTemplateCall(c)
	}
}

func (t *GoatsTemplate) buildProcessorChain(preProcessor Processor, node *html.Node) {
	if node.Type == html.CommentNode {
		if t.Parser.Settings.KeepComments {
			processor := NewCommentProcessor(node.Data)
			preProcessor.SetNext(processor)
			preProcessor = processor
		}
	} else if node.Type == html.TextNode {
		processor := NewTextProcessor(node.Data)
		preProcessor.SetNext(processor)
		preProcessor = processor
	} else if node.Type == html.ElementNode {
		goAttrs := t.getAttrMap(node)

		if val, ok := goAttrs["go:var"]; ok {
			varProcessor := NewVarsProcessor(val)
			preProcessor.SetNext(varProcessor)
			preProcessor = varProcessor
		}

		if val, ok := goAttrs["go:if"]; ok {
			ifProcessor := NewIfProcessor(val)
			preProcessor.SetNext(ifProcessor)
			preProcessor = ifProcessor
		}

		if val, ok := goAttrs["go:for"]; ok {
			forProcessor := NewForProcessor(val)
			preProcessor.SetNext(forProcessor)
			preProcessor = forProcessor
		}

		if val, ok := goAttrs["go:case"]; ok {
			caseProcessor := NewCaseProcessor(val)
			preProcessor.SetNext(caseProcessor)
			preProcessor = caseProcessor
		} else if _, ok := goAttrs["go:default"]; ok {
			defaultProcessor := NewDefaultProcessor()
			preProcessor.SetNext(defaultProcessor)
			preProcessor = defaultProcessor
		}

		if _, ok := goAttrs["go:switch"]; ok {
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
			callProcessor := NewCallProcessor("", val, ParseArgDefs(goAttrs["go:arg"]), nil, node.Attr)
			preProcessor.SetNext(callProcessor)
			preProcessor = callProcessor
			return
		}

		if val, ok := goAttrs["go:replaceable"]; ok && node != t.RootNode {
			replaceableProcessor := NewReplaceableProcessor(t.Name, val, ParseArgDefs(goAttrs["go:arg"]))
			preProcessor.SetNext(replaceableProcessor)
			preProcessor = replaceableProcessor
		}

		if val, ok := goAttrs["go:call"]; ok {
			if !strings.Contains(val, "#") {
				log.Fatal("Call to template must contain a \"#\".")
			}
			parts := strings.Split(val, "#")

			var replacements []*Replacement
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					if len(TrimWhiteSpaces(c.Data)) != 0 {
						log.Fatal("Node with go:call can only contain nodes with go:replace or spaces.")
					}
					continue
				}

				var found bool
				replacement := &Replacement{
					Args: []*Argument{},
				}
				for _, attr := range c.Attr {
					if attr.Key == "go:replace" {
						found = true
						head := NewHeadProcessor()
						t.buildProcessorChain(head, c)
						replacement.Name = attr.Val
						replacement.Head = head
						replacements = append(replacements, replacement)
					} else if attr.Key == "go:arg" {
						replacement.Args = append(replacement.Args, NewArgDef(attr.Val))
					}
				}
				if !found {
					log.Fatal("Node with go:call can only contain nodes with go:replace.")
				}
			}

			callProcessor := NewCallProcessor(
				parts[0], parts[1], ParseArgCalls(goAttrs["go:arg"]), replacements, node.Attr)
			preProcessor.SetNext(callProcessor)
			preProcessor = callProcessor

			return
		}
		t.handleTag(preProcessor, node, TagProcessingGoRegular)
	}
}

func (t *GoatsTemplate) handleTag(preProcessor Processor, node *html.Node, tagProcessingType int) {
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
	tagProcessor := NewTagProcessor(node.Data, omitTag, firstTag, !voidElements[node.Data], node.Attr)
	preProcessor.SetNext(tagProcessor)
	preProcessor = tagProcessor

	if tagProcessingType == TagProcessingGoSwitch {
		switchProcessor := NewSwitchProcessor(goAttrs["go:switch"])
		preProcessor.SetNext(switchProcessor)
		preProcessor = switchProcessor
	}

	if val, ok := goAttrs["go:content"]; ok {
		contentProcessor := NewContentProcessor(val)
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

			head := NewHeadProcessor()
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

type GoatsParser struct {
	Settings     *ParserSettings
	ModTime      time.Time
	Pkg          string
	OutputPath   string
	HtmlFilePath string
	Doc          *html.Node
	DocTypeTag   string
	DocTypeAttrs []html.Attribute
	Templates    map[string]*GoatsTemplate
	Imports      map[string]*PkgImport
}

func (p *GoatsParser) loadFile() {
	reader, err := os.Open(p.HtmlFilePath)
	if err != nil {
		log.Fatal("Failed to open file " + p.HtmlFilePath)
	}
	defer reader.Close()

	p.Doc, err = html.Parse(reader)
	if err != nil {
		log.Fatal("Failed to open file " + p.HtmlFilePath)
	}

	p.FindTemplates(p.Doc)

	// For each template search its go:call and add corresponding template.
	for _, tmpl := range p.Templates {
		tmpl.findTemplateCall(tmpl.RootNode)
	}
}

func (p *GoatsParser) FindTemplates(node *html.Node) {
	if node.Type == html.DoctypeNode {
		p.DocTypeTag = node.Data
		p.DocTypeAttrs = node.Attr
	} else if node.Type == html.ElementNode {
		// Collect imports.
		var pkgImport *PkgImport = nil
		for _, attr := range node.Attr {
			if attr.Key == "go:import" && node.Data == "html" {
				pkgImport = NewPkgImport(attr.Val)
			} else if attr.Key == "go:call" {
				pkgImport = NewPkgImportFromCall(p.Settings.OutputDir, attr.Val)
			}
			if pkgImport != nil {
				p.Imports[pkgImport.Alias] = pkgImport
			}
		}

		templateName := ""
		needsDocType := false
		args := []*Argument{}
		imports := map[string]*PkgImport{}
		for _, attr := range node.Attr {
			if attr.Key == "go:template" {
				templateName = attr.Val
				if node.Data == "html" {
					needsDocType = true
				}
			} else if attr.Key == "go:arg" {
				arg := NewArgDef(attr.Val)
				args = append(args, arg)
				if pkgImport, ok := p.Imports[arg.PkgName]; ok {
					imports[arg.PkgName] = pkgImport
				}
			}
		}
		if templateName != "" {
			template := NewGoatsTemplate(p, templateName, args, node, needsDocType, imports)
			p.Templates[templateName] = template
			p.findReplaceables(node, template)
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				for _, attr := range c.Attr {
					if attr.Key == "go:replace" {
						template.Replaces = append(template.Replaces,
							&GoatsReplace{
								Name:       attr.Val,
								HiddenName: ToHiddenName(attr.Val),
							})
						break
					}
				}
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		p.FindTemplates(c)
	}
}

func (p *GoatsParser) findReplaceables(node *html.Node, template *GoatsTemplate) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		var foundReplaceable bool
		var foundTemplate bool
		replaceable := &GoatsReplaceable{
			Args: []*Argument{},
		}
		for _, attr := range c.Attr {
			if attr.Key == "go:template" {
				if foundTemplate {
					log.Fatal("Found multiple go:template on the same node.")
				}
				foundTemplate = true
			} else if attr.Key == "go:replaceable" {
				if foundReplaceable {
					log.Fatal("Found multiple go:replaceable on the same node.")
				}
				foundReplaceable = true
				replaceable.Name = attr.Val
				replaceable.HiddenName = ToHiddenName(attr.Val)
				template.Replaceables = append(template.Replaceables, replaceable)
			} else if attr.Key == "go:arg" {
				replaceable.Args = append(replaceable.Args, NewArgDef(attr.Val))
			}
		}
		if foundTemplate && foundReplaceable {
			log.Fatal("go:template can not be on the same node which has go:replaceable.")
		}
		if !foundTemplate {
			p.findReplaceables(c, template)
		}
	}
}

func (p *GoatsParser) IsDirty() bool {
	_, err := os.Stat(p.OutputPath)
	if os.IsNotExist(err) {
		return true
	}

	if p.IsFileOld(MainFileName) {
		return true
	}

	for _, tmpl := range p.Templates {
		if tmpl.IsDirty() {
			return true
		}
	}

	return false
}

func (p *GoatsParser) IsFileOld(fileName string) bool {
	// Generated interface file.
	info, err := os.Stat(path.Join(p.OutputPath, fileName))
	if os.IsNotExist(err) {
		return true
	}

	return info.ModTime().Before(p.ModTime)
}

func (p *GoatsParser) Generate() {
	cmdPath := path.Join(p.OutputPath, "cmd")
	_, err := os.Stat(cmdPath)
	if os.IsNotExist(err) {
		os.MkdirAll(cmdPath, os.ModePerm)
	}

	for name, template := range p.Templates {
		fmt.Printf("    Generating template \"%s\":\n", name)
		fmt.Printf("        %s\n", template.OutputIfaceFile)
		template.generateInterface()
		fmt.Printf("        %s\n", template.OutputImplFile)
		template.generateImpl()
		fmt.Printf("        %s\n", template.OutputProxyFile)
		template.generateProxy()
	}
	fmt.Println("    Generating main file " + MainFileName)
	p.generateMain()
}

func (p *GoatsParser) generateMain() {
	goFilePath := path.Join(p.OutputPath, MainFileName)
	goFile, err := os.Create(goFilePath)
	if err != nil {
		log.Fatal("Failed to create file " + goFilePath)
	}
	defer goFile.Close()

	var buffer bytes.Buffer
	tmpl, err := txttpl.New("main").Parse(TemplateMainFile)
	if err != nil {
		log.Fatal("Failed to parse main template\n", err)
	}
	err = tmpl.Execute(&buffer, p)
	if err != nil {
		log.Fatal("Failed to generate file ", goFilePath, err)
	}
	formatSource(goFile, buffer.String())
}

func NewParser(parserSettings *ParserSettings, htmlFilePath string) *GoatsParser {
	info, err := os.Stat(htmlFilePath)
	if os.IsNotExist(err) {
		panic("Can not access template file " + htmlFilePath)
	}

	htmlFileName := path.Base(htmlFilePath)

	pkgAlias := strings.Replace(htmlFileName, ".", "_", -1)

	outputPath, err := filepath.Abs(
		path.Join(parserSettings.PkgRoot, parserSettings.OutputDir, pkgAlias))
	if err != nil {
		log.Fatal("Invalid output path: ", outputPath)
	}

	pkg := path.Join(parserSettings.OutputDir, pkgAlias)
	p := &GoatsParser{
		Settings:     parserSettings,
		ModTime:      info.ModTime(),
		Pkg:          pkg,
		OutputPath:   outputPath,
		HtmlFilePath: htmlFilePath,
		Templates:    map[string]*GoatsTemplate{},
		Imports:      map[string]*PkgImport{},
	}
	p.loadFile()
	return p
}
