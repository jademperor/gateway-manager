package healthchecking

import (
	"net/http"
	"time"

	"github.com/jademperor/gateway-manager/internal/logger"
)

// newHealthJob ....
func newHealthJob(healthCheckURL, key string, dur time.Duration) *HealthJob {
	return &HealthJob{
		ticker:        time.NewTicker(dur),
		lastCheckTime: time.Now(),
		TargetURL:     healthCheckURL,
		InstanceKey:   key,
	}
}

// HealthJob for healthchecking ...
type HealthJob struct {
	ticker        *time.Ticker // time.Ticker
	lastCheckTime time.Time    // maybe not need
	TargetURL     string       // server instance addr
	InstanceKey   string       // to find the instance and change it
	// IsAlive       bool         // flag to mark the instance is available or not
}

// HealthChecker ...
type HealthChecker struct {
	client *http.Client
}

type checkResult struct {
	IsAlive bool
	Key     string
}

// Check server instance is alive or not.
// HTTP response StatusOK(200) marked as success
func (checker *HealthChecker) Check(job *HealthJob, result chan<- checkResult) {
	// if checker.client == nil {
	// 	checker.client = &http.Client{Timeout: 5 * time.Second}
	// }
	select {
	case <-job.ticker.C:
		resp, err := checker.client.Get(job.TargetURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			logger.Logger.Infof("instance[URL: %s, key: %s] is healthy:", job.TargetURL, job.InstanceKey)
			result <- checkResult{IsAlive: true, Key: job.InstanceKey}
			return
		}
		if err != nil {
			logger.Logger.Errorf("(checker *HealthChecker) Check() got err: %v with TargetURL: [%s]", err, job.TargetURL)
		}
		if resp != nil {
			logger.Logger.Errorf("(checker *HealthChecker) Check() got StatusCode: %d", resp.StatusCode)
		}
		result <- checkResult{IsAlive: false, Key: job.InstanceKey}
	default:
	}

}

// Close ...
func (checker *HealthChecker) Close() {
	checker.client = nil
}
