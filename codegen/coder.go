package codegen

import "github.com/dave/jennifer/jen"

// Coder はjen.Codeに変換可能なインターフェイスです
type Coder interface {
	Code() jen.Code
}
