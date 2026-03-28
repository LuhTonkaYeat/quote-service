.PHONY: proto run build docker clean

# Переменные
PROTO_PATH = api/proto
PROTO_FILES = $(wildcard $(PROTO_PATH)/*.proto)
GO_OUT = api/proto

# Генерация proto файлов
proto:
	@echo "Generating protobuf files..."
	protoc --proto_path=$(PROTO_PATH) \
		--go_out=$(GO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Done!"

# Запуск сервера
run:
	go run cmd/server/main.go

# Сборка бинарника
build:
	go build -o bin/server cmd/server/main.go

# Сборка Docker образа
docker:
	docker build -t quote-service .

# Очистка
clean:
	rm -rf bin/
	go clean

# Запуск тестов
test:
	go test -v ./...

# Установка зависимостей
deps:
	go mod download
	go mod tidy