package gens

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

const (
	NewLine = "\n"
)

type InformationGenerator struct {
	format Format
	str    *Struct
	buf    bytes.Buffer
}

func (g *InformationGenerator) printf(s string, vals ...interface{}) {
	fmt.Fprintf(&g.buf, s, vals...)
}

func NewInformationGenerator(str *Struct, format Format) *InformationGenerator {
	return &InformationGenerator{str: str, format: format}
}

func (g *InformationGenerator) Generate(ctx context.Context) error {
	g.format.SetFields(g.str.Fields())
	g.format.SetMethods(g.str.Methods())
	g.printf(g.format.Format())

	return nil
}

func (g *InformationGenerator) Write(w io.Writer) error {
	bytes := g.buf.Bytes()
	_, err := w.Write(bytes)
	return err
}

func (g *InformationGenerator) ToString() string {
	return g.buf.String()
}

func (g *InformationGenerator) Bytes() []byte {
	return g.buf.Bytes()
}
