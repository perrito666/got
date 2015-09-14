package git

import (
	"log"
	"os"
	"os/exec"

	"github.com/juju/errors"
)

// Config holds the attributes required to invoke git.
type Config struct {
	SubCommand string
	Args       []string
}

// ErrSubCommandEmpty indicates that Config did not validate because
// of zero valued SubCommand.
var ErrSubCommandEmpty = errors.New("Git sub-command is empty")

func (c *Config) validate() error {
	if c.SubCommand == "" {
		return ErrSubCommandEmpty
	}
	return nil
}

func (c *Config) asArray() []string {
	if c.Args == nil || len(c.Args) == 0 {
		return []string{c.SubCommand}
	}
	args := make([]string, len(c.Args)+1)
	args[0] = c.SubCommand
	for i, v := range c.Args {
		args[i+1] = v
	}
	return args
}

// execCmd represents exec.Cmd
type execCmd interface {
	Run() error
}

type commandCraftFunc func([]string) execCmd

func command(args []string) execCmd {
	var cmd *exec.Cmd
	if args == nil || len(args) == 0 {
		cmd = exec.Command("git")
	} else {
		cmd = exec.Command("git", args...)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

type commandCraftAndRunFunc func([]string) (string, error)

func runCommand(args []string) (string, error) {
	var cmd *exec.Cmd
	if args == nil || len(args) == 0 {
		cmd = exec.Command("git")
	} else {
		cmd = exec.Command("git", args...)
	}
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Trace(err)
	}
	return string(out), nil
}

// Git makes a call to git with the given parameters.
func Git(c *Config) error {
	return git(c, command)
}

// LibGit will run a git command and return the stdout or error.
func LibGit(c *Config) (string, error) {
	return libGit(c, runCommand)
}

func libGit(c *Config, cmd commandCraftAndRunFunc) (string, error) {
	if c == nil {
		result, err := cmd(nil)
		if err != nil {
			log.Fatal(err)
			return "", errors.Trace(err)
		}
		return result, nil
	}
	if err := c.validate(); err != nil {
		return "", errors.Trace(err)
	}

	return cmd(c.asArray())
}

func git(c *Config, cmd commandCraftFunc) error {
	if c == nil {
		if err := cmd(nil).Run(); err != nil {
			log.Fatal(err)
			return errors.Trace(err)
		}
		return nil
	}
	if err := c.validate(); err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(cmd(c.asArray()).Run())
}
