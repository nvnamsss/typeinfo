package gens

import (
	"encoding/json"
	"strings"
)

type JSONFormat struct {
}

type structtype struct {
	Name        string
	Description string
}

type methodtype struct {
	Name        string
	Description string
	Params      []vartype
	Return      string
}

type vartype struct {
	Name string
	Type string
}

func (JSONFormat) Extension() string {
	return ".json"
}

func (JSONFormat) Start() string {
	return "{"
}

func (JSONFormat) Separate() string {
	return ","
}

func (JSONFormat) Struct(str *Struct) string {
	builder := strings.Builder{}
	strtype := structtype{
		Name:        str.Name,
		Description: str.Comment,
	}

	bytes, _ := json.Marshal(strtype)
	builder.WriteString("\"struct\":")
	builder.Write(bytes)
	return builder.String()
}

func (JSONFormat) Methods(methods []*Method) string {
	builder := strings.Builder{}
	ms := []methodtype{}
	for _, m := range methods {
		jm := methodtype{
			Name:        m.Name,
			Description: m.Comment,
		}
		paramsCount := m.Signature.Params().Len()
		for loop := 0; loop < paramsCount; loop++ {
			jm.Params = append(jm.Params, vartype{
				Name: m.Signature.Params().At(loop).Name(),
				Type: m.Signature.Params().At(loop).Type().String(),
			})
		}

		if m.Signature.Results().Len() > 0 {
			jm.Return = m.Signature.Results().At(0).Type().String()
		}
		ms = append(ms, jm)
	}

	bytes, _ := json.Marshal(ms)
	builder.WriteString("\"methods\":")
	builder.Write(bytes)

	return builder.String()
}

func (JSONFormat) Fields(fields []*Field) string {
	builder := strings.Builder{}
	vs := []vartype{}
	for _, f := range fields {
		vs = append(vs, vartype{Name: f.Name, Type: f.Var.Type().String()})
	}

	bytes, _ := json.Marshal(vs)
	builder.WriteString("\"fields\":")
	builder.Write(bytes)

	return builder.String()
}

func (JSONFormat) End() string {
	return "}"
}

func NewJSONFormat() *JSONFormat {
	return &JSONFormat{}
}
