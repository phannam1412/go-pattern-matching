package parser_test

//import (
//	. "github.com/phannam1412/go-pattern-matching"
//	"testing"
//)
//
//func testSqlParser(t *testing.T, formula Expression, tests []string) {
//	text := "INSERT INTO user(username, pass) VALUES ('phannam1412', 'testpass')"
//	SqlValue := ""
//	ColumnList := Combine(
//		Alphabet,
//		Any(Combine(Comma, Alphabet)),
//	)
//	InsertHead := Combine(
//		Text("INSERT INTO "),
//		Alphabet,
//		OpenParenthese,
//		ColumnList,
//		CloseParenthese,
//	)
//	InsertBody := Combine(
//		Text("VALUES"),
//		SomeWhitespaces,
//		OpenParenthese,
//		SqlValue,
//		Combine(Any(Comma, SqlValue)),
//		CloseParenthese,
//	)
//	InsertStatement := Combine(InsertHead, SomeWhitespaces, InsertBody)
//	formula := InsertStatement
//	res := formula(Tokenize(text), 0)
//	if res == nil {
//		t.Errorf("expected success, got error, test: %s", v)
//		return
//	}
//}
