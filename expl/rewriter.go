package expl

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"github.com/linuxerwang/goats-html/symbolmgr"
	"github.com/linuxerwang/goats-html/util"
)

// ExprRewriter parses an input expression and rewrite it according to the
// target lauguage format:
// - convert snake-case names to camel case for Go.
// - convert slices to slice function for closure library.
// - convert builtin functions to proper function call.
// - disallow non-builtin functions.
// - disallow 3-index slicing.
// - disallow vararg parameter passing.
// - register imports for Go.
type ExprRewriter struct {
	symMgr       *symbolmgr.SymbolMgr
	outputFormat string
	filterPrefix string

	// These are the flags to control identifier translation.

	shouldNotChange bool // Do not change the expression if shouldNotChange is true.
	inPbExpression  bool // Whether it's in a PB expression.
	inSelector      bool // Whether it's in a Go selector.
}

func (er *ExprRewriter) RewriteExpression(expr string) (string, error) {
	expr = util.ToGoString(expr)
	if expr == "" {
		return "", nil
	}

	f, err := parser.ParseExpr(expr)
	if err != nil {
		return "", err
	}

	return er.processExpr(f, false /* needQuote */)
}

// processExpr recursively traverses the given AST tree of the expression,
// processes each node according to the target lauguage format.
func (er *ExprRewriter) processExpr(n ast.Node, needQuote bool) (string, error) {
	switch node := n.(type) {
	case *ast.Ident:
		name := node.String()
		if !er.shouldNotChange && er.inPbExpression {
			switch er.outputFormat {
			case "go":
				return util.ToCamelCase(name), nil
			case "closure":
				r := []rune(name)
				if unicode.IsLower(r[0]) {
					return "get" + util.ToCamelCase(name) + "()", nil
				}
			}
		}
		return name, nil
	case *ast.BinaryExpr:
		x, err := er.processExpr(node.X, needQuote)
		if err != nil {
			return "", err
		}

		y, err := er.processExpr(node.Y, needQuote)
		if err != nil {
			return "", err
		}

		return x + node.Op.String() + y, nil
	case *ast.StarExpr:
		switch er.outputFormat {
		case "go":
			result, err := er.processExpr(node.X, needQuote)
			if err != nil {
				return "", err
			}

			return "*" + result, nil
		case "closure":
			return er.processExpr(node.X, needQuote)
		}
	case *ast.BasicLit:
		if needQuote && node.Kind == token.STRING {
			return strconv.Quote(node.Value), nil
		}
		return node.Value, nil
	case *ast.ParenExpr:
		x, err := er.processExpr(node.X, needQuote)
		if err != nil {
			return "", err
		}
		return "(" + x + ")", nil
	case *ast.SelectorExpr:
		// Suppose the selector expression is "a.b.c.d", the AST tree would
		// be like:
		//
		//       dot
		//      /   \
		//     dot   d
		//    /   \
		//   dot   c
		//  /   \
		// a     b
		//
		// Thus the right-most identifier will be visited first. To correctly
		// translate a pb/non-pb selector expression for both Go and Closure,
		// we have to set/reset the boolean flags when it comes to the
		// right-most identifier (which is always the first ast.SelectorExpr).
		// The field inSelector is introduced to help tracking whether it's
		// the right-most identifier.
		//
		// The left-most identifier determines whether the whole selection
		// expression is a protocol buffer (by looking up it as an alias to
		// a protocol buffer import). For each ast.SelectorExpr if it's an
		// *ast.Ident, we know it's the left-most identifier. We then set
		// the field inPbExpression to true. This flag will be propogated
		// during back tracking until it goes back to the right-most
		// identifier, where the flag will be reset.
		//
		// One special case we need to handle is when the right-most
		// identifier's first letter is capitalized, so the identifier
		// represents a protocol buffer enum, and we should keep the whole
		// selector expression as lower cased and snake-cased. Field shouldNotChange is
		// introduced as a flag which was set when it enters the right-most
		// identifier.

		if !er.inSelector {
			// Now it starts from the right-most identifier.
			er.inSelector = true

			// The right-most identifier's first letter determines if
			// translation should be performed for the whole selector
			// expression.
			r := []rune(node.Sel.Name)
			if unicode.IsUpper(r[0]) {
				er.shouldNotChange = true
			}

			// Reset the flags when it back tracked to the right-most
			// identifier.
			defer func() {
				er.shouldNotChange = false
				er.inSelector = false
				er.inPbExpression = false
			}()
		}

		x := ""
		var err error
		if id, ok := node.X.(*ast.Ident); ok {
			// Now it reaches the left-most identifier. Let's check if the
			// selector expression is a protocol buffer by looking up the
			// symbol map.
			//
			// The left-most identifier will never be translated.
			x = id.Name
			s, err := er.symMgr.GetSymbol(id.Name)
			if err != nil {
				return "", err
			}
			er.inPbExpression = s.IsPb
			if s.Type == symbolmgr.TypeImport {
				x = s.PkgImpt.PbPkg()
			}
		} else {
			// Not the left-most identifier, keep traversing the AST.
			x, err = er.processExpr(node.X, needQuote)
			if err != nil {
				return "", err
			}
		}

		sel, err := er.processExpr(node.Sel, needQuote)
		if err != nil {
			return "", err
		}

		return x + "." + sel, nil
	case *ast.IndexExpr:
		x, err := er.processExpr(node.X, needQuote)
		if err != nil {
			return "", err
		}

		idx, err := er.processExpr(node.Index, needQuote)
		if err != nil {
			return "", err
		}

		return x + "[" + idx + "]", nil

	case *ast.SliceExpr:
		if node.Slice3 {
			return "", errors.New("Only one colon is allowed for slices.")
		}

		x, err := er.processExpr(node.X, needQuote)
		if err != nil {
			return "", err
		}

		s := []string{x, "["}
		if node.Low != nil {
			low, err := er.processExpr(node.Low, needQuote)
			if err != nil {
				return "", err
			}

			s = append(s, low)
		}
		s = append(s, ":")
		if node.High != nil {
			high, err := er.processExpr(node.High, needQuote)
			if err != nil {
				return "", err
			}

			s = append(s, high)
		}
		s = append(s, "]")
		return strings.Join(s, ""), nil
	case *ast.CallExpr:
		fun, err := er.processExpr(node.Fun, true /* needQuote */)
		if err != nil {
			return "", err
		}

		if node.Ellipsis > 0 {
			return "", errors.New("Vararg is not allowed in goats: " + fun)
		}

		switch er.outputFormat {
		case "go":
			if f, ok := builtinFuncs[fun]; ok {
				fun = f
			} else {
				return "", errors.New("Function call is not allowed in goats: " + fun)
			}
		case "closure":
			// keep the all-lower-case name for javascript.
		}

		s := []string{er.filterPrefix, ".", fun, "("}
		for i, arg := range node.Args {
			a, err := er.processExpr(arg, needQuote)
			if err != nil {
				return "", err
			}

			s = append(s, a)
			if i < len(node.Args)-1 {
				s = append(s, ",")
			}
		}
		s = append(s, ")")
		return strings.Join(s, ""), nil
	case *ast.UnaryExpr:
		x, err := er.processExpr(node.X, needQuote)
		if err != nil {
			return "", err
		}

		return node.Op.String() + x, nil
	}
	return "", nil
}

// NewExprRewriter creates a new ExprRewriter for the given target format.
func NewExprRewriter(symMgr *symbolmgr.SymbolMgr, format string) ExprHandler {
	var prefix string
	switch format {
	case "go":
		prefix = "__impl"
	case "closure":
		prefix = "goats.runtime.filters"
	}
	return &ExprRewriter{
		symMgr:       symMgr,
		outputFormat: format,
		filterPrefix: prefix,
	}
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

func lowerFirstLetter(input string) string {
	s := []rune(input)
	s[0] = unicode.ToLower(s[0])
	return string(s)
}
