package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rigelrozanski/qi/lib"
)

// filestructure:
//                ./ideas/a,123456,YYYYMMDD,eYYYYMMDD,cYYYYMMDD,c432978,c543098...,tag1,tag2,tag3...
//                ./qi
//                ./log
//                ./config
//                ./working_files
//                ./working_content
// 123456 = id
// c123456 = consumes-id
// YYYYMMDD = creation date
// eYYYYMMDD = last edited date
// cYYYYMMDD = consumed date

//keywords used throughout qi
const (
	keyHelp1, keyHelp2 = "--help", "-h"
	keyCat             = "cat"
	keyScan            = "scan"
	keyTranscribe      = "transcribe"
	keyTranscribeMany  = "transcribe-many"
	keyConsume         = "consume"
	keyZombie          = "zombie"
	keyConsumes        = "consumes"
	keyNew             = "new"
	keyRm              = "rm"
	keyCp              = "cp"
	keyTags            = "tags"
	keyKillTag         = "kill-tag"
	keyRenameTag       = "rename-tag"
	keyAddTag          = "add-tag"
	keyDestroyTag      = "destroy-tag"
	keyLs              = "ls"
	keyLsFiles         = "ls-files"
	keyLog             = "log"
	keyPdfBackup       = "pdf-backup"

	help = `
/|||||\ |-o-o-~|
qi --------------------------------------> edit the tagless master quick ideas board in vim
qi [tags...] [entry] --------------------> quick entry to a new idea
qi [query] ------------------------------> open a vim tab with the contents of the query 
qi cat [query] --------------------------> print idea(s) contents' to console
qi scan [image_loco] [op_tag] -----------> scan a bunch of images as untranscribed ideas, optionally append a tag to all
qi transcribe [id] ----------------------> transcribe a random untranscribed image or a specific image by id
qi transcribe-many [op_tags...] ---------> transcribe many images one after another, optionally transcribe within a set of tags
qi consume [id] [entry] -----------------> quick consumes the given id into a new entry
qi consumes [consumed-id] [consumer-id] -> set the consumption of existing ideas
qi zombie [id] --------------------------> "unconsume" an idea based on id
qi new [tags...] ------------------------> create a new idea with the provided tags
qi rm [id] ------------------------------> remove an idea by id
qi cp [id] ------------------------------> duplicate an idea at the provided id
qi tags [id] ----------------------------> list the tags for a given id
qi kill-tag [id] ------------------------> remove a tag from an idea by id
qi rename-tag [id] ----------------------> rename all instances of a tag for all ideas
qi add-tag [id] -------------------------> add a tag from an idea by id
qi destroy-tag --------------------------> remove all instances of a tag for all ideas
qi ls -----------------------------------> list all tags used
qi ls-files -----------------------------> list all files
qi pdf-backup ---------------------------> backup best ideas to a printable pdf

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
	args := os.Args[1:]

	// for the master qi file for quick entry
	if len(args) == 0 {
		openText(lib.QiFile)
		return
	}

	var err error
	switch args[0] {
	case keyHelp1, keyHelp2:
		fmt.Println(help)
	case keyCat:
		QuickQuery(args[1])
	case keyScan:

	case keyTranscribe:

	case keyTranscribeMany:

	case keyConsume:

	case keyConsumes:

	case keyZombie:

	case keyNew:
		EnsureLen(args, 2)
		NewEmptyEntry(args[1])

	case keyRm:
		RemoveByID(args[1])

	case keyCp:
		CopyByID(args[1])

	case keyTags:

	case keyKillTag:

	case keyRenameTag:

	case keyAddTag:

	case keyDestroyTag:

	case keyLs:
		ListAllTags()

	case keyLsFiles:
		ListAllFiles()

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
		log.Fatalf("expected %v args", enLen)
	}
}
