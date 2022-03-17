package parse

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ToMd interface {
	ToMd(*strings.Builder)
}

type BBCode struct {
	Elems []*BBElem `@@*`
}

type BBElem struct {
	Word *string `  @Word`
	Tag  *BBTag  `| @@`
}

type BBTag struct {
	Name *string `"[" @Word "]"`
	// Attrs     []*BBAttr `@@* "]"`
	Body      *BBCode `@@`
	CloseName *string `"[/" @Word "]"`
}

// type BBAttr struct {
// 	Key   *string `@Word "="`
// 	Value *string `@Word`
// }

func BBCodeToMd(bbcRaw string) (string, error) {
	bbc, err := parseBBCode(bbcRaw)
	if err != nil {
		return "", nil
	}
	builder := new(strings.Builder)
	bbc.ToMd(builder)
	return builder.String(), nil
}

func parseBBCode(bbcRaw string) (*BBCode, error) {
	bbcodeLexer := lexer.MustSimple([]lexer.Rule{
		{"whitespace", `[\s\r\n]+`, nil}, // (Auto-ignores white-space)

		{"OpenStartTag", `\[`, nil},
		{"OpenEndTag", `\[/`, nil},
		{"CloseTag", `]`, nil},
		// {"Assign", `=`, nil},
		// {"OpenTag", `\[[^\]]+\]`, nil},
		// {"CloseTag", `\[/[^\]]+\]`, nil},
		{"Word", `[^\s[\]]+`, nil},
	})
	bbcodeParser := participle.MustBuild(&BBCode{}, participle.Lexer(bbcodeLexer))

	bbc := &BBCode{}
	err := bbcodeParser.ParseString("(raw bbcode)", bbcRaw, bbc)
	if err != nil {
		return nil, err
	}

	return bbc, nil
}

func (c *BBCode) ToMd(b *strings.Builder) {
	for i, e := range c.Elems {
		if i >= 1 {
			b.WriteString(" ")
		}
		e.ToMd(b)
	}
}

func (e *BBElem) ToMd(b *strings.Builder) {
	if e.Word != nil {
		b.WriteString(*e.Word)
	} else {
		e.Tag.ToMd(b)
	}
}

func (t *BBTag) ToMd(b *strings.Builder) {
	simpleNameMap := map[string]string{
		"b":    "**",
		"i":    "_",
		"code": "`",
	}
	if n, present := simpleNameMap[*t.Name]; present {
		b.WriteString(n)
		t.Body.ToMd(b)
		b.WriteString(n)
	} else {
		t.Body.ToMd(b)
	}
	// // This would be nicer using an interface to implement a type union
	// if t.Bold != nil {
	//	b.WriteString("**")
	//	t.Bold.ToMd(b)
	//	b.WriteString("**")
	// } else if t.Italic != nil {
	//	b.WriteString("_")
	//	t.Italic.ToMd(b)
	//	b.WriteString("_")
	// } else if t.Code != nil {
	//	b.WriteString("`")
	//	t.Code.ToMd(b)
	//	b.WriteString("`")
	// }
}
