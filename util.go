package main

import (
	"fmt"
	"log"
	"path"
	"strings"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/qi/lib"
)

func QuickQuery(unsplitTags string) {
	splitTags, err := strings.Split(unsplitTags, ",")
	if err != nil {
		log.Fatal(err)
	}
	ViewByTags(splitTags)
}

func QuickEntry(entry, unsplitTags string) {
	splitTags, err := strings.Split(unsplitTags, ",")
	if err != nil {
		log.Fatal(err)
	}
	Entry(entry, splitTags)
}

//__________________

func ListAllTags() {
	ideas := lib.GetIdeas(lib.IdeasDir)
	fmt.Println(ideas.GetUnique())
}

func ViewByTags(tags ...string) {
	content, found := lib.GetByTags(lib.IdeasDir)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	fmt.Printf("%v\n", content)
}

func RemoveById(id uint32) error {
}

//func edit(name string) (err error) {

//origWbBz, found := lib.GetWBRaw(name)
//if !found {
//name, err = getNameFromShortcut(name)
//if err != nil {
//return err
//}
//}

//wbPath, err := lib.GetWbPath(name)
//if err != nil {
//return err
//}

//cmd := exec.Command("vim", "-c", "+normal 1G1|", wbPath) //start in the upper left corner nomatter
//cmd.Stdin = os.Stdin
//cmd.Stdout = os.Stdout
//err = cmd.Run()
//if err != nil {
//return err
//}

//// log if there was a modification
//newWbBz, found := lib.GetWBRaw(name)
//if !found {
//panic("wuz found now isn't")
//}
//if bytes.Compare(origWbBz, newWbBz) != 0 {
//log("modified wb", name)
//}
//return nil
//}

func Entry(entry string, tags ...string) {

	idea := lib.NewNonConsumingIdea(tags)
	writePath := path.Join(lib.IdeasDir, idea.Filename)
	err := cmn.WriteLines([]string{entry}, writePath)
	if err != nil {
		log.Fatalf("error writing new file: %v", err)
	}
	lib.IncrementID()
}
