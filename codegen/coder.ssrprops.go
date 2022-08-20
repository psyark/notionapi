package codegen

import (
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

type MethodCoder struct {
	DocURL  string
	Props   SSRProps
	Returns string
}

func (c MethodCoder) Code() jen.Code {
	code := jen.Comment(c.Props.Doc.Title).Line().Comment(c.DocURL).Line()

	params := []jen.Code{
		jen.Id("ctx").Qual("context", "Context"),
	}
	for _, param := range c.getParams("path") {
		code := jen.Id(param.Name)
		switch param.Type {
		case "string":
			code.String()
		default:
			panic(param.Type)
		}
		params = append(params, code)
	}
	if c.hasOptions() {
		code := jen.Id("options").Op("*").Id(getMethodName(c.Props.Doc.Title) + "Options")
		params = append(params, code)
	}

	statements := []jen.Code{}

	{
		code := jen.Id("result").Op(":=").Op("&").Id(c.Returns).Block()
		statements = append(statements, code)
	}
	{
		pathParams := []jen.Code{}
		path := regexp.MustCompile(`\{\w+\}`).ReplaceAllStringFunc(c.Props.Doc.API.URL, func(s string) string {
			pathParams = append(pathParams, jen.Id(s[1:len(s)-1]))
			return "%v"
		})
		pathParams = append([]jen.Code{jen.Lit(path)}, pathParams...)

		options := jen.Nil()
		if c.hasOptions() {
			options = jen.Id("options")
		}

		code := jen.List(jen.Return().Id("result"), jen.Id("c").Dot("call").Call(
			jen.Id("ctx"),
			jen.Lit(strings.ToUpper(c.Props.Doc.API.Method)),
			jen.Qual("fmt", "Sprintf").Call(pathParams...),
			options,
			jen.Id("result"),
		))
		statements = append(statements, code)
	}

	code.Func().Params(jen.Id("c").Op("*").Id("Client")).Id(getMethodName(c.Props.Doc.Title)).Params(params...).Params(jen.Op("*").Id(c.Returns), jen.Error()).Block(statements...).Line()

	if c.hasOptions() {
		fields := []jen.Code{}
		for _, param := range c.getParams("body") {
			fields = append(fields, c.getOptionField(param))
		}
		code.Type().Id(getMethodName(c.Props.Doc.Title) + "Options").Struct(fields...).Line()
	}

	return code
}

func (c MethodCoder) getParams(in string) []SSRPropsDocAPIParam {
	params := []SSRPropsDocAPIParam{}
	for _, param := range c.Props.Doc.API.Params {
		if param.In == in {
			params = append(params, param)
		}
	}
	return params
}

func (c MethodCoder) hasOptions() bool {
	return len(c.getParams("body")) != 0
}

func (c MethodCoder) getOptionField(param SSRPropsDocAPIParam) jen.Code {
	code := jen.Id(getName(param.Name))
	switch param.Type {
	case "string":
		code.String()
	case "int":
		code.Int()
	case "boolean":
		switch param.Name {
		case "archived":
			code.Op("*").Bool()
		default:
			code.Bool()
		}
	case "json":
		switch param.Name {
		case "parent":
			code.Op("*").Id(getName(param.Name))
		case "properties":
			if strings.Contains(param.Desc, "and the values are [property values]") {
				code.Map(jen.String()).Id("PropertyValue")
			} else if strings.Contains(param.Desc, "and the values are [property schema objects]") {
				code.Map(jen.String()).Interface() // TODO
			} else {
				panic(param.Desc)
			}
		default:
			code.Map(jen.String()).Interface() // TODO
		}
	case "array_object", "array_mixed":
		switch {
		case strings.Contains(param.Desc, "An array of [rich text objects](ref:rich-text)"):
			code.Index().Id("RichText")
		default:
			code.Index().Interface()
		}
	default:
		panic(param.Type)
	}

	jsonTag := param.Name
	if !param.Required {
		jsonTag += ",omitempty"
	}
	code.Tag(map[string]string{"json": jsonTag})

	code.Comment(param.Desc)
	return code
}
