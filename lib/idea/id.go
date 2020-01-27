package idea

import (
	"fmt"
	"strconv"

	cmn "github.com/rigelrozanski/common"
)

func GetNextID() uint32 {
	lines, err := cmn.ReadLines(ConfigFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config, error: %v", err))
	}
	count, err := strconv.Atoi(lines[0])
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
