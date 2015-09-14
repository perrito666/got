package cli

import (
	"github.com/juju/errors"
	gc "github.com/rthornton128/goncurses"
)

// ChoiceMenu presents the user with a set of choices to toggle using
// a ncurses togle-able menu.
// TODO(perrito) current should be a list and used
func ChoiceMenu(items []string, one bool, curent int) ([]int, error) {
	stdscr, err := gc.Init()
	defer gc.End()
	if err != nil {
		return nil, errors.Trace(err)
	}

	gc.StartColor()
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)

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

	// This in true works very weirdly.
	menu.Option(gc.O_ONEVALUE, false)
	menu.Option(gc.O_SHOWDESC, true)

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
			if !one {
				menu.Driver(gc.REQ_TOGGLE)
			}
		case gc.KEY_RETURN, gc.KEY_ENTER:
			if one {
				menu.Driver(gc.REQ_TOGGLE)
			}
			var list []int
			for i, item := range menu.Items() {
				if item.Value() {
					list = append(list, i)
				}
			}
			return list, nil
		default:
			menu.Driver(gc.DriverActions[ch])
		}
	}
}
