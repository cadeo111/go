package main

import "goGame/internal/app/game"

func main() {
	//fmt.Println(sa)
	g := game.Game{}.New(9, 0)
	g.Move(0, 0, game.WHITE)
	g.Move(1, 0, game.WHITE)
	g.Move(2, 0, game.WHITE)
	g.Move(2, 1, game.WHITE)
	g.Move(2, 2, game.WHITE)
	g.Move(1, 2, game.WHITE)
	g.Move(0, 2, game.WHITE)
	g.Move(0, 1, game.WHITE)

	g.Move(0 + 5, 0, game.BLACK)
	g.Move(1 + 5, 0, game.BLACK)
	g.Move(2 + 5, 0, game.BLACK)
	g.Move(3 + 5, 3, game.BLACK)
	g.Move(3 + 5, 1, game.BLACK)
	g.Move(2 + 5, 2, game.BLACK)
	g.Move(1 + 5, 2, game.BLACK)
	g.Move(0 + 5, 2, game.BLACK)
	g.Move(0 + 5, 1, game.BLACK)

	g.Print()
	g.End()

}

