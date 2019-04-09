package main

import (
	"go/app/gogame"
	"fmt"
	"strings"
	"bufio"
	"os"
	"strconv"
	"bytes"
)

type cmd struct {
	gameSize int
}

func (c cmd) RequestNextMove(color int, board string) (kind int, x int, y int) { // gives color, expects(1:move, 2:pass, 3:forfeit), x,y
	var colorString string
	switch color {
	case gogame.WHITE:
		colorString = "White"
	case gogame.BLACK:
		colorString = "Black"
	}

	reader := bufio.NewReader(os.Stdin)

	for true {
		fmt.Println("What is " + colorString + "'s move?")
		fmt.Println("(forfeit, move, pass, show board)")
		text, _ := reader.ReadString('\n')
		switch strings.Replace(strings.ToLower(text), "\n", "", 1) {

		case "m","move":
			fmt.Println("what x position?")
			text, _ = reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", 1)
			x, errorx := strconv.Atoi(text)
			if errorx != nil {
				fmt.Print("x wasn't an integer")
				break
			}
			fmt.Println("what y position?")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", 1)
			y, errory := strconv.Atoi(text)
			if errory != nil {
				fmt.Println("y wasn't an integer")
				break
			}
			return 1, x, y
		case "p","pass":
			return 2, 0, 0
		case "f","forfeit":
			return 3, 0, 0
		case "s","b","sb","show board":
			c.ShowBoard(board, -1, -1)
		default:
			fmt.Println(
				"that is not one of the options\n" +
					"please try again")
		}
	}
	return 9, 9, 9
}

func (c cmd) ShowBoard(board string, capBlack int, capWhite int) {
	for i := 0; i < (c.gameSize)*2+1; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")

	sa := strings.Split(board, "/")
	f := func (s []string) string {
		var buffer bytes.Buffer
		for i, v := range s {

			switch v {
			case "n":
				buffer.WriteString("•")
			default:
				buffer.WriteString(s[i])
			}
			if i != len(s)-1 {
				buffer.WriteString(" ")
			}
		}
		return buffer.String()
	}

	for i := 0; i < c.gameSize; i++ {
		fmt.Print((c.gameSize - i -1)," ")
		fmt.Println(f(strings.Split(sa[i],"")))
	}
	fmt.Print("+ ")
	for i := 0; i < (c.gameSize); i++ {
		fmt.Print(i, " ")
	}
	fmt.Print("\n")
	for i := 0; i < (c.gameSize)*2 +1; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
	if(capBlack > 0 || capWhite > 0) {
		fmt.Println("White: ", capBlack, "\nBlack: ", capWhite)
	}
	fmt.Print("\n\n")
}

func (cmd) Message(error bool, kind string, message string) {
	fmt.Println(message)
}

func main() {
	input := cmd{9}

	g := gogame.Init(9, 0)
	g.Run(input)
}
