package correspondence

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/formsmigration"
)

type HeaderAndFooter struct {
	ReusableContentTypeDisplayId string `json:"reusable_content_type_display_id"`
	ReusableContentDisplayId     string `json:"reusable_content_display_id"`
	Name                         string `json:"name"`
	Header_footer_cursor         string `json:"header_footer_cursor"`
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

	formResult := strings.Join(formGroups, ",")
	return formResult, nil
}

func SourceCorrespondenceSelecitonSQL(sourceEnv, sourceUserName, sourceDBPassword, tenantID, formGroups string) []HeaderAndFooter {

	log.Printf("Getting Header and Footer Data for Correspondence Source Selection")

	// Get form group names from the file
	headerAndFooterList, _ := ReadFormGroupNames("correspondenceheaderandfooterlist.txt")
	fmt.Printf("Header and Footer List: %v\n", headerAndFooterList)

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

	sqlScript := fmt.Sprintf(`
	USE [%s];
	DECLARE      @TenantId       uniqueidentifier = '%s';
	-- Header and Footer

	DECLARE @HeaderList TABLE ( Headername    nvarchar(200)   );

	INSERT INTO @HeaderList ( Headername  )
	--select Corrname from @CorrtypeList
	values
	%s
	--use cursor data for section 1 and section 2 in header footer cursor 
	;WITH TL AS ( SELECT * FROM @HeaderList)
	SELECT 

		rct.ReusableContentTypeDisplayId
	,      rc.ReusableContentDisplayId
	,      rct.[Name]

	,      '('''
			+ isnull(rct.[Name],'') +''',''' 
			+ isnull(rct.[ReusableContentTypeCategory],'') +''',''' 
			+ isnull([Status],'') +''',''' 
			+ isnull(cast(rc.[Version] as nvarchar(10)),'') +''',''' 
			+ isnull(rc.[Description],'') +''',''' 
			+ isnull(rct.[ReusableContentTypeDisplayId],'') +''',''' 
			+ isnull(rc.[ReusableContentDisplayId],'') +''','''
			+ isnull(cast(rc.[ModuleId] as nvarchar(10)),'') +''',''' 
	--       +isnull(convert(varchar(max), rc.[ContentOpenXml], 1),'') +''',''' +''',''' 	
	+'''),' as Header_footer_cursor

	FROM [dbo].[ReusableContentType] rct
	JOIN [dbo].[ReusableContent]     rc  ON rct.ReusableContentTypeId=rc.ReusableContentTypeId
	join tl                              on  tl.Headername = rct.[Name]
	WHERE rct.TenantId= @TenantID
	and rc.status = 'ACTIVE'`, useCorrectDB, tenantID, headerAndFooterList)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	rows, err := db.Query(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("headerandfooter.json")
	if err != nil {
		log.Fatalf("Error creating headerandfooter.json file: %v", err)
	}
	defer file.Close()

	var results []HeaderAndFooter
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var (
			reusable_content_type_display_id string
			reusable_content_display_id      string
			name                             string
			header_footer_cursor             string
		)

		if err := rows.Scan(&reusable_content_type_display_id, &reusable_content_display_id, &name, &header_footer_cursor); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		header_footer_cursor = strings.TrimSuffix(header_footer_cursor, ",")

		results = append(results, HeaderAndFooter{
			ReusableContentTypeDisplayId: reusable_content_type_display_id,
			ReusableContentDisplayId:     reusable_content_display_id,
			Name:                         name,
			Header_footer_cursor:         header_footer_cursor,
		})
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
	log.Println("First data results written to headerandfooter.json")

	return results

}
