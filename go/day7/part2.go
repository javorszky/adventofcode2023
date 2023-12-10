package day7

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

var cardMapPart2 = map[string]int{
	"J": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"T": 10,
	"Q": 12,
	"K": 13,
	"A": 14,
}

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 7).Int("part", 2).Logger()

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
		class = classifyHandPart2(k)
		classification[class] = append(classification[class], k)
	}

	globalOrder := make([]string, 0)

	for _, currentClass := range tierOrder {
		localSlice := sortHandsPart2(classification[currentClass])
		globalOrder = append(globalOrder, localSlice...)
	}

	sum := 0

	for i, card := range globalOrder {
		sum += (i + 1) * hands[card]
	}

	solution := sum

	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("With the J rule, the total winnings are %d", solution)
}

func classifyHandPart2(hand string) handType {
	class := classifyHand(hand)

	jeez := strings.Count(hand, "J")

	if jeez == 0 || jeez == 5 {
		// no J in the hand, classify as normal
		return class
	}

	switch class {
	case fourOfAKind:
		// four of a kind and has a J would look either
		// QQQQJ or JJJJQ, and both cases they can become a five of a kind
		fallthrough
	case fullHouse:
		// full house with a J in it look either like
		// QQQJJ or QQJJJ, and both cases can become a five of a kind
		return fiveOfAKind
	case threeOfAKind:
		// this has not matched any of the above, so the three of a kind looks either of these:
		// QQQJA or JJJQA, and both of them can become a four of a kind
		return fourOfAKind
	case twoPair:
		// if one of the pairs of the two pair are the Js, then yeah
		// QQJJA -> four of a kind, but
		// QQAAJ -> full house
		switch jeez {
		case 1:
			return fullHouse
		case 2:
			return fourOfAKind
		default:
			panic(fmt.Sprintf("this should never have happened, twopair: jeez: %d, hand %s", jeez, hand))
		}
	case onePair:
		// not matched any of the above, there are two possibilities:
		// JJ234, in which case the 2 Js can turn into whatever else, or
		// 2234J, in which case the single J can turn into whatever the pair is
		return threeOfAKind
	case highCard:
		// not matched anything above, there's strictly 1 J present (we know because it hasn't matched
		// anything above), that J can take up whatever other form
		return onePair
	default:
		panic(fmt.Sprintf("this should never have happened, the class is %s, the hand is %s", class, hand))
	}
}

// This will order the hands in ascending order, which means the weakest hand is going to be at the
// beginning of the resulting slice.
func sortHandsPart2(hands []string) []string {
	slices.SortStableFunc(hands, cardSortFuncPart2)

	return hands
}

// cmp(a, b) should return
// * a negative number when a < b
// * a positive number when a > b
// * and zero when a == b.
// This will assume that the strings are always going to be 5 characters long.
func cardSortFuncPart2(a, b string) int {
	//fmt.Printf("a = %s, b = %s\n", a, b)
	for i := 0; i < 5; i++ {
		ca := cardMapPart2[string(a[i])]
		cb := cardMapPart2[string(b[i])]

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

	fmt.Printf("everything was same, returning 0\n")
	return 0
}
