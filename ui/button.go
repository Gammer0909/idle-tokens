package ui

type Button struct {
	Label   string
	pressed bool
}

func NewButton(label string) Button {
	return Button{
		Label:   label,
		pressed: false,
	}
}

func (b Button) Press() {
	b.pressed = true
}

func (b Button) Release() {
	b.pressed = false
}
