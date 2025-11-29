package graph

import (
	"sync"
)

type workerPool struct {
	workers   int
	taskQueue chan func()
	wg        sync.WaitGroup
}

func newWorkerPool(workers int, queueSize int) *workerPool {
	useQueueSize := queueSize
	if useQueueSize <= 0 {
		useQueueSize = 100
	}
	useWorkers := workers
	if useWorkers <= 0 {
		useWorkers = 4
	}

	pool := &workerPool{
		workers:   useWorkers,
		taskQueue: make(chan func(), useQueueSize),
	}
	pool.start()
	return pool
}

func (wp *workerPool) start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for task := range wp.taskQueue {
				task()
			}
		}()
	}
}

// Submit adds a task to the worker pool's task queue.
func (wp *workerPool) Submit(task func()) {
	wp.taskQueue <- task
}

// Shutdown gracefully shuts down the worker pool, waiting for all workers to finish.
func (wp *workerPool) Shutdown() {
	close(wp.taskQueue)
	wp.wg.Wait()
}
