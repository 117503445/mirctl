package executor

import "github.com/rs/zerolog/log"

type alpineExecutor struct {
}

func (a *alpineExecutor) Run() error {
	log.Info().Msg("alpine executor")
	return nil
}
