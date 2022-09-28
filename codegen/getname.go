package codegen

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	upperer = cases.Upper(language.Und)
	lowerer = cases.Lower(language.Und)
	titler  = cases.Title(language.Und)

	getNameRegex  = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	nf_snake_case = &nameFormatter{"_", lowerer, lowerer}
	nfCamelCase   = &nameFormatter{"", titler, upperer}
)

type nameFormatter struct {
	sep       string
	caser     cases.Caser
	abbrCaser cases.Caser
}

func (f *nameFormatter) String(name string) string {
	fields := []string{}
	for _, field := range getNameRegex.Split(name, -1) {
		switch upperer.String(field) {
		case "A":
		case "ID", "URL", "PDF":
			fields = append(fields, f.abbrCaser.String(field))
		default:
			fields = append(fields, f.caser.String(field))
		}
	}
	return strings.Join(fields, f.sep)
}
