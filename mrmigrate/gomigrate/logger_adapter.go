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
func NewLoggerAdapter(l mrlog.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		logger: l,
	}
}

// Printf - выводит лог информацию о миграции БД.
func (a *LoggerAdapter) Printf(format string, v ...any) {
	a.logger.Info(context.Background(), strings.TrimSpace(format), v...)
}

// Verbose - возвращает можно ли выводить лог миграций БД.
func (a *LoggerAdapter) Verbose() bool {
	return a.logger.Enabled(mrlog.LevelInfo)
}
