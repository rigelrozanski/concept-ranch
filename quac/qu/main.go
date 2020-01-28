package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rigelrozanski/wranch/quac/qu/lib"
)

// filestructure:
//                ./ideas/a,123456,YYYYMMDD,eYYYYMMDD,cYYYYMMDD,c432978,c543098...,tag1,tag2,tag3...
//                ./qu
//                ./log
//                ./config
//                ./working_files
//                ./working_content
// 123456 = id
// c123456 = consumes-id
// YYYYMMDD = creation date
// eYYYYMMDD = last edited date
// cYYYYMMDD = consumed date

//keywords used throughout qu
const (
	keyHelp1, keyHelp2 = "--help", "-h"
	keyCat             = "cat"
	keyScan            = "scan"
	keyTranscribe      = "transcribe"
	keyTranscribeMany  = "transcribe-many"
	keyConsume         = "consume"
	keyConsumes        = "consumes"
	keyZombie          = "zombie"
	keyLineage         = "lineage"
	keyNew             = "new"
	keyRm              = "rm"
	keyCp              = "cp"
	keyTags            = "tags"
	keyKillTag         = "kill-tag"
	keyRenameTag       = "rename-tag"
	keyAddTag          = "add-tag"
	keyDestroyTag      = "destroy-tag"
	keyLsTags          = "lst"
	keyLsFiles         = "lsf"
	keyPDFBackup       = "pdf-backup"

	help = `
/|||||\ |-o-o-~|
qu --------------------------------------> edit the tagless master quick ideas board in vim
qu [tags...] [entry] --------------------> quick entry to a new idea
qu [query] ------------------------------> open a vim tab with the contents of the query 
qu cat [query] --------------------------> print idea(s) contents' to console
qu scan [image_loco] [op_tag] -----------> scan a bunch of images as untranscribed ideas, optionally append a tag to all
qu transcribe [id] ----------------------> transcribe a random untranscribed image or a specific image by id
qu transcribe-many [op_tags...] ---------> transcribe many images one after another, optionally transcribe within a set of tags
qu consume [id] [op_entry] --------------> quick consumes the given id into a new entry
qu consumes [consumed-id] [consumer-id] -> set the consumption of existing ideas
qu zombie [id] --------------------------> "unconsume" an idea based on id
qu lineage [id] -------------------------> show the consumtion lineage  
qu new [tags...] ------------------------> create a new idea with the provided tags
qu rm [id] ------------------------------> remove an idea by id
qu cp [id] ------------------------------> duplicate an idea at the provided id
qu tags [id] ----------------------------> list the tags for a given id
qu kill-tag [id] [tag] ------------------> remove a tag from an idea by id
qu add-tag [id] [tag] -------------------> add a tag from an idea by id
qu rename-tag [from-tag] [to-tag]--------> rename all instances of a tag for all ideas
qu destroy-tag [tag] --------------------> remove all instances of a tag for all ideas
qu lst [op_tags] ------------------------> list all tags used, optionally which share a set of common tags
qu lsf [op_tags] ------------------------> list all files, optionally which contain provided tags
qu pdf-backup ---------------------------> backup best ideas to a printable pdf

Explanation of some terms:
[op_ ----------- An optional input 
[id] ----------- Either a 6 digit number (such as "123456") or the keyword "lastid" or "lastXid" where X is an integer
[query] -------- Either an [id] or a list of tags seperated by commas (such as "tag1,tag2,tag3") 
                     Additionally, the following special tags can be used:
                        consumed
                        created_afterYYYYMMDD (where YYYYMMDD is a date)
                        created_beforeYYYYMMDD 
                        edited_afterYYYYMMDD 
                        edited_beforeYYYYMMDD 
                        consumed_afterYYYYMMDD 
                        consumed_beforeYYYYMMDD 
[tag] ---------- A catagory to query or organize your ideas with
[tags...] ------ A list of tags seperated by commas (such as "tag1,tag2,tag3")
                     Additionally special tags can be used:
					    consumesXXXXXX (where XXXXXX is the id of the idea being consumed)
[entry] -------- Either raw input text or for untranscribed input a directory to an image/audio sample 
`
)

func main() {
	lib.Initialize(os.ExpandEnv("$HOME/.wranch_config.txt"))
	args := os.Args[1:]

	// for the master qu file for quick entry
	if len(args) == 0 {
		openText(lib.QuFile)
		return
	}

	var err error
	switch args[0] {
	case keyHelp1, keyHelp2:
		fmt.Println(help)
	case keyCat:
		QuickQuery(args[1])
	case keyScan:
		lib.Scan(args[1])

	case keyTranscribe:

	case keyTranscribeMany:

	case keyConsume:
		switch len(args) {
		case 2:
			Consume(args[1], "")
		case 3:
			Consume(args[1], args[2])
		default:
			EnsureLen(args, 2)
		}
	case keyConsumes:
		EnsureLen(args, 3)
		Consumes(args[1], args[2])
	case keyZombie:
		EnsureLen(args, 2)
		Zombie(args[1])
	case keyLineage:
		EnsureLen(args, 2)
		Lineage(args[1])
	case keyNew:
		EnsureLen(args, 2)
		NewEmptyEntry(args[1])
	case keyRm:
		EnsureLen(args, 2)
		RemoveByID(args[1])
	case keyCp:
		EnsureLen(args, 2)
		CopyByID(args[1])
	case keyTags:
		EnsureLen(args, 2)
		ListTagsByID(args[1])
	case keyKillTag:
		EnsureLen(args, 3)
		KillTagByID(args[1], args[2])
	case keyAddTag:
		EnsureLen(args, 3)
		AddTagByID(args[1], args[2])
	case keyRenameTag:
		EnsureLen(args, 3)
		RenameTag(args[1], args[2])
	case keyDestroyTag:
		EnsureLen(args, 2)
		DestroyTag(args[1])
	case keyLsTags:
		if len(args) == 1 {
			ListAllTags()
		} else {
			ListAllTagsWithTags(args[1])
		}
	case keyLsFiles:
		if len(args) == 1 {
			ListAllFiles()
		} else {
			ListAllFilesWithTags(args[1])
		}
	case keyPDFBackup:
		lib.ExportToPDF()
	default:
		if len(args) == 1 { // quick query
			MultiOpen(args[0])
		} else if len(args) >= 2 { // quick entry
			QuickEntry(args[0], strings.Join(args[1:], " "))
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func EnsureLen(args []string, enLen int) {
	if len(args) < enLen {
		log.Fatalf("expected at least %v args", enLen)
	}
}
