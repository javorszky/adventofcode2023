package day5

import (
	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
	"os"
	"sort"
)

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 5).Int("part", 1).Logger()

	fromLocation = make(map[string]string)

	// invert so we have the other way as well.
	for k, v := range fromSeed {
		fromLocation[v] = k
	}

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

	locations := make([]int, len(seeds))

	for i, seed := range seeds {
		locations[i] = walkSeed(m, seed)
	}

	sort.Ints(locations)

	solution := locations[0]
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Lowest location number is %d", solution)
}
