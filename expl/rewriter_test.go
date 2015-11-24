package expl

import (
	"runtime/debug"
	"testing"

	"github.com/linuxerwang/goats-html/symbolmgr"
)

func TestExprRewriterGo(t *testing.T) {

	er := NewExprRewriter(createSmgr(), "go")

	src := "'text-' + *files.group_a.input_file.file_name[2:5] + " +
		"unixdate(*files.group_a.input_file.file_size*-5) + " +
		"a.b.c.Category_Static + a.b.c.date_time"
	expected := "\"text-\"+*files.GroupA.InputFile.FileName[2:5]+" +
		"__impl.UnixDate(*files.GroupA.InputFile.FileSize*-5)+" +
		"a.b.c.Category_Static+a.B.C.DateTime"
	testExprRewriter(t, er, src, expected)
}

func TestExprRewriterClosure(t *testing.T) {
	er := NewExprRewriter(createSmgr(), "closure")

	src := "'text-' + *files.group_a.input_file.file_name[2:5] + " +
		"unixdate(*files.group_a.input_file.file_size * -5) + " +
		"a.b.c.Category_Static + a.b.c.date_time"
	expected := "\"text-\"+files.getGroupA().getInputFile().getFileName()[2:5]+" +
		"goats.runtime.filters.unixdate(files.getGroupA().getInputFile().getFileSize()*-5)+" +
		"x.y.a.b.c.Category_Static+x.y.a.getB().getC().getDateTime()"
	testExprRewriter(t, er, src, expected)
}

func createSmgr() *symbolmgr.SymbolMgr {
	builder := symbolmgr.NewSymbolMgrBuilder()
	builder.AddArg("files", true /* isPb */)
	builder.AddImport("a", true /* isPb */, "x/y/a", "x.y.a")
	return builder.Build()
}

func testExprRewriter(t *testing.T, er ExprHandler, src, expected string) {
	res, err := er.RewriteExpression(src)
	if err != nil {
		t.Errorf("Expr rewrite failed with error: %s", err.Error())
		return
	}

	if expected != res {
		debug.PrintStack()
		t.Errorf("wrong rendered text: %s,\nexpected: %s", res, expected)
	}
}
