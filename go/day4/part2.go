package day4

import (
	"os"

	"github.com/javorszky/adventofcode2023/inputs"
	"github.com/rs/zerolog"
)

func Task2(l zerolog.Logger) {
	localLogger := l.With().Int("day", 4).Int("part", 2).Logger()

	gog, err := inputs.ReadIntoLines("day4/input1.txt")
	if err != nil {
		localLogger.Err(err).Msg("could not read input file")
		os.Exit(1)
	}

	data, err := parse(gog)
	if err != nil {
		localLogger.Err(err).Msgf("parsing incoming lines failed")
	}

	// populate the numberOfCopies map to have 1 copy of everything at the start
	numberOfCopies := make(map[int]int)
	for i := 1; i <= len(data); i++ {
		numberOfCopies[i] = 1
	}

	maxLen := len(data)

	// loop through the cards in order. Need to use this for because the order of looping through a map is undefined.
	// In this case the keys of the map correspond to the integers of the cards. Could this have been an array?
	// Probably, but then I'd have needed to write another different parse function and I didn't want to.
	for i := 1; i <= len(data); i++ {
		// find how many winning numbers the card had. This will inform how many cards below this one will have copies
		// added to them.
		found := findUnions(data[i][0], data[i][1])

		// how many copies of this card do we already have? This will be used to increment the cards under this.
		haveCards, ok := numberOfCopies[i]
		if !ok {
			localLogger.Fatal().Msgf("numberOfCopies map should have had a record of card %d, but did not", i)
		}

		// if there are no unions, then skip to the next one, we're not going to increment anything.
		if found == 0 {
			continue
		}

		// add the copies to the cards below this
		for j := 0; j < found; j++ {
			next := i + j + 1
			// if we were to increment a card that is past the list, don't
			if next > maxLen {
				continue
			}

			// increment the number of cards that we have under this by the number of copies of the current card.
			numberOfCopies[next] += haveCards
		}
	}

	sum := 0
	for _, v := range numberOfCopies {
		sum += v
	}

	solution := sum
	s := localLogger.With().Int("solution", solution).Logger()
	s.Info().Msgf("After the new rules we are left with a total of %d cards", solution)
}
