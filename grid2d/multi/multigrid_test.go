package multi

import (
	"fmt"
	"testing"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_CreateDifferentGridSizes(t *testing.T) {
	t.Parallel()

	const (
		smallGridID uint8 = 1
		wideGridID  uint8 = 2
	)

	registry := NewRegistry[uint8, uint8, string]()
	small, err := registry.CreateGrid(smallGridID, 3, 2)
	require.NoError(t, err)
	wide, err := registry.CreateGrid(wideGridID, 12, 4)
	require.NoError(t, err)

	assert.Equal(t, 3, small.Width())
	assert.Equal(t, 2, small.Height())
	assert.Equal(t, 12, wide.Width())
	assert.Equal(t, 4, wide.Height())
	assert.Equal(t, 2, registry.GridCount())

	_, err = registry.CreateGrid(smallGridID, 10, 10)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridExists)
	assert.Equal(t, grid2d.CodeGridExists, grid2d.CodeOf(err))
}

func TestRegistry_UsesGridSpecificBounds(t *testing.T) {
	t.Parallel()

	const (
		gridAID  uint8 = 1
		gridBID  uint8 = 2
		stateKey uint8 = 1
	)

	registry := NewRegistry[uint8, uint8, string]()
	_, err := registry.CreateGrid(gridAID, 2, 2)
	require.NoError(t, err)
	_, err = registry.CreateGrid(gridBID, 5, 3)
	require.NoError(t, err)

	spaceA, err := registry.Space(gridAID)
	require.NoError(t, err)
	layerA, err := spaceA.Create(stateKey)
	require.NoError(t, err)

	spaceB, err := registry.Space(gridBID)
	require.NoError(t, err)
	layerB, err := spaceB.Create(stateKey)
	require.NoError(t, err)

	err = layerA.Set(grid2d.Position{X: 1, Y: 1}, "A")
	require.NoError(t, err)

	err = layerB.Set(grid2d.Position{X: 4, Y: 2}, "B")
	require.NoError(t, err)

	err = layerA.Set(grid2d.Position{X: 2, Y: 1}, "bad")
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrOutOfBounds)
}

func TestRegistry_SpaceIfExists_DoesNotCreate(t *testing.T) {
	t.Parallel()

	const (
		gridAID       uint8 = 1
		missingGridID uint8 = 99
	)

	registry := NewRegistry[uint8, uint8, string]()
	_, err := registry.CreateGrid(gridAID, 3, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	space, ok, err := registry.SpaceIfExists(gridAID)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, space)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	created, err := registry.Space(gridAID)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	space, ok, err = registry.SpaceIfExists(gridAID)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Same(t, created, space)

	space, ok, err = registry.SpaceIfExists(missingGridID)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, space)
}

func TestTurnbasedIntegration_Registry(t *testing.T) {
	t.Parallel()

	type playerID int
	type gridID uint8
	type layerKey uint8

	type action struct {
		Grid  gridID
		Pos   grid2d.Position
		State string
	}

	const (
		gridAlpha gridID   = 1
		gridBeta  gridID   = 2
		layerMain layerKey = 1
	)

	registry := NewRegistry[gridID, layerKey, string]()
	_, err := registry.CreateGrid(gridAlpha, 3, 2)
	require.NoError(t, err)
	_, err = registry.CreateGrid(gridBeta, 8, 4)
	require.NoError(t, err)

	alphaSpace, err := registry.Space(gridAlpha)
	require.NoError(t, err)
	alphaLayer, err := alphaSpace.Ensure(layerMain)
	require.NoError(t, err)
	betaSpace, err := registry.Space(gridBeta)
	require.NoError(t, err)
	betaLayer, err := betaSpace.Ensure(layerMain)
	require.NoError(t, err)

	layers := map[gridID]*grid2d.SparseLayer[string]{
		gridAlpha: alphaLayer,
		gridBeta:  betaLayer,
	}

	engine, err := turnbased.New[playerID, action]([]playerID{1, 2}, 1)
	require.NoError(t, err)

	resolve := func(actor playerID, act action) (turnbased.ActionOutcome[playerID], error) {
		layer, ok := layers[act.Grid]
		if !ok {
			return turnbased.ActionOutcome[playerID]{}, fmt.Errorf("%w: %v", grid2d.ErrGridNotFound, act.Grid)
		}
		if err := layer.Set(act.Pos, act.State); err != nil {
			return turnbased.ActionOutcome[playerID]{}, err
		}
		return turnbased.PassTurn[playerID](), nil
	}

	_, err = engine.Act(1, action{Grid: gridAlpha, Pos: grid2d.Position{X: 2, Y: 1}, State: "hit"}, resolve)
	require.NoError(t, err)
	assert.Equal(t, playerID(2), engine.CurrentPlayer())

	_, err = engine.Act(2, action{Grid: gridBeta, Pos: grid2d.Position{X: 7, Y: 3}, State: "burn"}, resolve)
	require.NoError(t, err)
	assert.Equal(t, playerID(1), engine.CurrentPlayer())

	value, ok, err := alphaLayer.Get(grid2d.Position{X: 2, Y: 1})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "hit", value)

	value, ok, err = betaLayer.Get(grid2d.Position{X: 7, Y: 3})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "burn", value)

	_, err = engine.Act(1, action{Grid: gridAlpha, Pos: grid2d.Position{X: 2, Y: 2}, State: "invalid"}, resolve)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrOutOfBounds)
	assert.Equal(t, playerID(1), engine.CurrentPlayer())
}
