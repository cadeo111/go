package main

import (
	"goGame/internal/app/gogame"
	"fmt"
	"strings"
	"bufio"
	"os"
	"strconv"
)

type cmdLED struct {
	gameSize int
	//serialConnection io.ReadWriteCloser
}

func (c cmdLED) RequestNextMove(color int, board string) (kind int, x int, y int) { // gives color, expects(1:move, 2:pass, 3:forfeit), x,y
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
		case "s","b","show board":
			c.ShowBoard(board, -1, -1)
		default:
			fmt.Println(
				"that is not one of the options\n" +
					"please try again")
		}
	}
	return 9, 9, 9
}

func (c cmdLED) ShowBoard(board string, capBlack int, capWhite int) {
	sa := strings.Split(strings.Replace(board, "/","",-1),"")
	fmt.Println(sa[0])
}

func (cmdLED) Message(error bool, kind string, message string) {
	fmt.Println(message)
}

func main() {
	input := cmdLED{9}

	g := gogame.Init(9, 0)
	g.Run(input)
}

