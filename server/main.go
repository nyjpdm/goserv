package main

import (
	"fmt"
)

// import (
// 	"bufio"
// 	"os"
// 	"strconv"
// 	"strings"
// )

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

	err = tree.MakeMove(0*7 + 0)  //B
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
	
	// fmt.Println("=== ИГРА ГО ===")
	// fmt.Println("Формат ввода: i j (например: 1 2)")
	// fmt.Println("Выход: quit")
	// fmt.Println()
	
	// // Показываем пустую доску в начале
	// printCurrentBoard(tree)
	
	// scanner := bufio.NewScanner(os.Stdin)
	// moveCount := 0
	
	// for {
	// 	moveCount++
	// 	fmt.Printf("Ход %d (%c) > ", moveCount, tree.Current.CurrentColor)
		
	// 	if !scanner.Scan() {
	// 		break
	// 	}
		
	// 	input := strings.TrimSpace(scanner.Text())
	// 	if input == "quit" || input == "exit" {
	// 		break
	// 	}
		
	// 	// Парсим ввод
	// 	parts := strings.Split(input, " ")
	// 	if len(parts) != 2 {
	// 		fmt.Println("❌ Ошибка: введите два числа через пробел (i j)")
	// 		moveCount-- // Не считаем этот ввод как ход
	// 		continue
	// 	}
		
	// 	i, err1 := strconv.Atoi(parts[0])
	// 	j, err2 := strconv.Atoi(parts[1])
		
	// 	if err1 != nil || err2 != nil {
	// 		fmt.Println("❌ Ошибка: i и j должны быть числами")
	// 		moveCount--
	// 		continue
	// 	}
		
	// 	// Проверяем границы
	// 	if i < 0 || i >= 7 || j < 0 || j >= 7 {
	// 		fmt.Println("❌ Ошибка: i и j должны быть от 0 до 6")
	// 		moveCount--
	// 		continue
	// 	}
		
	// 	// Вычисляем позицию и делаем ход
	// 	pos := i*7 + j
	// 	err := tree.MakeMove(pos)
	// 	if err != nil {
	// 		fmt.Printf("❌ Ошибка хода: %v\n\n", err)
	// 		moveCount-- // Не считаем неудачный ход
	// 	} else {
	// 		fmt.Printf("✅ Ход принят: %c -> (%d,%d)\n\n", 
	// 			tree.Current.Parent.CurrentColor, i, j)
	// 		printCurrentBoard(tree)
	// 	}
	// }
	
	// fmt.Println("Игра завершена!")
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
// package main

// import (
// 	"fmt"
// )

// // [НЕ ОТНОСИТСЯ К ОСНОВНОЙ РАЗРАБОТКЕ]

// // PrintBoardSimple выводит упрощенную версию доски
// func (node *GoNode) PrintBoardSimple(boardSize int) {
// 	fmt.Println("Доска Го:")
// 	for row := 0; row < boardSize; row++ {
// 		for col := 0; col < boardSize; col++ {
// 			index := row*boardSize + col
// 			Point := node.Position[index]

// 			switch Point {
// 			case black:
// 				fmt.Print("B ")
// 			case white:
// 				fmt.Print("W ")
// 			case empty:
// 				fmt.Print(". ")
// 			}
// 		}
// 		fmt.Println()
// 	}
// }

// func main() {
// 	fmt.Println("=== ТЕСТ ЗАХВАТА КАМНЕЙ В ГО ===")

// 	// Тест 1: Простой захват одного камня
// 	//fmt.Println("\n1. ТЕСТ: Захват одного камня в центре")
// 	//testSingleCapture()

// 	// Тест 2: Захват группы из нескольких камней
// 	//fmt.Println("\n2. ТЕСТ: Захват группы из двух камней")
// 	//testGroupCapture()

// 	// Тест 3: Самоубийственный ход
// 	//fmt.Println("\n3. ТЕСТ: Проверка самоубийственного хода")
// 	//testSuicideMove()

// 	// Тест 4: Правило ко
// 	fmt.Println("\n4. ТЕСТ: Проверка правила ко")
// 	testTrickyKo()
// 	fmt.Println("\n 5. защелка")
// 	testClick()
// 	fmt.Println("6. Вырожденная защелка")
// 	testSmallClk()
// }

// func testSingleCapture() {
// 	settings := GameSettings{BoardSize: 5}
// 	tree := NewGoTree(settings)

// 	fmt.Println("Доска 5x5: создаем ситуацию для захвата")

// 	// Белые ставят камень
// 	tree.MakeMove(12) // Центр доски 5x5 (позиция 12)
// 	fmt.Println("Белые -> позиция 12 (центр)")
// 	printCurrentBoard(tree)

// 	// Черные окружают белый камень
// 	tree.MakeMove(7)  // Сверху
// 	tree.MakeMove(9)  //white
// 	tree.MakeMove(11) // Слева
// 	tree.MakeMove(10)
// 	tree.MakeMove(13) // Справа
// 	tree.MakeMove(14)
// 	tree.MakeMove(17) // Снизу

// 	fmt.Println("Черные окружили белый камень:")
// 	printCurrentBoard(tree)

// 	// Проверяем, что белый камень еще на доске
// 	if tree.Current.Position[12] == white {
// 		fmt.Println("✅ Белый камень еще на доске (ожидаемо)")
// 	}

// 	// Черные завершают захват
// 	err := tree.MakeMove(12)
// 	if err != nil {
// 		fmt.Printf("❌ Ошибка захвата: %v\n", err)
// 		return
// 	}

// 	fmt.Println("Черные завершили захват:")
// 	printCurrentBoard(tree)

// 	// Проверяем результат
// 	if tree.Current.Position[12] == empty {
// 		fmt.Println("✅ УСПЕХ: Белый камень захвачен и удален!")
// 	} else {
// 		fmt.Printf("❌ ПРОВАЛ: Позиция 12 = %c (должна быть пустой)\n", tree.Current.Position[12])
// 	}

// 	if tree.Current.BlackCaptures == 1 {
// 		fmt.Println("✅ УСПЕХ: Счетчик захватов корректен!")
// 	} else {
// 		fmt.Printf("❌ ПРОВАЛ: Захвачено %d камней (должно быть 1)\n", tree.Current.BlackCaptures)
// 	}
// }

// func testClick() {
// 	settings := GameSettings{BoardSize: 7}
// 	tree := NewGoTree(settings)

// 	fmt.Println("Защелка")

// 	tree.MakeMove(4)
// 	tree.MakeMove(5)
// 	tree.MakeMove(11)
// 	tree.MakeMove(12)

// 	tree.MakeMove(19)
// 	tree.MakeMove(0)
// 	tree.MakeMove(20)
// 	tree.MakeMove(1)

// 	tree.MakeMove(6)
// 	printCurrentBoard(tree)
// 	tree.MakeMove(13)
// 	tree.MakeMove(6)
// 	printCurrentBoard((tree))
// 	if tree.Current.BlackCaptures == 3 {
// 		fmt.Println("✅ успех, защелкой захвачено 3 камня")
// 	}
// }

// func testSmallClk() {
// 	settings := GameSettings{BoardSize: 7}
// 	tree := NewGoTree(settings)

// 	fmt.Println("НЕ защелка")

// 	tree.MakeMove(5)
// 	tree.MakeMove(13)
// 	tree.MakeMove(12)
// 	tree.MakeMove(0)

// 	tree.MakeMove(20)
// 	tree.MakeMove(1)
// 	printCurrentBoard(tree)
// 	tree.MakeMove(6)
// 	if tree.Current.BlackCaptures == 1 {
// 		fmt.Println("✅ Успех, захвачен 1 камень")
// 	}
// 	if tree.MakeMove(13) == nil {
// 		fmt.Print("❌ суицидальный ход белых должен быть запрещен")
// 	}

// }

// func testTrickyKo() {
// 	settings := GameSettings{BoardSize: 5}
// 	tree := NewGoTree(settings)

// 	fmt.Println("Проверка правила ко")

// 	// Создаем ситуацию ко
// 	tree.MakeMove(7) // B
// 	tree.MakeMove(4) // W
// 	tree.MakeMove(1) // B
// 	tree.MakeMove(2) // W
// 	tree.MakeMove(0) // B
// 	tree.MakeMove(8) // W
// 	tree.MakeMove(24)
// 	tree.MakeMove(23)
// 	fmt.Println("Создана ситуация ко:")
// 	//	printCurrentBoard(tree)

// 	// Черные захватывают белый камень
// 	fmt.Print(tree.MakeMove(3)) // B - захватывает белый в 7
// 	fmt.Print(tree.MakeMove(2))
// 	fmt.Println("После захвата белого камня:")
// 	//	printCurrentBoard(tree)

// }

// func testGroupCapture() {
// 	settings := GameSettings{BoardSize: 5}
// 	tree := NewGoTree(settings)

// 	fmt.Println("Захват группы из двух соединенных камней")

// 	// Белые ставят два соединенных камня
// 	tree.MakeMove(6) // (1,1)
// 	tree.MakeMove(7) // (1,2) - соединен с первым
// 	printCurrentBoard(tree)

// 	// Черные окружают группу
// 	tree.MakeMove(1)  // Сверху-слева
// 	tree.MakeMove(2)  // Сверху
// 	tree.MakeMove(3)  // Сверху-справа
// 	tree.MakeMove(5)  // Слева
// 	tree.MakeMove(8)  // Справа
// 	tree.MakeMove(11) // Снизу-слева
// 	tree.MakeMove(12) // Снизу
// 	tree.MakeMove(13) // Снизу-справа

// 	fmt.Println("Черные окружили группу белых:")
// 	printCurrentBoard(tree)

// 	// Захватываем группу
// 	err := tree.MakeMove(10)
// 	if err != nil {
// 		fmt.Printf("❌ Ошибка захвата группы: %v\n", err)
// 		return
// 	}

// 	fmt.Println("После захвата группы:")
// 	printCurrentBoard(tree)

// 	// Оба белых камня должны быть удалены
// 	if tree.Current.Position[6] == empty && tree.Current.Position[7] == empty {
// 		fmt.Println("✅ УСПЕХ: Оба белых камня захвачены!")
// 	} else {
// 		fmt.Printf("❌ ПРОВАЛ: Камни не удалены: pos6=%c, pos7=%c\n",
// 			tree.Current.Position[6], tree.Current.Position[7])
// 	}

// 	if tree.Current.BlackCaptures == 2 {
// 		fmt.Println("✅ УСПЕХ: Захвачено 2 камня!")
// 	} else {
// 		fmt.Printf("❌ ПРОВАЛ: Захвачено %d камней (должно быть 2)\n", tree.Current.BlackCaptures)
// 	}
// }

// func testSuicideMove() {
// 	settings := GameSettings{BoardSize: 5}
// 	tree := NewGoTree(settings)

// 	fmt.Println("Проверка запрета самоубийственного хода")

// 	// Черные создают полностью окруженную точку
// 	tree.MakeMove(6)  //
// 	tree.MakeMove(7)  //
// 	tree.MakeMove(8)  //
// 	tree.MakeMove(11) //
// 	tree.MakeMove(14) //
// 	tree.MakeMove(13) //
// 	tree.MakeMove(19) //
// 	tree.MakeMove(17) //
// 	// Позиция 12 полностью окружена

// 	fmt.Println("Создана полностью окруженная позиция 12:")
// 	printCurrentBoard(tree)

// 	// Белые пытаются пойти в окруженную точку - должно быть самоубийство
// 	err := tree.MakeMove(12)
// 	if err != nil {
// 		fmt.Printf("✅ УСПЕХ: Самоубийственный ход запрещен: %v\n", err)
// 	} else {
// 		fmt.Println("❌ ПРОВАЛ: Самоубийственный ход разрешен!")
// 		printCurrentBoard(tree)
// 	}
// }

// func testKoRule() {
// 	settings := GameSettings{BoardSize: 5}
// 	tree := NewGoTree(settings)

// 	fmt.Println("Проверка правила ко")

// 	// Создаем ситуацию ко
// 	tree.MakeMove(7) // B
// 	tree.MakeMove(4) // W
// 	tree.MakeMove(1) // B
// 	tree.MakeMove(2) // W
// 	tree.MakeMove(0) // B
// 	tree.MakeMove(8) // W

// 	fmt.Println("Создана ситуация ко:")
// 	printCurrentBoard(tree)

// 	// Черные захватывают белый камень
// 	tree.MakeMove(3) // B - захватывает белый в 7
// 	tree.MakeMove(2)
// 	fmt.Println("После захвата белого камня:")
// 	printCurrentBoard(tree)

// 	// Белые немедленно пытаются вернуться - должно быть ко
// 	err := tree.MakeMove(3)
// 	if err != nil {
// 		fmt.Printf("✅ УСПЕХ: Правило ко работает: %v\n", err)
// 	} else {
// 		fmt.Println("❌ ПРОВАЛ: Правило ко не сработало!")
// 		printCurrentBoard(tree)
// 	}
// }

// func printCurrentBoard(tree *GoTree) {
// 	fmt.Println("Текущая доска:")
// 	boardSize := tree.BoardSize
// 	for i := 0; i < boardSize; i++ {
// 		for j := 0; j < boardSize; j++ {
// 			idx := i*boardSize + j
// 			fmt.Printf("%c ", tree.Current.Position[idx])
// 		}
// 		fmt.Println()
// 	}
// 	fmt.Printf("Ход: %d, Цвет: %c, Захваты: B=%d W=%d\n\n",
// 		tree.Current.MoveNumber, tree.Current.CurrentColor,
// 		tree.Current.BlackCaptures, tree.Current.WhiteCaptures)
// }
