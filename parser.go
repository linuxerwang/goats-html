package goats

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	txttpl "text/template"
	"time"

	"golang.org/x/net/html"
)

type ParserSettings struct {
	PkgRoot      string
	TemplateDir  string
	OutputDir    string
	Clean        bool
	KeepComments bool
	SampleData   bool
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
	PkgMgr       *PkgManager
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
				impt := TrimWhiteSpaces(attr.Val)
				if strings.Contains(impt, ":") {
					parts := strings.Split(impt, ":")
					p.PkgMgr.AddImport(TrimWhiteSpaces(parts[0]), TrimWhiteSpaces(parts[1]))
				} else {
					p.PkgMgr.AddImport(TrimWhiteSpaces(path.Base(impt)), TrimWhiteSpaces(impt))
				}
				break
			} else if attr.Key == "go:call" {
				pkgPath, _ := p.PkgMgr.ParseTmplCall(attr.Val)
				p.PkgMgr.AddImport("", pkgPath)
				break
			}
		}

		templateName := ""
		needsDocType := false
		args := []*Argument{}
		pkgRefs := p.PkgMgr.CreatePkgRefs()
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
					arg := NewArgDef(attr.Val)
					args = append(args, arg)
					if arg.PkgName != "" {
						pkgRefs.RefByAlias(arg.PkgName)
					}
				}
			}

			template := NewGoatsTemplate(p, templateName, args, node, needsDocType, pkgRefs)
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
					replaceable.Args = append(replaceable.Args, NewArgDef(attr.Val))
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
	goFilePath := filepath.Join(p.OutputPath, MainFileName)
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

	pkgMgr := NewPkgManager(path.Join(parserSettings.OutputDir, prefix))

	pkg := path.Join(parserSettings.OutputDir, prefix, pkgName)
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
	p.loadFile()
	return p
}
