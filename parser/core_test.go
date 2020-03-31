package parser_test

import (
	. "github.com/phannam1412/go-pattern-matching/parser"
	"testing"
)

func runTests(t *testing.T, formula Expression, tests []string) {
	var res *Res
	for _, v := range tests {
		res = formula(Tokenize(v), 0)
		if res == nil {
			t.Errorf("expected success, got error, test: %s", v)
			return
		}
		if res.Value != v {
			t.Errorf("expected '%s', got '%s'", v, res.Value)
			return
		}
	}
}

type Case struct {
	input    string
	expected []string
}

func TestTokenize(t *testing.T) {

	cases := []Case{
		{
			input:    "nam123bi456",
			expected: []string{"nam", "123", "bi", "456"},
		},
		{
			input:    "nam.bi.nhan123.4.",
			expected: []string{"nam", ".", "bi", ".", "nhan", "123",".","4","."},
		},
		{
			input: "7A/43/26, ThànhThái,P14,Q10",
			expected: []string{"7","A","/","43","/","26",","," ","ThànhThái",",","P","14",",","Q","10"},
		},
	}

	for _, v := range cases {
		res := Tokenize(v.input)
		if len(res) != len(v.expected) {
			t.Errorf("expected len: %d, got len: %d", len(res), len(v.expected))
		}
		for k, v2 := range v.expected {
			if v2 != res[k] {
				t.Errorf("expected token at index %d has value: %s, got value: %s", k, v2, res[k])
			}
		}
	}
}

func TestEmail(t *testing.T) {
	runTests(t, Email, []string{
		"example@yahoo.com",
		"example123@yahoo.com",
	})
}