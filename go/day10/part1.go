package day10

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const (
	pipeNone       = "."
	pipeVertical   = "|"
	pipeHorizontal = "-"
	pipeNorthEast  = "L"
	pipeNorthWest  = "J"
	pipeSouthWest  = "7"
	pipeSouthEast  = "F"
	pipeStart      = "S"

	north direction = "n"
	east  direction = "e"
	south direction = "s"
	west  direction = "w"
)

type direction string

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 10).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day10/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	pMap := newPipeMap(gog)

	var currentNode *node
	var lastMove direction

	for _, d := range []direction{west, south, east, north} {
		newNode, err := pMap.peek(d)
		if err != nil {
			localLogger.Err(err).Msgf("tried to peek into a direction of %s\n", d)
			continue
		}

		connections := connects(newNode)

		if op[d] == connections[0] || op[d] == connections[1] {
			// new node is good, let's fill out the details, attach the address of the start to the
			// right of new node
			newNode.rightNode = pMap.start

			// link that new node's address to the left of start
			pMap.start.leftNode = newNode

			// set the current node to the new node and the current move as the last move
			currentNode = newNode
			lastMove = d

			// actually move in the map
			err = pMap.move(d)
			if err != nil {
				localLogger.Err(err).Msgf("moving on pmap did not succeed, direction was %d", d)
			}
			break
		}
	}

	// we need to go leftNode
	elements := 0
	for {
		dirs := connects(currentNode)
		if op[lastMove] == dirs[0] {
			lastMove = dirs[1]
		} else {
			lastMove = dirs[0]
		}

		nextNode, err := pMap.peek(lastMove)
		if err != nil {
			localLogger.Err(err).Msgf("pMap.peek for the last move from %v", pMap.current)
		}

		if nextNode.nodeType == pipeStart {
			nextNode = pMap.start
		}

		err = pMap.move(lastMove)
		if err != nil {
			localLogger.Err(err).Msgf("actually moving was bad")
		}
		currentNode.leftNode = nextNode
		nextNode.rightNode = currentNode

		currentNode = nextNode

		if currentNode.nodeType == pipeStart {
			// we have arrived back to the start
			break
		}
		elements++
	}

	// okay, we have a ring at this point
	currentLeft := pMap.start.left()
	currentLeft.distanceFromS = 1

	currentRight := pMap.start.right()
	currentRight.distanceFromS = 1

	previousLeft, previousRight := pMap.start, pMap.start

	biggestDistance := 1
	distance := 1
	for {
		distance++
		biggestDistance = distance
		previousLeft = currentLeft
		currentLeft = currentLeft.left()

		previousRight = currentRight
		currentRight = currentRight.right()

		if currentLeft == currentRight {
			currentLeft.distanceFromS = distance
			fmt.Printf("we have found the same node coming from both directions!\n")
			break
		}

		if previousRight == currentLeft || previousLeft == currentRight {
			fmt.Printf("we have found an overlap, breaking without updating distance!\n")
			break
		}

		currentLeft.distanceFromS = distance
		currentRight.distanceFromS = distance
	}

	solution := biggestDistance
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Biggest distance is %d", solution)
}

type coordinate struct {
	row, col int
}

func (c coordinate) String() string {
	return fmt.Sprintf("%d-%d", c.row, c.col)
}

type pipeMap struct {
	layout     string
	rows, cols int
	start      *node
	current    coordinate
	lock       *sync.RWMutex
}

func (p *pipeMap) move(d direction) error {
	var (
		newRow int
		newCol int
	)

	switch d {
	case north:
		if p.current.row == 0 {
			return errors.New("we're on the top edge, can't peek north")
		}

		newRow = p.current.row - 1
		newCol = p.current.col
	case east:
		if p.current.col == p.cols-1 {
			return errors.New("we're on the rightNode edge, can't peek east")
		}

		newRow = p.current.row
		newCol = p.current.col + 1
	case south:
		if p.current.row == p.rows-1 {
			return errors.New("we're on the bottom edge, can't peek south")
		}

		newRow = p.current.row + 1
		newCol = p.current.col
	case west:
		if p.current.col == 0 {
			return errors.New("we're on the leftNode edge, can't peek west")
		}

		newRow = p.current.row
		newCol = p.current.col - 1
	default:
		panic("this should not have happened!")
	}

	p.current = coordinate{
		row: newRow,
		col: newCol,
	}

	return nil
}

func (p *pipeMap) normalize(c coordinate) {
	p.updateFieldTo(c, pipeNone)
}

func (p *pipeMap) flipStartTo(shape string) {
	si := strings.Index(p.layout, "S")

	p.layout = p.layout[:si] + shape + p.layout[si+1:]
}

func (p *pipeMap) updateFieldTo(c coordinate, to string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	i := c.row*p.cols + c.col

	p.layout = p.layout[:i] + to + p.layout[i+1:]
}

func (p *pipeMap) peek(d direction) (*node, error) {
	var (
		newRow int
		newCol int
	)

	switch d {
	case north:
		if p.current.row == 0 {
			return nil, errors.New("we're on the top edge, can't peek north")
		}

		newRow = p.current.row - 1
		newCol = p.current.col
	case east:
		if p.current.col == p.cols-1 {
			return nil, errors.New("we're on the rightNode edge, can't peek east")
		}

		newRow = p.current.row
		newCol = p.current.col + 1
	case south:
		if p.current.row == p.rows-1 {
			return nil, errors.New("we're on the bottom edge, can't peek south")
		}

		newRow = p.current.row + 1
		newCol = p.current.col
	case west:
		if p.current.col == 0 {
			return nil, errors.New("we're on the leftNode edge, can't peek west")
		}

		newRow = p.current.row
		newCol = p.current.col - 1
	default:
		panic("this should not have happened!")
	}

	return p.getNode(coordinate{
		row: newRow,
		col: newCol,
	}), nil
}

func (p *pipeMap) getNode(c coordinate) *node {
	i := c.row*p.cols + c.col
	return &node{
		nodeType:      string(p.layout[i]),
		distanceFromS: 0,
		coord: coordinate{
			row: c.row,
			col: c.col,
		},
		leftNode:  nil,
		rightNode: nil,
	}
}

func newPipeMap(gog []string) *pipeMap {
	rows := len(gog)
	cols := len(gog[0])

	layout := strings.Join(gog, "")
	i := strings.Index(layout, pipeStart)

	// row and col are 0-indexed
	sRow := i / cols
	sCol := i - (sRow * cols)

	p := pipeMap{
		layout: layout,
		rows:   rows,
		cols:   cols,
		start: &node{
			nodeType:      pipeStart,
			distanceFromS: 0,
			coord: coordinate{
				row: sRow,
				col: sCol,
			},
			leftNode:  nil,
			rightNode: nil,
		},
		current: coordinate{
			row: sRow,
			col: sCol,
		},
		lock: new(sync.RWMutex),
	}

	return &p
}

type node struct {
	nodeType            string
	distanceFromS       int
	coord               coordinate
	outside, inside     [2]direction
	leftNode, rightNode *node
}

func (n *node) left() *node {
	return n.leftNode
}

func (n *node) right() *node {
	return n.rightNode
}

func connects(n *node) []direction {
	switch n.nodeType {
	case pipeVertical:
		return []direction{north, south}
	case pipeHorizontal:
		return []direction{east, west}
	case pipeNorthEast:
		return []direction{north, east}
	case pipeNorthWest:
		return []direction{north, west}
	case pipeSouthWest:
		return []direction{south, west}
	case pipeSouthEast:
		return []direction{south, east}
	default:
		return nil
	}
}

var op = map[direction]direction{
	east:  west,
	west:  east,
	north: south,
	south: north,
}
