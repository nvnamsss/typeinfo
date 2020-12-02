package config

var SemVer = "0.0.0-dev"

type Config struct {
	Name                 string
	All                  bool
	Dir                  string
	FileName             string
	Case                 string
	Config               string
	Cpuprofile           string
	DisableVersionString bool `mapstructure:"disable-version-string"`
	StructName           string
	KeepTree             bool
	Recursive            bool
	Output               string
	LogLevel             string `mapstructure:"log-level"`
	Version              bool
	SrcPkg               string
	Profile              string
	Quiet                bool
}
