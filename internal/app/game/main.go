package game

import (
	"strconv"
	"fmt"
	"bytes"
	"log"
	"time"
)

//TODO refactor color into constants
//TODO separate into packages

/*
 -1 is black
Goal:
input x,y,color
return isPossible, currentBoardToDisplay

starting from the beginning of the game each move a frame is saved
a frame being a stored form of the current board state

A. before each move is made it must be checked to be legal
x .5 check if spot is unoccupied
x 1. first we must check that move satisfies the rule of ko (board must change between moves)
x 2. second we should check if the place is next to any open space, if true then move is legal
x3. if false we should check if it is next to any stones of the same color
x4. if false then it is an illegal move
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
	g := Game{}.new(9)
	g.Move(stone{0, 0, 1})
	g.Move(stone{1, 0, 1})
	g.Move(stone{2, 0, 1})
	g.Move(stone{2, 1, 1})
	g.Move(stone{2, 2, 1})
	g.Move(stone{1, 2, 1})
	g.Move(stone{0, 2, 1})
	g.Move(stone{0, 1, 1})

	g.Move(stone{0 + 5, 0, -1})
	g.Move(stone{1 + 5, 0, -1})
	g.Move(stone{2 + 5, 0, -1})
	g.Move(stone{2 + 5, 1, -1})
	g.Move(stone{2 + 5, 2, -1})
	g.Move(stone{1 + 5, 2, -1})
	g.Move(stone{0 + 5, 2, -1})
	g.Move(stone{0 + 5, 1, -1})

	g.currentFrame().print()

	g.end()

}

//---------------

//const size = 9

type posMap struct {
	m map[string]bool
}

func (p posMap) add(x int, y int) {
	str := strconv.Itoa(x) + "," + strconv.Itoa(y)
	p.m[str] = true
}

func (p posMap) get(x int, y int) bool {
	str := strconv.Itoa(x) + "," + strconv.Itoa(y)
	return p.m[str]
}

func (p posMap) stoneGet(s stone) bool {
	return p.get(s.x, s.y)
}

func (posMap) init() posMap {
	return posMap{m: make(map[string]bool)}
}

// stone is a representation of a spot on a board.
// it is derived from a frame.
type stone struct {
	x     int
	y     int
	color int // 0 none, -1 black, 1 white, 9 out of bounds
}

func (s stone) zero() stone {
	s.color = 0
	return s
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

func (f frame) print() {
	printAline := func(s []int) string {
		var buffer bytes.Buffer
		for i := 0; i < len(s); i++ {

			switch s[i] {
			case -1:
				buffer.WriteString("b")
			case 0:
				buffer.WriteString("•")
			case 1:
				buffer.WriteString("w")
			case 103: // for debuging
				buffer.WriteString("c")
			default:
				buffer.WriteString(strconv.Itoa(s[i]))
			}
			if i != len(s)-1 {
				buffer.WriteString(" ")
			}

		}
		return buffer.String()
	}
	for i := 0; i < f.size*2; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
	for i := 0; i < f.size; i++ {
		//fmt.Println(i*f.size," ",(i+1)*f.size)
		fmt.Println(printAline(f.board[i*f.size : (i+1)*f.size]))
		//fmt.Println(f.board[i*f.size:(i+1)*f.size])
	}
	for i := 0; i < f.size*2; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
}

//getStone is the function to retrieve the data of a position in a frame as a stone.
func (f frame) getStone(x int, y int) (stone) {
	pos := ((f.size * f.size) - (y * f.size) + x) - f.size
	stone := stone{x, y, f.board[pos]}
	return stone
}

func (f *frame) setStone(s stone) {
	pos := ((f.size * f.size) - (s.y * f.size) + s.x) - f.size
	f.board[pos] = s.color
}

func (f frame) getSurroundingStones(s stone) (stone, stone, stone, stone) {
	top, left, bottom, right := 9, 9, 9, 9
	if s.x+1 < f.size {
		right = f.getStone(s.x+1, s.y).color
	}
	if s.x > 0 {
		left = f.getStone(s.x-1, s.y).color
	}
	if s.y+1 < f.size {
		top = f.getStone(s.x, s.y+1).color
	}
	if s.y > 0 {
		bottom = f.getStone(s.x, s.y-1).color
	}

	return stone{s.x, s.y + 1, top},
		stone{s.x + 1, s.y, right},
		stone{s.x, s.y - 1, bottom},
		stone{s.x - 1, s.y, left}
	// returning 9 signifies out of bounds
}

func (f frame) processTakenStones(s stone) frame {
	nf := f.duplicate()
	nf.setStone(s)
	t, r, b, l := f.getSurroundingStones(s)

	oppositeColor := s.color * -1

	//first process opposite color
	if t.color == oppositeColor {
		if !nf.spotWouldHaveLiberty(t) {
			nf.removeConnected(t)
		}
		fmt.Println("t", nf.spotWouldHaveLiberty(t))
	}
	if r.color == oppositeColor {
		if !nf.spotWouldHaveLiberty(r) {
			nf.removeConnected(r)
		}
		fmt.Println("r", nf.spotWouldHaveLiberty(r))
	}
	if b.color == oppositeColor {
		if !nf.spotWouldHaveLiberty(b) {
			nf.removeConnected(b)
		}
		fmt.Println("b", nf.spotWouldHaveLiberty(b))
	}
	if l.color == oppositeColor {
		if !nf.spotWouldHaveLiberty(l) {
			nf.removeConnected(l)
		}
		fmt.Println("l", nf.spotWouldHaveLiberty(l))
	}
	return nf

}

func (f frame) removeConnected(s stone) {
	f.setStone(s.zero())
	t, r, b, l := f.getSurroundingStones(s)

	if t.color == s.color {
		f.removeConnected(t)
	}
	if r.color == s.color {
		f.removeConnected(r)
	}
	if b.color == s.color {
		f.removeConnected(b)
	}
	if l.color == s.color {
		f.removeConnected(l)
	}

}

func (f frame) spotWouldHaveLiberty(s stone) bool {
	ct, cr, cb, cl := f.getSurroundingStones(s)
	x, y, color := s.x, s.y, s.color
	pos := posMap{}.init()
	pos.add(x, y)

	top, right, bottom, left := false, false, false, false

	if ct.color == 0 || cr.color == 0 || cb.color == 0 || cl.color == 0 {
		return true
	}

	if ct.color == color {
		top = checkEachStone(x, y+1, color, f, pos)
	}
	if cr.color == color {
		right = checkEachStone(x+1, y, color, f, pos)
	}
	if cb.color == color {
		left = checkEachStone(x, y-1, color, f, pos)
	}
	if cl.color == color {
		bottom = checkEachStone(x-1, y, color, f, pos)
	}
	return top || right || bottom || left
}

func checkEachStone(x int, y int, color int, cf frame, pos posMap) bool {
	fmt.Println(x, y, pos.get(x, y))
	if pos.get(x, y) {
		return false
	} else {
		pos.add(x, y)
	}
	ct, cr, cb, cl := cf.getSurroundingStones(stone{x: x, y: y})
	fmt.Println(ct.color, cr.color, cb.color, cl.color)
	if (ct.color == 0 && !pos.get(ct.x, ct.y)) ||
		(cr.color == 0 && !pos.get(cr.x, cr.y)) ||
		(cb.color == 0 && !pos.get(cb.x, cb.y)) ||
		(cl.color == 0 && !pos.get(cl.x, cl.y)) {
		return true
	} else {
		top, right, bottom, left := false, false, false, false
		if ct.color == color {
			top = checkEachStone(x, y+1, color, cf, pos)
		}
		if cr.color == color {
			right = checkEachStone(x+1, y, color, cf, pos)
		}
		if cb.color == color {
			left = checkEachStone(x, y-1, color, cf, pos)
		}
		if cl.color == color {
			bottom = checkEachStone(x-1, y, color, cf, pos)
		}
		return top || right || bottom || left
	}
}

func (f frame) duplicate() frame {
	nf := frame{f.size, make([]int, f.size*f.size), f.turn}
	copy(nf.board, f.board)
	return nf
}

func newFrame(size int, turn int) frame {
	b := make([]int, size*size)
	return frame{size, b, turn}
}

type stoneStack struct {
	slice []stone
	size  int
}

func (st *stoneStack) push(sn stone) {
	st.slice = append(st.slice, sn)
	st.size += 1
}
func (st *stoneStack) pushSlice(ss []stone) {
	for _, s := range ss {
		st.push(s)
	}
}

func (st *stoneStack) peek() stone {
	return st.slice[len(st.slice)-1] // view last item in stack which was added last
}

func (st *stoneStack) pop() stone {
	if st.size > 0 {
		var r = st.peek()
		st.slice = st.slice[0 : len(st.slice)-1] // remove last item in slice the one that was peeked
		st.size -= 1
		return r
	} else {
		return stone{}
	}
}

//Game is the representation of the entire game and is how one should access everything
type Game struct {
	size int
	//frameRecord is the array of all the frames from each move in the game
	frameRecord  []frame
	currentTurn  int
	currentColor int // -1 black, 1 white
	gameOver     bool
	lastPass     int
	score        score
}

//blankFrame returns a new blank frame based on the standard frame size of the game
func (g Game) blankFrame() frame {
	return newFrame(g.size, g.currentTurn)
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
	fwst := cf.processTakenStones(s) //frame With Stones Taken
	x, y, color := s.x, s.y, s.color
	//oppositeColor := color * -1
	ct, cr, cb, cl := fwst.getSurroundingStones(s)

	lfs := lf.getStone(x, y)

	isOccupied := func() bool {
		return cf.getStone(x, y).color != 0
	}

	violatesKo := func() bool {

		if len(g.frameRecord) < 2 {
			// if there is less than two frames then is impossible to violate ko
			return false

		} else if lfs.color == color {
			t, r, b, l := lf.getSurroundingStones(s)

			if t.color == color || r.color == color || b.color == color || l.color == color {
				return false
			} else {
				return true
			}
		} else {
			// if the last frame didn't have a stone there or
			// if the last frame had a stone of the opposite color
			// then it is impossible to violate ko
			return false
		}

	}

	hasBlankNeighbor := func() bool {
		return ct.color == 0 || cr.color == 0 || cb.color == 0 || cl.color == 0
	}

	hasSameColorNeighbor := func() bool {
		if ct.color == color || cr.color == color || cb.color == color || cl.color == color {
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
		fmt.Println("spotWouldHaveLiberty", cf.spotWouldHaveLiberty(s))
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
		return cf.spotWouldHaveLiberty(s)
	}

} //has been validated with non exhaustive tests

func (g *Game) Move(s stone) bool {
	if g.isMovePossible(s) {
		f := g.currentFrame().processTakenStones(s)
		g.frameRecord = append(g.frameRecord, f)
		g.currentTurn++
		g.lastPass = 0
		g.currentColor *= -1
		return true
	} else {
		return false
	}
}

func (Game) new(size int) Game {
	return Game{size, []frame{newFrame(size, 0)}, 1, -1, false, 0}
}


type score struct {
	fr         frame  // current frame when score is calculated
	pm         posMap // posmap When score is calculated
	stack      stoneStack
	blackScore int
	whiteScore int
	debug      bool

	/*
	1. find a position that is blank
	2. check if it is next to a black stone
	3. check if it is next to a white stone
	4. check if it is next to any blank stones
	5. add all blank stones to stack
	6. add stone to recorded positions
	7. go to next blank stone in stack
	8. if no more stones in stack return
	*/
}

func (sc *score) checkNeighbors(x int, y int, color int) (bool) {
	ts, rs, bs, ls := sc.fr.getSurroundingStones(stone{x, y, color})

	if ts.color == color {
		return true
	}
	if rs.color == color {
		return true
	}
	if bs.color == color {
		return true
	}
	if ls.color == color {
		return true
	}
	return false
}
func (sc *score) findAllNeighborsBlank(x int, y int) ([]stone) {

	ts, rs, bs, ls := sc.fr.getSurroundingStones(stone{x, y, 0})
	var ss = []stone{ts, rs, bs, ls}
	var ret []stone
	for _, val := range ss {
		if val.color == 0 {
			ret = append(ret, val)
		}
	}
	return ret

}
func (sc *score) countOneStone(cs stone, touchedColor int, currentScore int) (int, int) { // score, color
	y := cs.y
	x := cs.x
	var next stone = stone{0, 0, 99}
	if !sc.pm.stoneGet(cs) {
		sc.fr.setStone(stone{x, y, 103})

		if sc.debug {
			sc.fr.print()
		}

		duration := time.Millisecond * 250
		time.Sleep(duration)

		currentScore += 1

		sc.pm.add(x, y) // save that stone has been checked

		blankNeighbors := sc.findAllNeighborsBlank(x, y) // find all neighbors that are blank
		if len(blankNeighbors) > 0 {
			next = blankNeighbors[0]
			for _, v := range blankNeighbors {

				sc.fr.setStone(stone{v.x, v.y, 3})
				sc.stack.push(v)
			}
		}
		if sc.debug {
			sc.fr.print() // for debug
		}

		nextToBlack := sc.checkNeighbors(x, y, -1)
		nextToWhite := sc.checkNeighbors(x, y, 1)

		// check if touching two colors invalidates group
		if touchedColor == 0 {
			if nextToBlack && nextToWhite {
				if sc.debug {
					fmt.Println("group touchedColor = 6346 (\"dead\") b/c both") // debug
				}
				touchedColor = 6346
			} else if nextToWhite {
				touchedColor = 1
			} else if nextToBlack {
				touchedColor = -1
			}
			//if touches neither, color is 0
		} else if touchedColor == -1 && nextToWhite {
			if sc.debug {
				fmt.Println("group touchedColor = 6346 (\"dead\") b/c white") // for debug
			}
			touchedColor = 6346
		} else if touchedColor == 1 && nextToBlack {
			if sc.debug {
				fmt.Println("group touchedColor = 6346 (\"dead\") b/c black") // for debug
			}
			touchedColor = 6346
		} else if touchedColor == 6346 {
			currentScore = 0;
		}
	}
	if sc.stack.size > 0 || (next != stone{color: 99}) {
		var s stone;
		if (next == stone{color: 99}) {
			s = sc.stack.pop()
		} else {
			s = next
		}
		score, color := sc.countOneStone(s, touchedColor, currentScore)

		return score, color

	}
	if sc.debug {
		fmt.Println("ended ran out of stack") // for debug
	}
	return currentScore, touchedColor
}


func (sc *score) countAllStones() {
	size := sc.fr.size
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			s := sc.fr.getStone(x, y)
			if s.color == 0 {
				//sc.stack.push(s);
				sc.fr.setStone(stone{s.x, s.y, 3})
				score, color := sc.countOneStone(s, 0, 0)

				switch color {
				case -1:
					sc.blackScore += score
				case 1:

					sc.whiteScore += score
				case 0, 6346:
					continue

				default:
					//print(color)
					log.Panic("Incorrect Color Return!")
				}
			}

		}
	}

	/*for(sc.stack.size > 0){
		s := sc.stack.pop()
		score, color := sc.countOneStone(s, 0, 0)


	}*/
}

func (g Game) end() {
	sc := score{fr: g.currentFrame(), pm: posMap{}.init()}

	//sc.debug = true // for debugging the score

	sc.countAllStones()
	fmt.Println("White: "+strconv.Itoa(sc.whiteScore), "Black: "+strconv.Itoa(sc.whiteScore))

}

func (g Game) Pass(color int) bool {
	if color == g.currentColor {
		if g.lastPass == g.currentColor * -1 {
			g.gameOver = true
			return true
		}
		g.currentTurn ++
		g.lastPass = color
		g.currentColor *= -1
		return true
	} else {
		return false
	}
}

func (g Game) Forfeit(color int) {

}

//cmdline specific
