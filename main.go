package main

import (
	"github.com/117503445/goutils"
	"github.com/117503445/mirctl/pkg/cfg"
	"github.com/117503445/mirctl/pkg/executor"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()
	cfg.Load()

	log.Debug().Msg("hello")
	executor.Run(cfg.Cfg.Repos[0], cfg.Cfg.Mirrors[0])
}
