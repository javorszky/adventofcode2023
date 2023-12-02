package day2

import (
	"os"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 2).Int("part", 2).Logger()

	gog, err := inputs.ReadIntoLines("day2/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	sum := 0

	for _, line := range gog {
		_, parsed, err := parse(line)
		if err != nil {
			localLogger.Fatal().Msgf("parsing line %s encountered error: %s", line, err.Error())
		}

		localMaxRed, localMaxGreen, localMaxBlue := 0, 0, 0

		for _, rgb := range parsed {
			if rgb[0] > localMaxRed {
				localMaxRed = rgb[0]
			}

			if rgb[1] > localMaxGreen {
				localMaxGreen = rgb[1]
			}

			if rgb[2] > localMaxBlue {
				localMaxBlue = rgb[2]
			}
		}

		sum += localMaxRed * localMaxGreen * localMaxBlue
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Sum of the power of minimum sets of cubes is %d", solution)
}
