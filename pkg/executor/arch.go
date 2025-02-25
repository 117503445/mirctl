package executor

import (
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/mirctl/pkg/utils"
	"github.com/rs/zerolog/log"
)

// https://mirrors.ustc.edu.cn/help/archlinux.html

type archExecutor struct {
}

func (e *archExecutor) PreCheck() bool {

	return true
}

func (e *archExecutor) Run() error {
	err := goutils.WriteText("/etc/pacman.d/mirrorlist", `Server = https://mirrors.ustc.edu.cn/archlinux/$repo/os/$arch`)
	if err != nil {
		log.Warn().Err(err).Msg("write pacman mirrorlist error")
		return err
	}

	cmd := utils.ExecGetCmd(strings.Split("pacman -Sy", " "))
	_, err = utils.RunCmdWithLog(cmd)
	if err != nil {
		log.Warn().Err(err).Msg("apk update error")
		return err
	}

	return nil
}
