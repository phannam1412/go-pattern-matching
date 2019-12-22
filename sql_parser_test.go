package parser_test

//import (
//	. "github.com/phannam1412/go-pattern-matching"
//	"testing"
//)
//
//func testSqlParser(t *testing.T, formula Expression, tests []string) {
//	text := "INSERT INTO user(username, pass) VALUES ('phannam1412', 'testpass')"
//	SqlValue := ""
//	ColumnList := And(
//		Alphabet,
//		Any(And(Comma, Alphabet)),
//	)
//	InsertHead := And(
//		Text("INSERT INTO "),
//		Alphabet,
//		OpenParenthese,
//		ColumnList,
//		CloseParenthese,
//	)
//	InsertBody := And(
//		Text("VALUES"),
//		SomeWhitespaces,
//		OpenParenthese,
//		SqlValue,
//		And(Any(Comma, SqlValue)),
//		CloseParenthese,
//	)
//	InsertStatement := And(InsertHead, SomeWhitespaces, InsertBody)
//	formula := InsertStatement
//	res := formula(Tokenize(text), 0)
//	if res == nil {
//		t.Errorf("expected success, got error, test: %s", v)
//		return
//	}
//}
