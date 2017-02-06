package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/codegangsta/cli"
)

func commandNotFound(c *cli.Context, command string) {
	app := c.App
	var choices []string
	for _, cmd := range app.Commands {
		choices = append(choices, cmd.Name)
	}
	//choices := globalOptionsNames(app)
	currentMin := 50
	bestSuggestion := ""
	for _, choice := range choices {
		similarity := levenshtein(choice, command)
		tmpMin := min(currentMin, similarity)
		if tmpMin < currentMin {
			bestSuggestion = choice
			currentMin = tmpMin
		}
	}

	suggestion := []string{fmt.Sprintf("Unrecognized command: %s", command),
		"",
		"Did you mean this?",
		fmt.Sprintf("\t%s\n", bestSuggestion),
		"",
		"You can use the --h or --help flag to see a list of options for any command.",
		"Examples:",
		"\track --h",
		"\track servers --h",
		"\track block-storage volumes --help",
		"",
	}

	fmt.Fprintf(c.App.Writer, strings.Join(suggestion, "\n"))
}

func levenshtein(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
