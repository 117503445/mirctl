package executor

import (
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/mirctl/pkg/utils"
	"github.com/rs/zerolog/log"
)

// https://mirrors.ustc.edu.cn/help/ubuntu.html

type ubuntuExecutor struct {
}

func (e *ubuntuExecutor) PreCheck() bool {
	// PRETTY_NAME="Ubuntu 24.04.1 LTS"
	// NAME="Ubuntu"
	// VERSION_ID="24.04"
	// VERSION="24.04.1 LTS (Noble Numbat)"
	// VERSION_CODENAME=noble
	// ID=ubuntu
	// ID_LIKE=debian
	// HOME_URL="https://www.ubuntu.com/"
	// SUPPORT_URL="https://help.ubuntu.com/"
	// BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
	// PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
	// UBUNTU_CODENAME=noble
	// LOGO=ubuntu-logo
	release, err := utils.ReadRelease()
	if err != nil {
		return false
	}

	return strings.Contains(release["NAME"], "Ubuntu")
}

func (e *ubuntuExecutor) Run() error {
	release, err := utils.ReadRelease()
	if err != nil {
		return err
	}

	codename, ok := release["UBUNTU_CODENAME"]
	if !ok {
		log.Error().Msg("ubuntu codename not found")
		return nil
	}
	log.Debug().Str("codename", codename).Msg("read ubuntu codename")

	writeContent := func(proto string) error {
		content, err := utils.RenderTemplate(`Types: deb
URIs: {{.proto}}://mirrors.ustc.edu.cn/ubuntu
Suites: {{.codename}} {{.codename}}-updates {{.codename}}-backports
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

Types: deb
URIs: {{.proto}}://mirrors.ustc.edu.cn/ubuntu
Suites: {{.codename}}-security
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg`, map[string]string{
			"codename": codename,
			"proto":    proto,
		})
		if err != nil {
			log.Error().Err(err).Msg("render template error")
			return err
		}

		fileSource := "/etc/apt/sources.list.d/ubuntu.sources"
		log.Debug().Str("content", content).Str("file", fileSource).Send()

		err = goutils.WriteText(fileSource, content)
		if err != nil {
			log.Error().Err(err).Msg("write /etc/apt/sources.list.d/ubuntu.sources error")
			return err
		}
		return nil
	}

	cmd := utils.ExecGetCmd(strings.Split("dpkg -s ca-certificates", " "))
	output, err := utils.RunCmdWithLog(cmd)

	if strings.Contains(output, "is not installed") {
		writeContent("http")

		cmd = utils.ExecGetCmd(strings.Split("apt-get update", " "))
		_, err = utils.RunCmdWithLog(cmd)
		if err != nil {
			return err
		}

		cmd = utils.ExecGetCmd(strings.Split("apt-get install -y ca-certificates", " "))
		_, err = utils.RunCmdWithLog(cmd)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	writeContent("https")

	cmd = utils.ExecGetCmd(strings.Split("apt-get update", " "))
	_, err = utils.RunCmdWithLog(cmd)
	if err != nil {
		return err
	}

	return nil
}
