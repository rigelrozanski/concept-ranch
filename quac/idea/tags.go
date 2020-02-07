package idea

import (
	"io/ioutil"
	"log"
	"path"
	"strings"
)

//_______________________________________________________

func ConcatAllContentFromTags(tags []string) (content []byte, found bool) {
	ideas := GetAllIdeasNonConsuming()
	subset := ideas.WithTags(tags).WithText()

	if len(subset) == 0 {
		return content, false
	}
	for _, idea := range subset {
		ideaContent, err := ioutil.ReadFile(path.Join(IdeasDir, idea.Filename))
		if err != nil {
			log.Fatal(err)
		}
		content = append(content, ideaContent...)
	}
	return content, true
}

//_______________________________________________________

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

// returns true if the idea has the tagged value
func (idea Idea) GetTaggedValue(tvName string) (val string, found bool) {
	for _, tag := range idea.Tags {
		ideaTvName := strings.Split(tag, ":")
		if len(ideaTvName) == 2 && ideaTvName[0] == tvName {
			return ideaTvName[1], true
		}
	}
	return "", false
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