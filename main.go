package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	codesection "idle-token/code-section"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
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
)

type PlayerInfo struct {
	width       int
	height      int
	tokens      int
	modifers    []int
	tps         int
	cursorPos   int
	RawCode     string
	CurrentCode codesection.CodeLib
}

func initialModel() PlayerInfo {
	filePath := "./data/code/python.txt"
	raw, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		os.Exit(1)
	}

	codeString := strings.TrimSpace(string(raw))
	runes := []rune(codeString)
	rawCode := make([]rune, len(runes)+1)
	colorChars := make([]rune, len(runes)+1)

	for i, ch := range runes {
		rawCode[i] = ch
		colorChars[i] = 'g'
	}

	return PlayerInfo{
		tokens:    10,
		modifers:  []int{1},
		tps:       1,
		cursorPos: 0,
		RawCode:   codeString,
		CurrentCode: codesection.CodeLib{
			RawChars:   rawCode,
			CharColors: colorChars,
			MaxTokens:  len(strings.Fields(codeString)),
		},
	}
}

func (p PlayerInfo) Init() tea.Cmd {
	return nil
}

func (p PlayerInfo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	rightChar := p.CurrentCode.RawChars[p.cursorPos]

	if unicode.IsSpace(rightChar) {
		p.cursorPos++
		return p, nil
	}

	color := 'e'
	cursorMove := 1

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "ctrl+x" {
			fmt.Println("\x1b[2J")
			return p, tea.Quit
		}

		if key == "backspace" {
			color = 'g'
			cursorMove = -1
		}

		if key != string(rightChar) {
			color = 'r'
		}
	}

	p.CurrentCode.CharColors[p.cursorPos] = color
	p.cursorPos += cursorMove
	return p, nil
}

func (p PlayerInfo) View() tea.View {
	s := ""
	runes := []rune(p.RawCode)

	for i, ch := range runes {

		if unicode.IsSpace(ch) {
			s += string(ch)
			continue
		}

		if i == p.cursorPos {
			s += cursorStyle.Render(string(ch))
			continue
		}

		rendered := ""
		switch p.CurrentCode.CharColors[i] {
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
	main := border.Render(s)

	rightTop := border.
		Width(14).
		Height(6).
		MarginLeft(2).
		Render("Placeholder")

	rightBottom := border.
		Width(28).
		MarginLeft(2).
		Render("Another\nplaceholder")
	rightBottom = border.
		Height(p.height - lipgloss.Height(rightBottom)).
		Render("Another\nPlaceholder")

	rightCol := lipgloss.JoinVertical(lipgloss.Left, rightTop, rightBottom)

	layout := lipgloss.JoinHorizontal(lipgloss.Top, main, rightCol)

	centeredLayout := lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(layout)

	return tea.NewView(centeredLayout)
}

func main() {

	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		fmt.Println("\nAn error occured getting the width and height of the terminal: ", err)
		return
	}

	app := tea.NewProgram(initialModel(), tea.WithWindowSize(width, height))
	if _, err := app.Run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
