package ui

type Menu struct {
	Buttons  []Button
	MenuOpen bool
}

func NewMenu(buttons []Button) Menu {
	return Menu{
		Buttons: buttons,
	}
}
