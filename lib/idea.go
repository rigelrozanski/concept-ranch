package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	layout = "2006-01-02" // time parse layout for YYYYMMDD

	kindAudio  = "audio"
	kindVisual = "visual"
	kindText   = "text"
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
	IsConsumed  bool
	Id          uint32
	ConsumesIds []uint32 // Id of idea which this idea consumes
	Kind        string   // kind of information
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
		IsConsumed:  false,
		Id:          GetNextID(),
		ConsumesIds: consumesIds,
		Kind:        kindText,
		Created:     todayDate,
		Edited:      todayDate,
		Consumed:    zeroDate,
		Tags:        tags,
	}

	(&idea).UpdateFilename()

	return idea
}

func NewIdeaFromFilename(filename string) (idea Idea) {
	idea.Filename = filename

	ext := path.Ext(filename)
	switch ext {
	case ".mp3", ".wav":
		idea.Kind = kindAudio
	case ".jpg", ".jpeg", ".tiff", ".png":
		idea.Kind = kindVisual
	case "", ".txt":
		idea.Kind = kindText
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
		idea.IsConsumed = false
	} else {
		idea.IsConsumed = true
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

// creates the filename based on idea information
func (idea *Idea) UpdateFilename() {

	prefix := "a"
	if idea.IsConsumed {
		prefix = "c"
	}

	strList := []string{
		prefix,
		strconv.Itoa(int(idea.Id)),
		idea.Created.Format(layout),
		"e" + idea.Edited.Format(layout)}
	if idea.IsConsumed {
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
	fmt.Printf("debug tagToKill: %v\n", tagToKill)
	fmt.Printf("debug idea: %v\n", idea)
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

func itoa(in []uint32) []string {
	out := make([]string, len(in))
	for i, el := range in {
		out[i] = strconv.Itoa(int(el))
	}
	return out[:]
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
