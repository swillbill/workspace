package platformconfiguration

import (
	"encoding/json"
	"fmt"
	"os"
)

func DoSourcePLCStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword *string) {
	var formResult string
	SourcePLCselectionCursorData1(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
}

func DoTargetPLCStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword *string) {

	var (
		getCursorData1 []CursorData1
		getCursorData2 []CursorData2
	)

	importData1, err := os.ReadFile("plcsourcecursordatasection1.json")
	if err != nil {
		fmt.Printf("Error cant read workflow cursor data file 1. %e\n", err)
	}

	importData2, err := os.ReadFile("plcsourcecursordatasection2.json")
	if err != nil {
		fmt.Printf("Error cant read workflow cursor data file 2. %e\n", err)
	}

	err = json.Unmarshal(importData1, &getCursorData1)
	if err != nil {
		fmt.Printf("Error cant unmarshal workflow CursorDataSection1 data files. %e\n", err)
	}

	err = json.Unmarshal(importData2, &getCursorData2)
	if err != nil {
		fmt.Printf("Error cant unmarshal workflow CursorDataSection2 data files. %e\n", err)
	}

	RunTargetPLCSQL(*targetEnv, *targetTenantID, *targetUserName, *targetDBPassword, getCursorData1, getCursorData2)
}
