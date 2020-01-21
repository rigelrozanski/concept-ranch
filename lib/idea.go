package lib

import (
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	layout    = "20060102" // time parse layout for YYYYMMDD
	alive     = "ALIVE"
	kindImage = "image"
	kindText  = "text"
)

type Ideas []Idea

func PathToIdeas(dir string) (ideas Ideas) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		ideas = append(ideas, FilenameToIdea(file.Name()))
	}
	return ideas
}

func (is Ideas) WithTags(tags ...string) (subset Ideas) {
	for _, idea := range is {
		if idea.HasTags(tags) {
			subset = append(subset, idea)
		}
	}
	return subset
}

func (is Ideas) UniqueTags() []string {
	tags := make(map[string]string)
	for _, idea := range is {
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

type Idea struct {
	Filename   string
	Id         uint64
	ConsumesId uint64 // Id of idea which this idea consumes
	Kind       string // kind of information
	IsConsumed bool
	Created    time.Time
	Edited     time.Time
	Consumed   time.Time
	Tags       []string
}

// creates the filename based on idea information
func (i *Idea) CreateFilename() {

	ConsumedDate := alive
	if !IsConsumed {
		ConsumedDate = Consumed.Format(layout)
	}

	strList := []string{
		strconv.Itoa(int(i.Id)),
		"c" + strconv.Itoa(int(i.ConsumesId)),
		Created.Format(layout),
		Edited.Format(layout),
		ConsumedDate}
	strList = append(strList, tags...)

	i.Filename = strings.Join(strList, ",")
}

func NewNonConsumingIdea(tags []string) Idea {
	return NewAliveIdea(0, tags)
}

// NewIdea creates a new Idea object
func NewAliveIdea(consumesId uint64, tags []string) Idea {

	idea := Idea{
		Id:         GetNextID(),
		ConsumesId: consumesId,
		Kind:       kindText,
		IsConsumed: false,
		Created:    created,
		Edited:     created,
		Consumed:   alive,
		Tags:       tags,
	}

	(&idea).CreateFilename()

	return idea
}

// returns true if the idea contains all the input tags
func (i Idea) HasTags(tags ...string) bool {
	for _, tag := range tags {
		hasTag := false
		for _, ideaTag := range Idea.Tags {
			if tag == ideaTag {
				continue
			}
		}
		if !hasTag {
			return false
		}
	}
	return true
}

func FilenameToIdea(filename string) (idea Idea) {
	idea.Filename = filename

	ext := path.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".tiff", ".png":
		idea.Kind = kindImage
	case "":
		idea.Kind = kindText
	}

	base := strings.TrimSuffix(filename, path.Ext(filename))
	split := strings.Split(base, ",")
	if len(split) < 6 { // must have at minimum: Id, ConsumesId, Created, Edited, Consumed and a Tag
		log.Fatalf("bad filename at %v", filename)
	}

	// Get id
	idea.Id = uint64(strconv.Atoi(split[0]))

	// Get consumes id
	if !strings.HasPrefix(split[1], "c") {
		log.Fatalf("bad consumes-id file format at %v", filename)
	}
	idea.ConsumesId = uint64(strconv.Atoi(strings.TrimPrefix(split[1], "c")))

	// Get Creation Date
	created, err := time.Parse(layout, split[2])
	if err != nil {
		log.Fatalf("bad created date file format at %v: %v", filename, err)
	}
	idea.Created = created

	// Get Edit Date
	if !strings.HasPrefix(split[3], "e") {
		log.Fatalf("bad edit date file format at %v", filename)
	}
	created, err := time.Parse(layout, strings.TrimPrefix(split[3], "e"))
	if err != nil {
		log.Fatalf("bad created date file format at %v: %v", filename, err)
	}
	idea.Edited = created

	// Get Consumed Date
	if split[4] == alive {
		idea.IsConsumed = false
	} else {
		idea.IsConsumed = true
		if !strings.HasPrefix(split[4], "c") {
			log.Fatalf("bad consumed date file format at %v", filename)
		}
		created, err := time.Parse(layout, strings.TrimPrefix(split[4], "c"))
		if err != nil {
			log.Fatalf("bad consumed date file format at %v: %v", filename, err)
		}
		idea.Consumed = created
	}

	// Get Tags
	for i := 5; i < len(split); i++ {
		idea.Tags = append(idea.Tags, split[i])
	}
}
