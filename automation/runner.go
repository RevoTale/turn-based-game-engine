package automation

import "fmt"

// Config configures a generic multi-actor automation loop.
type Config[P any, R any] struct {
	// MaxIterations is the maximum number of outer loop passes.
	MaxIterations int

	// ForEachActor visits each actor once per iteration.
	//
	// Returning false from visit stops actor iteration for the current pass.
	ForEachActor func(visit func(actor P) bool)

	// ProcessActor processes one actor step.
	//
	// emit can be called zero or more times to publish produced results.
	// The returned values mean:
	// - acted: actor performed any state-changing action
	// - stop: stop the whole automation loop after this actor step
	// - err:  step failed; loop continues after OnError callback
	ProcessActor func(actor P, emit func(result R)) (acted bool, stop bool, err error)

	// OnEmit is called for each emitted result in emit order.
	OnEmit func(result R)

	// OnError is called when ProcessActor returns an error.
	OnError func(actor P, err error)
}

// Run executes a bounded automation loop and returns emitted results.
//
// The loop stops when:
// - one ProcessActor step requests stop
// - one full iteration performs no actions
//
// Run returns ErrMaxIterationsReached when all iterations were used and the
// loop still kept making progress.
func Run[P any, R any](config Config[P, R]) ([]R, error) {
	if config.MaxIterations <= 0 {
		return nil, ErrInvalidMaxIterations
	}
	if config.ForEachActor == nil {
		return nil, ErrNilForEachActor
	}
	if config.ProcessActor == nil {
		return nil, ErrNilProcessActor
	}

	results := make([]R, 0)
	emit := func(result R) {
		results = append(results, result)
		if config.OnEmit != nil {
			config.OnEmit(result)
		}
	}

	for i := 0; i < config.MaxIterations; i++ {
		actedInIteration := false
		stop := false

		config.ForEachActor(func(actor P) bool {
			acted, shouldStop, err := config.ProcessActor(actor, emit)
			if err != nil && config.OnError != nil {
				config.OnError(actor, err)
			}
			if acted {
				actedInIteration = true
			}
			if shouldStop {
				stop = true
				return false
			}
			return true
		})

		if stop || !actedInIteration {
			return results, nil
		}
	}

	return results, fmt.Errorf("%w: %d", ErrMaxIterationsReached, config.MaxIterations)
}
