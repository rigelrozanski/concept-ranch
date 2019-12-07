package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rigelrozanski/wt/lib"
)

//keywords used throughout wt
const (
	keyView       = "cat"
	keyNew        = "new"
	keyCopy       = "cp"
	keyRemove     = "rm"
	keyTags       = "tags"
	keyKillTag    = "kill-tag"
	keyAddTag     = "add-tag"
	keyRenameTag  = "rename-tag"
	keyDestroyTag = "destroy-tag"
	keyRecover    = "recover"
	keyEmptyTrash = "empty-trash"
	keyListTrash  = "ls-trash"
	keyList       = "ls"
	keyLog        = "log"
	keyPush       = "push"

	keyHelp1 = "--help"
	keyHelp2 = "-h"

	logWT       = "loglog"
	shortcutsWT = "shortcuts"

	// filestructure: 123456_YYYYMMDD_eYYYYMMDD_cYYYYMMDD_tag_tag_tag_tag...
	// 123456 = id
	// YYYYMMDD = creation date
	// eYYYYMMDD = edited date
	// cYYYYMMDD = consumed date

	help = `
/|||||\ |-o-o-~|
Usage: 
here [query] field can either be populated with an id aka a 6 digit number such
as "123456" or a list of tags seperated by commas such as "tag1,tag2,tag3..."
wt [query]                 -> open a vim tab with the contents of the query 
wt cat [query]             -> print wt contents to console with provide tags
wt [tag1,tag2...] [entry]  -> fast entry appended as new line in wt
wt new [tag1,tag2...]      -> create a new wt with the provided tags
wt cp [id]                 -> duplicate a wt at the provided id
wt rm [id]                 -> remove a wt by id (add to the trash)
wt tags [id]               -> list the tags at an id 
wt kill-tag [id]           -> remove a tag from a wt by id
wt add-tag [id]            -> add a tag from a wt by id
wt rename-tag [id]         -> rename all instances of a tag for all wbs
wt destroy-tag [id]        -> remove all instances of a tag for all wbs
wt recover [id]            -> recover a wt from trash
wt empty-trash             -> empty trash
wt ls-trash                -> list trash ids
wt ls                      -> list all the wt tags 
wt log                     -> list the log
wt push [msg]              -> git push the boards directory
notes:
- if a tag is not provided then the the default empty "" tag is used,
  the default board named 'wt' will be used
- special reserved wt names: wt, lsls, loglog
- tags cannot include underscores (_)
`
)

func main() {
	args := os.Args[1:]

	var err error
	switch len(args) {
	case 0:
		err = edit(defaultWT)
	case 1:
		err = handle1Args(args)
	case 2:
		err = handle2Args(args)
	case 3:
		err = handle3Args(args)
	default:
		name := args[0]
		entry := strings.Join(args[1:], " ")
		err = fastEntry(name, entry)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func handle1Args(args []string) (err error) {
	if len(args) != 1 {
		panic("improper args")
	}
	switch args[0] {
	case keyHelp1, keyHelp2:
		fmt.Println(help)
	case keyPush:
		err = push(fmt.Sprintf("%v", time.Now()))
		lib.MustClearWT(logWT)
	case keyView:
		err = view(defaultWT)
	case keyList:
		err = list()
	case keyEmptyTrash:
		err = emptyTrash()
	case keyLog:
		err = listLog()
	case keyStats:
		err = listStats()
	case keyNew, keyRemove:
		fmt.Println("invalid argments, must specify name of board")
	default:
		err = edit(args[0])
	}

	return err
}

// TODO this is spagetti - fix!
func handle2Args(args []string) error {
	if len(args) != 2 {
		panic("improper args")
	}

	//edit/delete/create-new board
	Bview, Bdelete, Brecover, Bnew := false, false, false, false
	noRsrvArgs := 0

	boardArg := -1
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case keyView:
			Bview = true
			noRsrvArgs++
		case keyRemove:
			Bdelete = true
			noRsrvArgs++
		case keyRecover:
			Brecover = true
			noRsrvArgs++
		case keyNew:
			Bnew = true
			noRsrvArgs++
		case keyList, keyEmptyTrash:
			break
		default:
			boardArg = i
		}
	}

	switch {
	case Bnew:
		name := args[boardArg]
		return freshWT(name)
	case Bview:
		return view(args[boardArg])
	case Bdelete:
		name := args[boardArg]
		return remove(name)
	case Brecover:
		name := args[boardArg]
		return recoverWb(name)
	default:
		name := args[0]
		entry := args[1]
		return fastEntry(name, entry)
	}

	return nil
}

func handle3Args(args []string) error {
	if len(args) != 3 {
		panic("improper args")
	}

	switch args[0] {
	case keyCopy:
		return duplicate(args[1], args[2])
	case keyRename:
		return rename(args[1], args[2])
	default:
		name := args[0]
		entry := strings.Join(args[1:], " ")
		return fastEntry(name, entry)
	}
	return nil
}
