package main

import (
	"bytes"
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var servCmd = &Command{
	Name:      "serv",
	UsageLine: "goats serv [--package_root <path>] --template_dir <path> [--output_dir <path>] [--keep_comments] [--port <port>]",
	Short:     "run serv on all template files under template root",
	Long:      "run serv on all template files under template root",
}

var (
	servPort    = servCmd.Flag.Uint("port", 8192, "Port to listen on.")
	servPkgRoot = servCmd.Flag.String(
		"package_root", ".", "go packages root directory containing templates.")
	servTemplateDir = servCmd.Flag.String(
		"template_dir", "", "Template directory relative to package root.")
	servOutputDir = servCmd.Flag.String(
		"output_dir", "", "Output directory relative to package root.")
	servKeepComments = servCmd.Flag.Bool("keep_comments", false, "Keep comment in output HTML.")
)

var (
	hasDirtyTemplate = true
	templateWatcher  *fsnotify.Watcher
)

func runServ(cmd *Command, args []string) {
	if *servTemplateDir == "" {
		fmt.Fprintf(os.Stderr, "flag template_dir is required.\n\n")
		os.Exit(2)
	} else if *servOutputDir == "" {
		servOutputDir = servTemplateDir
	}

	fmt.Fprintf(os.Stderr, "Package root directory: %s\n", *servPkgRoot)
	fmt.Fprintf(os.Stderr, "Template directory: %s\n", *servTemplateDir)
	fmt.Fprintf(os.Stderr, "Output directory: %s\n", *servOutputDir)
	fmt.Fprintf(os.Stderr, "Listen and serve on port: %d\n\n", *servPort)

	go watchTemplates()
	startHttpServer()
}

func init() {
	servCmd.Run = runServ
}

func startHttpServer() {
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/template/", templateHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *servPort), nil)
	if err != nil {
		log.Fatal("Error starting GOATS development server.", err)
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	// serve the main page
	io.WriteString(w, "main page")
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	if hasDirtyTemplate {
		fmt.Fprint(os.Stderr, "Regen template go files.")

		// Reprocess template files.
		cmd := exec.Command("goats", "gen",
			"--package_root", *servPkgRoot,
			"--template_dir", *servTemplateDir,
			"--output_dir", *servOutputDir,
			"--keep_comments", fmt.Sprintf("%b", *servKeepComments))
		var out bytes.Buffer
		cmd.Stderr = &out
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			fmt.Fprint(w, "Error to run regenerate go files from templates.\n")
			fmt.Println(out.String())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hasDirtyTemplate = false
	}

	fmt.Fprintf(os.Stderr, "\nrequest: %v\n\n", *r)

	if r.Method == "POST" {
		path := r.URL.Path[len("/template/"):]
		pkg := filepath.Dir(path)
		template := filepath.Base(path)
		cmd := exec.Command("go", "run", fmt.Sprintf("%s/cmd/main.go", pkg), template)
		cmd.Stdin = r.Body
		var out bytes.Buffer
		cmd.Stderr = &out
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error to run command \"go run %s/cmd/main.go %s\"\n", pkg, template)
			fmt.Println(out.String())
			w.WriteHeader(http.StatusInternalServerError)
		}

		io.WriteString(w, out.String())
	}
}

func watchTemplates() {
	var err error
	templateWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Process events
	go handleNotification()

	templatePath := filepath.Join(*servPkgRoot, *servTemplateDir)
	err = templateWatcher.Watch(templatePath)
	if err != nil {
		log.Fatal("Watch templates directory failed: ", templatePath, err)
	}

	var x int
	fmt.Scan(&x)
}

func handleNotification() {
	for {
		select {
		case e := <-templateWatcher.Event:
			if e.IsCreate() || e.IsDelete() || e.IsModify() || e.IsRename() {
				hasDirtyTemplate = true
			}
		case err := <-templateWatcher.Error:
			log.Println("error:", err)
		}
	}
}
