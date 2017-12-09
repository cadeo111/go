package main

/*

Goal:
input x,y,color
return isPossible, currentBoardToDisplay

starting from the beginning of the game each move a frame is saved
a frame being a stored form of the current board state

A. before each move is made it must be checked to be legal
x .5 check if spot is un ocupied
x 1. first we must check that move satisfies the rule of ko (board must change between moves)
2. second we should check if the place is next to any open space, if true then move is legal
3. if false we should check if it is next to any stones of the same color
4. if false then it is an illegal move
5. if true we must check to see if each separate same color stone touches whitespace
6. if not then we must check if it touches a stone of the same color
6a. if stone1 does touch stone2 of same color add id of stone1 to hashset
7. and repeat from 4 until all stones are checked
8. if no stone same color stone touches white space then the move is illegal


* when checking stones have a record of which positions have been checked (hashset maybe?)

B. after each move is made, touching other color stones must have liberties checked

1. last placed stone is checked for surrounding stones of opposite color
2. each stone found is then checked to find if it touches whitespace
3. if not then checked if it touches stone of same color




*/

func main() {
	//fmt.Println(sa)
	g := Game{size: 9}

}

//---------------

const size = 19

// stone is a representation of a spot on a board.
// it is derived from a frame.
type stone struct {
	x     int
	y     int
	color int // 0 none, -1 black, 1 white, 9 out of bounds
}

//frame is a representation of the board at a certain point in the game.
type frame struct {
	//size is the max x and y of a position on the board
	size int
	//board is a way the stone position is stored
	board []int
	//turn is the turn associated with the frame
	turn int
}

//getStone is the function to retrieve the data of a position in a frame as a stone.
func (f frame) getStone(x int, y int) (stone) {
	pos := ((f.size * f.size) - (y * f.size) + x) - f.size
	stone := stone{x, y, f.board[pos]}
	return stone;
}

func (f frame) getSurroundingStonesColor(s stone) (int, int, int, int) {
	top, left, bottom, right := 9, 9, 9, 9
	if s.x < size-1 {
		right = f.getStone(s.x+1, s.y).color
	}
	if s.x > 0 {
		left = f.getStone(s.x-1, s.y).color
	}
	if s.y < size-1 {
		top = f.getStone(s.x, s.y+1).color
	}
	if s.y > 0 {
		bottom = f.getStone(s.x, s.y-1).color
	}

	return top, right, bottom, left
	// returning 9 signifies out of bounds
}

//Game is the representation of the entire game and is how one should access everything
type Game struct {
	size int
	//frameRecord is the array of all the frames from each move in the game
	frameRecord []frame
	currentTurn int
}

//blankFrame returns a new blank frame based on the standard frame size of the game
func (g Game) blankFrame() frame {
	b := make([]int, g.size*g.size)
	return frame{g.size, b, 0}
}

//currentFrame get the current frame if it exists.
//if not create a new blank frame and append it to the frame record
func (g Game) currentFrame() frame {
	if len(g.frameRecord) > 0 {
		return g.frameRecord[len(g.frameRecord)-1]
	} else {
		f := g.blankFrame()
		g.frameRecord = append(g.frameRecord, f)
		return f
	}
}

func (g Game) lastFrame() frame {
	if len(g.frameRecord) > 1 {
		return g.frameRecord[len(g.frameRecord)-2]
	} else {
		return g.frameRecord[0]
	}
}

func (g Game) isMovePossible(s stone) bool {
	cf := g.currentFrame()
	lf := g.lastFrame()
	x, y, color := s.x, s.y, s.color
	oppositeColor := color * -1

	lfs := lf.getStone(x, y)

	isUnoccupied := func() bool {
		return cf.getStone(x, y).color == 0
	}

	violatesKo := func() bool {

		if len(g.frameRecord) < 2 {
			// if there is less than two frames then is impossible to violate ko
			return false

		} else if lfs.color == color {
			t, r, b, l := lf.getSurroundingStonesColor(s)

			if t == color || r == color || b == color || l == color {
				return false
			}else {
				return true
			}
		}else{
			// if the last frame didn't have a stone there or
			// if the last frame had a stone of the opposite color
			// then it is impossible to violate ko
			return false
		}

	}

}

//getCurrentFrame is the function to get the Frames Struct representation of a
//string representation of a Go Board

/*

func isPossible(x int, y int, color int, currentFrame []int, previousFrame []int) bool {
	oppositeColor := 0 - color

	// basic surrounding
	ct, cr, cb, cl := getSurroundingStones(x, y, currentFrame) //  current bottom, top, right, left

	if getStoneState(x, y, previousFrame) == color {
		return false
	}

	if ct == 9 && cr == 9 { //top and right out of bounds
		if cb == oppositeColor && cl == oppositeColor { // surround on bottom and left

			return false

		} else if cb != 0 && cl != 0 { // need to search around to find if connected stones

			return hasLiberty(x, y, color, currentFrame, 0)

		}

	} else if ct == 9 && cl == 9 {// top and left out of bounds
		if cb == oppositeColor && cr == oppositeColor {

			return false

		} else if cb != 0 && cr != 0 {

			return hasLiberty(x, y, color, currentFrame, 0)

		}

	} else if cb == 9 && cr == 9 { // bottom and right out of bounds
		if ct == oppositeColor && cl == oppositeColor { //surrounded on left and and top

			return false

		} else if ct != 0 && cl != 0 { // top and left are not empty

			return hasLiberty(x, y, color, currentFrame, 0)

		}
	} else if cb == 9 && cl == 9 { // bottom and left out of bounds
		if cb == oppositeColor && cr == oppositeColor {

			return false

		} else if ct != 0 && cr != 0 {

			return hasLiberty(x, y, color, currentFrame, 0)

		}

	} else if ct == 9 { // top is out of bounds
		if cb == oppositeColor && cr == oppositeColor && cl == oppositeColor { // bottom left and right sides are surrounded

			return false

		} else if cl != 0 && cr != 0 && cb != 0 {

			return hasLiberty(x, y, color, currentFrame, 0)

		}

	} else if cr == 9 {
		if cb == oppositeColor && ct == oppositeColor && cl == oppositeColor { // bottom left and top sides are surrounded

			return false

		} else if ct != 0 && cl != 0 && cb != 0 {

			return hasLiberty(x, y, color, currentFrame, 0)

		}

	} else if cb == 9 {
		if ct == oppositeColor && cr == oppositeColor && cl == oppositeColor { // top left and right sides are surrounded

			return false

		} else if cl != 0 && cr != 0 && ct != 0 {

			return hasLiberty(x, y, color, currentFrame, 0)

		}

	} else if cl == 9 {
		if cb == oppositeColor && cr == oppositeColor && ct == oppositeColor { // bottom top and right sides are surrounded

			return false

		} else if ct != 0 && cr != 0 && cb != 0 {

			return hasLiberty(x, y, color, currentFrame, 0)

		}
	}
	return true
}

func hasLiberty(x int, y int, color int, aFrame []int, exclude int) bool {


	// exclude a certain direction
	//0 = none, 1 = top, 2 = right, 3 = bottom, 4 = left

	oppositeColor := 0 - color

	t, r, b, l := getSurroundingStones(x, y, aFrame)

	if t == 0 || r == 0 || b == 0 || l == 0 {
		return true
	}
	if t == color && exclude != 1 {
		return hasLiberty(x+1, y, color, aFrame, 4)
	} else if r == color && exclude != 2 {
		return hasLiberty(x+1, y, color, aFrame, 4)
	} else if b == color && exclude != 3 {
		return hasLiberty(x+1, y, color, aFrame, 1)
	} else if l == color && exclude != 4 {
		return hasLiberty(x+1, y, color, aFrame, 2)
	}
	return false
}

func getSurroundingStones(x int, y int, aFrame []int) (int, int, int, int) {
	top, left, bottom, right := 9, 9, 9, 9
	if x < size - 1 {
		right = getStoneState(x+1, y, aFrame)
	}
	if x > 0 {
		left = getStoneState(x-1, y, aFrame)
	}
	if y < size - 1 {
		top = getStoneState(x, y+1, aFrame)
	}
	if y > 0 {
		bottom = getStoneState(x, y-1, aFrame)
	}

	return top, right, bottom, left
	// returning 9 signifies out of bounds
}*/

/*
12345678901234567890
"0000000000000000000
0002000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000
0000000000000000000"
 */
