package model

import (
	codesection "idle-token/code-section"
	"strings"
)

type PlayerInfo struct {
	tokens      int
	modifers    []int
	tps         int
	RawCode     string
	CurrentCode codesection.CodeLib
}

func NewPlayer(codeString string, rawCode []rune, colorChars []rune) PlayerInfo {
	return PlayerInfo{
		tokens:   10,
		modifers: []int{1},
		RawCode:  codeString,
		CurrentCode: codesection.CodeLib{
			RawChars:   rawCode,
			CharColors: colorChars,
			MaxTokens:  len(strings.Fields(codeString)),
		},
	}
}
