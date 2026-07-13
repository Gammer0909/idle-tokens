package ui

import (
	"fmt"
	"idle-token/model"
	"log"
	"os"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone/v2"
)

var (
	grayStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B3B3B"))
	redStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F24646"))
	greenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#82ED7B"))
	border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(2, 2)
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#F2F2F2"))
	colorBg     = lipgloss.Color("#1E1E2E")
	colorFg     = lipgloss.Color("#CDD6F4")
	colorAccent = lipgloss.Color("#F5C2E7")
	colorGreen  = lipgloss.Color("#A6E3A1")
	colorRed    = lipgloss.Color("#F38BA8")
	colorYellow = lipgloss.Color("#F9E2AF")
	colorMuted  = lipgloss.Color("#6C7086")
)

type Renderer struct {
	width     int
	height    int
	cursorPos int
	Button    Button
	Player    model.PlayerInfo
}

func NewRenderer(codeString string, rawCode []rune, colorChars []rune) Renderer {
	player := model.NewPlayer(codeString, rawCode, colorChars)
	button := NewButton("Open Menu")

	return Renderer{
		cursorPos: 0,
		Button:    button,
		Player:    player,
	}
}

func (r Renderer) Init() tea.Cmd {
	return nil
}

func (r Renderer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	rightChar := r.Player.CurrentCode.RawChars[r.cursorPos]

	if unicode.IsSpace(rightChar) {
		r.cursorPos++
		return r, nil
	}

	color := 'e'
	cursorMove := 1

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		return r, nil
	case tea.MouseReleaseMsg:
		if msg.Button != tea.MouseLeft {
			return r, nil
		}

		if zone.Get("menu").InBounds(msg) {
			r.Button.Press()
		}
	case tea.KeyPressMsg:
		key := msg.String()
		log.Printf("Pressed key: %s", key)
		if key == "ctrl+c" || key == "ctrl+x" || key == "esc" {
			fmt.Println("\x1b[2J")
			return r, tea.Quit
		}

		if key == "backspace" {
			color = 'g'
			if unicode.IsSpace(r.Player.CurrentCode.RawChars[r.cursorPos-1]) {
				r.cursorPos -= 2
				r.Player.CurrentCode.CharColors[r.cursorPos] = color
				return r, nil
			}
			r.Player.CurrentCode.CharColors[r.cursorPos] = color
			r.cursorPos--
			return r, nil
		}

		if key != string(rightChar) {
			color = 'r'
		}
		r.Player.CurrentCode.CharColors[r.cursorPos] = color
		r.cursorPos += cursorMove
		return r, nil
	}

	return r, nil
}

func (r Renderer) View() tea.View {
	s := ""
	runes := []rune(r.Player.RawCode)

	for i, ch := range runes {

		if unicode.IsSpace(ch) {
			s += string(ch)
			continue
		}

		if i == r.cursorPos {
			s += cursorStyle.Render(string(ch))
			continue
		}

		rendered := ""
		switch r.Player.CurrentCode.CharColors[i] {
		case 'g':
			rendered = grayStyle.Render(string(ch))
		case 'r':
			rendered = redStyle.Render(string(ch))
		case 'e':
			rendered = greenStyle.Render(string(ch))
		default:
			rendered = grayStyle.Render(string(ch))
		}

		if rendered == "" {
			fmt.Println("\nError: rendered was an empty string.")
			os.Exit(1)
		}

		s += rendered

	}
	// Main View setup
	main := border.
		BorderForeground(colorAccent).
		Render(s)

	rightTop := border.
		BorderForeground(colorAccent).
		Width(28).
		Height(6).
		Render("Open Menu")

	rightBottom := border.
		BorderForeground(colorAccent).
		Width(28).
		Height((lipgloss.Height(main) - 2) - lipgloss.Height(rightTop)).
		Render("Another\nplaceholder")

	rightCol := lipgloss.JoinVertical(lipgloss.Left, zone.Mark("menu", rightTop), rightBottom)

	layout := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(colorFg).
		Padding(1, 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, main, rightCol))

	centeredLayout := lipgloss.NewStyle().
		Foreground(colorFg).
		Width(r.width).
		Height(r.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(layout)

	var view tea.View
	// Ensure that alt-screen is enabled, as bubblezone will only work in alt-screen mode.
	view.AltScreen = true
	// Enable mouse motion tracking.
	view.MouseMode = tea.MouseModeCellMotion
	// Wrap view in [zone.Scan].
	view.SetContent(zone.Scan(centeredLayout))

	return view
}
