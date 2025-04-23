package workflow

import (
	"encoding/json"
	"fmt"
	"os"
)

func DoSourceWorkflowSQLStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword *string) {
	var formResult string
	_, err := ReadFormGroupNames("wflbundleNameList.txt")
	if err != nil {
		fmt.Printf("Error cant read bundle of names. %e\n", err)
	}

	SourceWorkflowSelecitonCursorData1(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
	SourceWorkflowSelecitonCursorData2(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID)
	SourceWorkflowSelecitonCursorData3(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
}

func DoTargetWorkflowSQLStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword *string) {

	var (
		getCursorData1 []CursorDataSection1
		getCursorData2 []CursorDataSection2
		getCursorData3 []CursorDataSection3
	)

	importData1, err := os.ReadFile("workflowsourcecursordatasection1.json")
	if err != nil {
		fmt.Printf("Error cant read workflow cursor data file 1. %e\n", err)
	}

	importData2, err := os.ReadFile("workflowsourcecursordatasection2.json")
	if err != nil {
		fmt.Printf("Error cant read workflow cursor data file 2. %e\n", err)
	}

	importData3, err := os.ReadFile("workflowsourcecursordatasection3.json")
	if err != nil {
		fmt.Printf("Error cant read workflow cursor data file 3. %e\n", err)
	}

	err = json.Unmarshal(importData1, &getCursorData1)
	if err != nil {
		fmt.Printf("Error cant unmarshal workflow CursorDataSection1 data files. %e\n", err)
	}

	err = json.Unmarshal(importData2, &getCursorData2)
	if err != nil {
		fmt.Printf("Error cant unmarshal workflow CursorDataSection2 data files. %e\n", err)
	}

	err = json.Unmarshal(importData3, &getCursorData3)
	if err != nil {
		fmt.Printf("Error cant unmarshal workflow CursorDataSection3 data files. %e\n", err)
	}
	RunTargetWorkflowSQL(*targetEnv, *targetTenantID, *targetUserName, *targetDBPassword, getCursorData1, getCursorData2, getCursorData3)

}
