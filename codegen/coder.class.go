package codegen

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

var _ Coder = &Class{}

// Class はクラスを表し、Coderを実装します
type Class struct {
	Name    string
	Comment string
	Fields  []Coder
}

func (c *Class) AddField(fields ...Coder) *Class {
	c.Fields = append(c.Fields, fields...)
	return c
}

func (c *Class) AddParams(params ...ObjectDocParameter) error {
	for _, param := range params {
		prop, err := param.Property()
		if err != nil {
			return err
		}
		c.Fields = append(c.Fields, prop)
	}
	return nil
}

// Deprecated:
func (c *Class) AddDocProps(props ...DocProp) *Class {
	for _, p := range props {
		p := p
		c.Fields = append(c.Fields, &p)
	}
	return c
}

// Deprecated:
// AddConfiguration はNotion API特有の、特定のtypeに応じたプロパティを追加します
func (c *Class) AddConfiguration(tagName string, className string, comment string) *Class {
	p := &Property{Name: tagName, OmitEmpty: true, Description: comment}
	if className != "" {
		p.Type = jen.Op("*").Id(className)
	} else {
		p.Type = jen.Op("*").Struct()
	}
	c.Fields = append(c.Fields, p)
	return c
}

// AddConfiguration2 はNotion API特有の、特定のtypeに応じたプロパティを追加します
func (c *Class) AddConfiguration2(tagName string, className string, comment string) *Class {
	p := &Property{Name: tagName, TypeSpecific: true, Description: comment}
	if className != "" {
		p.Type = jen.Id(className)
	} else {
		p.Type = jen.Struct()
	}
	c.Fields = append(c.Fields, p)
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
		if p, ok := f.(Property); ok && p.TypeSpecific {
			return true
		}
		if p, ok := f.(*Property); ok && p.TypeSpecific {
			return true
		}
	}
	return false
}
