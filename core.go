package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func PrintJson(test interface{}) {
	tmp, _ := json.MarshalIndent(test, "", "    ")
	fmt.Printf("%s\n", tmp)
}

func JsonEncode(test interface{}) string {
	tmp, _ := json.Marshal(test)
	return string(tmp)
}

func isset(theMap map[string]string, key string) bool {
	_, ok := theMap[key]
	return ok
}




const separator = "+ _)(*&^%#@!~[];/.,\\{}|:\"<>?`-='\n\r\t$"
var Whitespace Expression
var SomeWhitespaces Expression
var AnyWhitespaces Expression
var Comma Expression
var Dollar Expression
var Hyphen Expression
var Dot Expression
var At Expression
var Equal Expression
var Colon Expression
var Tab Expression
var Plus Expression
var Backsplash Expression
var NumberOrAlphabet Expression
var NewLine Expression

func init() {

	Plus = Text("+")
	Comma = Text(",")
	Dollar = Text("$")
	Whitespace = Text(" ")
	Hyphen = Text("-")
	Dot = Text(".")
	At = Text("@")
	Equal = Text("=")
	Colon = Text(":")
	Tab = Text("\t")
	NewLine = Text("\n")
	Backsplash = Text("/")
	NumberOrAlphabet = Or(Number, Alphabet)
	SomeWhitespaces = Label("some whitespaces", Some(Whitespace))
	AnyWhitespaces = Label("any whitespaces", Any(Whitespace))
}

func Alphabet(tokens []string, pos int) *Res {
	res := strings.Contains(separator, tokens[pos])
	if res == true {
		return nil
	}
	return &Res{
		Pos: pos + 1,
		Expr: "alphabet",
		Value: tokens[pos],
	}
}

func Tokenize(text string) []string {
	var output []string
	currentWord := ""
	separators := "+ _)(*&^%#@!~[];/.,\\{}|:\"<>?`-='\n\r\t$"
	var prevChar uint8
	for a := 0; a < len(text); a++ {
		currentChar := text[a : a + 1]

		// separator
		if strings.Contains(separators, currentChar) {
			if len(currentWord) > 0 {
				output = append(output, currentWord)
			}
			output = append(output, currentChar)
			currentWord = ""
			prevChar = currentChar[0]
			continue
		}

		// transition from alphabet -> number and number -> alphabet
		if (isNumber(prevChar) && !isNumber(currentChar[0])) || (!isNumber(prevChar) && isNumber(currentChar[0])) {
			if len(currentWord) > 0 {
				output = append(output, currentWord)
			}
			currentWord = currentChar
			prevChar = currentChar[0]
			continue
		}

		currentWord += currentChar
		prevChar = currentChar[0]
	}
	if len(currentWord) > 0 {
		output = append(output, currentWord)
	}
	return output
}

type Expression func(tokens []string, pos int) *Res

type Res struct {
	Pos int `json:"-",omitempty"`
	Value string `json:",omitempty"`
	Expr string `json:",omitempty"`
	Children []Res `json:",omitempty"`
	Data map[string]string `json:",omitempty"`
}

func And(expressions ...Expression) Expression {
	return func(tokens []string, pos int) *Res {
		if pos >= len(tokens) {
			return nil
		}
		var children []Res
		var value []string
		for _, exp := range expressions {
			res := exp(tokens, pos)
			if res == nil {
				return nil
			}
			children = append(children, *res)
			pos = res.Pos
			value = append(value, res.Value)
		}
		return &Res{
			Pos: pos,
			Expr: "and",
			Children: children,
			Value: strings.Join(value, ""),
		}
	}
}

func AndBut(should Expression, shouldNot ...Expression) Expression {
	return func(tokens []string, pos int) *Res {
		if pos >= len(tokens) {
			return nil
		}
		first := should(tokens, pos)
		if first == nil {
			return nil
		}
		for _, v := range shouldNot {
			second := v(tokens, pos)
			if second != nil{
				return nil
			}
		}
		return first
	}
}

// Match all posibilities of all positions in the token list.
func FullSearch(expr Expression, limit int) Expression {
	return func(tokens []string, pos int) *Res {
		var children []Res
		var value []string

		// Travel through every possible position of token list to find as much matches as possible.
		for a := 0; a < len(tokens) && pos < len(tokens); a++ {

			// Is there any match at this position ?
			res := expr(tokens, pos)

			// No match ? Find at next position.
			if res == nil {
				pos++
				continue
			}
			pos = res.Pos
			children = append(children, *res)
			value = append(value, res.Value)
			if len(value) == limit && limit != -1 {
				break
			}
		}

		if value == nil {
			return nil
		}

		return &Res{
			Pos: pos,
			Expr: "full_search",
			Children: children,
			Value: strings.Join(value, ""),
		}
	}
}

func LookupForOne(expr Expression) Expression {
	return func(tokens []string, pos int) *Res {
		for a := 0; a < len(tokens) && pos < len(tokens); a++ {
			res := expr(tokens, pos)
			if res != nil {
				return &Res{
					Pos: res.Pos,
					Expr: "lookup_for_one",
					Value: res.Value,
				}
			}
			pos++
		}
		return nil
	}
}

func Some(expr Expression) Expression {
	return SomeWithMin(expr, 1)
}

func SomeWithMin(expr Expression, min int) Expression {
	return func(tokens []string, pos int) *Res {
		var value []string
		for a := 0; a < len(tokens) && pos < len(tokens); a++ {
			res := expr(tokens, pos)
			
			// Cannot match anymore ?
			if res == nil {
				if len(value) < min {
					return nil
				}
				return &Res{
					Pos: pos,
					Expr: "some",
					Value: strings.Join(value, ""),
				}
			}

			// Continue matching
			value = append(value, res.Value)
			pos = res.Pos
		}
		if len(value) < min {
			return nil
		}
		return &Res{
			Pos: pos,
			Expr: "some",
			Value: strings.Join(value, ""),
		}
	}
}

// Match from 0 -> as much as possible
// Match all posibilities starting at the current pos, limit by min.
func Any(expr Expression) Expression {
	return func(tokens []string, pos int) *Res {
		var values []string
		var children []Res
		a := 0
		for a = 0; a < len(tokens) && pos < len(tokens); a++ {
			res := expr(tokens, pos)

			// Cannot match anymore ?
			if res == nil {
				return &Res{
					Pos: pos,
					Expr: "any",
					Children: children,
					Value: strings.Join(values, ""),
				}
			}

			values = append(values, res.Value)
			children = append(children, *res)
			pos = res.Pos
		}

		return &Res{
			Pos: pos,
			Expr: "any",
			Children: children,
			Value: strings.Join(values, ""),
		}
	}
}

func Text(str string) Expression {
	thisTokens := Tokenize(str)
	return func(tokens []string, pos int) *Res {
		if pos + len(thisTokens) > len(tokens) {
			return nil
		}
		for a := 0; a < len(thisTokens); a++ {
			if thisTokens[a] != tokens[pos + a] {
				return nil
			}
		}
		return &Res{
			Pos: pos + len(thisTokens),
			Expr: "text",
			Value: str,
		}
	}
}

func CaseInsensitive(token string) Expression {
	tokensForMatch := Tokenize(token)
	return func(tokens []string, pos int) *Res {
		if pos + len(tokensForMatch) > len(tokens) {
			return nil
		}
		var matched []string
		for a := 0; a < len(tokensForMatch); a++ {
			if strings.ToLower(tokensForMatch[a]) != strings.ToLower(tokens[pos + a]) {
				return nil
			}
			matched = append(matched, tokens[pos + a])
		}
		return &Res{
			Pos: pos + len(tokensForMatch),
			Expr: "case_insensitive",
			Value: strings.Join(matched, ""),
		}
	}
}

func Or(expressions ...Expression) Expression {
	return func(tokens []string, pos int) *Res {
		if pos >= len(tokens) {
			return nil
		}
		for _, exp := range expressions {
			res := exp(tokens, pos)
			if res != nil {
				return &Res{
					Pos: res.Pos,
					Expr: "or",
					Children: []Res{*res},
					Value: res.Value,
				}
			}
		}
		return nil
	}
}

func TextUntilEndAt(matchingForEnd Expression) Expression {
	return func(tokens []string, pos int) *Res {
		var result []string
		a := 0
		for a = 0; a + pos < len(tokens); a++ {
			if matchingForEnd(tokens, pos + a) != nil {
				break
			}
			result = append(result, tokens[pos + a])
		}
		if len(result) == 0 {
			return nil
		}
		return &Res{
			Pos: pos + a,
			Expr: "text_until_end_at",
			Value: strings.Join(result, ""),
		}
	}
}

func TextUntilLineEnd(tokens []string, pos int) *Res {
	tmp := TextUntilEndAt(Text("\n"))
	res := tmp(tokens, pos)
	if res == nil {
		return nil
	}
	res.Expr = "text_until_line_end"
	return res
}

func OneTokenExceptLineBreak(tokens []string, pos int) *Res {
	if tokens[pos] == "\n" {
		return nil
	}
	return &Res{
		Pos: pos + 1,
		Expr: "one_token_except_line_break",
		Value: tokens[pos],
	}
}

func PairSeparator(tokens []string, pos int) *Res {
	main := Or(
		And(Any(Or(Whitespace,Tab)), Colon, Any(Or(Whitespace, Tab))),
		SomeWithMin(Or(Whitespace, Tab), 3),
	)
	res := main(tokens, pos)
	if res == nil {
		return nil
	}
	return &Res{
		Pos: res.Pos,
		Expr: "pair_separator",
		Value: res.Value,
	}
}

func Token(tokens []string, pos int) *Res {
	return &Res{
		Pos: pos + 1,
		Expr: "token",
		Value: tokens[pos],
	}
}

func SingleAlphabet(tokens []string, pos int) *Res {
	if len(tokens[pos]) > 1 {
		return nil
	}
	char := tokens[pos][0]
	if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
		return &Res{
			Pos: pos + 1,
			Expr: "token",
			Value: tokens[pos],
		}
	}
	return nil
}

func Ucfirst(tokens []string, pos int) *Res {
	if pos >= len(tokens) {
		return nil
	}
	if len(tokens[pos]) == 0 {
		return nil
	}
	if tokens[pos][0] < 'A' || tokens[pos][0] > 'Z' {
		return nil
	}
	return &Res{
		Pos: pos + 1,
		Expr: "token",
		Value: tokens[pos],
	}
}

func AllUppercases(tokens []string, pos int) *Res {
	if pos >= len(tokens) {
		return nil
	}
	if strings.ToUpper(tokens[pos]) == tokens[pos] {
		return &Res{
			Pos:   pos + 1,
			Expr:  "token",
			Value: tokens[pos],
		}
	}
	return nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isNumber(char uint8) bool {
	return char >= '0' && char <= '9';
}

func Number(tokens []string, pos int) *Res {
	if !isNumeric(tokens[pos]) {
		return nil
	}
	return &Res{
		Pos: pos + 1,
		Expr: "token",
		Value: tokens[pos],
	}
}

func LengthAtleast(expression Expression, minLength int) Expression {
	return func(tokens []string, pos int) *Res {
		res := expression(tokens, pos)
		if res == nil {
			return nil
		}
		if len(res.Value) < minLength {
			return nil
		}
		return res
	}
}

func NotToken(token string) Expression {
	return func(tokens []string, pos int) *Res {
		if pos >= len(tokens) {
			return nil
		}
		if tokens[pos] == token {
			return nil
		}
		return &Res{
			Pos: pos + 1,
			Expr: "not_token",
			Value: tokens[pos],
		}
	}
}

func Email(tokens []string, pos int) *Res {
	formula := And(
		Some(NumberOrAlphabet),
		Any(And(Dot, Some(NumberOrAlphabet))),
		Text("@"),
		Some(NumberOrAlphabet),
		Text("."),
		Alphabet,
	)
	return formula(tokens, pos)
}

func Label(label string, expr Expression) Expression {
	return func(tokens []string, pos int) *Res {
		res := expr(tokens, pos)
		if res == nil {
			return nil
		}
		res.Expr = label
		return res
	}
}