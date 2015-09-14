package git

import (
	"log"
	"os"
	"os/exec"

	"github.com/juju/errors"
)

const (
	CMD_GIT = "git"
)

const (
	SCMD_CHECKOUT string = "checkout"
	SCMD_BRANCH   string = "branch"
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
	Output() ([]byte, error)
}

type commandCraftFunc func([]string) execCmd

func command(args []string) execCmd {
	var cmd *exec.Cmd
	if args == nil || len(args) == 0 {
		cmd = exec.Command(CMD_GIT)
	} else {
		cmd = exec.Command(CMD_GIT, args...)
	}

	return cmd
}

func stdCommand(args []string) execCmd {
	cmd := command(args).(*exec.Cmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

// Git returns a exec.Cmd for git ready with the passed configuration.
func Git(c *Config) (*exec.Cmd, error) {
	cmd, err := git(c, command)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return cmd.(*exec.Cmd), nil
}

// Run runs git according to the passed config and
// returns err if it fails.
func Run(c *Config) error {
	cmd, err := git(c, stdCommand)
	if err != nil {
		log.Fatal(err)
		return errors.Trace(err)
	}
	return cmd.Run()
}

// Checkout runs git checkout on the current git repo.
func Checkout(branch string) error {
	c := &Config{
		SubCommand: SCMD_CHECKOUT,
		Args:       []string{branch},
	}
	return errors.Annotatef(Run(c), "cannot checkout %q", branch)
}

func git(c *Config, cmd commandCraftFunc) (execCmd, error) {
	if c == nil {
		return cmd(nil), nil
	}
	if err := c.validate(); err != nil {
		return nil, errors.Trace(err)
	}

	return cmd(c.asArray()), nil
}
