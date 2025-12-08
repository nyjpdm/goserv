// Глобальные переменные
        let boardState = [];
        let currentPlayer = 'B'; // 'B' для черных, 'W' для белых
        const BOARD_SIZE = 9;
        let moveHistory = [];
        let score = ""
        let socket;
        
        function connect() {
            // Создаем WebSocket соединение
            socket = new WebSocket("ws://localhost:8080/ws");
            
            socket.onopen = function(e) {
                console.log("Connected to server");
            };
            
            socket.onmessage = function(event) {
                const data = JSON.parse(event.data);
                console.log("Received:", data);
                handleServerResponse(data)
            };
            
            socket.onclose = function(event) {
                console.log("Connection closed");
                // Пытаемся переподключиться через 3 секунды
                setTimeout(connect, 3000);
            };
            
            socket.onerror = function(error) {
                console.log("Error:", error);
            };
        }
        window.onload = connect;
        function sendJSON(data, actq) {
            if (socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({
                        name: 'hello',
                        move: data,
                        action: actq,
                    }));
                return true;
            }
            return false;
        }

        // Инициализация доски
        function initializeBoard() {
            const board = document.getElementById('board');
            board.innerHTML = '';
            
            // Добавляем координаты
            addCoordinates();
            
            // Создаем клетки доски
            for (let y = 0; y < BOARD_SIZE; y++) {
                for (let x = 0; x < BOARD_SIZE; x++) {
                    const cell = document.createElement('div');
                    cell.className = 'cell';
                    cell.dataset.x = x;
                    cell.dataset.y = y;
                    
                    // Добавляем обработчики hover
                    cell.addEventListener('mouseenter', () => handleCellHover(x, y));
                    cell.addEventListener('mouseleave', handleCellHoverLeave);
                    
                    cell.onclick = () => handleCellClick(x, y);
                    board.appendChild(cell);
                }
            }

            
            // Инициализируем начальное состояние доски
            resetBoard();
            updateBoardDisplay();
            updatePlayerIndicator();
        }

        // Добавление координат на доску
        function addCoordinates() {
            const board = document.getElementById('board');
            const coordinates = 'ABCDEFGHIJKLMNOPQRST'.slice(0, BOARD_SIZE);
            
            // Горизонтальные координаты (буквы) - ПЕРЕМЕСТИТЬ ВНУТРЬ ДОСКИ
            for (let i = 0; i < BOARD_SIZE; i++) {
                const coordX = document.createElement('div');
                coordX.className = 'coordinates coord-x';
                coordX.textContent = coordinates[i];
                // ИЗМЕНИТЬ: использовать 20px вместо 20px, чтобы были ближе к центру
                coordX.style.left = `${(i * 40) + 20}px`;
                coordX.style.bottom = '5px'; // ПОДВИНУТЬ ВВЕРХ ВНУТРИ ДОСКИ
                board.appendChild(coordX);
            }
            
            // Вертикальные координаты (цифры) - ПЕРЕМЕСТИТЬ ВНУТРЬ ДОСКИ
            for (let i = 0; i < BOARD_SIZE; i++) {
                const coordY = document.createElement('div');
                coordY.className = 'coordinates coord-y';
                coordY.textContent = (i + 1).toString();
                // ИЗМЕНИТЬ: использовать 20px вместо 20px, чтобы были ближе к центру
                coordY.style.top = `${(i * 40) + 20}px`;
                coordY.style.right = '5px'; // ПОДВИНУТЬ ВЛЕВО ВНУТРИ ДОСКИ
                board.appendChild(coordY);
            }
        }

        // Сброс доски к пустому состоянию
        function resetBoard() {
            boardState = [];
            for (let y = 0; y < BOARD_SIZE; y++) {
                const row = [];
                for (let x = 0; x < BOARD_SIZE; x++) {
                    row.push('.');
                }
                boardState.push(row);
            }
            moveHistory = [];
            updateMoveLog();
        }

        // Обработчик наведения на клетку
        function handleCellHover(x, y) {
            const cell = document.querySelector(`.cell[data-x="${x}"][data-y="${y}"]`);
            
            // Проверяем, что клетка пуста
            if (boardState[y][x] === '.') {
                // Добавляем класс hover в зависимости от текущего игрока
                if (currentPlayer === 'B') {
                    cell.classList.add('hover-black');
                } else {
                    cell.classList.add('hover-white');
                }
            }
        }

        // Обработчик ухода мыши с клетки
        function handleCellHoverLeave() {
            // Убираем все hover-классы со всех клеток
            const cells = document.querySelectorAll('.cell');
            cells.forEach(cell => {
                cell.classList.remove('hover-black', 'hover-white');
            });
        }

        // Обновите функцию handleCellClick чтобы убирать hover при клике
        async function handleCellClick(x, y) {
            // Сначала убираем hover-эффект
            handleCellHoverLeave();
            
            // Проверяем, что клетка пуста
            if (boardState[y][x] !== '.') {
                console.log(`Клетка (${x}, ${y}) уже занята!`);
                return;
            }
            
            // Остальной код обработки клика остается без изменений
            boardState[y][x] = currentPlayer;
            
            // Логируем ход
            const coordinates = 'ABCDEFGHIJKLMNOPQRST';
            const move = {
                player: currentPlayer,
                x: x,
                y: y,
                coordinate: `${coordinates[x]}${y + 1}`,
                timestamp: new Date().toLocaleTimeString()
            };
            sendJSON(move, "none");
           // const data = await sendRequest(move, "none");
          //  handleServerResponse(data);
            moveHistory.push(move);
            
            console.log(`Ход отправлен: ${currentPlayer} на (${x}, ${y})`);
            
            // Обновляем отображение
            updateBoardDisplay();
            updateMoveLog();
        }

        // Обновление отображения доски
        function updateBoardDisplay() {
            const cells = document.querySelectorAll('.cell');
            
            cells.forEach(cell => {
                // Очищаем клетку от камней
                const existingStone = cell.querySelector('.stone');
                if (existingStone) {
                    existingStone.remove();
                }
                
                const x = parseInt(cell.dataset.x);
                const y = parseInt(cell.dataset.y);
                
                // Добавляем камень если нужно
                if (boardState[y][x] !== '.') {
                    const stone = document.createElement('div');
                    stone.className = `stone ${boardState[y][x] === 'B' ? 'black' : 'white'}`;
                    cell.appendChild(stone);
                }
            });
        }

        // Загрузка доски из текстового поля
        function loadBoardFromText() {
            const textArea = document.getElementById('boardInput');
            const lines = textArea.value.trim().split('\n');
            
            if (lines.length !== BOARD_SIZE) {
                alert(`Ожидается ${BOARD_SIZE} строк, получено ${lines.length}`);
                return;
            }
            
            boardState = [];
            for (let i = 0; i < BOARD_SIZE; i++) {
                const line = lines[i].trim();
                if (line.length !== BOARD_SIZE) {
                    alert(`Строка ${i + 1} должна содержать ${BOARD_SIZE} символов`);
                    return;
                }
                
                const row = line.split('');
                // Проверяем корректность символов
                for (let char of row) {
                    if (!['.', 'B', 'W'].includes(char)) {
                        alert(`Недопустимый символ: ${char}. Допустимы: '.', 'B', 'W'`);
                        return;
                    }
                }
                boardState.push(row);
            }
            
            updateBoardDisplay();
            console.log('Доска загружена из текста');
        }

        // Очистка доски
        async function clearBoard() {

           // const data = await sendRequest(null, "clear board");
           // handleServerResponse(data);
           sendJSON(null, "clear board");
            // Убираем hover-эффект
            handleCellHoverLeave();
            
            resetBoard();
            updateBoardDisplay();
            console.log('Доска очищена');
        }

       // Смена игрока
        function switchPlayer() {
            // Убираем hover-эффект при смене игрока
            handleCellHoverLeave();
            
            currentPlayer = currentPlayer === 'B' ? 'W' : 'B';
            updatePlayerIndicator();
            console.log(`Текущий игрок: ${currentPlayer === 'B' ? 'Черные' : 'Белые'}`);
        }

        // Пас (пропуск хода)
        async function passMove() {
            moveHistory.push({
                player: currentPlayer,
                x: -1,
                y: -1,
                coordinate: 'PASS',
                timestamp: new Date().toLocaleTimeString()
            });
             const move = {
                player: currentPlayer,
                x: -1,
                y: 0,
                coordinate: 'PASS',
                timestamp: new Date().toLocaleTimeString()
            };
            sendJSON(move, "none");
           // const data = await sendRequest(move, "none");
           // handleServerResponse(data);
            console.log(`Игрок ${currentPlayer} пропустил ход`);
            //switchPlayer();
            updateMoveLog();
        }


        // Обновление индикатора текущего игрока
        function updatePlayerIndicator() {
            const indicator = document.getElementById('currentPlayerIndicator');
            const text = document.getElementById('currentPlayerText');
            
            if (currentPlayer === 'B') {
                indicator.className = 'current-player player-black';
                text.textContent = 'Черные; ' + score;
            } else {
                indicator.className = 'current-player player-white';
                text.textContent = 'Белые;' + score;
            }
        }

        // Обновление лога ходов
        function updateMoveLog() {
            const log = document.getElementById('moveLog');
            log.innerHTML = '';
            
            // Показываем последние 10 ходов
            const recentMoves = moveHistory.slice(-10);
            
            recentMoves.forEach(move => {
                const entry = document.createElement('div');
                entry.className = `log-entry ${move.player === 'B' ? 'log-black' : 'log-white'}`;
                
                const playerText = move.player === 'B' ? 'Чёрные' : 'Белые';
                const moveText = move.coordinate === 'PASS' ? 'Пас' : move.coordinate;
                
                entry.textContent = `${move.timestamp} - ${playerText}: ${moveText}`;
                log.appendChild(entry);
            });
            
            // Прокручиваем вниз
            log.scrollTop = log.scrollHeight;
        }

        // Вывод состояния доски в консоль
        function printBoardState() {
            console.log('Текущее состояние доски:');
            console.log('Текстовое представление:');
            console.log(getBoardAsText());
            console.log('Массив:', boardState);
        }

        // Получение текстового представления доски
        function getBoardAsText() {
            return boardState.map(row => row.join('')).join('\n');
        }

        // Показать текстовое представление
        function showBoardText() {
            const textArea = document.getElementById('boardInput');
            textArea.value = getBoardAsText();
        }


        async function sendRequest(move, actq) {
            try {
                const response = await fetch('http://localhost:8080/api', {//await fetch('https://jenyasemagogame.loca.lt/api', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: 'hello',
                        move: move,
                        action: actq,
                    })
                });

                const data = await response.json();
                return data; // просто возвращаем результат
            } catch (error) {
                alert('Ошибка при отправке запроса:', error)
                console.error('Ошибка при отправке запроса:', error);
                return null; // или бросать ошибку, если хочешь обрабатывать выше
            }
        }
        function handleServerResponse(data) {
            if (!data) return;

            console.log('Получен ответ от сервера:', data);

            if(data.msg == "your username"){
                if(data.playingColor != currentPlayer){
                    switchPlayer();
                }
                console.log("username:");
                console.log(data.username);
                
                document.getElementById('boardInput').value = data.boardState;
                loadBoardFromText();

                return;
            }

            if (data.status === "ok" && data.boardState) {
                document.getElementById('boardInput').value = data.boardState;
                loadBoardFromText();
                console.log("Доска обновлена по данным сервера");

                if(data.playingColor != currentPlayer) {
                    switchPlayer();
                }

                score = `captured - B:${data.blackScore}; W:${data.whiteScore}`;
                updatePlayerIndicator();

            } else {
                console.error("Ошибка в ответе сервера:", data);
            }
        }

        // Инициализация при загрузке страницы
        document.addEventListener('DOMContentLoaded', () => {
            initializeBoard();
            console.log('🎮 Go Board Frontend initialized');
            console.log('Доступные команды в консоли:');
            console.log('- loadBoardFromText() - загрузить доску из текста');
            console.log('- clearBoard() - очистить доску');
            console.log('- switchPlayer() - сменить игрока');
            console.log('- printBoardState() - вывести состояние в консоль');
        });

        // Глобальные функции для отладки из консоли
        window.debug = {
            getBoardState: () => boardState,
            setBoardState: (newState) => {
                boardState = newState;
                updateBoardDisplay();
            },
            setCell: (x, y, value) => {
                if (x >= 0 && x < BOARD_SIZE && y >= 0 && y < BOARD_SIZE && ['.', 'B', 'W'].includes(value)) {
                    boardState[y][x] = value;
                    updateBoardDisplay();
                }
            },
            getMoveHistory: () => moveHistory,
            clearHistory: () => {
                moveHistory = [];
                updateMoveLog();
            }
        };