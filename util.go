package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	cmn "github.com/rigelrozanski/common"
	"github.com/rigelrozanski/qi/lib"
)

func Consume(consumedID, optionalEntry string) {
	consumed, err := strconv.Atoi(consumedID)
	if err != nil {
		log.Fatalf("bad id %v", consumedID)
	}
	consumerFilepath := lib.SetConsume(uint32(consumed), optionalEntry)
	if optionalEntry == "" {
		openText(consumerFilepath)
	}
}

func Consumes(consumedID, consumesID string) {
	consumed, err := strconv.Atoi(consumedID)
	if err != nil {
		log.Fatalf("bad id %v", consumedID)
	}
	consumes, err := strconv.Atoi(consumesID)
	if err != nil {
		log.Fatalf("bad id %v", consumesID)
	}
	lib.SetConsumes(uint32(consumed), uint32(consumes))
}

func Zombie(zombieID string) {
	zombie, err := strconv.Atoi(zombieID)
	if err != nil {
		log.Fatalf("bad id %v", zombieID)
	}
	lib.SetZombie(uint32(zombie))
}

func Lineage(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("bad id %v", idStr)
	}
	fmt.Print(lib.GetLineage(uint32(id)))
}

func QuickQuery(unsplitTagsOrID string) {
	id, err := strconv.Atoi(unsplitTagsOrID)
	if err == nil {
		ViewByID(uint32(id))
		return
	}
	splitTags := strings.Split(unsplitTagsOrID, ",")
	ViewByTags(splitTags)
}

func NewEmptyEntry(unsplitTags string) {
	splitTags := strings.Split(unsplitTags, ",")
	idear := lib.NewNonConsumingTextIdea(splitTags)
	writePath := path.Join(lib.IdeasDir, idear.Filename)
	lib.IncrementID()
	openText(writePath)
}

func QuickEntry(unsplitTags, entry string) {
	splitTags := strings.Split(unsplitTags, ",")
	Entry(entry, splitTags)
}

func MultiOpen(unsplitTagsOrID string) {
	id, err := strconv.Atoi(unsplitTagsOrID)
	if err == nil {
		filePath := lib.GetFilepathByID(uint32(id))
		open(filePath)
		return
	}
	splitTags := strings.Split(unsplitTagsOrID, ",")
	MultiOpenByTags(splitTags)
}

func parseIdStr(idStr string) uint32 {
	idI, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("error parsing id, error: %v", err)
	}
	return uint32(idI)
}

func RemoveByID(idStr string) {
	id := parseIdStr(idStr)
	lib.RemoveByID(id)
}

func CopyByID(idStr string) {
	id := parseIdStr(idStr)
	open(lib.CopyByID(id))
}

func ListTagsByID(idStr string) {
	id := parseIdStr(idStr)
	idea := lib.GetIdeaByID(id)
	fmt.Println(idea.Tags)
}

func KillTagByID(idStr, tagToKill string) {
	id := parseIdStr(idStr)
	idea := lib.GetIdeaByID(id)
	origFilename := idea.Filename
	(&idea).RemoveTag(tagToKill)
	(&idea).UpdateFilename()

	origPath := path.Join(lib.IdeasDir, origFilename)
	newPath := path.Join(lib.IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func AddTagByID(idStr, tagToAdd string) {
	id := parseIdStr(idStr)
	idea := lib.GetIdeaByID(id)
	origFilename := idea.Filename
	idea.Tags = append(idea.Tags, tagToAdd)
	(&idea).UpdateFilename()

	origPath := path.Join(lib.IdeasDir, origFilename)
	newPath := path.Join(lib.IdeasDir, idea.Filename)
	err := os.Rename(origPath, newPath)
	if err != nil {
		log.Fatal(err)
	}
}

func RenameTag(from, to string) {
	files, err := ioutil.ReadDir(lib.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		origFn := file.Name()
		if !strings.Contains(origFn, from) {
			continue
		}
		idea := lib.NewIdeaFromFilename(origFn)
		(&idea).RenameTag(from, to)
		(&idea).UpdateFilename()

		// perform the file rename
		origPath := path.Join(lib.IdeasDir, origFn)
		newPath := path.Join(lib.IdeasDir, idea.Filename)
		err := os.Rename(origPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DestroyTag(tag string) {
	files, err := ioutil.ReadDir(lib.IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		origFn := file.Name()
		if !strings.Contains(origFn, tag) {
			continue
		}
		idea := lib.NewIdeaFromFilename(origFn)
		(&idea).RemoveTag(tag)
		(&idea).UpdateFilename()

		// perform the file rename
		origPath := path.Join(lib.IdeasDir, origFn)
		newPath := path.Join(lib.IdeasDir, idea.Filename)
		err := os.Rename(origPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//__________________

func ListAllTags() {
	ideas := lib.GetAllIdeas(lib.IdeasDir)
	fmt.Println(ideas.UniqueTags())
}

func ListAllFiles() {
	ideas := lib.GetAllIdeas(lib.IdeasDir)
	for _, idea := range ideas {
		fmt.Println(idea.Filename)
	}
}

func ViewByID(id uint32) {
	content, found := lib.GetContentByID(id)
	if !found {
		fmt.Println("nothing found with that id")
	}
	fmt.Printf("%s\n", content)
}

func ViewByTags(tags []string) {
	content, found := lib.ConcatAllContentFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	fmt.Printf("%s\n", content)
}

func MultiOpenByTags(tags []string) {
	found, maxFNLen, singleReturn := lib.WriteWorkingContentAndFilenamesFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
		return
	}
	// if only a single entry is found then open only it!
	if singleReturn != "" {
		open(singleReturn)
		return
	}
	openTextSplit(lib.WorkingFnsFile, lib.WorkingContentFile, maxFNLen)
	lib.SaveFromWorkingFiles()
}

func RemoveById(id uint32) error {
	return nil
}

// create an entry
func Entry(entryOrPath string, tags []string) {

	if cmn.FileExists(entryOrPath) { // is a path

		idea := lib.NewIdeaFromFile(tags, entryOrPath)

		err := cmn.Copy(entryOrPath, idea.Path())
		if err != nil {
			log.Fatal(err)
		}
		lib.IncrementID()

		return
	}

	idea := lib.NewNonConsumingTextIdea(tags)
	err := cmn.WriteLines([]string{entryOrPath}, idea.Path())
	if err != nil {
		log.Fatalf("error writing new file: %v", err)
	}
	lib.IncrementID()
}

//_______________________________________________________________________________________________________

// open supported files
func open(pathToOpen string) {
	ext := path.Ext(pathToOpen)

	switch lib.GetKind(ext) {
	case lib.KindText:
		openText(pathToOpen)
	case lib.KindImage:
		viewImage(pathToOpen)
	case lib.KindAudio:
		listenAudio(pathToOpen)
	}
}

func viewImage(pathToOpen string) {

	fmt.Println(path.Base(pathToOpen))
	cmd := exec.Command("imgcat", pathToOpen)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func listenAudio(pathToOpen string) {

	fmt.Println(path.Base(pathToOpen))
	cmd := exec.Command("afplay", pathToOpen)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func openText(pathToOpen string) {

	// ignore error, allow for no file to be present
	origBz, _ := ioutil.ReadFile(pathToOpen)

	cmd := exec.Command("vim", "-c", "+normal 1G1|", pathToOpen) //start in the upper left corner nomatter
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	finalBz, err := ioutil.ReadFile(pathToOpen)
	if err != nil {
		log.Fatal(err)
	}
	if bytes.Compare(origBz, finalBz) != 0 {
		lib.UpdateEditedDateNow(pathToOpen)
	}
}

func openTextSplit(pathToOpenLeft, pathToOpenRight string, maxFNLen int) {

	cmd := exec.Command("vim",
		"-c", "vertical resize "+strconv.Itoa(maxFNLen+4)+" | execute \"normal \\<C-w>\\<C-l>\"",
		"-O", pathToOpenLeft, pathToOpenRight)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
