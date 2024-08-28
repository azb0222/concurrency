package pipelines

import "context"

// TODO: Error handling
// TODO: should all of these be public
type Pipeline[T any] struct {
	Generator IGenerator[T]
	Stages    map[string]*Stage[T]
}

// TODO: how to ensure type saftey in the passed in functions
func NewPipeline[T any](ctx context.Context, generator IGenerator[T], stageFns map[string]func(value T, args ...any) T) *Pipeline[T] {
	stages := map[string]*Stage[T]{}
	for k, fn := range stageFns {
		stages[k] = newStage(ctx, fn)
	}

	return &Pipeline[T]{
		Generator: generator,
		Stages:    stages,
	}
}
