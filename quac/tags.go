package quac

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/rigelrozanski/thranch/quac/idea"
)

func RemoveTagByIdea(idea *Idea, tagToRemove string) {
	origFilename := (*idea).Filename
	idea.RemoveTag(ParseTagFromString(tagToRemove))
	idea.UpdateFilename()
	origPath := path.Join(IdeasDir, origFilename)
	newPath := path.Join(IdeasDir, (*idea).Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func AddTagByIdea(idea *Idea, tagToAdd string) {
	origFilename := (*idea).Filename

	idea.AddTag(ParseTagFromString(tagToAdd))
	idea.UpdateFilename()
	origPath := path.Join(IdeasDir, origFilename)
	newPath := path.Join(IdeasDir, (*idea).Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func MultiOpenByTags(tags []idea.Tag, forceSplitView bool) {
	found, maxFNLen, singleReturn :=
		WriteWorkingContentAndFilenamesFromTags(tags, forceSplitView)
	if !found {
		fmt.Println("nothing found with those tags")
		return
	}
	// if only a single entry is found then Open only it!
	if singleReturn != "" && !forceSplitView {
		fmt.Println(path.Base(singleReturn))
		Open(singleReturn)
		return
	}
	origBzFN, origBzContent := GetOrigWorkingFileBytes()
	OpenTextSplit(WorkingFnsFile, WorkingContentFile, maxFNLen)
	SaveFromWorkingFiles(origBzFN, origBzContent)
}

func MultiOpenByRange(startId, endId uint32, forceSplitView bool) {
	found, maxFNLen, singleReturn :=
		WriteWorkingContentAndFilenamesFromRange(startId, endId, forceSplitView)
	if !found {
		fmt.Println("nothing found with those tags")
		return
	}
	// if only a single entry is found then Open only it!
	if singleReturn != "" && !forceSplitView {
		fmt.Println(path.Base(singleReturn))
		Open(singleReturn)
		return
	}
	origBzFN, origBzContent := GetOrigWorkingFileBytes()
	OpenTextSplit(WorkingFnsFile, WorkingContentFile, maxFNLen)
	SaveFromWorkingFiles(origBzFN, origBzContent)
}
