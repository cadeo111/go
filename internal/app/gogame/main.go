package gogame

import (
	gm "goGame/internal/pkg/game"
)

const BLACK = gm.BLACK
const WHITE = gm.WHITE

type UserInterface interface {
	RequestNextMove(color int, board string) (kind int, x int, y int) // gives color, expects(1:move, 2:pass, 3:forfeit), x,y
	ShowBoard(board string, capBlack int, capWhite int)
	Message(error bool, kind string, message string)
}
type GoGame struct {
	game gm.Game
}

func Init(size int, handicap int) GoGame {
	return GoGame{gm.Game{}.New(9)}
}

func (g *GoGame) Run(ui UserInterface) {
	for !g.game.GameOver {
		color := g.game.CurrentColor
		kind, x, y := ui.RequestNextMove(color, g.game.BoardString())
		switch kind {
		case 1:
			if g.game.Move(x, y, color) {
				ui.Message(false, "moveComplete", "move completed")
			} else {
				ui.Message(true, "moveIllegal", "move not completed")
			}
		case 2:
			g.game.Pass(color)
			ui.Message(false, "movePassed", "turn passed")
		case 3:
			g.game.Forfeit(color)
			ui.Message(false, "moveForfeit", "game forfeited")
		}
		ui.ShowBoard(g.game.BoardString(), g.game.CapturedBlackStones, g.game.CapturedWhiteStones)
	}
	g.game.End()

	ui.Message(false, "gameOver", "game is over")
}
