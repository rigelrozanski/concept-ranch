package quac

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gookit/color"
)

func Browse() {
	fmt.Println("Please enter the entry text here: (or just hit enter to open editor)")

	//tickTime := 100 * time.Millisecond
	//ticker := time.Tick(tickTime)
	highlighted := false
	for {

		c := color.New(color.FgDefault, color.BgDefault).Render
		if highlighted {
			c = color.New(color.FgBlue, color.BgYellow).Render
		}

		fmt.Printf("\rOn %v/10", c("bloop"))

		reader := bufio.NewReader(os.Stdin)
		char, err := reader.ReadByte()
		if err == nil {
			fmt.Printf("debug char: %s\n", char)
		}

		//<-ticker
	}
}
