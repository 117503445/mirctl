package executor

import (
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/mirctl/pkg/utils"
	"github.com/rs/zerolog/log"
)

// https://mirrors.ustc.edu.cn/help/alpine.html

type alpineExecutor struct {
}

func (a *alpineExecutor) PreCheck() bool {
	log.Info().Msg("alpine PreCheck")

	return true
}

func (a *alpineExecutor) Run() error {
	log.Info().Msg("alpine Run")

	release, err := goutils.ReadText("/etc/os-release")
	if err != nil {
		log.Error().Err(err).Msg("read /etc/os-release error")
		return err
	}
	// NAME="Alpine Linux"
	// ID=alpine
	// VERSION_ID=3.21.3
	// PRETTY_NAME="Alpine Linux v3.21"
	// HOME_URL="https://alpinelinux.org/"
	// BUG_REPORT_URL="https://gitlab.alpinelinux.org/alpine/aports/-/issues"

	lines := strings.Split(release, "\n")
	found := false
	ver := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "VERSION_ID=") {
			version := strings.TrimPrefix(line, "VERSION_ID=")
			ver = strings.Trim(version, "\"")
			// log.Info().Msgf("alpine version: %s", version)

			found = true
			break
		}
	}
	if !found {
		log.Error().Msg("alpine version not found")
		return nil
	}
	// 3.21.3
	log.Debug().Str("ver", ver).Send()

	// 3.21.3 -> v3.21
	ver = "v" + strings.Join(strings.Split(ver, ".")[:2], ".")

	content, err := utils.RenderTemplate(`{{.url}}/{{.ver}}/main
{{.url}}/{{.ver}}/community`, map[string]string{
		"url": "https://mirrors.ustc.edu.cn/alpine",
		"ver": ver,
	})
	if err != nil {
		log.Error().Err(err).Msg("render template error")
		return err
	}
	log.Debug().Str("content", content).Send()

	err = goutils.CopyFile("/etc/apk/repositories", "/etc/apk/repositories"+goutils.TimeStrMilliSec()+".bak")
	if err != nil {
		log.Error().Err(err).Msg("backup /etc/apk/repositories error")
		return err
	}

	err = goutils.WriteText("/etc/apk/repositories", content)
	if err != nil {
		log.Error().Err(err).Msg("write /etc/apk/repositories error")
		return err
	}

	cmd := utils.ExecGetCmd([]string{"apk", "update"})
	_, err = utils.RunCmdWithLog(cmd)
	if err != nil {
		log.Warn().Err(err).Msg("apk update error")
		return err
	}

	return nil
}
