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
	UsageLine: "goats gen [--package_root <path>] [--template_dir <path>] [--output_dir <path>] [--clean] [--keep_comments] [--sample_data] [file1.html file2.html ...] [--output-format=<format>] [--output-pkg-prefix=<prefix>]",
	Short:     "Generate go files from goats templates.",
	Long:      "If specific goats html files are provided, generate go files for them; otherwise, generate go files for all goats html files under template_dir.",
}

var (
	genPkgRoot = genCmd.Flag.String(
		"package_root", ".", "go packages root directory containing templates.")
	genTemplateDir = genCmd.Flag.String(
		"template_dir", "", "Template directory relative to package root, required if no specific goats html file is given.")
	genOutputDir = genCmd.Flag.String(
		"output_dir", "", "Output directory relative to package root.")
	genOutputFormat = genCmd.Flag.String(
		"output-format", "go", "Output format, can be go (default) or closure.")
	genOutputPkgPrefix = genCmd.Flag.String(
		"output-pkg-prefix", "", "Output package prefix.")
	genOutputExport = genCmd.Flag.Bool(
		"output-export", true, "Output closure export, used for output format closure.")
	genClean        = genCmd.Flag.Bool("clean", true, "Clean unexisting *_html directories.")
	genKeepComments = genCmd.Flag.Bool("keep_comments", false, "Keep comment in output HTML.")
	genSampleData   = genCmd.Flag.Bool("sample_data", false, "Keep comment in output HTML.")

	parserSettings *goats.ParserSettings
)

func checkFlags() {
	if *genOutputFormat == "closure" && *genOutputPkgPrefix == "" {
		fmt.Println("Flag --output-pkg-prefix is required for output format closure.")
		os.Exit(2)
	}
}

func runGen(cmd *Command, args []string) {
	checkFlags()

	fmt.Printf("Package root directory: %s\n", *genPkgRoot)
	fmt.Printf("Templates directory: %s\n", *genTemplateDir)
	fmt.Printf("Output directory: %s\n\n", *genOutputDir)

	if *genOutputDir == "" {
		genOutputDir = genTemplateDir
	}

	parserSettings = &goats.ParserSettings{
		PkgRoot:         *genPkgRoot,
		TemplateDir:     *genTemplateDir,
		OutputDir:       *genOutputDir,
		OutputFormat:    *genOutputFormat,
		OutputPkgPrefix: *genOutputPkgPrefix,
		OutputExport:    *genOutputExport,
		Clean:           *genClean,
		KeepComments:    *genKeepComments,
		SampleData:      *genSampleData,
	}

	if len(args) > 0 {
		// Template files are specified.
		for _, templateFile := range args {
			tf, err := filepath.Abs(templateFile)
			if err == nil {
				parseTemplateFile(tf)
			} else {
				panic("Invalid template file: " + templateFile)
			}
		}
	} else {
		// All template files.
		tmplDir, err := filepath.Abs(filepath.Join(*genPkgRoot, *genTemplateDir))
		fmt.Println(tmplDir)
		if err == nil {
			filepath.Walk(tmplDir, walkFunc)
		} else {
			panic("Invalid template dir: " + tmplDir)
		}
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
