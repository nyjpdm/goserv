package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type User struct {
	Name   string          `json:"name"`
	Status string          `json:"status"`
	Conn   *websocket.Conn `json:"-"`
}
type GameRoom struct {
	board       *GoTree
	whitePlayer *User
	blackPlayer *User
	watching    []*User
	Status      string
	Name        string
}
type MoveQuery struct {
	Player     string `json:"player"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Coordinate string `json:"coordinate"`
	Timestamp  string `json:"timestamp"`
}
type ApiRequest struct {
	Name   string    `json:"name"`
	Move   MoveQuery `json:"move"`
	Action string    `json:"action"`
}

type Server struct {
	players []*User
	games   []*GameRoom
}

var ourServer Server

//var ourgame = NewGoTree(GameSettings{BoardSize: 9})

func applyMoveQuery(request ApiRequest, tree *GoTree) error {
	if request.Move.X < -2 || request.Move.X > tree.BoardSize-1 || request.Move.Y < 0 || request.Move.Y > tree.BoardSize-1 {
		//	http.Error(w, `{"error": "move outside bounds"}`, http.StatusBadRequest)
		return nil
	}
	if request.Action == "clear board" {
		tree = NewGoTree(GameSettings{})
		return nil
	} else {
		move_result := tree.MakeMove(request.Move.Y*tree.BoardSize + request.Move.X)
		return move_result
	}

}

// Получаем абсолютный путь к корневой папке проекта
func getRootPath() string {
	// Получаем путь к текущему файлу (server/ folder)
	_, filename, _, _ := runtime.Caller(0)
	// Возвращаемся на одну папку вверх (из server/ в корень проекта)
	return filepath.Dir(filepath.Dir(filename))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Проверка Origin - разрешаем только localhost и наш домен
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		// Разрешаемые источники
		allowedOrigins := []string{
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"https://mydomain.com",
		}
		return true //DEBUG FIX LATER
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}

		// В режиме разработки можно логировать, но запрещать
		log.Printf("Запрещённый origin: %s", origin)
		return false

		// ИЛИ для разработки можно временно разрешить все:
		// return true
	},
}

// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	rootPath := getRootPath()
	// Теперь путь будет корень_проекта/index.html
	htmlPath := filepath.Join(rootPath, "index.html")

	// Для отладки можно добавить:
	fmt.Printf("Serving HTML from: %s\n", htmlPath)

	http.ServeFile(w, r, htmlPath)
}

func runServer() {
	r := mux.NewRouter()

	// Регистрируем обработчики

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/ws", handleWebSocket)
	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Структура проекта:")
	fmt.Printf("Корень: %s\n", getRootPath())
	var firstRoom = GameRoom{Name: "our Game", Status: "ongoing", blackPlayer: nil, whitePlayer: nil}
	firstRoom.board = NewGoTree(GameSettings{BoardSize: 9})
	ourServer = Server{games: []*GameRoom{}, players: []*User{}}
	ourServer.games = append(ourServer.games, &firstRoom)
	http.ListenAndServe(":8080", r)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	thisPlayer := User{}
	thisPlayer.Conn = conn
	thisPlayer.Name = "player_" + strconv.Itoa(len(ourServer.players))
	ourServer.players = append(ourServer.players, &thisPlayer)
	fmt.Println("number of players:", len(ourServer.players))
	var playColor PointColor //создаем игрока (пока при каждом новом обращении к вебсокету)
	if ourServer.games[0].blackPlayer == nil {
		ourServer.games[0].blackPlayer = &thisPlayer
		playColor = black
	} else if ourServer.games[0].whitePlayer == nil {
		ourServer.games[0].whitePlayer = &thisPlayer
		playColor = white
	}
	mess := map[string]interface{}{ //отправляем данные игрока, которые (пока) создает сервер
		"status":       "ok",
		"msg":          "your username",
		"playingColor": playColor,
		"username":     thisPlayer.Name,
		"boardState":   stringBoard(ourServer.games[0].board),
	}
	messg, _ := json.Marshal(mess)
	conn.WriteMessage(websocket.TextMessage, messg)

	go handleClient(conn, thisPlayer)

}

func handleClient(conn *websocket.Conn, thisPlayer User) {

	defer conn.Close()

	for {
		// Чтение сообщения от клиента
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Парсим JSON от клиента
		var request ApiRequest
		if err := json.Unmarshal(msgBytes, &request); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}
		fmt.Println(request.Move)
		res := applyMoveQuery(request, ourServer.games[0].board)
		fmt.Println(res)

		fmt.Println("Отправляем на клиент:")
		fmt.Print(stringBoard(ourServer.games[0].board))
		fmt.Println(string(ourServer.games[0].board.CurrentNode.LastMoveColor.Opposite()))
		if ourServer.games[0].blackPlayer != nil {
			sendBoard(ourServer.games[0].board, ourServer.games[0].board.CurrentNode.LastMoveColor.Opposite(), ourServer.games[0].blackPlayer.Conn, res == nil)
		}
		if ourServer.games[0].whitePlayer != nil {
			sendBoard(ourServer.games[0].board, ourServer.games[0].board.CurrentNode.LastMoveColor, ourServer.games[0].whitePlayer.Conn, res == nil)
		}

		sendBoard(ourServer.games[0].board, ourServer.games[0].board.CurrentNode.LastMoveColor.Opposite(), conn, res == nil)
		// response := map[string]interface{}{
		// 	"status":       "ok",
		// 	"msg":          "updating board",
		// 	"accepted":     res == nil,
		// 	"boardState":   stringBoard(ourgame.board), // возвращаем клиенту
		// 	"playingColor": string(ourgame.board.CurrentNode.LastMoveColor.Opposite()),
		// 	"blackScore":   ourgame.board.CurrentNode.BlackCaptures,
		// 	"whiteScore":   ourgame.board.CurrentNode.WhiteCaptures,
		// }
		// resmess, _ := json.Marshal(response)
		// conn.WriteMessage(websocket.TextMessage, resmess)

	}

}

func sendBoard(board *GoTree, color PointColor, conn *websocket.Conn, accepted bool) {
	response := map[string]interface{}{
		"status":       "ok",
		"msg":          "updating board",
		"accepted":     accepted,
		"boardState":   stringBoard(board), // возвращаем клиенту
		"playingColor": color,
		"blackScore":   board.CurrentNode.BlackCaptures,
		"whiteScore":   board.CurrentNode.WhiteCaptures,
	}
	resmess, _ := json.Marshal(response)
	conn.WriteMessage(websocket.TextMessage, resmess)
}
