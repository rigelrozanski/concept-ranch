package idea

import (
	"log"
	"path"
	"strconv"
	"strings"
	"time"
)

type Idea struct {
	Filename    string
	Cycle       int // alive/consumed/zombie
	Id          uint32
	ConsumesIds []uint32 // Id of idea which this idea consumes
	Kind        int      // kind of information
	Ext         string   // file extension, blank is assumed to be text
	Created     time.Time
	Edited      time.Time
	Consumed    time.Time
	Tags        []string
}

func NewNonConsumingTextIdea(tags []string) Idea {
	return NewTextIdea([]uint32{}, tags)
}

// NewAliveIdea creates a new idea object
func NewTextIdea(consumesIds []uint32, tags []string) Idea {

	todayDate := TodayDate()

	idea := Idea{
		Cycle:       CycleAlive,
		Id:          GetNextID(),
		ConsumesIds: consumesIds,
		Kind:        KindText,
		Ext:         "",
		Created:     todayDate,
		Edited:      todayDate,
		Consumed:    zeroDate,
		Tags:        tags,
	}

	(&idea).UpdateFilename()
	return idea
}

func NewIdeaFromFile(tags []string, filepath string, tagFilename bool) Idea {
	todayDate := TodayDate()

	ext := path.Ext(filepath)
	if tagFilename {
		pathbase := strings.TrimSuffix(path.Base(filepath), ext)
		tags = append(tags, pathbase)
	}

	idea := Idea{
		Cycle:       CycleAlive,
		Id:          GetNextID(),
		ConsumesIds: []uint32{},
		Kind:        GetKind(ext),
		Ext:         ext,
		Created:     todayDate,
		Edited:      todayDate,
		Consumed:    zeroDate,
		Tags:        tags,
	}

	(&idea).UpdateFilename()
	return idea
}

// NewAliveIdea creates a new idea object
func NewConsumingTextIdea(consumesIdea Idea) Idea {

	todayDate := TodayDate()

	consumesIdCp := make([]uint32, len(consumesIdea.ConsumesIds))
	consumesTagCp := make([]string, len(consumesIdea.Tags))
	copy(consumesIdCp, consumesIdea.ConsumesIds)
	copy(consumesTagCp, consumesIdea.Tags)

	idea := Idea{
		Cycle:       CycleAlive,
		Id:          GetNextID(),
		ConsumesIds: append(consumesIdCp, consumesIdea.Id),
		Kind:        KindText,
		Created:     todayDate,
		Edited:      todayDate,
		Consumed:    zeroDate,
		Tags:        consumesTagCp,
	}

	(&idea).UpdateFilename()
	return idea
}

func NewIdeaFromFilename(filename string) (idea Idea) {
	idea.Filename = filename

	ext := path.Ext(filename)
	idea.Ext = ext
	idea.Kind = GetKind(ext)

	base := strings.TrimSuffix(filename, path.Ext(filename))
	split := strings.Split(base, ",")
	if len(split) < 5 { // must have at minimum: ConsumedPrefix, Id, Created, Edited,and a Tag
		log.Fatalf("bad filename at %v", filename)
	}

	// get consumption prefix
	if split[0] == "a" {
		idea.Cycle = CycleAlive
	} else {
		idea.Cycle = CycleConsumed
	}

	// Get id
	id, err := ParseID(split[1])
	if err != nil {
		log.Fatal(err)
	}
	idea.Id = id

	// get creation date
	created, err := time.Parse(layout, split[2])
	if err != nil {
		log.Fatalf("bad created date file format at %v: %v", filename, err)
	}
	idea.Created = created

	// get edit date
	if !strings.HasPrefix(split[3], "e") {
		log.Fatalf("bad edit date file format at %v", filename)
	}
	edited, err := time.Parse(layout, strings.TrimPrefix(split[3], "e"))
	if err != nil {
		log.Fatalf("bad created date file format at %v: %v", filename, err)
	}
	idea.Edited = edited

	// rolling index
	ri := 4

	// get any consumed date
	if strings.HasPrefix(split[ri], "c") {
		// ignore error
		consumed, err := time.Parse(layout, strings.TrimPrefix(split[4], "c"))
		if err == nil {
			idea.Consumed = consumed
			ri++
		}
	}

	// get any consumes id(s)
	for ; ; ri++ {
		if !rxConsumedId.MatchString(split[ri]) {
			break
		}
		// special case to not log so don't use ParseID
		id, err := strconv.Atoi(strings.TrimPrefix(split[ri], "c"))
		if err != nil {
			log.Fatal(err)
		}
		idea.ConsumesIds = append(idea.ConsumesIds, uint32(id))
	}

	// get tag(s)
	if ri == len(split) {
		log.Fatalf("no tags on file: %v", filename)
	}
	for ; ri < len(split); ri++ {
		idea.Tags = append(idea.Tags, split[ri])
	}

	return idea
}
