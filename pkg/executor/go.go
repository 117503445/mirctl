package executor

import (
	"strings"

	"github.com/117503445/mirctl/pkg/utils"
	"github.com/rs/zerolog/log"
)

// https://goproxy.cn/

type goExecutor struct {
}

func (e *goExecutor) PreCheck() bool {

	return utils.CommandExists("go")
}

func (e *goExecutor) Run() error {
	cmd := utils.ExecGetCmd(strings.Split("go env -w GOPROXY=https://goproxy.cn,direct", " "))
	_, err := utils.RunCmdWithLog(cmd)
	if err != nil {
		log.Warn().Err(err).Msg("setup goproxy error")
		return err
	}

	return nil
}
