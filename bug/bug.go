package bug

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/perrito666/got/git"
	"github.com/perrito666/got/registry"
)

func init() {
	if err := registry.RegisterNewCommand("bug", NewBugCommand); err != nil {
		panic(err)
	}
}

// NewBugCommand is the constructor for the "bug" subcommand.
func NewBugCommand() (registry.Command, error) {
	return &Command{}, nil
}

var usageDoc = `
usage: 
got bug: prints this help

bug assumes that you have a project that has one or several maintenance branches
such as: master, 1.2, 1.3
It allows you to:
  - Specify a base branch in which al bugs are by default worked (ie: if your default
branch is 1.2, "bug fix 488992" will do a "git -b fix_1.2_488992 1.2" and then co
that same branch for you to work on, if not you will be prompted for a choice.
  - Forward/Back port fixes from one branch to the other (With limitations)
  - Quickly switch to the default branch for a bug.

bug available subcommands:

  fix [-b base_branch] <bug id>
    will create a new branch to fix <bug id>, if base branch is provided it will
    be used as the base for the fix branch, otherwise there will be an interactive
    prompt.

  port [-m] [-b target_branch]
    will try to cherry pick the commits from this branch to the target if possible.
    if -m is specified the following will happen:
      - will checkout the origin branch.
      - will pull -r from the upstream branch.
      - will will be prompted to choose a commit to do the merge (all possible
      magic will be worked to try to do this, even if it is a github merge.

  work
    will list the bugs you can work on and the branches for wich you can do it.

  work [-b base_branch] <number>
    will checkout the branch for that bug, if base branch is specified it will
    try to checkout the fix branch for that maintenance branch oterwise the default
    specified will be checked out.

  link
    will print the link to the bug in the issue tracker: (current supported trackers
    are launchpad and github)
 
  default [base_work_branch]
    will set the given base branch as the default or prompt you for a new one.

`

// Command provides the "bug" subcommand to got.
type Command struct {
}

// Run implements registry.Command
func (b *Command) Run(args []string) error {
	if len(args) == 0 {
		fmt.Print(usageDoc)
		return nil
	}
	switch args[0] {
	case "work":
		if len(args) == 1 {
			return workListSubCommand()
		}
	}
	return nil
}

func utilListFixes() (map[string][]string, error) {
	c := &git.Config{
		SubCommand: "branch",
	}
	branchList, err := git.LibGit(c)
	if err != nil {
		return nil, errors.Annotatef(err, "could not obtain list of branches.")
	}
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

func workListSubCommand() error {
	fixes, err := utilListFixes()
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println("Available bugs and their target versions:")
	for bug, targets := range fixes {
		fmt.Println(bug)
		for _, target := range targets {
			fmt.Println(fmt.Sprintf("  - %s", target))
		}
	}
	return nil
}
