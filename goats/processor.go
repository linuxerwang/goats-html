package goats

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
)

type TagContext struct {
	imports    map[string]bool
	fitlers    map[string]*RegisteredFilter
	AutoEscape bool
}

func (ctx *TagContext) GetImports() map[string]bool {
	return ctx.imports
}

func (ctx *TagContext) GetFilters() map[string]*RegisteredFilter {
	return ctx.fitlers
}

// Deprecated
func (ctx *TagContext) MaybeAddImports(expression string) {
	expression = TrimWhiteSpaces(expression)
	if expression == "" {
		return
	}

	f, err := parser.ParseExpr(ToGoString(expression))
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			selector := expression[x.Pos()-1 : x.End()-1]
			ctx.imports[strings.Split(selector, ".")[0]] = true
			return false
		}
		return true
	})
}

type FilterInstance struct {
	Name   string
	Start  token.Pos
	End    token.Pos
	Pivot  token.Pos // Unused
	Params string
}

type SelectorInstance struct {
	Name    string
	NewName string
	Start   token.Pos
	End     token.Pos
}

func (ctx *TagContext) RewriteExpression(originalExpr string) string {
	originalExpr = ToGoString(TrimWhiteSpaces(originalExpr))
	if originalExpr == "" {
		return ""
	}

	expr := originalExpr
	ctx.MaybeAddImports(expr)

	// Non builtin function calls in template should be disallowed.
	if f := ctx.findNonBuiltinFunc(expr); f != "" {
		panic("Function call is not allowed in goats: " + f)
	}

	// Convert snake-case to camel case so that template writer can write protocol buffer fields.
	for selector := ctx.findSelector(expr); selector != nil; selector = ctx.findSelector(expr) {
		if selector.Name != selector.NewName {
			expr = expr[0:selector.Start-1] + selector.NewName + expr[selector.End-1:]
		}
	}

	// Convert filter to a function call.
	for filter := ctx.findFilter(expr); filter != nil; filter = ctx.findFilter(expr) {
		var converted string
		if filter.Params != "" {
			converted = fmt.Sprintf("__impl.%s(%s)", filter.Name, filter.Params)
		} else {
			converted = fmt.Sprintf("__impl.%s(%s)", filter.Name, expr[filter.Start-1:filter.Pivot-1])
		}
		expr = expr[0:filter.Start-1] + converted + expr[filter.End-1:]
	}

	return expr
}

// TODO: Replace with a more effecient design.
func (ctx *TagContext) findFilter(expr string) *FilterInstance {
	f, err := parser.ParseExpr(expr)
	if err != nil {
		panic(err)
	}

	var filter *FilterInstance = nil
	ast.Inspect(f, func(n ast.Node) bool {
		switch y := n.(type) {
		case *ast.Ident:
			if registeredFilter, ok := RegisteredFilters[y.Name]; ok {
				filter = &FilterInstance{
					Name:   registeredFilter.VarName + "." + builtinFuncs[y.Name],
					Start:  n.Pos(),
					End:    n.End(),
					Params: "",
				}
				ctx.fitlers[registeredFilter.VarName] = registeredFilter
			}
		case *ast.CallExpr:
			name := ""
			switch z := y.Fun.(type) {
			case *ast.Ident:
				name = z.Name
			case *ast.SelectorExpr:
				// Must be already processed built-in function, continue
				return true
			default:
				panic("Function call has no name.  This is impossible.")
			}

			if registeredFilter, ok := RegisteredFilters[name]; ok {
				filter = &FilterInstance{
					Name:   registeredFilter.VarName + "." + builtinFuncs[name],
					Start:  n.Pos(),
					End:    n.End(),
					Params: strings.Trim(expr[y.Lparen:y.Rparen-1], " \t\n"),
				}
				ctx.fitlers[registeredFilter.VarName] = registeredFilter
			}
		}
		return filter == nil
	})
	return filter
}

func (ctx *TagContext) findNonBuiltinFunc(expr string) string {
	f, err := parser.ParseExpr(expr)
	if err != nil {
		panic(err)
	}

	funcName := ""
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			switch y := x.Fun.(type) {
			case *ast.Ident:
				if _, found := builtinFuncs[y.Name]; found {
					return true
				} else {
					funcName = y.Name
					return false
				}
			default:
				panic("Function call has no name.  This is impossible.")
			}
		default:
			return true
		}
	})
	return funcName
}

// TODO: Replace with a more effecient design.
func (ctx *TagContext) findSelector(expr string) *SelectorInstance {
	f, err := parser.ParseExpr(expr)
	if err != nil {
		panic(err)
	}

	var selector *SelectorInstance
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			if StartsWithLower(x.Sel.Name) {
				selector = &SelectorInstance{
					Name:    x.Sel.Name,
					NewName: ToCamelCase(x.Sel.Name),
					Start:   x.Sel.NamePos,
					End:     x.Sel.NamePos + token.Pos(len(x.Sel.Name)),
				}
			}
		}
		return selector == nil
	})
	return selector
}

func NewTagContext() *TagContext {
	return &TagContext{
		imports:    map[string]bool{},
		fitlers:    map[string]*RegisteredFilter{},
		AutoEscape: true,
	}
}

type Processor interface {
	HasNext() bool
	GetNext() Processor
	SetNext(next Processor)
	Process(writer io.Writer, context *TagContext)
}

type BaseProcessor struct {
	next Processor
}

func (b *BaseProcessor) HasNext() bool {
	return b.next != nil
}

func (b *BaseProcessor) GetNext() Processor {
	return b.next
}

func (b *BaseProcessor) SetNext(next Processor) {
	b.next = next
}

type GoHeadProcessor struct {
	BaseProcessor
}

func (i *GoHeadProcessor) Process(writer io.Writer, context *TagContext) {
	// Do nothing, only pass to next processor.
	if i.next != nil {
		i.next.Process(writer, context)
	}
}

func NewHeadProcessor() *GoHeadProcessor {
	return &GoHeadProcessor{}
}

type ProcessorCallbackFunc func()

/*
 * Callback processor gives a chance to execute a special processing when the processing chain
 * reaches its end but before the call stack starts to go back.
 */
type GoCallbackProcessor struct {
	BaseProcessor
	callback ProcessorCallbackFunc
}

func (c *GoCallbackProcessor) Process(writer io.Writer, context *TagContext) {
	// Do nothing, only call the callback function
	c.callback()
}

func NewCallbackProcessor(callback ProcessorCallbackFunc) *GoCallbackProcessor {
	return &GoCallbackProcessor{
		callback: callback,
	}
}

type RegisteredFilter struct {
	PkgName string
	Type    string
	VarName string
}

var builtinFilter = &RegisteredFilter{
	PkgName: "goats-html/goats/runtime",
	Type:    "BuiltinFilter",
	VarName: "builtinFilter",
}

var RegisteredFilters = map[string]*RegisteredFilter{
	"capfirst":     builtinFilter,
	"center":       builtinFilter,
	"cut":          builtinFilter,
	"debug":        builtinFilter,
	"floatformat":  builtinFilter,
	"format":       builtinFilter,
	"join":         builtinFilter,
	"ljust":        builtinFilter,
	"rjust":        builtinFilter,
	"title":        builtinFilter,
	"quote":        builtinFilter,
	"unixdate":     builtinFilter,
	"unixnanodate": builtinFilter,
}

var builtinFuncs = map[string]string{
	"center":       "Center",
	"cut":          "Cut",
	"debug":        "Debug",
	"floatformat":  "FloatFormat",
	"format":       "Format",
	"join":         "Join",
	"len":          "Length",
	"ljust":        "Ljust",
	"rjust":        "Rjust",
	"title":        "Title",
	"quote":        "Quote",
	"unixdate":     "UnixDate",
	"unixnanodate": "UnixNanoDate",
}
