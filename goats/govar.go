package goats

import (
	"fmt"
	"io"
	"strings"
)

type VarDef struct {
	Name string
	Val  string
}

type GoVarProcessor struct {
	BaseProcessor
	vars []*VarDef
}

func (v *GoVarProcessor) Process(writer io.Writer, context *TagContext) {
	// Start of a local scope
	io.WriteString(writer, "{\n")

	for _, varDef := range v.vars {
		expr := context.RewriteExpression(varDef.Val)
		io.WriteString(writer, fmt.Sprintf("%s := %s", varDef.Name, expr))
		io.WriteString(writer, "\n")
	}

	if v.next != nil {
		v.next.Process(writer, context)
	}

	// End of a local scope
	io.WriteString(writer, "}\n")
}

func newVarDef(varStr string) *VarDef {
	varName, varVal := SplitVarDef(varStr)
	return &VarDef{
		Name: TrimWhiteSpaces(varName),
		Val:  TrimWhiteSpaces(varVal),
	}
}

func NewVarProcessor(varDef string) *GoVarProcessor {
	processor := &GoVarProcessor{
		vars: []*VarDef{newVarDef(varDef)},
	}
	return processor
}

func NewVarsProcessor(varStrs string) *GoVarProcessor {
	varDefs := []*VarDef{}
	for _, varStr := range strings.Split(varStrs, ";") {
		varDefs = append(varDefs, newVarDef(varStr))
	}
	processor := &GoVarProcessor{
		vars: varDefs,
	}
	return processor
}
