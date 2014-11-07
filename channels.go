package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

var debugFlag bool
var verboseFlag bool
var delayFlag int
var iterationsFlag int
var queuesizeFlag int
var workersFlag int
var blockingFlag bool
var numcpusFlag int
var workers int

func msg(msg string) {
	if verboseFlag {
		fmt.Printf("%s", msg)
	}
}

func msgD(msg string) {
	if debugFlag {
		fmt.Printf("%s", msg)
	}
}

func worker(queue chan int, quit chan bool) {
	workers++
	var count = 0
	var delayDur = time.Duration(delayFlag)
	for {
		select {
		case <-quit:
			msg(fmt.Sprintf("\tworker exiting (%d jobs processed)\n", count))
			workers--
			return
		case _ = <-queue:
			msgD("-")
			count++
			time.Sleep(time.Microsecond * delayDur)
		}
	}
}

func init() {
	flag.IntVar(&delayFlag, "d", 0, "microsecond delay after worker processes job before getting another one")
	flag.IntVar(&iterationsFlag, "i", 100, "number of jobs to run")
	flag.IntVar(&queuesizeFlag, "q", 10, "number of jobs to hold in the queue")
	flag.IntVar(&workersFlag, "w", 10, "number of workers threads")
	flag.BoolVar(&blockingFlag, "b", false, "enable blocking queue behavior (unbuffered channel)")
	flag.IntVar(&numcpusFlag, "p", 1, "number of logical CPUs to use (0 means use ALL)")
	flag.BoolVar(&verboseFlag, "v", false, "detailed output")
	flag.BoolVar(&debugFlag, "D", false, "debug output")
}

func main() {
	flag.Parse()

	// If numcpus set to zero or larger than system logical cpus
	if numcpusFlag == 0 || numcpusFlag > runtime.NumCPU() {
		numcpusFlag = runtime.NumCPU()
	}

	// Set channel (queue) size appropriately
	if blockingFlag {
		queuesizeFlag = 1
	}

	// Create job queue channel
	queue := make(chan int, queuesizeFlag)

	// Set the number of available CPUs
	oldprocs := runtime.GOMAXPROCS(numcpusFlag)

	// Print parameters
	msg(fmt.Sprintf("# workers: %d\n", workersFlag))
	msg(fmt.Sprintf("Worker delay: %d microsecond(s)\n", delayFlag))
	msg(fmt.Sprintf("Queue size: %d\n", queuesizeFlag))
	msg(fmt.Sprintf("Blocking: %v\n", blockingFlag))
	msg(fmt.Sprintf("Iterations: %d\n", iterationsFlag))
	msg(fmt.Sprintf("# CPUs: %d (was %d)\n", numcpusFlag, oldprocs))
	msg("\n")

	// Create worker quit channel
	quit := make(chan bool)

	// Spawn worker threads
	msg(fmt.Sprintf("\tspawning %d workers\n", workersFlag))
	for i := 0; i < workersFlag; i++ {
		go worker(queue, quit)
	}

	// Send jobs
	msg(fmt.Sprintf("\tsending %d jobs to queue(s)\n", iterationsFlag))
	start := time.Now()
	for i := 0; i < iterationsFlag; i++ {
		msgD("+")
		queue <- i
	}

	// Wait for jobs to finish processing
	msg(fmt.Sprintf("\n\tWaiting for jobs to complete...\n"))
	for len(queue) > 0 {
		time.Sleep(time.Millisecond * 1)
	}
	end := time.Now()
	msg("\n")

	// Kill workers (closing the channel will instruct all workers to exit)
	close(quit)

	// Wait for workers to all be dead (we want their exit messages to display)
	for workers > 0 {
		time.Sleep(time.Microsecond * 100)
	}

	// Calculate elapsed time (just for job processing) and print out results
	elapsedSeconds := end.Sub(start).Seconds()
	iterPerSec := float64(iterationsFlag) / elapsedSeconds
	fmt.Printf("\nElapsed time: %.4f secs\n", elapsedSeconds)
	fmt.Printf("Jobs per second: %.4f\n", iterPerSec)
}
