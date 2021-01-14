package gens

import (
	"fmt"
	"strings"
)

type TextFormatter struct {
}

func (TextFormatter) Start() string {
	return "Generated in txt format" + NewLine
}

func (TextFormatter) Separate() string {
	return NewLine
}

func (TextFormatter) Struct(str *Struct) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Struct: %v", str.Name))
	return builder.String()
}

func (TextFormatter) Fields(fields []*Field) string {
	builder := strings.Builder{}
	builder.WriteString(NewLine)
	builder.WriteString("Fields: ")

	for _, f := range fields {
		builder.WriteString(NewLine)
		builder.WriteString(fmt.Sprintf("- %v: %v", f.Name(), f.Type()))
	}

	return builder.String()
}

func (TextFormatter) Methods(methods []*Method) string {
	builder := strings.Builder{}
	builder.WriteString(NewLine)
	builder.WriteString("Methods: ")

	for _, m := range methods {
		builder.WriteString(NewLine)
		builder.WriteString(fmt.Sprintf("- %v: %v %v", m.Name(), m.signature.Params(), m.signature.Results()))
	}

	return builder.String()
}

func (TextFormatter) End() string {
	return NewLine
}

func (TextFormatter) Extension() string {
	return ".txt"
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}
