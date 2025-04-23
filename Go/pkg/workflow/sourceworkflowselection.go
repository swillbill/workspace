package workflow

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	//this import grabs the connection strings and env vars
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/formsmigration"
)

type WorkflowSubtypeDetails struct {
	WorkflowSubtype string `json:"workflow_subtype"`
	QueueName       string `json:"queue_name"`
}

type CursorDataSection1 struct {
	CursorData_Section1 string `json:"cursor_data_section1"`
}

type CursorDataSection2 struct {
	CursorData_Section2 string `json:"cursor_data_section2"`
}

type CursorDataSection3 struct {
	CursorData_Section3 string `json:"cursor_data_section3"`
}

// reads from Namelist.txt file
func ReadFormGroupNames(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var formGroups []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			formGroups = append(formGroups, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read lines from file: %w", err)
	}

	return strings.Join(formGroups, ","), nil
}

func SourceWorkflowSelecitonCursorData1(sourceEnv, sourceUserName, sourceDBPassword, tenantID, formGroups string) []CursorDataSection1 {

	log.Printf("Getting Cursor Data Section 1 for Workflow Source Selection")

	// Get form group names from the file
	formGroups, _ = ReadFormGroupNames("wflbundleNameList.txt")
	fmt.Printf("WFL Bundle List of Names: %v\n", formGroups)

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

	var useCorrectDB = "WorkflowEngineDB"
	// get sourceEnv and add as a prefix to the DB
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB1: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "WorkflowEngineDB"
			fmt.Printf("Use Correct DB2: %s\n", useCorrectDB)
		}
	}

	sqlScript := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';
		DECLARE @IdList TABLE ( 
			WorkflowList    nvarchar(500), 
			WorkflowType    nvarchar(500), 
			WorkflowSubtype nvarchar(500) 
		);

		INSERT INTO @IdList ( WorkflowList, WorkflowType, WorkflowSubtype )
		VALUES %s 
		;WITH T AS ( SELECT * FROM @IdList )
		SELECT
			'(''' + ISNULL(wfv.code, '') + ''',''' + ISNULL(wfg.code, '') + ''',''' + ISNULL(wfg.LongDescription, '') + ''',''' + ISNULL(wfg.ShortDescription, '') + ''',''' + 
			ISNULL(wfv.LongDescription, '') + ''',''' + ISNULL(wfv.ShortDescription, '') + ''',''' + ISNULL(wfv.Module, '') + ''',''' + ISNULL(wfg.prefixId, '') + ''',''' + 
			ISNULL(CAST(wfv.[ModuleId] AS nvarchar(10)), '') + ''',''' + ISNULL(CAST(wfg.[ModuleId] AS nvarchar(10)), '') + '''),' AS cursordata_section1
		FROM [WorkflowVariants] wfv
		JOIN [WorkflowGroups] wfg ON wfv.WorkflowGroupId = wfg.WorkflowGroupId
		JOIN [WorkflowVariantObjects] wfo ON wfo.WorkflowVariantId = wfv.WorkflowVariantId
		JOIN T ON T.WorkflowType = wfg.code
		WHERE
			wfv.TenantId = wfo.TenantId
			AND wfv.TenantId = @TenantId
			AND wfg.code = T.WorkflowType
			AND wfv.code = T.WorkflowSubtype;`, useCorrectDB, tenantID, formGroups)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	rows, err := db.Query(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("workflowsourcecursordatasection1.json")
	if err != nil {
		log.Fatalf("Error creating workflowsourcecursordatasection1.json file: %v", err)
	}
	defer file.Close()

	var results []CursorDataSection1
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var cursorData string

		if err := rows.Scan(&cursorData); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")

		results = append(results, CursorDataSection1{CursorData_Section1: cursorData})
	}

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
	log.Println("First data results written to workflowsourcecursordatasection1.json")

	return results

}

func SourceWorkflowSelecitonCursorData2(sourceEnv string, sourceUserName string, sourceDBPassword string, tenantID string) []CursorDataSection2 {

	log.Printf("Getting Cursor Data Section 2 for Workflow Source Selection")

	// Get form group names from the file
	formGroups, _ := ReadFormGroupNames("wflbundleNameList.txt")
	fmt.Printf("WFL Bundle List of Names: %v\n", formGroups)
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

	var useCorrectDB = "WorkflowEngineDB"
	// get sourceEnv and add as a prefix to the DB
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB1: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "WorkflowEngineDB"
			fmt.Printf("Use Correct DB2: %s\n", useCorrectDB)
		}
	}

	sqlScript := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';
		DECLARE @IdList TABLE ( 
			WorkflowList    nvarchar(500), 
			WorkflowType    nvarchar(500), 
			WorkflowSubtype nvarchar(500) 
		);

		INSERT INTO @IdList ( WorkflowList, WorkflowType, WorkflowSubtype )
		VALUES %s;
		;WITH T AS ( SELECT * FROM @IdList )
		SELECT
			'(''' + ISNULL(q.Name, '') + ''',''' + ISNULL(q.[LongDescription], '') + ''',''' + ISNULL(q.[ShortDescription], '') + ''',''' + ISNULL(s.code, '') + '''),' AS cursordata_section2
		FROM [Queue] q
		JOIN [QueueStatus] s ON q.[QueueStatusId] = s.[QueueStatusId]
		WHERE q.TenantId = @TenantId
		AND q.name IN (
			SELECT qq.name
			FROM queue qq
			JOIN WorkflowVariantObjects vo ON qq.[QueueId] = vo.[QueueId]
			JOIN [WorkflowVariants] wv ON wv.[WorkflowVariantId] = vo.[WorkflowVariantId]
			JOIN t ON t.WorkflowSubtype = wv.[Code]
			WHERE vo.TenantId = @TenantId
		);`, useCorrectDB, tenantID, formGroups)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	rows, err := db.Query(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("workflowsourcecursordatasection2.json")
	if err != nil {
		log.Fatalf("Error creating workflowsourcecursordatasection2.json file: %v", err)
	}
	defer file.Close()

	var results []CursorDataSection2
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var cursorData string

		if err := rows.Scan(&cursorData); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")

		results = append(results, CursorDataSection2{CursorData_Section2: cursorData})
	}

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
	log.Println("First data results written to workflowsourcecursordatasection2.json")

	return results

}

func SourceWorkflowSelecitonCursorData3(sourceEnv, sourceUserName, sourceDBPassword, tenantID, formGroups string) []CursorDataSection3 {

	log.Printf("Getting Cursor Data Section 3 for Workflow Source Selection")

	// Get form group names from the file
	formGroups, _ = ReadFormGroupNames("wflbundleNameList.txt")
	fmt.Printf("WFL Bundle List of Names: %v\n", formGroups)

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

	var useCorrectDB = "WorkflowEngineDB"
	// get sourceEnv and add as a prefix to the DB
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB1: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "WorkflowEngineDB"
			fmt.Printf("Use Correct DB2: %s\n", useCorrectDB)
		}
	}

	sqlScript := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';
		DECLARE @IdList TABLE ( 
			WorkflowList    nvarchar(500), 
			WorkflowType    nvarchar(500), 
			WorkflowSubtype nvarchar(500) 
		);

		INSERT INTO @IdList ( WorkflowList, WorkflowType, WorkflowSubtype )
		VALUES %s;

		;WITH T AS ( SELECT * FROM @IdList )
		SELECT
			'(''' + ISNULL(o.[VariantObject], '') + ''',''' + ISNULL(q.[Name], '') + ''',''' + ISNULL(s.[WorkflowSchemaType], '') + ''',''' + ISNULL(v.Code, '') + ''',''' + 
			ISNULL(CAST(o.ActiveFlag AS nvarchar(16)), '') + ''',''' + ISNULL(CAST(o.[ModuleId] AS nvarchar(10)), '') + '''),' AS cursordata_section3
		FROM [WorkflowVariantObjects] o
		JOIN [WorkflowVariants] v ON o.[WorkflowVariantId] = v.[WorkflowVariantId]
		JOIN [Queue] q ON q.[QueueId] = o.[QueueId]
		JOIN [WorkflowSchemas] s ON s.[WorkflowSchemaId] = o.[WorkflowSchemaId]
		WHERE o.ActiveFlag = 1
		AND o.TenantId = @TenantId
		AND v.[Code] IN (SELECT t.WorkflowSubtype FROM t)`, useCorrectDB, tenantID, formGroups)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	rows, err := db.Query(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("workflowsourcecursordatasection3.json")
	if err != nil {
		log.Fatalf("Error creating workflowsourcecursordatasection3.json file: %v", err)
	}
	defer file.Close()

	var results []CursorDataSection3
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var cursorData string

		if err := rows.Scan(&cursorData); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")

		results = append(results, CursorDataSection3{CursorData_Section3: cursorData})
	}

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
	log.Println("First data results written to workflowsourcecursordatasection3.json")

	return results

}
