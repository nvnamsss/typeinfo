package gens

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"sort"
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
	conf.Mode = packages.LoadSyntax
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
	// To support relative paths to mock targets w/ vendor deps, we need to provide eventual
	// calls to build.Context.Import with an absolute path. It needs to be absolute because
	// Import will only find the vendor directory if our target path for parsing is under
	// a "root" (GOROOT or a GOPATH). Only absolute paths will pass the prefix-based validation.
	//
	// For example, if our parse target is "./ifaces", Import will check if any "roots" are a
	// prefix of "ifaces" and decide to skip the vendor search.
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

func (p *Parser) Find(name string) (*Interface, error) {
	for _, entry := range p.entries {
		for _, iface := range entry.interfaces {
			if iface == name {
				list := p.packageInterfaces(entry.pkg.Types, entry.fileName, []string{name}, nil)
				if len(list) > 0 {
					return list[0], nil
				}
			}
		}
	}
	return nil, ErrNotInterface
}

func (p *Parser) Interfaces() []*Interface {
	ifaces := make(sortableIFaceList, 0)
	for _, entry := range p.entries {
		declaredIfaces := entry.interfaces
		ifaces = p.packageInterfaces(entry.pkg.Types, entry.fileName, declaredIfaces, ifaces)
	}

	sort.Sort(ifaces)
	return ifaces
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

		obj.Type()
		typ, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}

		name = typ.Obj().Name()

		if typ.Obj().Pkg() == nil {
			continue
		}

		elem := &Struct{
			Name:     name,
			Pkg:      pkg,
			FileName: fileName,
			named:    typ,
			methods:  []*Method{},
		}
		// str, ok := typ.Underlying().(*types.Struct)
		n2 := typ.NumMethods()
		for loop := 0; loop < n2; loop++ {
			mm := typ.Method(loop)
			sig, ok := mm.Type().Underlying().(*types.Signature)
			if ok {
				fmt.Println(sig.Params())
			}
			method := &Method{Name: mm.Name(), Signature: sig}

			if index := searchComment(comments, int(mm.Pos())); index != -1 {
				method.Comment = comments[index].Text()
			}

			elem.methods = append(elem.methods, method)
		}

		structs = append(structs, elem)
	}

	return structs
}

func (p *Parser) packageInterfaces(
	pkg *types.Package,
	fileName string,
	declaredInterfaces []string,
	ifaces []*Interface) []*Interface {
	scope := pkg.Scope()
	for _, name := range declaredInterfaces {
		obj := scope.Lookup(name)
		if obj == nil {
			continue
		}

		typ, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}

		name = typ.Obj().Name()

		if typ.Obj().Pkg() == nil {
			continue
		}

		elem := &Interface{
			Name:          name,
			Pkg:           pkg,
			QualifiedName: pkg.Path(),
			FileName:      fileName,
			NamedType:     typ,
		}

		iface, ok := typ.Underlying().(*types.Interface)
		if ok {
			elem.IsFunction = false
			elem.ActualInterface = iface
		} else {
			sig, ok := typ.Underlying().(*types.Signature)
			if !ok {
				continue
			}

			elem.IsFunction = true
			elem.SingleFunction = &Method{Name: "Execute", Signature: sig}
		}

		ifaces = append(ifaces, elem)
	}

	return ifaces
}

func searchComment(comments []*ast.CommentGroup, pos int) int {
	length := len(comments)
	for loop := length - 1; loop >= 0; loop-- {
		if int(comments[loop].Pos()) < pos {
			return loop
		}
	}

	return -1
}
