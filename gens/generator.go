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
	g.printf(g.format.Start())
	// g.printf(g.format.Struct(g.str))
	// g.printf(g.format.Separate())
	g.printf(g.format.Fields(g.str.Fields()))
	g.printf(g.format.Separate())
	g.printf(g.format.Methods(g.str.Methods()))
	g.printf(g.format.End())
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
