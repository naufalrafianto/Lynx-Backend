package cache

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

type MetricsCacheService struct {
	cache         CacheService
	hits          uint64
	misses        uint64
	errors        uint64
	totalRequests uint64
	totalDuration int64
}

func NewMetricsCacheService(cache CacheService) *MetricsCacheService {
	return &MetricsCacheService{
		cache: cache,
	}
}

func (m *MetricsCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	atomic.AddUint64(&m.totalRequests, 1)
	start := time.Now()

	err := m.cache.Get(ctx, key, dest)

	atomic.AddInt64(&m.totalDuration, time.Since(start).Microseconds())

	if err != nil {
		if err == redis.Nil {
			atomic.AddUint64(&m.misses, 1)
		} else {
			atomic.AddUint64(&m.errors, 1)
		}
		return err
	}

	atomic.AddUint64(&m.hits, 1)
	return nil
}

func (m *MetricsCacheService) GetMetrics() map[string]interface{} {
	totalReq := atomic.LoadUint64(&m.totalRequests)
	avgDuration := float64(0)
	if totalReq > 0 {
		avgDuration = float64(atomic.LoadInt64(&m.totalDuration)) / float64(totalReq)
	}

	return map[string]interface{}{
		"hits":            atomic.LoadUint64(&m.hits),
		"misses":          atomic.LoadUint64(&m.misses),
		"errors":          atomic.LoadUint64(&m.errors),
		"total_requests":  totalReq,
		"hit_ratio":       float64(atomic.LoadUint64(&m.hits)) / float64(totalReq),
		"avg_duration_μs": avgDuration,
	}
}
