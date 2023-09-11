package worker

import (
	"time"
)

type Message[T any] struct {
	WorkerId int
	Duration time.Duration
	Result   T
}

type Task[T any] func() T

type IWorker[T any] interface {
	Work()
	Stop()
	AddTask(t Task[T])
}

func NewWorker[T any](numberOfWorkers int, messageChan chan Message[T]) IWorker[T] {
	return &Worker[T]{
		NumberOfWorkers: numberOfWorkers,
		messageChan:     messageChan,
		taskChan:        make(chan Task[T]),
		stopWorkerChan:  make(chan bool, numberOfWorkers),
	}
}

var _ IWorker[any] = (*Worker[any])(nil)

type Worker[T any] struct {
	NumberOfWorkers int
	taskChan        chan Task[T]
	messageChan     chan Message[T]
	stopWorkerChan  chan bool
}

func (w *Worker[T]) Work() {
	for i := 1; i <= w.NumberOfWorkers; i += 1 {
		go func(workerId int) {
			select {
			case t := <-w.taskChan:
				st := time.Now()
				r := t()
				w.messageChan <- Message[T]{WorkerId: workerId, Duration: time.Now().Sub(st), Result: r}
			case <-w.stopWorkerChan:
				break
			}
		}(i)
	}
}

func (w *Worker[T]) Stop() {
	for i := 0; i < w.NumberOfWorkers; i += 1 {
		w.stopWorkerChan <- true
	}
}

func (w *Worker[T]) AddTask(t Task[T]) {
	w.taskChan <- t
}
