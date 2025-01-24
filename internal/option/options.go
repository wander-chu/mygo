package option

import "errors"

var (
	errNotPrepared = errors.New("options not prepared")
)

type Options struct {
	ConfigFile string `toml:"-"`
	DBDir      string `toml:"-"`
	LogLevel   string `toml:"-"`
	Http       struct {
		Port int `toml:"port"`
	} `toml:"http"`
}
