package healthchecking

import (
	"testing"
	"time"
)

func Test_Checker(t *testing.T) {
	job := newHealthJob("http://127.0.0.1:9091/health", "/cluster/1/1", 5*time.Second)
	checker := &HealthChecker{}

	cr := make(chan checkResult)
	checker.Check(job, cr)
	close(cr)

	if r := <-cr; r.IsAlive {
		t.Errorf("want false, got: %v", r.IsAlive)
	}
}
