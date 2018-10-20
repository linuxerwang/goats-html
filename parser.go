package goats

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	txttpl "text/template"
	"time"

	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/processors"
	"github.com/linuxerwang/goats-html/util"
	"golang.org/x/net/html"
)

type ParserSettings struct {
	PkgRoot         string
	TemplateDir     string
	OutputDir       string
	OutputFormat    string
	OutputPkgPrefix string
	OutputExport    bool
	Clean           bool
	KeepComments    bool
	SampleData      bool
	GenMergedFile   bool
}

// The goats template parser. There is one parser per template file.
type GoatsParser struct {
	Settings     *ParserSettings
	ModTime      time.Time
	Pkg          string
	OutputPath   string
	HtmlFilePath string
	RelativePath string
	Doc          *html.Node
	DocTypeTag   string
	DocTypeAttrs []html.Attribute
	Templates    map[string]*GoatsTemplate
	PkgMgr       *pkgmgr.PkgManager
	PkgRefs      *pkgmgr.PkgRefs
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
}

func (p *GoatsParser) FindTemplates(node *html.Node) {
	if node.Type == html.DoctypeNode {
		p.DocTypeTag = node.Data
		p.DocTypeAttrs = node.Attr
	} else if node.Type == html.ElementNode {
		// Collect imports.
		for _, attr := range node.Attr {
			if attr.Key == "go:import" && node.Data == "html" {
				impt := util.TrimWhiteSpaces(attr.Val)
				if strings.Contains(impt, ":") {
					parts := strings.Split(impt, ":")
					pbPkg := ""
					if i := strings.Index(parts[1], "[pb]"); i > -1 {
						pbPkg = parts[1][i+4:]
						parts[1] = parts[1][:i]
					}
					p.PkgMgr.AddImport(util.TrimWhiteSpaces(parts[0]), util.TrimWhiteSpaces(parts[1]), util.TrimWhiteSpaces(pbPkg))
				} else {
					pbPkg := ""
					if i := strings.Index(impt, "[pb]"); i > -1 {
						pbPkg = impt[i+4:]
						impt = impt[:i]
					}
					p.PkgMgr.AddImport(util.TrimWhiteSpaces(path.Base(impt)), util.TrimWhiteSpaces(impt), util.TrimWhiteSpaces(pbPkg))
				}
			} else if attr.Key == "go:call" {
				if node.Data == "html" {
					panic("Attr go:call is not allowed on <html> tag.")
				}
				pkgPath, _, _ := p.PkgMgr.ParseTmplCall(attr.Val)
				// TODO: handle pb variation.
				p.PkgMgr.AddImport("", pkgPath, "")
				break
			}
		}

		templateName := ""
		needsDocType := false
		args := []*processors.Argument{}

		var pkgRefs *pkgmgr.PkgRefs
		if p.Settings.GenMergedFile {
			pkgRefs = p.PkgRefs
		} else {
			pkgRefs = p.PkgMgr.CreatePkgRefs()
		}

		for _, attr := range node.Attr {
			if attr.Key == "go:template" {
				templateName = attr.Val
				if node.Data == "html" {
					needsDocType = true
				}
			}
		}
		if templateName != "" {
			for _, attr := range node.Attr {
				if attr.Key == "go:arg" {
					arg := processors.NewArgDef(attr.Val)
					args = append(args, arg)
					if arg.PkgName != "" {
						pkgRefs.RefByAlias(arg.PkgName, true)
					}
				}
			}

			template := NewGoatsTemplate(p, templateName, args, node, needsDocType, pkgRefs)
			if p.Settings.OutputFormat == "closure" {
				prefix := p.Settings.OutputPkgPrefix
				pkgName := ""
				for i, part := range strings.Split(prefix, ".") {
					if i == 0 {
						pkgName = part
					} else {
						pkgName += "." + part
					}
				}
				pkgName += "." + template.PkgName

				template.ClosurePkgName = pkgName
			}

			p.Templates[templateName] = template
			p.findReplaceables(node, template)
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				for _, attr := range c.Attr {
					if attr.Key == "go:replace" {
						template.Replaces = append(template.Replaces,
							&GoatsReplace{
								Name:       attr.Val,
								HiddenName: util.ToHiddenName(attr.Val),
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
			Args: []*processors.Argument{},
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
				replaceable.Name = strings.Title(attr.Val)
				replaceable.HiddenName = util.ToHiddenName(attr.Val)
				template.Replaceables = append(template.Replaceables, replaceable)
			}
		}

		if foundTemplate && foundReplaceable {
			log.Fatal("go:template can not be on the same node which has go:replaceable.")
		}

		if !foundTemplate {
			p.findReplaceables(c, template)
		}

		if foundReplaceable {
			for _, attr := range c.Attr {
				if attr.Key == "go:arg" {
					replaceable.Args = append(replaceable.Args, processors.NewArgDef(attr.Val))
				}
			}
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
	info, err := os.Stat(filepath.Join(p.OutputPath, fileName))
	if os.IsNotExist(err) {
		return true
	}

	return info.ModTime().Before(p.ModTime)
}

func (p *GoatsParser) Generate() {
	cmdPath := filepath.Join(p.OutputPath, "cmd")
	_, err := os.Stat(cmdPath)
	if os.IsNotExist(err) {
		os.MkdirAll(cmdPath, os.ModePerm)
	}

	switch p.Settings.OutputFormat {
	case "go":
		p.genGoSource()
	case "closure":
		p.genClosureSource()
	}
}

func (p *GoatsParser) genGoSource() {
	if p.Settings.GenMergedFile {
		p.genMergedFile()
	} else {
		p.genMultiGoFile()
	}
}

func (p *GoatsParser) genMergedFile() {
	p.genMergedIfaceFile()
	p.genMergedImplFile()
	p.genMergedProxyFile()
	p.genMainFile()
}

func (p *GoatsParser) genMergedIfaceFile() {
	fmt.Printf("    Generating merged Go interface file \"interfaces.go\":\n")

	var bufBody bytes.Buffer
	var bufOther bytes.Buffer

	isFirst := true
	for _, t := range p.Templates {
		if isFirst {
			t.genIfacePkgDecl(&bufOther)
			io.WriteString(&bufOther, "import (\n")
			t.genIfaceImports(&bufOther)
			io.WriteString(&bufOther, ")\n\n")
			isFirst = false
		}
		t.genIfaceBody(&bufBody)
	}

	p.genFile(p.OutputPath, "interfaces.go", func(output io.Writer) {
		io.WriteString(output, formatSource(bufOther.String()+bufBody.String()))
	})
}

func (p *GoatsParser) genMergedImplFile() {
	fmt.Printf("    Generating merged Go implementation file \"implementations.go\":\n")

	var bufBody bytes.Buffer
	var bufOther bytes.Buffer

	isFirst := true
	for _, t := range p.Templates {
		// Gen impl body first to collect imports.
		t.genImplBody(&bufBody)

		if isFirst {
			t.genImplPkgDecl(&bufOther)
			io.WriteString(&bufOther, "import (\n")
			t.genImplImports(&bufOther)
			io.WriteString(&bufOther, ")\n\n")
			isFirst = false
		}
	}

	p.genFile(p.OutputPath, "implementations.go", func(output io.Writer) {
		io.WriteString(output, formatSource(bufOther.String()+bufBody.String()))
	})
}

func (p *GoatsParser) genMergedProxyFile() {
	fmt.Printf("    Generating merged Go proxy file \"proxies.go\":\n")

	var bufBody bytes.Buffer
	var bufOther bytes.Buffer

	isFirst := true
	for _, t := range p.Templates {
		if isFirst {
			t.genProxyPkgDecl(&bufOther)
			io.WriteString(&bufOther, "import (\n")
			t.genProxyImports(&bufOther)
			io.WriteString(&bufOther, ")\n\n")
			isFirst = false
		}
		t.genProxyBody(&bufOther)
	}

	p.genFile(p.OutputPath, "proxies.go", func(output io.Writer) {
		io.WriteString(output, formatSource(bufOther.String()+bufBody.String()))
	})
}

func (p *GoatsParser) genMainFile() {
	p.generateMain()
}

func (p *GoatsParser) genMultiGoFile() {
	for name, t := range p.Templates {
		fmt.Printf("    Generating template \"%s\":\n", name)

		var bufBody bytes.Buffer
		var bufOther bytes.Buffer

		// Gen iface body first to collect imports.
		t.genIfaceBody(&bufBody)

		t.genIfacePkgDecl(&bufOther)
		io.WriteString(&bufOther, "import (\n")
		t.genIfaceImports(&bufOther)
		io.WriteString(&bufOther, ")\n\n")

		p.genFile(p.OutputPath, t.OutputIfaceFile, func(output io.Writer) {
			fmt.Printf("        %s\n", t.OutputIfaceFile)
			io.WriteString(output, formatSource(bufOther.String()+bufBody.String()))
		})

		bufBody.Reset()
		bufOther.Reset()

		// Gen impl body first to collect imports.
		t.genImplBody(&bufBody)

		t.genImplPkgDecl(&bufOther)
		io.WriteString(&bufOther, "import (\n")
		t.genImplImports(&bufOther)
		io.WriteString(&bufOther, ")\n\n")

		p.genFile(p.OutputPath, t.OutputImplFile, func(output io.Writer) {
			fmt.Printf("        %s\n", t.OutputImplFile)
			io.WriteString(output, formatSource(bufOther.String()+bufBody.String()))
		})

		bufBody.Reset()
		bufOther.Reset()

		// Gen proxy body first to collect imports.
		t.genProxyBody(&bufBody)

		t.genProxyPkgDecl(&bufOther)
		io.WriteString(&bufOther, "import (\n")
		t.genProxyImports(&bufOther)
		io.WriteString(&bufOther, ")\n\n")

		p.genFile(p.OutputPath, t.OutputProxyFile, func(output io.Writer) {
			fmt.Printf("        %s\n", t.OutputProxyFile)
			io.WriteString(output, formatSource(bufOther.String()+bufBody.String()))
		})
	}

	// Generate main.
	fmt.Println("    Generating main file " + MainFileName)
	p.generateMain()
}

func (p *GoatsParser) genClosureSource() {
	if p.Settings.GenMergedFile {
		p.genMergedClosureFile()
	} else {
		p.genMultiClosureFile()
	}
}

func (p *GoatsParser) genMergedClosureFile() {
	p.genFile(p.OutputPath, "closure-all.js", func(output io.Writer) {
		fmt.Printf("    Generating template \"closure-all.js\":\n")

		// Sort: guarantee output is reproducible (same checksum for same source code).
		keys := make([]string, 0, len(p.Templates))
		for key := range p.Templates {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		var buf bytes.Buffer
		for _, key := range keys {
			t := p.Templates[key]
			t.genClosureBody(&buf)
		}

		isFirst := true
		for _, key := range keys {
			t := p.Templates[key]
			if isFirst {
				t.genClosurePkgDoc(output)
				isFirst = false
			}
		}

		for _, key := range keys {
			t := p.Templates[key]
			t.genClosureProvides(output)
		}

		requires := make(map[string]bool)
		for _, t := range p.Templates {
			t.dumpClosureCommonRequires(requires)
			t.dumpClosureRequires(requires)
		}

		keys = make([]string, 0, len(requires))
		for key := range requires {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			io.WriteString(output, fmt.Sprintf("goog.require('%s');\n", key))
		}

		io.WriteString(output, buf.String())
	})
}

func (p *GoatsParser) genMultiClosureFile() {
	for name, t := range p.Templates {
		p.genFile(p.OutputPath, t.OutputImplFile, func(output io.Writer) {
			fmt.Printf("    Generating template \"%s\":\n", name)
			fmt.Printf("        %s\n", t.OutputImplFile)

			var buf bytes.Buffer
			t.genClosureBody(&buf)

			t.genClosurePkgDoc(output)
			t.genClosureProvides(output)
			t.genClosureCommonRequires(output)
			t.genClosureRequires(output)

			io.WriteString(output, buf.String())
		})
	}
}

func (p *GoatsParser) genFile(dir string, fn string, callback func(output io.Writer)) {
	filePath := filepath.Join(dir, fn)
	output, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Failed to create file " + filePath)
	}
	defer output.Close()

	callback(output)
}

func (p *GoatsParser) generateMain() {
	goFilePath := filepath.Join(p.OutputPath, MainFileName)
	goFile, err := os.Create(goFilePath)
	if err != nil {
		log.Fatal("Failed to create file " + goFilePath)
	}
	defer goFile.Close()

	var buffer bytes.Buffer
	tmpl, err := txttpl.New("main").Parse(tmplMainFile)
	if err != nil {
		log.Fatal("Failed to parse main template\n", err)
	}
	err = tmpl.Execute(&buffer, p)
	if err != nil {
		log.Fatal("Failed to generate file ", goFilePath, err)
	}
	source := formatSource(buffer.String())
	io.WriteString(goFile, source)
}

func NewParser(parserSettings *ParserSettings, htmlFilePath string) *GoatsParser {
	info, err := os.Stat(htmlFilePath)
	if os.IsNotExist(err) {
		panic("Can not access template file " + htmlFilePath)
	}

	tmplDir, err := filepath.Abs(filepath.Join(parserSettings.PkgRoot, parserSettings.TemplateDir))
	if err != nil {
		log.Fatal("Invalid template path: ", parserSettings.TemplateDir)
	}

	prefix, err := filepath.Rel(tmplDir, filepath.Dir(htmlFilePath))
	if err != nil {
		log.Fatalf("Can't make relative path \"%s\" vs. \"%s\".\n", filepath.Dir(htmlFilePath), tmplDir)
	}

	htmlFileName := filepath.Base(htmlFilePath)
	pkgName := strings.Replace(htmlFileName, ".", "_", -1)

	outputPath, err := filepath.Abs(
		filepath.Join(parserSettings.PkgRoot, parserSettings.OutputDir, prefix, pkgName))
	if err != nil {
		log.Fatal("Invalid output path: ", outputPath)
	}

	pkgMgr := pkgmgr.New(path.Join(parserSettings.OutputDir, prefix), parserSettings.OutputPkgPrefix)

	pkg := path.Join(parserSettings.OutputPkgPrefix, parserSettings.OutputDir, prefix, pkgName)
	p := &GoatsParser{
		Settings:     parserSettings,
		ModTime:      info.ModTime(),
		Pkg:          pkg,
		OutputPath:   outputPath,
		HtmlFilePath: htmlFilePath,
		RelativePath: prefix,
		Templates:    map[string]*GoatsTemplate{},
		PkgMgr:       pkgMgr,
	}
	if parserSettings.GenMergedFile {
		p.PkgRefs = p.PkgMgr.CreatePkgRefs()
	}
	p.loadFile()
	return p
}
