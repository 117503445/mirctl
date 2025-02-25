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

func (e *alpineExecutor) PreCheck() bool {
	return true
}

func (e *alpineExecutor) Run() error {
	fileRelease := "/etc/os-release"

	release, err := goutils.ReadText(fileRelease)
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
	log.Debug().Str("ver", ver).Str("fileRelease", fileRelease).Str("content", release).Msg("read alpine version")

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

	src := "/etc/apk/repositories"
	dst := src + "." + goutils.TimeStrMilliSec() + ".bak"
	log.Debug().Str("src", src).Str("dst", dst).Msg("backup")
	err = goutils.CopyFile(src, dst)
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
