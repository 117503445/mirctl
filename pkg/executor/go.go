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

	return true
}

func (e *goExecutor) Run() error {
	cmd := utils.ExecGetCmd(strings.Split("go env -w GOPROXY=https://goproxy.cn,direct", " "))
	_, err := utils.RunCmdWithLog(cmd)
	if err != nil {
		log.Warn().Err(err).Msg("apk update error")
		return err
	}

	return nil
}
