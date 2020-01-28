package idea

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	cmn "github.com/rigelrozanski/common"
)

func GetNextID() uint32 {
	lines, err := cmn.ReadLines(ConfigFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config, error: %v", err))
	}
	count, err := ParseID(lines[0])
	if err != nil {
		panic(fmt.Sprintf("error reading id_counter, error: %v", err))
	}
	return uint32(count + 1)
}

func IncrementID() {
	err := cmn.WriteLines([]string{IdStr(GetNextID())}, ConfigFile)
	if err != nil {
		panic(err)
	}
}

// parse the last id, if no error add to the last ids file
func ParseID(idStr string) (uint32, error) {
	return ParseIDOp(idStr, true)
}

func ParseIDNoLogLast(idStr string) (uint32, error) {
	return ParseIDOp(idStr, false)
}

// parse the last id, if no error add to the last ids file
func ParseIDOp(idStr string, logLast bool) (uint32, error) {

	if len(strings.Split(idStr, ",")) > 1 {
		errors.New("not an id, contains commas")
	}

	// read in the lastIDs
	lastIDs, err := cmn.ReadLines(LastIdFile)
	if err != nil {
		return 0, err
	}

	// get the idea from the last list
	var parsedID uint32
	if strings.HasPrefix(idStr, Last) {
		remainder := strings.TrimPrefix(idStr, Last)
		switch len(remainder) {
		case 0:
			id, err := strconv.Atoi(lastIDs[0])
			if err != nil {
				return 0, err
			}
			parsedID = uint32(id)
		case 1:
			lastNo, err := strconv.Atoi(remainder)
			if err != nil {
				return 0, err
			}
			if lastNo > len(lastIDs) {
				return 0, errors.New("insufficient last lines saved")
			}

			id, err := strconv.Atoi(lastIDs[lastNo-1])
			if err != nil {
				return 0, err
			}
			parsedID = uint32(id)
		default:
			return 0, errors.New("can only return up to the 9th previous last id")
		}
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return 0, err
		}
		parsedID = uint32(id)
	}

	if logLast {
		// Prepend retrieved id to the "last" list
		// and trim the list to the appropriate length
		parsedIDStr := strconv.Itoa(int(parsedID))
		lastIDs = append([]string{parsedIDStr}, lastIDs...)
		if len(lastIDs) > 9 {
			lastIDs = lastIDs[:9]
		}
		err = cmn.WriteLines(lastIDs, LastIdFile)
		if err != nil {
			return 0, err
		}
	}

	return parsedID, nil
}