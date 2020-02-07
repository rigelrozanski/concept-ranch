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
	"github.com/rigelrozanski/thranch/quac"
	"github.com/rigelrozanski/thranch/quac/idea"
)

func Consume(consumedID, optionalEntry string) {
	consumed, err := quac.ParseID(consumedID)
	if err != nil {
		log.Fatalf("bad id %v", consumedID)
	}
	consumerFilepath := quac.SetConsume(uint32(consumed), optionalEntry)
	if optionalEntry == "" {
		quac.OpenText(consumerFilepath)
	}
}

func Consumes(consumedID, consumesID string) {
	consumed, err := quac.ParseID(consumedID)
	if err != nil {
		log.Fatalf("bad id %v", consumedID)
	}
	consumes, err := quac.ParseID(consumesID)
	if err != nil {
		log.Fatalf("bad id %v", consumesID)
	}
	quac.SetConsumes(uint32(consumed), uint32(consumes))
}

func Zombie(zombieID string) {
	zombie, err := quac.ParseID(zombieID)
	if err != nil {
		log.Fatalf("bad id %v", zombieID)
	}
	quac.SetZombie(uint32(zombie))
}

func Lineage(idStr string) {
	id, err := quac.ParseID(idStr)
	if err != nil {
		log.Fatalf("bad id %v", idStr)
	}
	fmt.Print(quac.GetLineage(uint32(id)))
}

func Transcribe(optionalQuery string) {

	if optionalQuery == "" {

		// TODO
		return
	}

	consumed, err := quac.ParseID(optionalQuery)
	if err == nil {
		idea := quac.GetIdeaByID(consumed)
		if !(idea.IsImage() || idea.IsAudio()) {
			fmt.Println("this idea is not an image or audio cannot be transcribed")
			os.Exit(1)
		}
		quac.Open(idea.Path())

		// read input from console
		fmt.Println("Please enter the entry text here: (or just hit enter to open editor)")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		optionalEntry := consoleScanner.Text()

		consumerFilepath := quac.SetConsume(uint32(consumed), optionalEntry)
		if optionalEntry == "" {
			quac.OpenText(consumerFilepath)
		}
		return
	}

	subsetTagsImages := quac.GetAllIdeasNonConsuming().
		WithTags(parseTags(optionalQuery)).
		WithImage()

	if len(subsetTagsImages) == 0 {
		fmt.Println("no active images to transcribe with those tags")
		os.Exit(1)
	}

	for _, idea := range subsetTagsImages {

		quac.Open(idea.Path())

		// read input from console
		fmt.Println("Please enter the entry text here: (or just hit enter to open editor)")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		optionalEntry := consoleScanner.Text()

		consumerFilepath := quac.SetConsume(idea.Id, optionalEntry)
		if optionalEntry == "" {
			quac.OpenText(consumerFilepath)
		}
	}
}

func QuickQuery(unsplitTagsOrID string) {
	id, err := quac.ParseID(unsplitTagsOrID)
	if err == nil {
		ViewByID(uint32(id))
		return
	}
	splitTags := parseTags(unsplitTagsOrID)
	ViewByTags(splitTags)
}

func NewEmptyEntry(unsplitTags string) {
	splitTags := parseTags(unsplitTags)
	idear := quac.NewNonConsumingTextIdea(splitTags)
	writePath := path.Join(quac.IdeasDir, idear.Filename)
	quac.IncrementID()
	quac.OpenText(writePath)
}

func SetEncryption(idStr string) {
	id, err := quac.ParseID(idStr)
	if err != nil {
		log.Fatalf("error parsing id, error: %v", err)
	}

	quac.SetEncryptionById(id)
}

func QuickEntry(unsplitTags, entry string) {
	splitTags := parseTags(unsplitTags)
	Entry(entry, splitTags)
}

func MultiOpen(unsplitTagsOrID string, forceSplitView bool) {
	id, err := quac.ParseID(unsplitTagsOrID)
	if err == nil {
		filePath := quac.GetFilepathByID(uint32(id))
		if forceSplitView {
			maxFNLen := quac.WriteWorkingContentAndFilenamesFromFilePath(filePath)
			quac.OpenTextSplit(quac.WorkingFnsFile, quac.WorkingContentFile, maxFNLen)
			quac.SaveFromWorkingFiles()
			return
		}
		quac.Open(filePath)
		return
	}
	splitTags := parseTags(unsplitTagsOrID)
	MultiOpenByTags(splitTags, forceSplitView)
}

func parseIdStr(idStr string) uint32 {
	idI, err := quac.ParseID(idStr)
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
	quac.RemoveByID(id)
}

func RemoveAcrossIDs(idStr, idStr2 string) {
	id := parseIdStr(idStr)
	id2 := parseIdStr(idStr2)
	if id >= id2 {
		log.Fatalf("second id must be greater first id")
	}
	for i := id; i <= id2; i++ {
		quac.RemoveByID(i)
	}
}

func CopyByID(idStr string) {
	id := parseIdStr(idStr)
	quac.Open(quac.CopyByID(id))
}

func ListTagsByID(idStr string) {
	id := parseIdStr(idStr)
	idea := quac.GetIdeaByID(id)
	fmt.Println(idea.Tags)
}

func KillTagByID(idStr, tagToKill string) {
	id := parseIdStr(idStr)
	idea := quac.GetIdeaByID(id)
	origFilename := idea.Filename
	(&idea).RemoveTag(tagToKill)
	(&idea).UpdateFilename()

	origPath := path.Join(quac.IdeasDir, origFilename)
	newPath := path.Join(quac.IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func AddTagByID(idStr, tagToAdd string) {
	id := parseIdStr(idStr)
	idea := quac.GetIdeaByID(id)
	origFilename := idea.Filename
	idea.Tags = append(idea.Tags, tagToAdd)
	(&idea).UpdateFilename()

	origPath := path.Join(quac.IdeasDir, origFilename)
	newPath := path.Join(quac.IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func RenameTag(from, to string) {
	files, err := ioutil.ReadDir(quac.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		origFn := file.Name()
		if !strings.Contains(origFn, from) {
			continue
		}
		idea := quac.NewIdeaFromFilename(origFn)
		(&idea).RenameTag(from, to)
		(&idea).UpdateFilename()

		// perform the file rename
		origPath := path.Join(quac.IdeasDir, origFn)
		newPath := path.Join(quac.IdeasDir, idea.Filename)
		err := os.Rename(origPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DestroyTag(tag string) {
	files, err := ioutil.ReadDir(quac.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		origFn := file.Name()
		if !strings.Contains(origFn, tag) {
			continue
		}
		idea := quac.NewIdeaFromFilename(origFn)
		(&idea).RemoveTag(tag)
		(&idea).UpdateFilename()

		// perform the file rename
		origPath := path.Join(quac.IdeasDir, origFn)
		newPath := path.Join(quac.IdeasDir, idea.Filename)
		err := os.Rename(origPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//__________________

func ListAllTags() {
	ideas := quac.GetAllIdeas()
	fmt.Println(ideas.UniqueTags())
}

func ListAllTagsWithTags(tagsGrouped string) {
	ideas := quac.GetAllIdeas()
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
	ideas := quac.GetAllIdeas()
	if len(ideas) == 0 {
		fmt.Println("no ideas found")
	}
	for _, idea := range ideas {
		fmt.Println(idea.Filename)
	}
}

func ListAllFilesWithTags(tagsGrouped string) {
	ideas := quac.GetAllIdeas()
	subset := ideas.WithTags(parseTags(tagsGrouped))
	if len(subset) == 0 {
		fmt.Println("no ideas found with those tags")
	}
	for _, idea := range subset {
		fmt.Println(idea.Filename)
	}
}

func ListAllFilesLast() {
	ids := idea.GetLastIDs()
	for _, id := range ids {
		fmt.Println(quac.GetFilenameByID(id))
	}
}

func ViewByID(id uint32) {
	content, found := quac.GetContentByID(id)
	if !found {
		fmt.Println("nothing found with that id")
	}
	fmt.Printf("%s\n", content)
}

func ViewByTags(tags []string) {
	content, found := quac.ConcatAllContentFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	fmt.Printf("%s\n", content)
}

func MultiOpenByTags(tags []string, forceSplitView bool) {
	found, maxFNLen, singleReturn := quac.WriteWorkingContentAndFilenamesFromTags(tags, forceSplitView)
	if !found {
		fmt.Println("nothing found with those tags")
		return
	}
	// if only a single entry is found then quac.Open only it!
	if singleReturn != "" && !forceSplitView {
		quac.Open(singleReturn)
		return
	}
	quac.OpenTextSplit(quac.WorkingFnsFile, quac.WorkingContentFile, maxFNLen)
	quac.SaveFromWorkingFiles()
}

func RemoveById(id uint32) error {
	return nil
}

// create an entry
func Entry(entryOrPath string, tags []string) {

	// TODO this logic should exist in the library
	hasFN := false
	for _, tag := range tags {
		if strings.Contains(tag, "FILENAME") {
			hasFN = true
			break
		}
	}

	if cmn.FileExists(entryOrPath) { // is a path

		fod, err := os.Stat(entryOrPath)
		if err != nil {
			log.Fatal(err)
		}
		var filepaths []string

		if fod.Mode().IsDir() {
			files, err := ioutil.ReadDir(entryOrPath)
			if err != nil {
				log.Fatal(err)
			}

			for _, file := range files {
				if !file.IsDir() {
					filepath := path.Join(entryOrPath, file.Name())
					filepaths = append(filepaths, filepath)
				}
			}
			if len(filepaths) == 0 {
				log.Fatal("directory is empty")
			}
		} else {
			filepaths = []string{entryOrPath}
		}

		for _, filepath := range filepaths {

			// skip if folder
			fod2, err := os.Stat(filepath)
			if err != nil {
				log.Fatal(err)
			}
			if fod2.Mode().IsDir() {
				continue
			}

			// TODO this logic should exist in the lib
			tags2 := make([]string, len(tags))
			copy(tags2, tags)
			if hasFN {
				for i, tag := range tags2 {
					if strings.Contains(tag, "FILENAME") {
						filebase := strings.TrimSuffix(path.Base(filepath), path.Ext(filepath))
						tags2[i] = strings.Replace(tag, "FILENAME", filebase, 2)
					}
				}
			}

			idea := quac.NewIdeaFromFile(tags2, filepath)
			err = cmn.Copy(filepath, idea.Path())
			if err != nil {
				log.Fatal(err)
			}
			quac.IncrementID()
		}
		return
	}

	if hasFN {
		log.Fatal("the tag \"FILENAME\" is reserved for file entry not raw-text-entry")
	}

	idea := quac.NewNonConsumingTextIdea(tags)
	err := cmn.WriteLines([]string{entryOrPath}, idea.Path())
	if err != nil {
		log.Fatalf("error writing new file: %v", err)
	}
	quac.IncrementID()
}
