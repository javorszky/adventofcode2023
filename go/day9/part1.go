package day9

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 9).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day9/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	lines := make([][]int, len(gog))

	for i, line := range gog {
		numbers, err := parseIntoInts(line)
		if err != nil {
			localLogger.Err(err).Msgf("Turning line '%s' into ints died", line)
		}

		lines[i] = numbers
	}

	sum := 0

	for _, line := range lines {
		p, err := predict(line)
		if err != nil {
			localLogger.Err(err).Msgf("trying this line: %v, the predict error happened", lines[0])
		}

		sum += p
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Sum of predictions is %d", solution)
}

func parseIntoInts(line string) ([]int, error) {
	numbersAsStrings := strings.Split(line, " ")

	nums := make([]int, len(numbersAsStrings))
	for i, s := range numbersAsStrings {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrapf(err, "strconv.Atoi %s", s)
		}

		nums[i] = n
	}

	return nums, nil
}

func getDifferences(in []int) ([]int, error) {
	if len(in) < 2 {
		return nil, fmt.Errorf("can't take the differences of incoming slice: %v", in)
	}

	diffs := make([]int, len(in)-1)

	for i := 1; i < len(in); i++ {
		diffs[i-1] = in[i] - in[i-1]
	}

	return diffs, nil
}

func predict(in []int) (int, error) {
	//fmt.Printf("======== start ========\n"+
	//	"predicting next for %v\n", in)

	diffs, err := getDifferences(in)
	if err != nil {
		return 0, errors.Wrapf(err, "getDifferences for slice %v", in)
	}

	allZero := true
	for _, v := range diffs {
		if v != 0 {
			//fmt.Printf("diffs not all zero (%v)\n", diffs)
			allZero = false
			break
		}
	}

	delta := 0

	if !allZero {
		d, err := predict(diffs)
		if err != nil {
			return 0, errors.Wrapf(err, "predict returned error for this slice: %v", diffs)
		}

		//fmt.Printf("predicted %d for the diffs\n", d)

		delta = d
	}

	//fmt.Printf("returning %d as the prediction\n", delta+diffs[len(diffs)-1])
	return delta + in[len(in)-1], nil
}
