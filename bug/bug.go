package bug

import (
	"flag"
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/perrito666/got/cli"
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
	if err := registry.RegisterNewCommand("bug", NewBugCommand); err != nil {
		panic(err)
	}
	flagSet = flag.NewFlagSet("bug", flag.ExitOnError)
	shortDescription := "show only the bug numbers omitting the target branches"
	flagSet.BoolVar(&callConfig.abbreviateList, "s", false, shortDescription)
	flagSet.BoolVar(&callConfig.abbreviateList, "short", false, shortDescription)
	flagSet.BoolVar(&callConfig.interactive, "i", false, "prompt the bug with a menu.")
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

  work [-s|--short]
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
	subC := args[0]
	subArgs := args[1:]
	if err := flagSet.Parse(subArgs); err != nil {
		return errors.Annotate(err, "error parsing arguments")
	}
	switch subC {
	case "work":
		return handleWorkSubCommand()

	}
	return nil
}

func handleWorkSubCommand() error {
	args := flagSet.Args()
	if len(args) == 0 {
		if callConfig.interactive {
			branch, err := workPickBug()
			if err != nil {
				return errors.Annotate(err, "error selecting a bug to work on")
			}
			if branch == "" {
				return nil
			}
			return errors.Annotatef(git.Checkout(branch), "cannot switch to branch %q", branch)
		}
		return workListSubCommand(callConfig.abbreviateList)
	}
	return nil
}

func utilListFixes() (map[string][]string, error) {
	c := &git.Config{
		SubCommand: "branch",
	}
	cmd, err := git.Git(c)
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

func workListSubCommand(short bool) error {
	fixes, err := utilListFixes()
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println("Available bugs and their target versions:")
	for bug, targets := range fixes {
		fmt.Println(bug)
		if short {
			continue
		}
		for _, target := range targets {
			fmt.Println(fmt.Sprintf("  - %s", target))
		}
	}
	return nil
}

func craftFixBranch(bugno, target string) string {
	return fmt.Sprintf("fix_%s_%s", target, bugno)
}

func workPickBug() (string, error) {
	fixes, err := utilListFixes()
	if err != nil {
		return "", errors.Trace(err)
	}
	choices := []string{}
	index := []string{}
	for bug, targets := range fixes {
		for _, target := range targets {
			choices = append(choices, fmt.Sprintf("%q (%s)", bug, target))
			index = append(index, craftFixBranch(bug, target))
		}
	}
	chosen, err := cli.ChoiceMenu(choices, true, -1)
	if err != nil {
		return "", errors.Annotate(err, "interactive bug choice failed")
	}
	if len(chosen) == 0 {
		return "", nil
	}
	return index[chosen[0]], nil
}
