package day9

import (
	"github.com/pkg/errors"
	"os"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 9).Int("part", 2).Logger()

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
	//p, err := predictFront(lines[2])
	//if err != nil {
	//	localLogger.Err(err).Msgf("trying this line: %v, the predict error happened", lines[0])
	//}
	//
	//fmt.Printf("prediction is %d\n", p)

	for _, line := range lines {
		p, err := predictFront(line)
		if err != nil {
			localLogger.Err(err).Msgf("trying this line: %v, the predict error happened", lines[0])
		}

		sum += p
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Predicting backwards yields %d", solution)
}

func predictFront(in []int) (int, error) {
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
		d, err := predictFront(diffs)
		if err != nil {
			return 0, errors.Wrapf(err, "predict returned error for this slice: %v", diffs)
		}

		//fmt.Printf("predicted %d for the diffs\n", d)

		delta = d
	} else {
		//fmt.Printf("diffs all zero: %v\n", diffs)
	}

	//fmt.Printf("returning %d as the prediction\n", in[0]-delta)
	return in[0] - delta, nil
}
