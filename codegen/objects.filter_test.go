package codegen

import (
	"bytes"
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
	implementFilter := func(name string) {
		builder.Add(RawCoder{jen.Func().Params(jen.Id("f").Id(name)).Id("filter").Call().Block()})
		builder.Add(RawCoder{jen.Var().Id("_").Id("Filter").Op("=").Id(name).Block()})
	}

	sections, err := ParseObjectDoc(url)
	if err != nil {
		t.Fatal(err)
	}

	for _, section := range sections {
		heading := section.Heading

		if heading == nil {
			desc := section.ParagraphText()
			builder.Add(RawCoder{jen.Comment(desc).Line().Type().Id("Filter").Interface(
				jen.Id("filter").Params(),
			)})
		} else {
			title := heading.Text

			switch title {
			case "Property filter object", "Timestamp filter object":
				desc := section.FirstParagraphText()
				name := getName(strings.TrimSuffix(title, " object"))
				err := builder.AddClass(name, desc).AddParams(section.Parameters()...)
				if err != nil {
					t.Fatal(err)
				}
				implementFilter(name)
			case "Compound filter object":
				desc := section.FirstParagraphText()
				obj := builder.AddClass(getName(strings.TrimSuffix(title, " object")), desc)
				for _, param := range section.Parameters() {
					obj.AddField(Property{
						Name:        param.Name,
						Type:        jen.Index().Id("Filter"),
						Description: param.Description,
						OmitEmpty:   true,
					})
				}
				implementFilter(obj.Name)
			case "Type-specific filter conditions": // ignore
			default:
				if !strings.HasSuffix(title, "filter condition") {
					t.Error(title)
				}

				desc := section.ParagraphText()

				match := descRegex.FindStringSubmatch(desc)
				typesStr := "[" + strings.Replace(match[1], ", and ", ", ", 1) + "]"
				types := []string{}

				if err := json.Unmarshal([]byte(typesStr), &types); err != nil {
					t.Fatal(err)
				}

				object := builder.AddClass(getName(title), desc)
				for _, param := range section.Parameters() {
					prop, err := param.Property()
					if err != nil {
						t.Fatal(err)
					}

					prop.OmitEmpty = true
					if !strings.HasPrefix(prop.Name, "is_") { // is_ 以外のプロパティで
						buffer := bytes.NewBuffer(nil)
						jen.Add(prop.Type).Render(buffer)
						code := string(buffer.Bytes())
						if code == "" {
							// インターフェイス
						} else if strings.HasPrefix(code, "*") {
							// 既に参照型
						} else {
							prop.Type = jen.Op("*").Add(prop.Type) // それ以外を参照型にする (omitempty)
						}
					}

					object.AddField(prop)
				}
				for _, propName := range types {
					prop := Property{Name: propName, Type: jen.Op("*").Id(object.Name), OmitEmpty: true}
					builder.GetClass("PropertyFilter").AddField(prop)
				}
			}
		}
	}

	if err := builder.Build("../types.filter.go"); err != nil {
		t.Fatal(err)
	}
}
