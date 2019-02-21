package healthchecking

import (
	"net/http"
	"testing"
	"time"
)

func Benchmark_PoolGet(b *testing.B) {
	fac := func() (*HealthChecker, error) {
		return &HealthChecker{
			client: &http.Client{Timeout: 5 * time.Second},
		}, nil
	}
	pool, err := newCheckerPool(100, 1000, fac)
	if err != nil {
		b.Errorf("NewCheckerPool(100, 1000, fac) got: %v", err)
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {
		if _, err := pool.Get(); err != nil {
			b.Errorf("pool.Get() got: %v", err)
			b.FailNow()
		}
	}
}
