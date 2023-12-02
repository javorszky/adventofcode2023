package day1

import (
	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
	"os"
	"regexp"
	"strconv"
)

var (
	reOne   = regexp.MustCompile("one|1")
	reTwo   = regexp.MustCompile("two|2")
	reThree = regexp.MustCompile("three|3")
	reFour  = regexp.MustCompile("four|4")
	reFive  = regexp.MustCompile("five|5")
	reSix   = regexp.MustCompile("six|6")
	reSeven = regexp.MustCompile("seven|7")
	reEight = regexp.MustCompile("eight|8")
	reNine  = regexp.MustCompile("nine|9")

	reMap = map[int]*regexp.Regexp{
		1: reOne,
		2: reTwo,
		3: reThree,
		4: reFour,
		5: reFive,
		6: reSix,
		7: reSeven,
		8: reEight,
		9: reNine,
	}
)

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 1).Int("part", 2).Logger()

	gog, err := inputs.ReadIntoLines("day1/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	sum := 0

	for _, l := range gog {
		digits, err := replaceSpelled(l)
		if err != nil {
			localLogger.Fatal().Msgf("line %s encountered an issue %s", l, err.Error())
		}

		sum += digits
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Once written replaced, the sum becomes %d", solution)
}

func replaceSpelled(line string) (int, error) {
	zeroes := make([]int, len(line), len(line))

	for k, r := range reMap {
		bla := r.FindAllStringSubmatchIndex(line, -1)

		for _, m := range bla {
			zeroes[m[0]] = k
		}
	}

	nslice := make([]int, 0)
	for _, d := range zeroes {
		switch d {
		case 0:
			continue
		default:
			nslice = append(nslice, d)
		}
	}

	return strconv.Atoi(strconv.Itoa(nslice[0]) + strconv.Itoa(nslice[len(nslice)-1]))
}
