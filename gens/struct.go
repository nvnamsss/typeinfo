package gens

import (
	"go/ast"
	"go/types"
)

type Struct struct {
	Name     string
	FileName string
	Comment  string
	methods  []*Method
	File     *ast.File
	Pkg      *types.Package
	named    *types.Named
}

func (structs *Struct) Methods() []*Method {
	return structs.methods
}

func (structs *Struct) Fields() []*Field {
	fields := make([]*Field, 0)
	str, ok := structs.named.Underlying().(*types.Struct)
	if !ok {
		return fields
	}

	count := str.NumFields()

	for loop := 0; loop < count; loop++ {
		v := str.Field(loop)
		fields = append(fields, &Field{Name: v.Name(), Var: v})
	}

	return fields
}
