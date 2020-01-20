package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	pathL "path"
	"path/filepath"
	"strconv"
	"strings"

	cmn "github.com/rigelrozanski/common"

	dbm "github.com/tendermint/tm-db"
)

var db dbm.DB

func init() {
	var err error
	db, err = NewLevelDB("wt", "./wtdb/")
	if err != nil {
		panic(err)
	}
}

// NewLevelDB instantiate a new LevelDB instance according to DBBackend.
func NewLevelDB(name, dir string) (db dbm.DB, err error) {
	backend := dbm.GoLevelDBBackend
	if DBBackend == string(dbm.CLevelDBBackend) {
		backend = dbm.CLevelDBBackend
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("couldn't create db: %v", r)
		}
	}()
	return dbm.NewDB(name, backend, dir), err
}

// directory name where boards are stored in repo
var (
	BoardsDir = "boards"
	TrashDir  = "trash"
)

// get the contents of a local wb
func GetWB(name string) (content []string, found bool) {

	path, err := GetWbPath(name)
	if err != nil {
		return content, false
	}
	if !cmn.FileExists(path) {
		return content, false
	}

	content, err = cmn.ReadLines(path)
	if err != nil {
		return content, false
	}

	return content, true
}

// get the contents of a local wb in bytes
func GetWBRaw(name string) (content []byte, found bool) {

	path, err := GetWbPath(name)
	if err != nil {
		return content, false
	}
	if !cmn.FileExists(path) {
		return content, false
	}

	content, err = ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return content, true
}

func todoStr(name string) string {
	return fmt.Sprintf("TODO add: [%v]", name)
}

// get the contents of a local wb
func RemoveFromLS(lsName, remove string) error {

	path, err := GetWbPath(lsName)
	if err != nil {
		return err
	}
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	fileStr := string(bz)

	// replace any todo lines
	fileStr = strings.Replace(fileStr, todoStr(remove)+"\n", "", 1)

	// alternatively replace w whitespace
	whitespace := fmt.Sprintf(`%`+strconv.Itoa(len(remove)+2)+`v`, " ")
	fileStr = strings.Replace(fileStr, fmt.Sprintf("[%s]", remove), whitespace, 1)

	err = ioutil.WriteFile(path, []byte(fileStr), 0644)
	if err != nil {
		return err
	}

	return nil
}

// get the contents of a local wb
func AddToLS(lsName, newWB string) error {
	return PrependWB(lsName, todoStr(newWB))
}

// replace a string in ls
func ReplaceInLS(lsName, oldName, newName string) error {
	path, err := GetWbPath(lsName)
	if err != nil {
		return err
	}
	oldName = fmt.Sprintf("[%v]", oldName)
	newName = fmt.Sprintf("[%v]", newName)
	return cmn.ReplaceAllStringInFile(path, oldName, newName)
}

// nolint
func MustPrependWB(name, entry string) {
	err := PrependWB(name, entry)
	if err != nil {
		panic(err)
	}
}

// prepend a string to a new top line within a wb
func PrependWB(name, entry string) error {

	path, err := GetWbPath(name)
	if err != nil {
		return err
	}
	if !cmn.FileExists(path) {
		return fmt.Errorf("file as %v does not exist", path)
	}

	content, err := cmn.ReadLines(path)
	if err != nil {
		return err
	}

	content = append([]string{entry}, content...)

	err = cmn.WriteLines(content, path)
	if err != nil {
		return err
	}

	return nil
}

// append a string to a new top line within a wb
func AppendWB(name, entry string) error {
	path, err := GetWbPath(name)
	if err != nil {
		return err
	}
	if !cmn.FileExists(path) {
		return fmt.Errorf("file as %v does not exist", path)
	}

	content, err := cmn.ReadLines(path)
	if err != nil {
		return err
	}

	content = append(content, entry)

	err = cmn.WriteLines(content, path)
	if err != nil {
		return err
	}
	return nil
}

// nolint
func MustClearWB(name string) {
	err := ClearWB(name)
	if err != nil {
		panic(err)
	}
}

// remove all the content of wb
func ClearWB(name string) error {
	path, err := GetWbPath(name)
	if err != nil {
		return err
	}
	if !cmn.FileExists(path) {
		return err
	}
	err = cmn.WriteLines([]string{}, path)
	if err != nil {
		return err
	}
	return nil
}

// get the contents of a local wb
func WbExists(name string) (found bool) {

	path, err := GetWbPath(name)
	if err != nil {
		return false
	}
	if !cmn.FileExists(path) {
		return false
	}
	return true
}

// get the full path of a wb
func GetWbPath(name string) (string, error) {
	return GetRepoPath(name, BoardsDir)
}

// get the full path of a wb in the trash
func GetWbInTrashPath(name string) (string, error) {
	return GetRepoPath(name, TrashDir)
}

// get the full path of a wb in the trash
func GetTrashPath() (string, error) {
	return GetRepoPath("", TrashDir)
}

// get the full path of a wb in a directory
func GetRepoPath(name, dir string) (string, error) {
	wbBackupRepoPath, err := GetWbBackupRepoPath()
	if err != nil {
		return "", err
	}
	return pathL.Join(wbBackupRepoPath, dir, name), nil
}

// move a wb from the boards dir to the trash DIR
func MoveWbToTrash(name string) error {

	wbBoardsPath, err := GetWbPath(name)
	if err != nil {
		return err
	}
	if !cmn.FileExists(wbBoardsPath) { //does the whiteboard not exist
		return errors.New("error can't delete non-existent whiteboard")
	}

	// make the trash path if it doesn't exist
	trashPath, err := GetTrashPath()
	if err != nil {
		return err
	}
	os.MkdirAll(trashPath, os.ModePerm)

	wbTrashPath, err := GetWbInTrashPath(name)
	if err != nil {
		return err
	}
	err = os.Rename(wbBoardsPath, wbTrashPath)
	if err != nil {
		return err
	}
	return nil
}

// delete all wbs in the trash folder
func EmptyTrash() error {

	trashPath, err := GetTrashPath()
	if err != nil {
		return err
	}

	d, err := os.Open(trashPath)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(trashPath, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// move a wb from the boards dir to the trash DIR
func RecoverWbFromTrash(name string) error {

	wbBoardsPath, err := GetWbPath(name)
	if err != nil {
		return err
	}
	wbTrashPath, err := GetWbInTrashPath(name)
	if err != nil {
		return err
	}
	if !cmn.FileExists(wbTrashPath) { //does the whiteboard not exist
		return errors.New("error can't recover - non-existent in trash")
	}
	err = os.Rename(wbTrashPath, wbBoardsPath)
	if err != nil {
		return err
	}
	return nil
}

// get the full path of a wb backup repo
func GetWbBackupRepoPath() (string, error) {
	configPath, err := cmn.GetRelPath("/src/github.com/rigelrozanski/wb", "config.txt")
	if err != nil {
		return "", fmt.Errorf("missing config.txt file in root of repo, error: %v", err)
	}
	lines, err := cmn.ReadLines(configPath)
	if err != nil {
		return "", fmt.Errorf("error reading config, error: %v", err)
	}
	return lines[0], nil
}

// get the full path of a wb backup repo
func GetBoardsGITDir() (string, error) {
	wbBackupRepoPath, err := GetWbBackupRepoPath()
	if err != nil {
		return "", err
	}
	return pathL.Join(wbBackupRepoPath, ".git"), nil
}

// function for iterating over the wbs
type IterateFn func(name, relPath string) (stop bool)

// perform the provided iterateFn for all wbs
func IterateWBs(iterFn IterateFn) error {

	boardPath, err := GetWbPath("")
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(boardPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		relPath, err := GetWbPath(name)
		if err != nil {
			return err
		}
		if iterFn(name, relPath) {
			break
		}
	}
	return nil
}
