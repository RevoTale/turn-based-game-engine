package events

import "fmt"

type queuedEvent struct {
	id      EventID
	payload any
}

// Emission is one follow-up event queued by a handler or hook.
type Emission struct {
	id      EventID
	payload any
}

// Next creates one typed follow-up emission.
func Next[T any](event Event[T], payload T) Emission {
	return Emission{
		id:      event.id,
		payload: payload,
	}
}

type dispatcher struct {
	defs map[EventID]definition
}

func (d *dispatcher) dispatch(eventID EventID, payload any) error {
	if d == nil {
		return ErrNilRuntime
	}
	if _, ok := d.defs[eventID]; !ok {
		return fmt.Errorf("%w: id=%d", ErrUnknownEvent, eventID)
	}

	queue := []queuedEvent{{id: eventID, payload: payload}}
	ctx := &Context{
		current: eventID,
		queue:   &queue,
		defs:    d.defs,
	}

	for i := 0; i < len(queue); i++ {
		current := queue[i]
		ctx.current = current.id

		def, ok := d.defs[current.id]
		if !ok {
			return fmt.Errorf("%w: id=%d", ErrUnknownEvent, current.id)
		}
		if err := def.run(ctx, current.payload); err != nil {
			return fmt.Errorf("event %d failed: %w", current.id, err)
		}
	}

	return nil
}

// Runtime executes previously registered commands and emitted events.
type Runtime struct {
	dispatcher dispatcher
}

func (r *Runtime) execute(commandID EventID, payload any) error {
	if r == nil {
		return ErrNilRuntime
	}
	return r.dispatcher.dispatch(commandID, payload)
}

// Execute runs one root command and drains all emitted child events before it returns.
func Execute[T any](r *Runtime, command Command[T], payload T) error {
	if r == nil {
		return ErrNilRuntime
	}
	if _, ok := r.dispatcher.defs[command.id]; !ok {
		return fmt.Errorf("%w: id=%d", ErrUnknownCommand, command.id)
	}
	return r.execute(command.id, payload)
}

// Context carries dispatch state while handlers and hooks execute.
type Context struct {
	current EventID
	queue   *[]queuedEvent
	defs    map[EventID]definition
	emitErr error
}

// CurrentEventID returns the event id currently being processed.
func (c *Context) CurrentEventID() EventID {
	if c == nil {
		return 0
	}
	return c.current
}

// Emit appends one follow-up emission to the current event stage.
func (c *Context) Emit(emission Emission) {
	if c == nil || c.queue == nil {
		return
	}
	if c.emitErr != nil {
		return
	}
	if _, ok := c.defs[emission.id]; !ok {
		c.emitErr = fmt.Errorf("%w: id=%d", ErrUnknownEvent, emission.id)
		return
	}
	*c.queue = append(*c.queue, queuedEvent{
		id:      emission.id,
		payload: emission.payload,
	})
}

func (c *Context) beginEmissionScope() int {
	if c == nil || c.queue == nil {
		return 0
	}
	return len(*c.queue)
}

func (c *Context) rollbackEmissionScope(mark int) {
	if c == nil || c.queue == nil {
		return
	}
	if mark < 0 || mark > len(*c.queue) {
		return
	}
	*c.queue = (*c.queue)[:mark]
	c.emitErr = nil
}

func (c *Context) consumeEmitError() error {
	if c == nil {
		return nil
	}
	err := c.emitErr
	c.emitErr = nil
	return err
}
