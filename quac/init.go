package quac

import (
	"fmt"
	"os"
	"path"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/thranch/quac/idea"
)

// directory name where boards are stored in repo
var QuDir, TrashCanDir, QuFile, LogFile, WorkingFnsFile, WorkingContentFile string

// load config and set global file directories
func Initialize(thranchConfigPath string) {

	lines, err := cmn.ReadLines(thranchConfigPath)
	if err != nil {
		panic(fmt.Sprintf("error reading %v, error: %v", thranchConfigPath, err))
	}

	QuDir = lines[0]
	idea.IdeasDir = path.Join(QuDir, "ideas")
	TrashCanDir = path.Join(QuDir, "trash")
	QuFile = path.Join(QuDir, "qu")
	LogFile = path.Join(QuDir, "log")
	idea.ConfigFile = path.Join(QuDir, "config")
	WorkingFnsFile = path.Join(QuDir, "working_files")
	WorkingContentFile = path.Join(QuDir, "working_content")
	idea.LastIdFile = path.Join(QuDir, "last")

	EnsureBasics()

	IdeasDir = idea.IdeasDir
	ConfigFile = idea.ConfigFile
	LastIdFile = idea.LastIdFile
}

func EnsureBasics() {
	if !cmn.FileExists(QuDir) {
		panic("directory specified in thranch config does not exist")
	}

	_ = os.Mkdir(idea.IdeasDir, os.ModePerm)
	_ = os.Mkdir(TrashCanDir, os.ModePerm)

	if !cmn.FileExists(QuFile) {
		err := cmn.CreateEmptyFile(QuFile)
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
