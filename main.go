package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rigelrozanski/wb/lib"
)

//keywords used throughout wb
const (
	keyNew        = "new"
	keyCopy       = "cp"
	keyRename     = "rename"
	keyView       = "cat"
	keyRemove     = "rm"
	keyRecover    = "recover"
	keyEmptyTrash = "empty-trash"
	keyList       = "ls"
	keyLog        = "log"
	keyStats      = "stats"
	keyPush       = "push"

	keyHelp1 = "--help"
	keyHelp2 = "-h"

	defaultWB   = "wb"
	lsWB        = "lsls"
	logWB       = "loglog"
	shortcutsWB = "shortcuts"

	help = `
/|||||\ |-o-o-~|

Usage: 

wb [name]               -> vim into a wb
wb [name] [entry]       -> fast entry appended as new line in wb
wb new [name]           -> create a new wb
wb cp [copy] [name]     -> duplicate a wb
wb rename [old] [name]  -> rename a wb
wb cat [name]           -> print wb contents to console
wb rm [name]            -> remove a wb (add to trash)
wb recover [name]       -> remove a wb from trash
wb empty-trash          -> empty trash
wb ls                   -> list all the wb in console
wb log                  -> list the log
wb stats                -> list git statistics per wb
wb push [msg]           -> git push the boards directory

notes:
- if the [name] is not provided, 
  the default board named 'wb' will be used
- special reserved wb names: wb, lsls, loglog

`
)

func main() {
	args := os.Args[1:]

	var err error
	switch len(args) {
	case 0:
		err = edit(defaultWB)
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
		lib.MustClearWB(logWB)
	case keyView:
		err = view(defaultWB)
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
		return freshWB(name)
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
