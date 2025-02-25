package utils

import (
	"strings"

	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

var readReleaseResult map[string]string
var readReleaseErr error

func ReadRelease() (map[string]string, error) {
	if readReleaseResult != nil || readReleaseErr != nil {
		return readReleaseResult, readReleaseErr
	}

	content, err := goutils.ReadText("/etc/os-release")
	if err != nil {
		log.Error().Err(err).Msg("read /etc/os-release error")
		readReleaseErr = err
		return nil, err
	}
	log.Debug().CallerSkipFrame(1).Str("content", content).Msg("read /etc/os-release")

	r := make(map[string]string)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		splits := strings.Split(line, "=")
		if len(splits) == 2 {
			key := strings.TrimSpace(splits[0])
			value := strings.TrimSpace(splits[1])
			r[key] = value
		}
	}

	readReleaseResult = r
	return r, nil
}

func Backup(file string) error {
	src := file
	dst := src + "." + goutils.TimeStrMilliSec() + ".bak"
	log.Debug().CallerSkipFrame(1).
		Str("src", src).Str("dst", dst).Msg("backup")
	err := goutils.CopyFile(src, dst)
	if err != nil {
		log.Error().CallerSkipFrame(1).Err(err).Msg("backup error")
		return err
	}
	return nil
}
