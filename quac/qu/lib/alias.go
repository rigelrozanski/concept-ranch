// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/rigelrozanski/wranch/quac/qu/lib/idea
package lib

import (
	"github.com/rigelrozanski/wranch/quac/qu/lib/idea"
)

const (
	Last          = idea.Last
	CycleAlive    = idea.CycleAlive
	CycleConsumed = idea.CycleConsumed
	CycleZombie   = idea.CycleZombie
	KindText      = idea.KindText
	KindImage     = idea.KindImage
	KindAudio     = idea.KindAudio
)

var (
	// functions aliases
	GetNextID                = idea.GetNextID
	IncrementID              = idea.IncrementID
	ParseID                  = idea.ParseID
	ParseIDNoLogLast         = idea.ParseIDNoLogLast
	ParseIDOp                = idea.ParseIDOp
	NewNonConsumingTextIdea  = idea.NewNonConsumingTextIdea
	NewTextIdea              = idea.NewTextIdea
	NewIdeaFromFile          = idea.NewIdeaFromFile
	NewConsumingTextIdea     = idea.NewConsumingTextIdea
	NewIdeaFromFilename      = idea.NewIdeaFromFilename
	GetAllIdeasNonConsuming  = idea.GetAllIdeasNonConsuming
	GetAllIdeas              = idea.GetAllIdeas
	ConcatAllContentFromTags = idea.ConcatAllContentFromTags
	TodayDate                = idea.TodayDate
	GetKind                  = idea.GetKind
	IdStr                    = idea.IdStr

	// variable aliases
	IdeasDir   = idea.IdeasDir
	ConfigFile = idea.ConfigFile
	LastIdFile = idea.LastIdFile
)

type (
	Idea  = idea.Idea
	Ideas = idea.Ideas
)