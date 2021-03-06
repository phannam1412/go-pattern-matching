# EBNF-inspired pattern matching written in golang

Following are some applications for pattern matching:

### Wrap all strings for translate

**input**
```
function sayHello() {
    console.log("hello \"world\"");
    console.log('how \'are\' you "doing" today');
}
```

**output**
```
function sayHello() {
    console.log(_t("hello \"world\""));
    console.log(_t('how \'are\' you "doing" today'));
}
```

**Full code** 
```
package main

import (
	. "github.com/phannam1412/go-pattern-matching/parser"
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

### Extract url

**input**
```$xslt
extract some url from https://example.com 
or http://www.w3school.com
or with hyphen http://www.w3-school.com
``` 

**output**
```$xslt
https://example.com
http://www.w3school.com
http://www.w3-school.com
```

**formula**
```$xslt
domainName := And(
    Alphabet,
    Any(And(Hyphen, Alphabet)),
)
formula := FullSearch(And(
    Or(Text("http"),Text("https")),
    Text("://"),
    Any(Text("www.")),
    domainName,
    Text("."),
    Alphabet,
), -1)
```

### Extract emails from a bulk of text 

**input**
```$xslt
extract some emails like contact-me-me@yahoo.com, dev.to@gmail.com, 
and exclude invalid emails like ^hello@wer.co 
```

**output**
```$xslt
contact-me-me@yahoo.com
to@gmail.com
hello@wer.co
```

**formula**
```$xslt
name := And(Alphabet, Any(And(Hyphen, Alphabet)))
formula := FullSearch(And(name,At,Alphabet,Dot,Alphabet,), -1)
```