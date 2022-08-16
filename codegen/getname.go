package codegen

import (
	"strings"
)

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
			if len(field) > 1 {
				field = strings.ToUpper(field[0:1]) + field[1:]
			} else {
				field = strings.ToUpper(field)
			}
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
			fields[i] = strings.ToUpper(field[0:1]) + field[1:]
		}
	}
	return strings.Join(fields, "")
}
