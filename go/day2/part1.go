package day2

import (
	"os"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 2).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day2/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	// code goes here

	solution := 2
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("-- change this for the part 1 message -- %d", solution)
}
