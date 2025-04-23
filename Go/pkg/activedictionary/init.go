package activedictionary

// import forms migration package
import (
	"encoding/json"
	//"fmt"
	"log"
	"os"
	//"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/formsmigration"
)

func DoSourceADSQLStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword *string) {
	// Run the Active Dictionary SQL selection script
	RunSourceActiveDirectorySQL1(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID)
	RunSourceActiveDirectorySQL2(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID)
	RunSourceActiveDirectorySQL3(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID)
	RunSourceActiveDirectorySQL4(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID)
}

func DoTargetADSQLStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword *string) {
	// Run the Active Dictionary SQL selection script
	importData1, err := os.ReadFile("ADsourcecursordata1.json")
	if err != nil {
		log.Fatal("Can't read ADsourcecursordata1.json!")
	}

	var getCursorData1 []SystemFieldDictionary
	err = json.Unmarshal(importData1, &getCursorData1)
	if err != nil {
		log.Fatal("Can't unmarshal SystemFieldDictionary results into cursordata!")
	}

	//formsmigration.Debugtown(fmt.Sprintf("init.go: cursordata: %+v\n", getCursorData1))
	RunTargetActiveDirectorySQL1(*targetEnv, *targetTenantID, *targetUserName, *targetDBPassword, getCursorData1)

	importData2, err := os.ReadFile("ADsourcecursordata2.json")
	if err != nil {
		log.Fatal("Can't read ADsourcecursordata2.json!")
	}

	var getCursorData2 []SystemFieldDictionaryOOTBSection
	err = json.Unmarshal(importData2, &getCursorData2)
	if err != nil {
		log.Fatal("Can't unmarshal SystemFieldDictionaryOOTBSection results into cursordata!")
	}

	RunTargetActiveDirectorySQL2(*targetEnv, *targetTenantID, *targetUserName, *targetDBPassword, getCursorData2)

	importData3, err := os.ReadFile("ADsourcecursordata3.json")
	if err != nil {
		log.Fatal("Can't read ADsourcecursordata3.json!")
	}

	var getCursorData3 []SystemFieldDictionarySection
	err = json.Unmarshal(importData3, &getCursorData2)
	if err != nil {
		log.Fatal("Can't unmarshal SystemFieldDictionary results into cursordata!")
	}

	RunTargetActiveDirectorySQL3(*targetEnv, *targetTenantID, *targetUserName, *targetDBPassword, getCursorData3)

	importData4, err := os.ReadFile("ADsourcecursordata4.json")
	if err != nil {
		log.Fatal("Can't read ADsourcecursordata4.json!")
	}

	var getCursorData4 []SystemFieldDictionaryOOTBField
	err = json.Unmarshal(importData4, &getCursorData4)
	if err != nil {
		log.Fatal("Can't unmarshal SystemFieldDictionaryOOTBField results into cursordata!")
	}

	RunTargetActiveDirectorySQL4(*targetEnv, *targetTenantID, *targetUserName, *targetDBPassword, getCursorData4)

}
