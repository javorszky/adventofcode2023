package day2

import (
	"os"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 2).Int("part", 2).Logger()

	_, err := inputs.ReadIntoLines("day2/input2.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	// code goes here

	solution := 2
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("-- change this to the part 2 message -- %d", solution)
}
