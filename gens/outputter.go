package gens

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vektra/mockery/v2/pkg/logging"
	"gitlab.id.vin/nam.nguyen10/typeinfo/config"
)

type Cleanup func() error

type OutputStreamProvider interface {
	GetInterfaceWriter(context.Context, *Interface, string) (io.Writer, error, Cleanup)
	GetStructWriter(context.Context, *Struct, string) (io.Writer, error, Cleanup)
}

type StdoutStreamProvider struct {
}

func (this *StdoutStreamProvider) GetWriter(ctx context.Context, iface *Interface) (io.Writer, error, Cleanup) {
	return os.Stdout, nil, func() error { return nil }
}

type FileOutputStreamProvider struct {
	Config                    config.Config
	BaseDir                   string
	InPackage                 bool
	TestOnly                  bool
	Case                      string
	KeepTree                  bool
	KeepTreeOriginalDirectory string
	FileName                  string
}

func (this *FileOutputStreamProvider) GetInterfaceWriter(ctx context.Context, iface *Interface, extension string) (io.Writer, error, Cleanup) {
	log := zerolog.Ctx(ctx).With().Str(logging.LogKeyInterface, iface.Name).Logger()
	ctx = log.WithContext(ctx)

	var path string

	caseName := iface.Name
	if this.Case == "underscore" || this.Case == "snake" {
		caseName = this.underscoreCaseName(caseName)
	}

	if this.KeepTree {
		absOriginalDir, err := filepath.Abs(this.KeepTreeOriginalDirectory)
		if err != nil {
			return nil, err, func() error { return nil }
		}
		relativePath := strings.TrimPrefix(
			filepath.Join(filepath.Dir(iface.FileName), this.filename(caseName, extension)),
			absOriginalDir)
		path = filepath.Join(this.BaseDir, relativePath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	} else if this.InPackage {
		path = filepath.Join(filepath.Dir(iface.FileName), this.filename(caseName, extension))
	} else {
		path = filepath.Join(this.BaseDir, this.filename(caseName, extension))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	}

	log = log.With().Str(logging.LogKeyPath, path).Logger()
	ctx = log.WithContext(ctx)

	log.Debug().Msgf("creating writer to file")
	f, err := os.Create(path)
	if err != nil {
		return nil, err, func() error { return nil }
	}

	return f, nil, func() error {
		return f.Close()
	}
}

func (this *FileOutputStreamProvider) GetStructWriter(ctx context.Context, str *Struct, extension string) (io.Writer, error, Cleanup) {
	log := zerolog.Ctx(ctx).With().Str(logging.LogKeyInterface, str.Name).Logger()
	ctx = log.WithContext(ctx)

	var path string

	caseName := str.Name
	if this.Case == "underscore" || this.Case == "snake" {
		caseName = this.underscoreCaseName(caseName)
	}

	if this.KeepTree {
		absOriginalDir, err := filepath.Abs(this.KeepTreeOriginalDirectory)
		if err != nil {
			return nil, err, func() error { return nil }
		}
		relativePath := strings.TrimPrefix(
			filepath.Join(filepath.Dir(str.FileName), this.filename(caseName, extension)),
			absOriginalDir)
		path = filepath.Join(this.BaseDir, relativePath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	} else if this.InPackage {
		path = filepath.Join(filepath.Dir(str.FileName), this.filename(caseName, extension))
	} else {
		path = filepath.Join(this.BaseDir, this.filename(caseName, extension))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	}

	log = log.With().Str(logging.LogKeyPath, path).Logger()
	ctx = log.WithContext(ctx)

	log.Debug().Msgf("creating writer to file")
	f, err := os.Create(path)
	if err != nil {
		return nil, err, func() error { return nil }
	}

	return f, nil, func() error {
		return f.Close()
	}
}

func (this *FileOutputStreamProvider) filename(name string, extension string) string {
	if this.FileName != "" {
		return this.FileName
	}
	return name + extension
}

// shamelessly taken from http://stackoverflow.com/questions/1175208/elegant-python-function-to-convert-camelcase-to-camel-caseo
func (this *FileOutputStreamProvider) underscoreCaseName(caseName string) string {
	rxp1 := regexp.MustCompile("(.)([A-Z][a-z]+)")
	s1 := rxp1.ReplaceAllString(caseName, "${1}_${2}")
	rxp2 := regexp.MustCompile("([a-z0-9])([A-Z])")
	return strings.ToLower(rxp2.ReplaceAllString(s1, "${1}_${2}"))
}
