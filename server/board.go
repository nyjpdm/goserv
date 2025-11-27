package main

import (
	"fmt"
	"time"
)

type Stone byte

const (
	empty Stone = '.'
	black Stone = 'B'
	white Stone = 'W'
)

type GoNode struct {
	Parent   *GoNode
	Children []*GoNode

	Position []Stone
	Move     int

	MoveNum int
	IsBlack bool

	BlackCaptures int
	WhiteCaptures int

	// Комментарии, чат партии и время добавим очень не скоро! [не смотреть сюда]
}

type GoTree struct {
	Root    *GoNode
	Current *GoNode

	GameSettings // Все поля GameSettings теперь встроены в GoTree

	// Системные метаданные (изменим в конце партии)
	Result    string
	Winner    string
	CreatedAt time.Time // Дата создания GoTree
}

type GameSettings struct {
	BoardSize int
	Komi      float64
	Rules     string
	Handicap  int
	BlackName string
	WhiteName string
	Event     string
}

func NewGoTree(settings GameSettings) *GoTree {
	// Дефолтные значения
	if settings.BoardSize == 0 {
		settings.BoardSize = 19
	}
	if settings.Komi == 0 {
		settings.Komi = 6.5
	}
	if settings.Rules == "" {
		settings.Rules = "japanese"
	}
	if settings.BlackName == "" {
		settings.BlackName = "Black"
	}
	if settings.WhiteName == "" {
		settings.WhiteName = "White"
	}

	// Создаем доску и root
	position := make([]Stone, settings.BoardSize*settings.BoardSize)
	for i := range position {
		position[i] = empty
	}

	root := &GoNode{
		Position: position,
		Move:     -1,
		MoveNum:  0,
		IsBlack:  true,
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
		if child.Move == move {
			// Ход уже рассчитан! Значит просто переходим к нему
			tree.Current = child
			return nil
		}
	}

	// #2 Базовая валидация (проверка на диапозон и свободный пункт)
	if move < 0 || move >= tree.BoardSize*tree.BoardSize {
		return fmt.Errorf("move %d out of board range", move)
	}
	if tree.Current.Position[move] != empty {
		return fmt.Errorf("position %d alreday occupied", move)
	}

	// ...

	return fmt.Errorf("advanced Go logic not implemented yet") // Пока оставим такой return
}

type Group struct {
	Graph          []bool
	Dame           int
	Color          Stone
	NumberOfStones int
}

func NewGroupFromPosition(position []Stone, startPos int, boardSize int) *Group {
	// Проверяем, что начальная позиция не пустая
	if position[startPos] == empty {
		return nil
	}

	group := &Group{
		Graph:          make([]bool, boardSize*boardSize),
		Dame:           0,
		Color:          position[startPos],
		NumberOfStones: 0,
	}

	// Запускаем поиск в ширину/глубину для нахождения всех связанных камней
	group.findConnectedStones(position, startPos, boardSize)

	// Вычисляем количество дыханий
	group.calculateDame(position, boardSize)

	return group
}

// findConnectedStones находит все камни, связанные с начальной позицией
func (g *Group) findConnectedStones(position []Stone, startPos int, boardSize int) {
	visited := make([]bool, boardSize*boardSize)
	queue := []int{startPos}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		// Добавляем камень в группу
		g.Graph[current] = true
		g.NumberOfStones++

		// Проверяем всех соседей того же цвета
		neighbors := getNeighbors(current, boardSize)
		for _, neighbor := range neighbors {
			if !visited[neighbor] && position[neighbor] == g.Color {
				queue = append(queue, neighbor)
			}
		}
	}
}

// calculateDame вычисляет количество дыханий группы
func (g *Group) calculateDame(position []Stone, boardSize int) {
	dameCount := 0
	visitedDame := make([]bool, boardSize*boardSize)

	// Проходим по всем камням группы
	for pos := 0; pos < len(g.Graph); pos++ {
		if g.Graph[pos] {
			// Для каждого камня проверяем соседей
			neighbors := getNeighbors(pos, boardSize)
			for _, neighbor := range neighbors {
				// Если сосед пустой и мы еще не считали это дыхание
				if position[neighbor] == empty && !visitedDame[neighbor] {
					dameCount++
					visitedDame[neighbor] = true
				}
			}
		}
	}

	g.Dame = dameCount
}

// IsAlive проверяет, жива ли группа (есть ли дыхания)
func (g *Group) IsAlive() bool {
	return g.Dame > 0
}

// IsDead проверяет, мертва ли группа
func (g *Group) IsDead() bool {
	return g.Dame == 0
}

// GetGraph возвращает карту группы (копию для безопасности)
func (g *Group) GetGraph() []bool {
	graph := make([]bool, len(g.Graph))
	copy(graph, g.Graph)
	return graph
}

// GetStonePositions возвращает список позиций камней группы
func (g *Group) GetStonePositions() []int {
	var positions []int
	for pos, exists := range g.Graph {
		if exists {
			positions = append(positions, pos)
		}
	}
	return positions
}

// getNeighbors возвращает индексы соседних клеток для данной позиции
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

// [ДАЛЕЕ: НЕ ОТНОСИТСЯ К ОСНОВНОЙ РАЗРАБОТКЕ]

// PrintBoardSimple выводит упрощенную версию доски
func (node *GoNode) PrintBoardSimple(boardSize int) {

	fmt.Println("Доска Го:")
	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			index := row*boardSize + col
			stone := node.Position[index]

			switch stone {
			case black:
				fmt.Print("B ")
			case white:
				fmt.Print("W ")
			case empty:
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}
