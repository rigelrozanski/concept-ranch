package quac

import (
	"fmt"
	"path"
)

func MultiOpenByTags(tags []string, forceSplitView bool) {
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
