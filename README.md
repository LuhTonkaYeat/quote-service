# Quote Service

gRPC сервис для работы с цитатами, реализованный на Go с использованием Clean Architecture и SQLite.

## О проекте

Quote Service - это микросервис для хранения и получения мотивационных цитат. Проект демонстрирует:
- Clean Architecture подход
- gRPC коммуникацию
- Работу с SQLite
- CLI клиент для тестирования
- Docker контейнеризацию

## Функциональность

- **GetRandom** - получить случайную цитату
- **GetByCategory** - получить случайную цитату из категории
- **AddQuote** - добавить новую цитату


## Быстрый старт

### Требования
- Go 1.21 или выше
- protoc (опционально, для генерации proto файлов)

### Установка и запуск

```bash
# Клонируем репозиторий
git clone https://github.com/LuhTonkaYeat/quote-service.git
cd quote-service

# Устанавливаем зависимости
go mod download

# Генерируем proto файлы (если нужно)
make proto

# Запускаем сервер
make run

# Сервер запустится на localhost:50051
```

### Использование CLI клиента (В другом терминале)

```bash
# Случайная цитата
go run client/main.go random

# Цитата из категории
go run client/main.go category motivation

# Добавить новую цитату
go run client/main.go add "Keep coding!" "Developer" "motivation"

# Помощь
go run client/main.go help
```

### Использование Makefile
```bash
make proto                           # Сгенерировать proto файлы
make run                             # Запустить сервер
make build                           # Собрать бинарник
make random                          # Получить случайную цитату
make category <cat>                  # Получить цитату из категории
make add "text" "author" "category"  # Добавить цитату
make clean                           # Очистить сборку