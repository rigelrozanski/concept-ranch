package quac

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/thranch/quac/idea"
)

func GetContentByID(id uint32) (content []byte, found bool) {
	filepath, found := GetFilepathByID(id)
	if !found {
		return content, false
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("problem reading filepath %v: %v\n", filepath, err)
		return content, false
	}
	return content, true
}

func GetFilepathByID(id uint32) (filepath string, found bool) {
	filename := GetFilenameByID(id)
	if filename == "" {
		return "", false
	}
	return path.Join(idea.IdeasDir, filename), true
}

func GetTrashcanFilepathsByID(id uint32) (currFilepath, trashCanFilePath string, found bool) {
	filename := GetFilenameByID(id)
	if filename == "" {
		return "", "", false
	}
	currFilepath = path.Join(idea.IdeasDir, filename)
	trashCanFilePath = path.Join(TrashCanDir, filename)
	return currFilepath, trashCanFilePath, true
}

func GetFilenameByID(id uint32) (fileName string) {
	files, err := ioutil.ReadDir(idea.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}

	idStr := idea.IdStr(id)
	for _, file := range files {
		fn := file.Name()
		if strings.HasPrefix(fn[2:], idStr) {
			fileName = fn
			break
		}
	}
	return fileName
}

func GetIdeaByID(id uint32, loglast bool) idea.Idea {
	fn := GetFilenameByID(id)
	return idea.NewIdeaFromFilename(fn, loglast)
}

func RemoveByID(id uint32) {
	existingFp, trashFp, found := GetTrashcanFilepathsByID(id)
	if !found {
		fmt.Println("nothing found at that ID")
		os.Exit(1)
	}
	if err := os.Rename(existingFp, trashFp); err != nil {
		log.Fatal(err)
	}
}

// copy an idea by the id
func CopyByID(id uint32) (newFilepath string) {
	fn := GetFilenameByID(id)
	newFilename := ReserveCopyFilename(fn, []string{})

	// perform the copy
	srcPath := path.Join(idea.IdeasDir, fn)
	writePath := path.Join(idea.IdeasDir, newFilename)
	cmn.Copy(srcPath, writePath)

	return writePath
}

func ReserveCopyFilename(oldFilename string, newTags []string) (newFilename string) {

	// remove the id, add in a new id
	idear := idea.NewIdeaFromFilename(oldFilename, true)
	idear.Id = idea.GetNextID()
	idear.Created = idea.TodayDate()
	idear.Tags = append(idear.Tags, newTags...)
	(&idear).UpdateFilename()

	idea.IncrementID()

	return idear.Filename
}
