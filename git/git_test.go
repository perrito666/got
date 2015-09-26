package git

import (
	"os/exec"
	"strings"
	"testing"
)

func TestConfigAsArrayReturnsProperArrayWithArgs(t *testing.T) {
	c := Call{
		subCommand: "blah",
		args:       []string{"string1", "string2"},
	}
	arr := c.asArray()
	if len(arr) != 3 {
		t.Logf("expected 3 got %d", len(arr))
		t.Fail()
	}

	if arr[0] != "blah" {
		t.Logf("expected %q got %q", c.subCommand, arr[0])
		t.Fail()
	}

	if arr[1] != "string1" {
		t.Logf("expected %q got %q", c.args[0], arr[1])
		t.Fail()
	}

	if arr[2] != "string2" {
		t.Logf("expected %q got %q", c.args[1], arr[2])
		t.Fail()
	}
}

func TestConfigAsArrayReturnsProperArrayWithoutArgs(t *testing.T) {
	c := Call{
		subCommand: "blah",
		args:       []string{},
	}
	arr := c.asArray()
	if len(arr) != 1 {
		t.Logf("expected 1 got %d", len(arr))
		t.Fail()
	}

	if arr[0] != "blah" {
		t.Logf("expected %q got %q", c.subCommand, arr[0])
		t.Fail()
	}
}

func TestConfigCommandCraftsProperCommandWithArgs(t *testing.T) {
	c := Call{
		subCommand: "blah",
		args:       []string{"string1", "string2"},
	}
	com := command(c.asArray()).(*exec.Cmd)

	if !strings.HasSuffix(com.Path, "git") {
		t.Logf("expected \"<path>/git\" obtained %q", com.Path)
		t.Fail()
	}

	if len(com.Args) != 4 {
		t.Logf("expected 4 got %d", len(com.Args))
		t.Fail()
	}

	if com.Args[0] != "git" {
		t.Logf("expected \"git\" got %q", com.Args[0])
		t.Fail()
	}

	if com.Args[1] != "blah" {
		t.Logf("expected %q got %q", c.subCommand, com.Args[1])
		t.Fail()
	}

	if com.Args[2] != "string1" {
		t.Logf("expected %q got %q", c.args[0], com.Args[2])
		t.Fail()
	}

	if com.Args[3] != "string2" {
		t.Logf("expected %q got %q", c.args[1], com.Args[3])
		t.Fail()
	}
}

func TestConfigCommandCraftsProperCommandWithoutArgs(t *testing.T) {
	com := command(nil).(*exec.Cmd)

	if !strings.HasSuffix(com.Path, "git") {
		t.Logf("expected \"<path>/git\" obtained %q", com.Path)
		t.Fail()
	}

	if len(com.Args) != 1 {
		t.Logf("expected 1 got %d", len(com.Args))
		t.Fail()
	}

	if com.Args[0] != "git" {
		t.Logf("expected \"git\" got %q", com.Args[0])
		t.Fail()
	}
}

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

func TestGitCallsProrperlyWithConfig(t *testing.T) {
	var gotStrings []string
	cmd := fakeCmd{0}
	cmdFunc := func(s []string) ExecCmd {
		gotStrings = s
		return &cmd
	}

	c := Call{
		subCommand: "blah",
		args:       []string{"string1", "string2"},
	}
	var gitCmd ExecCmd
	var err error
	if gitCmd, err = c.git(cmdFunc); err != nil {
		t.Logf("git failed, got error: %v", err)
		t.Fail()
	}

	// Test config has been properly passed
	if err := gitCmd.Run(); err != nil {
		t.Logf("Run failed, got error: %v", err)
		t.Fail()
	}

	if len(gotStrings) != 3 {
		t.Logf("expected 3 got %d", len(gotStrings))
		t.Fail()
	}

	if gotStrings[0] != "blah" {
		t.Logf("expected %q got %q", c.subCommand, gotStrings[0])
		t.Fail()
	}

	if gotStrings[1] != "string1" {
		t.Logf("expected %q got %q", c.args[0], gotStrings[1])
		t.Fail()
	}

	if gotStrings[2] != "string2" {
		t.Logf("expected %q got %q", c.args[1], gotStrings[2])
		t.Fail()
	}

	if cmd.ran != 1 {
		t.Logf("cmd was run %d times, 1 expected", cmd.ran)
		t.Fail()
	}

}

func TestGitCallsProrperlyWithoutConfig(t *testing.T) {
	var gotStrings []string
	cmd := fakeCmd{0}
	cmdFunc := func(s []string) ExecCmd {
		gotStrings = s
		return &cmd
	}

	var gitCmd ExecCmd
	var err error
	c := Call{}
	if gitCmd, err = c.git(cmdFunc); err != nil {
		t.Logf("git failed, got error: %v", err)
		t.Fail()
	}
	gitCmd.Run()

	if len(gotStrings) != 0 {
		t.Logf("expected 0 got %d", len(gotStrings))
		t.Fail()
	}

	if cmd.ran != 1 {
		t.Logf("cmd was run %d times, 1 expected", cmd.ran)
		t.Fail()
	}

}
