package gens

import (
	"context"
	"encoding/json"
	"strings"
)

const (
	fieldName  = `"field":`
	methodName = `"function":`
)

type JSONFormat struct {
}

type jsonFormat struct {
	Field    map[string]varType
	Function []methodType
}

type structType struct {
	Name        string
	Description string
}

type methodType struct {
	Name        string `json:"-"`
	Description string
	Params      []varType
	Return      varType
}

type varType struct {
	Name string
	Type interface{}
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

func (this *JSONFormat) Struct(str *Struct) string {
	builder := strings.Builder{}
	strtype := structType{
		Name:        str.Name,
		Description: str.Comment,
	}

	bytes, _ := json.Marshal(strtype)
	builder.WriteString("\"struct\":")
	builder.Write(bytes)
	return builder.String()
}

func (this *JSONFormat) Methods(methods []*Method) string {
	builder := strings.Builder{}
	ms := make(map[string]methodType)
	// ms := []methodType{}
	for _, m := range methods {
		jm := methodType{
			Name:        m.Name(),
			Description: m.Comment,
		}

		params := m.Params()
		jm.Params = make([]varType, 0, len(params))
		for _, p := range params {
			jm.Params = append(jm.Params, varType{
				Name: p.Name(),
				Type: p.Type().String(),
			})
		}

		if r := m.Return(); r != nil {
			jm.Return = varType{
				Name: r.Name(),
				Type: r.Type().String(),
			}
		}

		ms[jm.Name] = jm
		// ms = append(ms, jm)
	}

	bytes, _ := json.Marshal(ms)
	builder.WriteString(methodName)
	builder.Write(bytes)

	return builder.String()
}

func (this *JSONFormat) Fields(fields []*Field) string {
	var (
		builder  = strings.Builder{}
		varTypes = make(map[string]interface{}, len(fields))
	)

	for _, f := range fields {
		if str := f.Struct(); str != nil {
			g := NewInformationGenerator(str, this)
			g.Generate(context.TODO())

			jsonFormat := jsonFormat{}
			json.Unmarshal(g.buf.Bytes(), &jsonFormat)
			varTypes[f.Name()] = jsonFormat
		} else {
			varTypes[f.Name()] = f.Type().String()
		}

	}

	bytes, _ := json.Marshal(varTypes)
	builder.WriteString(fieldName)
	builder.Write(bytes)

	return builder.String()
}

func (this *JSONFormat) End() string {
	return "}"
}

func NewJSONFormat() *JSONFormat {
	return &JSONFormat{}
}
