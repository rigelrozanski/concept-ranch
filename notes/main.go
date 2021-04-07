package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
	wb "github.com/rigelrozanski/wb/lib"
)

// name of the wb used to hold command notes
const wbName = "notes"

// structure of notes to be stored
type Notes struct {
	Aliases map[string]string `yaml:"aliases"`
	Goodies map[string]string `yaml:"goodies"`
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("please provide the command to retrieve notes for")
		return
	}
	search := strings.Join(args, " ")

	// get the notes
	notesBz, found := wb.GetWBRaw(wbName)
	if !found {
		fmt.Println("can't find notes wb")
		return
	}

	// parse the notes file
	var notes Notes
	err := yaml.Unmarshal(notesBz, &notes)
	if len(notes.Goodies) == 0 || err != nil {
		fmt.Printf("debug err: %v\n", err)
		fmt.Printf("debug notes: %v\n", notes)
		fmt.Println("notes wb is not properly formatted")
		return
	}

	// print the notes
	note := notes.Goodies[search]

	// attempt to find an alias
	if len(note) == 0 {
		revisedSearch := notes.Aliases[search]
		if len(revisedSearch) > 0 {
			note = notes.Goodies[revisedSearch]
		}
	}
	if len(note) == 0 {
		fmt.Println("can't find note")
		return
	}

	fmt.Printf("\n%v\n", note)
}
