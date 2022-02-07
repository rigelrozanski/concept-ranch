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

// XXX delete this
//// these ideas will be sorted from oldest to newest
//func GetAllIdeasInRange(idStart, idEnd uint32) (ideas Ideas) {
//    files, err := ioutil.ReadDir(IdeasDir)
//    if err != nil {
//        log.Fatal(err)
//    }
//    // TODO there has got to be a faster way of doing this!!!!
//    for _, file := range files {
//        id, skip := GetIdByFilename(file.Name())
//        if skip {
//            continue
//        }
//        if id >= idStart && id <= idEnd {
//            if path.Ext(file.Name()) == ".swp" {
//                continue
//            }
//            ideas = append(ideas, NewIdeaFromFilename(file.Name(), false))
//        }
//    }
//    return ideas
//}

// inclusive range
func (ideas Ideas) InRange(idStart, idEnd uint32) (subset Ideas) {
	for _, idea := range ideas {
		if idStart <= idea.Id && idea.Id <= idEnd {
			subset = append(subset, idea)
		}
	}
	return subset
}

func (ideas Ideas) WithTag(tag Tag) (subset Ideas) {
	return ideas.WithTags([]Tag{tag})
}

func (ideas Ideas) WithTags(tags []Tag) (subset Ideas) {
	for _, idea := range ideas {
		if idea.HasTags(tags) {
			subset = append(subset, idea)
		}
	}
	return subset
}

func (ideas Ideas) WithAnyOfTags(tags []Tag) (subset Ideas) {
	for _, idea := range ideas {
		if idea.HasAnyOfTags(tags) {
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

func (ideas Ideas) UniqueTags() []Tag {
	tags := make(map[string]Tag)
	for _, idea := range ideas {
		for _, tag := range idea.Tags {
			tags[tag.String()] = tag
		}
	}
	var out []Tag
	for _, tag := range tags {
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
