package formsmigration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Creates source_documents.json file
func DoSourceMongoStuff(sourceEnv, sourceTenantID, fileName, sourceDB, sourceUserName, sourceDBPassword *string) {
	sourceConfig := GetConfig(*sourceEnv)
	sourceClient, _ := ConnectToMongoDB(*sourceEnv, sourceConfig)
	// Gets form group names for getting the things from mongo
	formGroups, err := ReadFormGroupNames(*fileName)
	if err != nil {
		log.Fatalf("Error reading form group names: %v", err)
	}
	// exports the source stuff as JSON, saves to file. Saves as source_documents.json
	//docs, err := ExportDocumentsAsJSON(sourceClient, sourceConfig, formGroups, *sourceTenantID)

	_, err = exportAllDocumentsAsJSON(sourceClient, sourceConfig, formGroups, *sourceTenantID)
	if err != nil {
		log.Fatalf("Error exporting documents: %v", err)
	}

}

// Updates source document with target _ids and saves as a new target_documents.json file
// Imports source_documents.json, exports sqlids.json
func DoTargetMongoStuff(targetEnv, targetTenantID, fileName, targetUserName, targetDBPassword, sourceTenantID *string) {

	targetConfig := GetConfig(*targetEnv)
	targetClient, _ := ConnectToMongoDB(*targetEnv, targetConfig)

	updatedResults, err := UpdateDocumentsWithTargetIDs(targetClient, targetConfig, *targetTenantID)
	if err != nil {
		log.Fatalf("Failed to update documents in target environment: %v", err)
	}

	if err := DeleteExistingDocuments(targetClient, targetConfig, updatedResults); err != nil {
		log.Fatalf("Failed to delete existing documents: %v", err)
	}
	Debugtown(fmt.Sprintf("DoMongoStuffTarget: updatedResults: %+v", updatedResults))
	if err := InsertDocumentsIntoTarget(targetClient, targetConfig, updatedResults); err != nil {
		log.Fatalf("Failed to insert documents into target: %v", err)
	}

	file, err := os.OpenFile("SQLids.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Unable to create SQLids.json!")
	}
	defer file.Close()

	jsonData, err := json.Marshal(updatedResults)
	if err != nil {
		log.Fatalf("Failed to marshal json data: %v", err)
	}

	_, lolerr := file.Write(jsonData)
	if lolerr != nil {
		log.Fatalf("Could'nt write to SQLids.json: %v", err)
	}
}

// Gets SQL ids and adds them to a SQL query to get the cursor data field
// Imports SQL Ids as SQLids.json
// Exports cursordata.json file
func DoSourceSQLStuff(sourceEnv, sourceTenantID, fileName, sourceDB, sourceUserName, sourceDBPassword *string) {
	// Run the SQL selection script

	importData, err := os.ReadFile("target_documents.json")
	if err != nil {
		log.Fatal("Can't read target_documents.json!")
	}
	resultids := []MongoJsonWithID{}

	err = json.Unmarshal(importData, &resultids)
	if err != nil {
		log.Fatal("Can't unmarshal results into resultids!")
	}
	//import SQLids.json and unmarshal it into the a new jsonids struct. modify the runselection to use the struct results after unmarshall,
	runMeToo := RunSelectionScript(resultids, *sourceTenantID, *sourceEnv, *sourceDB, *sourceUserName, *sourceDBPassword)

	// marshall runmetoo into a json file and save it, then we need to import it into the do targetsql stuff

	cursorData, err := json.Marshal(runMeToo)
	if err != nil {
		log.Fatalf("Failed to marshal cursordata data: %v", err)
	}

	if err := os.WriteFile("cursordata.json", cursorData, 0644); err != nil {
		log.Fatalf("Failed to write cursordata.json file: %v", err)
	}
}

func DoTargetSQLStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword *string) {

	importData, err := os.ReadFile("cursordata.json")
	if err != nil {
		log.Fatal("Can't read cursordata.json!")
	}

	var getCursorData []LayoutSelectionConfigData
	err = json.Unmarshal(importData, &getCursorData)
	if err != nil {
		log.Fatal("Can't unmarshal results into cursordata!")
	}

	err = RunMigrationScript(*targetEnv, *targetTenantID, *targetUserName, *targetDBPassword, getCursorData)
	if err != nil {
		log.Fatalf("Failed to run migration script: %v", err)
	}
}

func Debugtown(debugstring string) {
	cYellow := "\033[33m"
	reset := "\033[0m"
	// if *debug
	if true {
		fmt.Printf("%sDebugtown: %s%s\n", cYellow, debugstring, reset)
	}
}

func Debugging(debugstring string) {
	cPurple := "\033[35m"
	reset := "\033[0m"

	if true {
		fmt.Printf("%sDebugging: %s%s\n", cPurple, debugstring, reset)
	}
}
