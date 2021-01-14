package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.id.vin/nam.nguyen10/typeinfo/config"
	"gitlab.id.vin/nam.nguyen10/typeinfo/gens"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	rootCmd = &cobra.Command{
		Use:   "typeinfo",
		Short: "Generate information for your struct",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := GetRootAppFromViper(viper.GetViper())
			if err != nil {
				printStackTrace(err)
				return err
			}
			return r.Run()
		},
	}
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func printStackTrace(e error) {
	fmt.Printf("%v\n", e)
	if err, ok := e.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	}

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.String("name", "", "name or matching regular expression of interface to generate info for")
	pFlags.String("output", "./infos", "directory to write generated infos to")
	pFlags.String("format", "json", "file format info will be saved to")
	pFlags.String("dir", ".", "directory to search for generating struct")
	pFlags.BoolP("recursive", "r", false, "recurse search into sub-directories")
	pFlags.Bool("all", false, "generates info for all struct that found in directory")
	pFlags.String("case", "camel", "naming the generated file following convention [camel, snake, underscore]")
	pFlags.Bool("version", false, "prints the installed version of tinfo")
	pFlags.Bool("keeptree", false, "keep the hierarchy tree of the generated files that same as the original")
	pFlags.String("filename", "", "name of generated file (only works with --name and no regex)")

	_ = viper.BindPFlags(pFlags)
}

const regexMetadataChars = "\\.+*?()|[]{}^$"

type RootApp struct {
	config.Config
}

func GetRootAppFromViper(v *viper.Viper) (*RootApp, error) {
	r := &RootApp{}
	if err := v.UnmarshalExact(&r.Config); err != nil {
		return nil, errors.Wrapf(err, "failed to get config")
	}
	return r, nil
}

func (r *RootApp) Run() error {
	var recursive bool
	var filter *regexp.Regexp
	var err error
	var limitOne bool

	log, err := getLogger("info")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		return err
	}
	log.Info().Msgf("Starting tinfo")
	ctx := log.WithContext(context.Background())

	if r.Config.Version {
		fmt.Println(config.SemVer)
		return nil
	} else if r.Config.Name != "" && r.Config.All {
		log.Fatal().Msgf("Should specify only --name or --all")
	} else if r.Config.Name != "" {
		recursive = r.Config.Recursive
		if strings.ContainsAny(r.Config.Name, regexMetadataChars) {
			if filter, err = regexp.Compile(r.Config.Name); err != nil {
				log.Fatal().Err(err).Msgf("Invalid regular expression provided to -name")
			}
		} else {
			filter = regexp.MustCompile(fmt.Sprintf("^%s$", r.Config.Name))
			limitOne = true
		}
	} else if r.Config.All {
		recursive = true
		filter = regexp.MustCompile(".*")
	} else {
		log.Fatal().Msgf("Use --name to specify the name of the struct or --all for all structs found")
	}

	if r.Config.Format == "" {
		log.Warn().Msgf("Format is empty, default value json will be used instead")
	}

	osp := &gens.FileOutputStreamProvider{
		Config: r.Config,
	}

	baseDir := r.Config.Directory

	visitor := &gens.GeneratorVisitor{
		Config: r.Config,
		Osp:    osp,
	}

	walker := gens.Walker{
		Config:    r.Config,
		BaseDir:   baseDir,
		Recursive: recursive,
		Filter:    filter,
		LimitOne:  limitOne,
	}

	generated := walker.Walk(ctx, visitor)

	if r.Config.Name != "" && !generated {
		log.Fatal().Msgf("Unable to find '%s' in any go files under this path", r.Config.Name)
	}
	return nil
}

type timeHook struct{}

func (t timeHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Time("time", time.Now())
}

func getLogger(levelStr string) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		return zerolog.Logger{}, errors.Wrapf(err, "Couldn't parse log level")
	}
	out := os.Stderr
	writer := zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC822,
	}
	if !terminal.IsTerminal(int(out.Fd())) {
		writer.NoColor = true
	}
	log := zerolog.New(writer).
		Hook(timeHook{}).
		Level(level).
		With().
		Str("version", config.SemVer).
		Logger()

	return log, nil
}
