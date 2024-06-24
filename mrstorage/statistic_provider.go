package mrstorage

import (
	"time"
)

type (
	// DBStatProvider - провайдер статистики работы DB.
	DBStatProvider interface {
		AcquireCount() int64
		AcquireDuration() time.Duration
		AcquiredConns() int32
		CanceledAcquireCount() int64
		ConstructingConns() int32
		EmptyAcquireCount() int64
		IdleConns() int32
		MaxConns() int32
		TotalConns() int32
		NewConnsCount() int64
		MaxLifetimeDestroyCount() int64
		MaxIdleDestroyCount() int64
	}
)
