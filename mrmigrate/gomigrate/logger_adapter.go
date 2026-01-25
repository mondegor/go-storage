package gomigrate

import (
	"context"
	"strings"

	"github.com/mondegor/go-sysmess/mrlog"
)

type (
	// LoggerAdapter - адаптер для работы с логами golang-migrate.
	LoggerAdapter struct {
		logger mrlog.Logger
	}
)

// NewLoggerAdapter - создаёт объект LoggerAdapter.
func NewLoggerAdapter(logger mrlog.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		logger: logger,
	}
}

// Printf - выводит лог информацию о миграции БД.
func (l *LoggerAdapter) Printf(format string, v ...any) {
	l.logger.Info(context.Background(), strings.TrimSpace(format), v...)
}

// Verbose - возвращает можно ли выводить лог миграций БД.
func (l *LoggerAdapter) Verbose() bool {
	return mrlog.InfoEnabled(l.logger)
}
