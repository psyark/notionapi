package codegen

import (
	"regexp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	titler       = cases.Title(language.Und)
	upperer      = cases.Upper(language.Und)
	getNameRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

func getName(name string) string {
	result := ""
	for _, field := range getNameRegex.Split(name, -1) {
		switch upperer.String(field) {
		case "A":
		case "ID", "URL", "PDF":
			result += upperer.String(field)
		default:
			result += titler.String(field)
		}
	}
	return result
}
