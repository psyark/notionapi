package codegen

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	titler  = cases.Title(language.Und)
	upperer = cases.Upper(language.Und)
)

func getName(name string) string {
	if name == "" {
		return ""
	}

	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	fields := []string{}
	for _, field := range strings.Fields(name) {
		switch upperer.String(field) {
		case "ID", "URL", "PDF":
			field = upperer.String(field)
		default:
			field = titler.String(field)
		}
		fields = append(fields, field)
	}
	return strings.Join(fields, "")
}

func getMethodName(title string) string {
	fields := strings.Fields(title)
	for i, field := range fields {
		if field == "a" {
			fields[i] = ""
		} else {
			fields[i] = strings.Title(field)
		}
	}
	return strings.Join(fields, "")
}
