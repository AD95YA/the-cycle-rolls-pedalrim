package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"pedalrim/internal/game"
)

func init() {
	g, err := game.NewGame()
	if err != nil {
		panic(err)
	}
	mobile.SetGame(g)
}

// Dummy is an exported function ensuring package symbols can be bound by ebitenmobile.
func Dummy() {}

