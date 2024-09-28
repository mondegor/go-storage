package mrstorage

import (
	"context"
	"errors"
	"fmt"
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
		return fmt.Errorf("file provider '%s' is already registered", name)
	}

	p.providers[name] = provider

	return nil
}

// ProviderAPI - возвращает API провайдера по его имени или ошибку, если он не был найден.
func (p *FileProviderPool) ProviderAPI(name string) (FileProviderAPI, error) {
	if provider, ok := p.providers[name]; ok {
		return provider, nil
	}

	return nil, fmt.Errorf("file provider '%s' is not registered", name)
}

// Ping - проверяет работоспособность всех зарегистрированных провайдеров.
func (p *FileProviderPool) Ping(ctx context.Context) error {
	for name, provider := range p.providers {
		if err := provider.Ping(ctx); err != nil {
			return fmt.Errorf("file provider '%s' ping error: %w", name, err)
		}
	}

	return nil
}

// Close - закрывает текущие соединения всех зарегистрированных провайдеров.
func (p *FileProviderPool) Close() (err error) {
	for name, provider := range p.providers {
		if providerErr := provider.Close(); providerErr != nil {
			if err == nil {
				err = errors.New("file provider poll")
			}

			err = fmt.Errorf("%w; provider '%s' close error: %w", err, name, providerErr)
		}
	}

	return err
}
