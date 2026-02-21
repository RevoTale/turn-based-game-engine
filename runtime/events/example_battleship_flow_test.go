package events_test

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/runtime/events"
)

func ExampleExecuteCommand_battleshipFlow() {
	type fireCommand struct {
		Coordinate string
		Hit        bool
	}
	type hitEvent struct {
		Coordinate string
	}
	type missEvent struct {
		Coordinate string
	}

	builder := events.NewBuilder()

	hit, err := events.RegisterEvent(builder, func(_ *events.Context, payload hitEvent) error {
		fmt.Printf("hit %s\n", payload.Coordinate)
		return nil
	}, events.Hooks[hitEvent]{})
	if err != nil {
		return
	}

	miss, err := events.RegisterEvent(builder, func(_ *events.Context, payload missEvent) error {
		fmt.Printf("miss %s\n", payload.Coordinate)
		return nil
	}, events.Hooks[missEvent]{})
	if err != nil {
		return
	}

	fire, err := events.RegisterCommand(builder, func(_ *events.Context, payload fireCommand) error {
		fmt.Printf("fire %s\n", payload.Coordinate)
		return nil
	}, events.Hooks[fireCommand]{
		After: func(ctx *events.Context, payload fireCommand) error {
			if payload.Hit {
				ctx.Emit(events.Next(hit, hitEvent{Coordinate: payload.Coordinate}))
				return nil
			}
			ctx.Emit(events.Next(miss, missEvent{Coordinate: payload.Coordinate}))
			return nil
		},
	})
	if err != nil {
		return
	}

	runtime, err := builder.Build()
	if err != nil {
		return
	}
	_ = events.ExecuteCommand(runtime, fire, fireCommand{Coordinate: "B4", Hit: true})
	_ = events.ExecuteCommand(runtime, fire, fireCommand{Coordinate: "A1", Hit: false})

	// Output:
	// fire B4
	// hit B4
	// fire A1
	// miss A1
}
