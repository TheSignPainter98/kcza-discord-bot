package parse

import (
	"regexp"
	"testing"
)

func TestBBCodeSimple(t *testing.T) {
	raw := "Hello, world!"
	bbc, err := BBCodeToMd(raw)
	if err != nil && bbc != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}
	if raw != bbc {
		t.Fatalf("Different bbcode returned from simple markdown: expected %#v but got %#v", raw, bbc)
	}
}

func TestBBCodeBold(t *testing.T) {
	raw := "[b]something bold[/b]"
	bbc, err := BBCodeToMd(raw)
	if err != nil && bbc != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}

	exp := `\*\*something bold\*\*`
	want := regexp.MustCompile(exp)
	if err == nil {
		if !want.MatchString(bbc) {
			t.Fatalf("Output does not contain input: %#v doesn't contain %#v", bbc, exp)
		}
	}
}

func TestBBCodeItalic(t *testing.T) {
	raw := "[i]something bold[/i]"
	bbc, err := BBCodeToMd(raw)
	if err != nil && bbc != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}

	exp := `^_something bold_$`
	want := regexp.MustCompile(exp)
	if err == nil {
		if !want.MatchString(bbc) {
			t.Fatalf("Output does not contain input: %#v doesn't contain %#v", bbc, exp)
		}
	}
}
