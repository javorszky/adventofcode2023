package day2

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const (
	maxRed   = 12
	maxGreen = 13
	maxBlue  = 14
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 2).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day2/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	sum := 0

Line:
	for _, line := range gog {
		game, parsed, err := parse(line)
		if err != nil {
			localLogger.Fatal().Msgf("parsing line %s encountered error: %s", line, err.Error())
		}

		for _, rgb := range parsed {
			if rgb[0] > maxRed || rgb[1] > maxGreen || rgb[2] > maxBlue {
				continue Line
			}
		}

		sum += game
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Sum of games for possible plays is %d", solution)
}

func parse(line string) (int, [][3]int, error) {
	// trim 'Game ' from the front
	line = strings.TrimPrefix(line, "Game ")

	// explode into number left and collections right
	parts := strings.Split(line, ":")
	game, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, nil, errors.Wrap(err, "strconv.Atoi parts[0]")
	}

	parsed := make([][3]int, 0)

	// explode into collections
	collections := strings.Split(parts[1], ";")
	for _, c := range collections {
		red, green, blue := 0, 0, 0
		// for each collection, explode into parts
		colors := strings.Split(c, ",")

		for _, color := range colors {
			// for each color data, explode into number and color
			data := strings.Split(strings.TrimSpace(color), " ")

			// convert the number to an integer
			num, err := strconv.Atoi(data[0])
			if err != nil {
				return 0, nil, errors.Wrapf(err, "strconv.Atoi, data[0]: %s", color)
			}

			// assign that integer to either red green blue based on what color it is
			switch data[1] {
			case "red":
				red = num
			case "green":
				green = num
			case "blue":
				blue = num
			}
		}

		// record the data for the collection
		parsed = append(parsed, [3]int{red, green, blue})

	}

	return game, parsed, nil
}
