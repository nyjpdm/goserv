package main

// GetChain возвращает карту группы (копию для безопасности)
func (ch *Chain) GetChain() map[int]bool {
    graph := make(map[int]bool, len(ch.ChainMap))
    for pos, exists := range ch.ChainMap {
        graph[pos] = exists
    }
    return graph
}

// GetPointboards возвращает список позиций камней группы
func (ch *Chain) GetPointboards() []int {
    boards := make([]int, 0, len(ch.ChainMap))
    for pos := range ch.ChainMap {
        boards = append(boards, pos)
    }
    return boards
}