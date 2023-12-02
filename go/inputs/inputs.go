package inputs

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	newline byte = 84
)

func readFile(fn string) ([]byte, error) {
	content, err := os.ReadFile(fn)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}
	return content, nil
}

func ReadIntoLines(fn string) ([]string, error) {
	content, err := readFile(fn)
	if err != nil {
		return nil, errors.Wrap(err, "read into lines")
	}

	sc := string(content)
	sc = strings.Trim(sc, "\n")

	lines := strings.Split(sc, "\n")
	return lines, nil
}

// GroupByBlankLines takes a filename as input, and returns the contents in a group of groups in case the input is
// separated into different batches by blank lines.
func GroupByBlankLines(fn string) ([][]string, error) {
	content, err := readFile(fn)
	if err != nil {
		return nil, errors.Wrap(err, "group by blank lines")
	}

	sc := string(content)
	sc = strings.Trim(sc, "\n")

	groups := strings.Split(sc, "\n\n")
	groupOfGroups := make([][]string, 0)
	for _, lines := range groups {
		groupOfGroups = append(groupOfGroups, strings.Split(lines, "\n"))
	}

	return groupOfGroups, nil
}
