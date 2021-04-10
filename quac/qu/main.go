package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rigelrozanski/thranch/quac"
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
	keyStats           = "stats"
	keyScan            = "scan"
	keyTranscribe      = "transcribe"
	keyRetag           = "retag"
	keyWaterCloset     = "wc"
	keyConsume         = "consume"
	keyConsumes        = "consumes"
	keyZombie          = "zombie"
	keyLineage         = "lineage"
	keyNew             = "new"
	keyManualEntry     = "manual-entry"
	keySetEncryption   = "set-encryption"
	keyRm              = "rm"
	keyCp              = "cp"
	keyTags            = "tags"
	keyKillTag         = "kill-tag"
	keyRemoveTag       = "rm-tag"
	keyRenameTag       = "rename-tag"
	keyRenameTag2      = "tag-rename"
	keyAddTag          = "add-tag"
	keyAddTagToMany    = "add-tag-to-many"
	keyAddTags         = "add-tags"
	keyDestroyTag      = "destroy-tag"
	keyLsTags          = "lst"
	keyLsFiles         = "lsf"
	keyLsLocationFiles = "lsfl"
	keyLs              = "ls"
	keySelectFiles     = "sel"
	keyPDFBackup       = "pdf-backup"
	keyForceSplit      = "force-split"
	keyOpenWorking     = "open-working"
	keySaveWorking     = "save-working"

	help = `
/|||||\ |-o-o-~|
🦆 🦆 🦆 ✍️  🏁

qu --------------------------------------> edit the tagless master quick ideas board in vim
qu <tags...> <entry> --------------------> quick entry to a new idea
qu [force-split] <query> ----------------> open a vim tab with the contents of the query 
qu cat <query> --------------------------> print idea(s) contents' to console
qu scan <dir/file> [tags] ---------------> scan a bunch of images as untranscribed ideas, optionally append tags to all
qu transcribe [query] -------------------> transcribe a random untranscribed image or specific image(s) by query
qu retag --------------------------------> iterate and add tags to ideas with the tag "UNTAGGED"
qu wc -----------------------------------> while sitting on the toilet, either retag or if all are already tagged then transcribe
qu consume <id> [entry] -----------------> quick consumes the given id into a new entry
qu consumes <consumed-id> <consumer-id> -> set the consumption of existing ideas
qu zombie <id> --------------------------> "unconsume" an idea based on id
qu lineage <id> -------------------------> show the consumtion lineage  
qu new <tags> ---------------------------> create a new idea with the provided tags
qu manual-entry [tags] ------------------> interactive manual entry common tags may be entered with this command
qu set-encryption <id> ------------------> set encryption of existing idea
qu rm <id1-id2> -------------------------> remove an idea by id
qu cp <id> ------------------------------> duplicate an idea at the provided id
qu tags <id> ----------------------------> list the tags for a given id
qu kill-tag <id> <tag> ------------------> remove a tag from an idea by id (alternate cmd rm-tag)
qu add-tags <id> <tags> -----------------> add (a) tag(s) to an idea by id (alternate cmd add-tag)
qu add-tag-to-many <newtag> <tags..> ----> add a <newtag> to all ideas with any of <tags...>
qu rename-tag <from-tag> <to-tag>--------> rename all instances of a tag for all ideas
qu destroy-tag <tag> --------------------> remove all instances of a tag for all ideas
qu open-working -------------------------> open the working files to manually correct mistakes
qu save-working -------------------------> save the working files to manually correct mistakes
qu lst  [tags] --------------------------> list all tags used, optionally which share a set of common tags
qu lsf  [query] -------------------------> list all files, optionally which contain provided tags, or the last 9 viewed
qu lsfl [query] -------------------------> list all files by file location
qu ls   [tags]---------------------------> browse all tags
qu sel  [tags]---------------------------> select the idea from the tags (in cui)
qu pdf-backup ---------------------------> backup best ideas to a printable pdf
qu stats --------------------------------> statistics on your ideas

Explanation of some terms:
[...], <...> --- optional input, required input
id ------------- either a 6 digit number (such as "123456") or the keyword "lastid" or "lastXid" where X is an integer
id1-id2 -------- either just an [id] or a range of ids in the form 123456-222000
query ---------- either an [id], [id1-id2], or a list of tags seperated by commas (such as "tag1,tag2,tag3") 
tag ------------ a catagory to query or organize your ideas with
                   special tags: FILENAME - if this tag is used the filename of each entry file will be included 
				                            as a tag per idea. errors if entry is raw text input
tags ----------- a list of tags seperated by commas (such as "tag1,tag2,tag3")
entry ---------- either raw input text or source input as a file or directory
force-split ---- if the text "force-split" is included, split view will be used even if only one entry is found 
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
	case keyScan:
		switch len(args) {
		case 1:
			quac.ScanManual("")
		case 2:
			quac.ScanManual(args[1])
		default:
			EnsureLen(args, 1)
		}
	case keyTranscribe:
		switch len(args) {
		case 1:
			Transcribe("")
		case 2:
			Transcribe(args[1])
		default:
			EnsureLen(args, 1)
		}
	case keyRetag:
		Retag()
	case keyWaterCloset:
		WaterCloset()
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
	case keyManualEntry:
		if len(args) == 1 {
			ManualEntry("")
		} else {
			ManualEntry(args[1])
		}
	case keySetEncryption:
		EnsureLen(args, 2)
		SetEncryption(args[1])
	case keyRm:
		EnsureLen(args, 2)
		RemoveByID(args[1])
	case keyCp:
		EnsureLen(args, 2)
		CopyByID(args[1])
	case keyTags:
		EnsureLen(args, 2)
		ListTagsByID(args[1])
	case keyKillTag, keyRemoveTag:
		EnsureLen(args, 3)
		KillTagByID(args[1], args[2])
	case "rm-tags":
		fmt.Println("didn't you mean rm-tag???")
	case "kill-tags":
		fmt.Println("didn't you mean kill-tag???")
	case keyAddTag, keyAddTags:
		EnsureLen(args, 3)
		AddTagByID(args[1], args[2])
	case keyAddTagToMany:
		EnsureLen(args, 3)
		AddTagToMany(args[1], args[2])
	case keyRenameTag, keyRenameTag2:
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
			ListAllFilesWithQuery(args[1])
		}
	case keyLsLocationFiles:
		if len(args) == 1 {
			ListAllFilesByLocation()
		} else {
			ListAllFilesByLocationWithQuery(args[1])
		}
	case "lsd":
		fmt.Println("fking rip'd")
		os.Exit(0)
	case keyLs:
		if len(args) == 1 {
			quac.Ls("")
		} else {
			quac.Ls(args[1])
		}
	case keySelectFiles:
		if len(args) == 1 {
			EnsureLen(args, 2)
		} else {
			ListSelectAllFilesWithQuery(args[1])
		}
	case keyPDFBackup:
		quac.ExportToPDF()
	case keyStats:
		quac.GetStats()
	case keyForceSplit:
		EnsureLen(args, 2)
		MultiOpen(args[1], true) // quick entry force split view
	case keyOpenWorking:
		OpenWorking()
	case keySaveWorking:
		SaveWorking()
	default:
		if len(args) == 1 { // quick query
			//MultiOpen(args[0], false)
			ListSelectAllFilesWithQuery(args[0])
		} else if len(args) >= 2 { // quick entry
			QuickEntry(args[0], strings.Join(args[1:], " "))
		}
	}
}

func EnsureLen(args []string, enLen int) {
	if len(args) < enLen {
		log.Fatalf("expected at least %v args", enLen)
	}
}
