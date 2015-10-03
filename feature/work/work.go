// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package work

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/perrito666/got/git"
	"github.com/perrito666/got/interfaces"
	"github.com/perrito666/got/util"
)

// Command holds the configuration and methods required for the work
// sub-command.
type Command struct {
	Args        []string
	Interactive bool
	Short       bool
	UI          interfaces.UI
	NewGit      git.CompatibleConstructor
}

// Handle is the entry point for the Work sub-command.
func (w *Command) Handle() error {
	if len(w.Args) == 0 {
		if !w.Interactive {
			return w.List()
		}

		branch, err := w.PickFeature()
		if err != nil {
			return errors.Annotate(err, "could not select a bug to work on")
		}

		// no branch chosen, user most likely hit Esc.
		if branch == "" {
			return nil
		}
		g := git.Call{}

		if err := g.Checkout(branch); err != nil {
			return errors.Annotatef(err, "cannot switch to branch %q", branch)
		}
		return nil
	}

	// Handle the args, particularly 0 which should be the subcommand
	return nil
}

// List prints a list of the existing feature branches listing
// the feature name and underneath all the target branches.
// If short is provided, the list will show only features names.
// For a feature to show in this list, the branch name must
// be of the form feature_<target>_<feature name>
func (w *Command) List() error {
	fmt.Println("Available bugs and their target versions:")
	return util.PrintList(w.NewGit, util.FeatureType, w.Short)
}

// PickFeature will prompt the user which bug branch to
// checkout.
func (w *Command) PickFeature() (string, error) {
	return util.Picker(util.FeatureType, w.NewGit, w.UI)
}
