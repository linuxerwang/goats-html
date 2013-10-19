package goats

import (
	"bytes"
	"go/ast"
	"go/parser"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func StartsWithLower(input string) bool {
	if len(input) == 0 {
		return false
	}
	r, _ := utf8.DecodeRuneInString(input)
	return unicode.IsLower(r)
}

func IsSnakeCase(input string) bool {
	return StartsWithLower(input) && strings.Contains(input, "_")
}

func ToSnakeCase(camelCase string) string {
	var buffer bytes.Buffer
	for idx, ch := range camelCase {
		if idx > 0 && unicode.IsUpper(ch) {
			buffer.WriteRune('_')
		}
		buffer.WriteRune(unicode.ToLower(ch))
	}
	return buffer.String()
}

func ToSnakeExpr(camelExpr string) string {
	var buffer bytes.Buffer
	for idx, ch := range camelExpr {
		if unicode.IsUpper(ch) {
			if idx > 0 && camelExpr[idx-1] != '.' {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(ch))
		} else {
			buffer.WriteRune(ch)
		}
	}
	return buffer.String()
}

func ToCamelCase(snakeCase string) string {
	var buffer bytes.Buffer
	var afterUnderScore bool
	for idx, ch := range snakeCase {
		if idx == 0 {
			buffer.WriteRune(unicode.ToUpper(ch))
		} else if ch == '_' {
			afterUnderScore = true
		} else {
			if afterUnderScore && unicode.IsLower(ch) {
				buffer.WriteRune(unicode.ToUpper(ch))
				afterUnderScore = false
			} else {
				buffer.WriteRune(ch)
			}
		}
	}
	return buffer.String()
}

func ToCamelExpr(snakeExpr string) string {
	var buffer bytes.Buffer
	var afterUnderScore bool
	for idx, ch := range snakeExpr {
		if ch == '_' {
			afterUnderScore = true
		} else {
			if afterUnderScore {
				if !unicode.IsLower(ch) {
					buffer.WriteRune('_')
				}
				buffer.WriteRune(unicode.ToUpper(ch))
				afterUnderScore = false
			} else if idx > 0 && snakeExpr[idx-1] == '.' {
				buffer.WriteRune(unicode.ToUpper(ch))
			} else {
				buffer.WriteRune(ch)
			}
		}
	}
	return buffer.String()
}

func ToHiddenName(name string) string {
	var buffer bytes.Buffer
	for idx, ch := range name {
		if idx == 0 {
			buffer.WriteRune(unicode.ToLower(ch))
		} else {
			buffer.WriteRune(ch)
		}
	}
	return buffer.String()
}

func ToPublicName(name string) string {
	var buffer bytes.Buffer
	for idx, ch := range name {
		if idx == 0 {
			buffer.WriteRune(unicode.ToUpper(ch))
		} else {
			buffer.WriteRune(ch)
		}
	}
	return buffer.String()
}

func TrimWhiteSpaces(input string) string {
	return strings.Trim(input, " \t\n")
}

func ToGoString(input string) string {
	return strings.Replace(input, "'", "\"", -1)
}

func NormalizeText(input string) string {
	text := strings.Replace(input, "\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Replace(text, "\t", " ", -1)
	reg := regexp.MustCompile(`\s+`)
	text = reg.ReplaceAllString(text, " ")
	return TrimWhiteSpaces(text)
}

func HasLeadingSpace(input string) bool {
	ch := input[0]
	return ch == ' ' || ch == '\n' || ch == '\t'
}

func HasTrailingSpace(input string) bool {
	ch := input[len(input)-1]
	return ch == ' ' || ch == '\n' || ch == '\t'
}

func ExtractSelector(input string) (string, bool) {
	// Create the AST by parsing src.
	f, err := parser.ParseExpr(input)
	if err != nil {
		panic(err)
	}

	selector := ""
	found := false

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			selector = flatten(x)
			if strings.Index(selector, ".") > -1 {
				found = true
			}
			return false
		}
		return true
	})

	return selector, found
}

func flatten(s *ast.SelectorExpr) string {
	var buffer bytes.Buffer
	switch x := s.X.(type) {
	case *ast.SelectorExpr:
		buffer.WriteString(flatten(x))
	case *ast.Ident:
		buffer.WriteString(x.Name)
	}
	return buffer.String() + "." + s.Sel.Name
}

func SplitVarDef(varDef string) (string, string) {
	idx := strings.Index(varDef, ":")
	return varDef[:idx], varDef[idx+1:]
}
