package processors

import (
	"fmt"
	"io"
	"strings"

	"github.com/linuxerwang/goats-html/symbolmgr"
	"github.com/linuxerwang/goats-html/util"
)

type VarDef struct {
	Name string
	Val  string
	IsPb bool
}

type GoVarProcessor struct {
	BaseProcessor
	vars []*VarDef
}

func (gvp *GoVarProcessor) Process(writer io.Writer, ctx *TagContext) {
	sm := ctx.symMgr.Peek()
	for _, v := range gvp.vars {
		sm[v.Name] = &symbolmgr.Symbol{
			Name: v.Name,
			Type: symbolmgr.TypeVar,
			IsPb: v.IsPb,
		}
	}

	// Start of a local scope
	io.WriteString(writer, "{\n")

	for _, varDef := range gvp.vars {
		ctx.ExprParser.Evaluate(varDef.Val, writer, func(expr string) {
			switch ctx.OutputFormat {
			case "go":
				io.WriteString(writer, fmt.Sprintf("%s := %s\n", varDef.Name, expr))
			case "closure":
				io.WriteString(writer, fmt.Sprintf("var %s = %s;\n", varDef.Name, expr))
			}
		})
	}

	if gvp.next != nil {
		gvp.next.Process(writer, ctx)
	}

	// End of a local scope
	io.WriteString(writer, "}\n")
}

func newVarDef(varStr string) *VarDef {
	varName, varVal := util.SplitVarDef(varStr)
	isPb := false
	if i, j := strings.Index(varName, "["), strings.Index(varName, "]"); i > -1 && j > -1 {
		varKind := varName[i+1 : j]
		varName = varName[:i]
		switch varKind {
		case "pb":
			isPb = true
		}
	}
	return &VarDef{
		Name: util.TrimWhiteSpaces(varName),
		Val:  util.TrimWhiteSpaces(varVal),
		IsPb: isPb,
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
