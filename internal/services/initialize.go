package services

import (
	"github.com/jademperor/common/etcdutils"
	// "github.com/jademperor/common/pkg/utils"
)

var (
	store *etcdutils.EtcdStore
	// gLockCluster  *gLock
	// gLockAPIs     *gLock
	// gLockRoutings *gLock
)

// Init store
func Init(addrs []string) (err error) {
	store, err = etcdutils.NewEtcdStore(addrs)

	// gLockCluster = newGLock()
	// gLockAPIs = newGLock()
	// gLockRoutings = newGLock()

	// control UUID length
	// utils.SetUUIDBytesLen(8)

	return
}

// func newGLock() *gLock {
// 	return &gLock{
// 		changed:          false,
// 		lastModifiedTime: time.Now(),
// 		mutex:            &sync.RWMutex{},
// 	}
// }

// type gLock struct {
// 	changed          bool          // changed flag
// 	lastModifiedTime time.Time     // last modified time
// 	mutex            *sync.RWMutex // rw mutex
// }

// func (l *gLock) setChangedFlag(c bool) {
// 	l.mutex.Lock()
// 	defer l.mutex.Unlock()

// 	l.changed = c
// }

// func (l *gLock) hasChanged() bool {
// 	l.mutex.RLock()
// 	defer l.mutex.RUnlock()

// 	c := l.changed
// 	return c
// }
