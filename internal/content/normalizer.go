package content

import (
	"regexp"
	"strings"
	"unicode"
)

type DefaultNormalizer struct{}

func (n *DefaultNormalizer) Normalize(text string) string {
	s := removeInvisibleChars(text)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = cleanHTMLEntities(s)
	s = collapseInlineWhitespace(s)
	s = collapseBlankLines(s)
	s = trimLines(s)
	return strings.TrimSpace(s)
}

func removeInvisibleChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == '\n' || r == '\r' || r == '\t' || r == ' ':
			b.WriteRune(r)
		case unicode.IsControl(r):
			continue
		case r == '\u200B' || r == '\u200C' || r == '\u200D' || r == '\uFEFF' || r == '\u00AD':
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

var htmlEntityReplacer = strings.NewReplacer(
	"&amp;", "&",
	"&lt;", "<",
	"&gt;", ">",
	"&quot;", "\"",
	"&#39;", "'",
	"&apos;", "'",
	"&nbsp;", " ",
	"&mdash;", "—",
	"&ndash;", "–",
	"&hellip;", "...",
	"&lsquo;", "'",
	"&rsquo;", "'",
	"&ldquo;", "\"",
	"&rdquo;", "\"",
	"&bull;", "•",
)

var (
	numericEntityRe    = regexp.MustCompile(`&#\d+;`)
	namedEntityRe      = regexp.MustCompile(`&[a-zA-Z]+;`)
	inlineWhitespaceRe = regexp.MustCompile(`[^\S\n]+`)
	multiBlankLineRe   = regexp.MustCompile(`\n{3,}`)
)

func cleanHTMLEntities(s string) string {
	s = htmlEntityReplacer.Replace(s)
	s = numericEntityRe.ReplaceAllString(s, "")
	return namedEntityRe.ReplaceAllString(s, "")
}

func collapseInlineWhitespace(s string) string {
	return inlineWhitespaceRe.ReplaceAllString(s, " ")
}

func collapseBlankLines(s string) string {
	return multiBlankLineRe.ReplaceAllString(s, "\n\n")
}

func trimLines(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
