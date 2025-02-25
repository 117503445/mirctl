package executor

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/117503445/goutils"
	"github.com/117503445/mirctl/pkg/utils"
	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
)

// https://mirrors.ustc.edu.cn/help/rust.html

type rustExecutor struct {
}

func (e *rustExecutor) PreCheck() bool {

	return utils.CommandExists("cargo")
}

func (e *rustExecutor) Run() error {
	cmd := utils.ExecGetCmd(strings.Split("cargo --version", " "))
	output, err := utils.RunCmdWithLog(cmd)
	if err != nil {
		log.Warn().Err(err).Msg("pip update error")
		return err
	}
	// cargo 1.85.0 (d73d2caf9 2024-12-31)
	verStr := strings.Split(output, " ")[1]

	log.Info().Str("cargo version", verStr).Send()

	ver, err := semver.NewVersion(verStr)
	if err != nil {
		return err
	}

	home := os.Getenv("HOME")
	if home == "" {
		log.Error().Msg("HOME env not found")
		return err
	}

	fileCfg := ""
	if ver.GreaterThan(semver.MustParse("1.38.0")) {
		fileCfg = home + "/.cargo/config.toml"
	} else {
		fileCfg = home + "/.cargo/config"
	}

	parts := map[string]string{
		"source.crates-io": `[source.crates-io]
replace-with = 'ustc'`,
		"source.ustc": `[source.ustc]
registry = "sparse+https://mirrors.ustc.edu.cn/crates.io-index/"`,
	}
	if ver.GreaterThan(semver.MustParse("1.68.0")) {
		parts["registries.ustc"] = `[registries.ustc]
index = "sparse+https://mirrors.ustc.edu.cn/crates.io-index/"`
	} else {
		parts["registries.ustc"] = `[registries.ustc]
index = "https://mirrors.ustc.edu.cn/crates.io-index/"`
	}

	exists := make(map[string]any)
	content := ""
	if goutils.FileExists(fileCfg) {
		err := utils.Backup(fileCfg)
		if err != nil {
			return err
		}
		content, err = goutils.ReadText(fileCfg)
		if err != nil {
			log.Warn().Err(err).Msg("read file error")
		} else {
			checkKeyExistence := func(data map[string]any, keyPath string) bool {
				keys := strings.Split(keyPath, ".")
				currentMap := data

				for _, key := range keys {
					if next, ok := currentMap[key].(map[string]any); ok {
						currentMap = next // Move down one level in the map hierarchy
					} else if _, ok := currentMap[key]; !ok {
						// If the key doesn't exist or it's not a map (meaning we've reached the end of the path)
						return false
					}
					// If we are at the last key and it exists, return true
					if key == keys[len(keys)-1] {
						return true
					}
				}
				return false
			}
			var configMap map[string]any
			if _, err := toml.Decode(content, &configMap); err != nil {
				log.Warn().Err(err).Msg("read file error")
				return err
			}
			for key := range parts {
				if checkKeyExistence(configMap, key) {
					exists[key] = nil
				}
			}
		}
	}

	existsList := []string{}
	for key := range exists {
		existsList = append(existsList, key)
	}
	log.Debug().Interface("exists", existsList).Send()
	for key, value := range parts {
		if _, ok := exists[key]; !ok {
			if content != "" {
				content += "\n"
			}
			content += value + "\n"
		}
	}

	log.Debug().Str("content", content).Str("file", fileCfg).Msg("write file")
	if err := goutils.WriteText(fileCfg, content); err != nil {
		log.Warn().Err(err).Msg("write file error")
		return err
	}

	return nil
}
