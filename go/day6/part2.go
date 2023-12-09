package day6

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 6).Int("part", 2).Logger()

	gog, err := inputs.ReadIntoLines("day6/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	times, err := parseLinePart2("Time:", gog[0])
	if err != nil {
		localLogger.Err(err).Msgf("parseLine for Time:")
	}

	distances, err := parseLinePart2("Distance:", gog[1])
	if err != nil {
		localLogger.Err(err).Msgf("parseLine for Distance:")
	}

	product := 1

	for i := 0; i < len(times); i++ {
		product *= findX(times[i], distances[i])
	}

	solution := product
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Margin of error for part 2 is %d", solution)
}

func parseLinePart2(prefix, line string) ([]int, error) {
	noPrefix := strings.TrimPrefix(line, prefix)

	removeKerning := strings.ReplaceAll(noPrefix, " ", "")

	n, err := strconv.Atoi(removeKerning)
	if err != nil {
		return nil, errors.Wrapf(err, "strconv.Atoi: %s", removeKerning)
	}

	return []int{n}, nil
}
