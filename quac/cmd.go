package quac

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
)

// open supported files
func Open(pathToOpen string) {
	ext := path.Ext(pathToOpen)

	id := GetIdByFilename(pathToOpen)
	PrependLast(id)

	switch GetKind(ext) {
	case KindText:
		OpenText(pathToOpen)
	case KindEnText:
		OpenText(pathToOpen)
	case KindImage:
		ViewImage(pathToOpen)
	case KindAudio:
		ListenAudio(pathToOpen)
	}
}

// open supported files
func View(pathToOpen string) {
	ext := path.Ext(pathToOpen)

	id := GetIdByFilename(pathToOpen)
	PrependLast(id)

	switch GetKind(ext) {
	case KindText:
		ViewText(pathToOpen)
	case KindImage:
		ViewImage(pathToOpen)
	case KindAudio:
		ListenAudio(pathToOpen)
	}
}

func ViewImage(pathToOpen string) {
	fmt.Println(path.Base(pathToOpen))
	ViewImageNoFilename(pathToOpen)
}

func ViewImageNoFilename(pathToOpen string) {
	//fmt.Printf("debug pathToOpen: %v\n", pathToOpen)
	cmd := exec.Command("kitty", "+kitten", "icat", pathToOpen) // using kitty command line
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func ListenAudio(pathToOpen string) {

	fmt.Println(path.Base(pathToOpen))
	cmd := exec.Command("afplay", pathToOpen)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func ViewText(filepath string) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", content)
}

func OpenText(pathToOpen string) {

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
		UpdateEditedDateNow(pathToOpen)
	}
}

func SetEncryptionById(id uint32) {

	pathToOpen, found := GetFilepathByID(id)
	if !found {
		fmt.Println("nothing found at that ID")
		os.Exit(1)
	}
	enPath := UpdateFilepathToEncrypted(pathToOpen)

	// ignore error, allow for no file to be present
	cmd := exec.Command("vim", "-c", "X", enPath) //start in the upper left corner nomatter
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func OpenTextSplit(pathToOpenLeft, pathToOpenRight string, maxFNLen int) {

	// limit the split
	if maxFNLen > 65 {
		maxFNLen = 65
	}

	cmd := exec.Command("vim",
		"-c", "vertical resize "+strconv.Itoa(maxFNLen+4)+
			" | set scb!"+ // set the scrollbind
			" | execute \"normal \\<C-w>\\<C-l>\""+
			" | set scb!", "-O", pathToOpenLeft, pathToOpenRight)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
