package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

/*
Pool is a fixed-size worker pool built on three pieces working
together — not three separate features:

  - CHANNEL:      jobs, below. The queue itself — safe for many
                  goroutines to read from at once with no manual
                  locking.
  - CONCURRENCY:  the worker() loop — each goroutine independently
                  waits, receives, and processes, structured to
                  handle many jobs without coordinating in lockstep.
  - PARALLELISM:  Start(n) launching n goroutines — on a multi-core
                  machine, several of these literally run at the
                  same instant, not just interleaved. DOWNLOAD_WORKERS
                  in .env is what sets n.
*/
type Pool struct {
	jobs    chan Job // <-- the channel
	process ProcessFunc
	wg      sync.WaitGroup

	mu      sync.Mutex
	cancels map[string]context.CancelFunc // jobID -> cancel, for future /cancel endpoint
}

func NewPool(bufferSize int, process ProcessFunc) *Pool {
	return &Pool{
		jobs:    make(chan Job, bufferSize),
		process: process,
		cancels: make(map[string]context.CancelFunc),
	}
}

/*
Start launches n worker goroutines. <-- PARALLELISM happens here.
Each one runs worker() independently; the Go runtime schedules them
across available CPU cores.
*/
func (p *Pool) Start(n int) {
	for i := 0; i < n; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	log.Printf("worker pool started with %d workers", n)
}

/*
worker is the CONCURRENCY piece — an independent loop pulling from
the shared channel. `for job := range p.jobs` blocks until a job
arrives or the channel is closed (Shutdown), never busy-waits.
*/
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for job := range p.jobs {
		start := time.Now()
		log.Printf("worker %d: starting job %s", id, job.ID)

		if err := p.process(job); err != nil {
			log.Printf("worker %d: job %s failed: %v", id, job.ID, err)
		} else {
			log.Printf("worker %d: job %s completed in %s", id, job.ID, time.Since(start))
		}

		p.mu.Lock()
		delete(p.cancels, job.ID)
		p.mu.Unlock()

		job.Cancel() // release context resources regardless of outcome
	}
}

/*
Submit builds the job's context (30-minute timeout, workme §7) and
pushes it onto the channel. Non-blocking: if the channel is full,
returns false immediately instead of blocking the HTTP request that
called this — that's what service/download.go's ErrQueueFull expects.
*/
func (p *Pool) Submit(url, jobID string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	job := Job{ID: jobID, URL: url, Ctx: ctx, Cancel: cancel}

	select {
	case p.jobs <- job: // <-- sending into the channel
		p.mu.Lock()
		p.cancels[jobID] = cancel
		p.mu.Unlock()
		return true
	default:
		cancel()
		return false
	}
}

/*
Cancel looks up a running job's cancel func and calls it — this is
what POST /api/downloads/:id/cancel will call once that endpoint is
built. Returns false if the job isn't currently tracked (already
finished, or never existed).
*/
func (p *Pool) Cancel(jobID string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	cancel, ok := p.cancels[jobID]
	if !ok {
		return false
	}
	cancel()
	delete(p.cancels, jobID)
	return true
}

/*
Shutdown closes the channel (no more jobs accepted) and blocks until
every in-flight worker finishes its current job — called from
main.go's graceful-shutdown path so a deploy doesn't kill a download
mid-write.
*/
func (p *Pool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
	log.Println("worker pool shut down")
}
