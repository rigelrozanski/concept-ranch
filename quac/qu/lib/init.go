package lib

import (
	"fmt"
	"os"
	"path"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/wranch/quac/qu/lib/idea"
)

// directory name where boards are stored in repo
var QuDir, QuFile, LogFile, WorkingFnsFile, WorkingContentFile string

// load config and set global file directories
func Initialize(wranchConfigPath string) {

	lines, err := cmn.ReadLines(wranchConfigPath)
	if err != nil {
		panic(fmt.Sprintf("error reading %v, error: %v", wranchConfigPath, err))
	}

	QuDir = lines[0]
	idea.IdeasDir = path.Join(QuDir, "ideas")
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
		panic("directory specified in wranch config does not exist")
	}
	_ = os.Mkdir(idea.IdeasDir, os.ModePerm)
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
