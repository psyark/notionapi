package codegen

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestFilterObject(t *testing.T) {
	t.Parallel()

	url := "https://developers.notion.com/reference/post-database-query-filter"
	builder := NewBuilder().Add(CommentWithBreak(url))

	descRegex := regexp.MustCompile(`can be applied to database properties of (?:type|types) ([^\.]+)\.`)

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	{
		desc := sections[0].AllParagraphText()
		builder.Add(RawCoder{jen.Comment(desc).Line().Type().Id("Filter").Interface(jen.Id("filter").Params())})
	}

	for _, section := range sections[1:] {
		title := section.Heading.Text

		switch title {
		case "Property filter object", "Timestamp filter object":
			desc := section.FirstParagraphText()
			name := nfCamelCase.String(strings.TrimSuffix(title, " object"))
			err := builder.AddClass(name, desc).Implement("filter").AddParams(nil, section.Parameters()...)
			if err != nil {
				t.Fatal(err)
			}
		case "Compound filter object":
			desc := section.FirstParagraphText()
			obj := builder.AddClass(nfCamelCase.String(strings.TrimSuffix(title, " object")), desc).Implement("filter")
			for _, param := range section.Parameters() {
				obj.AddField(&Property{
					Name:        param.Name,
					Type:        jen.Index().Id("Filter"),
					Description: param.Description,
					OmitEmpty:   true,
				})
			}
		case "Type-specific filter conditions": // ignore
		default:
			if !strings.HasSuffix(title, "filter condition") {
				t.Error(title)
			}

			desc := section.AllParagraphText()

			match := descRegex.FindStringSubmatch(desc)
			typesStr := "[" + strings.Replace(match[1], ", and ", ", ", 1) + "]"
			types := []string{}

			if err := json.Unmarshal([]byte(typesStr), &types); err != nil {
				t.Fatal(err)
			}

			object := builder.AddClass(nfCamelCase.String(title), desc)
			for _, param := range section.Parameters() {
				opt := &PropertyOption{
					OmitEmpty: true,
					Nullable:  !strings.HasPrefix(param.Name, "is_"), // is_ 以外のプロパティはNullable
				}

				if err := object.AddParams(opt, param); err != nil {
					t.Fatal(err)
				}
			}
			for _, propName := range types {
				prop := &Property{Name: propName, Type: jen.Op("*").Id(object.Name), OmitEmpty: true}
				builder.GetClass("PropertyFilter").AddField(prop)
			}
		}
	}

	if err := builder.Build("../types.filter.go"); err != nil {
		t.Fatal(err)
	}
}
