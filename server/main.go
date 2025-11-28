package main

import (
	"fmt"
)

func main() {
	settings := GameSettings{BoardSize: 7}
	tree := NewGoTree(settings)

	err := tree.MakeMove(1*7 + 1) //B
	if err != nil {
		fmt.Println(err)
	} else {
		printCurrentBoard(tree)
	}
	err = tree.MakeMove(1*7 + 0) //W
	if err != nil {
		fmt.Println(err)
	} else {
		printCurrentBoard(tree)
	}
	err = tree.MakeMove(0*7 + 2) //B
	if err != nil {
		fmt.Println(err)
	} else {
		printCurrentBoard(tree)
	}
	err = tree.MakeMove(0*7 + 1) //W
	if err != nil {
		fmt.Println(err)
	} else {
		printCurrentBoard(tree)
	}

	err = tree.MakeMove(0*7 + 0) //B
	if err != nil {
		fmt.Println(err)
	} else {
		printCurrentBoard(tree)
	}

	err = tree.MakeMove(0*7 + 1) //W
	if err != nil {
		fmt.Println(err)
	} else {
		printCurrentBoard(tree)
	}
	runServer()
}
func stringBoard(tree *GoTree) string {
	res := ""
	boardSize := tree.BoardSize
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			idx := i*boardSize + j
			res += string(tree.CurrentNode.Board[idx])
		}
		res += "\n"
	}
	return res
}

func printCurrentBoard(tree *GoTree) {
	fmt.Println("  0 1 2 3 4 5 6") // заголовок столбцов
	boardSize := tree.BoardSize
	for i := 0; i < boardSize; i++ {
		fmt.Printf("%d ", i) // номер строки
		for j := 0; j < boardSize; j++ {
			idx := i*boardSize + j
			fmt.Printf("%c ", tree.CurrentNode.Board[idx])
		}
		fmt.Println()
	}
	fmt.Printf("\nХод: %d, Следующий: %c, Захваты: B=%d W=%d\n\n",
		tree.CurrentNode.NodeOrder, tree.CurrentNode.LastMoveColor,
		tree.CurrentNode.BlackCaptures, tree.CurrentNode.WhiteCaptures)
}
