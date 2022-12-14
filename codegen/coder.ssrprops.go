package codegen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
)

type MethodCoder struct {
	DocURL string
	Props  SSRProps
	Type   MethodCoderType
}

type MethodCoderType interface {
	New() jen.Code
	Returns() jen.Code
	Access(name string) jen.Code
}

type ReturnsStructRef struct {
	Name string
}

func (r *ReturnsStructRef) New() jen.Code {
	return jen.Op("&").Id(r.Name)
}
func (r *ReturnsStructRef) Returns() jen.Code {
	return jen.Op("*").Id(r.Name)
}
func (r *ReturnsStructRef) Access(name string) jen.Code {
	return jen.Id(name)
}

type ReturnsInterface struct {
	Name string
}

func (r *ReturnsInterface) New() jen.Code {
	return jen.Op("&").Id(fmt.Sprintf("_%vUnmarshaller", r.Name))
}
func (r *ReturnsInterface) Returns() jen.Code {
	return jen.Id(r.Name)
}
func (r *ReturnsInterface) Access(name string) jen.Code {
	return jen.Id(name).Dot(r.Name)
}

func (c MethodCoder) Code() jen.Code {
	code := jen.Comment(c.Props.Doc.Title).Line().Comment(c.DocURL).Line()

	callingParams := []jen.Code{jen.Id("ctx")}
	typedParams := []jen.Code{jen.Id("ctx").Qual("context", "Context")}

	// APIのパスパラメータを引数化
	for _, param := range c.getParams("path") {
		if param.Type != "string" {
			panic(param.Type)
		}

		callingParams = append(callingParams, jen.Id(param.Name))
		typedParams = append(typedParams, jen.Id(param.Name).String())
	}
	// オプション構造体引数
	if c.hasOptions() {
		callingParams = append(callingParams, jen.Id("options"))
		typedParams = append(typedParams, jen.Id("options").Op("*").Id(nfCamelCase.String(c.Props.Doc.Title)+"Options"))
	}

	callingParams = append(callingParams, jen.Nil())
	statements := []jen.Code{
		jen.Id("result").Op(":=").Add(c.Type.New()).Block(),
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

		code := jen.List(jen.Return().Add(c.Type.Access("result")), jen.Id("c").Dot("call").Call(
			jen.Id("ctx"),
			jen.Lit(strings.ToUpper(c.Props.Doc.API.Method)),
			jen.Qual("fmt", "Sprintf").Call(pathParams...),
			options,
			jen.Id("result"),
			jen.Id("bodyWriter"),
		))
		statements = append(statements, code)
	}

	methodName := nfCamelCase.String(c.Props.Doc.Title)
	// 公開関数
	code.Func().Params(jen.Id("c").Op("*").Id("Client")).Id(methodName).Params(typedParams...).Params(c.Type.Returns(), jen.Error()).Block(
		jen.Return(jen.Id("c").Dot("_" + methodName).Call(callingParams...)),
	).Line()

	typedParams = append(typedParams, jen.Id("bodyWriter").Qual("io", "Writer"))

	// 非公開関数（テスト用）
	code.Func().Params(jen.Id("c").Op("*").Id("Client")).Id("_"+methodName).Params(typedParams...).Params(c.Type.Returns(), jen.Error()).Block(statements...).Line()

	if c.hasOptions() {
		fields := []jen.Code{}
		for _, param := range c.getParams("body") {
			fields = append(fields, c.getOptionField(param))
		}
		code.Type().Id(nfCamelCase.String(c.Props.Doc.Title) + "Options").Struct(fields...).Line()
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
	code := jen.Id(nfCamelCase.String(param.Name))
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
			code.Op("*").Id(nfCamelCase.String(param.Name))
		case "filter":
			code.Id(nfCamelCase.String(param.Name))
		case "icon":
			code.Id("FileOrEmoji")
		case "cover":
			code.Op("*").Id("File")
		case "properties":
			if strings.Contains(param.Desc, "and the values are [property values]") {
				code.Id("PropertyValueMap")
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
			code.Id("RichTextArray")
		case strings.Contains(param.Desc, "an array of [block objects](ref:block)"):
			code.Index().Id("Block")
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
