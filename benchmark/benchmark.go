package benchmark

import (
	"fmt"
	"htcache/client"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"
	"strconv"
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
	Pipeline int
	Klen int
	Vlen int
	Result     *Result
	Elapsed    time.Duration
}

func NewBenchmark(protocol string, addr string, concurrent int, total int, operation string, pipeline int, klen int, vlen int) *Benchmark {
	return &Benchmark{
		Protocol:   protocol,
		Addr:       addr,
		Concurrent: concurrent,
		Total:      total,
		Pipeline:pipeline,
		Klen: klen,
		Vlen: vlen,
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
			if b.Pipeline > 1 {
				results <- b.PipelineRun()
			} else {
				results <- b.Run()
			}

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
	key := strings.Repeat(strconv.Itoa(rand.Intn(31) + int('a')), b.Klen)
	value := strings.Repeat(strconv.Itoa(rand.Intn(31) + int('A')), b.Vlen)

	var command *client.Command

	result := NewResult()
	cli := client.NewClient(b.Protocol, b.Addr)
	defer cli.Close()

	for i := 0; i < b.Total/b.Concurrent; i++ {
		start := time.Now()

		rtype := b.Operation
		if rtype == "mixed" {
			rtype = "set"
			if rand.Int()%2 == 0 {
				rtype = "get"
			}
		}
		switch rtype {
		case "set":
			command = client.NewCommand(rtype, fmt.Sprintf("%s_%d", key, i), []byte(fmt.Sprintf("%s_%d", value, i)))
			cli.Run(command)
		case "get":
			command = client.NewCommand(rtype, fmt.Sprintf("%s_%d", key, i), nil)
			cli.Run(command)
			if len(command.Value) == 0 {
				rtype = "miss"
			}
		}

		elapsed := time.Now().Sub(start)

		result.AddDuration(elapsed, rtype)
	}

	return result
}

func (b *Benchmark) PipelineRun() *Result {
	key := strings.Repeat(strconv.Itoa(rand.Intn(31) + int('a')), b.Klen)
	value := strings.Repeat(strconv.Itoa(rand.Intn(31) + int('A')), b.Vlen)

	commands := make([]*client.Command, 0)
	result := NewResult()
	cli := client.NewClient(b.Protocol, b.Addr)
	defer cli.Close()

	start := time.Now()
	for i := 0; i < b.Total/b.Concurrent; i++ {
		rtype := b.Operation
		if rtype == "mixed" {
			rtype = "set"
			if rand.Int()%2 == 0 {
				rtype = "get"
			}
		}
		switch rtype {
		case "set":
			commands = append(commands, client.NewCommand(rtype, fmt.Sprintf("%s_%d", key, i), []byte(fmt.Sprintf("%s_%d", value, i))))
		case "get":
			commands = append(commands, client.NewCommand(rtype, fmt.Sprintf("%s_%d", key, i), nil))
		}

		if len(commands) >= b.Pipeline {
			cli.Pipeline(commands)
			elapsed := time.Duration(float64(time.Now().Sub(start)) / float64(b.Pipeline))
			for _, command := range commands {
				switch command.Name {
				case "set":
					result.AddDuration(elapsed, command.Name)
				case "get":
					if len(command.Value) == 0 {
						result.AddDuration(elapsed, "miss")
					} else {
						result.AddDuration(elapsed, command.Name)
					}
				}
			}
			commands = make([]*client.Command, 0)
			start = time.Now()
		}

	}
	if len(commands) >= 0 {
		cli.Pipeline(commands)
		elapsed := time.Duration(float64(time.Now().Sub(start)) / float64(b.Pipeline))
		for _, command := range commands {
			switch command.Name {
			case "set":
				result.AddDuration(elapsed, command.Name)
			case "get":
				if len(command.Value) == 0 {
					result.AddDuration(elapsed, "miss")
				} else {
					result.AddDuration(elapsed, command.Name)
				}
			}
		}
		commands = make([]*client.Command, 0)
	}
	return result
}

func (b *Benchmark) Output(writer io.Writer) {
	fmt.Fprintf(writer, "elapsed: %s\n", b.Elapsed)
	fmt.Fprintf(writer, "pipline: %d\n", b.Pipeline)
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
