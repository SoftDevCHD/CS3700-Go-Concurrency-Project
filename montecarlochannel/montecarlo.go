package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

const MINT int = 1
const MAXT int = 10
const MINP int = 10
const MAXP int = 1000000

func simulate(chan_a chan<- float64, iterations int, wg *sync.WaitGroup) {
	var part_a float64 = 0
	for i := 0; i < iterations; i++ {
		x := rand.Float64()
		y := rand.Float64()
		r := math.Sqrt(x*x + y*y)

		if r <= 1.0 {
			part_a++
		}
	}
	chan_a <- part_a
	wg.Done()
}

func main() {
	fmt.Println()
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Program requires 2 arguments: number of threads and number of data points")
		os.Exit(1)
	}

	numThreads, err1 := strconv.Atoi(args[0])
	numDataPoints, err2 := strconv.Atoi(args[1])

	if err1 != nil || numThreads < MINT || numThreads > MAXT {
		fmt.Println("Number of threads must be between 1 and 10")
		os.Exit(1)
	}

	if err2 != nil || numDataPoints < MINP || numDataPoints > MAXP {
		fmt.Println("Number of data points must be between 10 and 1000000")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	chan_a := make(chan float64, numThreads)

	start := time.Now()

	var slice int = numDataPoints / numThreads
	var remainder int = numDataPoints % numThreads
	var startIdx int = 0

	for i := 0; i < numThreads; i++ {
		wg.Add(1)

		if i == numThreads-1 {
			slice += remainder
		}

		go simulate(chan_a, slice, &wg)

		fmt.Printf("Goroutine/Thread %d: %d to %d\n", i, startIdx, startIdx+slice)
		startIdx += slice
	}

	wg.Wait()

	close(chan_a)

	var total_a float64 = 0
	for n := range chan_a {
		total_a += n
	}

	pi := 4.0 * total_a / float64(numDataPoints)
	delta := math.Abs(pi - math.Pi)

	elapsed := time.Since(start)

	fmt.Printf("Pi: %f\n", pi)
	fmt.Printf("Delta: %f\n", delta)
	fmt.Printf("Elapsed Time: %s\n", elapsed)
	fmt.Println()
}
