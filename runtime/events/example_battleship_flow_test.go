package events_test

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/runtime/events"
)

func ExampleExecuteCommand_battleshipFlow() {
	type state struct{}
	type patch struct {
		coord string
		hit   bool
		trace []string
	}
	type input struct {
		Coordinate string
		Hit        bool
	}

	var hitEv events.Event[*state, patch]
	var missEv events.Event[*state, patch]

	resolveEv, _ := events.DefineEvent[*state, patch](func(ctx events.Context[*state, patch]) error {
		p := ctx.Patch()
		if p.hit {
			ctx.Emit(hitEv)
			return nil
		}
		ctx.Emit(missEv)
		return nil
	})

	hitEv, _ = events.DefineEvent[*state, patch](func(ctx events.Context[*state, patch]) error {
		p := ctx.Patch()
		p.trace = append(p.trace, "hit "+p.coord)
		return nil
	})

	missEv, _ = events.DefineEvent[*state, patch](func(ctx events.Context[*state, patch]) error {
		p := ctx.Patch()
		p.trace = append(p.trace, "miss "+p.coord)
		return nil
	})

	fireCmd, _ := events.DefineCommand[*state, patch, input](func(ctx events.Context[*state, patch], in input) error {
		p := ctx.Patch()
		p.coord = in.Coordinate
		p.hit = in.Hit
		p.trace = append(p.trace, "fire "+in.Coordinate)
		ctx.Emit(resolveEv)
		return nil
	})

	runtime := events.NewRuntime()
	s := &state{}

	first, _ := events.ExecuteCommand(runtime, s, fireCmd, input{Coordinate: "B4", Hit: true}, func() *patch {
		return &patch{trace: make([]string, 0, 4)}
	})
	second, _ := events.ExecuteCommand(runtime, s, fireCmd, input{Coordinate: "A1", Hit: false}, func() *patch {
		return &patch{trace: make([]string, 0, 4)}
	})

	for _, line := range first.trace {
		fmt.Println(line)
	}
	for _, line := range second.trace {
		fmt.Println(line)
	}

	// Output:
	// fire B4
	// hit B4
	// fire A1
	// miss A1
}
