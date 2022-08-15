package codegen

import "github.com/dave/jennifer/jen"

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

func (c *Class) AddDocProps(props ...DocProp) *Class {
	for _, p := range props {
		p := p
		c.Fields = append(c.Fields, &p)
	}
	return c
}

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

func (c *Class) Code() jen.Code {
	fields := []jen.Code{}
	for _, f := range c.Fields {
		fields = append(fields, f.Code())
	}

	code := jen.Empty()
	if c.Comment != "" {
		code.Comment(c.Comment).Line()
	}
	return code.Type().Id(c.Name).Struct(fields...).Line()
}
