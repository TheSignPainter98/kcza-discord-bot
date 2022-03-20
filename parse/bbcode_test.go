package parse

import (
	"fmt"
	"regexp"
	"strings"
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

func tagTest(t *testing.T, tag, mdTag, mdPat string) {
	testText := "something to test"
	raw := fmt.Sprintf("[%s]%s[/%s]", tag, testText, tag)
	md, err := BBCodeToMd(raw)
	if err != nil && md != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}

	if err == nil {
		want := regexp.MustCompile(mdPat + testText + mdPat)
		if !want.MatchString(md) {
			exp := mdTag + testText + mdTag
			t.Fatalf("Incorrect output: expected %#v but got %#v", exp, md)
		}
	}
}

func TestBBCodeBold(t *testing.T) {
	tagTest(t, "b", "**", `\*\*`)
}

func TestBBCodeItalic(t *testing.T) {
	tagTest(t, "i", "_", "_")
}

func TestBBCodeCode(t *testing.T) {
	tagTest(t, "code", "`", "`")
}

func TestImgSimple(t *testing.T) {
	imgLink := "https://kcza.net/img.png"
	raw := fmt.Sprintf("[img]%s[/img]", imgLink)
	md, err := BBCodeToMd(raw)
	if err != nil && md != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}

	if err == nil {
		want := regexp.MustCompile(fmt.Sprintf(`!\[%s\]\(%s\)`, imgLink, imgLink))
		if !want.MatchString(md) {
			t.Fatalf("Incorrect output: expected %#v but got %#v", want, md)
		}
	}
}

func TestLinkSimple(t *testing.T) {
	testLink := "https://kcza.net"
	raw := fmt.Sprintf("[link]%s[/link]", testLink)
	md, err := BBCodeToMd(raw)
	if err != nil && md != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}

	if err == nil {
		want := regexp.MustCompile(fmt.Sprintf(`\[%s\]\(%s\)`, testLink, testLink))
		if !want.MatchString(md) {
			t.Fatalf("Incorrect output: expected %#v but got %#v", want, md)
		}
	}
}

func TestLinkWithUrl(t *testing.T) {
	testLink := "https://kcza.net"
	testText := "some text to test"
	raw := fmt.Sprintf("[link url=%s]%s[/link]", testLink, testText)
	md, err := BBCodeToMd(raw)
	if err != nil && md != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}

	if err == nil {
		want := regexp.MustCompile(fmt.Sprintf(`\[%s\]\(%s\)`, testText, testLink))
		if !want.MatchString(md) {
			t.Fatalf("Incorrect output: expected %#v but got %#v", want, md)
		}
	}
}

func TestListNumeric(t *testing.T) {
	testListType(t, "1", `[0-9]+\.`)
}

func TestListBullet(t *testing.T) {
	testListType(t, "", `[-*]`)
}

func TestListUnsupported(t *testing.T) {
	unsupportedPat := `([0-9]+\.|[-*])`
	testListType(t, "a", unsupportedPat)
	testListType(t, "A", unsupportedPat)
	testListType(t, "i", unsupportedPat)
	testListType(t, "I", unsupportedPat)
}

func testListType(t *testing.T, listType, mdItemPat string) {
	hdrText := ""
	if len(listType) > 0 {
		hdrText = "=" + listType
	}
	// lines := []string{"Hello", "World!", "How", "are", "you?", "", ""}
	lines := []string{"Hello", "World!", "How", "are", "you?"}
	bbcListBody := make([]string, len(lines))
	for i, l := range lines {
		bbcListBody[i] = "[*]" + l
	}
	raw := fmt.Sprintf("[list%s]%s[/list]", hdrText, strings.Join(bbcListBody, ""))
	md, err := BBCodeToMd(raw)
	if err != nil && md != "" {
		t.Fatalf("Expected either non-nil error or non-empty string and non-empty input")
	}

	if err == nil {
		mdListBody := make([]string, len(lines))
		for i, l := range lines {
			mdListBody[i] = mdItemPat + " " + l
		}
		want := regexp.MustCompile(strings.Join(mdListBody, "\n"))
		if !want.MatchString(md) {
			t.Fatalf("Incorrect output: expected %#v but got %#v", want, md)
		}
	} else {
		t.Fatalf("Should have been able to parse %#v, got %#v: %s", raw, md, err)
	}
}
