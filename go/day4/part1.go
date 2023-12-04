package day4

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 4).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day4/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	data, err := parse(gog)
	if err != nil {
		localLogger.Err(err).Msgf("parsing incoming lines failed")
	}

	score := 0

	for _, card := range data {
		found := findUnions(card[0], card[1])
		if found == 0 {
			continue
		}

		score += 1 << (found - 1)
	}

	solution := score
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Score of all scratchcards is %d", solution)
}

func parse(lines []string) (map[int][2][]int, error) {
	m := make(map[int][2][]int)

	for _, line := range lines {
		// separate by : to get "card xxx" and the numbers
		parts := strings.Split(line, ":")

		// remove Card from the first part, trim, and convert to an int to get card number
		n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(parts[0], "Card")))
		if err != nil {
			return nil, errors.Wrapf(err, "strconv.Atoi for line: '%s', trimmed: '%s', spacetrimmed: '%s'",
				parts[0],
				strings.TrimPrefix(parts[0], "Card"),
				strings.TrimSpace(strings.TrimPrefix(parts[0], "Card")),
			)
		}

		// separate the numbers by the pipe
		numberHalves := strings.Split(parts[1], "|")

		wins, err := parseNumbers(numberHalves[0])
		if err != nil {
			return nil, errors.Wrapf(err, "parsing first half of numbers: %s", numberHalves[0])
		}

		haves, err := parseNumbers(numberHalves[1])
		if err != nil {
			return nil, errors.Wrapf(err, "parsing second half of numbers: %s", numberHalves[1])
		}

		m[n] = [2][]int{
			wins,
			haves,
		}
	}

	return m, nil
}

func parseNumbers(line string) ([]int, error) {
	// first replace the double spaces with a single space "  " -> " "
	line = strings.ReplaceAll(strings.TrimSpace(line), "  ", " ")

	// then explode it into constituent elements by the space
	numbersString := strings.Split(line, " ")

	nums := make([]int, len(numbersString))
	// then convert it into numbers
	for i, n := range numbersString {
		parsedNumber, err := strconv.Atoi(n)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing number %s", n)
		}

		nums[i] = parsedNumber
	}

	return nums, nil
}

func findUnions(a, b []int) int {
	am := make(map[int]struct{})
	for _, el := range a {
		am[el] = struct{}{}
	}

	found := 0

	for _, el := range b {
		_, ok := am[el]
		if ok {
			found += 1
		}
	}

	return found
}
