// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package interfaces

// UI reprsents a CLI/GUI
type UI interface {
	ChoiceMenu(items []string, one bool, curent int) ([]int, error)
}
