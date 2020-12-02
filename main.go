package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"runtime/pprof"
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
	cfgFile = ""
	rootCmd = &cobra.Command{
		Use:   "tinfo",
		Short: "Generate mock objects for your Golang interfaces",
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
	// visitor := &gens.GeneratorVisitor{}
	// walker := gens.Walker{
	// 	BaseDir:   "examples",
	// 	Recursive: true,
	// }
	// walker.Walk(context.Background(), visitor)
}

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringVar(&cfgFile, "config", "", "config file to use")
	pFlags.String("name", "", "name or matching regular expression of interface to generate mock for")
	pFlags.String("output", "./mocks", "directory to write mocks to")
	pFlags.String("dir", ".", "directory to search for interfaces")
	pFlags.BoolP("recursive", "r", false, "recurse search into sub-directories")
	pFlags.Bool("all", false, "generates mocks for all found interfaces in all sub-directories")
	pFlags.String("case", "camel", "name the mocked file using casing convention [camel, snake, underscore]")
	pFlags.String("cpuprofile", "", "write cpu profile to file")
	pFlags.Bool("version", false, "prints the installed version of mockery")
	pFlags.Bool("quiet", false, `suppresses logger output (equivalent to --log-level="")`)
	pFlags.Bool("keeptree", false, "keep the tree structure of the original interface files into a different repository. Must be used with XX")
	pFlags.String("filename", "", "name of generated file (only works with -name and no regex)")
	pFlags.String("structname", "", "name of generated struct (only works with -name and no regex)")
	pFlags.String("log-level", "info", "Level of logging")
	pFlags.Bool("disable-version-string", false, "Do not insert the version string into the generated mock file.")

	viper.BindPFlags(pFlags)
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

	if r.Quiet {
		// if "quiet" flag is set, disable logging
		r.Config.LogLevel = ""
	}

	log, err := getLogger(r.Config.LogLevel)
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
		log.Fatal().Msgf("Specify --name or --all, but not both")
	} else if (r.Config.FileName != "" || r.Config.StructName != "") && r.Config.All {
		log.Fatal().Msgf("Cannot specify --filename or --structname with --all")
	} else if r.Config.Dir != "" && r.Config.Dir != "." && r.Config.SrcPkg != "" {
		log.Fatal().Msgf("Specify -dir or -srcgens, but not both")
	} else if r.Config.Name != "" {
		recursive = r.Config.Recursive
		if strings.ContainsAny(r.Config.Name, regexMetadataChars) {
			if filter, err = regexp.Compile(r.Config.Name); err != nil {
				log.Fatal().Err(err).Msgf("Invalid regular expression provided to -name")
			} else if r.Config.FileName != "" || r.Config.StructName != "" {
				log.Fatal().Msgf("Cannot specify --filename or --structname with regex in --name")
			}
		} else {
			filter = regexp.MustCompile(fmt.Sprintf("^%s$", r.Config.Name))
			limitOne = true
		}
	} else if r.Config.All {
		recursive = true
		filter = regexp.MustCompile(".*")
	} else {
		log.Fatal().Msgf("Use --name to specify the name of the interface or --all for all interfaces found")
	}

	if r.Config.Profile != "" {
		f, err := os.Create(r.Config.Profile)
		if err != nil {
			return errors.Wrapf(err, "Failed to create profile file")
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var osp gens.OutputStreamProvider
	osp = &gens.FileOutputStreamProvider{
		Config:                    r.Config,
		BaseDir:                   r.Config.Output,
		KeepTree:                  r.Config.KeepTree,
		KeepTreeOriginalDirectory: r.Config.Dir,
		FileName:                  r.Config.FileName,
	}

	baseDir := r.Config.Dir

	visitor := &gens.GeneratorVisitor{
		Config:     r.Config,
		Osp:        osp,
		StructName: r.Config.StructName,
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
