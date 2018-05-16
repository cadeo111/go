package stone

type Stone struct {
	X     int
	Y     int
	Color int // 0 none, -1 black, 1 white, 9 out of bounds
}


// stone is a representation of a spot on a board.
// it is derived from a position.Map
func (s Stone) Zero() Stone {
	s.Color = 0
	return s
}


type Stack struct {
	slice []Stone
	Size  int
}

func (st *Stack) Push(sn Stone) {
	st.slice = append(st.slice, sn)
	st.Size += 1
}
func (st *Stack) PushSlice(ss []Stone) {
	for _, s := range ss {
		st.Push(s)
	}
}

func (st *Stack) Peek() Stone {
	return st.slice[len(st.slice)-1] // view last item in stack which was added last
}

func (st *Stack) Pop() Stone {
	if st.Size > 0 {
		var r = st.Peek()
		st.slice = st.slice[0 : len(st.slice)-1] // remove last item in slice the one that was peeked
		st.Size -= 1
		return r
	} else {
		return Stone{}
	}
}