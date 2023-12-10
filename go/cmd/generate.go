package cmd

/*
Copyright © 2023 Gabor Javorszky <gabor@javorszky.co.uk>
*/

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type dayData struct {
	Pkg string
	Day int
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a folder with files with correct package names",
	Long: `Takes the argument or the flag value and generates a new folder for that day adjusting the package names as
needed. If both an argument and a flag are supplied, the flag gets higher priority.

Usage:
$ aoc23 generate 21
// generates folder 'day21' with package name 'day21' set

$ aoc23 generate --day 19
// generates folder 'day19', same as above

$ aoc23 generate 9 --day 14
// generates folder 'day14'`,
	Run: func(cmd *cobra.Command, args []string) {
		day := time.Now().Day()

		if len(args) > 0 {
			maybe, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("First argument is not an integer: %s", err.Error())
			}

			if maybe < 1 || maybe > 31 {
				log.Fatalf("Specified day is outside of the valid 1-31 range. Got %d", maybe)
			}
			log.Printf("Day was specified, continuing with value %d\n", maybe)
			day = maybe
		} else {
			log.Printf("No day was specified, using today's: %d\n", day)
		}

		dirName := fmt.Sprintf("day%d", day)

		// dir exists
		_, err := os.ReadDir(dirName)
		if err == nil {
			log.Fatalf("Directory for day %d already exists, aborting", day)
		}

		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("Reading directory encountered an unexpected error, aborting.Error got: %s", err.Error())
		}

		templateData := dayData{
			Pkg: dirName,
			Day: day,
		}

		tps, err := template.ParseFiles("dayn/part1.go.tpl", "dayn/part2.go.tpl", "dayn/readme.md.tpl")
		if err != nil {
			log.Fatalf("Parsing template files encountered an error: %s", err.Error())
		}

		err = os.Mkdir(dirName, 0755)
		if err != nil {
			log.Fatalf("Creating the directory %s failed: %s", dirName, err.Error())
		}

		p1, err := os.OpenFile(dirName+"/part1.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("Creating file 'part1.go' failed: %s", err.Error())
		}
		defer p1.Close()

		p2, err := os.OpenFile(dirName+"/part2.go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("Creating file 'part2.go' failed: %s", err.Error())
		}
		defer p2.Close()

		rdme, err := os.OpenFile(dirName+"/readme.md", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("Creating file 'readme.md' failed: %s", err.Error())
		}
		defer rdme.Close()

		i1, err := os.OpenFile(dirName+"/input1.txt", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("Creating file 'readme.md' failed: %s", err.Error())
		}
		defer i1.Close()

		i2, err := os.OpenFile(dirName+"/example.txt", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("Creating file 'readme.md' failed: %s", err.Error())
		}
		defer i2.Close()

		err = tps.ExecuteTemplate(p1, "part1.go.tpl", templateData)
		if err != nil {
			log.Fatalf("Executing template for 'part1.go' failed: %s", err)
		}

		err = tps.ExecuteTemplate(p2, "part2.go.tpl", templateData)
		if err != nil {
			log.Fatalf("Executing template for 'part2.go' failed: %s", err)
		}

		err = tps.ExecuteTemplate(rdme, "readme.md.tpl", templateData)
		if err != nil {
			log.Fatalf("Executing template for 'readme.md' failed: %s", err)
		}

		log.Printf("Folder for %s generated, all good!", dirName)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
