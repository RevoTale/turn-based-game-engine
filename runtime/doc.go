// Package engruntime provides optional runtime primitives for integration
// layers.
//
// It currently offers typed topic-based event streams with explicit
// subscription ownership and publish delivery controls, plus keyed delay
// policy helpers for randomized action pacing.
//
// Core engine logic can remain deterministic while runtime concerns such as
// fan-out and subscriber lifecycle are handled here.
package engruntime
