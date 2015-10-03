// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package newfeature

import (
	"flag"
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/perrito666/got/git"
	"github.com/perrito666/got/interfaces"
	"github.com/perrito666/got/util"
)

type config struct {
	featureName   string
	featureTarget string
}

var callConfig = &config{}

// Flags sets the flags for this feature.
func Flags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&callConfig.featureName, "name", "", "the name of the feature, will be sanitized to name the branch.")
	flagSet.StringVar(&callConfig.featureTarget, "target", "", "the target of the feature, this branch should already exist.")
}

// Command holds the configuration and methods required for the new
// sub-command.
type Command struct {
	Args   []string
	UI     interfaces.UI
	NewGit git.CompatibleConstructor
}

// Handle is the entry point for the Work sub-command.
func (w *Command) Handle() error {
	if callConfig.featureName == "" && callConfig.featureTarget == "" {

		// TODO(perrito666) create an interactive mode for feature
		fmt.Println("please specify target branch and branch name")
		return nil
	}
	name := callConfig.featureName
	for _, invalidChar := range []string{",", "-", ".", " "} {
		name = strings.Replace(name, invalidChar, "_", -1)
	}
	created, err := util.NewBranch(w.NewGit, util.FeatureType, callConfig.featureTarget, name)
	if err != nil {
		return errors.Annotate(err, "cannot create new branch")
	}
	fmt.Printf("now working in %q \n", created)
	return nil
}
