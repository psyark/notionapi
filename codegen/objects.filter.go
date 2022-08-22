package codegen

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

func BuildFilter() error {
	url := "https://developers.notion.com/reference/post-database-query-filter"
	builder := NewBuilder().Add(CommentWithBreak(url))

	implementFilter := func(name string) {
		builder.Add(RawCoder{jen.Func().Params(jen.Id("f").Id(name)).Id("filter").Call().Block()})
		builder.Add(RawCoder{jen.Var().Id("_").Id("Filter").Op("=").Id(name).Block()})
	}

	builder.Add(RawCoder{jen.Type().Id("Filter").Interface(
		jen.Id("filter").Params(),
	)})

	descRegex := regexp.MustCompile(`can be applied to database properties of (?:type|types) ([^\.]+)\.`)

	err := Parse(url, func(title, desc string, props []DocProp) error {
		if title == "Property filter object" {
			name := getName(strings.TrimSuffix(title, " object"))
			builder.AddClass(name, desc).AddDocProps(props...)
			implementFilter(name)
		} else if title == "Timestamp filter object" {
			name := getName(strings.TrimSuffix(title, " object"))
			builder.AddClass(name, desc).AddDocProps(props...)
			implementFilter(name)
		} else if title == "Compound filter object" {
			object := builder.AddClass(getName(strings.TrimSuffix(title, " object")), desc)
			for _, dp := range props {
				object.AddField(Property{
					Name:        dp.Name,
					Type:        jen.Index().Id("Filter"),
					Description: dp.Description,
					OmitEmpty:   true,
				})
			}
		} else if title == "Type-specific filter conditions" {
		} else if strings.HasSuffix(title, "filter condition") {
			match := descRegex.FindStringSubmatch(desc)
			typesStr := "[" + strings.Replace(match[1], ", and ", ", ", 1) + "]"
			types := []string{}

			if err := json.Unmarshal([]byte(typesStr), &types); err != nil {
				panic(err)
			}

			objName := getName(title)
			builder.AddClass(objName, desc).AddDocProps(props...)
			for _, propName := range types {
				prop := Property{Name: propName, Type: jen.Op("*").Id(objName), OmitEmpty: true}
				builder.GetClass("PropertyFilter").AddField(prop)
			}
		} else {
			panic(title)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return builder.Build("../types.filter.go")
}
