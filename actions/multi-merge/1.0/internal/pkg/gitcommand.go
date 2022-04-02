package pkg

import (
	"fmt"

	"github.com/erda-project/erda-actions/pkg/command"
	"github.com/erda-project/erda/pkg/crypto/uuid"
)

type Command interface {
	Execute() error
}

type GitInitCommand struct {
	cmds *command.Cmd
}

func (g *GitInitCommand) Execute() error {
	return g.cmds.Run()
}

func NewGitInitCommand(args ...string) Command {
	c := &GitInitCommand{
		cmds: command.NewCmd("/bin/bash", toolScript),
	}
	for _, arg := range args {
		c.cmds.Add(arg)
	}
	return c
}

type GitCloneCommand struct {
	cmds []*command.Cmd
}

func (g *GitCloneCommand) Execute() error {
	for _, cmd := range g.cmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func NewGitCloneCommand(repo string, branch string) Command {
	c := &GitCloneCommand{}
	initC := command.NewCmd("git", "init")
	remoteC := command.NewCmd("git", "remote", "add", "origin", repo)
	fetchC := command.NewCmd("git", "fetch", "origin", branch)
	checkoutC := command.NewCmd("git", "checkout", "FETCH_HEAD")
	c.cmds = []*command.Cmd{
		initC,
		remoteC,
		fetchC,
		checkoutC,
	}
	return c
}

type GitMergeCommand struct {
	cmds []*command.Cmd
}

func (g *GitMergeCommand) Execute() error {
	for _, cmd := range g.cmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func NewGitMergeCommand(repo string, branches []string) Command {
	c := &GitMergeCommand{cmds: make([]*command.Cmd, 0)}
	newRemote := uuid.New()
	remoteC := command.NewCmd("git", "remote", "add", newRemote, repo)
	c.cmds = append(c.cmds, remoteC)
	for _, branch := range branches {
		c.cmds = append(c.cmds, command.NewCmd("git", "fetch", newRemote, branch))
		c.cmds = append(c.cmds, command.NewCmd("git", "merge", newRemote+"/"+branch, "-m", fmt.Sprintf("Merge branch '%s' into head", branch)))
	}
	return c
}

type GitLoginCommand struct {
	cmd *command.Cmd
}

func (l *GitLoginCommand) Execute() error {
	return l.cmd.Run()
}

func NewGitLoginCommand(username string, password string) Command {
	c := &GitLoginCommand{
		cmd: command.NewCmd("echo"),
	}
	c.cmd.Add(fmt.Sprintf("\"default login %s password %s\"", username, password))
	c.cmd.Add(">")
	c.cmd.Add("$HOME/.netrc")
	return c
}

type GitConfigCommand struct {
	cmds []*command.Cmd
}

func (g *GitConfigCommand) Execute() error {
	for _, cmd := range g.cmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func NewGitConfigCommand(configs []GitConfig) Command {
	c := &GitConfigCommand{cmds: make([]*command.Cmd, 0)}
	for _, config := range configs {
		c.cmds = append(c.cmds, command.NewCmd("git", "config", "--global", config.Name, config.Value))
	}
	return c
}
