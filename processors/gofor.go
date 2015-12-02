package processors

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/linuxerwang/goats-html/symbolmgr"
	"github.com/linuxerwang/goats-html/util"
)

type GoForProcessor struct {
	BaseProcessor
	varName        string
	collectionName string
}

func (f *GoForProcessor) Process(writer io.Writer, ctx *TagContext) {
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
				loopVarName = util.TrimWhiteSpaces(val[1:])
			} else {
				if itemVar1 == "" {
					itemVar1 = util.TrimWhiteSpaces(val)
				} else if itemVar2 == "" {
					itemVar2 = util.TrimWhiteSpaces(val)
				} else {
					log.Fatalf("Only two loop var can be defined in go:for, found three \"@%s\", \"@%s\" and \"%s\"\n",
						itemVar1, itemVar2, val)
				}
			}
		}
	} else {
		itemVar1 = "_"
		itemVar2 = util.TrimWhiteSpaces(f.varName)
	}

	if itemVar1 == "" {
		log.Fatalf("Loop var is required for go:for but none was found.\n")
	}
	if itemVar2 == "" {
		itemVar2 = itemVar1
		itemVar1 = "_"
	}

	expr, err := ctx.RewriteExpression(f.collectionName)
	if err != nil {
		panic(err)
	}

	// Create a new symbol map.
	sm := make(map[string]*symbolmgr.Symbol)
	if itemVar1 != "" && itemVar1 != "_" {
		isPb := false
		if strings.HasSuffix(itemVar1, "[pb]") {
			isPb = true
			itemVar1 = util.TrimWhiteSpaces(itemVar1[:len(itemVar1)-4])
		}
		s := &symbolmgr.Symbol{
			Name: itemVar1,
			Type: symbolmgr.TypeFor,
			IsPb: isPb,
		}
		sm[itemVar1] = s
	}
	if itemVar2 != "" && itemVar2 != "_" {
		isPb := false
		if strings.HasSuffix(itemVar2, "[pb]") {
			isPb = true
			itemVar2 = util.TrimWhiteSpaces(itemVar2[:len(itemVar2)-4])
		}
		s := &symbolmgr.Symbol{
			Name: itemVar2,
			Type: symbolmgr.TypeFor,
			IsPb: isPb,
		}
		sm[itemVar2] = s
	}
	if loopVarName != "" && loopVarName != "_" {
		s := &symbolmgr.Symbol{
			Name: loopVarName,
			Type: symbolmgr.TypeFor,
			IsPb: false, // loop var is always not PB.
		}
		sm[loopVarName] = s
	}
	ctx.symMgr.Push(sm)

	switch ctx.OutputFormat {
	case "go":
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
			fmt.Sprintf("for %s, %s := range __loopItems {\n", util.TrimWhiteSpaces(itemVar1), util.TrimWhiteSpaces(itemVar2)))

		if loopVarName != "" {
			io.WriteString(writer, fmt.Sprintf("%s.Counter++\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.Counter0++\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.RevCounter--\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.RevCounter0--\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.Last = (%s.Total == %s.Counter)\n", loopVarName, loopVarName, loopVarName))
		}
	case "closure":
		if itemVar1 == "_" {
			itemVar1 = "__index"
		}
		io.WriteString(writer, fmt.Sprintf("var __loopItems = %s;\n", expr))
		io.WriteString(writer, "if (__loopItems) {\n")

		if loopVarName != "" {
			io.WriteString(writer, fmt.Sprintf("var __loopTotal = __loopItems.length;\n"))
			io.WriteString(writer,
				fmt.Sprintf("var %s = {\n"+
					"total: __loopTotal,\n"+
					"counter: 0,\n"+
					"counter0: -1,\n"+
					"revCounter: __loopTotal + 1,\n"+
					"revCounter0: __loopTotal,\n"+
					"first: __loopTotal > 1,\n"+
					"last: false\n"+
					"};\n", loopVarName))
		}

		// Start of for loop.
		io.WriteString(writer, fmt.Sprintf("goog.array.forEach(__loopItems, function(%s, %s) {\n", util.TrimWhiteSpaces(itemVar2), util.TrimWhiteSpaces(itemVar1)))

		if loopVarName != "" {
			io.WriteString(writer, fmt.Sprintf("%s.counter++;\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.counter0++;\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.revCounter--;\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.revCounter0--;\n", loopVarName))
			io.WriteString(writer, fmt.Sprintf("%s.last = (%s.total == %s.counter);\n", loopVarName, loopVarName, loopVarName))
		}
	}

	if f.next != nil {
		f.next.Process(writer, ctx)
	}

	if loopVarName != "" {
		switch ctx.OutputFormat {
		case "go":
			io.WriteString(writer, fmt.Sprintf("\n%s.First = false\n", loopVarName))
		case "closure":
			io.WriteString(writer, fmt.Sprintf("\n%s.first = false;\n", loopVarName))
		}
	}

	// End of for loop.
	switch ctx.OutputFormat {
	case "go":
		io.WriteString(writer, "}\n")
	case "closure":
		io.WriteString(writer, "}, this);\n")
		io.WriteString(writer, "}\n")
	}

	ctx.symMgr.Pop()

	// End of the local scope
	io.WriteString(writer, "}\n")
}

func NewForProcessor(forLine string) *GoForProcessor {
	varName, varVal := util.SplitVarDef(forLine)
	processor := &GoForProcessor{
		varName:        util.TrimWhiteSpaces(varName),
		collectionName: util.TrimWhiteSpaces(varVal),
	}
	return processor
}
