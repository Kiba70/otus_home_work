package hw06pipelineexecution

import "time"

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
			chClosed := false
			toClose := func() {
				if !chClosed {
					close(chOut)
					chClosed = true
				}
			}
			for {
				select {
				// Решение не однозначное - оно позволяет максимально быстро (одновременно)
				// закрыть все stage, но в период между закрытием канала done и закрытием
				// канала in будет высокая нагрузка на ЦПУ в связи постоянным срабатыванием
				// пункта "case <-done:".
				// Для уменьшения этого эффекта сделал паузу в данном пункте.
				// Как альтернатива данному решению - можно объединить все stage из pipeline_test.go
				// в одну цепочку и при закрытии канала done закрыть входной канал цепечки и
				// канал Out, при этом оставаясь на вычитке из цепочки для корректного завершения
				// всех горутин.
				case <-done:
					toClose()                         // Закрываем выходной канал, но не выходим (для вычитки из stage)
					time.Sleep(time.Millisecond * 10) // Не грузим CPU
				case v, ok := <-in:
					if !ok {
						toClose()
						return
					}
					if !chClosed {
						chOut <- v
					}
				}
			}
		}()
		return stages[0](chOut)
	}

	return ExecutePipeline(stages[0](in), done, stages[1:]...)
}
