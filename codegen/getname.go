package codegen

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	getNameRegex  = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	nf_snake_case = &nameFormatter{"_", toLower, toLower}
	nfCamelCase   = &nameFormatter{"", toTitle, toUpper}
)

type nameFormatter struct {
	sep           string
	transform     func(string) string
	abbrTransform func(string) string
}

func toUpper(src string) string {
	return cases.Upper(language.Und).String(src)
}
func toLower(src string) string {
	return cases.Lower(language.Und).String(src)
}
func toTitle(src string) string {
	return cases.Title(language.Und).String(src)
}

func (f *nameFormatter) String(name string) string {
	fields := []string{}
	for _, field := range getNameRegex.Split(name, -1) {
		switch toUpper(field) {
		case "A":
		case "ID", "URL", "PDF":
			fields = append(fields, f.abbrTransform(field))
		default:
			fields = append(fields, f.transform(field))
		}
	}
	return strings.Join(fields, f.sep)
}
