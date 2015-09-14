package cli

import (
	"github.com/juju/errors"
	gc "github.com/rthornton128/goncurses"
)

// ChoiceMenu presents the user with a set of choices to toggle using
// a ncurses togle-able menu.
func ChoiceMenu(items []string, one bool, curent int) ([]string, error) {
	stdscr, err := gc.Init()
	defer gc.End()
	if err != nil {
		return nil, errors.Trace(err)
	}

	menuItems := make([]*gc.MenuItem, len(items))
	for i, val := range items {
		menuItems[i], _ = gc.NewItem(val, "")
		defer menuItems[i].Free()
	}

	// create the menu
	menu, err := gc.NewMenu(menuItems)
	defer menu.Free()
	if err != nil {
		return nil, errors.Trace(err)
	}

	menu.Option(gc.O_ONEVALUE, one)

	menu.Post()
	defer menu.UnPost()

	for {
		gc.Update()
		ch := stdscr.GetChar()

		switch ch {
		// This is Esc at least in linux vtx.... why is it not a const?
		case 27:
			return nil, nil
		case ' ':
			menu.Driver(gc.REQ_TOGGLE)
		case gc.KEY_RETURN, gc.KEY_ENTER:
			var list []string
			for _, item := range menu.Items() {
				if item.Value() {
					list = append(list, item.Name())
				}
			}
			return list, nil
		default:
			menu.Driver(gc.DriverActions[ch])
		}
	}
}
