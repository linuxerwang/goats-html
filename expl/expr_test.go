package expl

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"testing"
)

func TestExprParser_0(t *testing.T) {
	src := `format("hello, %s", a.Name) + b.CallIt() + a if (a > 5) elif (a > 2) format("world, %s", a.Name) + a
		elif (a + c > 100) "exit" else "!"`
	expected := `if (a>5) {
attrVar = format("hello, %s",a.Name)+b.CallIt()+a
} else if a>2 {
attrVar = format("world, %s",a.Name)+a
} else if a+c>100 {
attrVar = "exit"
} else {
attrVar = "!"
}

`
	testExprParser(t, src, expected)
}

func TestExprParser_1(t *testing.T) {
	src := `"hello" if ((a + b) > 100) elif (a > 50) ("world" + " war") elif (a > 30) "haha" else "done"`
	expected := `if ((a+b)>100) {
attrVar = "hello"
} else if a>50 {
attrVar = ("world"+" war")
} else if a>30 {
attrVar = "haha"
} else {
attrVar = "done"
}

`
	testExprParser(t, src, expected)
}

func TestExprParser_2(t *testing.T) {
	src := `format("hello, %s", a.Name) + b.CallIt() + a`
	expected := `attrVar = format("hello, %s",a.Name)+b.CallIt()+a

`

	testExprParser(t, src, expected)
}

func TestExprParser_3(t *testing.T) {
	src := `format("hello, %s", a.Name) + b.CallIt() + a if (a > 30)`
	expected := `if (a>30) {
attrVar = format("hello, %s",a.Name)+b.CallIt()+a
}

`
	testExprParser(t, src, expected)
}

func TestExprParser_4(t *testing.T) {
	src := `format("hello, %s", a.Name) + b.CallIt() + a if (a > 30) else "great" + a.Name`
	expected := `if (a>30) {
attrVar = format("hello, %s",a.Name)+b.CallIt()+a
} else {
attrVar = "great"+a.Name
}

`
	testExprParser(t, src, expected)
}

func TestExprParser_5(t *testing.T) {
	src := `format("hello, %s", a.Name) + b.CallIt() + a if (a > 30) elif (a > 10) "great" + a.Name`
	expected := `if (a>30) {
attrVar = format("hello, %s",a.Name)+b.CallIt()+a
} else if a>10 {
attrVar = "great"+a.Name
}

`
	testExprParser(t, src, expected)
}

type fakeExprHandler struct{}

func (fe *fakeExprHandler) RewriteExpression(originalExpr string) string {
	return originalExpr
}

func testExprParser(t *testing.T, src, expected string) {
	var buf bytes.Buffer

	p := NewExprParser(new(fakeExprHandler))

	p.Evaluate(src, &buf, func(expr string) {
		(&buf).WriteString(fmt.Sprintf("attrVar = %s\n", expr))
	})

	if expected != buf.String() {
		debug.PrintStack()
		t.Errorf("wrong rendered text: %s,\nexpected: %s", buf.String(), expected)
	}
}
