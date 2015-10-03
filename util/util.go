// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package util

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/perrito666/got/git"
	"github.com/perrito666/got/interfaces"
)

const (
	// FeatureType is the branch prefix for features.
	FeatureType = "feature"
	// FixType is the branch prefix for bug fixes.
	FixType = "fix"
)

// PrintList will print a list of items.
func PrintList(newGit git.CompatibleConstructor, branchType string, short bool) error {
	fixes, err := ListBranches(newGit, branchType)
	if err != nil {
		return errors.Trace(err)
	}

	for branch, targets := range fixes {
		fmt.Println(branch)
		if short {
			continue
		}
		for _, target := range targets {
			fmt.Println(fmt.Sprintf("  - %s", target))
		}
	}
	return nil
}

// ListBranches returns a map containing all open bug fixes and a list
// of the target branches being develped for each.
func ListBranches(newGit git.CompatibleConstructor, branchType string) (map[string][]string, error) {
	c := newGit("branch", []string{"--list", "--no-color"})
	cmd, err := c.Git()
	if err != nil {
		return nil, errors.Annotate(err, "cannot create git command caller")
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Annotatef(err, "calling git %q failed: %v", c.SubCommand, out)
	}

	branchList := string(out)

	branches := strings.Split(branchList, "\n")
	foundBranches := make(map[string][]string)
	for _, branch := range branches {
		// current might be marked with an asterisk.
		if strings.HasPrefix(strings.TrimSpace(branch), "*") {
			branch = branch[1:]
		}

		if strings.HasPrefix(strings.TrimSpace(branch), fmt.Sprintf("%s_", branchType)) {
			parts := strings.SplitN(branch, "_", 3)
			// this is not one of ours.
			if len(parts) != 3 {
				continue
			}
			target := parts[1]
			branchName := parts[2]
			targets := foundBranches[branchName]
			targets = append(targets, target)
			foundBranches[branchName] = targets
		}
	}
	return foundBranches, nil
}

// CraftBranchName returns the name of a branch using the params.
func CraftBranchName(branchType, reference, target string) string {
	return fmt.Sprintf("%s_%s_%s", branchType, target, reference)
}

// Picker presents a choice between branches of a type
func Picker(branchType string, newGit git.CompatibleConstructor, ui interfaces.UI) (string, error) {
	fixes, err := ListBranches(newGit, branchType)
	if err != nil {
		return "", errors.Trace(err)
	}
	choices := []string{}
	index := []string{}
	for bug, targets := range fixes {
		for _, target := range targets {
			choices = append(choices, fmt.Sprintf("%q (%s)", bug, target))
			index = append(index, CraftBranchName(branchType, bug, target))
		}
	}
	chosen, err := ui.ChoiceMenu(choices, true, -1)
	if err != nil {
		return "", errors.Annotate(err, "interactive branch choice failed")
	}
	if len(chosen) == 0 {
		return "", nil
	}
	return index[chosen[0]], nil
}

// NewBranch creates a new branch from the given target.
func NewBranch(newGit git.CompatibleConstructor, branchType, target, name string) (string, error) {
	branch := CraftBranchName(branchType, name, target)
	c := newGit("checkout", []string{"-b", branch, target})
	cmd, err := c.Git()
	if err != nil {
		return "", errors.Annotate(err, "cannot create git command caller")
	}

	out, err := cmd.Output()
	if err != nil {
		return "", errors.Annotatef(err, "calling git %q %v failed: %v", c.SubCommand(), c.Args(), out)
	}
	printable := string(out)
	if printable != "" {
		fmt.Println(printable)
	}
	return branch, nil
}
