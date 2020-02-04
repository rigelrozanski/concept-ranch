package idea

import (
	"log"
	"os"
	"path"
)

// rename the tag on this idea
func (idea *Idea) SetConsumed() {
	origFilename := idea.Filename
	idea.Cycle = CycleConsumed
	idea.Consumed = TodayDate()
	idea.UpdateFilename()

	srcPath := path.Join(IdeasDir, origFilename)
	writePath := path.Join(IdeasDir, idea.Filename)
	err := os.Rename(srcPath, writePath)
	if err != nil {
		log.Fatal(err)
	}
}

// rename the tag on this idea
func (idea *Idea) SetZombie() {
	origFilename := idea.Filename
	idea.Cycle = CycleZombie
	idea.UpdateFilename()

	srcPath := path.Join(IdeasDir, origFilename)
	writePath := path.Join(IdeasDir, idea.Filename)
	err := os.Rename(srcPath, writePath)
	if err != nil {
		log.Fatal(err)
	}
}
