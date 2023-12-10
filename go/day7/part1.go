package day7

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

const (
	fiveOfAKind  handType = "five"
	fourOfAKind  handType = "four"
	fullHouse    handType = "fullHouse"
	threeOfAKind handType = "three"
	twoPair      handType = "two"
	onePair      handType = "one"
	highCard     handType = "highCard"
)

var tierOrder = []handType{
	highCard,
	onePair,
	twoPair,
	threeOfAKind,
	fullHouse,
	fourOfAKind,
	fiveOfAKind,
}
var cardMap = map[string]int{
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"T": 10,
	"J": 11,
	"Q": 12,
	"K": 13,
	"A": 14,
}

type handType string

func Task1(l zerolog.Logger) {
	localLogger := l.With().Int("day", 7).Int("part", 1).Logger()

	gog, err := inputs.ReadIntoLines("day7/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	hands, err := parseHands(gog)
	if err != nil {
		localLogger.Err(err).Msg("parseHands")
	}

	classification := map[handType][]string{
		fiveOfAKind:  {},
		fourOfAKind:  {},
		fullHouse:    {},
		threeOfAKind: {},
		twoPair:      {},
		onePair:      {},
		highCard:     {},
	}

	var class handType

	for k := range hands {
		class = classifyHand(k)
		classification[class] = append(classification[class], k)
	}

	globalOrder := make([]string, 0)

	for _, currentClass := range tierOrder {
		localSlice := sortHands(classification[currentClass])

		fmt.Printf("okay, so apparently sorted %s\n%v\n\n", currentClass, localSlice)
		globalOrder = append(globalOrder, localSlice...)
	}

	sum := 0

	for i, card := range globalOrder {
		sum += (i + 1) * hands[card]
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("Sum of all the bids multiplied by the hands' ranks is %d", solution)
}

func parseHands(gog []string) (map[string]int, error) {
	m := make(map[string]int)

	for _, line := range gog {
		parts := strings.Split(line, " ")

		bid, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, errors.Wrapf(err, "strconv.Atoi: %s", parts[1])
		}

		m[parts[0]] = bid
	}

	return m, nil
}

func classifyHand(hand string) handType {
	cards := make(map[string]int)

	for _, c := range hand {
		cards[string(c)]++
	}

	switch len(cards) {
	case 1:
		return fiveOfAKind
	case 5:
		return highCard
	case 4:
		return onePair
	case 2:
		// this could be either a full house, or four of a kind
		p := 1
		for _, v := range cards {
			p *= v
		}

		switch p {
		case 4:
			return fourOfAKind
		case 6:
			return fullHouse
		default:
			panic(fmt.Sprintf("len 2: this should not have happened with card: %v and p %d", cards, p))
		}
	case 3:
		// this could be three of a kind, or two pairs
		p := 1

		for _, v := range cards {
			p *= v
		}
		switch p {
		case 3:
			return threeOfAKind
		case 4:
			return twoPair
		default:
			panic(fmt.Sprintf("len 3: this should not have happened with card %v and p %d", cards, p))
		}
	default:
		panic(fmt.Sprintf("len weird: cards %v and len %d", cards, len(cards)))
	}
}

// This will order the hands in ascending order, which means the weakest hand is going to be at the
// beginning of the resulting slice.
func sortHands(hands []string) []string {
	slices.SortFunc(hands, cardSortFunc)

	return hands
}

// cmp(a, b) should return
// * a negative number when a < b
// * a positive number when a > b
// * and zero when a == b.
// This will assume that the strings are always going to be 5 characters long.
func cardSortFunc(a, b string) int {
	//fmt.Printf("a = %s, b = %s\n", a, b)
	for i := 0; i < 5; i++ {
		ca := cardMap[string(a[i])]
		cb := cardMap[string(b[i])]

		//fmt.Printf("character %d\n"+
		//	"a[i] = %s => ca = %d\n"+
		//	"b[i] = %s => cb = %d\n",
		//	i, string(a[i]), ca, string(b[i]), cb)

		switch {
		case ca == cb:
			//fmt.Printf("character %d is same, continue to next character\n", i)
			continue
		case ca < cb:
			//fmt.Printf("character %d, a is lower, returning -1\n", i)
			return -1
		case cb < ca:
			//fmt.Printf("character %d, b is lower, returning 1\n", i)
			return 1
		}
	}

	//fmt.Printf("everything was same, returning 0\n")
	return 0
}
