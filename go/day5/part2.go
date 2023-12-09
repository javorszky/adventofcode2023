package day5

import (
	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"sort"
	"strconv"
	"strings"
)

type section struct {
	sourceStart, sourceEnd, destinationStart, destinationEnd, delta int
}

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 5).Int("part", 2).Logger()

	fromLocation = make(map[string]string)

	// invert so we have the other way as well.
	for k, v := range fromSeed {
		fromLocation[v] = k
	}

	groups, err := inputs.GroupByBlankLines("day5/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	// extract the seeds first
	seeds, err := parseSeeds(groups[0][0])
	if err != nil {
		localLogger.Err(err).Msg("parseSeeds")
	}

	seedSections := parseSeedsToSections(seeds)

	sections, err := parseGroupsToSections(groups)
	if err != nil {
		localLogger.Err(err).Msg("parseGroupsToSections")
	}

	sections["seeds"] = seedSections

	keys := make([]string, 0)
	for k, _ := range sections {
		keys = append(keys, k)
	}

	var (
		previous []section
	)

	previous = sections["seeds"]

	for i := 1; i < len(groupMap); i++ {
		//fmt.Printf("\n\n================ Parsing %s ================\n\n", groupMap[i])
		//fmt.Printf("previous section: %#v\nnew section: %#v\n", previous, sections[groupMap[i]])

		previous = overlaps(previous, sections[groupMap[i]])
	}

	lowestDestination := 2 << 61

	for _, sect := range previous {
		if sect.destinationStart < lowestDestination {
			lowestDestination = sect.destinationStart
		}
	}

	solution := lowestDestination
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Lowest location number is %d", solution)
}

func parseSeedsToSections(seeds []int) []section {
	m := make([]section, 0)

	for i := 0; i < len(seeds); i = i + 2 {
		m = append(m, section{
			sourceStart:      seeds[i],
			sourceEnd:        seeds[i] + seeds[i+1],
			destinationStart: seeds[i],
			destinationEnd:   seeds[i] + seeds[i+1],
			delta:            0,
		})
	}

	return m
}

func parseGroupsToSections(group [][]string) (map[string][]section, error) {
	m := make(map[string][]section)

	// extract the rest of the groups
	for i := 1; i < len(group); i++ {
		newMap, err := groupToSection(groupMap[i], group[i])
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %s: %v", groupMap[i], group[i])
		}

		for k, v := range newMap {
			m[k] = v
		}
	}

	return m, nil
}

func groupToSection(label string, group []string) (map[string][]section, error) {
	m := make(map[string][]section)
	m[label] = make([]section, 0)

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

		m[label] = append(m[label], section{
			sourceStart:      dataSlice[1],
			sourceEnd:        dataSlice[1] + dataSlice[2] - 1,
			destinationStart: dataSlice[0],
			destinationEnd:   dataSlice[0] + dataSlice[2] - 1,
			delta:            dataSlice[0] - dataSlice[1],
		})
	}

	return m, nil
}

// overlaps will return a slice of section such that the target section only covers the connecting side of the source
// sections. That means that any "from" part in the target section that wouldn't take input are discarded, and any
// missing sections where the target isn't covering anything will have a 1-1 section added to it.
func overlaps(sourceSections []section, targetSections []section) []section {
	// let's get the from-to of the source sections:
	sourcePairs := make(map[int]int)
	sourceOrder := make([]int, len(sourceSections))
	for i, s := range sourceSections {
		sourceOrder[i] = s.destinationStart
		sourcePairs[s.destinationStart] = s.destinationEnd
	}

	sort.Ints(sourceOrder)

	start, end := 0, 0

	overlapSections := make([]section, 0)

	for _, sourceStart := range sourceOrder {
		start = sourceStart
		end = sourcePairs[sourceStart]

		//fmt.Printf("- outer: start %d / end %d, stepvalue set to start\n", start, end)
		stepValue := start

		i := 0

	fragmentLoop:
		for {
			i++

			covers, next, delta := doesSectionCoverThisValue(stepValue, targetSections)
			//fmt.Printf("\n  :: inner iter %d: stepvalue %d, covers next delta are %t, %d, %d\n", i, stepValue, covers, next, delta)

			switch covers {
			case false:
				//fmt.Printf("  :: inner::: iter %d: covers false\n", i)

				switch {
				case next == 0:
					fallthrough
				case next > end:
					//fmt.Printf("  :: inner iter %d: no cover, next is either 0, or less than the end (next/end: %d/%d)\n", i, next, end)
					// does not cover it, and the next one is past the end of the current one or
					// does not cover it and there are no more sections in the receiving end to deal with.
					// most straightforward, we need to add a 1:1 section that covers from stepValue to end
					overlapSections = append(overlapSections, section{
						sourceStart:      stepValue,
						sourceEnd:        end,
						destinationStart: stepValue,
						destinationEnd:   end,
						delta:            0,
					})

					//fmt.Printf("  :: inner iter %d: added section source / end / delta: %d / %d / %d\n", i, stepValue, end, 0)

					//fmt.Printf("  :: inner iter %d: breaking fragmentloop, it should start a new start/end iteration\n\n", i)
					// once we're here, we need to tick the sending sections over
					break fragmentLoop
				case next <= end:
					// does not cover it, but the next is within the current section
					// create a section from current stepValue to next
					overlapSections = append(overlapSections, section{
						sourceStart:      stepValue,
						sourceEnd:        next - 1,
						destinationStart: stepValue,
						destinationEnd:   next - 1,
						delta:            0,
					})

					//fmt.Printf("  :: inner iter %d: added section source / end / delta: %d / %d / %d\n", i, stepValue, next-1, delta)

					stepValue = next

					//fmt.Printf("  :: inner iter %d: not cover it, next is smaller than or equal to end, so setting step value "+
					//	"to be next + 1 (stepvalue, next: %d, %d), but not breaking!\n", i, stepValue, next)
				default:
					panic("you dun effed up")
				}
			case true:
				//fmt.Printf("  :: inner::: iter %d: covers true\n", i)

				switch {
				case next <= end:
					// covers it, and the next one is either at, or before the end one, which means we need a section
					// the current stepValue to the next value and set stepValue to be next +1 and loop
					overlapSections = append(overlapSections, section{
						sourceStart:      stepValue,
						sourceEnd:        next,
						destinationStart: stepValue + delta,
						destinationEnd:   next + delta,
						delta:            delta,
					})

					//fmt.Printf("  :: inner iter %d: added section source / end / delta: %d / %d / %d\n", i, stepValue, next, delta)

					stepValue = next + 1
					//fmt.Printf("  :: inner iter %d: covers, next is smaller equal to end (%d/%d), setting stepvalue to be "+
					//	"next + 1: %d/%d\n", i, next, end, stepValue, next)
				default:
					// covers it, but next is bigger than end, so we need a section to cover until the end, and then
					// break
					overlapSections = append(overlapSections, section{
						sourceStart:      stepValue,
						sourceEnd:        end,
						destinationStart: stepValue + delta,
						destinationEnd:   end + delta,
						delta:            delta,
					})

					//fmt.Printf("  :: inner iter %d: added section source / end / delta: %d / %d / %d\n", i, stepValue, end, delta)

					//fmt.Printf("  :: inner iter %d: covers, next %d is bigger than end %d, so breaking the inner loop"+
					//	" which should start a new start/end loop\n\n", i, next, end)

					break fragmentLoop
				}
			}

			//fmt.Printf("  :: inner iter %d: outside of switch, uh...\n", i)

			if i > 10 {
				break fragmentLoop
			}

			if stepValue > end {
				//fmt.Printf("  :: inner iter %d: stepvalue %d is bigger than end %d, breaking\n", i, stepValue, end)
				break fragmentLoop
			}
		}
	}

	return overlapSections
}

// doesSectionCoverThisValue returns whether a given value is covered by any of the section ranges. The possible return
// values are:
//   - true, some integer > value, which means that value IS covered, and here's the last value that's covered in that
//     section
//   - false, some integer > value, which means the value is NOT covered, but the next section that DOES cover it begins
//     with the returned integer, or
//   - false, 0, which means the value is NOT covered, and will NOT be covered at all
func doesSectionCoverThisValue(value int, sections []section) (covers bool, next int, delta int) {
	m := make(map[int]section)
	order := make([]int, len(sections))

	for i, section := range sections {
		m[section.sourceStart] = section
		order[i] = section.sourceStart
	}
	sort.Ints(order)

	if len(order) == 0 {
		return false, 0, 0
	}

	if value < order[0] {
		return false, m[order[0]].sourceStart, 0
	}

	lastReturn := 0
	delta = 0

	for i, start := range order {
		if value >= start {
			if value <= m[start].sourceEnd {
				// covers it, should return delta as well
				return true, m[start].sourceEnd, m[start].delta
			}

			// does not cover it
			if i < len(order)-1 {
				// there is a new section
				lastReturn = m[order[i+1]].sourceStart
			} else {
				lastReturn = 0
			}
		} else {
			return false, lastReturn, 0
		}
	}

	return false, lastReturn, 0
}
