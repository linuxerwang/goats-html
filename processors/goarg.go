package processors

import (
	"fmt"
	"io"
	"strings"

	"github.com/linuxerwang/goats-html/util"
)

/*
 * There are two different places a go:arg can appear:
 * - A template definition
 *     <div go:template="<template name>" go:arg="<name>: <type> [= expression]"></div>
 * - A template call
 *     <div go:call="#<template name>" go:arg="<name>: <expression>"></div>
 */
type Argument struct {
	Name    string
	PkgName string
	Type    string
	Val     string
	Declare string
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
	}
}

func NewArgCall(argCall string) *Argument {
	var argName, argVal string
	colon := strings.Index(argCall, ":")
	argName = util.TrimWhiteSpaces(argCall[:colon])
	argVal = util.ToCamelExpr(util.ToGoString(util.TrimWhiteSpaces(argCall[colon+1:])))

	return &Argument{
		Name: argName,
		Val:  argVal,
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
	args []*Argument
}

func (a *GoArgProcessor) Process(writer io.Writer, context *TagContext) {
	for _, arg := range a.args {
		io.WriteString(writer, fmt.Sprintf("%s := __args.%s\n", arg.Name, util.ToPublicName(arg.Name)))
	}

	if a.next != nil {
		a.next.Process(writer, context)
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
