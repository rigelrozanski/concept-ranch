package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/qi/lib"
)

func QuickQuery(unsplitTags string) {
	splitTags := strings.Split(unsplitTags, ",")
	ViewByTags(splitTags)
}

func NewEmptyEntry(unsplitTags string) {
	splitTags := strings.Split(unsplitTags, ",")
	idea := lib.NewNonConsumingTextIdea(splitTags)
	writePath := path.Join(lib.IdeasDir, idea.Filename)
	lib.IncrementID()
	openText(writePath)
}

func QuickEntry(unsplitTags, entry string) {
	splitTags := strings.Split(unsplitTags, ",")
	Entry(entry, splitTags)
}

func MultiOpen(unsplitTags string) {
	splitTags := strings.Split(unsplitTags, ",")
	MultiOpenByTags(splitTags)
}

func RemoveByID(idStr string) {
	idI, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("error parsing id, error: %v", err)
	}
	id := uint32(idI)

	lib.RemoveByID(id)
}

//__________________

func ListAllTags() {
	ideas := lib.PathToIdeas(lib.IdeasDir)
	fmt.Println(ideas.UniqueTags())
}

func ViewByTags(tags []string) {
	content, found := lib.ConcatAllContentFromTags(lib.IdeasDir, tags)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	fmt.Printf("%s\n", content)
}

func MultiOpenByTags(tags []string) {
	found, maxFNLen := lib.WriteWorkingContentAndFilenamesFromTags(lib.IdeasDir, tags)
	if !found {
		fmt.Println("nothing found with those tags")
		return
	}
	openTextSplit(lib.WorkingFnsFile, lib.WorkingContentFile, maxFNLen)
}

func RemoveById(id uint32) error {
	return nil
}

func Entry(entry string, tags []string) {

	idea := lib.NewNonConsumingTextIdea(tags)
	writePath := path.Join(lib.IdeasDir, idea.Filename)
	err := cmn.WriteLines([]string{entry}, writePath)
	if err != nil {
		log.Fatalf("error writing new file: %v", err)
	}
	lib.IncrementID()
}

//_______________________________________________________________________________________________________

func openText(pathToOpen string) {

	cmd := exec.Command("vim", "-c", "+normal 1G1|", pathToOpen) //start in the upper left corner nomatter
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func openTextSplit(pathToOpenLeft, pathToOpenRight string, maxFNLen int) {

	cmd := exec.Command("vim",
		"-c", "vertical resize "+strconv.Itoa(maxFNLen+4)+" | execute \"normal \\<C-w>\\<C-l>\"",
		"-O", pathToOpenLeft, pathToOpenRight)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
