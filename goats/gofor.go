package goats

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type GoForProcessor struct {
	BaseProcessor
	varName        string
	collectionName string
}

func (f *GoForProcessor) Process(writer io.Writer, context *TagContext) {
	// Start of a local scope
	io.WriteString(writer, "{\n")

	loopVarName := ""
	itemVar1 := ""
	itemVar2 := ""
	if strings.Contains(f.varName, ",") {
		parts := strings.Split(f.varName, ",")
		for _, val := range parts {
			if strings.Index(val, "@") > -1 {
				if loopVarName != "" {
					log.Fatalf("Only one @ loop var can be defined in go:for, found two \"@%s\" and \"%s\"\n",
						loopVarName, val)
				}
				loopVarName = TrimWhiteSpaces(val[1:])
			} else {
				if itemVar1 == "" {
					itemVar1 = val
				} else if itemVar2 == "" {
					itemVar2 = val
				} else {
					log.Fatalf("Only two loop var can be defined in go:for, found three \"@%s\", \"@%s\" and \"%s\"\n",
						itemVar1, itemVar2, val)
				}
			}
		}
	} else {
		itemVar1 = "_"
		itemVar2 = f.varName
	}

	if itemVar1 == "" {
		log.Fatalf("Loop var is required for go:for but none was found.\n")
	}
	if itemVar2 == "" {
		itemVar2 = itemVar1
		itemVar1 = "_"
	}

	expr := context.RewriteExpression(f.collectionName)

	io.WriteString(writer, fmt.Sprintf("__loopItems := %s\n", expr))
	if loopVarName != "" {
		io.WriteString(writer, fmt.Sprintf("__loopTotal := len(__loopItems)\n"))
		io.WriteString(writer,
			fmt.Sprintf("%s := &runtime.LoopVar{\n"+
				"Total: __loopTotal,\n"+
				"Counter0: -1,\n"+
				"RevCounter: __loopTotal + 1,\n"+
				"RevCounter0: __loopTotal,\n"+
				"First: __loopTotal > 1,\n"+
				"}\n", loopVarName))
	}

	// Start of for loop.
	io.WriteString(writer,
		fmt.Sprintf("for %s, %s := range __loopItems {\n",
			TrimWhiteSpaces(itemVar1), TrimWhiteSpaces(itemVar2)))

	if loopVarName != "" {
		io.WriteString(writer, fmt.Sprintf("%s.Counter++\n", loopVarName))
		io.WriteString(writer, fmt.Sprintf("%s.Counter0++\n", loopVarName))
		io.WriteString(writer, fmt.Sprintf("%s.RevCounter--\n", loopVarName))
		io.WriteString(writer, fmt.Sprintf("%s.RevCounter0--\n", loopVarName))
		io.WriteString(writer,
			fmt.Sprintf("%s.Last = (%s.Total == %s.Counter)\n", loopVarName, loopVarName, loopVarName))
	}

	if f.next != nil {
		f.next.Process(writer, context)
	}

	if loopVarName != "" {
		io.WriteString(writer, fmt.Sprintf("\n%s.First = false\n", loopVarName))
	}

	// End of for loop.
	io.WriteString(writer, "}\n")

	// End of the local scope
	io.WriteString(writer, "}\n")
}

func NewForProcessor(forLine string) *GoForProcessor {
	parts := strings.Split(forLine, ":")
	processor := &GoForProcessor{
		varName:        TrimWhiteSpaces(parts[0]),
		collectionName: TrimWhiteSpaces(parts[1]),
	}
	return processor
}
