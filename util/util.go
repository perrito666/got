// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package util

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/perrito666/got/git"
)

// ListFixes returns a map containing all open bug fixes and a list
// of the target branches being develped for each.
func ListFixes(newGit git.CompatibleConstructor) (map[string][]string, error) {
	c := newGit("branch", nil)
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
	fixes := make(map[string][]string)
	for _, branch := range branches {
		if strings.HasPrefix(strings.TrimSpace(branch), "fix_") {
			parts := strings.SplitN(branch, "_", 3)
			// this is not one of ours.
			if len(parts) != 3 {
				continue
			}
			target := parts[1]
			bugno := parts[2]
			targets := fixes[bugno]
			targets = append(targets, target)
			fixes[bugno] = targets
		}
	}
	return fixes, nil
}

// CraftFixBranch returns the name of a bug fix branch using the bug
// number and target branch.
func CraftFixBranch(bugno, target string) string {
	return fmt.Sprintf("fix_%s_%s", target, bugno)
}
