package ui

type Menu struct {
	MenuName string
	Buttons  []Button
	MenuOpen bool
}

func NewMenu(name string, buttons []Button) Menu {
	return Menu{
		MenuName: name,
		Buttons:  buttons,
	}
}

func (m Menu) MenuView(r Renderer) string {

	s := "Yooo wsp"

	menu := border.Render(s)

	return menu

}
