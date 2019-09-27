package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

func PrintJson(test interface{}) {
	tmp, _ := json.MarshalIndent(test, "", "    ")
	fmt.Printf("%s\n", tmp)
}

func JsonEncode(test interface{}) string {
	tmp, _ := json.MarshalIndent(test, "", "    ")
	return string(tmp)
}

func isset(theMap map[string]string, key string) bool {
	_, ok := theMap[key]
	return ok
}




const separator = "+ _)(*&^%#@!~[];/.,\\{}|:\"<>?`-='\n\r\t$"
var tries Expression
var mapping map[string]string
var Whitespace Expression
var SomeWhitespaces Expression
var AnyWhitespaces Expression
var Comma Expression
var Dollar Expression
var Hyphen Expression
var Dot Expression
var At Expression

func init() {

	Comma = Text(",")
	Dollar = Text("$")
	Whitespace = Text(" ")
	Hyphen = Text("-")
	Dot = Text(".")
	At = Text("@")
	SomeWhitespaces = Some(Whitespace)
	AnyWhitespaces = Any(Whitespace)

	mapping = map[string]string {
		"day of birth": "birthday",
		"date of birth": "birthday",
		"birthday": "birthday",

		"e-mail address": "email",
		"e-mail": "email",
		"email": "email",

		"home address": "address",
		"contact address": "address",
		"address": "address",

		"height": "height",
		"weight": "weight",

		"full name": "fullname",

		"phone number": "phone",
		"hand phone": "phone",
		"contact phone": "phone",
		"tel": "phone",
		"phone": "phone",
		"mobile": "phone",

		"sex": "sex",
		"marital": "marital",
		"dob": "birthday",
		"position": "position",
		"job level": "job_level",
		"work place": "work_place",
		"job category": "job_category",
		"salary expected": "salary",
		"salary": "salary",


		"years of experience": "experience",
		"experience": "experience",

		"interests": "interest",


		"highest degree": "highest_degree",
		"language proficiency": "language_proficiency",
		"most recent job": "most_recent_job",
		"most recent company": "most_recent_company",
		"current job level": "current_job_level",
		"first name": "first_name",
		"last name": "last_name",
		"gender": "gender",
		"nationality": "nationality",
		"website": "website",
		"web site": "website",
		"birthplace": "birthplace",
		"place of birth": "birthplace",
		"languages": "language",
		"language": "language",
		"yahoo": "yahoo",
		"skype": "skype",
	}

	var tmp []Expression
	for wordForDetect := range mapping {
		tmp = append(tmp, CaseInsensitive(Tokenize(wordForDetect)))
	}
	tries = Or(tmp...)
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
	var res []string
	word := ""
	chars := "+ _)(*&^%#@!~[];/.,\\{}|:\"<>?`-='\n\r\t$"
	for a := 0; a < len(text); a++ {
		char := text[a : a + 1]
		if !strings.Contains(chars, char) {
			word += char
		} else {
			if len(word) > 0 {
				res = append(res, word)
			}
			res = append(res, char)
			word = ""
		}
	}
	if len(word) > 0 {
		res = append(res, word)
	}
	return res
}

type Expression func(tokens []string, pos int) *Res

type Res struct {
	Pos int
	Value string
	Expr string
	Children []Res
	Data map[string]string
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

func AtLeast(expr Expression, count int) Expression {
	return func(tokens []string, pos int) *Res {
		var value []string
		for a := 0; a < len(tokens) && pos < len(tokens); a++ {
			res := expr(tokens, pos)
			if res != nil {
				return nil
			}
			value = append(value, res.Value)
			if a >= count {
				return &Res{
					Pos: res.Pos,
					Expr: "at_least",
					Value: strings.Join(value, ""),
				}
			}
			pos++
		}
		return nil
	}
}

func Some(expr Expression) Expression {
	return func(tokens []string, pos int) *Res {
		var value []string
		for a := 0; a < len(tokens) && pos < len(tokens); a++ {
			res := expr(tokens, pos)
			if res == nil {
				if len(value) > 0 {
					return &Res{
						Pos: pos,
						Expr: "at_least",
						Value: strings.Join(value, ""),
					}
				}
				return nil
			}
			value = append(value, res.Value)
			pos++
		}
		if len(value) > 0 {
			return &Res{
				Pos: pos,
				Expr: "at_least",
				Value: strings.Join(value, ""),
			}
		}
		return nil
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
			Expr: "greedy",
			Children: children,
			Value: strings.Join(values, ""),
		}
	}
}

// Match all posibilities starting at the current pos, limit by min.
func Greedy(expr Expression, min int) Expression {
	return func(tokens []string, pos int) *Res {
		var values []string
		var children []Res
		a := 0
		for a = 0; a < len(tokens) && pos < len(tokens); a++ {
			res := expr(tokens, pos)

			// Cannot match anymore ?
			if res == nil {
				if a < min {
					return nil
				}
				return &Res{
					Pos: pos,
					Expr: "greedy",
					Children: children,
					Value: strings.Join(values, ""),
				}
			}

			values = append(values, res.Value)
			children = append(children, *res)
			pos++
		}

		// Match until the end.
		if a < min {
			return nil
		}

		return &Res{
			Pos: pos,
			Expr: "greedy",
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

func CaseInsensitive(tokensForMatch []string) Expression {
	return func(tokens []string, pos int) *Res {
		if pos + len(tokensForMatch) > len(tokens) {
			return nil
		}
		for a := 0; a < len(tokensForMatch); a++ {
			if !strings.EqualFold(tokensForMatch[a], tokens[pos + a]) {
				return nil
			}
		}
		return &Res{
			Pos: pos + len(tokensForMatch),
			Expr: "case_insensitive",
			Value: strings.Join(tokensForMatch, ""),
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
		Greedy(Or(Whitespace, Text("\t"), Text(":")), 2),
		Text(":"),
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
	formula := And(Token, Text("@"), Token, Text("."), Token)
	return formula(tokens, pos)
}

func Keyword(tokens []string, pos int) *Res {
	res := tries(tokens, pos)
	if res == nil {
		return nil
	}
	keyword := res.Value
	candidateField := mapping[keyword]
	return &Res{
		Pos: res.Pos,
		Expr: "keyword",
		Value: candidateField,
	}
}

func OnePair(tokens []string, pos int) *Res {
	main := And(Keyword, PairSeparator, TextUntilLineEnd)
	return main(tokens, pos)
}
