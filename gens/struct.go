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
	pkg      *types.Package
	named    *types.Named

	comments []*ast.CommentGroup
}

func (this *Struct) Methods() []*Method {
	methods := make([]*Method, 0, len(this.methods))
	for _, m := range this.methods {
		if m._func.Exported() {
			methods = append(methods, m)
		}
	}

	return methods
}

func (this *Struct) Fields() []*Field {
	fields := make([]*Field, 0)
	str, ok := this.named.Underlying().(*types.Struct)
	if !ok {
		return fields
	}

	count := str.NumFields()

	for loop := 0; loop < count; loop++ {
		v := str.Field(loop)
		if v.Exported() {
			fields = append(fields, &Field{_var: v})
		}
	}
	return fields
}

func (this *Struct) String() string {
	return this.named.String()
}

// func NewStruct(obj types.Object) *Struct {
// 	if obj == nil {
// 		return nil
// 	}

// 	obj.Type()
// 	typ, ok := obj.Type().(*types.Named)
// 	if !ok || typ.Obj().Pkg() == nil {
// 		return nil
// 	}

// 	name = typ.Obj().Name()

// 	str := &Struct{
// 		Name:     name,
// 		pkg:      pkg,
// 		FileName: fileName,
// 		named:    typ,
// 		methods:  []*Method{},
// 	}
// 	return &Struct{}
// }
