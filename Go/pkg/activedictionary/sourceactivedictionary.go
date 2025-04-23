package activedictionary

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	// import forms migration package
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/formsmigration"
)

// type SystemFieldDictionary struct {
// 	CursorData string `json:"cursordata"`
// }

type SystemFieldDictionary struct {
	FixedId                       int    `json:"fixed_id"`
	SystemFieldDictionaryId       int    `json:"system_field_dictionary_id"`
	ModuleId                      int    `json:"module_id"`
	Version                       int    `json:"version"`
	Status                        string `json:"status"`
	LogixInputSchemaJson          string `json:"logix_input_schema_json"`
	SystemFieldDictionaryUniqueId string `json:"system_field_dictionary_unique_id"`
	CursorData                    string `json:"cursordata"`
}

type SystemFieldDictionaryOOTBSection struct {
	SectionJson                      string `json:"section_json"`
	ContextName                      string `json:"context_name"`
	SectionName                      string `json:"section_name"`
	ModuleId                         int    `json:"module_id"`
	SystemFieldDictionaryOOTBSection string `json:"system_field_dictionary_ootb_section"`
	SystemFieldDictionaryId          int    `json:"system_field_dictionary_id"`
	CursorData                       string `json:"cursordata"`
}

type SystemFieldDictionarySection struct {
	SectionJson             string `json:"section_json"`
	ContextName             string `json:"context_name"`
	SystemFieldDictionaryId int    `json:"system_field_dictionary_id"`
	ModuleId                int    `json:"module_id"`
	CursorData              string `json:"cursordata"`
}

type SystemFieldDictionaryOOTBField struct {
	FieldJson   string `json:"field_json"`
	SectionName string `json:"section_name"`
	FieldName   string `json:"field_name"`
	ModuleId    int    `json:"module_id"`
	CursorData  string `json:"cursordata"`
}

func RunSourceActiveDirectorySQL1(sourceEnv string, sourceUserName string, sourceDBPassword string, tenantID string) []SystemFieldDictionary {
	// getconfiguration based on the environment, source username, and source password flags
	sqlconfig, ok := formsmigration.GetSelectionSQLConfig(sourceEnv, sourceUserName, sourceDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", sqlconfig)
	}

	fmt.Printf("Get DB ENV %v\n", sqlconfig.Database)

	// connection string from the config for SQL Server connection OPEN
	db, err := sql.Open("sqlserver", sqlconfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	fmt.Println("SQL PING!")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	var useCorrectDB = "NgCorrespondence"
	// get sourceEnv and add as a prefix to the DB
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB1: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use Correct DB2: %s\n", useCorrectDB)
		}
	}

	sqlScript := fmt.Sprintf(`USE [%s];
		SELECT [FixedId]                                                                 
,      [Status]                                                                  
,      [Version]                                                                 
,      [LogixInputSchemaJson]                                                    
,      SystemFieldDictionaryUniqueId                                             
,      SystemFieldDictionaryId                                                   
,      ModuleId                                                                  
,      '('''+isnull(cast([FixedId] AS nvarchar(10)),'') +''','''+ isnull([Status],'') +''','''+ isnull(cast([Version] AS nvarchar(10)),'') +''',''' + isnull(cast([ModuleId] AS nvarchar(10)),'') 
		+''','''+ isnull(replace([LogixInputSchemaJson],'''',''''''),'')
		+''','''+ isnull(cast([SystemFieldDictionaryUniqueId] AS nvarchar(50)),'')+'''),' cursordata
FROM [SystemFieldDictionary]
WHERE [TenantId] = '%s'
	AND status IN ( 'ACTIVE') 
	--AND status IN ( 'ACTIVE', 'DRAFT')`, useCorrectDB, tenantID)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	rows, err := db.Query(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	// open a file to write the JSON data
	file, err := os.Create("ADsourcecursordata1.json")
	if err != nil {
		log.Fatalf("Error creating ADsourcecursordata1.json file: %v", err)
	}
	defer file.Close()

	var results []SystemFieldDictionary
	var rownumbah int = 0

	for rows.Next() {
		rownumbah++
		var (
			fixedId                       int
			systemFieldDictionaryId       int
			moduleId                      int
			version                       int
			status                        string
			logixInputSchemaJson          string
			systemFieldDictionaryUniqueId string
			cursorData                    string
		)
		//err := rows.Scan(&cursorData)
		if err := rows.Scan(&fixedId, &status, &version, &logixInputSchemaJson, &systemFieldDictionaryUniqueId, &systemFieldDictionaryId, &moduleId, &cursorData); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")

		results = append(results, SystemFieldDictionary{
			FixedId:                       fixedId,
			Status:                        status,
			Version:                       version,
			LogixInputSchemaJson:          logixInputSchemaJson,
			SystemFieldDictionaryUniqueId: systemFieldDictionaryUniqueId,
			SystemFieldDictionaryId:       systemFieldDictionaryId,
			ModuleId:                      moduleId,
			CursorData:                    cursorData,
		})
	}
	fmt.Printf("Number of rows: %d\n", rownumbah)
	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred while iterating rows: %v", err)
	}

	// marshal the results into JSON format
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to JSON file: %v", err)
	}
	log.Println("First data results written to ADsourcecursordata1.json")

	return results
}

func RunSourceActiveDirectorySQL2(sourceEnv string, sourceUserName string, sourceDBPassword string, tenantID string) []SystemFieldDictionaryOOTBSection {
	log.Println("Starting second active directory source sql query")

	// getconfiguration based on the environment, source username, and source password flags
	sqlconfig, ok := formsmigration.GetSelectionSQLConfig(sourceEnv, sourceUserName, sourceDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", sqlconfig)
	}

	fmt.Printf("Get DB ENV %v\n", sqlconfig.Database)

	// connection string from the config for SQL Server connection OPEN
	db, err := sql.Open("sqlserver", sqlconfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	fmt.Println("SQL PING!")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	// get sourceEnv and add as a prefix to the DB
	var useCorrectDB = "NgCorrespondence"
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use REF DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use DB: %s\n", useCorrectDB)
		}
	}

	sqlScript2 := fmt.Sprintf(
		`USE [%s];
		SELECT o.[SectionJson]
		,      o.[ContextName]
		,      o.[SectionName]
		,      o.[ModuleId]
		,      '('''+Isnull(o.[SectionJson],'') +''','''+ Isnull(o.[ContextName],'') +''','''+ Isnull(o.[SectionName],'')+''','''+ Isnull(Cast(o.[ModuleId] AS NVARCHAR(10)),'') + '''),'cursordata
		FROM [SystemFieldDictionaryOOTBSection] o
		JOIN [SystemFieldDictionary]            d ON d.[SystemFieldDictionaryId] = o.[SystemFieldDictionaryId]
		WHERE o.[TenantId] = '%s' 
			AND d.[Status] = 'ACTIVE'`, useCorrectDB, tenantID)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript2))

	rows, err := db.Query(sqlScript2)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	// open a file to write the JSON data
	file, err := os.Create("ADsourcecursordata2.json")
	if err != nil {
		log.Fatalf("Error creating ADsourcecursordata2.json file: %v", err)
	}
	defer file.Close()

	// create a slice to hold the results
	var results []SystemFieldDictionaryOOTBSection
	var rownumbah int = 0

	for rows.Next() {
		rownumbah++
		var (
			sectionJson, contextName, sectionName, cursorData string
			moduleId                                          int
		)
		if err := rows.Scan(&sectionJson, &contextName, &sectionName, &moduleId, &cursorData); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")
		results = append(results, SystemFieldDictionaryOOTBSection{
			SectionJson: sectionJson,
			ContextName: contextName,
			SectionName: sectionName,
			ModuleId:    moduleId,
			CursorData:  cursorData,
		})
	}
	fmt.Printf("Number of rows: %d\n", rownumbah)
	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred while iterating rows: %v", err)
	}

	// marshal the results into JSON format
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to JSON file: %v", err)
	}

	log.Println("Second data results written to ADsourcecursordata2.json")
	return results

}

func RunSourceActiveDirectorySQL3(sourceEnv string, sourceUserName string, sourceDBPassword string, tenantID string) []SystemFieldDictionarySection {
	log.Println("Starting third active directory source sql query")

	// getconfiguration based on the environment, source username, and source password flags
	sqlconfig, ok := formsmigration.GetSelectionSQLConfig(sourceEnv, sourceUserName, sourceDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", sqlconfig)
	}

	// connection string from the config for SQL Server connection OPEN
	db, err := sql.Open("sqlserver", sqlconfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	fmt.Println("SQL PING!")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	// get sourceEnv and add as a prefix to the DB
	var useCorrectDB = "NgCorrespondence"
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use REF DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use DB: %s\n", useCorrectDB)
		}
	}

	sqlScript3 := fmt.Sprintf(
		`USE [%s];
		SELECT s.[SectionJson]                                                                                                                          
		,      s.[ContextName]                                                                                                                          
		,      d. [SystemFieldDictionaryId]                                                                                                             
		,      s. ModuleId                                                                                                                              
		,      '('''+isnull(s.[SectionJson],'') +''','''+ isnull(s.[ContextName],'')+''',''' + isnull(cast(s.[ModuleId] AS nvarchar(10)),'') +'' +'''),' cursordata
		
		FROM [SystemFieldDictionarySection] s
		JOIN [SystemFieldDictionary]        d ON d.[SystemFieldDictionaryId] = s.[SystemFieldDictionaryId]
		WHERE d.[TenantId] = '%s'
			AND d.[Status] = 'ACTIVE'`, useCorrectDB, tenantID)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript3))

	rows, err := db.Query(sqlScript3)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	// open a file to write the JSON data
	file, err := os.Create("ADsourcecursordata3.json")
	if err != nil {
		log.Fatalf("Error creating ADsourcecursordata3.json file: %v", err)
	}
	defer file.Close()

	var results []SystemFieldDictionarySection
	var rownumbah int = 0

	for rows.Next() {
		rownumbah++
		var (
			sectionJson, contextName, cursorData string
			systemFieldDictionaryId, moduleId    int
		)
		if err := rows.Scan(&sectionJson, &contextName, &systemFieldDictionaryId, &moduleId, &cursorData); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")
		results = append(results, SystemFieldDictionarySection{
			SectionJson:             sectionJson,
			ContextName:             contextName,
			SystemFieldDictionaryId: systemFieldDictionaryId,
			ModuleId:                moduleId,
			CursorData:              cursorData,
		})
	}
	fmt.Printf("Number of rows: %d\n", rownumbah)
	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred while iterating rows: %v", err)
	}

	// marshal the results into JSON format
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to JSON file: %v", err)
	}

	log.Println("Second data results written to ADsourcecursordata3.json")
	return results

}

func RunSourceActiveDirectorySQL4(sourceEnv string, sourceUserName string, sourceDBPassword string, tenantID string) []SystemFieldDictionaryOOTBField {
	log.Println("Starting fourth active directory source sql query")

	// getconfiguration based on the environment, source username, and source password flags
	sqlconfig, ok := formsmigration.GetSelectionSQLConfig(sourceEnv, sourceUserName, sourceDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", sqlconfig)
	}

	fmt.Printf("Get DB ENV %v\n", sqlconfig.Database)

	// connection string from the config for SQL Server connection OPEN
	db, err := sql.Open("sqlserver", sqlconfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	fmt.Println("SQL PING!")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	// get sourceEnv and add as a prefix to the DB
	var useCorrectDB = "NgCorrespondence"
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use REF DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use DB: %s\n", useCorrectDB)
		}
	}

	sqlScript4 := fmt.Sprintf(
		`USE [%s];
		SELECT f.[FieldJson]                                                                                                                                                         
		,      f.[SectionName]                                                                                                                                                       
		,      f.[FieldName]                                                                                                                                                         
		,      f.ModuleId                                                                                                                                                            
		,      '('''+isnull(f.[FieldJson],'') +''','''+ isnull(f.[SectionName],'') +''','''+ isnull(f.[FieldName],'') +''',''' + isnull(cast(f.[ModuleId] AS nvarchar(10)),'')+'''),' cursordata
		FROM [SystemFieldDictionaryOOTBField]   f 
		JOIN [SystemFieldDictionaryOOTBSection] os ON os.[SystemFieldDictionaryOOTBSectionId] = f.[SystemFieldDictionaryOOTBSectionId]
		JOIN [SystemFieldDictionary]            s  ON s.[SystemFieldDictionaryId] = os.[SystemFieldDictionaryId]
		WHERE f.[TenantId] = '%s'
			AND s.[Status] = 'ACTIVE'`, useCorrectDB, tenantID)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript4))

	rows, err := db.Query(sqlScript4)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	// open a file to write the JSON data
	file, err := os.Create("ADsourcecursordata4.json")
	if err != nil {
		log.Fatalf("Error creating ADsourcecursordata4.json file: %v", err)
	}
	defer file.Close()

	var results []SystemFieldDictionaryOOTBField
	var rownumbah int = 0

	for rows.Next() {
		rownumbah++
		var (
			fieldJson, sectionName, fieldName, cursorData string
			moduleId                                      int
		)
		if err := rows.Scan(&fieldJson, &sectionName, &fieldName, &moduleId, &cursorData); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")
		results = append(results, SystemFieldDictionaryOOTBField{
			FieldJson:   fieldJson,
			SectionName: sectionName,
			FieldName:   fieldName,
			ModuleId:    moduleId,
			CursorData:  cursorData,
		})
	}

	fmt.Printf("Number of rows: %d\n", rownumbah)
	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred while iterating rows: %v", err)
	}

	// marshal the results into JSON format
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to JSON file: %v", err)
	}

	log.Println("Second data results written to ADsourcecursordata4.json")
	return results
}
