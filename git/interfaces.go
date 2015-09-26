// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package git

// ExecCmd represents exec.Cmd
type ExecCmd interface {
	Run() error
	Output() ([]byte, error)
}

// CommandCraftFunc represents a function that creates an
// ExecCmd from a command/args string slice.
type CommandCraftFunc func([]string) ExecCmd

// Compatible represents a struct that can perform
// all the tasks that git should.
type Compatible interface {
	Git() (ExecCmd, error)
	Run() error
	Checkout(string) error

	SubCommand() string
	SetSubCommand(string)

	Args() []string
	SetArgs([]string)
}

// CompatibleConstructor returns a new Git Compatible struct.
type CompatibleConstructor func(string, []string) Compatible
