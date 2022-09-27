package codegen

import (
	"strings"
)

// TODO prop_nameとPropNameの切り替えサポート
func getName(name string) string {
	if name == "" {
		return ""
	}

	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	fields := []string{}
	for _, field := range strings.Fields(name) {
		switch field {
		case "id", "url":
			field = strings.ToUpper(field)
		default:
			field = strings.Title(field)
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
