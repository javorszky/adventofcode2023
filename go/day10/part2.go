package day10

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const (
	pipeOutside = " "
	pipeInside  = "o"
)

var ErrOnEdge = errors.New("can't determine inside/outside, direction is on edge")

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 10).Int("part", 2).Logger()

	/**
	Read the input into lines
	*/
	gog, err := inputs.ReadIntoLines("day10/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	/**
	Construct a new pipe map that only has the layout, start, and its position.
	*/
	pMap := newPipeMap(gog)

	/**
	Figure out the first connection on the left of the start node. We'll need this later.
	*/
	var currentNode *node
	var lastMove direction

	for _, d := range []direction{west, south, east, north} {
		newNode, err := pMap.peek(d)
		if err != nil {
			localLogger.Err(err).Msgf("tried to peek into a direction of %s\n", d)
			continue
		}
		if newNode.nodeType == pipeNone {
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
				localLogger.Err(err).Msgf("moving on pmap did not succeed, direction was %s", d)
			}
			break
		}
	}

	/**
	Start traversing the pipes from start, going left, and then moving the same direction every time.
	While there, also find one of each edge tile.
	*/
	elements := 0
	var northernMostElement *node
	northernMostRow := len(gog)

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

		if nextNode.coord.row < northernMostRow {
			northernMostRow = nextNode.coord.row
			northernMostElement = nextNode
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

	/**
	Create a map of the loop, one that we can use to say "nah, do not turn these tiles into a normalized
	ground tile for the purposes of finding infills.
	*/
	loop := make(map[coordinate]*node)
	currentNode = pMap.start
	loop[currentNode.coord] = currentNode

	for {
		currentNode = currentNode.left()
		if currentNode.nodeType == pipeStart {
			break
		}
		loop[currentNode.coord] = currentNode
	}

	/**
	Normalize the map, as in turn every piece of non-connected pipe into a ground field: .
	*/
	for currentRow := 0; currentRow < pMap.rows; currentRow++ {
		for currentCol := 0; currentCol < pMap.cols; currentCol++ {
			c := coordinate{
				row: currentRow,
				col: currentCol,
			}

			if _, ok := loop[c]; ok {
				continue
			}
			pMap.normalize(c)
		}
	}

	/**
	Draw a pretty map of the normalised state, with the start
	*/
	//visualize(pMap)

	/**
	Find what the original pipe the start tile is by looking at where its left and right nodes are in relation
	to it, and then figuring out what shape that corresponds to. Then update the start node and the layout in
	the map to that new shape.

	This is needed so when we traverse the loop and keep the outside / inside on one side, it's not going to
	break with the S tile.
	*/
	leftFromStart := pMap.start.left()
	rightFromStart := pMap.start.right()
	d1, err := dir(pMap.start.coord, leftFromStart.coord)
	if err != nil {
		localLogger.Err(err).Msgf("checking direction from start to left of start")
	}

	d2, err := dir(pMap.start.coord, rightFromStart.coord)
	if err != nil {
		localLogger.Err(err).Msgf("checking direction from start to right of start")
	}
	shape := shapeForDirections(d1, d2)
	//fmt.Printf("the shape of start is actually %s\n", shape)
	pMap.start.nodeType = shape
	pMap.flipStartTo(shape)

	/**
	Draw the map again, this time with the start tile replaced.
	*/
	//fmt.Printf("start is at %#v\n", pMap.start.coord)
	//visualize(pMap)

	/**
	Let's flood fill the outside.
	*/
	visited := make(map[string]struct{})
	var f func(*node)
	directionSlice := []direction{north, east, south, west}

	didTheObvious := false

	f = func(n *node) {
		if n.nodeType != pipeNone {
			return
		}

		pMap.current = n.coord
		_, ok := visited[n.coord.String()]
		if ok {
			return
		}

		if didTheObvious {
			fmt.Printf("clearing the outside, but within the loop at %s\n", n.coord.String())
		}

		visited[n.coord.String()] = struct{}{}
		pMap.updateFieldTo(n.coord, pipeOutside)

		for _, d := range directionSlice {
			pMap.current = n.coord
			next, err := pMap.peek(d)
			if err != nil {
				continue
				// silently do nothing
			}

			f(next)
		}
	}

	f(pMap.getNode(coordinate{
		row: 0,
		col: 0,
	}))

	//visualize(pMap)

	//didTheObvious = true

	/**
	Let's start with the inside/outside thing.
	*/
	currentNode = northernMostElement
	previousNode := currentNode
	outside1 := north
	outside2, err := otherOutside(currentNode, north)
	if err != nil {
		localLogger.Err(err).Msgf("getting other outside for node %v from north", *currentNode)
	}

	var keep1, keep2 direction

	for {
		previousNode = currentNode
		currentNode = currentNode.left()

		if currentNode == northernMostElement {
			// we've come full circle, break
			break
		}

		directionWeWent, err := dir(previousNode.coord, currentNode.coord)
		if err != nil {
			localLogger.Err(err).Msgf("tried to figure out the direction between nodes %v and %v", *previousNode, *currentNode)
		}

		un1, un2 := unaffectedOutside(directionWeWent)
		//fmt.Printf("node %s: outside 1 %s and outside 2 %s\n"+
		//	"we went %s, new node is type %s\n"+
		//	"unaffected 1 and 2 are %s, %s\n", previousNode.nodeType, outside1, outside2,
		//	directionWeWent, currentNode.nodeType,
		//	un1, un2)

		if outside1 == un1 || outside1 == un2 {
			keep1 = outside1
		} else {
			keep1 = outside2
		}

		keep2, err = otherOutside(currentNode, keep1)
		if err != nil {
			localLogger.Err(err).Msgf("otherOutside for node %v and direction %s", *currentNode, keep1)
		}

		outside1, outside2 = keep1, keep2

		pMap.current = currentNode.coord

		outsideNode1, err := pMap.peek(outside1)
		if err == nil {
			if outsideNode1.nodeType == pipeNone {
				f(outsideNode1)
			}
		}

		if outside1 != outside2 {
			outsideNode2, err := pMap.peek(outside2)
			if err == nil {
				if outsideNode2.nodeType == pipeNone {
					f(outsideNode2)
				}
			}
		}
	}

	//fmt.Printf("\nOkay, apparently we've also done the outsides, but on the inside\n")
	//visualize(pMap)

	insides := strings.Count(pMap.layout, pipeNone)

	solution := insides
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Inside of the loop there are %d non-pipe tiles", solution)
}

func visualize(m *pipeMap) {
	var sb strings.Builder
	i := 0
	for row := 0; row < m.rows; row++ {
		for col := 0; col < m.cols; col++ {
			sb.WriteByte(m.layout[i])
			i++
		}
		sb.WriteString("\n")
	}

	fmt.Printf(sb.String())
}

func otherOutside(n *node, d direction) (direction, error) {
	switch n.nodeType {
	case pipeHorizontal:
		switch d {
		case north:
			return north, nil
		case east:
			return "", ErrOnEdge
		case south:
			return south, nil
		case west:
			return "", ErrOnEdge
		}
	case pipeVertical:
		switch d {
		case north:
			return "", ErrOnEdge
		case east:
			return east, nil
		case south:
			return "", ErrOnEdge
		case west:
			return west, nil
		}
	case pipeNorthEast:
		switch d {
		case north:
			return east, nil
		case east:
			return north, nil
		case south:
			return west, nil
		case west:
			return south, nil
		}
	case pipeNorthWest:
		switch d {
		case north:
			return west, nil
		case east:
			return south, nil
		case south:
			return east, nil
		case west:
			return north, nil
		}
	case pipeSouthEast:
		switch d {
		case north:
			return west, nil
		case east:
			return south, nil
		case south:
			return east, nil
		case west:
			return north, nil
		}
	case pipeSouthWest:
		switch d {
		case north:
			return east, nil
		case east:
			return north, nil
		case south:
			return west, nil
		case west:
			return south, nil
		}
	}

	panic(fmt.Sprintf("this should never have happened for direction %s and node type %s\n", d, n.nodeType))
}

func dir(previous, next coordinate) (direction, error) {
	dCol := previous.col - next.col
	dRow := previous.row - next.row

	if dCol < -1 || dCol > 1 || dRow < -1 || dRow > 1 || dRow*dCol != 0 {
		return "", fmt.Errorf("two coordinates are not one tile away: %v and %v", previous, next)
	}

	switch {
	case dCol == -1:
		return east, nil
	case dCol == 1:
		return west, nil
	case dRow == -1:
		return south, nil
	case dRow == 1:
		return north, nil
	}

	panic(fmt.Sprintf("direction between two coordinates, unhandled! %v and %v", previous, next))
}

func unaffectedOutside(d direction) (direction, direction) {
	switch d {
	case north:
		fallthrough
	case south:
		return east, west
	default:
		return north, south
	}
}

func shapeForDirections(d1, d2 direction) string {
	switch d1 {
	case north:
		switch d2 {
		case east:
			return pipeNorthEast
		case south:
			return pipeVertical
		case west:
			return pipeNorthWest
		}
	case east:
		switch d2 {
		case north:
			return pipeNorthEast
		case south:
			return pipeSouthEast
		case west:
			return pipeHorizontal
		}
	case south:
		switch d2 {
		case north:
			return pipeVertical
		case east:
			return pipeSouthEast
		case west:
			return pipeSouthWest
		}
	case west:
		switch d2 {
		case north:
			return pipeNorthWest
		case east:
			return pipeHorizontal
		case south:
			return pipeSouthWest
		}
	}

	panic(fmt.Sprintf("shapeForDirections, inputs %s and %s", d1, d2))
}
