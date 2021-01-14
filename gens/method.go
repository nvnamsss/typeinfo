package gens

import "go/types"

type Method struct {
	Comment   string
	_func     *types.Func
	signature *types.Signature
}

func (m *Method) Name() string {
	return m._func.Name()
}

func (m *Method) Params() []*Field {
	var (
		params = []*Field{}
	)

	paramsCount := m.signature.Params().Len()
	for loop := 0; loop < paramsCount; loop++ {
		v := m.signature.Params().At(loop)
		params = append(params, &Field{
			_var: v,
		})
	}
	return params
}

func (m *Method) Return() *Field {
	if m.signature.Results().Len() > 0 {
		v := m.signature.Results().At(0)
		return &Field{
			_var: v,
		}
	}

	return nil
}
