package executor

import (
	"strings"

	"github.com/117503445/mirctl/pkg/utils"
	"github.com/rs/zerolog/log"
)

// https://mirrors.ustc.edu.cn/help/pypi.html

type pipExecutor struct {
}

func (e *pipExecutor) PreCheck() bool {

	return true
}

func (e *pipExecutor) Run() error {
	cmd := utils.ExecGetCmd(strings.Split("pip config set global.index-url https://mirrors.ustc.edu.cn/pypi/simple", " "))
	_, err := utils.RunCmdWithLog(cmd)
	if err != nil {
		log.Warn().Err(err).Msg("pip update error")
		return err
	}

	return nil
}
