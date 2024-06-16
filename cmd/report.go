package cmd

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Report struct {
	executionStartTime  time.Time
	executionEndTime    time.Time
	TotalRequests       atomic.Int32
	TotalUndefinedError atomic.Int32
	StatusCounts        map[int]int64
	mu                  sync.Mutex
}

func (r *Report) StartExecution() {
	r.executionStartTime = time.Now()
	r.StatusCounts = make(map[int]int64)
}

func (r *Report) EndExecution() {
	r.executionEndTime = time.Now()
}

func (r *Report) Show() {
	log.Println("Report:")
	log.Printf("Total requests: %d", r.TotalRequests.Load())
	for statusCode, count := range r.StatusCounts {
		log.Printf("Total Status Code: %d, Count: %d\n", statusCode, count)
	}
	log.Printf("Total undefined errors: %d", r.TotalUndefinedError.Load())
	log.Printf("Execution time: %d milliseconds", r.executionEndTime.Sub(r.executionStartTime).Milliseconds())
}

func (r *Report) IncrementStatusCount(statusCode int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.StatusCounts[statusCode]++
}
