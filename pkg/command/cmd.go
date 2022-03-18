package command

import (
	"os"
	"os/exec"
)

type Args []string

func (a *Args) Add(argument string) {
	*a = append(*a, argument)
}

type Cmd struct {
	Name string
	Args
	exec.Cmd
}

func NewCmd(name string, args ...string) *Cmd {
	c := &Cmd{
		Name: name,
		Args: make(Args, 0),
	}
	for _, arg := range args {
		c.Add(arg)
	}
	c.Cmd = *exec.Command(name, c.Args...)
	return c
}

func (c *Cmd) SetDir(dir string) {
	c.Cmd.Dir = dir
}

func (c *Cmd) Run() error {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Cmd.Args = append([]string{c.Name}, c.Args...)
	return c.Cmd.Run()
}
