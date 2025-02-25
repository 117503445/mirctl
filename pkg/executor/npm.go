package executor

import (
	"strings"

	"github.com/117503445/mirctl/pkg/utils"
	"github.com/rs/zerolog/log"
)

// https://mirrors.ustc.edu.cn/help/pypi.html

type npmExecutor struct {
}

func (e *npmExecutor) PreCheck() bool {

	return true
}

func (e *npmExecutor) Run() error {
	cmd := utils.ExecGetCmd(strings.Split("npm config set registry https://npmreg.proxy.ustclug.org/", " "))
	_, err := utils.RunCmdWithLog(cmd)
	if err != nil {
		log.Warn().Err(err).Msg("npm update error")
		return err
	}

	return nil
}
