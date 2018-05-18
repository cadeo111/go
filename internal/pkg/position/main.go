package position

import (
	"strconv"
	"bytes"
	"fmt"
	sn "goGame/internal/pkg/stone"
)
//Map is a representation of the board at a certain point in the game.
type Map struct {
	m map[string]bool
}

func (p Map) Add(x int, y int) {
	str := strconv.Itoa(x) + "," + strconv.Itoa(y)
	p.m[str] = true
}

func (p Map) Get(x int, y int) bool {
	str := strconv.Itoa(x) + "," + strconv.Itoa(y)
	return p.m[str]
}

func (p Map) StoneGet(s sn.Stone) bool {
	return p.Get(s.X, s.Y)
}

func (Map) Init() Map {
	return Map{m: make(map[string]bool)}
}



const DEBUG_FRAME = false

type Frame struct {
	//Size is the max x and y of a position on the board
	Size int
	//board is a way the stone position is stored
	board []int
	//turn is the turn associated with the position
	turn int
}


//Public

func (f Frame) Print() {
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
	for i := 0; i < f.Size*2; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
	for i := 0; i < f.Size; i++ {
		//fmt.Println(i*f.Size," ",(i+1)*f.Size)
		fmt.Println(printAline(f.board[i*f.Size : (i+1)*f.Size]))
		//fmt.Println(f.board[i*f.Size:(i+1)*f.Size])
	}
	for i := 0; i < f.Size*2; i++ {
		fmt.Print("-")
	}

	fmt.Print("\n")
}
func (f Frame) toFormattedString() string {
	var final bytes.Buffer
	aLine := func(s []int) string {
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


	for i := 0; i < f.Size*2; i++ {
		final.WriteString("-")
	}
	final.WriteString("\n")
	for i := 0; i < f.Size; i++ {

		final.WriteString(aLine(f.board[i*f.Size : (i+1)*f.Size]))

	}
	for i := 0; i < f.Size*2; i++ {
		final.WriteString("-")
	}
	final.WriteString("\n")
	return final.String()
}


func (f Frame) ToString() string {
	var final bytes.Buffer
	aLine := func(s []int) string {
		var buffer bytes.Buffer
		for i := 0; i < len(s); i++ {

			switch s[i] {
			case -1:
				buffer.WriteString("b")
			case 0:
				buffer.WriteString("n")
			case 1:
				buffer.WriteString("w")
			default:
				buffer.WriteString(strconv.Itoa(s[i]))
			}
			if i != len(s)-1 {
				buffer.WriteString("")
			}

		}
		return buffer.String()
	}
	for i := 0; i < f.Size; i++ {

		final.WriteString(aLine(f.board[i*f.Size : (i+1)*f.Size]))
		if i != f.Size-1 {
			final.WriteString("/")
		}
	}
	return final.String()
}



//GetStone is the function to retrieve the data of a position in a position as a stone.
func (f Frame) GetStone(x int, y int) (sn.Stone) {
	pos := ((f.Size * f.Size) - (y * f.Size) + x) - f.Size
	stone := sn.Stone{x, y, f.board[pos]}
	return stone
}

func (f *Frame) SetStone(s sn.Stone) {
	pos := ((f.Size * f.Size) - (s.Y * f.Size) + s.X) - f.Size
	f.board[pos] = s.Color
}

func (f Frame) GetSurroundingStones(s sn.Stone) (sn.Stone, sn.Stone, sn.Stone, sn.Stone) {
	top, left, bottom, right := 9, 9, 9, 9
	if s.X+1 < f.Size {
		right = f.GetStone(s.X+1, s.Y).Color
	}
	if s.X > 0 {
		left = f.GetStone(s.X-1, s.Y).Color
	}
	if s.Y+1 < f.Size {
		top = f.GetStone(s.X, s.Y+1).Color
	}
	if s.Y > 0 {
		bottom = f.GetStone(s.X, s.Y-1).Color
	}

	return sn.Stone{s.X, s.Y + 1, top},
		sn.Stone{s.X + 1, s.Y, right},
		sn.Stone{s.X, s.Y - 1, bottom},
		sn.Stone{s.X - 1, s.Y, left}
	// returning 9 signifies out of bounds
}

func (f Frame) ProcessTakenStones(s sn.Stone) (Frame Frame, numRemoved int) {
	nf := f.duplicate()
	nf.SetStone(s)
	t, r, b, l := f.GetSurroundingStones(s)

	oppositeColor := s.Color * -1

	takenStones := 0

	//first process opposite color
	if t.Color == oppositeColor {
		if !nf.SpotWouldHaveLiberty(t) {
			takenStones += nf.removeConnected(t)
		}
		if DEBUG_FRAME {
		fmt.Println("t", nf.SpotWouldHaveLiberty(t))
		}
	}
	if r.Color == oppositeColor {
		if !nf.SpotWouldHaveLiberty(r) {
			takenStones += nf.removeConnected(r)
		}
		if DEBUG_FRAME {
			fmt.Println("r", nf.SpotWouldHaveLiberty(r))
		}
	}
	if b.Color == oppositeColor {
		if !nf.SpotWouldHaveLiberty(b) {
			takenStones += nf.removeConnected(b)
		}
		if DEBUG_FRAME {
			fmt.Println("b", nf.SpotWouldHaveLiberty(b))
		}
	}
	if l.Color == oppositeColor {
		if !nf.SpotWouldHaveLiberty(l) {
			takenStones += nf.removeConnected(l)
		}
		if DEBUG_FRAME {
			fmt.Println("l", nf.SpotWouldHaveLiberty(l))
		}
	}
	return nf, takenStones

}

func (f Frame) SpotWouldHaveLiberty(s sn.Stone) bool {
	ct, cr, cb, cl := f.GetSurroundingStones(s)
	x, y, color := s.X, s.Y, s.Color
	pos := Map{}.Init()
	pos.Add(x, y)

	top, right, bottom, left := false, false, false, false

	if ct.Color == 0 || cr.Color == 0 || cb.Color == 0 || cl.Color == 0 {
		return true
	}

	if ct.Color == color {
		top = checkEachStone(x, y+1, color, f, pos)
	}
	if cr.Color == color {
		right = checkEachStone(x+1, y, color, f, pos)
	}
	if cb.Color == color {
		left = checkEachStone(x, y-1, color, f, pos)
	}
	if cl.Color == color {
		bottom = checkEachStone(x-1, y, color, f, pos)
	}
	return top || right || bottom || left
}

func NewFrame(size int, turn int) Frame {
	b := make([]int, size*size)
	return Frame{size, b, turn}
}

//private

func (f Frame) removeConnected(s sn.Stone)  (removedNum int) {
	f.SetStone(sn.Stone{})
	t, r, b, l := f.GetSurroundingStones(s)
	tot := 1

	if t.Color == s.Color {
		tot += f.removeConnected(t)
	}
	if r.Color == s.Color {
		tot += f.removeConnected(r)
	}
	if b.Color == s.Color {
		tot += f.removeConnected(b)
	}
	if l.Color == s.Color {
		tot+= f.removeConnected(l)
	}
	return tot
}

func (f Frame) duplicate() Frame {
	nf := Frame{f.Size, make([]int, f.Size*f.Size), f.turn}
	copy(nf.board, f.board)
	return nf
}

func checkEachStone(x int, y int, color int, cf Frame, pos Map) bool {
	fmt.Println(x, y, pos.Get(x, y))
	if pos.Get(x, y) {
		return false
	} else {
		pos.Add(x, y)
	}
	ct, cr, cb, cl := cf.GetSurroundingStones(sn.Stone{X: x, Y: y})
	fmt.Println(ct.Color, cr.Color, cb.Color, cl.Color)
	if (ct.Color == 0 && !pos.Get(ct.X, ct.Y)) ||
		(cr.Color == 0 && !pos.Get(cr.X, cr.Y)) ||
		(cb.Color == 0 && !pos.Get(cb.X, cb.Y)) ||
		(cl.Color == 0 && !pos.Get(cl.X, cl.Y)) {
		return true
	} else {
		top, right, bottom, left := false, false, false, false
		if ct.Color == color {
			top = checkEachStone(x, y+1, color, cf, pos)
		}
		if cr.Color == color {
			right = checkEachStone(x+1, y, color, cf, pos)
		}
		if cb.Color == color {
			left = checkEachStone(x, y-1, color, cf, pos)
		}
		if cl.Color == color {
			bottom = checkEachStone(x-1, y, color, cf, pos)
		}
		return top || right || bottom || left
	}
}
