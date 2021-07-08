package gens

import (
	"encoding/json"
	"strings"
)

const (
	Depth = 4
)

type Format interface {
	SetMethods([]*Method)
	SetFields([]*Field)
	Format() string
	Extension() string
}

type methodType struct {
	Name        string `json:"-"`
	Description string
	Params      []fieldType
	Return      fieldType
	Variadic    bool
}

type fieldType struct {
	Name string
	Type interface{}
}

// Format json 1
type JsonFormat1 struct {
	// Field    jf1field
	// Function []methodType
	fields  []*Field
	methods []*Method
}

func (JsonFormat1) Extension() string {
	return ".json"
}

func (this *JsonFormat1) recursiveField(m map[string]string, name string, field *Field, depth int) {
	if depth >= Depth {
		return
	}

	if field.Struct() == nil {
		k := field.Name()
		if name != "" {
			k = name + "." + k
		}

		v := field.Type().String()
		m[k] = v
		return
	}

	str := field.Struct()

	fields := str.Fields()
	if name != "" {
		name = name + "." + field.Name()
	} else {
		name = field.Name()
	}

	for _, f := range fields {
		this.recursiveField(m, name, f, depth+1)
	}
}

func (JsonFormat1) End() string {
	return "}"
}

func (f *JsonFormat1) SetMethods(methods []*Method) {
	f.methods = methods
}

func (f *JsonFormat1) SetFields(fields []*Field) {
	f.fields = fields
}

func (this *JsonFormat1) Format() string {
	type str struct {
		Field    map[string]string
		Function map[string]methodType
	}

	var (
		sb strings.Builder
		mf map[string]string = make(map[string]string)
		mm                   = make(map[string]methodType)
	)

	for _, m := range this.methods {
		jm := methodType{
			Name:        m.Name(),
			Description: m.Comment,
		}

		params := m.Params()
		jm.Params = make([]fieldType, 0, len(params))
		for _, p := range params {

			jm.Params = append(jm.Params, fieldType{
				Name: p.Name(),
				Type: p.Type().String(),
			})
		}

		if r := m.Return(); r != nil {
			jm.Return = fieldType{
				Name: r.Name(),
				Type: r.Type().String(),
			}
		}

		mm[jm.Name] = jm
	}

	for _, field := range this.fields {
		this.recursiveField(mf, "", field, 1)
	}

	st := str{
		Field:    mf,
		Function: mm,
	}

	fs, _ := json.Marshal(st)
	sb.Write(fs)

	return sb.String()
}

func NewJF1() Format {
	return &JsonFormat1{}
}

// Format json 1
type JsonFormat2 struct {
	fields  []*Field
	methods []*Method
}

func (JsonFormat2) Extension() string {
	return ".json"
}

func (this *JsonFormat2) recursiveField(m map[string]interface{}, field *Field, depth int) {
	if field.Struct() == nil {
		k := field.Name()
		v := field.Type().String()
		m[k] = v
		return
	}

	if depth >= Depth-1 {
		return
	}
	str := field.Struct()
	mm := make(map[string]interface{})
	fields := str.Fields()

	for _, f := range fields {
		this.recursiveField(mm, f, depth+1)
	}
	m[field.Name()] = mm
}

func (f *JsonFormat2) SetMethods(methods []*Method) {
	f.methods = methods
}

func (f *JsonFormat2) SetFields(fields []*Field) {
	f.fields = fields
}

func (this *JsonFormat2) Format() string {
	type str struct {
		Field    map[string]interface{}
		Function map[string]methodType
	}

	var (
		sb strings.Builder
		mf map[string]interface{} = make(map[string]interface{})
		mm                        = make(map[string]methodType)
	)

	for _, m := range this.methods {
		jm := methodType{
			Name:        m.Name(),
			Description: m.Comment,
			Variadic:    m.signature.Variadic(),
		}

		params := m.Params()
		jm.Params = make([]fieldType, 0, len(params))
		for _, p := range params {
			jm.Params = append(jm.Params, fieldType{
				Name: p.Name(),
				Type: p.Type().String(),
			})
		}

		if r := m.Return(); r != nil {
			jm.Return = fieldType{
				Name: r.Name(),
				Type: r.Type().String(),
			}
		}

		mm[jm.Name] = jm
	}

	for _, field := range this.fields {
		this.recursiveField(mf, field, 1)
	}

	st := str{
		Field:    mf,
		Function: mm,
	}

	fs, _ := json.Marshal(st)
	sb.Write(fs)

	return sb.String()
}

func NewJF2() Format {
	return &JsonFormat2{}
}
