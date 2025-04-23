package correspondence

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func DoSourceCorrespondenceStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword *string) {
	var formResult string
	SourceCorrespondenceSelecitonSQL(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
}

func DoSourceCorrespondenceTemplateStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword *string) {
	var formResult string
	CorrespondenceTemplate1(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
	CorrespondenceTemplate2(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
	CorrespondenceTemplate3(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
	CorrespondenceTemplate4(*sourceEnv, *sourceUserName, *sourceDBPassword, *sourceTenantID, formResult)
}

func DoTargetCorrespondenceCommercialStuff(targetEnv, sourceEnv, sourceTenantID, targetTenantID, targetUserName, targetDBPassword *string) {

	var gov = (*targetEnv == "stg" || *targetEnv == "vaprod" || *targetEnv == "txprod" ||
		*targetEnv == "stl" || *targetEnv == "stp" || *targetEnv == "devgov" || *targetEnv == "qagov")
	fmt.Printf("SOURCE ENV: %v\n", *targetEnv)

	if *targetEnv != "" && gov {
		log.Print("Executing Target Correspondence GOV SQL Selection Query!")

		// Run the Active Dictionary SQL selection script
		importData, err := os.ReadFile("headerandfooter.json")
		if err != nil {
			log.Fatal("Can't read headerandfooter.json!")
		}

		var getHeaderandFooter []HeaderAndFooter
		err = json.Unmarshal(importData, &getHeaderandFooter)
		if err != nil {
			log.Fatal("Can't unmarshal Header And Footer results!")
		}

		TargetCorrespondenceMigrationGOVSQL(*targetEnv, *sourceTenantID, *targetTenantID, *targetUserName, *targetDBPassword, getHeaderandFooter)
	}

	// Run the Corr SQL selection script
	importData, err := os.ReadFile("headerandfooter.json")
	if err != nil {
		log.Fatal("Can't read headerandfooter.json!")
	}

	var getHeaderandFooter []HeaderAndFooter
	err = json.Unmarshal(importData, &getHeaderandFooter)
	if err != nil {
		log.Fatal("Can't unmarshal Header And Footer results!")
	}

	TargetCorrespondenceMigrationCommercialSQL(*targetEnv, *sourceTenantID, *targetTenantID, *targetUserName, *targetDBPassword, getHeaderandFooter)

}

func DoTargetCorrespondenceTemplateCommercialStuff(targetEnv, sourceEnv, sourceTenantID, targetTenantID, targetUserName, targetDBPassword *string) {
	// Define GOV environments
	govEnvs := map[string]bool{
		"stg": true, "vaprod": true, "txprod": true,
		"stl": true, "stp": true, "devgov": true, "qagov": true,
	}

	fmt.Printf("SOURCE ENV: %v\n", *targetEnv)

	// Check if it's a GOV environment
	if *targetEnv != "" && govEnvs[*targetEnv] {
		log.Print("Executing Target Correspondence GOV SQL Selection Query!")

		// Define data variables for unmarshalling
		var getCursorData1 []CursorData1
		var getCursorData2 []CursorData2
		var getCursorData3 []CursorData3
		var getCursorData4 []CursorData4

		// Define file paths to read from
		files := []string{
			"corrsourcecursordatasection1.json",
			"corrsourcecursordatasection2.json",
			"corrsourcecursordatasection3.json",
			"corrsourcecursordatasection4.json",
		}

		// Define corresponding data variables
		dataVars := []interface{}{&getCursorData1, &getCursorData2, &getCursorData3, &getCursorData4}

		// Iterate over each file and unmarshal its contents
		for i, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				log.Fatalf("Error reading file %s: %v", file, err)
			}

			// Unmarshal data into the corresponding structure
			err = json.Unmarshal(data, dataVars[i])
			if err != nil {
				log.Fatalf("Error unmarshaling data from file %s: %v", file, err)
			}
		}

		// Pass the unmarshalled data to the target migration function
		TargetCorrespondenceTemplateMigrationGOVSQL(
			*targetEnv, *sourceTenantID, *targetTenantID, *targetUserName, *targetDBPassword,
			getCursorData1, getCursorData2, getCursorData3, getCursorData4,
		)
	} else {
		// Similar logic for non-GOV environments
		var getCursorData1 []CursorData1
		var getCursorData2 []CursorData2
		var getCursorData3 []CursorData3
		var getCursorData4 []CursorData4

		// Define file paths to read from
		files := []string{
			"corrsourcecursordatasection1.json",
			"corrsourcecursordatasection2.json",
			"corrsourcecursordatasection3.json",
			"corrsourcecursordatasection4.json",
		}

		// Define corresponding data variables
		dataVars := []interface{}{&getCursorData1, &getCursorData2, &getCursorData3, &getCursorData4}

		// Iterate over each file and unmarshal its contents
		for i, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				log.Fatalf("Error reading file %s: %v", file, err)
			}

			// Unmarshal data into the corresponding structure
			err = json.Unmarshal(data, dataVars[i])
			if err != nil {
				log.Fatalf("Error unmarshaling data from file %s: %v", file, err)
			}
		}

		// Call the template migration function with the unmarshalled data
		TargetCorrespondenceMigrationCommercialSQLTemplate(
			*targetEnv, *sourceTenantID, *targetTenantID, *targetUserName, *targetDBPassword,
			getCursorData1, getCursorData2, getCursorData3, getCursorData4,
		)
	}
}
