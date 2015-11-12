package expl

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

// ExprHandler represents any object which can rewrite an expression
// according to certain rules.
type ExprHandler interface {
	RewriteExpression(originalExpr string) (string, error)
}

type OnEvalExpr func(expr string)

type ExprParser struct {
	eh         ExprHandler
	s          scanner.Scanner
	token_type []rune
	token_text []string
	ifPos      int
	elsePos    int
	elifPos    []int
}

func (ep *ExprParser) Evaluate(src string, w io.Writer, callback OnEvalExpr) {
	err := ep.parse(src)
	if err != nil {
		panic(err)
	}

	err = ep.render(w, callback)
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte("\n"))
	if err != nil {
		panic(err)
	}
}

func (ep *ExprParser) parse(src string) error {
	ep.s.Init(strings.NewReader(src))
	ep.token_type = nil
	ep.token_text = nil
	ep.ifPos = -1
	ep.elsePos = -1
	ep.elifPos = nil

	for {
		tok := ep.s.Scan()
		if tok == scanner.EOF {
			return nil
		}

		if tok == scanner.Ident && ep.s.TokenText() == "if" {
			if len(ep.token_type) == 0 {
				return errors.New("Syntax error: no expression before 'if'.")
			}
			if ep.ifPos > -1 {
				return errors.New("Syntax error: more than one 'if' was found.")
			} else if len(ep.elifPos) > 0 {
				return errors.New("Syntax error: found 'if' after 'elif'.")
			} else if ep.elsePos > -1 {
				return errors.New("Syntax error: found 'if' after 'else'.")
			}
			ep.ifPos = len(ep.token_type)
		} else if tok == scanner.Ident && ep.s.TokenText() == "elif" {
			if ep.ifPos == -1 {
				return errors.New("Syntax error: expecting 'if' but found 'elif'.")
			} else if ep.elsePos > -1 {
				return errors.New("Syntax error: found 'elif' after 'else'.")
			}
			ep.elifPos = append(ep.elifPos, len(ep.token_type))
		} else if tok == scanner.Ident && ep.s.TokenText() == "else" {
			if ep.ifPos == -1 {
				return errors.New("Syntax error: expecting 'if' but found 'else'.")
			}
			if ep.elsePos > -1 {
				return errors.New("Syntax error: more than one 'else' was found.")
			}
			ep.elsePos = len(ep.token_type)
		}

		ep.token_type = append(ep.token_type, tok)
		ep.token_text = append(ep.token_text, ep.s.TokenText())
	}

	return nil
}

func (ep *ExprParser) render(w io.Writer, callback OnEvalExpr) error {
	foundIf := ep.ifPos > -1
	foundElif := len(ep.elifPos) > 0
	foundElse := ep.elsePos > -1

	if !(foundIf || foundElif || foundElse) {
		expr, err := ep.eh.RewriteExpression(strings.Join(ep.token_text, ""))
		if err != nil {
			return err
		}
		callback(expr)
	} else {
		if ep.token_text[ep.ifPos+1] != "(" {
			return errors.New("Syntax Error: conditional must be enclosed in ()")
		}

		n := -1
		if foundElif {
			n = ep.elifPos[0]
		} else if foundElse {
			n = ep.elsePos
		} else {
			n = len(ep.token_text)
		}
		if ep.token_text[n-1] != ")" {
			return errors.New("Syntax Error: conditional must be enclosed in ()")
		}

		expr, err := ep.eh.RewriteExpression(strings.Join(ep.token_text[ep.ifPos+1:n], ""))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(fmt.Sprintf("if %s {\n", expr)))
		if err != nil {
			return err
		}

		expr, err = ep.eh.RewriteExpression(strings.Join(ep.token_text[:ep.ifPos], ""))
		if err != nil {
			return err
		}

		callback(expr)

		for i, pos := range ep.elifPos {
			n = len(ep.token_text) - 1
			if i < len(ep.elifPos)-1 {
				n = ep.elifPos[i+1] - 1
			} else if i == len(ep.elifPos)-1 && ep.elsePos > -1 {
				n = ep.elsePos - 1
			}
			if ep.token_text[pos+1] != "(" {
				return errors.New("Syntax Error: conditional must be enclosed in ()")
			}

			cond, expr, err := ep.parseCondExpr(pos+2, n)
			if err != nil {
				return err
			}

			_, err = w.Write([]byte(fmt.Sprintf("} else if (%s) {\n", cond)))
			if err != nil {
				return err
			}

			expr, err = ep.eh.RewriteExpression(expr)
			if err != nil {
				return err
			}

			callback(expr)
		}
		if ep.elsePos > -1 {
			_, err = w.Write([]byte("} else {\n"))
			if err != nil {
				return err
			}
			expr, err = ep.eh.RewriteExpression(strings.Join(ep.token_text[ep.elsePos+1:], ""))
			if err != nil {
				return err
			}
			callback(expr)
		}
		_, err = w.Write([]byte("}\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (ep *ExprParser) parseCondExpr(start, end int) (string, string, error) {
	c := 1
	d := end
	for i := start + 1; i <= end; i++ {
		if ep.token_text[i] == "(" {
			c += 1
		} else if ep.token_text[i] == ")" {
			c -= 1
			if c == 0 {
				d = i
				break
			}
		}
	}

	if d == end {
		return "", "", errors.New("Syntax Error: no expression found after conditional.")
	}

	e1, err := ep.eh.RewriteExpression(strings.Join(ep.token_text[start:d], ""))
	if err != nil {
		return "", "", err
	}

	e2, err := ep.eh.RewriteExpression(strings.Join(ep.token_text[d+1:end+1], ""))
	if err != nil {
		return "", "", err
	}

	return e1, e2, nil
}

func NewExprParser(eh ExprHandler) *ExprParser {
	return &ExprParser{
		eh: eh,
	}
}
