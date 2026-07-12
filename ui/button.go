package ui

import tea "charm.land/bubbletea/v2"

type Button struct {
	Label      string
	X, Y, W, H int
	pressed    bool
}

func NewButton(label string, x int, y int, w int, h int) Button {
	return Button{
		Label: label,
		X:     x,
		Y:     y,
		W:     w,
		H:     h,
	}
}

func (b Button) Update(msg tea.Msg) (Button, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Mouse().Button == tea.MouseLeft &&
			msg.Mouse().X >= b.X &&
			msg.Mouse().Y < b.X+b.W &&
			msg.Mouse().Y >= b.Y &&
			msg.Mouse().Y < b.Y+b.H {
			b.pressed = true
			return b, func() tea.Msg { return ButtonClickedMsg{Label: b.Label} }
		}
	}
}

type ButtonClickedMsg struct {
	Label string
}
