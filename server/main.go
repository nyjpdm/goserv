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

// Получаем абсолютный путь к корневой папке проекта
func getRootPath() string {
	// Получаем путь к текущему файлу
	_, filename, _, _ := runtime.Caller(0)
	// Возвращаемся на две папки вверх (из server/ в корень)
	return filepath.Dir(filepath.Dir(filename))
}

// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	rootPath := getRootPath()
	htmlPath := filepath.Join(rootPath, "static", "../index.html")
	http.ServeFile(w, r, htmlPath)
}

// Обработчик API
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Гость"
	}
	fmt.Println("client board")
	fmt.Println(name)
	responseData := map[string]string{
		"message": fmt.Sprintf("returning board %s", name),
		"status":  "ok",
	}
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
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

func main() {
	// Регистрируем обработчики
	initBoard()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/calc", calcHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Структура проекта:")
	fmt.Printf("Корень: %s\n", getRootPath())
	//http.ListenAndServe(":8080", nil)
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

func initBoard() {
	tree := NewGoTree(GameSettings{})

	// Создаем тестовую позицию с группой
	tree.Current.Position[3*19+3] = black // D4
	tree.Current.Position[3*19+4] = black // E4
	tree.Current.Position[4*19+3] = black // D5

	fmt.Println("Тестовая доска:")
	tree.Current.PrintBoardSimple(19)

	// Тестируем создание группы
	group := NewGroupFromPosition(tree.Current.Position, 3*19+3, 19)
	if group != nil {
		fmt.Printf("\n=== Информация о группе ===\n")
		fmt.Printf("Цвет: %c\n", group.Color)
		fmt.Printf("Камней: %d\n", group.NumberOfStones)
		fmt.Printf("Дыханий: %d\n", group.Dame)
		fmt.Printf("Жива: %t\n", group.IsAlive())
		fmt.Printf("Мертва: %t\n", group.IsDead())

	}

	// Тестируем с мертвой группой (окруженной)
	fmt.Println("\n=== Тест мертвой группы ===")
	tree2 := NewGoTree(GameSettings{BoardSize: 5})
	// Создаем окруженную группу
	tree2.Current.Position[1*5+1] = black
	tree2.Current.Position[1*5+2] = white
	tree2.Current.Position[1*5+3] = white
	tree2.Current.Position[2*5+1] = white
	tree2.Current.Position[2*5+2] = black
	tree2.Current.Position[2*5+3] = white
	tree2.Current.Position[3*5+1] = white
	tree2.Current.Position[3*5+2] = white
	tree2.Current.Position[3*5+3] = white

	tree2.Current.PrintBoardSimple(5)

	deadGroup := NewGroupFromPosition(tree2.Current.Position, 2*5+2, 5)
	if deadGroup != nil {
		fmt.Printf("Мертвая группа - дыханий: %d, жива: %t\n",
			deadGroup.Dame, deadGroup.IsAlive())
	}
}
