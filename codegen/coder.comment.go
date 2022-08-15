package codegen

import "github.com/dave/jennifer/jen"

var _ Coder = Comment("")

type Comment string

func (f Comment) Code() jen.Code {
	return jen.Comment(string(f))
}

type CommentWithBreak string

func (f CommentWithBreak) Code() jen.Code {
	return jen.Comment(string(f)).Line()
}
