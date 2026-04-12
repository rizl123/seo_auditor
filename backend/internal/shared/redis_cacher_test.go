package shared

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisCacher(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	cacher := NewRedisCacher(mr.Addr())
	defer cacher.Close()

	ctx := context.Background()

	t.Run("DataIntegrity_JSONSerialization", func(t *testing.T) {
		type complexData struct {
			Map   map[string]int `json:"map"`
			Slice []string       `json:"slice"`
		}
		input := complexData{
			Map:   map[string]int{"key": 1},
			Slice: []string{"a", "b"},
		}

		err := cacher.Store(ctx, "integrity", "id1", input, time.Minute)
		assert.NoError(t, err)

		var output complexData
		err = cacher.Fetch(ctx, "integrity", "id1", &output)
		assert.NoError(t, err)
		assert.Equal(t, input, output)
	})

	t.Run("Concurrency_ThreadSafety", func(t *testing.T) {
		const goroutines = 10
		const iterations = 20
		var wg sync.WaitGroup

		wg.Add(goroutines)
		for i := range goroutines {
			go func(id int) {
				defer wg.Done()
				for j := range iterations {
					key := "perf_key"
					_ = cacher.Store(ctx, "perf", key, id+j, time.Minute)
					var res int
					_ = cacher.Fetch(ctx, "perf", key, &res)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("ErrorHandling_BrokenData", func(t *testing.T) {
		mr.Set("broken:key", "invalid-json-content")

		var res map[string]string
		err := cacher.Fetch(ctx, "broken", "key", &res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal")
	})

	t.Run("Reliability_ConnectionLoss", func(t *testing.T) {
		tmpMr, err := miniredis.Run()
		assert.NoError(t, err)

		tmpCacher := NewRedisCacher(tmpMr.Addr())

		err = tmpCacher.PingWithTimeout(time.Second)
		assert.NoError(t, err)

		tmpMr.Close()

		err = tmpCacher.PingWithTimeout(time.Millisecond * 50)
		assert.Error(t, err)
	})

	t.Run("Lifecycle_Expiration", func(t *testing.T) {
		key := "exp_key"
		err := cacher.Store(ctx, "lifecycle", key, "val", time.Second)
		assert.NoError(t, err)

		mr.FastForward(2 * time.Second)

		var res string
		err = cacher.Fetch(ctx, "lifecycle", key, &res)
		assert.ErrorIs(t, err, ErrCacheMiss)
	})

	t.Run("Lifecycle_ContextCancellation", func(t *testing.T) {
		cCtx, cancel := context.WithCancel(context.Background())
		cancel()

		err := cacher.Store(cCtx, "ctx", "key", 1, time.Minute)
		assert.Error(t, err)
	})

	t.Run("Collision_Analysis", func(t *testing.T) {
		_ = cacher.Store(ctx, "user:1", "data", "val1", time.Minute)
		_ = cacher.Store(ctx, "user", "1:data", "val2", time.Minute)

		var res string
		_ = cacher.Fetch(ctx, "user:1", "data", &res)
		assert.Equal(t, "val2", res, "Current implementation is subject to key collisions")
	})
}
