package gens

import "go/types"

type Method struct {
	Name      string
	Comment   string
	Signature *types.Signature
}
