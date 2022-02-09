package quac

import (
	"fmt"

	"github.com/rigelrozanski/thranch/quac/idea"
)

func GetStats() {

	idears := idea.GetAllIdeas()
	ncidears := idea.GetAllIdeasNonConsuming()
	imgidears := ncidears.WithImage()
	wot, _ := idea.NewTagWithout("DNT", "")
	untranscribed := imgidears.WithTags(wot)
	utags := ncidears.UniqueTags()

	fmt.Println("\n ~ IDEA STATISTICS ~ ")
	fmt.Printf("ideas:\t\t%v\n", len(ncidears))
	fmt.Printf("consumed:\t%v\n", len(idears)-len(ncidears))
	fmt.Printf("unique tags:\t%v\n", len(utags))
	fmt.Printf("images:\t\t%v\n", len(imgidears))
	fmt.Printf("untranscribed:\t%v\n", len(untranscribed))
	fmt.Println("")
}
