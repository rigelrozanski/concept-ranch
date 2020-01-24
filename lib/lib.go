package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	cmn "github.com/rigelrozanski/common"
)

// directory name where boards are stored in repo
var QiDir, IdeasDir, QiFile, LogFile, ConfigFile, WorkingFnsFile, WorkingContentFile, LastIdFile string

// load config and set global file directories
func init() {

	rootConfigPath := os.ExpandEnv("$HOME/.qi_config.txt")
	lines, err := cmn.ReadLines(rootConfigPath)
	if err != nil {
		panic(fmt.Sprintf("error reading ~/.qi_config.txt, error: %v", err))
	}

	QiDir = lines[0]
	IdeasDir = path.Join(QiDir, "ideas")
	QiFile = path.Join(QiDir, "qi")
	LogFile = path.Join(QiDir, "log")
	ConfigFile = path.Join(QiDir, "config")
	WorkingFnsFile = path.Join(QiDir, "working_files")
	WorkingContentFile = path.Join(QiDir, "working_content")
	LastIdFile = path.Join(QiDir, "last")

	EnsureBasics()
}

func EnsureBasics() {
	if !cmn.FileExists(QiDir) {
		panic("directory specified in ~/.qi_config.txt does not exist")
	}
	_ = os.Mkdir(IdeasDir, os.ModePerm)
	if !cmn.FileExists(QiFile) {
		err := cmn.CreateEmptyFile(QiFile)
		if err != nil {
			panic(err)
		}
	}
	if !cmn.FileExists(LogFile) {
		err := cmn.CreateEmptyFile(LogFile)
		if err != nil {
			panic(err)
		}
	}
	if !cmn.FileExists(WorkingFnsFile) {
		err := cmn.CreateEmptyFile(WorkingFnsFile)
		if err != nil {
			panic(err)
		}
	}
	if !cmn.FileExists(WorkingContentFile) {
		err := cmn.CreateEmptyFile(WorkingContentFile)
		if err != nil {
			panic(err)
		}
	}
	if !cmn.FileExists(ConfigFile) {
		err := cmn.WriteLines([]string{"000001"}, ConfigFile)
		if err != nil {
			panic(err)
		}
	}
	if !cmn.FileExists(LastIdFile) {
		err := cmn.WriteLines([]string{"000000"}, LastIdFile)
		if err != nil {
			panic(err)
		}
	}
}

func GetNextID() uint32 {
	lines, err := cmn.ReadLines(ConfigFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config, error: %v", err))
	}
	count, err := strconv.Atoi(lines[0])
	if err != nil {
		panic(fmt.Sprintf("error reading id_counter, error: %v", err))
	}
	return uint32(count + 1)
}

func IncrementID() {
	err := cmn.WriteLines([]string{strconv.Itoa(int(GetNextID()))}, ConfigFile)
	if err != nil {
		panic(err)
	}
}

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
	return path.Join(IdeasDir, filename)
}

func GetFilenameByID(id uint32) (filepath string) {
	files, err := ioutil.ReadDir(IdeasDir)
	if err != nil {
		log.Fatal(err)
	}

	fileName := ""
	idStr := strconv.Itoa(int(id))
	for _, file := range files {
		fn := file.Name()
		if strings.HasPrefix(fn[2:], idStr) {
			fileName = fn
			break
		}
	}
	return fileName
}

func GetIdeaByID(id uint32) Idea {
	fn := GetFilenameByID(id)
	return NewIdeaFromFilename(fn)

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

// create an empty file in the ideas Dir based on the filename
func WriteIdea(filename, entry string) {
	filepath := path.Join(IdeasDir, filename)
	err := ioutil.WriteFile(filepath, []byte(entry), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
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

	// remove the id, add in a new id
	newID := GetNextID()
	split := strings.SplitN(fn, ",", 3)

	newFilename := strings.Join([]string{
		split[0], strconv.Itoa(int(newID)), split[2]}, ",")
	IncrementID()

	// actually perform the copy
	srcPath := path.Join(IdeasDir, fn)
	writePath := path.Join(IdeasDir, newFilename)
	cmn.Copy(srcPath, writePath)

	return writePath
}

// copy an idea by the id
func Zombie(zombieId uint32) {
	consumedIdea := GetIdeaByID(zombieId)
	consumedIdea.SetZombie()
}

// display the lineage of the
func Lineage(id uint32) (compiled string) {
	lineageIdea := GetIdeaByID(id)
	fmt.Printf("debug lineageIdea: %v\n", lineageIdea)
	fmt.Printf("debug ConsumesIds: %v\n", lineageIdea.ConsumesIds)
	for _, consume := range lineageIdea.ConsumesIds {
		fmt.Printf("debug consume: %v\n", consume)
		fn := GetFilenameByID(consume)
		content, found := GetContentByID(consume)
		if !found {
			log.Fatalf("child not found: %v", consume)
		}
		compiled = fmt.Sprintf("%v\n%v\n%s", compiled, fn, content)
	}
	return compiled
}

// copy an idea by the id
func Consume(consumedId uint32, entry string) (consumerFilepath string) {
	consumedIdea := GetIdeaByID(consumedId)

	// consumer: remove the id, add in a new id, add the consumes id
	consumerIdea := NewConsumingIdea(consumedIdea)
	IncrementID()
	WriteIdea(consumerIdea.Filename, entry)

	consumedIdea.SetConsumed()
	return consumerIdea.Path()
}

func Consumes(consumedId, consumesId uint32) {
	consumedIdea := GetIdeaByID(consumedId)
	consumesIdea := GetIdeaByID(consumesId)

	// consumer: remove the id, add in a new id, add the consumes id
	consumesIdea.ConsumesIds = append(consumesIdea.ConsumesIds, consumedId)
	srcPath := consumesIdea.Path()
	(&consumesIdea).UpdateFilename()
	writePath := path.Join(IdeasDir, consumesIdea.Filename)
	err := os.Rename(srcPath, writePath)
	if err != nil {
		log.Fatal(err)
	}

	consumedIdea.SetConsumed()
}

func ConcatAllContentFromTags(tags []string) (content []byte, found bool) {
	ideas := PathToIdeas(IdeasDir)
	subset := ideas.WithTags(tags)

	if len(subset) == 0 {
		return content, false
	}
	for _, idea := range subset {
		ideaContent, err := ioutil.ReadFile(path.Join(IdeasDir, idea.Filename))
		if err != nil {
			log.Fatal(err)
		}
		content = append(content, ideaContent...)
	}
	return content, true
}

func WriteWorkingContentAndFilenamesFromTags(tags []string) (found bool, maxFNLen int, singleReturn string) {
	ideas := PathToIdeas(IdeasDir)
	subset := ideas.WithTags(tags)

	switch len(subset) {
	case 0:
		return false, 0, ""
	case 1:
		// if only one found, return its path
		return true, 1, path.Join(QiDir, subset[0].Filename)
	default:
		// write working contents and filenames from tags
		var contentBz, fnBz []byte
		for _, idea := range subset {
			icontentBz, err := ioutil.ReadFile(path.Join(IdeasDir, idea.Filename))
			if err != nil {
				log.Fatal(err)
			}

			noLines := bytes.Count(icontentBz, []byte{'\n'})

			if len(idea.Filename)+2 > maxFNLen {
				maxFNLen = len(idea.Filename) + 2
			}
			fnBz = append(fnBz, []byte(idea.Filename+strings.Repeat("\n", noLines))...)
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

func SaveFromWorkingFiles() {
	fnLines, err := cmn.ReadLines(WorkingFnsFile)
	if err != nil {
		log.Fatal(err)
	}
	contentLines, err := cmn.ReadLines(WorkingContentFile)
	if err != nil {
		log.Fatal(err)
	}

	for start, fnLine := range fnLines {
		if fnLine == "" {
			continue
		}
		end := start + 1
		// keep adding to end unless the next line is not empty
		// or is not a part of the array!
		for ; !(end >= len(fnLines) || fnLines[end] != ""); end++ {
		}

		// get id and orig bytes
		id := GetIdByFilename(fnLine)
		origBz, found := GetContentByID(id)
		if !found {
			log.Fatal("not found when should be")
		}

		// remove the old file by id (may have been renamed)
		RemoveByID(id)

		// create the new file
		filepath := path.Join(IdeasDir, fnLine)
		err := cmn.WriteLines(contentLines[start:end], filepath)
		if err != nil {
			log.Fatal(err)
		}

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
