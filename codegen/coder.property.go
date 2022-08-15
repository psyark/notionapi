package codegen

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

var _ Coder = &Property{}

// Property は（DocPropと対照的に）自由に設定できるプロパティのCoderです
type Property struct {
	Name        string
	Type        jen.Code
	Description string
	OmitEmpty   bool
}

func (f Property) Code() jen.Code {
	jsonTag := f.Name
	if f.OmitEmpty {
		jsonTag += ",omitempty"
	}
	code := jen.Id(getName(f.Name)).Add(f.Type)
	code.Tag(map[string]string{"json": jsonTag})
	if f.Description != "" {
		code.Comment(strings.ReplaceAll(f.Description, "\n", " "))
	}
	return code
}
