package main

import (
	"context"

	"gitlab.id.vin/nam.nguyen10/typeinfo/gens"
)

func main() {
	visitor := &gens.GeneratorVisitor{}
	walker := gens.Walker{
		BaseDir:   "examples",
		Recursive: true,
	}
	walker.Walk(context.Background(), visitor)
}
