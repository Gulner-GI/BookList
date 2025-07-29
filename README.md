# Simple BookList

Этот проект - это элементарный RESTful API для создания и управления списком книг. Реализован на Golang с использованием фреймворка Gin, с базой данных SQlite3.
Проект был сделан в учебных целях, для оттачивания собственных навыков программирования
Стек технологий:
 Язык: Go  
 Веб‑фреймворк: Gin  
 БД: SQLite3  
 Документация: swaggo (Swagger UI)  
 Тестирование: встроенный testing  
 CLI: Bash скрипты
 
Установка и запуск
   bash
1. Клонировать репозиторий
git clone https://github.com/Gulner-GI/BookList.git
cd BookList
2. Загрузить зависимости
go mod download
3. Сделать скрипт исполняемым
chmod +x run
4. Запустить сервер
./run start

Структура проекта:
 BookList/
 config/config.go                      # конфигурация через ENV
 db/db.go                              # инициализация и миграции базы
 docs/docs.go                          # сгенерированная Swagger‑документация
     /swagger.json
     /swagger.yaml
 handlers/handlers.go                   # HTTP‑хендлеры и тесты
         /helpers_test.go
         /handlers_integration_test.go
 loggers/logger.go                      # настройка логирования
 models/models.go                       # структуры Book, Update, ErrorResponse
 routes/router.go                       # роутинг Gin
 run                                    # обёртка для запуска скриптов
 scripts/start.sh                       # CLI‑скрипты: get.sh, add.sh, update.sh, delete.sh
        /help.sh
        /get.sh
        /create.sh
        /update.sh
        /delete.sh
        /end.sh
 main.go                                # точка входа
 books.db                               # база данных
 go.sum            
 go.mod            
 server.log                             # лог файл (генерируется на запуске)
 README.md                              # этот файл

Функционал
Реализован полный CRUD функционал для получения списка книг, добавления книг в список, изменения данных о книгах, удаления книг из списка.
Для прямой работы через HTTP можно использовать curl (см. документацию Swagger UI)
Реализован CLI-интерфейс на bash:
./run start                       # запускает сервер
./run help                        # выдаёт список возможных команд
./run get   [<id>]                # выдаёт список из всех книг, либо книгу с указанным id
./run create <title> <year> <genre> <true|false> [--link <url>]     # создаёт новую книгу
./run update --id <id> [--title <title>] [--year <year>] [--genre <genre>] [--status <true|false>] [--link <url>]    # частично изменяет данные
./run delete <id>                     # удаляет книгу по указанному id
./run end                         # завершает работу сервера
Доступна swagger-документация по адресу http://localhost:8080/swagger/index.html
Реализовано логирование, все логи пишутся в server.log
Реализованы unit-тесты для вспомогательных функций (helpers_test.go), и интеграционные тесты полного CRUD цикла (handlers_integration_test.go)

Конфигурация
По умолчанию: порт 8080, путь к файлу SQLite - DB_PATH указывает на books.db, режим работы GIN_MODE=debug

