package codegen

import (
	"bytes"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var _ Coder = &Class{}

// Class はクラスを表し、Coderを実装します
type Class struct {
	Name       string
	Comment    string
	Fields     []Coder
	implements []string
}

func (c *Class) AddField(fields ...Coder) *Class {
	c.Fields = append(c.Fields, fields...)
	return c
}

func (c *Class) AddParams(opt *PropertyOption, params ...ObjectDocParameter) error {
	for _, param := range params {
		prop, err := param.Property(opt)
		if err != nil {
			return err
		}
		c.Fields = append(c.Fields, prop)
	}
	return nil
}

// AddConfiguration はNotion API特有の、特定のtypeに応じたプロパティを追加します
func (c *Class) AddConfiguration(tagName string, className string, comment string) *Class {
	p := &Property{Name: tagName, TypeSpecific: true, Description: comment}
	if className != "" {
		p.Type = jen.Op("*").Id(className)
	} else {
		p.Type = jen.Struct()
	}
	c.Fields = append(c.Fields, p)
	return c
}

func (c *Class) Implement(method string) *Class {
	c.implements = append(c.implements, method)
	return c
}

func (c *Class) Code() jen.Code {
	fields := []jen.Code{}
	for _, f := range c.Fields {
		fields = append(fields, f.Code())
	}

	code := jen.Empty()
	if c.Comment != "" {
		code.Comment(strings.TrimSpace(c.Comment)).Line()
	}
	code.Type().Id(c.Name).Struct(fields...).Line()

	for _, method := range c.implements {
		code.Func().Params(jen.Id("c").Op("*").Id(c.Name)).Id(method).Params().Block().Line()
	}

	if unions := c.unionProperties(); len(unions) != 0 {
		unmFields := jen.Statement{}

		for _, prop := range unions {
			buffer := bytes.NewBuffer(nil)
			jen.Var().Id("_").Add(prop.Type).Render(buffer)
			unionName := strings.TrimPrefix(buffer.String(), "var _ ")

			title := cases.Title(language.Und).String(prop.Name)
			unmFields.Add(
				jen.Id("p").Dot(title).Op("=").Id("new"+unionName).Call(jen.Id("data"), jen.Lit(prop.Name)),
			)
		}

		unmFields.Add(
			jen.Type().Id("Alias").Id(c.Name),
			jen.Return().Qual("encoding/json", "Unmarshal").Call(
				jen.Id("data"),
				jen.Call(jen.Op("*").Id("Alias")).Call(jen.Id("p")),
			),
		)

		code.Func().Params(jen.Id("p").Op("*").Id(c.Name)).Id("UnmarshalJSON").Params(jen.Id("data").Index().Byte()).Params(
			jen.Error(),
		).Block(
			unmFields...,
		)
	}

	if c.hasTypeSpecificProperty() {
		code.Func().Params(jen.Id("p").Id(c.Name)).Id("MarshalJSON").Params().Params(
			jen.Index().Byte(), jen.Error(),
		).Block(
			jen.Type().Id("Alias").Id(c.Name),
			jen.Return().Id("marshalByType").Call(
				jen.Id("Alias").Call(jen.Id("p")),
				jen.Id("p").Dot("Type"),
			),
		)
	}

	return code
}

func (c *Class) hasTypeSpecificProperty() bool {
	for _, f := range c.Fields {
		if p, ok := f.(*Property); ok && p.TypeSpecific {
			return true
		}
	}
	return false
}

func (c *Class) unionProperties() []*Property {
	unionProps := []*Property{}
	for _, f := range c.Fields {
		if p, ok := f.(*Property); ok {
			if p.IsUnion {
				unionProps = append(unionProps, p)
			}
		}
	}
	return unionProps
}
