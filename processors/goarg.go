package processors

import (
	"fmt"
	"io"
	"strings"

	"github.com/linuxerwang/goats-html/symbolmgr"
	"github.com/linuxerwang/goats-html/util"
)

/*
 * There are three different places a go:arg can appear:
 * - A template definition
 *     <div go:template="<template name>" go:arg="<name>[pb]: <type> [= expression]"></div>
 * - A replaceable definition
 *     <div go:replaceable="<replaceable name>" go:arg="<name>[pb]: <expression>"></div>
 * - A template call
 *     <div go:call="#<template name>" go:arg="<name>[pb]: <expression>"></div>
 */
type Argument struct {
	Name    string
	PkgName string
	Type    string
	Val     string
	Declare string
	IsPb    bool
}

func NewArgDef(argDef string) *Argument {
	var argName, pkgName, argType, argVal string
	colon := strings.Index(argDef, ":")
	equal := strings.LastIndex(argDef, "=")
	if colon > -1 {
		argName = util.TrimWhiteSpaces(argDef[:colon])
		if equal > -1 && equal < len(argDef) {
			argType = util.TrimWhiteSpaces(argDef[colon+1 : equal])
			argVal = util.ToGoString(util.TrimWhiteSpaces(argDef[equal+1:]))
		} else {
			argType = util.TrimWhiteSpaces(argDef[colon+1:])
		}
	}

	isPb := false
	if i, j := strings.Index(argName, "["), strings.Index(argName, "]"); i > -1 && j > -1 {
		switch argName[i+1 : j] {
		case "pb":
			isPb = true
		}
		argName = argName[:i]
	}

	// Package name.
	pkgName, found := util.ExtractSelector(argType)
	if found {
		dot := strings.Index(pkgName, ".")
		if dot > -1 {
			pkgName = pkgName[:dot]
		}
	}

	return &Argument{
		Name:    argName,
		PkgName: pkgName,
		Type:    argType,
		Val:     argVal,
		Declare: util.ToPublicName(argName) + " " + argType,
		IsPb:    isPb,
	}
}

func NewArgCall(argCall string) *Argument {
	var argName, argVal string
	colon := strings.Index(argCall, ":")
	argName = util.TrimWhiteSpaces(argCall[:colon])
	argVal = util.TrimWhiteSpaces(argCall[colon+1:])

	isPb := false
	if i, j := strings.Index(argName, "["), strings.Index(argName, "]"); i > -1 && j > -1 {
		switch argName[i+1 : j] {
		case "pb":
			isPb = true
		}
		argName = argName[:i]
	}

	return &Argument{
		Name: argName,
		Val:  argVal,
		IsPb: isPb,
	}
}

func ParseArgDefs(argDefs string) []*Argument {
	var args []*Argument
	for _, argDef := range strings.Split(argDefs, ";") {
		if argDef != "" {
			args = append(args, NewArgDef(argDef))
		}
	}
	return args
}

type GoArgProcessor struct {
	BaseProcessor
	args  []*Argument
	where int
}

func (a *GoArgProcessor) Process(writer io.Writer, ctx *TagContext) {
	sm := ctx.symMgr.Peek()
	for _, arg := range a.args {
		sm[arg.Name] = &symbolmgr.Symbol{
			Name: arg.Name,
			Type: symbolmgr.TypeArg,
			IsPb: arg.IsPb,
		}
	}

	switch ctx.OutputFormat {
	case "go":
		for _, arg := range a.args {
			io.WriteString(writer, fmt.Sprintf("%s := __args.%s\n", arg.Name, util.ToPublicName(arg.Name)))
		}
	case "closure":
		for _, arg := range a.args {
			t := ctx.symMgr.ExpandType(arg.Type)

			if t == arg.Type {
				io.WriteString(writer, fmt.Sprintf("var %s = __args[\"%s\"];\n", arg.Name, arg.Name))
			} else {
				ctx.pkgRefs.RefClosureRequire(t)
				io.WriteString(writer, fmt.Sprintf("var %s = new %s(__args[\"%s\"]);\n", arg.Name, t, arg.Name))
			}
		}
	}

	if a.next != nil {
		a.next.Process(writer, ctx)
	}
}

func NewArgProcessor(args []*Argument) *GoArgProcessor {
	processor := &GoArgProcessor{
		args: args,
	}
	return processor
}

func ParseArgCalls(argCalls string) []*Argument {
	var args []*Argument
	for _, argCall := range strings.Split(argCalls, ";") {
		if argCall != "" {
			args = append(args, NewArgCall(argCall))
		}
	}
	return args
}
