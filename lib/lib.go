package lib

import (
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
var QiDir, IdeasDir, ConsumedDir, QiFile, LogFile, ConfigFile, WorkingFile string

// load config and set global file directories
func Init() {

	rootConfigPath := os.ExpandEnv("$HOME/.qi_config.txt")
	lines, err := cmn.ReadLines(rootConfigPath)
	if err != nil {
		panic(fmt.Sprintf("error reading ~/.qi_config.txt, error: %v", err))
	}

	QiDir = lines[0]
	IdeasDir = path.Join(QiDir, "ideas")
	ConsumedDir = path.Join(QiDir, "consumed")
	QiFile = path.Join(QiDir, "qi")
	LogFile = path.Join(QiDir, "log")
	ConfigFile = path.Join(QiDir, "config")
	WorkingFile = path.Join(QiDir, "working")

	EnsureBasics()
}

func EnsureBasics() {
	if !cmn.FileExists(QiDir) {
		panic("directory specified in ~/.qi_config.txt does not exist")
	}
	_ = os.Mkdir(IdeasDir, os.ModePerm)
	_ = os.Mkdir(ConsumedDir, os.ModePerm)
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
	if !cmn.FileExists(WorkingFile) {
		err := cmn.CreateEmptyFile(LogFile)
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

}

func GetNextID() uint32 {
	lines, err = cmn.ReadLines(ConfigFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config, error: %v", err))
	}
	IdCounter, err = uint32(strconv.Atoi(lines[0]))
	if err != nil {
		panic(fmt.Sprintf("error reading id_counter, error: %v", err))
	}
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
	idStr := strconv.Itoa(id)
	for _, file := range files {
		fn := file.Name()
		if strings.HasPrefix(idStr) {
			fileName = fn
			break
		}
	}
	if fileName == "" {
		return content, false
	}

	content, err = ioutil.ReadFile(path.Join(IdeasDir, fn))
	if err != nil {
		log.Fatal(err)
	}

	return content, true
}

func GetByTags(tags ...string) (content []byte, found bool) {
	ideas := PathToIdeas(IdeasDir)
	subset := idea.WithTags(tags)

	if len(subset) == 0 {
		return content, false
	}
	for _, idea := range subset {
		ideaContent, err = ioutil.ReadFile(path.Join(IdeasDir, idea.Filename))
		if err != nil {
			log.Fatal(err)
		}
		content += ideaContent
	}
	return content, found
}
