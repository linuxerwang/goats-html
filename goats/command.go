package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// A Command is an implementation of a goats command like goats build or goats fix.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	Name string

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

func (c *Command) Usage() {
	fmt.Printf(strings.TrimSpace(c.Short))
	fmt.Printf(`
Usage:

    %s

`, c.UsageLine)
	fmt.Printf("%s\n\n", strings.TrimSpace(c.Long))
	c.Flag.PrintDefaults()
	os.Exit(2)
}

func (c *Command) Exec(args []string) {
	c.Flag.Parse(args)

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(fatal); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	c.Run(c, c.Flag.Args())
}

type fatal struct{}
