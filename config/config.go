package config

var SemVer = "0.1.0"

type Config struct {
	Name      string
	All       bool
	Directory string `mapstructure:"dir"`
	FileName  string
	Case      string
	KeepTree  bool
	Recursive bool
	Output    string
	Version   bool
	Format    string
}
