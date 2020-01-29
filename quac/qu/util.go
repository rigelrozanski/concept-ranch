package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/wranch/quac/qu/lib"
)

func Consume(consumedID, optionalEntry string) {
	consumed, err := lib.ParseID(consumedID)
	if err != nil {
		log.Fatalf("bad id %v", consumedID)
	}
	consumerFilepath := lib.SetConsume(uint32(consumed), optionalEntry)
	if optionalEntry == "" {
		lib.OpenText(consumerFilepath)
	}
}

func Consumes(consumedID, consumesID string) {
	consumed, err := lib.ParseID(consumedID)
	if err != nil {
		log.Fatalf("bad id %v", consumedID)
	}
	consumes, err := lib.ParseID(consumesID)
	if err != nil {
		log.Fatalf("bad id %v", consumesID)
	}
	lib.SetConsumes(uint32(consumed), uint32(consumes))
}

func Zombie(zombieID string) {
	zombie, err := lib.ParseID(zombieID)
	if err != nil {
		log.Fatalf("bad id %v", zombieID)
	}
	lib.SetZombie(uint32(zombie))
}

func Lineage(idStr string) {
	id, err := lib.ParseID(idStr)
	if err != nil {
		log.Fatalf("bad id %v", idStr)
	}
	fmt.Print(lib.GetLineage(uint32(id)))
}

func Transcribe(optionalQuery string) {

	if optionalQuery == "" {

		// TODO
		return
	}

	consumed, err := lib.ParseID(optionalQuery)
	if err == nil {
		idea := lib.GetIdeaByID(consumed)
		if !(idea.IsImage() || idea.IsAudio()) {
			fmt.Println("this idea is not an image or audio cannot be transcribed")
			os.Exit(1)
		}
		lib.Open(idea.Path())

		// read input from console
		fmt.Println("Please enter the entry text here: (or just hit enter to open editor)")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		optionalEntry := consoleScanner.Text()

		consumerFilepath := lib.SetConsume(uint32(consumed), optionalEntry)
		if optionalEntry == "" {
			lib.OpenText(consumerFilepath)
		}
		return
	}

	subsetTagsImages := lib.GetAllIdeasNonConsuming().
		WithTags(parseTags(optionalQuery)).
		WithImage()

	if len(subsetTagsImages) == 0 {
		fmt.Println("no active images to transcribe with those tags")
		os.Exit(1)
	}

	for _, idea := range subsetTagsImages {

		lib.Open(idea.Path())

		// read input from console
		fmt.Println("Please enter the entry text here: (or just hit enter to open editor)")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		optionalEntry := consoleScanner.Text()

		consumerFilepath := lib.SetConsume(idea.Id, optionalEntry)
		if optionalEntry == "" {
			lib.OpenText(consumerFilepath)
		}
	}
}

func QuickQuery(unsplitTagsOrID string) {
	id, err := lib.ParseID(unsplitTagsOrID)
	if err == nil {
		ViewByID(uint32(id))
		return
	}
	splitTags := parseTags(unsplitTagsOrID)
	ViewByTags(splitTags)
}

func NewEmptyEntry(unsplitTags string) {
	splitTags := parseTags(unsplitTags)
	idear := lib.NewNonConsumingTextIdea(splitTags)
	writePath := path.Join(lib.IdeasDir, idear.Filename)
	lib.IncrementID()
	lib.OpenText(writePath)
}

func QuickEntry(unsplitTags, entry string) {
	splitTags := parseTags(unsplitTags)
	Entry(entry, splitTags)
}

func MultiOpen(unsplitTagsOrID string) {
	id, err := lib.ParseID(unsplitTagsOrID)
	if err == nil {
		filePath := lib.GetFilepathByID(uint32(id))
		lib.Open(filePath)
		return
	}
	splitTags := parseTags(unsplitTagsOrID)
	MultiOpenByTags(splitTags)
}

func parseIdStr(idStr string) uint32 {
	idI, err := lib.ParseID(idStr)
	if err != nil {
		log.Fatalf("error parsing id, error: %v", err)
	}
	return uint32(idI)
}

func parseTags(tagsGrouped string) []string {
	return strings.Split(tagsGrouped, ",")
}

func RemoveByID(idStr string) {
	id := parseIdStr(idStr)
	lib.RemoveByID(id)
}

func CopyByID(idStr string) {
	id := parseIdStr(idStr)
	lib.Open(lib.CopyByID(id))
}

func ListTagsByID(idStr string) {
	id := parseIdStr(idStr)
	idea := lib.GetIdeaByID(id)
	fmt.Println(idea.Tags)
}

func KillTagByID(idStr, tagToKill string) {
	id := parseIdStr(idStr)
	idea := lib.GetIdeaByID(id)
	origFilename := idea.Filename
	(&idea).RemoveTag(tagToKill)
	(&idea).UpdateFilename()

	origPath := path.Join(lib.IdeasDir, origFilename)
	newPath := path.Join(lib.IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func AddTagByID(idStr, tagToAdd string) {
	id := parseIdStr(idStr)
	idea := lib.GetIdeaByID(id)
	origFilename := idea.Filename
	idea.Tags = append(idea.Tags, tagToAdd)
	(&idea).UpdateFilename()

	origPath := path.Join(lib.IdeasDir, origFilename)
	newPath := path.Join(lib.IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func RenameTag(from, to string) {
	files, err := ioutil.ReadDir(lib.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		origFn := file.Name()
		if !strings.Contains(origFn, from) {
			continue
		}
		idea := lib.NewIdeaFromFilename(origFn)
		(&idea).RenameTag(from, to)
		(&idea).UpdateFilename()

		// perform the file rename
		origPath := path.Join(lib.IdeasDir, origFn)
		newPath := path.Join(lib.IdeasDir, idea.Filename)
		err := os.Rename(origPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DestroyTag(tag string) {
	files, err := ioutil.ReadDir(lib.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		origFn := file.Name()
		if !strings.Contains(origFn, tag) {
			continue
		}
		idea := lib.NewIdeaFromFilename(origFn)
		(&idea).RemoveTag(tag)
		(&idea).UpdateFilename()

		// perform the file rename
		origPath := path.Join(lib.IdeasDir, origFn)
		newPath := path.Join(lib.IdeasDir, idea.Filename)
		err := os.Rename(origPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//__________________

func ListAllTags() {
	ideas := lib.GetAllIdeas()
	fmt.Println(ideas.UniqueTags())
}

func ListAllTagsWithTags(tagsGrouped string) {
	ideas := lib.GetAllIdeas()
	queryTags := parseTags(tagsGrouped)
	subset := ideas.WithTags(queryTags)
	uniqueTags := subset.UniqueTags()
	outTags := make([]string, len(uniqueTags))

	// remove the query tags from this list
	i := 0
	for _, uTag := range uniqueTags {
		isQTag := false
		for _, qTag := range queryTags {
			if uTag == qTag {
				isQTag = true
			}
		}
		if !isQTag {
			outTags[i] = uTag
			i++
		}
	}

	fmt.Println(outTags)
}

func ListAllFiles() {
	ideas := lib.GetAllIdeas()
	for _, idea := range ideas {
		fmt.Println(idea.Filename)
	}
}

func ListAllFilesWithTags(tagsGrouped string) {
	ideas := lib.GetAllIdeas()
	subset := ideas.WithTags(parseTags(tagsGrouped))
	for _, idea := range subset {
		fmt.Println(idea.Filename)
	}
}

func ViewByID(id uint32) {
	content, found := lib.GetContentByID(id)
	if !found {
		fmt.Println("nothing found with that id")
	}
	fmt.Printf("%s\n", content)
}

func ViewByTags(tags []string) {
	content, found := lib.ConcatAllContentFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	fmt.Printf("%s\n", content)
}

func MultiOpenByTags(tags []string) {
	found, maxFNLen, singleReturn := lib.WriteWorkingContentAndFilenamesFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
		return
	}
	// if only a single entry is found then lib.Open only it!
	if singleReturn != "" {
		lib.Open(singleReturn)
		return
	}
	lib.OpenTextSplit(lib.WorkingFnsFile, lib.WorkingContentFile, maxFNLen)
	lib.SaveFromWorkingFiles()
}

func RemoveById(id uint32) error {
	return nil
}

// create an entry
func Entry(entryOrPath string, tags []string) {

	if cmn.FileExists(entryOrPath) { // is a path

		idea := lib.NewIdeaFromFile(tags, entryOrPath)

		err := cmn.Copy(entryOrPath, idea.Path())
		if err != nil {
			log.Fatal(err)
		}
		lib.IncrementID()

		return
	}

	idea := lib.NewNonConsumingTextIdea(tags)
	err := cmn.WriteLines([]string{entryOrPath}, idea.Path())
	if err != nil {
		log.Fatalf("error writing new file: %v", err)
	}
	lib.IncrementID()
}
