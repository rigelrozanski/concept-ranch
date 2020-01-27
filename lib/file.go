package lib

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

// create an empty file in the ideas Dir based on the filename
func WriteIdea(filename, entry string) {
	filepath := path.Join(IdeasDir, filename)
	err := ioutil.WriteFile(filepath, []byte(entry), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateEditedDateNow(updatePath string) {
	origFilename := path.Base(updatePath)
	idea := NewIdeaFromFilename(origFilename)
	idea.Edited = TodayDate()
	(&idea).UpdateFilename()
	origPath := path.Join(IdeasDir, origFilename)
	newPath := path.Join(IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}
