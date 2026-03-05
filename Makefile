### Команды для запуска ВСЕЙ системы (с postgresql) во время локальной разработки ###

# запуск приложения, postgresql и применение миграций
# команда позволяет запустить систему
startall:
	docker compose up --build

# полная остановка после startall
# останавливает систему и удаляет контейнеры
stopall:
	docker compose down


### Проход по всем тестам ###

test:
	go test -v ./...


### Инструкции для сборки и запуска образа конкретно приложения ###

# переменные приложения
APP_NAME = app
TAG = latest
PORT = 8000

# параметры флагов (-repo и -config)
REPO ?= "inmemory" # inmemory || postgresql
# при изменении на postgresql важно проверить файлы конфигурации и переменных окружения
# все переменные для подключения к БД должны быть корректными
CONFIG_DIR_HOST = $(shell pwd)/configs
CONFIG_DIR_CONTAINER = /app/configs

# сборка образа приложения
buildapp:
	docker build -t $(APP_NAME):$(TAG) .

# запуск образа приложения
runapp:
	docker run --rm -p $(PORT):$(PORT) \
		-v $(CONFIG_DIR_HOST):$(CONFIG_DIR_CONTAINER) \
		$(APP_NAME):$(TAG) \
		-repo=$(REPO) \
		-config=$(CONFIG_DIR_CONTAINER)/config.yaml

# сразу сборка и последующий запуск
upapp: buildapp runapp
