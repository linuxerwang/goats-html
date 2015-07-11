package util

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	verifyToSnakeCase(t, "CamelCase", "camel_case")
	verifyToSnakeCase(t, "camelCase", "camel_case")
}

func TestToSnakeExpr(t *testing.T) {
	verifyToSnakeExpr(t, "test.BookShelf.CamelCase", "test.book_shelf.camel_case")
	verifyToSnakeExpr(t, "test.bookShelf.camelCase", "test.book_shelf.camel_case")
	verifyToSnakeExpr(t, "test.book_shelf.camel_case", "test.book_shelf.camel_case")
}

func TestToCamelCase(t *testing.T) {
	verifyToCamelCase(t, "camel_case", "CamelCase")
	verifyToCamelCase(t, "Camel_case", "CamelCase")
}

func TestToCamelExpr(t *testing.T) {
	verifyToCamelExpr(t, "test.book_shelf.camel_case", "test.BookShelf.CamelCase")
	verifyToCamelExpr(t, "test.Book_shelf.camel_case", "test.BookShelf.CamelCase")
	verifyToCamelExpr(t, "test.BookShelf.CamelCase", "test.BookShelf.CamelCase")
	// This is to make protocol buffer enums working.
	verifyToCamelExpr(t, "proto1.Book_SCIFI", "proto1.Book_SCIFI")
}

//----------- Private functions -----------

func verifyToSnakeCase(t *testing.T, input string, expected string) {
	actual := ToSnakeCase(input)
	if actual != expected {
		t.Errorf("ToSnakeCase(%s) != %s:  %s", input, expected, actual)
	}
}

func verifyToCamelCase(t *testing.T, input string, expected string) {
	actual := ToCamelCase(input)
	if actual != expected {
		t.Errorf("ToCamelCase(%s) != %s:  %s", input, expected, actual)
	}
}

func verifyToSnakeExpr(t *testing.T, input string, expected string) {
	actual := ToSnakeExpr(input)
	if actual != expected {
		t.Errorf("ToSnakeExpr(%s) != %s:  %s", input, expected, actual)
	}
}

func verifyToCamelExpr(t *testing.T, input string, expected string) {
	actual := ToCamelExpr(input)
	if actual != expected {
		t.Errorf("ToCamelExpr(%s) != %s:  %s", input, expected, actual)
	}
}
