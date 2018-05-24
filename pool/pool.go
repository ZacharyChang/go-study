package pool

import (
	"log"
)

type Pool struct {
	Queue  chan func() error
	Number int
	Size   int

	result         chan error
	finishCallback func()
}

func (pool *Pool) Init(number int, size int) {
	pool.Queue = make(chan func() error, size)
	pool.Number = number
	pool.Size = size
	pool.result = make(chan error, size)
}

// 开启线程池
func (pool *Pool) Start() {
	log.Println("[POOL]: Starting...")
	// 开启goroutine
	for i := 0; i < pool.Number; i++ {
		go func() {
			for {
				task, ok := <-pool.Queue
				if !ok {
					break
				}

				err := task()
				pool.result <- err
			}
		}()
	}
	// 获取执行结果
	for j := 0; j < pool.Size; j++ {
		res, ok := <-pool.result
		if !ok {
			break
		}
		if res != nil {
			log.Println("[POOL]: Task fail： ", res)
		}
	}
	// 回调函数
	if pool.finishCallback != nil {
		pool.finishCallback()
	}
}

// 关闭线程池
func (pool *Pool) Stop() {
	log.Println("[POOL]: Stopping...")
	close(pool.Queue)
	close(pool.result)
}

// 添加任务
func (pool *Pool) AddTask(task func() error) {
	pool.Queue <- task
}

// 设置回调函数
func (pool *Pool) SetFinishCallback(callback func()) {
	pool.finishCallback = callback
}
