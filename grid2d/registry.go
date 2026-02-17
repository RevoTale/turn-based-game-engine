package grid2d

// LayerRegistry provides multi-grid typed layer management.
// It is a convenience runtime wrapper over MultiLayerSpace.
type LayerRegistry[G comparable, K comparable, T comparable] struct {
	spaces *MultiLayerSpace[G, K, T]
}

func NewLayerRegistry[G comparable, K comparable, T comparable](gridSet *GridSet[G]) (*LayerRegistry[G, K, T], error) {
	if gridSet == nil {
		return nil, ErrNilGridSet
	}
	spaces, err := NewMultiLayerSpace[G, K, T](gridSet)
	if err != nil {
		return nil, err
	}
	return &LayerRegistry[G, K, T]{spaces: spaces}, nil
}

// Ensure returns an existing layer or creates it lazily if missing.
func (r *LayerRegistry[G, K, T]) Ensure(gridID G, layerKey K) (*SparseLayer[T], error) {
	space, err := r.ensureSpace(gridID)
	if err != nil {
		return nil, err
	}
	if layer, ok := space.Get(layerKey); ok {
		return layer, nil
	}
	return space.Create(layerKey)
}

// Get returns a layer if it exists.
// It never creates layer space or layer entries.
func (r *LayerRegistry[G, K, T]) Get(gridID G, layerKey K) (*SparseLayer[T], bool, error) {
	space, ok, err := r.existingSpace(gridID)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	layer, ok := space.Get(layerKey)
	return layer, ok, nil
}

func (r *LayerRegistry[G, K, T]) Set(gridID G, layerKey K, pos Position, value T) error {
	layer, err := r.existingLayer(gridID, layerKey)
	if err != nil {
		return err
	}
	return layer.Set(pos, value)
}

func (r *LayerRegistry[G, K, T]) DeleteValue(gridID G, layerKey K, pos Position) error {
	layer, err := r.existingLayer(gridID, layerKey)
	if err != nil {
		return err
	}
	return layer.Delete(pos)
}

// EnsureSet lazily creates a layer and writes a value in a single call.
func (r *LayerRegistry[G, K, T]) EnsureSet(gridID G, layerKey K, pos Position, value T) error {
	layer, err := r.Ensure(gridID, layerKey)
	if err != nil {
		return err
	}
	return layer.Set(pos, value)
}

func (r *LayerRegistry[G, K, T]) DeleteLayer(gridID G, layerKey K) (bool, error) {
	space, ok, err := r.existingSpace(gridID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return space.Delete(layerKey), nil
}

func (r *LayerRegistry[G, K, T]) DeleteGrid(gridID G) bool {
	if r == nil || r.spaces == nil {
		return false
	}
	return r.spaces.DeleteGrid(gridID)
}

func (r *LayerRegistry[G, K, T]) LayerSpaceCount() int {
	if r == nil || r.spaces == nil {
		return 0
	}
	return r.spaces.Count()
}

func (r *LayerRegistry[G, K, T]) ensureSpace(gridID G) (*LayerSpace[K, T], error) {
	if r == nil || r.spaces == nil {
		return nil, ErrNilLayerRegistry
	}
	return r.spaces.Space(gridID)
}

func (r *LayerRegistry[G, K, T]) existingSpace(gridID G) (*LayerSpace[K, T], bool, error) {
	if r == nil || r.spaces == nil {
		return nil, false, ErrNilLayerRegistry
	}
	return r.spaces.SpaceIfExists(gridID)
}

func (r *LayerRegistry[G, K, T]) existingLayer(gridID G, layerKey K) (*SparseLayer[T], error) {
	space, ok, err := r.existingSpace(gridID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrLayerNotFound
	}
	layer, ok := space.Get(layerKey)
	if !ok {
		return nil, ErrLayerNotFound
	}
	return layer, nil
}
