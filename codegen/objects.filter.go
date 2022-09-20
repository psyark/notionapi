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
			implementFilter(object.Name)
		} else if title == "Type-specific filter conditions" {
		} else if strings.HasSuffix(title, "filter condition") {
			match := descRegex.FindStringSubmatch(desc)
			typesStr := "[" + strings.Replace(match[1], ", and ", ", ", 1) + "]"
			types := []string{}

			if err := json.Unmarshal([]byte(typesStr), &types); err != nil {
				panic(err)
			}

			object := builder.AddClass(getName(title), desc)
			for _, dp := range props {
				prop := Property{Name: dp.Name, Description: dp.Description, OmitEmpty: true}
				switch dp.Type {
				case "string":
					prop.Type = jen.String()
				case "boolean":
					prop.Type = jen.Op("*").Bool()
				case "boolean (only true)":
					prop.Type = jen.Bool()
				case "number":
					prop.Type = jen.Op("*").Float64()
				case "string (ISO 8601 date)":
					prop.Type = jen.Id("ISO8601String")
				case "string (UUIDv4)":
					prop.Type = jen.Id("UUIDString")
				case "object (empty)":
					prop.Type = jen.Struct()
				case "object (text filter condition)", "object (number filter condition)", "object (date filter condition)", "object (checkbox filter condition)":
					t := strings.TrimSuffix(strings.TrimPrefix(dp.Type, "object ("), ")")
					prop.Type = jen.Id(getName(t))
				case "object":
					prop.Type = jen.Interface()
				default:
					panic(dp.Type)
				}
				object.AddField(prop)
			}
			for _, propName := range types {
				prop := Property{Name: propName, Type: jen.Op("*").Id(object.Name), OmitEmpty: true}
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
