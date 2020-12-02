package gens

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vektra/mockery/v2/pkg/logging"
	"gitlab.id.vin/nam.nguyen10/typeinfo/config"
)

type Walker struct {
	config.Config
	BaseDir   string
	Recursive bool
	Filter    *regexp.Regexp
	LimitOne  bool
	BuildTags []string
}

type WalkerVisitor interface {
	VisitWalk(context.Context, *Interface) error
	VisitStruct(context.Context, *Struct) error
}

func (this *Walker) Walk(ctx context.Context, visitor WalkerVisitor) (generated bool) {
	log := zerolog.Ctx(ctx)
	ctx = log.WithContext(ctx)

	log.Info().Msgf("Walking")

	parser := NewParser(this.BuildTags)
	this.doWalk(ctx, parser, this.BaseDir, visitor)

	if err := parser.Load(); err != nil {
		fmt.Printf("Error walking: %v", err)
	}

	for _, iface := range parser.Interfaces() {
		if !this.Filter.MatchString(iface.Name) {
			continue
		}
		err := visitor.VisitWalk(ctx, iface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s: %s\n", iface.Name, err)
			os.Exit(1)
		}
		generated = true
		if this.LimitOne {
			return
		}
	}

	for _, str := range parser.Structs() {
		if this.Filter != nil && !this.Filter.MatchString(str.Name) {
			continue
		}

		if err := visitor.VisitStruct(ctx, str); err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s: %s\n", str.Name, err)
			os.Exit(1)
		}
	}

	return
}

func (this *Walker) doWalk(ctx context.Context, p *Parser, dir string, visitor WalkerVisitor) (generated bool) {
	log := zerolog.Ctx(ctx)
	ctx = log.WithContext(ctx)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "_") {
			continue
		}

		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			if this.Recursive {
				generated = this.doWalk(ctx, p, path, visitor) || generated
				if generated && this.LimitOne {
					return
				}
			}
			continue
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			continue
		}

		err = p.Parse(ctx, path)
		if err != nil {
			log.Err(err).Msgf("Error parsing file")
			continue
		}
	}

	return
}

type GeneratorVisitor struct {
	config.Config
	InPackage bool
	Note      string
	Osp       OutputStreamProvider
	// The name of the output package, if InPackage is false (defaults to "mocks")
	PackageName       string
	PackageNamePrefix string
	StructName        string
}

func (this *GeneratorVisitor) VisitWalk(ctx context.Context, iface *Interface) error {
	log := zerolog.Ctx(ctx).With().
		Str(logging.LogKeyInterface, iface.Name).
		Str(logging.LogKeyQualifiedName, iface.QualifiedName).
		Logger()
	ctx = log.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("Unable to generate information: %s", r)
			return
		}
	}()

	// generator := NewInformationGenerator()

	// err := generator.Generate()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (this *GeneratorVisitor) VisitStruct(ctx context.Context, str *Struct) error {
	log := zerolog.Ctx(ctx).With().
		Str(logging.LogKeyInterface, str.Name).
		Logger()
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("Unable to generate information: %s", r)
			return
		}
	}()
	var format Format = NewJSONFormat()
	switch this.Config.Format {
	case "json":
		format = NewJSONFormat()
	case "txt":
		format = NewTextFormatter()
	default:
		format = NewJSONFormat()
	}

	out, err, closer := this.Osp.GetStructWriter(ctx, str, format.Extension())
	if err != nil {
		log.Err(err).Msgf("Unable to get writer")
		os.Exit(1)
	}
	defer closer()

	generator := NewInformationGenerator(str, format)

	if err := generator.Generate(ctx); err != nil {
		log.Error().Msgf("Generate file error: %v", err)
	}

	if err := generator.Write(out); err != nil {
		log.Error().Msgf("Write file error: %v", err)
	} else {
		log.Info().Msgf("Write struct: %v", str.Name)
	}

	return nil
}
