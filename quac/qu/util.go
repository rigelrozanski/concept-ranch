package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"

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
	var ideaImages quac.Ideas
	if err == nil {
		idear := quac.GetIdeaByID(consumed, true)
		if !(idear.IsImage() || idear.IsAudio()) {
			fmt.Println("this idea is not an image or audio cannot be transcribed")
			os.Exit(1)
		}
		ideaImages = []quac.Idea{idear}
	} else { // not an id, get by tags
		wot, _ := idea.NewTagWithout("DNT", "")
		ideaImages = quac.GetAllIdeasNonConsuming().
			WithImage().WithTags(wot)
		if optionalQuery != "" {
			ideaImages = ideaImages.WithTags(idea.ParseClumpedTags(optionalQuery))
			if len(ideaImages) == 0 {
				fmt.Println("no active images to transcribe with those tags")
				os.Exit(1)
			}
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
	fmt.Println("         - ADDTAG <newtag> to add the newtag to the entry")
	fmt.Println("         - KILLTAG <tag> to remove the existing tag to the entry")
	fmt.Println("         - KILL to delete the entry")
	fmt.Println("         - SKIP to skip")
	fmt.Println("         - QUIT to quit")
	fmt.Println("         - UNDO to undo the previous transcription")
	fmt.Println("-------------------------------------------------------------\n")

IdeaLoop:
	for _, idea := range ideaImages {

		quac.Open(idea.Path())

	GETINPUT:

		// read input from console
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		optionalEntry := consoleScanner.Text()
		optionalEntry = strings.TrimSpace(optionalEntry)

		switch {
		case optionalEntry == "DNT":
			quac.AddTagByIdea(&idea, "DNT")
			fmt.Println("ol'right never transcribing again!")
			continue
		case optionalEntry == "SKIP":
			fmt.Println("skipp'd")
			continue
		case optionalEntry == "KILL":
			quac.RemoveByID(idea.Id)
			fmt.Println("killed it")
			continue
		case strings.HasPrefix(optionalEntry, "ADDTAG "):
			newTag := strings.SplitN(optionalEntry, " ", 1)
			quac.AddTagByIdea(&idea, newTag[1])
			fmt.Printf("added the tag! new filename:\n%v\n", idea.Filename)
			fmt.Println("continue transcription:")
			goto GETINPUT
		case strings.HasPrefix(optionalEntry, "KILLTAG "):
			newTag := strings.SplitN(optionalEntry, " ", 2)
			quac.RemoveTagByIdea(&idea, newTag[1])
			fmt.Printf("removed the tag! new filename:\n%v\n", idea.Filename)
			fmt.Println("continue transcription:")
			goto GETINPUT
		case optionalEntry == "QUIT":
			break IdeaLoop
		case optionalEntry == "UNDO":
			panic("unimplemented")
		}

		consumerFilepath := quac.SetConsume(idea.Id, optionalEntry)
		if optionalEntry == "" {
			quac.OpenText(consumerFilepath)
		}
		fmt.Printf("created: %v\n", consumerFilepath)
	}
}

func TagUntagged() {
	untaggedIdeas := idea.GetAllIdeasNonConsuming().WithTag(idea.MustNewTagReg("UNTAGGED", ""))

	fmt.Println("             ~ Instructions ~")
	fmt.Println("enter desired tags seperated by spaces")
	fmt.Println("alternatively, enter:")
	fmt.Println("         - KILL to delete the entry")
	fmt.Println("         - SKIP to skip")
	fmt.Println("         - QUIT to quit")

	for _, idear := range untaggedIdeas {
		quac.View(idear.Path())

		// read input from console
		fmt.Println("Desired tags:")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		spacedTags := consoleScanner.Text()
		tagsStr := strings.Fields(spacedTags)

		if len(tagsStr) == 1 && tagsStr[0] == "SKIP" {
			fmt.Println("k, skip'd")
			continue
		}
		if len(tagsStr) == 1 && tagsStr[0] == "KILL" {
			quac.RemoveByID(idear.Id)
			fmt.Println("killed it")
			continue
		}
		if len(tagsStr) == 1 && tagsStr[0] == "QUIT" {
			fmt.Println("goodbye")
			break
		}

		origPath := idear.Path()

		// add the tags
		idear.Tags = idea.ParseStringTags(tagsStr)
		(&idear).UpdateFilename()

		// perform the file rename
		err := os.Rename(origPath, idear.Path())
		fmt.Printf("retagged to:\n%v\n", idear.Filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WaterCloset() {
	if idea.TagUsedInNonConsuming("UNTAGGED") {
		TagUntagged()
	} else {
		Transcribe("")
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

func NewEmptyEntry(clumpedTags string) {
	idear := quac.NewNonConsumingTextIdea(clumpedTags)
	writePath := path.Join(quac.IdeasDir, idear.Filename)
	quac.IncrementID()
	fmt.Printf("created: %v\n", writePath)
	quac.OpenText(writePath)
}

func ManualEntry(commonTagsClumped string) {

	fmt.Println("_________________________________________")
	fmt.Println("INTERACTIVE MANUAL TRANSCRIPTION")
	fmt.Println(" - tags to be entered seperated by spaces")
	fmt.Println(" - quit with ctrl-c, or by typing QUIT during tag entry")
	fmt.Println("")
	for {

		fmt.Print("enter tags:\t")
		consoleScanner := bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		newStrTags := strings.Fields(consoleScanner.Text())
		if len(newStrTags) > 0 && newStrTags[0] == "QUIT" {
			break
		}
		newClumpedTags := strings.Join(newStrTags, ",")
		clumpedTags := quac.CombineClumpedTags(commonTagsClumped, newClumpedTags)

		fmt.Print("enter entry:\t")
		consoleScanner = bufio.NewScanner(os.Stdin)
		_ = consoleScanner.Scan()
		entry := consoleScanner.Text()

		Entry(entry, clumpedTags)
	}
}

func SetEncryption(idStr string) {
	id, err := quac.ParseID(idStr)
	if err != nil {
		log.Fatalf("error parsing id, error: %v", err)
	}

	quac.SetEncryptionById(id)
}

func QuickEntry(clumpedTags, entry string) {
	Entry(entry, clumpedTags)
}

func MultiOpen(unsplitTagsOrID string, forceSplitView bool) {

	if unsplitTagsOrID == "" {
		return
	}

	startID, endID, isRange := IsIDorIDRange(unsplitTagsOrID)
	if isRange {
		quac.MultiOpenByRange(startID, endID, forceSplitView)
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
	quac.MultiOpenByTags(splitTags, forceSplitView)
}

func OpenWorking() {
	quac.OpenTextSplit(quac.WorkingFnsFile, quac.WorkingContentFile, 50)
}

func SaveWorking() {
	quac.SaveFromWorkingFiles([]byte{}, []byte{})
}

func parseIdStr(idStr string) (uint32, error) {
	idI, err := quac.ParseID(idStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing id, error: %v", err)
	}
	return uint32(idI), nil
}

func RemoveByID(idOrIds string) {
	startID, endID, valid := IsIDorIDRange(idOrIds)
	if valid {
		RemoveAcrossIDs(startID, endID)
		fmt.Println("roger, moved to the trash-can")
	} else {
		fmt.Println("invalid remove range")
	}
}

func EmptyTrash() {
	if len(quac.TrashCanDir) < 5 { // NOTE vital, don't want to delete the root
		panic("TrashCanDir not set!")
	}
	names, err := ioutil.ReadDir(quac.TrashCanDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range names {
		idea.ValidateFilenameAsIdea(n.Name()) // panic on things which are not ideas

		fp := path.Join(quac.TrashCanDir, n.Name())
		err := os.Remove(fp)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("trash can emptied into the void")
}

func RemoveAcrossIDs(id1, id2 uint32) {
	for i := id1; i <= id2; i++ {
		quac.RemoveByID(i)
	}
}

func CopyByID(idStr string) {
	id, err := parseIdStr(idStr)
	if err != nil {
		log.Fatal(err)
	}
	quac.Open(quac.CopyByID(id))
}

// NOTE allows for reversed inputs
func RemoveTagByID(idStr, tagToRemove string) {
	id, err := parseIdStr(idStr)
	if err != nil {
		// swap inputs
		id, err = parseIdStr(tagToRemove)
		if err != nil {
			log.Fatal(err)
		}
		tagToRemove = idStr
	}

	idea := quac.GetIdeaByID(id, true)
	quac.RemoveTagByIdea(&idea, tagToRemove)
}

// NOTE allows for reversed inputs
func AddTagByID(idStr, tagToAdd string) {
	id, err := parseIdStr(idStr)
	if err != nil {
		// swap inputs
		id, err = parseIdStr(tagToAdd)
		if err != nil {
			log.Fatal(err)
		}
		tagToAdd = idStr
	}
	idea := quac.GetIdeaByID(id, true)
	quac.AddTagByIdea(&idea, tagToAdd)
}

func AddTagToMany(tagToAdd, manyTagsClumped string) {

	manyTags := idea.ParseClumpedTags(manyTagsClumped)
	ideas := quac.GetAllIdeas().WithAnyOfTags(manyTags)
	for _, idea := range ideas {
		quac.AddTagByIdea(&idea, tagToAdd)
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
		idea := quac.NewIdeaFromFilename(origFn, true)
		fromTag := quac.ParseFirstTagFromString(from)
		toTag := quac.ParseFirstTagFromString(to)
		(&idea).RenameTag(fromTag, toTag)
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
		idea := quac.NewIdeaFromFilename(origFn, true)
		(&idea).RemoveTags(quac.ParseTagFromString(tag))
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

func ListAllTagsWithTags(clumpedTags string) {
	ideas := quac.GetAllIdeas()
	queryTags := idea.ParseClumpedTags(clumpedTags)
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
			outTags[i] = uTag.String()
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

func IsIDorIDRange(query string) (idStart, idEnd uint32, isIDorIDRange bool) {
	id, err := strconv.Atoi(query)
	if err == nil {
		return uint32(id), uint32(id), true
	}

	sp := strings.Split(query, "-")
	if len(sp) != 2 {
		return 0, 0, false
	}

	idStart, err = quac.ParseID(sp[0])
	if err != nil {
		return 0, 0, false
	}
	idEnd, err = quac.ParseID(sp[1])
	if err != nil {
		return 0, 0, false
	}

	return idStart, idEnd, true
}

func ListAllFilesByLocation() {
	ideas := quac.GetAllIdeas()
	if len(ideas) == 0 {
		fmt.Println("no ideas found")
	}
	for _, idea := range ideas {
		fmt.Println(idea.Path())
	}
}

func ListSelectAllFilesWithQueryNoLast(query string) {
	if query == "last" {
		MultiOpen("last", false)
		return
	}
	refined := RefineQueryCUI(query)
	MultiOpen(refined, false)
}

func ListSelectAllFilesWithQuery(query string) {
	refined := RefineQueryCUI(query)
	MultiOpen(refined, false)
}

// in command line CUI for selecting based on the query provided
func RefineQueryCUI(query string) (refinedQuery string) {

	// get query ideas
	var ideas idea.Ideas
	idStart, idEnd, isRange := IsIDorIDRange(query)
	switch {
	case isRange:
		ideas = quac.GetAllIdeas().InRange(idStart, idEnd)
	case query == "last":
		ids := quac.GetLastIDs()
		for _, id := range ids {
			ideas = append(ideas, quac.GetIdeaByID(id, false))
		}
	default:
		ideas = quac.GetAllIdeas().WithTags(idea.ParseClumpedTags(query))
	}

	// skip this process if there is only one entry (or none)
	switch len(ideas) {
	case 0:
		fmt.Printf("nothing found for the query %v\n", query)
		return ""
	case 1:
		return query
	}

	termWidth := 100
	if term.IsTerminal(0) {
		var err error
		termWidth, _, err = term.GetSize(0)
		if err != nil {
			log.Fatal(err)
		}
	}

	// for fitting all the ideas into a small space
	maxDisplayIdeas := 30
	displayedIdeasStartIndex := 0
	displayedIdeasList := []string{">>> OPEN ALL SIMULTANEOUSLY"}
	for _, idea := range ideas {

		// Trucate the size if it's too big for the term
		// TODO wrap it and account for the diff
		toAdd := idea.Filename
		if len([]rune(toAdd)) >= termWidth {
			// NOTE need to subtract so much more than the
			// actual width because of invisible characters
			// used for the colouring of the line
			toAdd = string([]rune(toAdd)[:termWidth-14])
			toAdd += "..."
		}
		displayedIdeasList = append(displayedIdeasList, toAdd)
	}

	colorBlue := "\033[34m"
	colorReset := "\033[0m"
	lenList := len(ideas)
	if lenList > maxDisplayIdeas {
		lenList = maxDisplayIdeas
	}
	lenList += 2 // +1  for the "open all" line, +1 for usage line
	selected := 0

	fmt.Print("\033[?25l") // hide cursor
	defer func() {
		fmt.Print("\033[?25h") // show cursor

		// move the cursor to the bottom of the list
		fmt.Print(fmt.Sprintf("\033[%vB", lenList))
	}()

	fmt.Printf("\n")
	toPrintPrevious := []string{}
	for {
		toPrint := []string{"navigation: 'j'/'k', select: 'enter', escape: 'esc'/'q'", ""}

		startOffset, endOffset := 0, 0 // need this to keep everything with a constant maxDisplayIdeas
		displayStartExpand, displayEndExpand := false, false
		if displayedIdeasStartIndex > 0 {
			displayStartExpand = true
			startOffset = 1
			toPrint = append(toPrint, " ...")
		}
		if displayedIdeasStartIndex+maxDisplayIdeas < len(displayedIdeasList) {
			displayEndExpand = true
			endOffset = 1
		}

		for i := displayedIdeasStartIndex + startOffset; i < displayedIdeasStartIndex+maxDisplayIdeas-endOffset; i++ {
			if i >= len(displayedIdeasList) {
				break
			}

			textToDisplay := displayedIdeasList[i]
			if i == selected {
				toPrint = append(toPrint, " "+string(colorBlue)+textToDisplay+string(colorReset))
			} else {
				toPrint = append(toPrint, " "+textToDisplay)
			}
		}
		if displayEndExpand {
			toPrint = append(toPrint, " ...")
		}

		//whiteout the previous lines
		for _, tpl := range toPrintPrevious {
			fmt.Printf("%"+fmt.Sprintf("%v", len(tpl))+"x\n", "")
		}
		fmt.Print(fmt.Sprintf("\033[%vA", len(toPrintPrevious))) // move the cursor back up for the next print
		toPrintPrevious = toPrint

		//print these lines
		for _, tp := range toPrint {
			fmt.Printf("%s\n", tp)
		}
		fmt.Print(fmt.Sprintf("\033[%vA", len(toPrint))) // move the cursor back up for the next print

		// scan for user input
		ch, err := GetCharJustHit()
		if err != nil {
			return ""
		}

		// TODO get arrow keys to work
		switch {
		// up or left (arrow keys and vim)
		case ch == 'k' || ch == 'h':
			if selected > 0 {
				selected--
				if displayStartExpand && selected < displayedIdeasStartIndex+startOffset {
					displayedIdeasStartIndex = selected - startOffset
				}
			}
			// down or right (arrow keys and vim)
		case ch == 'j' || ch == 'l':
			if selected < len(ideas) {
				selected++
				if displayEndExpand && selected > displayedIdeasStartIndex+maxDisplayIdeas-endOffset-1 {
					displayedIdeasStartIndex = selected - maxDisplayIdeas + endOffset + 1
				}
			}
		case ch == '\r': // enter
			fmt.Printf("\n")
			switch selected {
			case 0:
				return query // "all" of the original query
			default:
				return fmt.Sprintf("%v", ideas[selected-1].Id)
			}
		case ch == '\033' || ch == 'q' || ch == 'c': // escape key or 'q'
			fmt.Printf("\n")
			return ""
		}
	}
}

func GetCharJustHit() (rune, error) {
	// switch stdin into 'raw' mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return ' ', err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return ' ', err
	}
	return rune(b[0]), nil
}

func ListAllFilesWithQuery(query string) {
	listAllFilesWithQuery(query, false)
}

func ListAllFilesByLocationWithQuery(query string) {
	listAllFilesWithQuery(query, true)
}

func listAllFilesWithQuery(query string, showFilepath bool) {
	if query == "last" {
		ListAllFilesLast(showFilepath)
		return
	}
	idStart, idEnd, isRange := IsIDorIDRange(query)
	if isRange {
		ListAllFilesIDRange(idStart, idEnd, showFilepath)
		return
	}
	ListAllFilesWithTags(query, showFilepath)
}

func ListAllFilesWithTags(tagsGrouped string, showFilepath bool) {
	ideas := quac.GetAllIdeas()
	subset := ideas.WithTags(idea.ParseClumpedTags(tagsGrouped))
	if len(subset) == 0 {
		fmt.Println("no ideas found with those tags")
		os.Exit(1)
	}
	for _, idea := range subset {
		quac.PrependLast(idea.Id)
		if showFilepath {
			fmt.Println(idea.Path())
		} else {
			fmt.Println(idea.Filename)
		}
	}
}

func ListAllFilesIDRange(idStart, idEnd uint32, showFilepath bool) {
	ideas := quac.GetAllIdeas()
	subset := ideas.InRange(idStart, idEnd)
	if len(ideas) == 0 {
		fmt.Println("no ideas found with in that range")
		os.Exit(1)
	}
	for _, idea := range subset {
		quac.PrependLast(idea.Id)
		if showFilepath {
			fmt.Println(idea.Path())
		} else {
			fmt.Println(idea.Filename)
		}
	}
}

func ListAllFilesLast(showFilepath bool) {
	ids := quac.GetLastIDs()
	for _, id := range ids {
		if showFilepath {
			fmt.Println(quac.GetFilenameByID(id))
		} else {
			fmt.Println(quac.GetIdeaByID(id, false).Path())
		}
	}
}

func ViewByID(id uint32) {
	content, found := quac.GetContentByID(id)
	if !found {
		fmt.Println("nothing found with that id")
	}
	fmt.Printf("%s\n", content)
}

func ViewByTags(tags []idea.Tag) {
	content, found := quac.ConcatAllContentFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	fmt.Printf("%s\n", content)
}

// TODO rewrite this function, it's confusing as heck
// TODO this logic should exist in the library
// create an entry
func Entry(entryOrPath string, clumpedTags string) {

	hasFN := false
	if strings.Contains(clumpedTags, "FILENAME") {
		hasFN = true
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
			if strings.Contains(clumpedTags, "FILENAME") {
				filebase := strings.TrimSuffix(path.Base(filepath), path.Ext(filepath))
				clumpedTags = strings.Replace(clumpedTags, "FILENAME", filebase, 2)
			}

			idea := quac.NewIdeaFromFile(clumpedTags, filepath)
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

	idea := quac.NewNonConsumingTextIdea(clumpedTags)
	err := cmn.WriteLines([]string{entryOrPath}, idea.Path())
	if err != nil {
		log.Fatalf("error writing new file: %v", err)
	}
	quac.PrependLast(idea.Id)
	quac.IncrementID()
}
