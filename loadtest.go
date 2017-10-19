package main
/* Simple program that receives a single url, number of jobs and a max number of workers and will test performance with a concurrent approach */

import (
	"fmt"
	"time"
	"flag"
	"sync"
	"net/http"
)

var waitgroup sync.WaitGroup

func load_test(work chan string, result chan uint64) {
	waitgroup.Add(1)
	defer waitgroup.Done()
	for url := range work {
		t1 := time.Now()
		_, err := http.Get(url)
		if err != nil{
			panic(err)
		}
		t2 := time.Now()
		result <- uint64(t2.Sub(t1).Nanoseconds())
	}
}

func parse_flags() (string, uint64, uint64) {
	/* Parse input flags to parametrize execution */
	url := flag.String("url", "https://google.com", "This is the url that will be tested")
	jobs := flag.Uint64("jobs", 10, "Jobs assigned to the url")
	workers := flag.Uint64("workers", 2, "Concurrent workers")
	flag.Parse()
	return *url, *jobs, *workers
}

func main(){
	url, jobs, workers := parse_flags()
	work := make(chan string, jobs)
	result := make(chan uint64, jobs)
	
	var i uint64
	
	// Launch workers
	for i = 0; i < workers; i++ {
		go load_test(work, result)
	}
	
	// Start work
	for i = 0; i < jobs; i++ {
		work <- url
	}
	
	// Get results
	var accum uint64 = 0
	for i = 0; i < jobs; i++ {
		accum += <- result
	}
	
	// Close channels
	close(work)
	close(result)
	
	waitgroup.Wait()
	
	fmt.Printf("Average time for %s (%d workers for %d jobs): %f\n", url, workers, jobs, (float64(accum) / float64(jobs)) / float64(1e9))
}
