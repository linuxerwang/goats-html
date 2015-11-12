package processors

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/linuxerwang/goats-html/expl"
	"github.com/linuxerwang/goats-html/pkgmgr"
	"github.com/linuxerwang/goats-html/symbolmgr"
	"github.com/linuxerwang/goats-html/util"
)

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

type TagContext struct {
	ContextId    int
	ExprParser   *expl.ExprParser
	pkgMgr       *pkgmgr.PkgManager
	pkgRefs      pkgmgr.AliasReferer
	symMgr       *symbolmgr.SymbolMgr
	fitlers      map[string]*RegisteredFilter
	AutoEscape   bool
	OutputFormat string
}

func (ctx *TagContext) GetFilters() map[string]*RegisteredFilter {
	return ctx.fitlers
}

func (ctx *TagContext) MaybeAddImports(expression string) {
	expression = util.TrimWhiteSpaces(expression)
	if expression == "" {
		return
	}

	f, err := parser.ParseExpr(util.ToGoString(expression))
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			selector := expression[x.Pos()-1 : x.End()-1]
			ctx.pkgRefs.RefByAlias(strings.Split(selector, ".")[0], false)
			return false
		}
		return true
	})
}

func (ctx *TagContext) RewriteExpression(originalExpr string) (string, error) {
	ctx.MaybeAddImports(originalExpr)
	er := expl.NewExprRewriter(ctx.symMgr, ctx.OutputFormat)
	return er.RewriteExpression(originalExpr)
}

var tagCount int = 0

func NewTagContext(pkgMgr *pkgmgr.PkgManager, pkgRefs pkgmgr.AliasReferer, outputFormat string) *TagContext {
	ctx := &TagContext{
		ContextId:    tagCount,
		pkgMgr:       pkgMgr,
		pkgRefs:      pkgRefs,
		symMgr:       symbolmgr.New(),
		fitlers:      map[string]*RegisteredFilter{},
		AutoEscape:   true,
		OutputFormat: outputFormat,
	}
	tagCount++

	sm := make(map[string]*symbolmgr.Symbol)
	if ctx.symMgr.Size() == 0 {
		for alias, pkgImpt := range pkgMgr.GetPkgsForAlias() {
			sm[alias] = &symbolmgr.Symbol{
				Name:    alias,
				Type:    symbolmgr.TypeImport,
				IsPb:    pkgImpt.PbPkg() != "",
				PkgImpt: pkgImpt,
			}
		}
	}
	ctx.symMgr.Push(sm)

	ctx.ExprParser = expl.NewExprParser(ctx)
	return ctx
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

func (gheadp *GoHeadProcessor) Process(writer io.Writer, ctx *TagContext) {
	sm := make(map[string]*symbolmgr.Symbol)
	ctx.symMgr.Push(sm)

	if gheadp.next != nil {
		gheadp.next.Process(writer, ctx)
	}

	ctx.symMgr.Pop()
}

func NewHeadProcessor() *GoHeadProcessor {
	return &GoHeadProcessor{}
}

type ProcessorCallbackFunc func()

/*
 * Callback processor gives a chance to execute a special processing when
 * the processing chain reaches its end but before the call stack starts to
 * go back.
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
