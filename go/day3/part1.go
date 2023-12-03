package day3

import (
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	"github.com/javorszky/adventofcode2023/inputs"
)

var (
	//rePossibleNonNumbers = regexp.MustCompile(`\.(\d+)\.|^\d+\.|\.\d+$`)
	//reDefiniteNumbers    = regexp.MustCompile(``)
	reNumbers = regexp.MustCompile(`\d+`)
	reSymbols = regexp.MustCompile(`[^\.\d]`)

	maxLen = 0
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 3).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day3/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	maxLen = len(gog[0])

	numbers, err := numberPositions(gog)
	if err != nil {
		localLogger.Err(err).Msgf("numberPositions")
	}

	sum := collectNumbers(gog, numbers)

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Sum of all part numbers is %d", solution)
}

func numberPositions(lines []string) (map[int][][3]int, error) {
	m := make(map[int][][3]int)

	// for each line
	for row, line := range lines {
		// find all the numbers and get their coordinates
		nums := reNumbers.FindAllStringSubmatchIndex(line, -1)
		for _, coords := range nums {
			// for each set of coordinates, get the actual number in there
			nString := line[coords[0]:coords[1]]

			// convert the string into a number
			n, err := strconv.Atoi(nString)
			if err != nil {
				return nil, errors.Wrapf(err, "strconv.Atoi for %s", nString)
			}

			m[n] = append(m[n], [3]int{row, coords[0], coords[1]})
		}
	}

	return m, nil
}

func collectNumbers(gog []string, m map[int][][3]int) int {
	sum := 0

	for number, coords := range m {
		for _, c := range coords {
			var sb strings.Builder
			if c[0] > 0 {
				sb.WriteString(gog[c[0]-1][lessOrZero(c[1]):moreOrMax(c[2])])
			}

			sb.WriteString(gog[c[0]][lessOrZero(c[1]):moreOrMax(c[2])])

			if c[0] < len(gog)-1 {
				sb.WriteString(gog[c[0]+1][lessOrZero(c[1]):moreOrMax(c[2])])
			}

			if reSymbols.MatchString(sb.String()) {
				sum += number
			}
		}
	}

	return sum
}

func lessOrZero(i int) int {
	j := i - 1
	if j < 0 {
		return 0
	}

	return j
}

func moreOrMax(i int) int {
	j := i + 1
	if j > maxLen {
		return maxLen
	}

	return j
}
