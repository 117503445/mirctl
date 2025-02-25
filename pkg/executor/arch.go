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
	// NAME="Arch Linux"
	// PRETTY_NAME="Arch Linux"
	// ID=arch
	// BUILD_ID=rolling
	// VERSION_ID=20250223.0.312761
	// ANSI_COLOR="38;2;23;147;209"
	// HOME_URL="https://archlinux.org/"
	// DOCUMENTATION_URL="https://wiki.archlinux.org/"
	// SUPPORT_URL="https://bbs.archlinux.org/"
	// BUG_REPORT_URL="https://gitlab.archlinux.org/groups/archlinux/-/issues"
	// PRIVACY_POLICY_URL="https://terms.archlinux.org/docs/privacy-policy/"
	// LOGO=archlinux-logo
	release, err := utils.ReadRelease()
	if err != nil {
		return false
	}

	return strings.Contains(release["NAME"], "Arch")
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
