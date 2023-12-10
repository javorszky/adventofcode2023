package day8

import (
	"fmt"
	"os"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const (
	left      instruction = "left"
	right     instruction = "right"
	startNode string      = "AAA"
	endNode   string      = "ZZZ"
)

type instruction string

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 8).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day8/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	instructions := parseInstructions(gog[0])

	nodes := parseNodes(gog)

	i := 0
	currentNode := startNode
	visited := make(map[string][]int)

	for {
		visited[currentNode] = append(visited[currentNode], i)

		currentNode = nodes[currentNode][pickDirection(i, instructions)]

		i++

		if currentNode == endNode {
			break
		}
	}

	solution := i
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("It took %d steps to reach ZZZ", solution)
}

func parseInstructions(line string) []instruction {
	ins := make([]instruction, len(line))

	for i, ch := range line {
		switch string(ch) {
		case "L":
			ins[i] = left
		case "R":
			ins[i] = right
		default:
			panic(fmt.Sprintf("this should never have happened: ch: %s T %D", string(ch), i))
		}
	}

	return ins
}

func parseNodes(lines []string) map[string]map[instruction]string {
	nodes := make(map[string]map[instruction]string)

	for i := 2; i < len(lines); i++ {
		from := lines[i][0:3]
		leftNode := lines[i][7:10]
		rightNode := lines[i][12:15]

		nodes[from] = map[instruction]string{
			left:  leftNode,
			right: rightNode,
		}
	}

	return nodes
}

func pickDirection(step int, list []instruction) instruction {
	n := step % len(list)
	return list[n]
}
