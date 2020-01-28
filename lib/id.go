package lib

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/qi/lib/idea"
)

func GetContentByID(id uint32) (content []byte, found bool) {
	filepath := GetFilepathByID(id)
	if filepath == "" {
		return content, false
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return content, true
}

func GetFilepathByID(id uint32) (filepath string) {
	filename := GetFilenameByID(id)
	if filename == "" {
		return ""
	}
	return path.Join(idea.IdeasDir, filename)
}

func GetFilenameByID(id uint32) (filepath string) {
	files, err := ioutil.ReadDir(idea.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}

	fileName := ""
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

func GetIdeaByID(id uint32) idea.Idea {
	fn := GetFilenameByID(id)
	return idea.NewIdeaFromFilename(fn)

}

func GetIdByFilename(filename string) (id uint32) {
	split := strings.SplitN(filename, ",", 3)
	if len(split) != 3 {
		log.Fatal("bad split")
	}
	idI, err := strconv.Atoi(split[1])
	if err != nil {
		log.Fatal(err)
	}
	return uint32(idI)
}

func RemoveByID(id uint32) {
	fp := GetFilepathByID(id)
	err := os.Remove(fp)
	if err != nil {
		log.Fatal(err)
	}
}

// copy an idea by the id
func CopyByID(id uint32) (newFilepath string) {
	fn := GetFilenameByID(id)
	newFilename := ReserveCopyFilename(fn)

	// perform the copy
	srcPath := path.Join(idea.IdeasDir, fn)
	writePath := path.Join(idea.IdeasDir, newFilename)
	cmn.Copy(srcPath, writePath)

	return writePath
}

func ReserveCopyFilename(oldFilename string) (newFilename string) {

	// remove the id, add in a new id
	newID := idea.GetNextID()
	split := strings.SplitN(oldFilename, ",", 3)

	newFilename = strings.Join([]string{
		split[0], idea.IdStr(newID), split[2]}, ",")
	idea.IncrementID()

	return newFilename
}
