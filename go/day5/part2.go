package day5

import (
	"os"
	"sort"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
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

	locations := make([]int, 0)

	for i := 0; i < len(seeds); i = i + 2 {
		for j := 0; j < seeds[i+1]; j++ {
			locations = append(locations, walkSeed(m, seeds[i]+j))
		}
	}

	sort.Ints(locations)
	solution := locations[0]
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Lowest location when the seeds are ranges is %d", solution)
}

// destination - source - length
