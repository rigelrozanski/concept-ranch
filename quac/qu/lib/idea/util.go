package idea

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	layout = "2006-01-02" // time parse layout for YYYYMMDD
	Last   = "last"       // last id keyword

	CycleAlive = iota
	CycleConsumed
	CycleZombie

	KindText = iota
	KindImage
	KindAudio
	KindEnText
)

var (
	IdeasDir, ConfigFile, LastIdFile string
	zeroDate                         time.Time
	rxConsumedId                     = regexp.MustCompile(`[c]\d{6,6}`)
)

func TodayDate() time.Time {
	todayDate, err := time.Parse(layout, time.Now().Format(layout))
	if err != nil {
		log.Fatal(err)
	}
	return todayDate
}

func GetKind(ext string) int {
	switch ext {
	case ".mp3", ".wav":
		return KindAudio
	case ".jpg", ".jpeg", ".tiff", ".png":
		return KindImage
	case "", ".txt":
		return KindText
	case ".en":
		return KindEnText
	default:
		log.Fatalf("unknown filetype: %v", ext)
	}
	return 0
}

func (idea Idea) Path() string {
	return path.Join(IdeasDir, idea.Filename)
}

func (idea Idea) GetContent() []byte {
	content, err := ioutil.ReadFile(idea.Path())
	if err != nil {
		log.Fatal(err)
	}
	return content
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

func (idea Idea) IsText() bool {
	return (idea.Kind == KindText)
}

func (idea Idea) IsImage() bool {
	return (idea.Kind == KindImage)
}

func (idea Idea) IsAudio() bool {
	return (idea.Kind == KindAudio)
}

// creates the filename based on idea information
func (idea *Idea) UpdateFilename() {

	prefix := idea.Prefix()

	strList := []string{
		prefix,
		IdStr(idea.Id),
		idea.Created.Format(layout),
		"e" + idea.Edited.Format(layout)}
	if idea.Cycle != CycleAlive {
		strList = append(strList, "c"+idea.Consumed.Format(layout))
	}
	strList = append(strList, itoa(idea.ConsumesIds)...)
	strList = append(strList, idea.Tags...)

	joined := strings.Join(strList, ",")
	joined += idea.Ext
	idea.Filename = joined
}

func itoa(in []uint32) []string {
	out := make([]string, len(in))
	for i, el := range in {
		out[i] = "c" + IdStr(el)
	}
	return out[:]
}

func IdStr(id uint32) string {
	return fmt.Sprintf("%06d", id)
}
