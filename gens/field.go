package gens

import (
	"go/types"
)

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
	var (
		named *types.Named
	)

	if ss, ok := f._var.Type().Underlying().(*types.Pointer); ok {
		if n, ok := ss.Elem().(*types.Named); ok {
			named = n
		}
	}

	if n, ok := f._var.Type().(*types.Named); ok {
		named = n
	}

	if named != nil {
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
			str.methods = append(str.methods, method)
		}
		return &Struct{
			named: named,
		}
	}
	// if named, ok := f._var.Type().(*types.Named); ok {

	// }

	return nil
}

func (f *Field) String() string {
	return f._var.Type().String()
}
