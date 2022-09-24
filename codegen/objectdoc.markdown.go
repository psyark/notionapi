package codegen

import (
	"regexp"
	"strings"
)

var (
	markdownLinkRegex    = regexp.MustCompile(`\[(.+?)\]\(.+?\)`) // [property values](ref:property-value-object)
	markdownEmRegex1     = regexp.MustCompile(`\*(.+?)\*`)
	markdownEmRegex2     = regexp.MustCompile(`(^|\W)_(.+?)_(\W|$)`)
	markdownStrongRegex1 = regexp.MustCompile(`\*\*(.+?)\*\*`)
	markdownStrongRegex2 = regexp.MustCompile(`(^|\W)__(.+?)__(\W|$)`)
)

func stripMarkdown(md string) string {
	md = strings.TrimPrefix(md, "# ")
	md = strings.TrimPrefix(md, "## ")
	md = strings.TrimPrefix(md, "### ")
	md = strings.ReplaceAll(md, "`", "")
	md = markdownLinkRegex.ReplaceAllString(md, "$1")
	md = markdownStrongRegex1.ReplaceAllString(md, "$1")
	md = markdownStrongRegex2.ReplaceAllString(md, "$1$2$3")
	md = markdownEmRegex1.ReplaceAllString(md, "$1")
	md = markdownEmRegex2.ReplaceAllString(md, "$1$2$3")
	return md
}
