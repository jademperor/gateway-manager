package healthchecking

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	errClosed = errors.New("pool is closed")
)

// Factory the generator to creat HealthChecker
type Factory func() (*HealthChecker, error)

func defaultChekerFactory() (*HealthChecker, error) {
	return &HealthChecker{
		client: &http.Client{Timeout: 5 * time.Second},
	}, nil
}

type chanCheckerPool struct {
	mutex    sync.RWMutex
	checkers chan *HealthChecker
	factory  Factory
}

// newCheckerPool ...
func newCheckerPool(initialCap, maxCap int, factory Factory) (*chanCheckerPool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	pool := &chanCheckerPool{
		checkers: make(chan *HealthChecker, maxCap),
		factory:  factory,
	}

	// create initial connections, if something goes wrong,
	// just close the pool error out.
	for i := 0; i < initialCap; i++ {
		checker, err := factory()
		if err != nil {
			checker.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		pool.checkers <- checker
	}

	return pool, nil
}

// getConnsAndFactory ...
func (p *chanCheckerPool) getConnsAndFactory() (chan *HealthChecker, Factory) {
	p.mutex.RLock()
	checkers, factory := p.checkers, p.factory
	p.mutex.RUnlock()
	return checkers, factory
}

// Close close the pool
func (p *chanCheckerPool) Close() {
	p.mutex.Lock()
	checkers := p.checkers
	p.checkers = nil
	p.factory = nil
	p.mutex.Unlock()

	if checkers == nil {
		return
	}

	close(checkers)
	for checker := range checkers {
		checker.Close()
	}
}

// Get
func (p *chanCheckerPool) Get() (*HealthChecker, error) {
	checkers, factory := p.getConnsAndFactory()
	if checkers == nil {
		return nil, errClosed
	}

	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case checker := <-checkers:
		return checker, nil
	default:
		checker, err := factory()
		if err != nil {
			return nil, err
		}
		return checker, nil
	}
}

// Put ...
func (p *chanCheckerPool) Put(checker *HealthChecker) error {
	if checker == nil {
		return errors.New("checker is nil. rejecting")
	}

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.checkers == nil {
		// pool is closed, close passed connection
		checker.Close()
		return nil
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case p.checkers <- checker:
		return nil
	default:
		// pool is full, close passed connection
		checker.Close()
		return nil
	}
}

// Len ...
func (p *chanCheckerPool) Len() int {
	checkers, _ := p.getConnsAndFactory()
	return len(checkers)
}
