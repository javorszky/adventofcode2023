package day5

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

var (
	groupMap = []string{
		"seeds",
		"seed2soil",
		"soil2fertilizer",
		"fertilizer2water",
		"water2light",
		"light2temperature",
		"temperature2humidity",
		"humidity2location",
	}

	seedOrder     = []string{"seed", "soil", "fertilizer", "water", "light", "temperature", "humidity", "location"}
	locationOrder = []string{"location", "humidity", "temperature", "light", "water", "fertilizer", "soil", "seed"}

	fromSeed = map[string]string{
		"seed":        "soil",
		"soil":        "fertilizer",
		"fertilizer":  "water",
		"water":       "light",
		"light":       "temperature",
		"temperature": "humidity",
		"humidity":    "location",
	}
	fromLocation map[string]string
)

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 5).Int("part", 1).Logger()

	fromLocation = make(map[string]string)

	// invert so we have the other way as well.
	for k, v := range fromSeed {
		fromLocation[v] = k
	}

	groups, err := inputs.GroupByBlankLines("day5/example.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	m, err := parseGroups(groups)
	if err != nil {
		localLogger.Err(err).Msg("parseGroups returned an error")
	}

	// extract the seeds first
	seeds, err := parseSeeds(groups[0][0])
	if err != nil {
		localLogger.Err(err).Msgf("parseSeeds")
	}

	locations := make([]int, len(seeds))

	for i, seed := range seeds {
		locations[i] = walkSeed(m, seed)
	}

	sort.Ints(locations)

	solution := locations[0]
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Lowest location number is %d", solution)
}

func parseSeeds(line string) ([]int, error) {
	seedsTrimPrefix := strings.TrimPrefix(line, "seeds: ")
	seedsTrimSpace := strings.TrimSpace(seedsTrimPrefix)

	seedsStringSlice := strings.Split(seedsTrimSpace, " ")
	seeds := make([]int, len(seedsStringSlice))

	for i, seedString := range seedsStringSlice {
		seedNumber, err := strconv.Atoi(seedString)
		if err != nil {
			return nil, errors.Wrapf(err, "strconv.Atoi: %s", seedString)
		}

		seeds[i] = seedNumber
	}

	return seeds, nil
}

func parseGroups(group [][]string) (map[string][]map[string]int, error) {
	m := make(map[string][]map[string]int)

	// extract the rest of the groups
	for i := 1; i < len(group); i++ {
		newMap, err := parseGroup(groupMap[i], group[i])
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %s: %v", groupMap[i], group[i])
		}

		for k, v := range newMap {
			m[k] = v
		}
	}

	return m, nil
}

func parseGroup(label string, group []string) (map[string][]map[string]int, error) {
	m := make(map[string][]map[string]int)

	labelParts := strings.Split(label, "2")
	flipLabel := fmt.Sprintf("%s2%s", labelParts[1], labelParts[0])

	m[label] = make([]map[string]int, 0)
	m[flipLabel] = make([]map[string]int, 0)

	dataSlice := [3]int{}

	for i := 1; i < len(group); i++ {
		numSlice := strings.Split(group[i], " ")
		for j, numString := range numSlice {
			n, err := strconv.Atoi(numString)
			if err != nil {
				return nil, errors.Wrapf(err, "strconv.Atoi: %s in %s", numString, group[i])
			}

			dataSlice[j] = n
		}

		m[label] = append(m[label], map[string]int{
			"source": dataSlice[1],
			"dest":   dataSlice[0],
			"len":    dataSlice[2],
		})

		m[flipLabel] = append(m[flipLabel], map[string]int{
			"source": dataSlice[0],
			"dest":   dataSlice[1],
			"len":    dataSlice[2],
		})
	}

	return m, nil
}

func walkSeed(groupMap map[string][]map[string]int, seed int) int {
	fromLocation := 0
	destination := seed

	for i := 0; i < len(seedOrder)-1; i++ {
		fromLocation = destination
		key := fmt.Sprintf("%s2%s", seedOrder[i], seedOrder[i+1])

		destination = findMap(fromLocation, groupMap[key])
	}

	return destination
}

func findMap(source int, group []map[string]int) int {
	for _, gMap := range group {

		if source < gMap["source"] {
			continue
		}

		if source-gMap["len"] > gMap["source"] {
			continue
		}

		return gMap["dest"] + (source - gMap["source"])
	}

	return source
}
