package executor

import (
	"strconv"
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/mirctl/pkg/utils"
	"github.com/rs/zerolog/log"
)

// https://mirrors.ustc.edu.cn/help/debian.html

type debianExecutor struct {
}

func (e *debianExecutor) PreCheck() bool {

	return true
}

func (e *debianExecutor) Run() error {
	release, err := utils.ReadRelease()
	if err != nil {
		return err
	}

	// PRETTY_NAME="Debian GNU/Linux 12 (bookworm)"
	// NAME="Debian GNU/Linux"
	// VERSION_ID="12"
	// VERSION="12 (bookworm)"
	// VERSION_CODENAME=bookworm
	// ID=debian
	// HOME_URL="https://www.debian.org/"
	// SUPPORT_URL="https://www.debian.org/support"
	// BUG_REPORT_URL="https://bugs.debian.org/"

	codename, ok := release["VERSION_CODENAME"]
	if !ok {
		log.Error().Msg("debian codename not found")
		return nil
	}
	versionIdStr, ok := release["VERSION_ID"]
	if !ok {
		log.Error().Msg("debian versionId not found")
		return nil
	}
	versionId, err := strconv.Atoi(strings.Trim(versionIdStr, "\""))
	if err != nil {
		log.Error().Err(err).Msg("fail to parse debian versionId")
		return nil
	}
	log.Debug().Str("codename", codename).Int("versionId", versionId).Msg("read debian release")

	nonfree := ""
	if versionId >= 12 {
		nonfree = "non-free-firmware"
	}

	writeContent := func(proto string) error {
		content, err := utils.RenderTemplate(`Types: deb
URIs: {{.proto}}://mirrors.ustc.edu.cn/debian
Suites: {{.codename}} {{.codename}}-updates
Components: main contrib non-free {{.nonfree}}
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg

Types: deb
URIs: {{.proto}}://mirrors.ustc.edu.cn/debian-security
Suites: {{.codename}}-security
Components: main contrib non-free {{.nonfree}}
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg`, map[string]string{
			"codename": codename,
			"proto":    proto,
			"nonfree":  nonfree,
		})
		if err != nil {
			log.Error().Err(err).Msg("render template error")
			return err
		}

		fileSource := "/etc/apt/sources.list.d/debian.sources"
		log.Debug().Str("content", content).Str("file", fileSource).Send()

		err = goutils.WriteText(fileSource, content)
		if err != nil {
			log.Error().Err(err).Msg("write /etc/apt/sources.list.d/debian.sources error")
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
