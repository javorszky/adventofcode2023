package cmd

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/javorszky/adventofcode2023/day1"
	"github.com/javorszky/adventofcode2023/day2"
	"github.com/javorszky/adventofcode2023/day3"
	"github.com/javorszky/adventofcode2023/day4"
	"github.com/javorszky/adventofcode2023/day5"
	"github.com/javorszky/adventofcode2023/day6"
	"github.com/javorszky/adventofcode2023/day7"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		l := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Str("module", "adventofcode").Int("year", 2023).Logger()
		l.Info().Msg("Welcome to Gabor Javorszky's Advent of Code 2023 solutions!")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		tasks := map[int][2]func(logger zerolog.Logger){
			1: {day1.Task1, day1.Task2},
			2: {day2.Task1, day2.Task2},
			3: {day3.Task1, day3.Task2},
			4: {day4.Task1, day4.Task2},
			5: {day5.Task1, day5.Task2},
			6: {day6.Task1, day6.Task2},
			7: {day7.Task1, day7.Task2},
		}

		lenT := len(tasks)
		day := 0
		parts := 2

		if len(args) > 0 {
			d, err := parseDay(args[0], lenT)
			if err != nil {
				l.Err(err).Msgf("parsing day argument '%s'", args[0])
			} else {
				day = d
			}
		}

		if len(args) > 1 {
			p, err := parsePart(args[1])
			if err != nil {
				l.Err(err).Msgf("parsing parts argument '%s'", args[1])
			} else {
				parts = p
			}
		}

		switch {
		case day > 0 && day <= lenT:
			if parts > 1 {
				tasks[day][0](l)
				tasks[day][1](l)
				return
			}

			tasks[day][parts](l)
		default:
			for i := 1; i <= len(tasks); i++ {
				tasks[i][0](l)
				tasks[i][1](l)
			}
		}
	},
}

func parseDay(arg string, max int) (int, error) {
	d, err := strconv.Atoi(arg)
	if err != nil {
		return 0, errors.Wrapf(err, "converting '%s' first argument into an integer. This should be a number between 1-31", arg)
	}

	if d < 0 || d > max {
		return 0, fmt.Errorf("can't call specified day. You wanted %d. Today it needs to be a number between 1 and %d, including those two", d, max)
	}

	return d, nil
}

func parsePart(arg string) (int, error) {
	d, err := strconv.Atoi(arg)
	if err != nil {
		return 0, errors.Wrapf(err, "converting '%s' first argument into an integer. This should be either the number 1, or the number 2", arg)
	}

	if d != 1 && d != 2 {
		return 0, fmt.Errorf("can't call specified part. You wanted %d. This needs to be either number 1 or number 2", d)
	}

	return d - 1, nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}

type NoDebugLog struct{}

func (n NoDebugLog) Run(_ *zerolog.Event, _ zerolog.Level, _ string) {
	// do nothing
}
