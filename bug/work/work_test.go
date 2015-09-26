// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package work

import (
	"testing"

	gtesting "github.com/perrito666/got/testing"
)

func TestHandleNoArgsNonInteractive(t *testing.T) {
	c := Command{
		Args:        []string{},
		Interactive: false,
		Short:       false,
		UI:          &gtesting.FakeUI{},
		NewGit:      gtesting.New,
	}
	c.Handle()
}

func TestHandleNoArgsInteractivePicksBug(t *testing.T) {
	c := Command{
		Args:        []string{},
		Interactive: false,
		Short:       false,
		UI:          &gtesting.FakeUI{},
		NewGit:      gtesting.New,
	}
	c.Handle()
}

func TestHandleNoArgsInteractiveDosntPicksBug(t *testing.T) {
	c := Command{
		Args:        []string{},
		Interactive: false,
		Short:       false,
		UI:          &gtesting.FakeUI{},
		NewGit:      gtesting.New,
	}
	c.Handle()
}
