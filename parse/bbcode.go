// NOTE: BBCode parsers exist for go, I just wanted to write one myself for my own interest
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
	List *BBList `  @@`
	Tag  *BBTag  `| @@`
	Word *string `| (@Word | @Assign)`
}

type BBList struct {
	Type  *string   `"[list" ("=" @Word)? "]"`
	Items []*BBCode `( "[*]" @@ )* "[/" "list" "]"`
}

type BBTag struct {
	Name     *string `"[" @Word`
	Attrs    BBAttrs `@@ "]"`
	Body     *BBCode `@@`
	CloseTag *string `"[/" @Word "]"`
}

type BBAttrs struct {
	RawAttrs []*BBAttr `@@*`
	Attrs    map[string]string
}

type BBAttr struct {
	Key   *string `@Word "="`
	Value *string `@Word`
}

func BBCodeToMd(bbcRaw string) (string, error) {
	bbc, err := parseBBCode(bbcRaw)
	if err != nil {
		return "", err
	}
	builder := new(strings.Builder)
	bbc.ToMd(builder)
	return builder.String(), nil
}

func parseBBCode(bbcRaw string) (*BBCode, error) {
	bbcodeLexer := lexer.MustSimple([]lexer.Rule{
		{"whitespace", `[\s\r\n]+`, nil}, // (Auto-ignores white-space)

		{"ListItem", `\[\*\]`, nil},
		{"OpenListTag", `\[list`, nil},
		{"OpenEndTag", `\[/`, nil},
		{"OpenStartTag", `\[`, nil},
		{"CloseTag", `\]`, nil},
		{"Assign", `=`, nil},

		{"Word", `[^\s\[=\]]+`, nil},
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
	} else if e.Tag != nil {
		e.Tag.ToMd(b)
	} else {
		e.List.ToMd(b)
	}
}

func (l *BBList) ToMd(b *strings.Builder) {
	var listMarker string
	if l.Type != nil && *l.Type == "1" {
		listMarker = "1. "
	} else {
		listMarker = "- "
	}
	for _, e := range l.Items {
		b.WriteString(listMarker)
		e.ToMd(b)
		b.WriteString("\n")
	}
}

type complicatedMdifierFunc func(*BBTag, *strings.Builder)

func (t *BBTag) ToMd(b *strings.Builder) {
	simpleNameMap := map[string]string{
		"b":    "**",
		"i":    "_",
		"code": "`",
	}
	complicatedNameMap := map[string]complicatedMdifierFunc{
		"link": makeLink,
		"img":  makeImg,
	}
	if f, present := complicatedNameMap[*t.Name]; present {
		t.Attrs.Attrs = make(map[string]string)
		for _, a := range t.Attrs.RawAttrs {
			t.Attrs.Attrs[*a.Key] = *a.Value
		}
		f(t, b)
	} else if n, present := simpleNameMap[*t.Name]; present {
		b.WriteString(n)
		t.Body.ToMd(b)
		b.WriteString(n)
	} else {
		t.Body.ToMd(b)
	}
}

func makeLink(t *BBTag, b *strings.Builder) {
	if l, p := t.Attrs.Attrs["url"]; p {
		b.WriteString("[")
		t.Body.ToMd(b)
		b.WriteString("](")
		b.WriteString(l)
		b.WriteString(")")
	} else {
		b.WriteString("[")
		t.Body.ToMd(b)
		b.WriteString("](")
		t.Body.ToMd(b)
		b.WriteString(")")
	}
}

func makeImg(t *BBTag, b *strings.Builder) {
	b.WriteString("![")
	t.Body.ToMd(b)
	b.WriteString("](")
	t.Body.ToMd(b)
	b.WriteString(")")
}
