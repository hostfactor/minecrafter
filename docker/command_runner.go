package docker

import "os/exec"

type CommandRunner interface {
	Run(cmd *exec.Cmd) error
}

type DefaultCommandRunner struct {
}

func (d *DefaultCommandRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}
