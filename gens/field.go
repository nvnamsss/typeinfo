package gens

import "go/types"

type Type interface {
	String() string
}

type Field struct {
	Comment string
	_var    *types.Var
}

// Name returns name of field
func (f *Field) Name() string {
	return f._var.Name()
}

// Type returns type of field
func (f *Field) Type() Type {

	return f
}

func (f *Field) Struct() *Struct {
	if f._var.Embedded() {
		if named, ok := f._var.Type().(*types.Named); ok {
			name := named.Obj().Name()

			str := &Struct{
				Name:    name,
				named:   named,
				methods: []*Method{},
			}

			n2 := named.NumMethods()
			// prevPos := 0
			for loop := 0; loop < n2; loop++ {
				f := named.Method(loop)
				sig, ok := f.Type().Underlying().(*types.Signature)
				if !ok {
					continue
				}

				method := &Method{_func: f, signature: sig}

				// if index := searchComment(comments, int(f.Pos()), prevPos); index != -1 {
				// 	method.Comment = comments[index].Text()
				// }

				str.methods = append(str.methods, method)
				// prevPos = int(f.Pos())
			}
			return &Struct{
				named: named,
			}
		}
	}
	return nil
}

func (f *Field) String() string {
	return f._var.Type().String()
}
