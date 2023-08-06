package installer

import (
	"io"
	"os/exec"
)

type CommandRunner interface {
	Command(name string, arg ...string) Runner
}

type Runner interface {
	StderrPipe() (io.ReadCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	Kill() error
}

var (
	_ CommandRunner = (*cmdRunner)(nil)
	_ Runner        = (*execRunner)(nil)
)

type execRunner struct {
	*exec.Cmd
}

func (e *execRunner) Kill() error {
	return e.Process.Kill()
}

type cmdRunner struct {
}

func (*cmdRunner) Command(name string, arg ...string) Runner {
	return &execRunner{exec.Command(name, arg...)}
}
