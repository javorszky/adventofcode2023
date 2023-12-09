package day6

import (
	"github.com/pkg/errors"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 6).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day6/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	times, err := parseLine("Time:", gog[0])
	if err != nil {
		localLogger.Err(err).Msgf("parseLine for Time:")
	}

	distances, err := parseLine("Distance:", gog[1])
	if err != nil {
		localLogger.Err(err).Msgf("parseLine for Distance:")
	}

	product := 1

	for i := 0; i < len(times); i++ {
		product *= findX(times[i], distances[i])
	}

	solution := product
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("The margin of error on the boat race is %d", solution)
}

func parseLine(prefix, line string) ([]int, error) {
	stringParts := strings.Split(strings.TrimPrefix(line, prefix), " ")

	stringNumbers := make([]string, 0)
	for _, el := range stringParts {
		if el != "" {
			stringNumbers = append(stringNumbers, el)
		}
	}

	nums := make([]int, len(stringNumbers))
	for i, s := range stringNumbers {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrapf(err, "strconv.Atoi: '%s'", s)
		}

		nums[i] = n
	}

	return nums, nil
}

func findX(time, distance int) int {
	tf := float64(time)
	df := float64(distance)

	res1 := (tf + math.Sqrt((tf*tf)-4*1*df)) / 2
	res2 := (tf - math.Sqrt((tf*tf)-4*1*df)) / 2

	var smaller, larger float64
	if res1 < res2 {
		smaller = res1
		larger = res2
	} else {
		smaller = res2
		larger = res1
	}

	lower := int(math.Ceil(smaller))
	upper := int(math.Floor(larger))

	//fmt.Printf(""+
	//	"uh, okay, so result 1 and 2 are %f and %f\n"+
	//	"smaller and larger are %f and %f\n"+
	//	"lower and upper are %d and %d, and their difference is %d\n\n",
	//	res1, res2,
	//	smaller, larger,
	//	lower, upper, upper-lower)

	return upper - lower + 1
}
