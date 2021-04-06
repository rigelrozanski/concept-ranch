package idea

import (
	"io/ioutil"
	"log"
	"path"
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
		ext := path.Ext(file.Name())
		if ext == ".swp" || ext == ".vim" {
			continue
		}
		ideas = append(ideas, NewIdeaFromFilename(file.Name(), false))
	}
	return ideas
}

func TagUsedInNonConsuming(tag string) bool {
	files, err := ioutil.ReadDir(IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fn := file.Name()
		if strings.HasPrefix(fn, "c") { // do not read from consumed ideas
			continue
		}
		if path.Ext(fn) == ".swp" {
			continue
		}
		if strings.Contains(fn, tag) {
			return true
		}
	}
	return false
}

// these ideas will be sorted from oldest to newest
func GetAllIdeas() (ideas Ideas) {
	files, err := ioutil.ReadDir(IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if path.Ext(file.Name()) == ".swp" {
			continue
		}
		ideas = append(ideas, NewIdeaFromFilename(file.Name(), false))
	}
	return ideas
}

// these ideas will be sorted from oldest to newest
func GetAllIdeasInRange(idStart, idEnd uint32) (ideas Ideas) {
	files, err := ioutil.ReadDir(IdeasDir)
	if err != nil {
		log.Fatal(err)
	}
	// TODO there has got to be a faster way of doing this!!!!
	for _, file := range files {
		id, skip := GetIdByFilename(file.Name())
		if skip {
			continue
		}
		if id >= idStart && id <= idEnd {
			if path.Ext(file.Name()) == ".swp" {
				continue
			}
			ideas = append(ideas, NewIdeaFromFilename(file.Name(), false))
		}
	}
	return ideas
}

func (ideas Ideas) WithTag(tag string) (subset Ideas) {
	return ideas.WithTags([]string{tag})
}

func (ideas Ideas) WithTags(tags []string) (subset Ideas) {
	for _, idea := range ideas {
		if idea.HasTags(tags) {
			subset = append(subset, idea)
		}
	}
	return subset
}

func (ideas Ideas) WithAnyOfTags(tags []string) (subset Ideas) {
	for _, idea := range ideas {
		if idea.HasAnyOfTags(tags) {
			subset = append(subset, idea)
		}
	}
	return subset
}

func (ideas Ideas) WithoutTag(tag string) (subset Ideas) {
	return ideas.WithoutTags([]string{tag})
}

func (ideas Ideas) WithoutTags(tags []string) (subset Ideas) {
	for _, idea := range ideas {
		if !idea.HasTags(tags) {
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

func (ideas Ideas) Paths() []string {
	fps := make([]string, len(ideas))
	for i, idea := range ideas {
		fps[i] = idea.Path()
	}
	return fps
}

func (ideas Ideas) Filenames() []string {
	fns := make([]string, len(ideas))
	for i, idea := range ideas {
		fns[i] = idea.Filename
	}
	return fns
}
