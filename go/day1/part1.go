package day1

import (
	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
	"os"
	"strconv"
	"strings"
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 1).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day1/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	sum := 0

	for _, l := range gog {
		filtered := filterForNumbers(l)
		n, err := twoDigit(filtered)
		if err != nil {
			localLogger.Fatal().Msgf("twoDigit: %s", err.Error())
		}

		sum += n
	}
	// code goes here

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()

	s.Info().Msgf("Sum of all two digit numbers from the first and last digits in each line is %d", solution)
}

func filterForNumbers(line string) string {
	var sb strings.Builder

	for _, c := range line {
		switch c {
		case 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
			sb.WriteByte(byte(c))
		}
	}

	return sb.String()
}

func twoDigit(line string) (int, error) {
	var sb strings.Builder
	sb.WriteByte(byte(line[0]))
	sb.WriteByte(byte(line[len(line)-1]))
	return strconv.Atoi(sb.String())
}
