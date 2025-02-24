package cfg

import (
	"github.com/alecthomas/kong"

	"github.com/rs/zerolog/log"
)

var Cfg struct {
	Mirrors []string `help:"mirror names" default:"ustc"`
	Repos   []string `help:"repos" default:"alpine"`
}

func Load() {
	kong.Parse(&Cfg)
	log.Info().Interface("config", Cfg).Send()
}
