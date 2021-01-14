package gens

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vektra/mockery/v2/pkg/logging"
	"golang.org/x/tools/go/packages"
)

type parserEntry struct {
	fileName   string
	pkg        *packages.Package
	syntax     *ast.File
	interfaces []string
	structs    []string
	comments   []*ast.CommentGroup
}

type Parser struct {
	entries           []*parserEntry
	entriesByFileName map[string]*parserEntry
	parserPackages    []*types.Package
	conf              packages.Config
}

func NewParser(buildTags []string) *Parser {
	var conf packages.Config
	conf.Mode = packages.NeedFiles | packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax
	if len(buildTags) > 0 {
		conf.BuildFlags = []string{"-tags", strings.Join(buildTags, ",")}
	}
	return &Parser{
		parserPackages:    make([]*types.Package, 0),
		entriesByFileName: map[string]*parserEntry{},
		conf:              conf,
	}
}

func (p *Parser) Parse(ctx context.Context, path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range files {
		log := zerolog.Ctx(ctx).With().
			Str(logging.LogKeyDir, dir).
			Str(logging.LogKeyFile, fi.Name()).
			Logger()
		ctx = log.WithContext(ctx)

		if filepath.Ext(fi.Name()) != ".go" || strings.HasSuffix(fi.Name(), "_test.go") {
			continue
		}

		log.Debug().Msgf("parsing")

		fname := fi.Name()
		fpath := filepath.Join(dir, fname)
		if _, ok := p.entriesByFileName[fpath]; ok {
			continue
		}

		pkgs, err := packages.Load(&p.conf, "file="+fpath)
		if err != nil {
			return err
		}
		if len(pkgs) == 0 {
			continue
		}
		if len(pkgs) > 1 {
			names := make([]string, len(pkgs))
			for i, p := range pkgs {
				names[i] = p.Name
			}
			panic(fmt.Sprintf("file %s resolves to multiple packages: %s", fpath, strings.Join(names, ", ")))
		}

		pkg := pkgs[0]
		if len(pkg.Errors) > 0 {
			return pkg.Errors[0]
		}
		if len(pkg.GoFiles) == 0 {
			continue
		}

		for idx, f := range pkg.GoFiles {
			if _, ok := p.entriesByFileName[f]; ok {
				continue
			}

			entry := parserEntry{
				fileName: f,
				pkg:      pkg,
				syntax:   pkg.Syntax[idx],
			}
			p.entries = append(p.entries, &entry)
			p.entriesByFileName[f] = &entry
		}
	}

	return nil
}

type NodeVisitor struct {
	declaredInterfaces []string
	declaredStructs    []string
	comments           []*ast.CommentGroup
}

func NewNodeVisitor() *NodeVisitor {
	return &NodeVisitor{
		declaredInterfaces: make([]string, 0),
	}
}

func (n *NodeVisitor) DeclaredInterfaces() []string {
	return n.declaredInterfaces
}

func (n *NodeVisitor) DeclaredStructs() []string {
	return n.declaredStructs
}

func (nv *NodeVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.InterfaceType, *ast.FuncType:
			nv.declaredInterfaces = append(nv.declaredInterfaces, n.Name.Name)
		case *ast.StructType:
			nv.declaredStructs = append(nv.declaredStructs, n.Name.Name)
		}
	case *ast.CommentGroup:
		nv.comments = append(nv.comments, n)
	}
	return nv
}

func (p *Parser) Load() error {
	for _, entry := range p.entries {
		nv := NewNodeVisitor()
		ast.Walk(nv, entry.syntax)

		entry.interfaces = nv.DeclaredInterfaces()
		entry.structs = nv.DeclaredStructs()
		entry.comments = nv.comments
	}
	return nil
}

func (p *Parser) Structs() []*Struct {
	structs := make([]*Struct, 0)
	for _, entry := range p.entries {
		declaredStructs := entry.structs
		structs = p.packageStructs(entry.pkg.Types, entry.fileName, declaredStructs, structs, entry.comments)
	}

	return structs
}

func (p *Parser) packageStructs(pkg *types.Package, fileName string, declaredStructs []string, structs []*Struct, comments []*ast.CommentGroup) []*Struct {
	scope := pkg.Scope()

	for _, name := range declaredStructs {
		obj := scope.Lookup(name)

		if obj == nil {
			continue
		}

		typ, ok := obj.Type().(*types.Named)
		if !ok || typ.Obj().Pkg() == nil {
			continue
		}

		name = typ.Obj().Name()

		str := &Struct{
			Name:     name,
			pkg:      pkg,
			FileName: fileName,
			named:    typ,
			methods:  []*Method{},
		}

		n2 := typ.NumMethods()
		prevPos := 0
		for loop := 0; loop < n2; loop++ {
			f := typ.Method(loop)
			sig, ok := f.Type().Underlying().(*types.Signature)
			if !ok {
				continue
			}

			method := &Method{_func: f, signature: sig}

			if index := searchComment(comments, int(f.Pos()), prevPos); index != -1 {
				method.Comment = comments[index].Text()
			}

			str.methods = append(str.methods, method)
			prevPos = int(f.Pos())
		}

		structs = append(structs, str)
	}

	return structs
}

// Naive search, will be improved later
func searchComment(comments []*ast.CommentGroup, pos int, upper int) int {
	length := len(comments)
	for loop := length - 1; loop >= 0; loop-- {
		commentPos := int(comments[loop].Pos())
		if commentPos <= pos && commentPos > upper {
			return loop
		}
	}
	return -1
}
