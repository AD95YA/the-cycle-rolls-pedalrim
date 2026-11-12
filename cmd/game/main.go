package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"pedalrim/internal/config"
	"pedalrim/internal/game"
)

func main() {
	ebiten.SetWindowSize(config.ScreenWidth*config.WindowScale, config.ScreenHeight*config.WindowScale)
	ebiten.SetWindowTitle(config.GameTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g, err := game.NewGame()
	if err != nil {
		log.Fatalf("failed to initialize game: %v", err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("error running game loop: %v", err)
	}
}

