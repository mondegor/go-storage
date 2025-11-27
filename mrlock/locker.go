package mrlock

import (
	"context"
	"time"
)

//go:generate mockgen -source=locker.go -destination=./mock/locker.go

const (
	// DefaultExpiry - истечение срока блокировки по умолчанию.
	DefaultExpiry = time.Second
)

type (
	// Locker - блокировщик указанного ключа, который возвращает функцию для его разблокирования.
	Locker interface {
		Lock(ctx context.Context, key string) (unlock func(), err error)
		LockWithExpiry(ctx context.Context, key string, expiry time.Duration) (unlock func(), err error)
	}
)
