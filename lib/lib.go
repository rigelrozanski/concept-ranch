package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	pathL "path"
	"strconv"
	"strings"

	cmn "github.com/rigelrozanski/common"
)

// directory name where boards are stored in repo
var QiDir, IdeasDir, ConsumedDir, QiFile, LogFile, ConfigFile string
var IdCounter uint32

// load config and set global file directories
func Init() {

	rootConfigPath := os.ExpandEnv("$HOME/.qi_config.txt")
	lines, err := cmn.ReadLines(rootConfigPath)
	if err != nil {
		panic(fmt.Sprintf("error reading ~/.qi_config.txt, error: %v", err))
	}

	QiDir = lines[0]
	IdeasDir = pathL.Join(QiDir, "ideas")
	ConsumedDir = pathL.Join(QiDir, "consumed")
	QiFile = pathL.Join(QiDir, "qi.txt")
	LogFile = pathL.Join(QiDir, "log.txt")
	ConfigFile = pathL.Join(QiDir, "config")

	qiConfigPath := os.ExpandEnv(pathL.Join(QiDir, "config"))
	lines, err = cmn.ReadLines(configPath)
	if err != nil {
		panic(fmt.Sprintf("error reading config, error: %v", err))
	}
	IdCounter, err = uint32(strconv.Atoi(lines[0]))
	if err != nil {
		panic(fmt.Sprintf("error reading id_counter, error: %v", err))
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

	content, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return content, true
}
