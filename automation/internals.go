package automation

import "sync"

type scopedGuard[K comparable] struct {
	mu     sync.Mutex
	active map[K]struct{}
}

func (g *scopedGuard[K]) tryAcquire(scope K) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.active == nil {
		g.active = make(map[K]struct{})
	}
	if _, exists := g.active[scope]; exists {
		return false
	}

	g.active[scope] = struct{}{}
	return true
}

func (g *scopedGuard[K]) release(scope K) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.active == nil {
		return
	}
	delete(g.active, scope)
}

func (g *scopedGuard[K]) isActive(scope K) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.active == nil {
		return false
	}
	_, ok := g.active[scope]
	return ok
}

type asyncTracker struct {
	wg sync.WaitGroup
}

func (t *asyncTracker) start() {
	t.wg.Add(1)
}

func (t *asyncTracker) done() {
	t.wg.Done()
}

func (t *asyncTracker) wait() {
	t.wg.Wait()
}
