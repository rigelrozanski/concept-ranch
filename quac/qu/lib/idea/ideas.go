package idea

import (
	"io/ioutil"
	"log"
	"strings"
)

type Ideas []Idea

func GetAllIdeasNonConsuming() (ideas Ideas) {
	files, err := ioutil.ReadDir(IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "c") { // do not read from consumed ideas
			continue
		}
		ideas = append(ideas, NewIdeaFromFilename(file.Name()))
	}
	return ideas
}

func GetAllIdeas() (ideas Ideas) {
	files, err := ioutil.ReadDir(IdeasDir)
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

func (ideas Ideas) WithText() (subset Ideas) {
	for _, idea := range ideas {
		if idea.IsText() {
			subset = append(subset, idea)
		}
	}
	return subset
}

func (ideas Ideas) WithImage() (subset Ideas) {
	for _, idea := range ideas {
		if idea.IsImage() {
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