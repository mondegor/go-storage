# Описание GoStorage v0.16.0
Этот репозиторий содержит описание библиотеки GoStorage.

## Статус библиотеки
Библиотека находится в стадии разработки.

## Описание библиотеки
Библиотека для работы с хранилищами данных.
На данный момент реализованы адаптеры для следующих клиентов:
- postgres (`pgx/v5`);
- rabbitmq (`amqp-go/v1`);
- redis (`go-redis/v9` + `redislock/v0.9`);
- `S3 minio` + `FileProvider`;
- `Native File System` + `FileProvider`;

Также реализован вспомогательный пакет для тестирования приложения без необходимости
разворачивания инфраструктуры (БД, S3 и т.д.), для этого используется библиотека `testcontainers`.
На базе этой библиотеки подготовлены следующие контейнеры:
- postgres;
- redis;
- minio;

## Подключение библиотеки
`go get -u github.com/mondegor/go-storage@v0.16.0`

## Установка библиотеки для её локальной разработки
- Выбрать рабочую директорию, где должна быть расположена библиотека
- `mkdir go-storage && cd go-storage` // создать и перейти в директорию проекта
- `git clone git@github.com:mondegor/go-storage.git .`
- `cp .env.dist .env`
- `mrcmd go-dev deps` // загрузка зависимостей проекта
- Для работы утилит `gofumpt`, `goimports`, `mockgen` необходимо в `.env` проверить
  значения переменных `GO_DEV_TOOLS_INSTALL_*` и запустить `mrcmd go-dev install-tools`

### Консольные команды используемые при разработке библиотеки

> Перед запуском консольных скриптов библиотеки необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

- `mrcmd go-dev help` // выводит список всех доступных go-dev команд;
- `mrcmd go-dev generate` // генерирует go файлы через встроенный механизм go:generate;
- `mrcmd go-dev gofumpt-fix` // исправляет форматирование кода (`gofumpt -l -w -extra ./`);
- `mrcmd go-dev goimports-fix` // исправляет imports, если это требуется (`goimports -d -local ${GO_DEV_IMPORTS_LOCAL_PREFIXES} ./`);
- `mrcmd golangci-lint check` // запускает линтеров для проверки кода (на основе `.golangci.yaml`);
- `mrcmd go-dev test` // запускает тесты библиотеки;
- `mrcmd go-dev test-report` // запускает тесты библиотеки с формированием отчёта о покрытии кода (`test-coverage-full.html`);
- `mrcmd plantuml build-all` // генерирует файлы изображений из `.puml` [подробнее](https://github.com/mondegor/mrcmd-plugins/blob/master/plantuml/README.md#%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0-%D1%81-%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D0%B5%D0%B9-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0-markdown--plantuml);

#### Короткий вариант выше приведённых команд (Makefile)
- `make deps` // аналог `mrcmd go-dev deps`
- `make generate` // аналог `mrcmd go-dev generate`
- `make fmt` // аналог `mrcmd go-dev gofumpt-fix`
- `make fmti` // аналог `mrcmd go-dev goimports-fix`
- `make lint` // аналог `mrcmd golangci-lint check`
- `make test` // аналог `mrcmd go-dev test`
- `make test-report` // аналог `mrcmd go-dev test-report`
- `make plantuml` // аналог `mrcmd plantuml build-all`

> Чтобы расширить список команд, необходимо создать Makefile.mk и добавить
> туда дополнительные команды, все они будут добавлены в единый список команд make утилиты.