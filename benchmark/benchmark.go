package benchmark

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"
)

type Statistics struct {
	Count int
	Time  time.Duration
}

type Result struct {
	Get     int
	Set     int
	Miss    int
	Total   int
	Buckets []Statistics
}

func NewResult() *Result {
	return &Result{
		Get:     0,
		Set:     0,
		Miss:    0,
		Buckets: []Statistics{},
	}
}

func (r *Result) AddDuration(elapsed time.Duration, rtype string) {
	index := int(elapsed / time.Millisecond)
	if index >= len(r.Buckets) {
		buckets := make([]Statistics, index+1)
		copy(buckets, r.Buckets)
		r.Buckets = buckets
	}
	bucket := r.Buckets[index]
	bucket.Count++
	bucket.Time += elapsed
	r.Buckets[index] = bucket

	r.Total++
	switch rtype {
	case "get":
		r.Get++
	case "set":
		r.Set++
	case "miss":
		r.Miss++
	}
}

func (r *Result) AddResult(rs *Result) {
	r.Get += rs.Get
	r.Set += rs.Set
	r.Miss += rs.Miss
	r.Total += rs.Total
	for index, bucket := range rs.Buckets {
		if index >= len(r.Buckets) {
			r.Buckets = append(r.Buckets, bucket)
		} else {
			tbucket := r.Buckets[index]
			tbucket.Count += bucket.Count
			tbucket.Time += bucket.Time
			r.Buckets[index] = tbucket
		}
	}
}

type Benchmark struct {
	Protocol   string
	Addr       string
	Concurrent int
	Total      int
	Operation  string
	Result     *Result
	Elapsed    time.Duration
}

func NewBenchmark(protocol string, addr string, concurrent int, total int, operation string) *Benchmark {
	return &Benchmark{
		Protocol:   protocol,
		Addr:       addr,
		Concurrent: concurrent,
		Total:      total,
		Operation:  operation,
		Result:     NewResult(),
	}
}

func (b *Benchmark) Execute() {
	start := time.Now()
	results := make(chan *Result, b.Concurrent)

	var wg sync.WaitGroup
	wg.Add(b.Concurrent)

	for i := 0; i < b.Concurrent; i++ {
		go func() {
			results <- b.Run()
			wg.Done()
		}()
	}

	wg.Wait()

	b.Elapsed = time.Now().Sub(start)
	for i := 0; i < b.Concurrent; i++ {
		b.Result.AddResult(<-results)
	}

}

func (b *Benchmark) Run() *Result {
	result := NewResult()
	for i := 0; i < b.Total / b.Concurrent; i++ {
		start := time.Now()
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		elapsed := time.Now().Sub(start)
		rtype := b.Operation
		result.AddDuration(elapsed, rtype)
	}

	return result
}

func (b *Benchmark) Output(writer io.Writer) {
	fmt.Fprintf(writer, "elapsed: %s\n", b.Elapsed)
	fmt.Fprintf(writer, "Set Count: %d\n", b.Result.Set)
	fmt.Fprintf(writer, "Get Count: %d\n", b.Result.Get)
	fmt.Fprintf(writer, "Miss Count: %d\n", b.Result.Miss)
	fmt.Fprintf(writer, "Total Count: %d\n", b.Result.Total)
	count := 0
	duration := time.Duration(0)
	for index, bucket := range b.Result.Buckets {
		if bucket.Count == 0 {
			continue
		}
		count += bucket.Count
		duration += bucket.Time
		fmt.Fprintf(writer, "%.2f%% requests < %d ms\n", float64(count)*100.0/float64(b.Result.Total), index+1)
	}

	fmt.Fprintf(writer, "qps: %.2f\n", float64(b.Result.Total)/float64(b.Elapsed.Seconds()))
	fmt.Fprintf(writer, "average request: %.2f ms\n", float64(duration.Seconds()*1000)/float64(count))
}

func init() {
	rand.Seed(time.Now().Unix())
}
