package internal

import (
	"os/exec"
	"path/filepath"
)

type Command struct {
	ffPath string
	cmd    *exec.Cmd
}

func NewCommand(ffPath string) *Command {
	return &Command{
		ffPath: ffPath,
	}
}

func (c *Command) ExecuteFFCommand(executable string, args []string) *exec.Cmd {
	cmd := exec.Command(filepath.Join(c.ffPath, executable), args...)
	c.cmd = cmd

	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Stdin = os.Stdin
	return cmd
}
