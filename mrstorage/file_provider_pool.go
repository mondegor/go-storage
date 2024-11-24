package mrstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
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
		return mrcore.ErrInternalWithDetails.New(fmt.Sprintf("file provider '%s' is already registered", name))
	}

	p.providers[name] = provider

	return nil
}

// ProviderAPI - возвращает API по имени провайдера или ошибку, если он не был найден.
func (p *FileProviderPool) ProviderAPI(name string) (FileProviderAPI, error) {
	if provider, ok := p.providers[name]; ok {
		return provider, nil
	}

	return nil, mrcore.ErrInternalWithDetails.New(fmt.Sprintf("file provider '%s' is not registered", name))
}

// Ping - проверяет работоспособность всех зарегистрированных провайдеров.
func (p *FileProviderPool) Ping(ctx context.Context) error {
	for name, provider := range p.providers {
		if err := provider.Ping(ctx); err != nil {
			return mrcore.ErrInternalWithDetails.Wrap(err, fmt.Sprintf("file provider '%s' ping error", name))
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
		return mrcore.ErrInternalWithDetails.Wrap(errors.Join(errs...), "file provider pool")
	}

	return nil
}
