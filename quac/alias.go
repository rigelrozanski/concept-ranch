// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/rigelrozanski/thranch/quac/idea
package quac

import (
	"github.com/rigelrozanski/thranch/quac/idea"
)

const (
	Last          = idea.Last
	CycleAlive    = idea.CycleAlive
	CycleConsumed = idea.CycleConsumed
	CycleZombie   = idea.CycleZombie
	KindText      = idea.KindText
	KindImage     = idea.KindImage
	KindAudio     = idea.KindAudio
	KindEnText    = idea.KindEnText
)

var (
	// functions aliases
	ValidateFilenameAsIdea   = idea.ValidateFilenameAsIdea
	GetIdByFilename          = idea.GetIdByFilename
	GetNextID                = idea.GetNextID
	IncrementID              = idea.IncrementID
	ParseID                  = idea.ParseID
	ParseIDNoLogLast         = idea.ParseIDNoLogLast
	ParseIDOp                = idea.ParseIDOp
	PrependLast              = idea.PrependLast
	GetLastIDs               = idea.GetLastIDs
	NewNonConsumingTextIdea  = idea.NewNonConsumingTextIdea
	NewNonConsumingAudioIdea = idea.NewNonConsumingAudioIdea
	NewTextIdea              = idea.NewTextIdea
	NewAudioIdea             = idea.NewAudioIdea
	NewIdea                  = idea.NewIdea
	NewIdeaFromFile          = idea.NewIdeaFromFile
	NewConsumingTextIdea     = idea.NewConsumingTextIdea
	NewIdeaFromFilepath      = idea.NewIdeaFromFilepath
	NewIdeaFromFilename      = idea.NewIdeaFromFilename
	GetAllIdeasNonConsuming  = idea.GetAllIdeasNonConsuming
	TagUsedInNonConsuming    = idea.TagUsedInNonConsuming
	GetAllIdeas              = idea.GetAllIdeas
	NewTagBase               = idea.NewTagBase
	NewTagReg                = idea.NewTagReg
	MustNewTagReg            = idea.MustNewTagReg
	NewTagWithout            = idea.NewTagWithout
	NewTagAll                = idea.NewTagAll
	NewTagContains           = idea.NewTagContains
	NewTagDates              = idea.NewTagDates
	ParseTagFromString       = idea.ParseTagFromString
	ParseFirstTagFromString  = idea.ParseFirstTagFromString
	ConcatAllContentFromTags = idea.ConcatAllContentFromTags
	ParseClumpedTags         = idea.ParseClumpedTags
	ParseStringTags          = idea.ParseStringTags
	CombineClumpedTags       = idea.CombineClumpedTags
	TodayDate                = idea.TodayDate
	GetKind                  = idea.GetKind
	IdStr                    = idea.IdStr

	// variable aliases
	WithoutKeyword       = idea.WithoutKeyword
	AllAliveKeyword      = idea.AllAliveKeyword
	AllConsumedKeyword   = idea.AllConsumedKeyword
	AllZombieKeyword     = idea.AllZombieKeyword
	ContainsKeyword      = idea.ContainsKeyword
	ContainsCIKeyword    = idea.ContainsCIKeyword
	NoContainsKeyword    = idea.NoContainsKeyword
	NoContainsCIKeyword  = idea.NoContainsCIKeyword
	CreatedDateKeyword   = idea.CreatedDateKeyword
	CreatedYearKeyword   = idea.CreatedYearKeyword
	CreatedDatesKeyword  = idea.CreatedDatesKeyword
	EditedDateKeyword    = idea.EditedDateKeyword
	EditedDatesKeyword   = idea.EditedDatesKeyword
	ConsumedDateKeyword  = idea.ConsumedDateKeyword
	ConsumedDatesKeyword = idea.ConsumedDatesKeyword
	IdeasDir             = idea.IdeasDir
	ConfigFile           = idea.ConfigFile
	LastIdFile           = idea.LastIdFile
)

type (
	Idea        = idea.Idea
	Ideas       = idea.Ideas
	Tag         = idea.Tag
	TagBase     = idea.TagBase
	TagReg      = idea.TagReg
	TagWithout  = idea.TagWithout
	TagAll      = idea.TagAll
	TagContains = idea.TagContains
	TagDates    = idea.TagDates
)
