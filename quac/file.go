package quac

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/rigelrozanski/thranch/quac/idea"
)

// create an empty file in the ideas Dir based on the filename
func WriteIdea(filename, entry string) {
	filepath := path.Join(idea.IdeasDir, filename)
	err := ioutil.WriteFile(filepath, []byte(entry), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateEditedDateNow(updatePath string) {
	origFilename := path.Base(updatePath)
	idear := idea.NewIdeaFromFilename(origFilename, true)
	idear.Edited = idea.TodayDate()
	(&idear).UpdateFilename()
	origPath := path.Join(idea.IdeasDir, origFilename)
	newPath := path.Join(idea.IdeasDir, idear.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateFilepathToEncrypted(Path string) string {
	err := os.Rename(Path, Path+".en")
	if err != nil {
		log.Fatal(err)
	}
	return Path + ".en"
}
