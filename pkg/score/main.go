package score

import (
	"go/pkg/position"
	"go/pkg/stone"
	"fmt"
	"log"
	"time"
)

type Score struct {
	Fr    position.Frame // current position when Score is calculated
	Pm    position.Map   // Map When Score is calculated
	stack stone.Stack
	Black int
	White int
	Debug bool

	/*
	1. find a position that is blank
	2. check if it is next to a black stone
	3. check if it is next to a white stone
	4. check if it is next to any blank stones
	5. Add all blank stones to stack
	6. Add stone to recorded positions
	7. go to next blank stone in stack
	8. if no more stones in stack return
	*/
}

func (sc *Score) checkNeighbors(x int, y int, color int) (bool) {
	ts, rs, bs, ls := sc.Fr.GetSurroundingStones(stone.Stone{x, y, color})

	if ts.Color == color {
		return true
	}
	if rs.Color == color {
		return true
	}
	if bs.Color == color {
		return true
	}
	if ls.Color == color {
		return true
	}
	return false
}
func (sc *Score) findAllNeighborsBlank(x int, y int) ([]stone.Stone) {

	ts, rs, bs, ls := sc.Fr.GetSurroundingStones(stone.Stone{x, y, 0})
	var ss = []stone.Stone{ts, rs, bs, ls}
	var ret []stone.Stone
	for _, val := range ss {
		if val.Color == 0 {
			ret = append(ret, val)
		}
	}
	return ret

}
func (sc *Score) countOneStone(cs stone.Stone, touchedColor int, currentScore int) (int, int) { // Score, color
	y := cs.Y
	x := cs.X
	var next stone.Stone = stone.Stone{0, 0, 99}
	if !sc.Pm.StoneGet(cs) {
		sc.Fr.SetStone(stone.Stone{x, y, 103})

		if sc.Debug {
			sc.Fr.Print()
			duration := time.Millisecond * 250
			time.Sleep(duration)
		}



		currentScore += 1

		sc.Pm.Add(x, y) // save that stone has been checked

		blankNeighbors := sc.findAllNeighborsBlank(x, y) // find all neighbors that are blank

		if sc.Debug {
			if len(blankNeighbors) > 0 {
				next = blankNeighbors[0]
				for _, v := range blankNeighbors {

					sc.Fr.SetStone(stone.Stone{v.X, v.Y, 3})

				}
			}
			sc.Fr.Print() // for Debug
		}

		sc.stack.PushSlice(blankNeighbors)



		nextToBlack := sc.checkNeighbors(x, y, -1)
		nextToWhite := sc.checkNeighbors(x, y, 1)

		// check if touching two colors invalidates group
		if touchedColor == 0 {
			if nextToBlack && nextToWhite {
				if sc.Debug {
					fmt.Println("group touchedColor = 6346 (\"dead\") b/c both") // Debug
				}
				touchedColor = 6346
			} else if nextToWhite {
				touchedColor = 1
			} else if nextToBlack {
				touchedColor = -1
			}
			//if touches neither, color is 0
		} else if touchedColor == -1 && nextToWhite {
			if sc.Debug {
				fmt.Println("group touchedColor = 6346 (\"dead\") b/c white") // for Debug
			}
			touchedColor = 6346
		} else if touchedColor == 1 && nextToBlack {
			if sc.Debug {
				fmt.Println("group touchedColor = 6346 (\"dead\") b/c black") // for Debug
			}
			touchedColor = 6346
		} else if touchedColor == 6346 {
			currentScore = 0;
		}
	}
	if sc.stack.Size > 0 || (next != stone.Stone{Color: 99}) {
		var s stone.Stone;
		if (next == stone.Stone{Color: 99}) {
			s = sc.stack.Pop()
		} else {
			s = next
		}
		score, color := sc.countOneStone(s, touchedColor, currentScore)

		return score, color

	}
	if sc.Debug {
		fmt.Println("ended ran out of stack") // for Debug
	}
	return currentScore, touchedColor
}
func (sc *Score) CountAllStones() {
	size := sc.Fr.Size
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			s := sc.Fr.GetStone(x, y)
			if s.Color == 0 {
				//sc.stack.push(s);
				sc.Fr.SetStone(stone.Stone{s.X, s.Y, 3})
				score, color := sc.countOneStone(s, 0, 0)

				switch color {
				case -1:
					sc.Black += score
				case 1:

					sc.White += score
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
		Score, color := sc.countOneStone(s, 0, 0)


	}*/
}