package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
)

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
}

var ourgame = NewGoTree(GameSettings{BoardSize: 9})

// Получаем абсолютный путь к корневой папке проекта
func getRootPath() string {
	// Получаем путь к текущему файлу (server/ folder)
	_, filename, _, _ := runtime.Caller(0)
	// Возвращаемся на одну папку вверх (из server/ в корень проекта)
	return filepath.Dir(filepath.Dir(filename))
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

// Обработчик API
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	//rootPath := getRootPath()
	//htmlPath := filepath.Join(rootPath, "index.html")
	//fmt.Printf("spi Serving HTML from: %s\n", htmlPath)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed, bro?"}`, http.StatusMethodNotAllowed)
		fmt.Print("wrong method")
		return
	}

	// читаем тело запроса
	var request ApiRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "badm json"}`, http.StatusBadRequest)
		return
	}

	// выводим на сервер
	fmt.Println("Получено от клиента:")
	fmt.Println("Name:", request.Name)
	fmt.Println("Action: ", request.Action)
	fmt.Printf("Move: %+v\n", request.Move)
	if request.Move.X < 0 || request.Move.X > ourgame.BoardSize-1 || request.Move.Y < 0 || request.Move.Y > ourgame.BoardSize-1 {
		http.Error(w, `{"error": "move outside bounds"}`, http.StatusBadRequest)
		return
	}
	move_result := err
	if request.Action == "clear board" {
		ourgame = NewGoTree(GameSettings{})
	} else {
		move_result = ourgame.MakeMove(request.Move.Y*ourgame.BoardSize + request.Move.X)
	}
	//index = request.Move %
	fmt.Println("Отправляем на клиент:")
	fmt.Print(stringBoard(ourgame))
	fmt.Println(string(ourgame.CurrentNode.LastMoveColor.Opposite()))
	// формируем ответ
	response := map[string]interface{}{
		"status":       "ok",
		"msg":          "json received",
		"accepted":     move_result == nil,
		"boardState":   stringBoard(ourgame), // возвращаем клиенту
		"playingColor": string(ourgame.CurrentNode.LastMoveColor.Opposite()),
	}

	json.NewEncoder(w).Encode(response)
}

// Обработчик формы
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	age := r.FormValue("age")
	tmpl := `
    <html>
        <body>
            <h1>Полученные данные:</h1>
            <p>Имя: {{.Name}}</p>
            <p>Возраст: {{.Age}}</p>
            <a href="/">Назад</a>
        </body>
    </html>
    `

	data := struct {
		Name string
		Age  string
	}{name, age}

	t := template.Must(template.New("result").Parse(tmpl))
	t.Execute(w, data)
}

// Обработчик калькулятора
func calcHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")

	a, _ := strconv.Atoi(aStr)
	b, _ := strconv.Atoi(bStr)

	result := a + b
	fmt.Fprintf(w, "Сумма %d + %d = %d", a, b, result)
}

func runServer() {
	// Регистрируем обработчики

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/calc", calcHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Структура проекта:")
	fmt.Printf("Корень: %s\n", getRootPath())
	http.ListenAndServe(":8080", nil)
}
