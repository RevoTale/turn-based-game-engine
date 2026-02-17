package grid2d

import (
	"fmt"
	"testing"

	"github.com/RevoTale/turn-based-game-engine/turnbased"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGridSet_CreateDifferentSizes(t *testing.T) {
	t.Parallel()

	grids := NewGridSet[string]()
	small, err := grids.Create("small", 3, 2)
	require.NoError(t, err)
	wide, err := grids.Create("wide", 12, 4)
	require.NoError(t, err)

	assert.Equal(t, 3, small.Width())
	assert.Equal(t, 2, small.Height())
	assert.Equal(t, 12, wide.Width())
	assert.Equal(t, 4, wide.Height())
	assert.Equal(t, 2, grids.Count())

	_, err = grids.Create("small", 10, 10)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGridExists)
	assert.Equal(t, CodeGridExists, CodeOf(err))
}

func TestMultiLayerSpace_UsesGridSpecificBounds(t *testing.T) {
	t.Parallel()

	grids := NewGridSet[string]()
	_, err := grids.Create("a", 2, 2)
	require.NoError(t, err)
	_, err = grids.Create("b", 5, 3)
	require.NoError(t, err)

	spaces, err := NewMultiLayerSpace[string, string, string](grids)
	require.NoError(t, err)

	spaceA, err := spaces.Space("a")
	require.NoError(t, err)
	layerA, err := spaceA.Create("state")
	require.NoError(t, err)

	spaceB, err := spaces.Space("b")
	require.NoError(t, err)
	layerB, err := spaceB.Create("state")
	require.NoError(t, err)

	err = layerA.Set(Position{X: 1, Y: 1}, "A")
	require.NoError(t, err)

	err = layerB.Set(Position{X: 4, Y: 2}, "B")
	require.NoError(t, err)

	err = layerA.Set(Position{X: 2, Y: 1}, "bad")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)
}

func TestMultiLayerSpace_SpaceIfExists_DoesNotCreate(t *testing.T) {
	t.Parallel()

	grids := NewGridSet[string]()
	_, err := grids.Create("a", 3, 2)
	require.NoError(t, err)

	spaces, err := NewMultiLayerSpace[string, string, string](grids)
	require.NoError(t, err)
	assert.Equal(t, 0, spaces.Count())

	space, ok, err := spaces.SpaceIfExists("a")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, space)
	assert.Equal(t, 0, spaces.Count())

	created, err := spaces.Space("a")
	require.NoError(t, err)
	assert.Equal(t, 1, spaces.Count())

	space, ok, err = spaces.SpaceIfExists("a")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Same(t, created, space)

	space, ok, err = spaces.SpaceIfExists("missing")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, space)
}

func TestTurnbasedIntegration_MultiGrid(t *testing.T) {
	t.Parallel()

	type playerID int
	type gridID string
	type layerKey string

	type action struct {
		Grid  gridID
		Pos   Position
		State string
	}

	grids := NewGridSet[gridID]()
	_, err := grids.Create("alpha", 3, 2)
	require.NoError(t, err)
	_, err = grids.Create("beta", 8, 4)
	require.NoError(t, err)

	spaces, err := NewMultiLayerSpace[gridID, layerKey, string](grids)
	require.NoError(t, err)

	alphaSpace, err := spaces.Space("alpha")
	require.NoError(t, err)
	alphaLayer, err := alphaSpace.Create("states")
	require.NoError(t, err)

	betaSpace, err := spaces.Space("beta")
	require.NoError(t, err)
	betaLayer, err := betaSpace.Create("states")
	require.NoError(t, err)

	layers := map[gridID]*SparseLayer[string]{
		"alpha": alphaLayer,
		"beta":  betaLayer,
	}

	engine, err := turnbased.New[playerID, action]([]playerID{1, 2}, 1)
	require.NoError(t, err)

	resolve := func(actor playerID, act action) (turnbased.ActionOutcome[playerID], error) {
		layer, ok := layers[act.Grid]
		if !ok {
			return turnbased.ActionOutcome[playerID]{}, fmt.Errorf("%w: %v", ErrGridNotFound, act.Grid)
		}
		if err := layer.Set(act.Pos, act.State); err != nil {
			return turnbased.ActionOutcome[playerID]{}, err
		}
		return turnbased.PassTurn[playerID](), nil
	}

	_, err = engine.Act(1, action{Grid: "alpha", Pos: Position{X: 2, Y: 1}, State: "hit"}, resolve)
	require.NoError(t, err)
	assert.Equal(t, playerID(2), engine.CurrentPlayer())

	_, err = engine.Act(2, action{Grid: "beta", Pos: Position{X: 7, Y: 3}, State: "burn"}, resolve)
	require.NoError(t, err)
	assert.Equal(t, playerID(1), engine.CurrentPlayer())

	value, ok, err := alphaLayer.Get(Position{X: 2, Y: 1})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "hit", value)

	value, ok, err = betaLayer.Get(Position{X: 7, Y: 3})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "burn", value)

	_, err = engine.Act(1, action{Grid: "alpha", Pos: Position{X: 2, Y: 2}, State: "invalid"}, resolve)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)
	assert.Equal(t, playerID(1), engine.CurrentPlayer())
}
