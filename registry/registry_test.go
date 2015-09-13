package registry

import "testing"

type fakeCommand struct {
	cmdName string
}

// Run satisfies Command interface
func (*fakeCommand) Run(_ []string) error {
	return nil
}

func TestRegistryAddsToGlobalRegistry(t *testing.T) {
	registry = make(map[string]CommandConstructor)
	if len(registry) != 0 {
		t.Fatalf("registry should be empty, instead has %d items: %v", len(registry), registry)
	}
	fakeConstructor := func() (Command, error) {
		return &fakeCommand{"test command"}, nil
	}

	RegisterNewCommand("command", fakeConstructor)
	if len(registry) != 1 {
		t.Logf("registry should have exactly 1 item, instead has %d items: %v", len(registry), registry)
		t.Fail()
	}
	command, ok := registry["command"]
	if !ok {
		t.Logf("the registered command is not present in the registry: %v", registry)
	}
	obtainedCommand, err := command()
	if err != nil {
		t.Fatal(err)
	}
	expectedCommand, err := fakeConstructor()
	if err != nil {
		t.Fatal(err)
	}

	if obtainedCommand.(*fakeCommand).cmdName != expectedCommand.(*fakeCommand).cmdName {
		t.Logf("registered and obtained command constructors differ")
		t.Fail()
	}
}

func TestCommandsReturnACopy(t *testing.T) {}

func TestCommandsAfterChangingRegistryRetursUpdated(t *testing.T) {}
