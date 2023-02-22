package extractor

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Command struct {
	executable string
	args       []string
}

func RunCmd(c Command) error {
	cmd := exec.Command(c.executable, c.args...) //#nosec G204
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create a pipe from stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cmd: %w", err)
	}

	logCmd(stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait cmd: %w", err)
	}

	return nil
}

// logCmd capture the buffer and log each line of it
func logCmd(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		log.Info().Msg(scanner.Text())
	}
}
