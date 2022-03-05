package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rigelrozanski/thranch/quac"
)

// filestructure:
//       ./ideas/a,123456,YYYYMMDD,eYYYYMMDD,cYYYYMMDD,c432978,c543098...,tag1,tag2,tag3...
//       ./qu
//       ./log
//       ./config
//       ./working_files
//       ./working_content
//
// 123456    = id
// c123456   = consumes-id
// YYYYMMDD  = creation date
// eYYYYMMDD = last edited date
// cYYYYMMDD = consumed date

//keywords used throughout qu
const (
	keyHelp1, keyHelp2 = "--help", "-h"
	keyQuickEntry      = "qe"
	keyCat             = "cat"
	keyStats           = "stats"
	keyScan            = "scan"
	keyWaterCloset     = "wc"
	keyTranscribe      = "transcribe"
	keyTagUntagged     = "tag-untagged"
	keyConsume         = "consume"
	keyConsumes        = "consumes"
	keyZombie          = "zombie"
	keyLineage         = "lineage"
	keyNew             = "new"
	keyManualEntry     = "manual-entry"
	keySetEncryption   = "set-encryption"
	keyRm              = "rm"
	keyEmptyTrash      = "empty-trash"
	keyCp              = "cp"
	keyRemoveTag       = "rm-tag"
	keyRenameTagToMany = "rename-tag-to-many"
	keyAddTag          = "add-tag"
	keyAddTagToMany    = "add-tag-to-many"
	keyAddTags         = "add-tags"
	keyDestroyTag      = "destroy-tag"
	keyCommonTags      = "common-tags"
	keyLS              = "ls"
	keyLSFile          = "lsfl"
	keySelectFiles     = "sel"
	keyPDFBackup       = "pdf-backup"
	keyForceSplit      = "force-split"
	keyOpenWorking     = "open-working"
	keySaveWorking     = "save-working"

	help = `
/|||||\ |-o-o-~|
🦆 🦆 🦆 ✍️  🏁

qu ---------------------------------------> edit the tagless master idea in vim
qu [force-split] <query> -----------------> open a vim tab with the contents of the query 
qu new <tags> ----------------------------> create a new idea with the provided tags
qu cat <query> ---------------------------> print idea(s) contents' to console
qu ls [query] ----------------------------> list ideas which match the [query], if no query is
                                              provided list recently opened ideas
qu cp <id> -------------------------------> duplicate an idea at the provided id
qu rm <id1-id2> --------------------------> remove an idea by id or id-range to the trash can
qu empty-trash ---------------------------> empty the trash can optionally appending tags to all
                                              or specific image(s) by query
-- ENTRY --
qu scan <dir/file> [tags] ----------------> add provided image(s) to untranscribed ideas, 
qu tag-untagged --------------------------> iterate and add tags to ideas with the tag "UNTAGGED"
qu transcribe [query] --------------------> transcribe either a random untranscribed image 
qu wc ------------------------------------> (water closet) tag-untagged all then transcribe
qu manual-entry [tags] -------------------> interactive manual entry common tags may be entered 
qu consume <id> [entry] ------------------> quick consumes the given id into a new entry
qu consumes <consumed-id> <consumer-id> --> set the consumption of existing ideas
qu zombie <id> ---------------------------> "unconsume" an idea based on id
qu lineage <id> --------------------------> show the consumtion lineage  

-- TAGS MANAGEMENT --
qu common-tags [tags] --------------------> list all tags which share a set of common [tags]
qu rm-tag <id> <tag> ---------------------> remove a tag from an idea by id
qu add-tags <id> <tags> ------------------> add (a) tag(s) to an idea by id (alias: add-tag)
qu add-tag-to-many <newtag> <tags..> -----> add a <newtag> to all ideas with any of <tags...>
qu rename-tag-to-many <from-tag> <to-tag> > rename all instances of a tag for all ideas
qu destroy-tag <tag> ---------------------> remove all instances of a tag for all ideas

-- OTHER --
qu qe <tags...> <entry> ------------------> quick entry to a new idea
qu set-encryption <id> -------------------> set encryption of existing idea
qu open-working --------------------------> open the working split files to manually correct mistakes
qu save-working --------------------------> save the working split files to manually correct mistakes
qu pdf-backup ----------------------------> backup active ideas to a printable pdf
qu stats ---------------------------------> statistics on your ideas
qu sel [tags]-----------------------------> select the idea from the tags (in cui)
qu lsfl [query] --------------------------> list all files by file location

Explanation of some terms:
[...], <...> --- optional input, required input
id ------------- either a 6 digit number (such as "123456") or the keyword "lastid" 
                   or "lastXid" where X is an integer
id1-id2 -------- either just an [id] or a range of ids in the form 123456-222000
query ---------- either an [id], [id1-id2], or a list of tags 
                   seperated by commas (such as "tag1,tag2,tag3") 
tag ------------ a catagory to query or organize your ideas with
                   special tags: FILENAME - if this tag is used the filename of 
				   each entry file will be included as a tag per idea, 
				   errors if entry is raw text input
tags ----------- a list of tags seperated by commas (such as "tag1,tag2,tag3")
                 SPECIAL TAGS: 
				   WITHOUT=foo        <- exclude tags 'foo' 
				   CONTAINS=foo       <- include ideas which contain the text 'foo' 
				   CONTAINS-CI=foo    <- same as CONTAINS but case-insensitive
				   NO-CONTAINS=foo    <- excludes ideas which contain the text 'foo' 
				   NO-CONTAINS-CI=foo <- same as NO-CONTAINS but case-insensitive
				   *NOTE: Within these examples 'foo' may also be an array 
				          in the format of ['foo','bar']
entry ---------- either raw input text or source input as a file or directory
force-split ---- if the text "force-split" is included, split view will be used 
                   even if only one entry is found 
`
)

func main() {
	quac.Initialize(os.ExpandEnv("$HOME/.thranch_config"))
	args := os.Args[1:]

	// for the master qu file for quick entry
	if len(args) == 0 {
		quac.OpenText(quac.QuFile)
		return
	}

	switch args[0] {
	case keyHelp1, keyHelp2:
		fmt.Println(help)
	case keyCat:
		QuickQuery(args[1])
	case keyQuickEntry:
		if len(args) >= 2 {
			QuickEntry(args[0], strings.Join(args[1:], " "))
		} else {
			fmt.Println("not enough arguments for a quick entry")
		}
	case keyScan:
		switch len(args) {
		case 1:
			quac.ScanManual("")
		case 2:
			quac.ScanManual(args[1])
		default:
			EnsureLenAtLeast(args, 1)
		}
	case keyTranscribe:
		switch len(args) {
		case 1:
			Transcribe("")
		case 2:
			Transcribe(args[1])
		default:
			EnsureLenAtLeast(args, 1)
		}
	case keyTagUntagged:
		TagUntagged()
	case keyWaterCloset:
		WaterCloset()
	case keyConsume:
		switch len(args) {
		case 2:
			Consume(args[1], "")
		case 3:
			Consume(args[1], args[2])
		default:
			EnsureLenAtLeast(args, 2)
		}
	case keyConsumes:
		EnsureLenAtLeast(args, 3)
		Consumes(args[1], args[2])
	case keyZombie:
		EnsureLenAtLeast(args, 2)
		Zombie(args[1])
	case keyLineage:
		EnsureLenAtLeast(args, 2)
		Lineage(args[1])
	case keyNew:
		EnsureLenAtLeast(args, 2)
		NewEmptyEntry(strings.Join(args[1:], " "))
	case keyManualEntry:
		if len(args) == 1 {
			ManualEntry("")
		} else {
			ManualEntry(args[1])
		}
	case keySetEncryption:
		EnsureLenAtLeast(args, 2)
		SetEncryption(args[1])
	case keyRm:
		EnsureLenAtLeast(args, 2)
		RemoveByID(args[1])
	case keyEmptyTrash:
		EnsureLenAtLeast(args, 1)
		EmptyTrash()
	case keyCp:
		EnsureLenAtLeast(args, 2)
		CopyByID(args[1])
	case keyRemoveTag:
		EnsureLenAtLeast(args, 3)
		RemoveTagByID(args[1], args[2])
	case "rm-tags":
		fmt.Println("didn't you mean rm-tag???")
	case keyAddTag, keyAddTags:
		EnsureLenAtLeast(args, 3)
		AddTagByID(args[1], args[2])
	case keyAddTagToMany:
		EnsureLenAtLeast(args, 3)
		AddTagToMany(args[1], args[2])
	case keyRenameTagToMany:
		EnsureLenAtLeast(args, 3)
		RenameTag(args[1], args[2])
	case keyDestroyTag:
		EnsureLenAtLeast(args, 2)
		DestroyTag(args[1])
	case keyCommonTags:
		if len(args) == 1 {
			ListAllTags()
		} else {
			ListAllTagsWithTags(args[1])
		}
	case keyLS:
		if len(args) == 1 {
			ListAllFilesLast(false)
		} else {
			ListAllFilesWithQuery(args[1])
		}
	case keyLSFile:
		if len(args) == 1 {
			ListAllFilesByLocation()
		} else {
			ListAllFilesByLocationWithQuery(args[1])
		}
	case keySelectFiles:
		if len(args) == 1 {
			EnsureLenAtLeast(args, 2)
		} else {
			ListSelectAllFilesWithQuery(args[1])
		}
	case keyPDFBackup:
		quac.ExportToPDF()
	case keyStats:
		quac.GetStats()
	case keyForceSplit:
		EnsureLenAtLeast(args, 2)
		MultiOpen(args[1], true) // quick entry force split view
	case keyOpenWorking:
		OpenWorking()
	case keySaveWorking:
		SaveWorking()
	default:
		if len(args) == 1 { // quick query
			ListSelectAllFilesWithQueryNoLast(args[0])
		} else {
			fmt.Println("unknown command")
		}
	}
}

func EnsureLenAtLeast(args []string, enLen int) {
	if len(args) < enLen {
		log.Fatalf("expected at least %v args", enLen)
	}
}
