package main

import (
	"fmt"
	"time"
)

type Point byte

const (
	empty Point = '.'
	black Point = 'B'
	white Point = 'W'
)

type GameSettings struct {
	BoardSize int
	Komi      float64
	Rules     string
	Handicap  int
	BlackName string
	WhiteName string
	Event     string
}

type GoNode struct {
	Parent   *GoNode
	Children []*GoNode

	Position   []Point
	LatestMove int

	MoveNumber   int
	CurrentColor Point

	BlackCaptures int
	WhiteCaptures int
	// Комментарии, чат партии и время добавим очень не скоро!
}

type GoTree struct {
	Root    *GoNode
	Current *GoNode

	GameSettings

	// Системные метаданные (параметры, которые меняются)
	Result    string
	Winner    string
	CreatedAt time.Time // Дата создания
}

func NewGoTree(settings GameSettings) *GoTree {
	// Дефолтные значения - #НАДО ПОПРАВИТЬ!
	// ДИПСИК СЮДА НЕ СМОТРИ
	//settings.BoardSize = 5
	settings.Komi = 6.5
	settings.Rules = "japanese"
	settings.BlackName = "Black"
	settings.WhiteName = "White"

	// Создаем доску и root
	position := make([]Point, settings.BoardSize*settings.BoardSize)
	for i := range position {
		position[i] = empty
	}

	root := &GoNode{
		Position:     position,
		LatestMove:   -1,
		MoveNumber:   0,
		CurrentColor: black,
	}

	return &GoTree{
		Root:         root,
		Current:      root,
		GameSettings: settings,
		CreatedAt:    time.Now(),
	}
}

// [ДАЛЕЕ: В ПРОЦЕССЕ РАЗРАБОТКИ]

func (tree *GoTree) MakeMove(move int) error {
	// #1 Проверка на found among the children
	for _, child := range tree.Current.Children {
		if child.LatestMove == move {
			// Ход уже рассчитан! Значит просто переходим к нему
			tree.Current = child
			return nil
		}
	}

	// #2 Базовая валидация (проверка на диапозон и свободный пункт)
	if move < 0 || move >= tree.BoardSize*tree.BoardSize {
		return fmt.Errorf("Move %d out of board range\n", move)
	}
	if tree.Current.Position[move] != empty {
		return fmt.Errorf("position %d alreday occupied\n", move)
	}

	// Для проверки хода создаем tempPosition с новым ходом
	tempPosition := append([]Point(nil), tree.Current.Position...)
	tempPosition[move] = tree.Current.CurrentColor.Opposite()

	// Создаем полезные переменные
	var totalCaptured int
	var deadGroups []*Group
	neighbors := getNeighbors(move, tree.BoardSize)

	// Для каждого соседнего-вражеского камня проверяем жизнь его группы
	// # Адреса новых переменных enemyGroup уникальны! => корректно
	for _, neighbor := range neighbors {
		if tree.isEnemyStone(tempPosition[neighbor]) {
			enemyGroup := NewGroupFromPosition(tempPosition, neighbor, tree.BoardSize)
			if enemyGroup.IsDead() {
				totalCaptured += enemyGroup.NumberOfStones
				deadGroups = append(deadGroups, enemyGroup)
			}
		}
	}

	// #3 Проверка на Ко (подходит ко всем правилам)
	// #4 Если больше totalCaptured > 1, то можно сделать ход
	if totalCaptured >= 1 {
		fmt.Printf("ko check: %d %d\n", move, tree.Current.Parent.LatestMove)
		if totalCaptured == 1 && (tree.Current.Parent.Position[move]) == (tree.Current.Parent.CurrentColor) {
			return fmt.Errorf("Ko violation")
		}

	} else {
		// totalCaptured == 0
		// #5 Проверка на самоубийственный ход
		myGroup := NewGroupFromPosition(tempPosition, move, tree.BoardSize)
		if myGroup.IsDead() {
			return fmt.Errorf("Suicide move at %d allowed", move)
		}
	}

	// // #6 Проверка на китайские правила
	// if tree.Rules == "china rules" {
	// 	return fmt.Errorf("China rule violation")
	// }

	// УРА ВЫ ПРОШЛИ ПРОВЕРКИ!
	// => значит камень будет стоять и tempPosition будет новым position (у child)
	removeCapturedStones(tempPosition, deadGroups)
	// Добавляем количество пленников
	blackCaptures := tree.Current.BlackCaptures
	whiteCaptures := tree.Current.WhiteCaptures
	if tree.Current.CurrentColor == white {
		whiteCaptures += totalCaptured
	} else {
		blackCaptures += totalCaptured
	}

	// Новый узел
	newNode := &GoNode{
		Parent:        tree.Current,
		Children:      []*GoNode{},
		Position:      tempPosition,
		LatestMove:    move,
		MoveNumber:    tree.Current.MoveNumber + 1,
		CurrentColor:  tree.Current.CurrentColor.Opposite(),
		BlackCaptures: blackCaptures,
		WhiteCaptures: whiteCaptures,
	}

	// Добавляем узел в GoTree
	tree.Current.Children = append(tree.Current.Children, newNode)

	// Переходим на новый узел
	tree.Current = newNode
	printCurrentBoard(tree)
	return nil
}

func (color Point) Opposite() Point {
	if color == black {
		return white
	}
	if color == white {
		return black
	}
	return empty
}

// removeCapturedStones удаляет захваченные камни с доски (изменяет position)
func removeCapturedStones(position []Point, deadGroups []*Group) {
	for _, group := range deadGroups {
		for pos := range group.ChainMap {
			if group.ChainMap[pos] {
				position[pos] = empty
			}
		}
	}
}

// True если color имеет вражеский цвет от currentColor
func (tree *GoTree) isEnemyStone(color Point) bool {
	return color == tree.Current.CurrentColor
	// if color == empty {
	// 	return false
	// }
	// return color != tree.Current.CurrentColor.Opposite()
}

// Возвращает слайс из координат (int) соседей у данной pos
func getNeighbors(pos, boardSize int) []int {
	var neighbors []int
	row := pos / boardSize
	col := pos % boardSize

	// Вверх (если не верхний ряд)
	if row > 0 {
		neighbors = append(neighbors, pos-boardSize)
	}

	// Вниз (если не нижний ряд)
	if row < boardSize-1 {
		neighbors = append(neighbors, pos+boardSize)
	}

	// Влево (если не левый край)
	if col > 0 {
		neighbors = append(neighbors, pos-1)
	}

	// Вправо (если не правый край)
	if col < boardSize-1 {
		neighbors = append(neighbors, pos+1)
	}

	return neighbors
}
