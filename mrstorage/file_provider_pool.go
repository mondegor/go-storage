package mrstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

type (
	// FileProviderPool - пул файловых провайдеров, позволяет
	// хранить файловые провайдеры предназначенные для различных целей.
	FileProviderPool struct {
		providers providerMap
	}

	providerMap map[string]FileProvider
)

// NewFileProviderPool - создаёт объект FileProviderPool.
func NewFileProviderPool() *FileProviderPool {
	return &FileProviderPool{
		providers: make(providerMap),
	}
}

// Register - регистрирует провайдера по его имени.
func (p *FileProviderPool) Register(name string, provider FileProvider) error {
	if _, ok := p.providers[name]; ok {
		return mr.ErrInternal.Wrap(fmt.Errorf("file provider '%s' is already registered", name))
	}

	p.providers[name] = provider

	return nil
}

// ProviderAPI - возвращает API по имени провайдера или ошибку, если он не был найден.
func (p *FileProviderPool) ProviderAPI(name string) (FileProviderAPI, error) {
	if provider, ok := p.providers[name]; ok {
		return provider, nil
	}

	return nil, mr.ErrInternal.Wrap(fmt.Errorf("file provider '%s' is not registered", name))
}

// Ping - сообщает, установлено ли соединение и является ли оно стабильным для всех зарегистрированных провайдеров.
func (p *FileProviderPool) Ping(ctx context.Context) error {
	for name, provider := range p.providers {
		if err := provider.Ping(ctx); err != nil {
			return ErrFileProviderPingError.Wrap(err, name)
		}
	}

	return nil
}

// Close - закрывает текущие соединения всех зарегистрированных провайдеров.
func (p *FileProviderPool) Close() error {
	var errs []error

	for name, provider := range p.providers {
		if providerErr := provider.Close(); providerErr != nil {
			errs = append(errs, fmt.Errorf("provider '%s' close error: %w", name, providerErr))
		}
	}

	if len(errs) > 0 {
		return mr.ErrInternalFailedToClose.Wrap(
			errors.Join(errs...),
			"source", "file provider pool",
		)
	}

	return nil
}
