EBNF-inspired pattern matching written in golang

Following are some application for pattern matching:

### Wrap all strings for translate

input:
```
function sayHello() {
	console.log("hello world");
	console.log('how are you doing today');
}
```

output:
```
function sayHello() {
	console.log(_t("hello world"));
	console.log(_t('how are you doing today'));
}
```

code: 
```
package main

import (
	. "github.com/phannam1412/go-pattern-matching"
	"strings"
)

func main() {
	input := `
		function sayHello() {
			console.log("hello \"world\"");
			console.log('how \'are\' you "doing" today');
		}
	`
	doubleQuote := And(
		Text(`"`),
		Any(Or(Text(`\"`),NotToken(`"`))),
		Text(`"`),
	)
	singleQuote := And(
		Text(`'`),
		Any(Or(Text(`\'`),NotToken(`'`))),
		Text(`'`),
	)
	formula := FullSearch(Or(doubleQuote, singleQuote), -1)

	tokens := Tokenize(input)

	parsed := formula(tokens, 0)

	output := input
	for _, v := range parsed.Children {
		output = strings.ReplaceAll(output, v.Value, "_t(" + v.Value + ")")
	}

	println(output)

}
```

