package main

import (
	"fmt"
	"os"
	"strings"
)

// filestructure:
//                ./ideas/123456,c234567,YYYYMMDD,eYYYYMMDD,cYYYYMMDD,tag1,tag2,tag3...
//                ./consumed/123456,c234567,YYYYMMDD,eYYYYMMDD,cYYYYMMDD,tag1,tag2,tag3...
//                ./qi
//                ./log
//                ./config
//                ./working
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
	keyLog             = "log"
	keyPdfBackup       = "pdf-backup"

	help = `
/|||||\ |-o-o-~|
qi --------------------------------------> edit the tagless master quick ideas board in vim
qi [tag1,tag2...] [entry] ---------------> quick entry to a new idea
qi [query] ------------------------------> open a vim tab with the contents of the query 
qi cat [query] --------------------------> print idea(s) contents' to console
qi scan [image_loco] [op_tag] -----------> scan a bunch of images as untranscribed ideas, optionally append a tag to all
qi transcribe [id] ----------------------> transcribe a random untranscribed image or a specific image by id
qi transcribe-many [op_tags] ------------> transcribe many images one after another, optionally transcribe within a set of tags
qi consume [id] [op_entry] --------------> consumes the given id into a 
qi consumes [consumed-id] [consumer-id] -> set the consumption of existing ideas
qi new [tags] ---------------------------> create a new idea with the provided tags
qi rm [id] ------------------------------> remove an idea by id
qi cp [id] ------------------------------> duplicate an idea at the provided id
qi tags [id] ----------------------------> list the tags for a given id
qi kill-tag [id] ------------------------> remove a tag from an idea by id
qi rename-tag [id] ----------------------> rename all instances of a tag for all ideas
qi add-tag [id] -------------------------> add a tag from an idea by id
qi destroy-tag --------------------------> remove all instances of a tag for all ideas
qi ls -----------------------------------> list all tags used
qi log ----------------------------------> list the log
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
[tags] --------- A list of tags seperated by commas (such as "tag1,tag2,tag3")
[entry] -------- Either raw input text or for untranscribed input a directory to an image
`
)

func main() {
	args := os.Args[1:]

	// for the master quick entry
	if len(args) == 1 {
		//err := edit(defaultQI)
		//if err != nil {
		//fmt.Println(err)
		//}
	}

	var err error
	switch args[0] {
	case keyHelp1, keyHelp2:
		fmt.Println(help)
	case keyCat:

	case keyScan:

	case keyTranscribe:

	case keyTranscribeMany:

	case keyConsume:

	case keyConsumes:

	case keyNew:

	case keyRm:

	case keyCp:

	case keyTags:

	case keyKillTag:

	case keyRenameTag:

	case keyAddTag:

	case keyDestroyTag:

	case keyLs:
		ListAllTags()

	case keyLog:

	default:
		if len(args) == 1 { // quick query
			QuickQuery(args[0])
		} else if len(args) >= 2 { // quick entry
			//var entry string
			//for i := 1; i < len(args); i++ {
			//entry += " " + args[i]
			//}
			QuickEntry(args[0], strings.Join(args[1:], " "))
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}
