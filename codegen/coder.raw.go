package codegen

import "github.com/dave/jennifer/jen"

var _ Coder = RawCoder{}

type RawCoder struct {
	Value jen.Code
}

func (f RawCoder) Code() jen.Code {
	return f.Value
}
