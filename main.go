package main

import (
	"game/engine"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCREENWIDTH  int = 490
	SCREENHEIGHT int = 300
)

func main() {
	ebiten.SetWindowSize(SCREENWIDTH*2, SCREENHEIGHT*2)
	ebiten.SetWindowTitle("Golang_DibaloⅡ")
	ebiten.SetMaxTPS(80)
	gameStart := engine.NewGame()
	if err := ebiten.RunGame(gameStart); err != nil {
		log.Fatal(err)
	}
}
