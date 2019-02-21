package healthchecking

import (
	"context"
	// "encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/jademperor/common/configs"
	"github.com/jademperor/common/etcdutils"
	"github.com/jademperor/common/models"
	"github.com/jademperor/gateway-manager/internal/logger"
	"go.etcd.io/etcd/client"
)

var (
	clusterWatcher *etcdutils.Watcher    // clusterWatcher for update taskQ
	taskQ          map[string]*HealthJob // taskQ is map of server instance health cheker
	taskQMutex     sync.RWMutex          // read write locker for taskQ
	// HealthCheckTTL (second)
	HealthCheckTTL = 10 * time.Second // default health job ticker duration
)

// Init ...
func Init(etcdAddrs []string, watchDuration time.Duration) {
	store, err := etcdutils.NewEtcdStore(etcdAddrs)
	if err != nil {
		panic(err)
	}

	taskQ = make(map[string]*HealthJob)
	taskQMutex = sync.RWMutex{}

	if err := initTaskQ(store.Kapi); err != nil {
		panic(err)
	}

	// while clusters instance changed
	clusterWatcher = etcdutils.NewWatcher(store.Kapi, watchDuration, configs.ClustersKey)
	go clusterWatcher.Watch(clusterWatchCallback)

	checkerPool, err := newCheckerPool(10, 100, defaultChekerFactory)
	if err != nil {
		panic(err)
	}

	go healthChecking(store, checkerPool)
}

// "/clusters/{clusterID}/{instanceID}"
func isInstanceKey(key string) bool {
	ks := strings.Split(key, "/")
	if len(ks) == 4 && ks[3] != configs.ClusterOptionsKey {
		return true
	}
	return false
}

func initTaskQ(kapi client.KeysAPI) error {
	resp, err := kapi.Get(context.Background(), configs.ClustersKey, nil)
	if err != nil {
		return err
	} else if !resp.Node.Dir {
		return errors.New("configs.ClustersKey is not a dir")
	}

	for _, clusterNode := range resp.Node.Nodes {
		// clusterID := strings.Split(clusterNode.Key, "/")[2]
		clsOpt := new(models.ClusterOption)
		if resp2, err := kapi.Get(context.Background(), clusterNode.Key, nil); err == nil && resp2.Node.Dir {
			for _, srvInsNode := range resp2.Node.Nodes {
				// skip the option node
				if strings.Split(srvInsNode.Key, "/")[3] == configs.ClusterOptionsKey {
					if err := etcdutils.Decode(srvInsNode.Value, clsOpt); err != nil {
						logger.Logger.Error(err)
					}
					continue
				}

				srvInsCfg := new(models.ServerInstance)
				if err := etcdutils.Decode(srvInsNode.Value, srvInsCfg); err != nil {
					logger.Logger.Error(err)
					continue
				}

				// if need check health of server instance
				if srvInsCfg.NeedCheckHealth {
					taskQ[srvInsNode.Key] = newHealthJob(srvInsCfg.HealthCheckURL, srvInsNode.Key, HealthCheckTTL)
				}
			}
		}
	}
	return nil
}

func clusterWatchCallback(op etcdutils.OpCode, key, v string) {
	// logger.Logger.Infof("op: %d, key: %s", op, key)
	if !isInstanceKey(key) {
		return
	}
	
	switch op {
	case etcdutils.SetOp:
		instance := new(models.ServerInstance)
		if err := etcdutils.Decode(v, instance); err != nil {
			logger.Logger.Errorf("etcdutils.Decode(v,instance) failed: err %v, v=[%s]", err, v)
			return
		}
	
		taskQMutex.RLock()
		job, ok := taskQ[key]
		taskQMutex.RUnlock()

		// existed and addr has no changed, skip this
		if ok {
			if job.TargetURL == instance.HealthCheckURL {
				return
			}
			if !instance.NeedCheckHealth {
				taskQMutex.Lock()
				delete(taskQ, key)
				taskQMutex.Unlock()
				return
			}
		}
		// else set the key with new value
		taskQMutex.Lock()
		taskQ[key] = newHealthJob(instance.HealthCheckURL, key, HealthCheckTTL)
		taskQMutex.Unlock()
	case etcdutils.DeleteOp:
		taskQMutex.Lock()
		delete(taskQ, key)
		taskQMutex.Unlock()
	default:
		return
	}
}

func healthChecking(store *etcdutils.EtcdStore, pool *chanCheckerPool) {
	chanCheckResult := make(chan checkResult, 100)

	go func() {
		for {
			select {
			case cr := <-chanCheckResult:
				v, err := store.Get(cr.Key)
				if err != nil {
					logger.Logger.Errorf("healthChecking() kapi.Get(context.Background(), cr.Key, nil) got err: %v", err)
					continue
				}

				ins := new(models.ServerInstance)
				if err := etcdutils.Decode(v, ins); err != nil {
					logger.Logger.Errorf("healthChecking() kapi.Get(context.Background(), cr.Key, nil) got err: %v", err)
					continue
				}
				// TODO: ignore no need change
				// no change do not set again
				// if ins.IsAlive == cr.IsAlive {
				// 	return
				// }
				ins.IsAlive = cr.IsAlive
				data, _ := etcdutils.Encode(ins)
				if err := store.Set(cr.Key, string(data), -1); err != nil {
					logger.Logger.Errorf("healthChecking() kapi.Set(context.Background(), cr.Key, string(data), nil) got err: %v", err)
					continue
				}
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	for {
		taskQMutex.RLock()
		for _, job := range taskQ {
			checker, err := pool.Get()
			if err != nil {
				logger.Logger.Errorf("healthChecking() pool.Get() got err: %v", err)
				continue
			}
			go checker.Check(job, chanCheckResult)
		}
		taskQMutex.RUnlock()
		time.Sleep(100 * time.Millisecond)
	}
}
