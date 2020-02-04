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

func ViewImage(pathToOpen string) {
	fmt.Println(path.Base(pathToOpen))
	ViewImageNoFilename(pathToOpen)
}

func ViewImageNoFilename(pathToOpen string) {
	cmd := exec.Command("imgcat", pathToOpen)
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

	pathToOpen := GetFilepathByID(id)
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
