// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package feature

import (
	"flag"
	"fmt"

	"github.com/juju/errors"
	"github.com/perrito666/got/cli"
	"github.com/perrito666/got/feature/newfeature"
	"github.com/perrito666/got/feature/work"
	"github.com/perrito666/got/git"
	"github.com/perrito666/got/registry"
)

var flagSet *flag.FlagSet

type config struct {
	abbreviateList bool
	interactive    bool
}

var callConfig = &config{}

func init() {
	if err := registry.RegisterNewCommand("feature", NewFeatureCommand); err != nil {
		panic(err)
	}
	flagSet = flag.NewFlagSet("feature", flag.ExitOnError)
	shortDescription := "show only the feature numbers omitting the target branches"
	flagSet.BoolVar(&callConfig.abbreviateList, "s", false, shortDescription)
	flagSet.BoolVar(&callConfig.abbreviateList, "short", false, shortDescription)
	flagSet.BoolVar(&callConfig.interactive, "i", false, "prompt the feature with a menu.")
	newfeature.Flags(flagSet)
}

// NewFeatureCommand is the constructor for the "feature" subcommand.
func NewFeatureCommand() (registry.Command, error) {
	return &Command{}, nil
}

var usageDoc = `
usage: 
got feature: prints this help
`

// Command provides the "feature" subcommand to got.
type Command struct {
}

// Run implements registry.Command
func (b *Command) Run(args []string) error {
	if len(args) == 0 {
		fmt.Print(usageDoc)
		return nil
	}
	subC := args[0]
	subArgs := args[1:]
	if err := flagSet.Parse(subArgs); err != nil {
		return errors.Annotate(err, "error parsing arguments")
	}
	switch subC {
	case "work":
		// TODO (perrito666) make command an interface
		// make command simpler and Flags instead
		w := work.Command{
			Args:        flagSet.Args(),
			Interactive: callConfig.interactive,
			Short:       callConfig.abbreviateList,
			UI:          &cli.UI{},
			NewGit:      git.New,
		}
		return w.Handle()
	case "new":
		n := newfeature.Command{
			Args:   flagSet.Args(),
			UI:     &cli.UI{},
			NewGit: git.New,
		}
		return n.Handle()
	}
	return nil
}
