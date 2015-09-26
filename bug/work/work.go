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

		branch, err := w.PickBug()
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

// List prints a list of the existing bug branches listing
// the bug number and underneath all the target branches.
// If short is provided, the list will show only bug numbers.
// For a bug to show in this list, the branch name must
// be of the form fix_<target>_<bug id>
func (w *Command) List() error {
	fixes, err := util.ListFixes(w.NewGit)
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println("Available bugs and their target versions:")
	for bug, targets := range fixes {
		fmt.Println(bug)
		if w.Short {
			continue
		}
		for _, target := range targets {
			fmt.Println(fmt.Sprintf("  - %s", target))
		}
	}
	return nil
}

// PickBug will prompt the user which bug branch to
// checkout.
func (w *Command) PickBug() (string, error) {
	fixes, err := util.ListFixes(w.NewGit)
	if err != nil {
		return "", errors.Trace(err)
	}
	choices := []string{}
	index := []string{}
	for bug, targets := range fixes {
		for _, target := range targets {
			choices = append(choices, fmt.Sprintf("%q (%s)", bug, target))
			index = append(index, util.CraftFixBranch(bug, target))
		}
	}
	chosen, err := w.UI.ChoiceMenu(choices, true, -1)
	if err != nil {
		return "", errors.Annotate(err, "interactive bug choice failed")
	}
	if len(chosen) == 0 {
		return "", nil
	}
	return index[chosen[0]], nil
}
