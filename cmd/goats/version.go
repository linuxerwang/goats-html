package main

import (
	"fmt"
	"os"
)

const VERSION = "0.1.0"

var versionCmd = &Command{
	Name:      "version",
	UsageLine: "goats version",
	Short:     "Show version",
	Long:      "Show version",
}

func runVersion(cmd *Command, args []string) {
	fmt.Fprintf(os.Stdout, "%s\n", VERSION)
}

func init() {
	versionCmd.Run = runVersion
}
