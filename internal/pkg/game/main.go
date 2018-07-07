package game

import (

	"goGame/internal/pkg/stone"
	"goGame/internal/pkg/position"
	"goGame/internal/pkg/score"
	"log"
)



/*
 -1 is black
Goal:
input x,y,color
return isPossible, currentBoardToDisplay

starting from the beginning of the game each move a position is saved
a position being a stored form of the current board state

A. before each move is made it must be checked to be legal
x .5 check if spot is unoccupied
x 1. first we must check that move satisfies the rule of ko (board must change between moves)
x 2. second we should check if the place is next to any open space, if true then move is legal
x3. if false we should check if it is next to any stones of the same color
x4. if false then it is an illegal move
5. if true we must check to see if each separate same color stone touches whitespace
6. if not then we must check if it touches a stone of the same color
6a. if stone1 does touch stone2 of same color Add id of stone1 to hashset
7. and repeat from 4 until all stones are checked
8. if no stone same color stone touches white space then the move is illegal


* when checking stones have a record of which positions have been checked (hashset maybe?)

B. after each move is made, touching other color stones must have liberties checked

1. last placed stone is checked for surrounding stones of opposite color
2. each stone found is then checked to find if it touches whitespace
3. if not then checked if it touches stone of same color




*/

//---------------

//const size = 9

//Game is the representation of the entire game and is how one should access everything
const BLACK = -1
const WHITE = 1

type Game struct {
	size int
	//frameRecord is the array of all the frames from each move in the game
	frameRecord  []position.Frame
	CurrentTurn  int
	CurrentColor int // -1 black, 1 white
	GameOver     bool
	lastPass     int
	CapturedWhiteStones int
	CapturedBlackStones int
	sc           score.Score
}

//blankFrame returns a new blank position based on the standard position size of the game
func (g Game) blankFrame() position.Frame {
	return position.NewFrame(g.size, g.CurrentTurn)
}

//currentFrame Get the current position if it exists.
//if not create a new blank position and append it to the position record
func (g Game) currentFrame() position.Frame {
	if len(g.frameRecord) > 0 {
		return g.frameRecord[len(g.frameRecord)-1]
	} else {
		f := g.blankFrame()
		g.frameRecord = append(g.frameRecord, f)
		return f
	}
}

func (g Game) lastFrame() position.Frame {
	if len(g.frameRecord) > 1 {
		return g.frameRecord[len(g.frameRecord)-2]
	} else {
		return g.frameRecord[0]
	}
}

func (g Game) isMovePossible(s stone.Stone) bool {
	cf := g.currentFrame()
	lf := g.lastFrame()
	fwst,_ := cf.ProcessTakenStones(s) //position With Stones Taken
	x, y, color := s.X, s.Y, s.Color
	//oppositeColor := color * -1
	ct, cr, cb, cl := fwst.GetSurroundingStones(s)

	lfs := lf.GetStone(x, y)

	isOccupied := func() bool {
		return cf.GetStone(x, y).Color != 0
	}

	violatesKo := func() bool {

		if len(g.frameRecord) < 2 {
			// if there is less than two frames then is impossible to violate ko
			return false

		} else if lfs.Color == color {
			t, r, b, l := lf.GetSurroundingStones(s)

			if t.Color == color || r.Color == color || b.Color == color || l.Color == color {
				return false
			} else {
				return true
			}
		} else {
			// if the last position didn't have a stone there or
			// if the last position had a stone of the opposite color
			// then it is impossible to violate ko
			return false
		}

	}

	hasBlankNeighbor := func() bool {
		return ct.Color == 0 || cr.Color == 0 || cb.Color == 0 || cl.Color == 0
	}

	hasSameColorNeighbor := func() bool {
		if ct.Color == color || cr.Color == color || cb.Color == color || cl.Color == color {
			return true
		} else {
			return false
		}
	}
	//debug
	/*
		fmt.Println("isOccupied", isOccupied())
		fmt.Println("violatesKo", violatesKo())
		fmt.Println("hasBlankNeighbor", hasBlankNeighbor())
		fmt.Println("hasSameColorNeighbor", hasSameColorNeighbor())
		fmt.Println("SpotWouldHaveLiberty", cf.SpotWouldHaveLiberty(s))
	*/

	if violatesKo() {
		return false
	} else if isOccupied() {
		return false
	} else if hasBlankNeighbor() {
		return true
	} else if !hasSameColorNeighbor() {
		return false
	} else {
		return cf.SpotWouldHaveLiberty(s)
	}

} //has been validated with non exhaustive tests

func (g *Game) Move(x int, y int, color int) bool {
	return g.move(stone.Stone{x, y, color})
}

func (g *Game) move(s stone.Stone) bool {
	if g.isMovePossible(s) {
		f , removed := g.currentFrame().ProcessTakenStones(s)
		if(s.Color == BLACK){
			g.CapturedWhiteStones += removed
		}else{
			g.CapturedBlackStones += removed
		}
		g.frameRecord = append(g.frameRecord, f)
		g.CurrentTurn++
		g.lastPass = 0
		g.CurrentColor *= -1
		return true
	} else {
		return false
	}
}

func (Game) New(size int) Game {
	return Game{size: size, frameRecord: []position.Frame{position.NewFrame(size, 0)}, CurrentTurn: 1, CurrentColor: -1, GameOver: false, lastPass: 0 }
}

func (g *Game) End() {
	sc := score.Score{Fr: g.currentFrame(), Pm: position.Map{}.Init()}

	//sc.debug = true // for debugging the Score

	sc.CountAllStones()
	//fmt.Println("White: "+strconv.Itoa(sc.White), "Black: "+strconv.Itoa(sc.Black))
	sc.Black += g.CapturedWhiteStones
	sc.White += g.CapturedBlackStones
	g.sc = sc

}

func (g *Game) Pass(color int) bool {
	if color == g.CurrentColor {
		if g.lastPass == g.CurrentColor * -1 {
			g.GameOver = true
			return true
		}
		g.CurrentTurn ++
		g.lastPass = color
		g.CurrentColor *= -1
		return true
	} else {
		return false
	}
}

func (g *Game) Forfeit(color int) {
	if color == BLACK {
		g.sc = score.Score{White: 0, Black: 9001}
	} else if color == WHITE {
		g.sc = score.Score{Black: 0, White: 9001}
	} else {
		log.Fatal("Incorrect Forfeit Color")
	}
}

func (g Game) Print() {
	g.currentFrame().Print()
}

func (g Game) BoardString() string{
	return g.currentFrame().ToString()
}

