package gens

type Format interface {
	Start() string
	Separate() string
	Struct(*Struct) string
	Methods([]*Method) string
	Fields([]*Field) string
	End() string
	Extension() string
}
