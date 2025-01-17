package config_test

import (
	"context"
	"testing"

	"github.com/nholuongut/terragrunt/config"
	"github.com/nholuongut/terragrunt/internal/cache"
	"github.com/stretchr/testify/assert"
)

const testCacheName = "TerragruntConfig"

func TestTerragruntConfigCacheCreation(t *testing.T) {
	t.Parallel()

	cache := cache.NewCache[config.TerragruntConfig](testCacheName)

	assert.NotNil(t, cache.Mutex)
	assert.NotNil(t, cache.Cache)

	assert.Empty(t, cache.Cache)
}

func TestTerragruntConfigCacheOperation(t *testing.T) {
	t.Parallel()

	testCacheKey := "super-safe-cache-key"

	ctx := context.Background()
	cache := cache.NewCache[config.TerragruntConfig](testCacheName)

	actualResult, found := cache.Get(ctx, testCacheKey)

	assert.False(t, found)
	assert.Empty(t, actualResult)

	stubTerragruntConfig := config.TerragruntConfig{
		IsPartial: true, // Any random property will be sufficient
	}

	cache.Put(ctx, testCacheKey, stubTerragruntConfig)
	actualResult, found = cache.Get(ctx, testCacheKey)

	assert.True(t, found)
	assert.NotEmpty(t, actualResult)
	assert.Equal(t, stubTerragruntConfig, actualResult)
}
