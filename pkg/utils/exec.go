package utils

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	if err != nil {
		return false
	}
	return true
}

func ExecGetCmd(cmds []string) *exec.Cmd {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	return cmd
}

func RunCmdWithLog(cmd *exec.Cmd) (string, error) {
	formatDuration := func(d time.Duration) string {
		// 将 duration 转换为秒
		sec := d.Seconds()

		// 确定合适的单位和数值范围
		if sec < 1 {
			ms := d.Milliseconds() // 毫秒
			return fmt.Sprintf("%dms", ms)
		} else if sec >= 1 && sec < 60 {
			return fmt.Sprintf("%.3fs", sec)
		} else {
			return fmt.Sprintf("%.3gs", sec)
		}
	}

	start := time.Now()
	log.Info().Str("cmd", cmd.String()).CallerSkipFrame(1).Send()

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)

	log.Info().Str("output", output).Err(err).Str("duration", formatDuration(time.Since(start))).CallerSkipFrame(1).Send()

	return output, err
}
