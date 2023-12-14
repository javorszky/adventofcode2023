package day11

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const (
	star  = "#"
	empty = "."
)

var reGalaxy = regexp.MustCompile(`#`)

type galaxy struct {
	row, col int
}

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 11).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day11/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	// grow that, ie for every horizontal line without a #, add another identical horizontal line
	grownColumns := growSpace(gog)

	// rotate back, so now we have double columns of empty space
	rowsWithGrownColumns := rowsToColumns(grownColumns)

	// repeat the empty line for the proper rows again
	grownRowsAndColumns := growSpace(rowsWithGrownColumns)

	// find all the galaxies in the now grown universe
	galaxies := parseGalaxies(grownRowsAndColumns)

	// pair them up
	pairs := pairGalaxies(galaxies)

	sum := 0
	for _, p := range pairs {
		sum += distance(p[0], p[1])
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Sum of distances between galaxies in a grown universe is %d", solution)
}

func rowsToColumns(gog []string) []string {
	cols := make([][]string, len(gog[0]))

	// for each element in gog, use the line as row
	for i, row := range gog {
		// for each character in row, use them
		for j, char := range row {
			if i == 0 {
				cols[j] = make([]string, len(gog))
			}

			cols[j][i] = string(char)
		}
	}

	colsAgain := make([]string, len(cols))
	for i, col := range cols {
		colsAgain[i] = strings.Join(col, "")
	}

	return colsAgain
}

func growSpace(in []string) []string {
	grown := make([]string, 0, len(in))

	for _, line := range in {
		grown = append(grown, line)
		if strings.Count(line, star) == 0 {
			grown = append(grown, line)
		}
	}

	return grown
}

func visualize(in []string) {
	var sb strings.Builder

	for i := 0; i < len(in); i++ {
		sb.WriteString(in[i])
		sb.WriteString("\n")
	}

	fmt.Printf(sb.String())
}

func parseGalaxies(in []string) []galaxy {
	gxs := make([]galaxy, 0)
	for rowIndex, row := range in {
		found := reGalaxy.FindAllStringIndex(row, -1)

		for _, pair := range found {
			gxs = append(gxs, galaxy{
				row: rowIndex,
				col: pair[0],
			})
		}
	}

	return gxs
}

func pairGalaxies(gxs []galaxy) [][2]galaxy {
	gxyPairs := make([][2]galaxy, 0)
	l := len(gxs)
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			gxyPairs = append(gxyPairs, [2]galaxy{gxs[i], gxs[j]})
		}
	}

	return gxyPairs
}

func distance(a, b galaxy) int {
	vert := a.row - b.row
	horz := a.col - b.col

	if vert < 0 {
		vert = -vert
	}

	if horz < 0 {
		horz = -horz
	}

	//fmt.Printf("          %+v : diff: %d\n", [2]galaxy{a, b}, vert+horz)

	return vert + horz
}
