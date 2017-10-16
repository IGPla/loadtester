package main
/* Simple program that receives a single url, a workload and a max number of workers and will test performance with a concurrent approach */

import "fmt"
import "time"
import "flag"
import "net/http"

func load_test(url string, ready chan bool, result chan uint64, closed chan bool) {
	for true{
		_, more := <- ready
		if(!more){
			break
		}
		t1 := time.Now()
		_, err := http.Get(url)
		if err != nil{
			panic(err)
		}
		t2 := time.Now()
		result <- uint64(t2.Sub(t1).Nanoseconds())
	}
	closed <- true
}

func parse_flags() (string, uint64, uint64) {
	/* Parse input flags to parametrize execution */
	url := flag.String("url", "test", "This is the url that will be tested")
	workload := flag.Uint64("workload", 10, "Workload for the url")
	workers := flag.Uint64("workers", 2, "Concurrent workers")
	flag.Parse()
	return *url, *workload, *workers
}

func main(){
	url, workload, workers := parse_flags()
	ready := make(chan bool, workload)
	result := make(chan uint64, workload)
	closed := make(chan bool, workers)

	var i uint64
	// Launch workers
	for i = 0; i < workers; i++ {
		go load_test(url, ready, result, closed)
	}
	// Start work
	for i = 0; i < workload; i++ {
		ready <- true
	}
	// Get results
	var accum uint64 = 0
	for i = 0; i < workload; i++ {
		accum += <- result
	}
	// Close channels
	close(ready)
	close(result)
	
	for i = 0; i < workers; i++{
		<- closed
	}
	fmt.Printf("Average time for %s (%d workers with %d workload): %f\n", url, workers, workload, (float64(accum) / float64(workload)) / float64(1e9))
}
