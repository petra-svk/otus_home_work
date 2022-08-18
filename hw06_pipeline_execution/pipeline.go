package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline builds pipeline from stages.
// stages should not be nil. Otherwise ExecutePipeline panics.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := exitStage(in, done)
	for _, stFunc := range stages {
		out = stFunc(exitStage(out, done))
	}
	return out
}

func exitStage(in, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()

	return out
}
