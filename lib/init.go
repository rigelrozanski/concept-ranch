package lib

import (
	"fmt"
	"os"
	"path"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/qi/lib/idea"
)

// directory name where boards are stored in repo
var QiDir, QiFile, LogFile, WorkingFnsFile, WorkingContentFile string

// load config and set global file directories
func init() {

	rootConfigPath := os.ExpandEnv("$HOME/.qi_config.txt")
	lines, err := cmn.ReadLines(rootConfigPath)
	if err != nil {
		panic(fmt.Sprintf("error reading ~/.qi_config.txt, error: %v", err))
	}

	QiDir = lines[0]
	idea.IdeasDir = path.Join(QiDir, "ideas")
	QiFile = path.Join(QiDir, "qi")
	LogFile = path.Join(QiDir, "log")
	idea.ConfigFile = path.Join(QiDir, "config")
	WorkingFnsFile = path.Join(QiDir, "working_files")
	WorkingContentFile = path.Join(QiDir, "working_content")
	idea.LastIdFile = path.Join(QiDir, "last")

	EnsureBasics()

	IdeasDir = idea.IdeasDir
	ConfigFile = idea.ConfigFile
	LastIdFile = idea.LastIdFile
}

func EnsureBasics() {
	if !cmn.FileExists(QiDir) {
		panic("directory specified in ~/.qi_config.txt does not exist")
	}
	_ = os.Mkdir(idea.IdeasDir, os.ModePerm)
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
	if !cmn.FileExists(idea.ConfigFile) {
		err := cmn.WriteLines([]string{"000001"}, idea.ConfigFile)
		if err != nil {
			panic(err)
		}
	}
	if !cmn.FileExists(idea.LastIdFile) {
		err := cmn.WriteLines([]string{"000000"}, idea.LastIdFile)
		if err != nil {
			panic(err)
		}
	}
}
