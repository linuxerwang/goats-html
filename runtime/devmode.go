// +build goats_devmod

package runtime

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var (
	goatsSettings    *GoatsSettings
	devServerProcess *os.Process
)

func NewGoatsSettings() *GoatsSettings {
	return &GoatsSettings{
		DevServerPort: 8192,
		PkgRoot:       ".",
		TemplateDir:   "",
		OutputDir:     "",
	}
}

func InitGoats(settings *GoatsSettings) {
	if settings == nil {
		goatsSettings = NewGoatsSettings()
	} else {
		goatsSettings = settings
	}

	// Start dev http server.
	startDevServer()
}

func startDevServer() {
	fmt.Println("Starting GOATS development server.\n")
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "."
	}

	procAttr := &os.ProcAttr{
		Dir:   workDir,
		Files: []*os.File{nil, os.Stdout, os.Stderr},
	}

	outputDir := goatsSettings.OutputDir
	if outputDir == "" {
		outputDir = goatsSettings.TemplateDir
	}

	executable, err := exec.LookPath("goats")
	if err != nil {
		executable = "goats"
	}
	devServerProcess, err = os.StartProcess(
		executable,
		[]string{
			"goats", "serv",
			fmt.Sprintf("--port=%d", goatsSettings.DevServerPort),
			fmt.Sprintf("--package_root=%s", goatsSettings.PkgRoot),
			fmt.Sprintf("--template_dir=%s", goatsSettings.TemplateDir),
			fmt.Sprintf("--output_dir=%s", outputDir),
			fmt.Sprintf("--keep_comments=true"),
		},
		procAttr)
	if err != nil {
		fmt.Println(err)
		panic("Failed to start goats dev server.")
	}

	// Sleep for one second to make sure server starts.
	time.Sleep(1000 * time.Millisecond)
}

func CallRpc(
	pkg,
	template string,
	settings *TemplateSettings,
	args interface{},
	writer io.Writer) error {
	fmt.Println("call rpc:", pkg, template)
	gob.Register(args)
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(settings)
	if err != nil {
		fmt.Println("===> Encode settings error:", err)
	}
	err = encoder.Encode(args)
	if err != nil {
		fmt.Println("===> Encode args error:", err)
	}

	response, err := http.Post(
		fmt.Sprintf("http://localhost:%d/template/%s/%s", goatsSettings.DevServerPort, pkg, template),
		"application/gob",
		&buffer)
	if err != nil {
		fmt.Println("Failed!", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("Dev server returned error.")
	}

	body, err := ioutil.ReadAll(response.Body)
	writer.Write(body)
	return err
}
