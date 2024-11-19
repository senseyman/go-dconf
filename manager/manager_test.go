package manager_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/senseyman/go-dconf/manager"
	"github.com/senseyman/go-dconf/manager/mock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	testCfg := struct {
		LogLevel string
		InitFee  float32
	}{
		LogLevel: "debug",
		InitFee:  0.3,
	}

	m := manager.New(repo, testCfg, time.Minute)

	require.NotNil(t, m)
}

func TestConfigManager_GetConfig(t *testing.T) {
	t.Parallel()

	testCfg := struct {
		LogLevel string
		InitFee  float32
	}{
		LogLevel: "debug",
		InitFee:  0.3,
	}

	m := manager.New(nil, testCfg, time.Minute)
	require.NotNil(t, m)

	cfg := m.GetConfig()
	require.Equal(t, testCfg, cfg)
}

func TestConfigManager_LoadConfig(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	repo := mock.NewMockRepository(ctrl)
	ctx := context.Background()

	testCfg := struct {
		LogLevel string
		InitFee  float32
	}{
		LogLevel: "debug",
		InitFee:  0.3,
	}

	repo.EXPECT().GetConfig(ctx, gomock.AnyOf(&struct {
		LogLevel string
		InitFee  float32
	}{}), gomock.AnyOf(testCfg)).Return(nil)

	m := manager.New(repo, testCfg, time.Minute)
	require.NotNil(t, m)

	err := m.LoadConfig(ctx)
	require.NoError(t, err)
}
