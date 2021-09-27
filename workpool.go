package workpool

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidWorkPoolCount = errors.New("WorkPool MaxWorkersCount Invalid")
	ErrWorkPoolClosed       = errors.New("WorkPool Closed")
)

const (
	RUNNING = 1
	STOPED  = 0
)

type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

// WorkPool
type WorkPool struct {
	MaxWorkersCount int32
	runningWorkers  int32
	status          int
	chTask          chan *Task
	PanicHandler    func(interface{})
	lock            sync.Mutex
}

// NewWorkPool init WorkPool
func NewWorkPool(count int32) (*WorkPool, error) {
	if count <= 0 {
		return nil, ErrInvalidWorkPoolCount
	}
	p := &WorkPool{
		MaxWorkersCount: count,
		status:          RUNNING,
		chTask:          make(chan *Task, count),
	}

	return p, nil
}

func (p *WorkPool) checkWorker() {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.runningWorkers == 0 && len(p.chTask) > 0 {
		p.run()
	}
}

func (p *WorkPool) GetMaxWorkersCount() int32 {
	return p.MaxWorkersCount
}

func (p *WorkPool) GetRunningWorkers() int32 {
	return atomic.LoadInt32(&p.runningWorkers)
}

func (p *WorkPool) incRunning() {
	atomic.AddInt32(&p.runningWorkers, 1)
}

func (p *WorkPool) decRunning() {
	atomic.AddInt32(&p.runningWorkers, -1)
}

func (p *WorkPool) Put(task *Task) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.status == STOPED {
		return ErrWorkPoolClosed
	}

	// run worker
	if p.GetRunningWorkers() < p.GetMaxWorkersCount() {
		p.run()
	}

	// send task
	if p.status == RUNNING {
		p.chTask <- task
	}

	return nil
}

func (p *WorkPool) run() {
	p.incRunning()

	go func() {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					log.Printf("WorkerPool panic: %s\n", r)
				}
			}
			p.checkWorker() // check worker avoid no worker running
		}()

		for {
			select {
			case task, ok := <-p.chTask:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

func (p *WorkPool) setStatus(status int) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.status == status {
		return false
	}

	p.status = status

	return true
}

func (p *WorkPool) Close() {

	if !p.setStatus(STOPED) { // stop put task
		return
	}

	for len(p.chTask) > 0 { // wait all task be consumed
		time.Sleep(1e6) // reduce CPU load
	}

	close(p.chTask)
}
