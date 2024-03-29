package quac

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/thranch/quac/idea"
)

const SPLIT = "SPLIT"

func WriteWorkingContentAndFilenamesFromRange(startId, endId uint32,
	forceSplitView bool) (found bool, maxFNLen int, singleReturn string) {

	ideas := GetAllIdeas().InRange(startId, endId)
	return WriteWorkingContentAndFilenamesFromIdeas(ideas, forceSplitView)
}

func WriteWorkingContentAndFilenamesFromTags(tags []idea.Tag, forceSplitView bool) (
	found bool, maxFNLen int, singleReturn string) {

	ideas := idea.GetAllIdeasNonConsuming().WithTags(tags)
	return WriteWorkingContentAndFilenamesFromIdeas(ideas, forceSplitView)
}

func WriteWorkingContentAndFilenamesFromIdeas(idears idea.Ideas,
	forceSplitView bool) (found bool, maxFNLen int, singleReturn string) {

	switch len(idears) {
	case 0:
		return false, 0, ""
	case 1:
		if !forceSplitView {
			// if only one found, return its path
			return true, 1, idears[0].Path()
		}
		fallthrough
	default:
		// write working contents and filenames from tags
		var contentBz, fnBz []byte
		for _, idear := range idears {
			// TODO de-dup code from below
			if !idear.IsText() {
				continue
			}
			icontentBz, err := ioutil.ReadFile(idear.Path())
			if err != nil {
				log.Fatal(err)
			}

			noLines := bytes.Count(icontentBz, []byte{'\n'})

			if len(idear.Filename)+2 > maxFNLen {
				maxFNLen = len(idear.Filename) + 2
			}
			fnBz = append(fnBz, []byte(idear.Filename+strings.Repeat("\n", noLines))...)
			contentBz = append(contentBz, icontentBz...)
		}

		err := ioutil.WriteFile(WorkingFnsFile, fnBz, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(WorkingContentFile, contentBz, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		return true, maxFNLen, ""
	}
}

func WriteWorkingContentAndFilenamesFromFilePath(filePath string) (maxFNLen int) {
	idear := idea.NewIdeaFromFilepath(filePath, true)

	// write working contents and filenames from tags
	var contentBz, fnBz []byte
	if !idear.IsText() {
		log.Fatal("file at this idea is not text")
	}
	icontentBz, err := ioutil.ReadFile(idear.Path())
	if err != nil {
		log.Fatal(err)
	}

	noLines := bytes.Count(icontentBz, []byte{'\n'})

	if len(idear.Filename)+2 > maxFNLen {
		maxFNLen = len(idear.Filename) + 2
	}
	fnBz = append(fnBz, []byte(idear.Filename+strings.Repeat("\n", noLines))...)
	contentBz = append(contentBz, icontentBz...)

	err = ioutil.WriteFile(WorkingFnsFile, fnBz, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(WorkingContentFile, contentBz, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return maxFNLen
}

// get the bytes of the working files (original)
func GetOrigWorkingFileBytes() (origBzFN, origBzContent []byte) {

	// do not save if no modifications have been made
	origBzFN, err := ioutil.ReadFile(WorkingFnsFile)
	if err != nil {
		panic(err)
	}
	origBzContent, err = ioutil.ReadFile(WorkingContentFile)
	if err != nil {
		panic(err)
	}

	return origBzFN, origBzContent
}

func SaveFromWorkingFiles(origBzFN, origBzContent []byte) {

	// do not save if no modifications have been made
	finalBzFN, err := ioutil.ReadFile(WorkingFnsFile)
	if err != nil {
		panic(err)
	}
	finalBzContent, err := ioutil.ReadFile(WorkingContentFile)
	if err != nil {
		panic(err)
	}
	if bytes.Compare(origBzFN, finalBzFN) == 0 &&
		bytes.Compare(origBzContent, finalBzContent) == 0 {
		return
	}

	fnLines, err := cmn.ReadLines(WorkingFnsFile)
	if err != nil {
		log.Fatal(err)
	}
	contentLines, err := cmn.ReadLines(WorkingContentFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(fnLines) != len(contentLines) {
		fmt.Println("Error! unequal number of lines in working files!" +
			" Correct manually with cmd: qu open-working")
		os.Exit(1)
	}

	var recentFileName string
	for startRange, fnLine := range fnLines {
		if fnLine == "" {
			continue
		}
		splitFile := false
		if strings.HasPrefix(fnLine, SPLIT) {
			splitFile = true
		} else {
			recentFileName = fnLine
		}

		endRange := startRange + 1
		// keep adding to end unless the next line is not empty
		// or is not a part of the array!
		for ; !(endRange >= len(fnLines) || fnLines[endRange] != ""); endRange++ {
		}

		var origBz []byte
		var filepath string

		if splitFile {
			if recentFileName == "" {
				log.Fatal("cannot split with no prior filename")
			}
			potentialTags := strings.TrimSpace(
				strings.TrimPrefix(fnLine, SPLIT))
			filename := ReserveCopyFilename(recentFileName, potentialTags)

			// create the split filepath but change the id
			filepath = path.Join(idea.IdeasDir, filename)

		} else {
			// get the orig bytes (non existant if a split)
			id, _ := idea.GetIdByFilename(fnLine)
			found := false
			origBz, found = GetContentByID(id)
			if !found {
				log.Fatal("not found when should be")
			}

			// remove the old file by id (may have been renamed)
			RemoveByID(id)

			// create the new file
			filepath = path.Join(idea.IdeasDir, fnLine)
		}

		// do not write the file if there is no content
		if endRange-startRange == 1 &&
			strings.TrimSpace(contentLines[startRange]) == "" {
			continue
		}

		// write the file
		err := cmn.WriteLines(contentLines[startRange:endRange], filepath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Split this out: %v\n", filepath)

		// check the content and possibly mark as edited
		finalBz, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		if bytes.Compare(origBz, finalBz) != 0 {
			UpdateEditedDateNow(filepath)
		}
	}
}
