package main 

type Group struct {
	ChainMap map[int]bool
	Dame int
	Color Point
	NumberOfStones int
}

func NewGroupFromPosition(position []Point, startPos int, boardSize int) *Group {
	// Проверяем, что начальная позиция не пустая
	if position[startPos] == empty {
		return nil
	}

	group := &Group{
		ChainMap:	make(map[int]bool, boardSize * boardSize),
		Dame:	0,
		Color:	position[startPos],
		NumberOfStones:	0,
	}

	// Запускаем поиск BFS для обхода цепи
	group.findConnectedStones(position, startPos, boardSize)

	// Вычисляем количество дыханий
	group.calculateDame(position, boardSize)

	return group
}


// findConnectedStones находит все камни, связанные с startPos
func (g *Group) findConnectedStones(position []Point, startPos int, boardSize int) {
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
        g.ChainMap[current] = true
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
func (g *Group) calculateDame(position []Point, boardSize int) {
    dameCount := 0
    visitedDame := make(map[int]bool, boardSize*boardSize)
    
    for pos := range g.ChainMap{
		// Для каждого камня проверяем соседей
        neighbors := getNeighbors(pos, boardSize)
        for _, neighbor := range neighbors {
            if position[neighbor] == empty && !visitedDame[neighbor] {
                dameCount++
                visitedDame[neighbor] = true
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
func (g *Group) GetGraph() map[int]bool {
    graph := make(map[int]bool, len(g.ChainMap))
    for pos, exists := range g.ChainMap {
        graph[pos] = exists
    }
    return graph
}

// GetPointPositions возвращает список позиций камней группы
func (g *Group) GetPointPositions() []int {
    positions := make([]int, 0, len(g.ChainMap))

    for pos := range g.ChainMap {
        positions = append(positions, pos)
    }
    return positions
}