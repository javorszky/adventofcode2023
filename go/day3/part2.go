package day3

import (
	"fmt"
	"os"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const gear = "*"

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 3).Int("part", 2).Logger()

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

	gears, err := collectGears(gog, numbers)
	if err != nil {
		localLogger.Err(err).Msgf("collectGears made a booboo")
	}

	ratios := calculateRatios(gears)

	solution := ratios
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Sum of all products for all gears with only two numbers next to them is %d", solution)
}

func collectGears(gog []string, m map[int][][3]int) (map[string][]int, error) {
	asterisks := make(map[string][]int)

	for number, coords := range m {
		for _, c := range coords {
			if c[0] > 0 {
				above := gog[c[0]-1][lessOrZero(c[1]):moreOrMax(c[2])]

				// if the line above has more than one gear, return an error, so we know we need to do something diff
				if strings.Count(above, gear) > 1 {
					return nil, fmt.Errorf("number %d has more than one asterisk above it: %s", number, above)
				}

				// get the index of the gear in the line above the number
				aboveGearIndex := strings.Index(above, gear)

				// if the index is not -1, ie there is one (and not more than one), then add the number to the map which
				// has a key of the coordinate of the gear
				if aboveGearIndex != -1 {
					key := fmt.Sprintf("%d-%d", c[0]-1, lessOrZero(c[1])+aboveGearIndex)
					asterisks[key] = append(asterisks[key], number)
				}
			}

			sameLine := gog[c[0]][lessOrZero(c[1]):moreOrMax(c[2])]
			if strings.Count(sameLine, gear) > 1 {
				return nil, fmt.Errorf("number %d has more than one asterisk above it: %s", number, sameLine)
			}

			// get the index of the gear in the line above the number
			sameLineGearIndex := strings.Index(sameLine, gear)

			if sameLineGearIndex != -1 {
				key := fmt.Sprintf("%d-%d", c[0], lessOrZero(c[1])+sameLineGearIndex)
				asterisks[key] = append(asterisks[key], number)
			}

			if c[0] < len(gog)-1 {
				below := gog[c[0]+1][lessOrZero(c[1]):moreOrMax(c[2])]

				if strings.Count(below, gear) > 1 {
					return nil, fmt.Errorf("number %d has more than one asterisk above it: %s", number, below)
				}

				// get the index of the gear in the line above the number
				belowGearIndex := strings.Index(below, gear)

				if belowGearIndex != -1 {
					key := fmt.Sprintf("%d-%d", c[0]+1, lessOrZero(c[1])+belowGearIndex)
					asterisks[key] = append(asterisks[key], number)
				}
			}
		}
	}

	return asterisks, nil
}

func calculateRatios(gears map[string][]int) int {
	sum := 0

	for _, numbers := range gears {
		if len(numbers) != 2 {
			continue
		}

		sum += numbers[0] * numbers[1]
	}

	return sum
}
