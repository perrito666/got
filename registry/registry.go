package registry

import (
	"github.com/juju/errors"
)

// Command represents one of the possible sub-commands of got.
type Command interface {
	Run([]string) error
}

// CommandConstructor is the type of the function you should implement to have your
// command be instantiated by the main thread..
type CommandConstructor func() (Command, error)

var registry map[string]CommandConstructor

func init() {
	registry = make(map[string]CommandConstructor)
}

// RegisterNewCommand will add a command for got to handle when called.
func RegisterNewCommand(name string, constructor CommandConstructor) error {
	existing, ok := registry[name]
	if ok {
		return errors.AlreadyExistsf("command %q already registered for %v", name, existing)
	}
	registry[name] = constructor
	return nil
}

// Commands returns a copy of the command registry.
func Commands() map[string]CommandConstructor {
	registryCopy := make(map[string]CommandConstructor, len(registry))
	for name, constructor := range registry {
		registryCopy[name] = constructor
	}
	return registryCopy
}
