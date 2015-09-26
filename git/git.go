// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package git

import (
	"log"
	"os"
	"os/exec"

	"github.com/juju/errors"
)

const (
	// CMDGit is the git shell command.
	CMDGit = "git"
)

const (
	// SCMDCheckout is the git sub command for checkout.
	SCMDCheckout string = "checkout"
	// SCMDBranch is the git sub command for branch.
	SCMDBranch string = "branch"
)

func command(args []string) ExecCmd {
	var cmd *exec.Cmd
	if args == nil || len(args) == 0 {
		cmd = exec.Command(CMDGit)
	} else {
		cmd = exec.Command(CMDGit, args...)
	}

	return cmd
}

func stdCommand(args []string) ExecCmd {
	cmd := command(args).(*exec.Cmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}

// ErrSubCommandEmpty indicates that Config did not validate because
// of zero valued SubCommand.
var ErrSubCommandEmpty = errors.New("Git sub-command is empty")

// Call holds the attributes required to invoke git.
type Call struct {
	subCommand string
	args       []string
}

func New(s string, a []string) Compatible {
	return &Call{
		subCommand: s,
		args:       a,
	}
}

func (c *Call) SubCommand() string {
	return c.subCommand
}

func (c *Call) SetSubCommand(s string) {
	c.subCommand = s
}

func (c *Call) Args() []string {
	return c.args
}
func (c *Call) SetArgs(a []string) {
	c.args = a
}

// Validate checks that Call has all the required attributes.
func (c *Call) Validate() error {
	return nil
}

func (c *Call) asArray() []string {
	if c.args == nil || len(c.args) == 0 {
		return []string{c.subCommand}
	}
	args := make([]string, len(c.args)+1)
	args[0] = c.subCommand
	for i, v := range c.args {
		args[i+1] = v
	}
	return args
}

// Git returns a exec.Cmd for git ready with the passed configuration.
func (c *Call) Git() (ExecCmd, error) {
	cmd, err := c.git(command)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return cmd, nil
}

// Run runs git according to the passed config and
// returns err if it fails.
func (c *Call) Run() error {
	cmd, err := c.git(stdCommand)
	if err != nil {
		log.Fatal(err)
		return errors.Trace(err)
	}
	return cmd.Run()
}

// Checkout runs git checkout on the current git repo.
func (c *Call) Checkout(branch string) error {
	c.subCommand = SCMDCheckout
	c.args = []string{branch}
	return errors.Annotatef(c.Run(), "cannot checkout %q", branch)
}

func (c *Call) git(cmd CommandCraftFunc) (ExecCmd, error) {
	if c.subCommand == "" {
		return cmd(nil), nil
	}
	return cmd(c.asArray()), nil
}
