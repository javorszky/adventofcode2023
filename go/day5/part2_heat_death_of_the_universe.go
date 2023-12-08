package day5

import (
	"fmt"
	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

func Task2HeatDeathOfTheUniverse(l zerolog.Logger) {
	localLogger := l.With().Int("day", 5).Int("part", 2).Logger()

	groups, err := inputs.GroupByBlankLines("day5/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	m, err := parseGroups(groups)
	if err != nil {
		localLogger.Err(err).Msg("parseGroups returned an error")
	}

	// extract the seeds first
	seeds, err := parseSeeds(groups[0][0])
	if err != nil {
		localLogger.Err(err).Msgf("parseSeeds")
	}

	locations := walkSeedsThreaded(m, seeds)

	sort.Ints(locations)
	solution := locations[0]
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Lowest location when the seeds are ranges is %d", solution)
}

func walkSeedsThreaded(m map[string][]map[string]int, seeds []int) []int {
	in := make(chan int)
	out := make(chan int)

	// launch n number of go funcs that read from a channel and call walkSeed.
	for i := 0; i < runtime.NumCPU()*10; i++ {
		go func(cpuThread int) {
			for {
				//fmt.Printf("- - (fn %d) grabbing number out from the thing\n", cpuThread)
				a, ok := <-in
				if !ok {
					return
				}

				//fmt.Printf("- - (fn %d) generating number and putting it into the out channel\n", cpuThread)
				out <- walkSeed(m, a)
				//fmt.Printf("- - (fn %d) generated a number %d using walkseed\n", cpuThread, a)
			}
		}(i)
	}

	loc := make([]int, 0)

	go func() {
		// launch a for loop that reads out from the channel and populates the end array
		for {
			//fmt.Printf("- getting number out of out channel\n")
			a, ok := <-out
			if !ok {
				//fmt.Printf("- <- out was not okay, returned...")
				return
			}

			loc = append(loc, a)
			//fmt.Printf("- appended the number to the slice, slice len is now %d\n", len(loc))
		}
	}()

	// write a function that receives from the out channel

	var wg sync.WaitGroup
	wg.Add(1)

	go func(wait *sync.WaitGroup) {
		//fmt.Printf("==== len seeds is %d\n", len(seeds))
		for i := 0; i < len(seeds); i = i + 2 {
			//fmt.Printf("==== seeds i + 1 is %d\n", seeds[i+1])
			for j := 0; j < seeds[i+1]; j++ {
				//fmt.Printf("==== putting seed %d (i %d j %d) into the channel\n", seeds[i]+j, i, j)
				in <- seeds[i] + j
				//fmt.Printf("==== put the number in the channel, advancing...\n")
			}
		}

		//fmt.Printf("==== we're done with shoving all the numbers in the channel\n")
		wait.Done()
	}(&wg)

	fmt.Printf("main thread, waiting for finish...\n")
	wg.Wait()

	time.Sleep(time.Second)

	return loc
}

// destination - source - length
