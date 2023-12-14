package day11

import (
	"os"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const expandRate = 1000000

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 11).Int("part", 2).Logger()

	gog, err := inputs.ReadIntoLines("day11/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	// grow that, ie for every horizontal line without a #, add another identical horizontal line
	rowMillionsAt := growSpace2(gog)

	// rotate back, so now we have double columns of empty space
	rotated := rowsToColumns(gog)

	// repeat the empty line for the proper rows again
	colMillionsAt := growSpace2(rotated)

	//fmt.Printf("row millions at\n%v\n\ncol millions at\n%v\n\n", rowMillionsAt, colMillionsAt)

	// find all the galaxies in the now grown universe
	galaxies := parseGalaxies(rotated)

	// pair them up
	pairs := pairGalaxies(galaxies)

	d2 := distance2(colMillionsAt, rowMillionsAt)

	sum := 0
	for _, p := range pairs {
		sum += d2(p)
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("With a expandRate empty lines, the distance sum is %d", solution)
}

func growSpace2(in []string) []int {
	whereSpace := make([]int, 0)

	for i, line := range in {
		if strings.Count(line, star) == 0 {
			whereSpace = append(whereSpace, i)
		}
	}

	return whereSpace
}

func distance2(rowMils, colMils []int) func(gxyPairs [2]galaxy) int {

	return func(gxyPairs [2]galaxy) int {
		var smolRow, bigRow, smolCol, bigCol int

		// figure out which one is the larger row
		if gxyPairs[0].row > gxyPairs[1].row {
			bigRow = gxyPairs[0].row
			smolRow = gxyPairs[1].row
		} else {
			smolRow = gxyPairs[0].row
			bigRow = gxyPairs[1].row
		}

		// figure out which one is the larger column
		if gxyPairs[0].col > gxyPairs[1].col {
			bigCol = gxyPairs[0].col
			smolCol = gxyPairs[1].col
		} else {
			smolCol = gxyPairs[0].col
			bigCol = gxyPairs[1].col
		}

		milRows := 0
		for _, rm := range rowMils {
			if smolRow < rm && rm < bigRow {
				milRows++
			}
		}

		bigRow = bigRow + milRows*(expandRate-1)

		milCols := 0
		for _, cm := range colMils {
			if smolCol < cm && cm < bigCol {
				milCols++
			}
		}

		bigCol = bigCol + milCols*(expandRate-1)

		// get the differences without the expandRate miles
		diffRow := bigRow - smolRow
		diffCol := bigCol - smolCol

		diff := diffRow + diffCol

		//fmt.Printf("originals: %+v\n"+
		//	"adjusted: %+v\n\n", gxyPairs, [2]galaxy{
		//	{row: smolRow, col: smolCol},
		//	{row: bigRow, col: bigCol},
		//})

		return diff
	}
}
