package manager

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type DynamicConfig struct {
}

type ConfigManager struct {
	repo Repository

	mx     sync.RWMutex
	dynCfg DynamicConfig
}

func New(repo Repository) *ConfigManager {
	return &ConfigManager{repo: repo}
}

func (m *ConfigManager) GetConfig() DynamicConfig {
	m.mx.RLock()
	defer m.mx.RUnlock()
	return m.dynCfg
}

func (m *ConfigManager) LoadConfig(ctx context.Context) error {
	var cfg DynamicConfig
	if err := m.repo.GetConfig(ctx, &cfg); err != nil {
		return err
	}

	m.mx.Lock()
	defer m.mx.Unlock()
	m.dynCfg = cfg

	return nil
}

func (m *ConfigManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				if err := m.LoadConfig(ctx); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}
