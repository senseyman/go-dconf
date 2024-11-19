package manager

import (
	"context"
	"log"
	"sync"
	"time"
)

type ConfigManager[T any] struct {
	repo         Repository
	scanInterval time.Duration

	mx     sync.Mutex
	dynCfg T
}

func New[T any](repo Repository, initCfg T, scanInterval time.Duration) *ConfigManager[T] {
	return &ConfigManager[T]{
		dynCfg:       initCfg,
		repo:         repo,
		scanInterval: scanInterval,
	}
}

func (m *ConfigManager[T]) GetConfig() T {
	m.mx.Lock()
	defer m.mx.Unlock()
	return m.dynCfg
}

func (m *ConfigManager[T]) LoadConfig(ctx context.Context) error {
	var cfg T
	if err := m.repo.GetConfig(ctx, &cfg, m.dynCfg); err != nil {
		return err
	}

	m.mx.Lock()
	defer m.mx.Unlock()
	m.dynCfg = cfg

	return nil
}

func (m *ConfigManager[T]) Run(ctx context.Context, wg *sync.WaitGroup) error {
	if err := m.LoadConfig(ctx); err != nil {
		log.Printf("Initial config load failed: %v", err)
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		errorCount := 0
		maxRetries := 5
		backoff := time.Second

		for {
			select {
			case <-ctx.Done():
				log.Println("Configuration manager stopped")
				return
			case <-time.After(m.scanInterval):
				if err := m.LoadConfig(ctx); err != nil {
					log.Printf("Error loading configuration: %v", err)
					errorCount++
					if errorCount >= maxRetries {
						log.Printf("Too many consecutive errors (%d), stopping updates", errorCount)
						return
					}
					time.Sleep(backoff)
					backoff *= 2 // Exponential backoff
				} else {
					errorCount = 0
					backoff = time.Second // Reset backoff on success
				}
			}
		}
	}()

	return nil
}
