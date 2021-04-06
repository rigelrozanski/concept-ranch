package quac

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/rigelrozanski/thranch/quac/idea"
)

// Display the immediate lineage of ideas
func GetLineage(id uint32) (compiled string) {
	lineageIdea := GetIdeaByID(id, false)
	for _, consume := range lineageIdea.ConsumesIds {
		fn := GetFilenameByID(consume)
		content, found := GetContentByID(consume)
		if !found {
			log.Fatalf("child not found: %v", consume)
		}
		compiled = fmt.Sprintf("%v\n%v\n%s", compiled, fn, content)
	}
	return compiled
}

// copy an idea by the id
func SetConsume(consumedId uint32, entry string) (consumerFilepath string) {
	consumedIdea := GetIdeaByID(consumedId, true)

	// consumer: remove the id, add in a new id, add the consumes id
	consumerIdea := idea.NewConsumingTextIdea(consumedIdea)
	idea.IncrementID()
	WriteIdea(consumerIdea.Filename, entry)

	consumedIdea.SetConsumed()
	return consumerIdea.Path()
}

func SetConsumes(consumedId, consumesId uint32) {

	consumedIdea := GetIdeaByID(consumedId, true)
	consumesIdea := GetIdeaByID(consumesId, true)

	// consumer: remove the id, add in a new id, add the consumes id
	consumesIdea.ConsumesIds = append(consumesIdea.ConsumesIds, consumedId)
	srcPath := consumesIdea.Path()
	(&consumesIdea).UpdateFilename()
	writePath := path.Join(idea.IdeasDir, consumesIdea.Filename)
	err := os.Rename(srcPath, writePath)
	if err != nil {
		log.Fatal(err)
	}

	consumedIdea.SetConsumed()
}

// Set a consumed idea to zombie
func SetZombie(zombieId uint32) {
	consumedIdea := GetIdeaByID(zombieId, true)
	consumedIdea.SetZombie()
}
