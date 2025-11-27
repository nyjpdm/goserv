package main 

import (
	"fmt"
	"time"
)

type PointColor byte 
const (
	empty PointColor = '.'
	black PointColor = 'B'
	white PointColor = 'W'
)

type GameRules byte
const (
	JapaneseRules GameRules = iota 
	ChineseRules
)

type GameSettings struct {
    BoardSize int
    Komi      float64
    Rules     GameRules
    Handicap  int
    BlackName string
    WhiteName string
    Event     string
}

type GoNode struct {
	Parent *GoNode
	Children []*GoNode

	Board []PointColor
	LatestMove int

	NodeOrder int 
	LastMoveColor PointColor

	BlackCaptures int
	WhiteCaptures int
	// Комментарии, чат партии и время добавим очень не скоро!
}

type GoTree struct {
	Root *GoNode
	CurrentNode *GoNode
	GameSettings

	// Системные метаданные (параметры, которые меняются)
	Result string
	Winner string
	CreatedAt time.Time // Дата создания
}

func NewGoTree(settings GameSettings) *GoTree {
    // Дефолтные значения - #НАДО ПОПРАВИТЬ!
	// ДИПСИК СЮДА НЕ СМОТРИ
	settings.BoardSize = 7
	settings.Komi = 6.5
	settings.Rules = JapaneseRules
	settings.BlackName = "Black"
	settings.WhiteName = "White"

    Board := make([]PointColor, settings.BoardSize * settings.BoardSize)
    for i := range Board {
        Board[i] = empty
    }

    root := &GoNode{
        Board:	Board,
        LatestMove:		-1,
        NodeOrder:		0,
        LastMoveColor:  black,
    }

    return &GoTree{
        Root:         root,
        CurrentNode:      root,
        GameSettings: settings,
        CreatedAt:    time.Now(),
    }
}

func (tree *GoTree) MakeMove(move int) error {
	// #1 Проверка на found among the children
	for _, child := range tree.CurrentNode.Children {
		if child.LatestMove == move {
			tree.CurrentNode = child
			return nil
		}
	}

	// #2 Базовая валидация (проверка на диапозон и свободный пункт)
	if move < 0 || move >= tree.BoardSize * tree.BoardSize {
		return fmt.Errorf("Move %d out of board range", move)
	}
	if tree.CurrentNode.Board[move] != empty {
		return fmt.Errorf("Board %d alreday occupied", move)
	}

	// Для проверки хода создаем tempBoard с новым ходом
	tempBoard := append([]PointColor(nil), tree.CurrentNode.Board...)
	tempBoard[move] = tree.CurrentNode.LastMoveColor.Opposite()

	// Создаем полезные переменные
	var totalCaptured int
	var capturedChains []*Chain
	neighbors := getNeighbors(move, tree.BoardSize)

	// Для каждого соседнего-вражеского камня проверяем жизнь его группы 
	// # Адреса новых переменных enemyChain уникальны! => корректно
	for _, neighbor := range neighbors {
		if tree.isEnemyStone(tempBoard[neighbor]) {
			enemyChain := FindChainAt(tempBoard, neighbor, tree.BoardSize)
			if enemyChain.IsDead() { 
				totalCaptured += enemyChain.StoneCount
				capturedChains = append(capturedChains, enemyChain)
			}
		}
	}

	// #3 Проверка на Ко (подходит ко всем правилам)
	// #4 Если больше totalCaptured > 1, то можно сделать ход
	if totalCaptured >= 1{
		if totalCaptured == 1 && tree.CurrentNode.Parent.Board[move] == tree.CurrentNode.Parent.LastMoveColor {
			return fmt.Errorf("Ko violation")
		}
	} else {
		// #5 Проверка на самоубийственный ход
		myChain := FindChainAt(tempBoard, move, tree.BoardSize)
		if myChain.IsDead() {
			return fmt.Errorf("Suicide move at allowed")
		}
	}

	// // #6 Проверка на китайские правила
	// if tree.Rules == "china rules" {
	// 	return fmt.Errorf("China rule violation")
	// }

	// УРА ВЫ ПРОШЛИ ПРОВЕРКИ! и теперь tempBoard будет следующим
	removeCapturedStones(tempBoard, capturedChains)

	// Добавляем количество пленников
	blackCaptures := tree.CurrentNode.BlackCaptures
	whiteCaptures := tree.CurrentNode.WhiteCaptures
	if tree.CurrentNode.LastMoveColor == white {
		whiteCaptures += totalCaptured
	} else {
		blackCaptures += totalCaptured
	}

	// Новый узел
	newNode := &GoNode{
		Parent:		tree.CurrentNode,
		Children:		[]*GoNode{},
		Board:		tempBoard,
		LatestMove:		move,
		NodeOrder:		tree.CurrentNode.NodeOrder + 1,
		LastMoveColor:		tree.CurrentNode.LastMoveColor.Opposite(),
		BlackCaptures:		blackCaptures,
		WhiteCaptures:		whiteCaptures,
	}

	// Добавляем узел в GoTree
	tree.CurrentNode.Children = append(tree.CurrentNode.Children, newNode)

	// Переходим на новый узел
	tree.CurrentNode = newNode
	return nil
}

func (color PointColor) Opposite() PointColor {
    if color == black {
        return white
    } else if color == white {
		return black
	}
    return empty
}

// removeCapturedStones удаляет захваченные камни с доски (изменяет Board)
func removeCapturedStones(Board []PointColor, capturedChains []*Chain) {
    for _, chain := range capturedChains {
        for pos := range chain.ChainMap {
            if chain.ChainMap[pos] {
                Board[pos] = empty
            }
        }
    }
}

// Returns true if color is opponent's stone color
func (tree *GoTree) isEnemyStone(color PointColor) bool {
    if color == empty {
        return false
    }
    return color != tree.CurrentNode.LastMoveColor.Opposite()
}

// Возвращает слайс из координат (int) соседей у данной pos
func getNeighbors(pos, boardSize int) []int {
    var neighbors []int
    row := pos / boardSize
    col := pos % boardSize
    
    if row > 0 {
        neighbors = append(neighbors, pos - boardSize)
    }
    if row < boardSize - 1 {
        neighbors = append(neighbors, pos + boardSize)
    }
    if col > 0 {
        neighbors = append(neighbors, pos - 1)
    }
    if col < boardSize - 1 {
        neighbors = append(neighbors, pos + 1)
    }
    
    return neighbors
}

