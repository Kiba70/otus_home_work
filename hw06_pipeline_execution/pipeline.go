package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out { //nolint: gocognit
	chOut := make(Bi)
	sLen := len(stages)

	if sLen == 0 { // Stages нет
		go func() {
			defer close(chOut)
			for {
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case chOut <- in:
					}
				}
			}
		}()
		return chOut
	}

	if sLen == 1 {
		go func() {
			for {
				select {
				case <-done:
					close(chOut) // Закрываем выходной канал, но не выходим (для вычитки из stage)
					for {        // Вычитываем входной канал и ждём закрытия
						if _, ok := <-in; !ok {
							return
						}
					}
				case v, ok := <-in:
					if !ok {
						close(chOut)
						return
					}
					chOut <- v
				}
			}
		}()
		return stages[0](chOut)
	}

	return ExecutePipeline(stages[0](in), done, stages[1:]...)
}
