package automation

import (
	"sync"
	"time"
)

// SessionConfig configures one scoped automation session.
type SessionConfig[K comparable, P any, R any] struct {
	// MaxIterations bounds the number of outer loop passes.
	MaxIterations int

	// ForEachActor visits all actors for the provided scope.
	ForEachActor func(scope K, visit func(actor P) bool)

	// ProcessActor executes one actor step for the provided scope.
	ProcessActor func(scope K, actor P, emit func(result R)) (acted bool, stop bool, err error)

	// DelayForActor optionally returns a delay before processing actor.
	DelayForActor func(scope K, actor P) time.Duration

	// OnEmit is called for each emitted result.
	OnEmit func(scope K, result R)

	// OnError is called when ProcessActor returns an error.
	OnError func(scope K, actor P, err error)
}

// SchedulerOptions configures scheduler runtime behavior.
type SchedulerOptions struct {
	// Sleep is used by DelayForActor. Defaults to time.Sleep.
	Sleep func(time.Duration)
}

// Runtime serializes automation execution per scope and exposes one
// synchronous session runner.
type Runtime[K comparable, P any, R any] interface {
	// Run executes one scoped automation session synchronously.
	Run(scope K, cfg SessionConfig[K, P, R]) ([]R, error)
	// IsActive reports whether the given scope is currently being processed.
	IsActive(scope K) bool
}

type scheduler[K comparable, P any, R any] struct {
	mu     sync.Mutex
	sleep  func(time.Duration)
	active map[K]struct{}
}

// NewScheduler creates a scheduler with optional runtime options.
func NewScheduler[K comparable, P any, R any](opts SchedulerOptions) Runtime[K, P, R] {
	sleep := opts.Sleep
	if sleep == nil {
		sleep = time.Sleep
	}
	return &scheduler[K, P, R]{
		sleep:  sleep,
		active: make(map[K]struct{}),
	}
}

// Run executes one scoped automation session synchronously.
//
// Run returns ErrScopeBusy when the same scope is already active.
func (s *scheduler[K, P, R]) Run(scope K, cfg SessionConfig[K, P, R]) ([]R, error) {
	if s == nil {
		return nil, ErrScopeBusy
	}

	s.mu.Lock()
	if _, busy := s.active[scope]; busy {
		s.mu.Unlock()
		return nil, ErrScopeBusy
	}
	s.active[scope] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.active, scope)
		s.mu.Unlock()
	}()

	return s.runSession(scope, cfg)
}

// IsActive reports whether the given scope is currently being processed.
func (s *scheduler[K, P, R]) IsActive(scope K) bool {
	if s == nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.active[scope]
	return ok
}

func (s *scheduler[K, P, R]) runSession(scope K, cfg SessionConfig[K, P, R]) ([]R, error) {
	if cfg.MaxIterations <= 0 {
		return nil, ErrInvalidMaxIterations
	}
	if cfg.ForEachActor == nil {
		return nil, ErrNilForEachActor
	}
	if cfg.ProcessActor == nil {
		return nil, ErrNilProcessActor
	}

	return Run(Config[P, R]{
		MaxIterations: cfg.MaxIterations,
		ForEachActor: func(visit func(actor P) bool) {
			cfg.ForEachActor(scope, visit)
		},
		ProcessActor: func(actor P, emit func(result R)) (bool, bool, error) {
			if cfg.DelayForActor != nil {
				if delay := cfg.DelayForActor(scope, actor); delay > 0 {
					s.sleep(delay)
				}
			}
			return cfg.ProcessActor(scope, actor, emit)
		},
		OnEmit: func(result R) {
			if cfg.OnEmit != nil {
				cfg.OnEmit(scope, result)
			}
		},
		OnError: func(actor P, err error) {
			if cfg.OnError != nil {
				cfg.OnError(scope, actor, err)
			}
		},
	})
}
