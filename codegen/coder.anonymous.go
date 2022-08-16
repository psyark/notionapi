package codegen

import "github.com/dave/jennifer/jen"

var _ Coder = AnonymousField("")

type AnonymousField string

func (f AnonymousField) Code() jen.Code {
	return jen.Id(string(f))
}
