package main

import (
	"fmt"
	"idle-token/ui"
	"log"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	zone "github.com/lrstanley/bubblezone/v2"
)

func initialModel() ui.Renderer {
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
	return ui.NewRenderer(codeString, rawCode, colorChars)
}

func main() {

	logFile, err := os.Create("out.log")
	if err != nil {
		os.Exit(1)
	}

	log.SetOutput(logFile)

	zone.NewGlobal()

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
