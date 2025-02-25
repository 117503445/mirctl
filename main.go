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

	for _, repo := range cfg.Cfg.Repos {
		log.Debug().Str("repo", repo).Send()
		executor.Run(repo, cfg.Cfg.Mirrors[0])
	}
}
