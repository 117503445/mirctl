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

	repos := cfg.Cfg.Repos
	if len(repos) == 1 && repos[0] == "AUTO" {
		repos = executor.PreCheck()
		log.Info().Strs("repos", repos).Msg("auto detect")
	}

	for _, repo := range repos {
		log.Info().Str("repo", repo).Send()
		executor.Run(repo, cfg.Cfg.Mirrors[0])
	}
}
