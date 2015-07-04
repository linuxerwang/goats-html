package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/linuxerwang/goats-html"
)

var genCmd = &Command{
	Name:      "gen",
	UsageLine: "goats gen [--package_root <path>] --template_dir <path> [--output_dir <path>] [--clean] [--keep_comments] [--sample_data]",
	Short:     "run gen on all template files under template root",
	Long:      "run gen on all template files under template root",
}

var (
	genPkgRoot = genCmd.Flag.String(
		"package_root", ".", "go packages root directory containing templates.")
	genTemplateDir = genCmd.Flag.String(
		"template_dir", "", "Template directory relative to package root.")
	genOutputDir = genCmd.Flag.String(
		"output_dir", "", "Output directory relative to package root.")
	genClean        = genCmd.Flag.Bool("clean", true, "Clean unexisting *_html directories.")
	genKeepComments = genCmd.Flag.Bool("keep_comments", false, "Keep comment in output HTML.")
	genSampleData   = genCmd.Flag.Bool("sample_data", false, "Keep comment in output HTML.")

	parserSettings *goats.ParserSettings
)

func runGen(cmd *Command, args []string) {
	if *genTemplateDir == "" {
		fmt.Fprintf(os.Stderr, "flag template_root is required.\n\n")
		os.Exit(2)
	} else if *genOutputDir == "" {
		genOutputDir = genTemplateDir
	}

	fmt.Fprintf(os.Stderr, "Package root directory: %s\n", *genPkgRoot)
	fmt.Fprintf(os.Stderr, "Templates directory: %s\n", *genTemplateDir)
	fmt.Fprintf(os.Stderr, "Output directory: %s\n\n", *genOutputDir)

	parserSettings = &goats.ParserSettings{
		PkgRoot:      *genPkgRoot,
		TemplateDir:  *genTemplateDir,
		OutputDir:    *genOutputDir,
		Clean:        *genClean,
		KeepComments: *genKeepComments,
		SampleData:   *genSampleData,
	}

	tmplDir, err := filepath.Abs(filepath.Join(*genPkgRoot, *genTemplateDir))
	if err == nil {
		filepath.Walk(tmplDir, walkFunc)
	} else {
		panic("Invalid template dir: " + tmplDir)
	}
}

func init() {
	genCmd.Run = runGen
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if info.IsDir() || !strings.HasSuffix(info.Name(), ".html") {
		return nil
	}

	parseTemplateFile(path)

	return nil
}

func parseTemplateFile(templateFile string) {
	parser := goats.NewParser(parserSettings, templateFile)
	fmt.Printf("Loading template file %s:", templateFile)
	if parser.IsDirty() {
		fmt.Println()
		parser.Generate()
	} else {
		fmt.Println(" not dirty")
	}
}
