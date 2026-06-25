package ui

import "github.com/lxn/walk"

func loadAppIcon() (*walk.Icon, error) {
	for _, id := range []int{3, 2, 1} {
		icon, err := walk.NewIconFromResourceId(id)
		if err == nil && icon != nil {
			return icon, nil
		}
	}
	for _, name := range []string{"APPICON", "IDI_ICON1"} {
		icon, err := walk.NewIconFromResource(name)
		if err == nil && icon != nil {
			return icon, nil
		}
	}
	return walk.NewIconFromSysDLL("shell32", 2)
}
