package codegen

import "github.com/dave/jennifer/jen"

type Builder struct {
	Coders []Coder
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Add(coder Coder) *Builder {
	b.Coders = append(b.Coders, coder)
	return b
}

func (b *Builder) AddClass(name string, comment string) *Class {
	cls := &Class{Name: name, Comment: comment}
	b.Add(cls)
	return cls
}

func (b *Builder) GetClass(name string) *Class {
	for _, coder := range b.Coders {
		if cls, ok := coder.(*Class); ok && cls.Name == name {
			return cls
		}
	}
	return nil
}

func (b *Builder) Build(fileName string) error {
	f := jen.NewFile("notionapi")
	f.Comment("Code generated by notion.codegen; DO NOT EDIT.").Line()
	for _, coder := range b.Coders {
		f.Add(coder.Code())
	}
	return f.Save(fileName)
}
