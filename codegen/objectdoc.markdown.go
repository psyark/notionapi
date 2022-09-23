package codegen

import (
	"regexp"
	"strings"
)

// [property values](ref:property-value-object)
var markdownLinkRegex = regexp.MustCompile(`\[(.+?)\]\(.+?\)`)

func stripMarkdown(md string) string {
	md = strings.TrimPrefix(md, "# ")
	md = strings.TrimPrefix(md, "## ")
	md = strings.TrimPrefix(md, "### ")
	md = strings.ReplaceAll(md, "`", "")
	md = strings.ReplaceAll(md, "**", "")
	md = markdownLinkRegex.ReplaceAllString(md, "$1")
	return md
}
