package main 

type Chain struct {
	ChainMap map[int]bool
	Dame int
	Color PointColor
	StoneCount int
}

func FindChainAt(board []PointColor, startPos int, boardSize int) *Chain {
	// Проверяем, что начальная пункт startPos не пуст
	if board[startPos] == empty {
		return nil
	}

	chain := &Chain{
		ChainMap:	make(map[int]bool, boardSize * boardSize),
		Dame:	0,
		Color:	board[startPos],
		StoneCount:	0,
	}

	// Запускаем поиск BFS для обхода цепи
	chain.findConnectedStones(board, startPos, boardSize)

	// Вычисляем количество дыханий
	chain.calculateDame(board, boardSize)

	return chain
}


// findConnectedStones находит все камни, связанные с startPos
func (ch *Chain) findConnectedStones(board []PointColor, startPos int, boardSize int) {
    visited := make(map[int]bool, boardSize*boardSize)
    queue := []int{startPos}
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        if visited[current] {
            continue
        }
        visited[current] = true
        
        // Добавляем камень в группу
        ch.ChainMap[current] = true
        ch.StoneCount++
        
        // Проверяем всех соседей того же цвета
        neighbors := getNeighbors(current, boardSize)
        for _, neighbor := range neighbors {
            if !visited[neighbor] && board[neighbor] == ch.Color {
                queue = append(queue, neighbor)
            }
        }
    }
}

// calculateDame вычисляет количество дыханий группы
func (ch *Chain) calculateDame(board []PointColor, boardSize int) {
    dameCount := 0
    visitedDame := make(map[int]bool, boardSize*boardSize)
    
    for pos := range ch.ChainMap{
		// Для каждого камня проверяем соседей
        neighbors := getNeighbors(pos, boardSize)
        for _, neighbor := range neighbors {
            if board[neighbor] == empty && !visitedDame[neighbor] {
                dameCount++
                visitedDame[neighbor] = true
            }
        }
    }
    
    ch.Dame = dameCount
}

func (ch *Chain) IsAlive() bool {
    return ch.Dame > 0
}

func (ch *Chain) IsDead() bool {
    return ch.Dame == 0
}