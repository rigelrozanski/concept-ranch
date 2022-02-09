package idea

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"

	cmn "github.com/rigelrozanski/common"
)

type Tag interface {
	GetName() string
	GetValue() string
	String() string
	Includes(Idea) bool
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
var st = &specialTags{} // global instantiation used throughout

type specialTags struct {
	sts []specialTagsRoute
}

type newTagFn func(keyword, name string) ([]Tag, error)

type specialTagsRoute struct {
	newFn    newTagFn
	keywords []string
}

// ReserveTag reserves a tag
func (s *specialTags) registerTags(fn newTagFn, tagKeywords ...string) {
	s.sts = append(s.sts, specialTagsRoute{fn, tagKeywords})
}

// ReserveTag reserves a tag
func (s *specialTags) verifyUnreserved(tagName string) error {
	for _, str := range s.sts {
		for _, keyword := range str.keywords {
			if tagName == keyword {
				return fmt.Errorf("reserved tag name %v used", tagName)
			}
		}
	}
	return nil
}

// ReserveTag reserves a tag
func (s *specialTags) getNewFn(tagKeyword string) (fn newTagFn, found bool) {
	for _, str := range s.sts {
		for _, keyword := range str.keywords {
			if tagKeyword == keyword {
				return str.newFn, true
			}
		}
	}
	return fn, false
}

// ------------------------------------------
type TagReg struct{ TagBase }

var _ Tag = TagReg{}

// NewTagWithValue creates a new Tag with a value
func NewTagReg(name, value string) ([]Tag, error) {
	err := st.verifyUnreserved(name)
	tb := NewTagBase(name, value)
	return []Tag{TagReg{tb}}, err
}

func MustNewTagReg(name, value string) Tag {
	t, err := NewTagReg(name, value)
	if err != nil {
		panic(err)
	}
	return t[0]
}

func (t TagReg) Includes(idea Idea) bool {
	for _, t2 := range idea.Tags {
		if t.Name == t2.GetName() && t.Value == t2.GetValue() {
			return true
		}
	}
	return false
}

// ------------------------------------------

func splitIfArray(in string) []string {
	inChs := []rune(in)
	if len(inChs) < 3 || inChs[0] != '[' ||
		inChs[len(inChs)-1] != ']' {
		return []string{in}
	}
	// remove curly brackets, and split
	inNoBrac := string(inChs[1 : len(inChs)-1])
	return strings.FieldsFunc(inNoBrac, func(ch rune) bool {
		return (ch == ',' || ch == ' ')
	})
}

// ------------------------------------------
type TagWithout struct{ TagBase }

var _ Tag = TagWithout{}

var WithoutKeyword = "WITHOUT"

func init() { st.registerTags(NewTagWithout, WithoutKeyword) }

func NewTagWithout(keyword, without string) ([]Tag, error) {
	cws := splitIfArray(without)
	tags := []Tag{}
	for _, cw := range cws {
		tags = append(tags,
			TagWithout{
				NewTagBase(keyword, cw),
			},
		)
	}
	return tags, nil
}

func (t TagWithout) Includes(idea Idea) bool {
	for _, t2 := range idea.Tags {
		if t.Value == t2.GetName() {
			return false
		}
	}
	return true
}

// ------------------------------------------
type TagAll struct{ TagBase }

var _ Tag = TagAll{}
var (
	AllAliveKeyword    = "ALL"
	AllConsumedKeyword = "ALL-CONSUMED"
	AllZombieKeyword   = "ALL-ZOMBIE"
)

func init() {
	st.registerTags(NewTagAll,
		AllAliveKeyword, AllConsumedKeyword, AllZombieKeyword)
}

func NewTagAll(keyword, _ string) ([]Tag, error) {
	return []Tag{TagAll{NewTagBase(keyword, "")}}, nil
}

func (t TagAll) Includes(idea Idea) bool {
	switch t.GetName() {
	case AllAliveKeyword:
		return idea.Cycle == CycleAlive
	case AllConsumedKeyword:
		return idea.Cycle == CycleConsumed
	case AllZombieKeyword:
		return idea.Cycle == CycleZombie
	}
	return true
}

// ------------------------------------------
type TagContains struct {
	TagBase
	DoesNotContain  bool
	CaseInsensitive bool
}

var _ Tag = TagContains{}
var (
	ContainsKeyword     = "CONTAINS"
	ContainsCIKeyword   = "CONTAINS-CI"
	NoContainsKeyword   = "NO-CONTAINS"
	NoContainsCIKeyword = "NO-CONTAINS-CI"
)

func init() {
	st.registerTags(NewTagContains,
		ContainsKeyword, ContainsCIKeyword,
		NoContainsKeyword, NoContainsCIKeyword)
}

func NewTagContains(keyword, containsWhat string) ([]Tag, error) {
	cws := splitIfArray(containsWhat)
	tags := []Tag{}
	for _, cw := range cws {
		tb := NewTagBase(keyword, cw)
		switch keyword {
		case ContainsKeyword:
			tags = append(tags, TagContains{tb, false, false})
		case ContainsCIKeyword:
			tags = append(tags, TagContains{tb, false, true})
		case NoContainsKeyword:
			tags = append(tags, TagContains{tb, true, false})
		case NoContainsCIKeyword:
			tags = append(tags, TagContains{tb, true, true})
		}
	}
	return tags, nil
}

func (t TagContains) Includes(idea Idea) bool {
	bz := idea.GetContent()
	fnToLower := func(in string) string { return in }
	if t.CaseInsensitive {
		fnToLower = strings.ToLower
	}
	res := strings.Contains(fnToLower(string(bz)), fnToLower(t.Value))
	if t.DoesNotContain {
		return !res
	}
	return res
}

// ------------------------------------------
type TagDates struct {
	TagBase
	startDate time.Time
	endDate   time.Time
}

var _ Tag = TagDates{}
var (
	CreatedDateKeyword   = "DATE"
	CreatedYearKeyword   = "YEAR"
	CreatedDatesKeyword  = "DATES"
	EditedDateKeyword    = "EDIT-DATE"
	EditedDatesKeyword   = "EDIT-DATES"
	ConsumedDateKeyword  = "CONSUMED-DATE"
	ConsumedDatesKeyword = "CONSUMED-DATES"
)

func init() {
	st.registerTags(NewTagDates,
		CreatedDateKeyword, CreatedYearKeyword,
		CreatedDatesKeyword, EditedDateKeyword, EditedDatesKeyword,
		ConsumedDateKeyword, ConsumedDatesKeyword)
}

// can take either a single date or a [date,range]
func NewTagDates(keyword, date string) ([]Tag, error) {
	dateRangeStr := splitIfArray(date)
	if len(dateRangeStr) != 2 {
		// assume a single date entered
		dateRangeStr = []string{date, date}
	}

	// special logic to allow for the date to just be a year
	if len(dateRangeStr[0]) == 4 {
		dateRangeStr[0] = dateRangeStr[0] + "-01-01"
	}
	if len(dateRangeStr[1]) == 4 {
		dateRangeStr[1] = dateRangeStr[1] + "-12-31"
	}

	dateTimeStart, err := cmn.ParseYYYYdMMdDD(dateRangeStr[0])
	if err != nil {
		return []Tag{}, err
	}
	dateTimeEnd, err := cmn.ParseYYYYdMMdDD(dateRangeStr[1])
	if err != nil {
		return []Tag{}, err
	}
	tb := NewTagBase(keyword, date)
	return []Tag{TagDates{
		TagBase:   tb,
		startDate: dateTimeStart,
		endDate:   dateTimeEnd,
	}}, nil
}

func (t TagDates) Includes(idea Idea) bool {
	var ideaDate time.Time
	switch t.Name {
	case CreatedDateKeyword, CreatedYearKeyword, CreatedDatesKeyword:
		ideaDate = idea.Created
	case EditedDateKeyword, EditedDatesKeyword:
		ideaDate = idea.Edited
	case ConsumedDateKeyword, ConsumedDatesKeyword:
		ideaDate = idea.Consumed
	default:
		panic("unknown date kind")
	}
	if (t.startDate.Before(ideaDate) || t.startDate.Equal(ideaDate)) &&
		(ideaDate.Before(t.endDate) || ideaDate.Equal(t.endDate)) {
		return true
	}
	return false
}

//_______________________________________________________

// NOTE all tag types must be registered within this function
func ParseTagFromString(in string) []Tag {
	keyword, value := in, ""
	splt := strings.Split(in, "=")
	if len(splt) == 2 {
		keyword, value = splt[0], splt[1]
	}
	fn, found := st.getNewFn(keyword)
	if !found {
		fn = NewTagReg
	}
	ts, err := fn(keyword, value)
	if err != nil {
		log.Fatal(err)
	}
	return ts
}

func ParseFirstTagFromString(in string) Tag {
	return ParseTagFromString(in)[0]
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

// parse clumped tags seperated by spaces or commas
func ParseClumpedTags(clumpedTags string) []Tag {
	trim := strings.TrimPrefix(clumpedTags, ",")
	trim = strings.TrimSuffix(trim, ",")
	trim = strings.TrimSuffix(trim, " ")
	trim = strings.TrimPrefix(clumpedTags, " ")

	// split strings by "," not contained within square brackets
	split := []string{}
	collecting := ""
	bracCount := 0
	for _, ch := range trim {
		if ch == '[' {
			bracCount++
		} else if ch == ']' {
			bracCount--
		}
		if bracCount == 0 && (ch == ',' || ch == ' ') {
			if len(collecting) > 0 {
				split = append(split, collecting)
			}
			collecting = ""
		} else {
			collecting += string(ch)
		}
	}
	if len(collecting) > 0 {
		split = append(split, collecting)
	}

	return ParseStringTags(split)
}

func ParseStringTags(strTags []string) []Tag {
	var out []Tag
	for _, s := range strTags {
		trim := strings.TrimSpace(s)
		if len(trim) > 0 {
			out = append(out, ParseTagFromString(trim)...)
		}
	}
	return out
}

//_______________________________________________________

func (idea Idea) HasTag(tag Tag) bool {
	return tag.Includes(idea)
}

// returns true if the idea contains all the input tags
func (idea Idea) HasTags(tags []Tag) bool {
	for _, tag := range tags {
		if !tag.Includes(idea) {
			return false
		}
	}
	return true
}

// returns true if the idea contains any of the input tags
func (idea Idea) HasAnyOfTags(tags []Tag) bool {
	for _, tag := range tags {
		if tag.Includes(idea) {
			return true
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
func (idea *Idea) RemoveTags(tagsToRemove []Tag) {
	for _, tagToRemove := range tagsToRemove {
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
}

// add the tag on this idea
func (idea *Idea) AddTags(tags []Tag) {
	for _, tag := range tags {
		if idea.HasTag(tag) {
			return
		}
		idea.Tags = append(idea.Tags, tag)
	}
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
