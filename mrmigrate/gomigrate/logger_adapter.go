package gomigrate

import (
	"context"
	"strings"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// LoggerAdapter - адаптер для интеграции логгера приложения с библиотекой golang-migrate.
	LoggerAdapter struct {
		logger mrlog.Logger
	}
)

// NewLoggerAdapter - создаёт объект LoggerAdapter для интеграции с golang-migrate.
func NewLoggerAdapter(logger mrlog.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		logger: logger,
	}
}

// Printf - реализует метод интерфейса migrate.Logger для вывода сообщений миграции.
// Использует уровень логирования Info для всех сообщений.
func (l *LoggerAdapter) Printf(format string, v ...any) {
	l.logger.Info(context.Background(), strings.TrimSpace(format), v...)
}

// Verbose - реализует метод интерфейса migrate.Logger для проверки уровня детализации логов.
// Возвращает true, если основной логгер поддерживает уровень Info, что позволяет
// библиотеке golang-migrate выводить подробные сообщения о процессе миграции.
func (l *LoggerAdapter) Verbose() bool {
	return mrlog.InfoEnabled(l.logger)
}
