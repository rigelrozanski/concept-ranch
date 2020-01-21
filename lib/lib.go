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

func GetByID(id uint32) (content []byte, found bool) {
	files, err := ioutil.ReadDir(IdeasDir)
	if err != nil {
		log.Fatal(err)
	}

	fileName := ""
	idStr := strconv.Itoa(int(id))
	for _, file := range files {
		fn := file.Name()
		if strings.HasPrefix(fn, idStr) {
			fileName = fn
			break
		}
	}
	if fileName == "" {
		return content, false
	}

	content, err = ioutil.ReadFile(path.Join(IdeasDir, fileName))
	if err != nil {
		log.Fatal(err)
	}

	return content, true
}

func ConcatAllContentFromTags(dir string, tags []string) (content []byte, found bool) {
	ideas := PathToIdeas(dir)
	subset := ideas.WithTags(tags)

	if len(subset) == 0 {
		return content, false
	}
	for _, idea := range subset {
		ideaContent, err := ioutil.ReadFile(path.Join(dir, idea.Filename))
		if err != nil {
			log.Fatal(err)
		}
		content = append(content, ideaContent...)
	}
	return content, true
}

func WriteWorkingContentAndFilenamesFromTags(dir string, tags []string) (found bool, maxFNLen int) {
	ideas := PathToIdeas(dir)
	subset := ideas.WithTags(tags)

	if len(subset) == 0 {
		return false, 0
	}

	var contentBz, fnBz []byte
	for _, idea := range subset {
		icontentBz, err := ioutil.ReadFile(path.Join(dir, idea.Filename))
		if err != nil {
			log.Fatal(err)
		}

		noLines := bytes.Count(icontentBz, []byte{'\n'})

		if len(idea.Filename)+2 > maxFNLen {
			maxFNLen = len(idea.Filename) + 2
		}
		fnBz = append(fnBz, []byte("["+idea.Filename+"]"+strings.Repeat("\n", noLines))...)
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
	return true, maxFNLen
}
