package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

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

	ideaImages := quac.GetAllIdeasNonConsuming().
		WithImage().WithoutTag("DNT")

	if optionalQuery != "" {
		ideaImages = ideaImages.WithTags(idea.ParseClumpedTags(optionalQuery))
		if len(ideaImages) == 0 {
			fmt.Println("no active images to transcribe with those tags")
			os.Exit(1)
		}
	}

	// shuffle the images to transcribe
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ideaImages), func(i, j int) { ideaImages[i], ideaImages[j] = ideaImages[j], ideaImages[i] })

	fmt.Println("\n-------------------------------------------------------------")
	fmt.Println("                     ~ Instructions ~")
	fmt.Println("       After each transcription item enter either:")
	fmt.Println("         - nothing to open up your editor where you")
	fmt.Println("             may enter the transcription text")
	fmt.Println("         - transcribed entry text")
	fmt.Println("         - DNT for do not transcribe (DNT is added as")
	fmt.Println("             a tag and never asked to transcribe again)")
	fmt.Println("         - KILL to delete the entry")
	fmt.Println("         - SKIP to skip")
	fmt.Println("         - QUIT to quit")
	fmt.Println("-------------------------------------------------------------\n")

	for _, idea := range ideaImages {

		quac.Open(idea.Path())

		// read input from console
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		optionalEntry := consoleScanner.Text()
		optionalEntry = strings.TrimSpace(optionalEntry)

		if optionalEntry == "DNT" {
			AddTagByIdea(idea, "DNT")
			fmt.Println("ol'right never transcribing again!")
			continue
		}
		if optionalEntry == "SKIP" {
			fmt.Println("skipp'd")
			continue
		}
		if optionalEntry == "KILL" {
			quac.RemoveByID(idea.Id)
			fmt.Println("killed it")
			continue
		}
		if optionalEntry == "QUIT" {
			break
		}

		consumerFilepath := quac.SetConsume(idea.Id, optionalEntry)
		if optionalEntry == "" {
			quac.OpenText(consumerFilepath)
		}
		fmt.Printf("created: %v\n", consumerFilepath)
	}
}

func Retag() {
	untaggedIdeas := idea.GetAllIdeasNonConsuming().WithTag("UNTAGGED")

	fmt.Println("             ~ Instructions ~")
	fmt.Println("enter desired tags seperated by spaces")
	fmt.Println("alternatively, enter SKIP to skip or QUIT to quit\n")

	for _, idea := range untaggedIdeas {
		quac.View(idea.Path())

		// read input from console
		fmt.Println("Desired tags:")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		spacedTags := consoleScanner.Text()
		tags := strings.Fields(spacedTags)

		if len(tags) == 1 && tags[0] == "SKIP" {
			continue
		}
		if len(tags) == 1 && tags[0] == "QUIT" {
			break
		}

		origPath := idea.Path()

		// add the tags
		idea.Tags = tags
		(&idea).UpdateFilename()

		// perform the file rename
		err := os.Rename(origPath, idea.Path())
		fmt.Printf("retagged to:\n%v\n", idea.Filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func QuickQuery(unsplitTagsOrID string) {
	id, err := quac.ParseID(unsplitTagsOrID)
	if err == nil {
		fp, found := quac.GetFilepathByID(uint32(id))
		if !found {
			fmt.Println("nothing found at that id")
			os.Exit(1)
		}
		quac.View(fp)
		return
	}
	splitTags := idea.ParseClumpedTags(unsplitTagsOrID)
	ViewByTags(splitTags)
}

func NewEmptyEntry(unsplitTags string) {
	splitTags := idea.ParseClumpedTags(unsplitTags)
	idear := quac.NewNonConsumingTextIdea(splitTags)
	writePath := path.Join(quac.IdeasDir, idear.Filename)
	quac.IncrementID()
	quac.OpenText(writePath)
}

func ManualEntry(commonTagsClumped string) {
	commonTags := idea.ParseClumpedTags(commonTagsClumped)

	fmt.Println("_________________________________________")
	fmt.Println("INTERACTIVE MANUAL TRANSCRIPTION")
	fmt.Println(" - tags to be entered seperated by spaces")
	fmt.Println(" - quit with ctrl-c, or by typing QUIT during tag entry")
	fmt.Println("")
	for {

		fmt.Print("enter tags:\t")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		newTags := strings.Fields(consoleScanner.Text())
		if len(newTags) > 0 && newTags[0] == "QUIT" {
			break
		}
		tags := append(commonTags, newTags...)

		fmt.Print("enter entry:\t")
		consoleScanner = bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		entry := consoleScanner.Text()

		Entry(entry, tags)
	}
}

func SetEncryption(idStr string) {
	id, err := quac.ParseID(idStr)
	if err != nil {
		log.Fatalf("error parsing id, error: %v", err)
	}

	quac.SetEncryptionById(id)
}

func QuickEntry(unsplitTags, entry string) {
	splitTags := idea.ParseClumpedTags(unsplitTags)
	Entry(entry, splitTags)
}

func MultiOpen(unsplitTagsOrID string, forceSplitView bool) {

	startID, endID, isRange := IsIDRange(unsplitTagsOrID)
	if isRange {
		MultiOpenByRange(startID, endID, forceSplitView)
		return
	}

	id, err := quac.ParseID(unsplitTagsOrID)
	if err == nil {
		filePath, found := quac.GetFilepathByID(uint32(id))
		if !found {
			fmt.Println("nothing found at that ID")
			os.Exit(1)
		}
		if forceSplitView {
			maxFNLen := quac.WriteWorkingContentAndFilenamesFromFilePath(filePath)
			origBzFN, origBzContent := quac.GetOrigWorkingFileBytes()
			quac.OpenTextSplit(quac.WorkingFnsFile, quac.WorkingContentFile, maxFNLen)
			quac.SaveFromWorkingFiles(origBzFN, origBzContent)
			return
		}
		quac.Open(filePath)
		return
	}
	splitTags := idea.ParseClumpedTags(unsplitTagsOrID)
	MultiOpenByTags(splitTags, forceSplitView)
}

func OpenWorking() {
	quac.OpenTextSplit(quac.WorkingFnsFile, quac.WorkingContentFile, 50)
}

func SaveWorking() {
	quac.SaveFromWorkingFiles([]byte{}, []byte{})
}

func parseIdStr(idStr string) uint32 {
	idI, err := quac.ParseID(idStr)
	if err != nil {
		log.Fatalf("error parsing id, error: %v", err)
	}
	return uint32(idI)
}

func RemoveByID(idOrIds string) {
	startID, endID, isRange := IsIDRange(idOrIds)
	if isRange {
		RemoveAcrossIDs(startID, endID)
		fmt.Println("roger, removed everything between those ids")
		return
	}

	id := parseIdStr(idOrIds)
	quac.RemoveByID(id)
	fmt.Println("roger, removed that id")
}

func RemoveAcrossIDs(id1, id2 uint32) {
	if id1 >= id2 {
		log.Fatalf("second id must be greater first id")
	}
	for i := id1; i <= id2; i++ {
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
	AddTagByIdea(idea, tagToAdd)
}

// TODO move to library
func AddTagByIdea(idea quac.Idea, tagToAdd string) {
	origFilename := idea.Filename
	(&idea).AddTag(tagToAdd)
	(&idea).UpdateFilename()
	origPath := path.Join(quac.IdeasDir, origFilename)
	newPath := path.Join(quac.IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func AddTagToMany(tagToAdd, manyTagsClumped string) {

	manyTags := idea.ParseClumpedTags(manyTagsClumped)
	ideas := quac.GetAllIdeas().WithAnyOfTags(manyTags)
	for _, idea := range ideas {
		AddTagByIdea(idea, tagToAdd)
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
	queryTags := idea.ParseClumpedTags(tagsGrouped)
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

func IsIDRange(query string) (idStart, idEnd uint32, isRange bool) {
	sp := strings.Split(query, "-")
	if len(sp) != 2 {
		return 0, 0, false
	}

	idStart, err := quac.ParseID(sp[0])
	if err != nil {
		return 0, 0, false
	}
	idEnd, err = quac.ParseID(sp[1])
	if err != nil {
		return 0, 0, false
	}

	return idStart, idEnd, true
}

func ListAllFilesWithQuery(query string) {
	if query == "last" {
		ListAllFilesLast()
		return
	}
	idStart, idEnd, isRange := IsIDRange(query)
	if isRange {
		ListAllFilesIDRange(idStart, idEnd)
		return
	}
	ListAllFilesWithTags(query)
}

func ListAllFilesWithTags(tagsGrouped string) {
	ideas := quac.GetAllIdeas()
	subset := ideas.WithTags(idea.ParseClumpedTags(tagsGrouped))
	if len(subset) == 0 {
		fmt.Println("no ideas found with those tags")
		os.Exit(1)
	}
	for _, idea := range subset {
		quac.PrependLast(idea.Id)
		fmt.Println(idea.Filename)
	}
}

func ListAllFilesIDRange(idStart, idEnd uint32) {
	ideas := quac.GetAllIdeasInRange(idStart, idEnd)
	if len(ideas) == 0 {
		fmt.Println("no ideas found with in that range")
		os.Exit(1)
	}
	for _, idea := range ideas {
		quac.PrependLast(idea.Id)
		fmt.Println(idea.Filename)
	}
}

func ListAllFilesLast() {
	ids := quac.GetLastIDs()
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
		fmt.Println(path.Base(singleReturn))
		quac.Open(singleReturn)
		return
	}
	origBzFN, origBzContent := quac.GetOrigWorkingFileBytes()
	quac.OpenTextSplit(quac.WorkingFnsFile, quac.WorkingContentFile, maxFNLen)
	quac.SaveFromWorkingFiles(origBzFN, origBzContent)
}

func MultiOpenByRange(startId, endId uint32, forceSplitView bool) {
	found, maxFNLen, singleReturn := quac.WriteWorkingContentAndFilenamesFromRange(startId, endId, forceSplitView)
	if !found {
		fmt.Println("nothing found with those tags")
		return
	}
	// if only a single entry is found then quac.Open only it!
	if singleReturn != "" && !forceSplitView {
		fmt.Println(path.Base(singleReturn))
		quac.Open(singleReturn)
		return
	}
	origBzFN, origBzContent := quac.GetOrigWorkingFileBytes()
	quac.OpenTextSplit(quac.WorkingFnsFile, quac.WorkingContentFile, maxFNLen)
	quac.SaveFromWorkingFiles(origBzFN, origBzContent)
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
			quac.PrependLast(idea.Id)
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
	quac.PrependLast(idea.Id)
	quac.IncrementID()
}
