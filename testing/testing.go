package testing

import "github.com/perrito666/got/git"

type fakeCmd struct {
	ran int
}

// Run implements execCmd
func (f *fakeCmd) Run() error {
	f.ran++
	return nil
}

// Output implements execCmd
func (f *fakeCmd) Output() ([]byte, error) {
	return nil, nil
}

type FakeGit struct {
	cmd *fakeCmd
}

func (f *FakeGit) Git() (git.ExecCmd, error) {
	return f.cmd, nil
}
func (f *FakeGit) Run() error {
	return nil
}

func (f *FakeGit) Checkout(string) error {
	return nil
}

func (f *FakeGit) SubCommand() string {
	return ""
}

func (f *FakeGit) SetSubCommand(string) {
}

func (f *FakeGit) Args() []string {
	return []string{}
}
func (f *FakeGit) SetArgs([]string) {
}

func New(s string, a []string) git.Compatible {
	return &FakeGit{&fakeCmd{}}
}

type FakeUI struct {
}

func (*FakeUI) ChoiceMenu(items []string, one bool, curent int) ([]int, error) {
	return nil, nil
}
