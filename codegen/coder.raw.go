package codegen

import "github.com/dave/jennifer/jen"

var _ Coder = RawCoder{}

type RawCoder struct {
	Value jen.Code
}

func (f RawCoder) Code() jen.Code {
	return f.Value
}

func AnonymousField(name string) RawCoder {
	return RawCoder{jen.Id(name)}
}

func Comment(text string) RawCoder {
	return RawCoder{jen.Comment(text)}
}

func CommentWithBreak(text string) RawCoder {
	return RawCoder{jen.Comment(text).Line()}
}
