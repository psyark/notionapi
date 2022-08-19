package codegen

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

var _ Coder = &Property{}

// Property は（DocPropと対照的に）自由に設定できるプロパティのCoderです
type Property struct {
	Name         string
	Type         jen.Code
	Description  string
	OmitEmpty    bool
	TypeSpecific bool
}

func (f Property) Code() jen.Code {
	tags := map[string]string{"json": f.Name}
	if f.OmitEmpty {
		tags["json"] += ",omitempty"
	}
	if f.TypeSpecific {
		tags["specific"] = "type"
	}

	code := jen.Id(getName(f.Name)).Add(f.Type)
	code.Tag(tags)
	if f.Description != "" {
		code.Comment(strings.ReplaceAll(f.Description, "\n", " "))
	}
	return code
}
