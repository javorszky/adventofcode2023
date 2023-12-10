package day8

import (
	"fmt"
	"os"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

type cycleData struct {
	node string
	step int
	nth  int
	zeds map[string][]int
}

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 8).Int("part", 2).Logger()

	gog, err := inputs.ReadIntoLines("day8/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	instructions := parseInstructions(gog[0])
	dir := pickDirection2(instructions)

	nodes := parseNodes(gog)
	starts := filterStarts(nodes)
	cycles := make(map[string]cycleData)

	for k := range starts {
		//encounteredZFirst := 0

		records := make(map[string]map[int][]int)
		i := 0
		currentNode := k
		var direction instruction
		var n int

		zeds := make(map[string][]int)

		for {
			direction, n = dir(i)
			currentNode = nodes[currentNode][direction]

			if records[currentNode] == nil {
				records[currentNode] = make(map[int][]int)
			}

			if strings.HasSuffix(currentNode, "Z") {
				zeds[currentNode] = append(zeds[currentNode], i)
				if len(zeds[currentNode]) > 1 {
					break
				}
			}
			//
			//records[currentNode][n] = append(records[currentNode][n], i)
			//if len(records[currentNode][n]) > 1 {
			//	fmt.Printf("we have visited the same node from the same start twice on the same step at these two iterations: %v\n", records[currentNode][n])
			//	break
			//}

			i++
			if i > 60000 {
				fmt.Printf("does not cycle in 30k iterations\n")
				break
			}
		}

		cycles[k] = cycleData{
			node: currentNode,
			step: i,
			nth:  n,
			zeds: zeds,
		}
	}

	cyclesInSlice := make([]uint64, 0)

	for _, v := range cycles {
		for _, w := range v.zeds {
			cyclesInSlice = append(cyclesInSlice, uint64(w[1]-w[0]))
		}
	}

	result := leastCommonMultiple(cyclesInSlice[0], cyclesInSlice[1], cyclesInSlice[2:]...)
	//result := uint64(1)
	solution := result
	s := localLogger.With().Uint64("solution", solution).Logger()
	s.Info().Msgf("They will all hit it at %d", solution)
}

func filterStarts(nodes map[string]map[instruction]string) map[string]map[instruction]string {
	start := make(map[string]map[instruction]string)

	for k, v := range nodes {
		if strings.HasSuffix(k, "A") {
			start[k] = v
		}
	}

	return start
}

func pickDirection2(list []instruction) func(int) (instruction, int) {
	return func(step int) (instruction, int) {
		n := step % len(list)
		return list[n], n
	}
}

func greatestCommonDivisor(a, b uint64) uint64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}

	return a
}

// leastCommonMultiple finds the Least Common Multiple via the greatestCommonDivisor.
func leastCommonMultiple(a, b uint64, integers ...uint64) uint64 {
	result := a * b / greatestCommonDivisor(a, b)

	for i := 0; i < len(integers); i++ {
		result = leastCommonMultiple(result, integers[i])
	}

	return result
}
