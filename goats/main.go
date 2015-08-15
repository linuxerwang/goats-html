package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var commands = []*Command{
	genCmd,
	servCmd,
	versionCmd,
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			cmd.Run(cmd, args)
			return
		}
	}

	// No subcommand is hit.
	fmt.Fprint(os.Stderr, "Unknown command.\n\n")
	usage()
}

func printUsage(writer io.Writer) {
	io.WriteString(writer, `goats is the command for GOATS (GO Attribute-based Template System).

Usage:

	goats command [arguments]

The commands are:

	gen	Generate go files for html templates.
	serv	Start an http server to serve the generated templates for development purpose.
	version	Show version.

Use "goats help [command]" for more information about a command.

`)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'goats help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: goats help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'goats help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name == arg {
			cmd.Usage()
			// not exit 2: succeeded at 'goats help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'goats help'.\n", arg)
	os.Exit(2) // failed at 'goats help cmd'
}
