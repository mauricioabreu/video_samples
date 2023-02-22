package extractor

import (
	"fmt"
	"os/exec"
)

type Command struct {
	executable string
	args       []string
}

func RunCmd(c Command) error {
	cmd := exec.Command(c.executable, c.args...) //#nosec G204
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cmd: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to wait cmd: %w", err)
	}

	return cmd.Wait()
}
