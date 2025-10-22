package tools

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func KillProcessListeningOnPort(port int) error {
	portArg := fmt.Sprintf("-iTCP:%d", port)
	cmd := exec.Command("lsof", "-nP", portArg, "-sTCP:LISTEN", "-t")
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) == 0 {
			return nil
		}
		return fmt.Errorf("running lsof: %w", err)
	}
	if len(out) == 0 {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, ln := range lines {
		if ln == "" {
			continue
		}
		pid, parseErr := strconv.Atoi(strings.TrimSpace(ln))
		if parseErr != nil {
			continue
		}
		if killErr := syscall.Kill(pid, syscall.SIGTERM); killErr != nil {
			if killErr2 := syscall.Kill(pid, syscall.SIGKILL); killErr2 != nil {
				return fmt.Errorf("failed to kill pid %d: %v, %v", pid, killErr, killErr2)
			}
		}
	}
	return nil
}

func KillOpenPorts(ports []int) error {
	var aggErr error
	for _, p := range ports {
		if err := KillProcessListeningOnPort(p); err != nil {
			aggErr = fmt.Errorf("port %d: %w; %v", p, err, aggErr)
		}
	}
	return aggErr
}
