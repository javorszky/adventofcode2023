package day5

import (
	"fmt"
	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
	"os"
	"runtime"
	"sort"
	"sync"
)

func Task2(l zerolog.Logger) {
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
	locations := make([]int, 0, 10000000)

	ch := make(chan int, 100)
	out := make(chan int, 100)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			a := <-ch
			out <- walkSeed(m, a)
			fmt.Printf("-- generated a number using walkseed\n")
		}()
	}

	go func() {
		for {
			fmt.Printf("- getting number out of out channel\n")
			a, ok := <-out
			if !ok {
				return
			}

			locations = append(locations, a)
		}
	}()

	// write a function that receives from the out channel

	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Printf("this is seeds\n%v\n", seeds)

	go func(wait *sync.WaitGroup) {

		for i := 0; i < len(seeds); i = i + 2 {
			for j := 0; j < seeds[i+1]; j++ {
				fmt.Printf("putting seed %d (i %d j %d) into the channel\n", seeds[i]+j, i, j)
				ch <- seeds[i] + j
			}
		}

		wait.Done()
	}(&wg)

	wg.Wait()

	return locations
}

// destination - source - length
