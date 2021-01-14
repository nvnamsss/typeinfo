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
	GetStructWriter(context.Context, *Struct, string) (io.Writer, error, Cleanup)
}

type FileOutputStreamProvider struct {
	Config config.Config
}

func (o *FileOutputStreamProvider) GetStructWriter(ctx context.Context, str *Struct, extension string) (io.Writer, error, Cleanup) {
	log := zerolog.Ctx(ctx).With().Str(logging.LogKeyInterface, str.Name).Logger()

	var path string

	caseName := str.Name
	if o.Config.Case == "underscore" || o.Config.Case == "snake" {
		caseName = o.underscoreCaseName(caseName)
	}

	if o.Config.KeepTree {
		absOriginalDir, err := filepath.Abs(o.Config.Directory)
		if err != nil {
			return nil, err, func() error { return nil }
		}
		relativePath := strings.TrimPrefix(
			filepath.Join(filepath.Dir(str.FileName), o.filename(caseName, extension)),
			absOriginalDir)
		path = filepath.Join(o.Config.Output, relativePath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	} else {
		path = filepath.Join(o.Config.Output, o.filename(caseName, extension))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	}

	log = log.With().Str(logging.LogKeyPath, path).Logger()

	log.Debug().Msgf("creating writer to file")
	f, err := os.Create(path)
	if err != nil {
		return nil, err, func() error { return nil }
	}

	return f, nil, func() error {
		return f.Close()
	}
}

func (o *FileOutputStreamProvider) filename(name string, extension string) string {
	if o.Config.FileName != "" {
		return o.Config.FileName
	}
	return name + extension
}

func (o *FileOutputStreamProvider) underscoreCaseName(caseName string) string {
	rxp1 := regexp.MustCompile("(.)([A-Z][a-z]+)")
	s1 := rxp1.ReplaceAllString(caseName, "${1}_${2}")
	rxp2 := regexp.MustCompile("([a-z0-9])([A-Z])")
	return strings.ToLower(rxp2.ReplaceAllString(s1, "${1}_${2}"))
}
