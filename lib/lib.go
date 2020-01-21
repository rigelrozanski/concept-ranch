package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	cmn "github.com/rigelrozanski/common"
)

// directory name where boards are stored in repo
var QiDir, IdeasDir, ConsumedDir, QiFile, LogFile, ConfigFile, WorkingFile string
var zeroDate time.Time

// load config and set global file directories
func Init() {

	zd, err := time.Parse(layout, "0")
	if err != nil {
		log.Fatal(err)
	}
	zeroDate = zd

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
	lines, err := cmn.ReadLines(ConfigFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config, error: %v", err))
	}
	count, err := strconv.Atoi(lines[0])
	if err != nil {
		panic(fmt.Sprintf("error reading id_counter, error: %v", err))
	}
	return uint32(count)
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

func GetByTags(tags ...string) (content []byte, found bool) {
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
	return content, found
}
