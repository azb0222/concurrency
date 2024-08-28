package tests

import (
	"asritha.dev/concurrency/pkg/patterns/channels"
	"asritha.dev/concurrency/pkg/patterns/pipelines"
	"context"
	"fmt"
)

// TODO: am i using context everywhere correctly?
func PipelineStreamTester() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Stream Testing
	nums := []int{1, 2, 3, 4}
	//	expectedNums := []int{4, 6, 8, 10}
	streamGenerator := pipelines.NewStreamGenerator[int](nums...)

	stageFns := map[string]func(value int, args ...any) int{
		"addByOne": func(value int, args ...any) int {
			return value + 1
		},
		"mulByTwo": func(value int, args ...any) int {
			return value * 2
		},
	}
	p := pipelines.NewPipeline[int](ctx, streamGenerator, stageFns)

	addByOneStage := p.Stages["addByOne"].Fn
	//mulByTwoStage := p.Stages["mulByTwo"].FanOutFanIn(ctx).Fn
	mulByTwoStage := p.Stages["mulByTwo"].Fn

	channels.Tee(ctx, addByOneStage(streamGenerator.GetValues(ctx)))

	outputNums := mulByTwoStage(addByOneStage(streamGenerator.GetValues(ctx)))

	for num := range outputNums {
		fmt.Println(num)
	}
}
