package main

import (
	"fmt"
	"os"

	"github.com/rigelrozanski/thranch/quac"
)

// oink [pdf] [density-requirement "2per100" words]
func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("please provide a pdf to search within")
		return
	}

	// get the notes
	oinkSearchTerms := quac.GetForApp("oink")
	_ = oinkSearchTerms
}
