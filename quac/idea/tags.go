package idea

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type Tag interface {
	GetName() string
	GetValue() string
	String() string
	Has(Tag) bool
}

// -----------------------------

// a tag can be either just a name
// or a name and a value.
type TagBase struct {
	Name  string
	Value string
}

// NewTagBase creates a new TagBase object
func NewTagBase(name, value string) TagBase {
	return TagBase{
		Name:  name,
		Value: value,
	}
}

func (t TagBase) GetName() string {
	return t.Name
}

func (t TagBase) GetValue() string {
	return t.Value
}

func (t TagBase) String() string {
	if t.Value == "" {
		return t.Name
	}
	return t.Name + "=" + t.Value
}

// ------------------------------------------
type TagReg struct{ TagBase }

var _ Tag = TagReg{}

var ReservedTagNames = []string{
	"CONTAINS",
	"WITHOUT",
}

// NewTagWithValue creates a new Tag with a value
func NewTagReg(name string) (Tag, error) {
	for _, rn := range ReservedTagNames {
		if name == rn {
			return TagReg{}, fmt.Errorf("reserved tag name %v used", rn)
		}
	}
	tb := NewTagBase(name, "")
	return TagReg{tb}, nil
}

func MustNewTagReg(name string) Tag {
	t, err := NewTagReg(name)
	if err != nil {
		panic(err)
	}
	return t
}

// NewTagWithValue creates a new Tag with a value
func NewTagRegWithValue(name, value string) (Tag, error) {
	for _, rn := range ReservedTagNames {
		if name == rn {
			return TagReg{}, fmt.Errorf("reserved tag name %v used", rn)
		}
	}
	tb := NewTagBase(name, value)
	return TagReg{tb}, nil
}

func MustNewTagRegWithValue(name, value string) Tag {
	t, err := NewTagRegWithValue(name, value)
	if err != nil {
		panic(err)
	}
	return t
}

func (t TagReg) Has(t2 Tag) bool {
	return t.Name == t2.GetName() && t.Value == t2.GetValue()
}

// ------------------------------------------
type TagWithout struct{ TagBase }

var _ Tag = TagWithout{}

var WithoutKeyword = "WITHOUT"

func NewTagWithout(withoutName string) Tag {
	tb := NewTagBase(WithoutKeyword, withoutName)
	return TagWithout{tb}
}

func (t TagWithout) Has(t2 Tag) bool {
	return !(t.Value == t2.GetName())
}

//_______________________________________________________

// NOTE all tag types must be registered within this function
func ParseTagFromString(in string) Tag {
	splt := strings.Split(in, "=")
	if len(splt) == 2 {
		switch splt[0] {
		case WithoutKeyword:
			return NewTagWithout(splt[1])
		default:
			t, err := NewTagRegWithValue(splt[0], splt[1])
			if err != nil {
				log.Fatal(err)
			}
			return t
		}
	}
	t, err := NewTagReg(in)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

//_______________________________________________________

func ConcatAllContentFromTags(tags []Tag) (content []byte, found bool) {
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

func ParseClumpedTags(clumpedTags string) []Tag {
	trim := strings.TrimPrefix(clumpedTags, ",")
	trim = strings.TrimSuffix(trim, ",")
	split := strings.Split(trim, ",")
	return ParseStringTags(split)
}

func ParseStringTags(strTags []string) []Tag {
	var out []Tag
	for _, s := range strTags {
		trim := strings.TrimSpace(s)
		if len(trim) > 0 {
			out = append(out, ParseTagFromString(trim))
		}
	}
	return out
}

//_______________________________________________________

func (idea Idea) HasTag(tag Tag) bool {
	for _, ideaTag := range idea.Tags {
		if tag.Has(ideaTag) {
			return true
		}
	}
	return false
}

// returns true if the idea contains all the input tags
func (idea Idea) HasTags(tags []Tag) bool {
	for _, tag := range tags {
		hasTag := false
		for _, ideaTag := range idea.Tags {
			if tag.Has(ideaTag) {
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

// returns true if the idea contains any of the input tags
func (idea Idea) HasAnyOfTags(tags []Tag) bool {
	for _, tag := range tags {
		for _, ideaTag := range idea.Tags {
			if tag.Has(ideaTag) {
				return true
			}
		}
	}
	return false
}

// get all the tags clumped together
func (idea Idea) GetClumpedTags() (out string) {
	for _, tag := range idea.Tags {
		if len(out) == 0 {
			out = tag.String()
		} else {
			out += "," + tag.String()
		}
	}
	return out
}

// rename the tag on this idea
func (idea *Idea) RenameTag(from, to Tag) {
	for i, tag := range idea.Tags {
		if tag.String() == from.String() {
			idea.Tags[i] = to
			break
		}
	}
}

// remove the tag on this idea
func (idea *Idea) RemoveTag(tagToRemove Tag) {
	if len(idea.Tags) == 1 && idea.Tags[0] == tagToRemove {
		log.Fatalf("cannot remove the final tag of %v, aborting", idea.Filename)
	}
	for i, tag := range idea.Tags {
		if tag.String() == tagToRemove.String() {
			idea.Tags = append(idea.Tags[:i], idea.Tags[i+1:]...)
			break
		}
	}
}

// add the tag on this idea
func (idea *Idea) AddTag(tag Tag) {
	if idea.HasTag(tag) {
		return
	}
	idea.Tags = append(idea.Tags, tag)
}

// -----------------------------------

func CombineClumpedTags(ct1, ct2 string) string {
	if ct1 == "" {
		return ct2
	} else if ct2 == "" {
		return ct1
	}
	return ct1 + "," + ct2
}
