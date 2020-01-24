package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	layout = "2006-01-02" // time parse layout for YYYYMMDD

	CycleAlive = iota
	CycleConsumed
	CycleZombie

	KindText = iota
	KindImage
	KindAudio
)

var (
	zeroDate     time.Time
	rxConsumedId = regexp.MustCompile(`[c]\d{6,6}`)
)

//_________________________________________________________________

type Ideas []Idea

func PathToIdeas(dir string) (ideas Ideas) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		ideas = append(ideas, NewIdeaFromFilename(file.Name()))
	}
	return ideas
}

func (ideas Ideas) WithTags(tags []string) (subset Ideas) {
	for _, idea := range ideas {
		if idea.HasTags(tags) {
			subset = append(subset, idea)
		}
	}
	return subset
}

func (ideas Ideas) UniqueTags() []string {
	tags := make(map[string]string)
	for _, idea := range ideas {
		for _, tag := range idea.Tags {
			tags[tag] = ""
		}
	}
	var out []string
	for tag, _ := range tags {
		out = append(out, tag)
	}
	return out
}

//_________________________________________________________________

type Idea struct {
	Filename    string
	Cycle       int // alive/consumed/zombie
	Id          uint32
	ConsumesIds []uint32 // Id of idea which this idea consumes
	Kind        int      // kind of information
	Created     time.Time
	Edited      time.Time
	Consumed    time.Time
	Tags        []string
}

func NewNonConsumingTextIdea(tags []string) Idea {
	return NewTextIdea([]uint32{}, tags)
}

func TodayDate() time.Time {
	todayDate, err := time.Parse(layout, time.Now().Format(layout))
	if err != nil {
		log.Fatal(err)
	}
	return todayDate
}

// NewAliveIdea creates a new idea object
func NewTextIdea(consumesIds []uint32, tags []string) Idea {

	todayDate := TodayDate()

	idea := Idea{
		Cycle:       CycleAlive,
		Id:          GetNextID(),
		ConsumesIds: consumesIds,
		Kind:        KindText,
		Created:     todayDate,
		Edited:      todayDate,
		Consumed:    zeroDate,
		Tags:        tags,
	}

	(&idea).UpdateFilename()
	return idea
}

// NewAliveIdea creates a new idea object
func NewConsumingIdea(consumesIdea Idea) Idea {

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
	switch ext {
	case ".mp3", ".wav":
		idea.Kind = KindAudio
	case ".jpg", ".jpeg", ".tiff", ".png":
		idea.Kind = KindImage
	case "", ".txt":
		idea.Kind = KindText
	default:
		log.Fatalf("unknown filetype: %v", ext)
	}

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
	id, err := strconv.Atoi(split[1])
	if err != nil {
		log.Fatal(err)
	}
	idea.Id = uint32(id)

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
		id, err = strconv.Atoi(strings.TrimPrefix(split[ri], "c"))
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

func (idea Idea) Path() string {
	return path.Join(IdeasDir, idea.Filename)
}

func (idea Idea) Prefix() (prefix string) {
	switch idea.Cycle {
	case CycleAlive:
		prefix = "a"
	case CycleConsumed:
		prefix = "c"
	case CycleZombie:
		prefix = "z"
	}
	return prefix
}

// creates the filename based on idea information
func (idea *Idea) UpdateFilename() {

	prefix := idea.Prefix()

	strList := []string{
		prefix,
		idStr(idea.Id),
		idea.Created.Format(layout),
		"e" + idea.Edited.Format(layout)}
	if idea.Cycle == CycleConsumed {
		strList = append(strList, "c"+idea.Consumed.Format(layout))
	}
	strList = append(strList, itoa(idea.ConsumesIds)...)
	strList = append(strList, idea.Tags...)

	idea.Filename = strings.Join(strList, ",")
}

// rename the tag on this idea
func (idea *Idea) RenameTag(from, to string) {
	for i, tag := range idea.Tags {
		if tag == from {
			idea.Tags[i] = to
			break
		}
	}
}

// rename the tag on this idea
func (idea *Idea) RemoveTag(tagToKill string) {
	if len(idea.Tags) == 1 && idea.Tags[0] == tagToKill {
		log.Fatalf("cannot remove the final tag of %v, aborting", idea.Filename)
	}
	for i, tag := range idea.Tags {
		if tag == tagToKill {
			idea.Tags = append(idea.Tags[:i], idea.Tags[i+1:]...)
			break
		}
	}
}

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

func itoa(in []uint32) []string {
	out := make([]string, len(in))
	for i, el := range in {
		out[i] = "c" + idStr(el)
	}
	return out[:]
}

func idStr(id uint32) string {
	return fmt.Sprintf("%06d", id)
}

// returns true if the idea contains all the input tags
func (idea Idea) HasTags(tags []string) bool {
	for _, tag := range tags {
		hasTag := false
		for _, ideaTag := range idea.Tags {
			if tag == ideaTag {
				hasTag = true
				continue
			}
		}
		if !hasTag {
			return false
		}
	}
	return true
}
