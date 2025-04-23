package platformconfiguration

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

type CursorData1 struct {
	PlatformConfigurationId  string `json:"PlatformConfigurationId"`
	Name                     string `json:"Name"`
	ConfigurationModule      int    `json:"ConfigurationModule"`
	ConfigurationDomain      int    `json:"ConfigurationDomain"`
	ConfigurationType        string `json:"ConfigurationType"`
	ConfigurationName        string `json:"ConfigurationName"`
	StateOf                  string `json:"StateOf"`
	ConfigurationDescription string `json:"ConfigurationDescription"`
	DisplayName              string `json:"displayname"`
	IsSchema                 bool   `json:"IsSchema"`
	Version                  int    `json:"Version"`
	ConfigurationInfo        string `json:"ConfigurationInfo"`
	IsOOTBEditable           int    `json:"IsOOTBEditable"`
	CursorData               string `json:"CursorData"`
}

type CursorData2 struct {
	CursorData string `json:"CursorData"`
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
	fmt.Printf("Bundle List of Names: %v\n", formResult)
	return formResult, nil
}

func SourcePLCselectionCursorData1(sourceEnv, sourceUserName, sourceDBPassword, sourceTenantID, formResult string) ([]CursorData1, []CursorData2) {

	log.Printf("Getting Cursor Data Section 1 for Workflow Source Selection")

	// Get form group names from the file
	formResult, _ = ReadFormGroupNames("plcbundleNameList.txt")

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

	var useCorrectDB = "NgPlatform"
	// get sourceEnv and add as a prefix to the DB
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB1: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgPlatform"
			fmt.Printf("Use Correct DB2: %s\n", useCorrectDB)
		}
	}

	sqlScript1 := fmt.Sprintf(`
	USE [%s];
	DECLARE      @TenantId       uniqueidentifier = '%s';
	WITH T AS 
	(SELECT sq1.*
	FROM (
	VALUES 

	--Section 1: Data for the selection
	--------------------------
	%s
	--------------------------
	--End of Section 1: Data for the selection

	) AS sq1 (Module, Domain, PCtype, PCname, IsSchema)
	)
	select distinct pc.PlatformConfigurationId, pc.Name,pc.ConfigurationModule, pc.ConfigurationDomain, pc.ConfigurationType, pc.ConfigurationName, 
	pc.StateOf,  pc.ConfigurationDescription,pc.displayname ,pc.IsSchema, 1/*pc.Version*/ [version],pc.ConfigurationInfo, 1 /*pc.IsOOTBEditable*/IsOOTBEditable,

	'('''+cast(pc.PlatformConfigurationId as varchar(100))+''','''+pc.Name+''','+cast(pc.ConfigurationModule as varchar(100))+','+cast(pc.ConfigurationDomain as varchar(100))+','''+
	pc.ConfigurationType+''','''+pc.ConfigurationName+''','''+isnull(pc.StateOf,'')+''','''+isnull(pc.ConfigurationDescription,'')+''',
	'''+isnull(pc.displayname,'')+''',
	'+cast(pc.IsSchema as varchar(100))+
	','+cast([version] as varchar(100))+','''+pc.ConfigurationInfo+''',1),' cursordata


	from t join 
	(
	select p.PlatformConfigurationId, g.Name, p.ConfigurationModule, p.ConfigurationDomain, p.ConfigurationType, p.ConfigurationName, 
	p.StateOf, p.ConfigurationDescription, p.IsSchema, p.Version,i.ConfigurationInfo, p.IsOOTBEditable, p.displayname
	from PlatformConfiguration p 
	join PlatformConfigurationInfo i on i.PlatformConfigurationId=p.PlatformConfigurationId
	left join PlatformConfigurationGroup g on g.PlatformConfigurationGroupId = p.PlatformConfigurationGroupId
	where p.TenantId= @TenantId -- CHES CONFIG

	--  'aa106e70-19f7-4269-9eeb-fd6e9bfdc88b' --STL RefConfig Revenue
	) pc on pc.ConfigurationDomain=t.Domain and pc.ConfigurationModule=t.Module and pc.ConfigurationType=t.PCtype 
			and pc.ConfigurationName=t.PCname and pc.IsSchema = t.IsSchema --and pc.name = t.GroupName
	where 1=1
	and (stateof is null or stateof <> 'Deleted' or stateof='')
	order by ConfigurationType, ConfigurationName`, useCorrectDB, sourceTenantID, formResult)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript1))
	rows, err := db.Query(sqlScript1)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("plcsourcecursordatasection1.json")
	if err != nil {
		log.Fatalf("Error creating plcsourcecursordatasection1.json file: %v", err)
	}
	defer file.Close()

	var results1 []CursorData1
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var (
			platformconfiguraitonid  string
			name                     string
			configurationmodule      int
			configurationdomain      int
			configurationtype        string
			configurationname        string
			stateof                  sql.NullString
			configurationdescription string
			displayname              sql.NullString
			isschema                 bool
			version                  int
			configurationinfo        string
			isootbeditable           int
			cursordata               string
		)

		if err := rows.Scan(&platformconfiguraitonid, &name, &configurationmodule, &configurationdomain, &configurationtype, &configurationname, &stateof, &configurationdescription, &displayname, &isschema, &version, &configurationinfo, &isootbeditable, &cursordata); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursordata = strings.TrimSuffix(cursordata, ",")
		results1 = append(results1, CursorData1{
			PlatformConfigurationId:  platformconfiguraitonid,
			Name:                     name,
			ConfigurationModule:      configurationmodule,
			ConfigurationDomain:      configurationdomain,
			ConfigurationType:        configurationtype,
			ConfigurationName:        configurationname,
			StateOf:                  stateof.String,
			ConfigurationDescription: configurationdescription,
			DisplayName:              stateof.String,
			IsSchema:                 isschema,
			Version:                  version,
			ConfigurationInfo:        configurationinfo,
			IsOOTBEditable:           isootbeditable,
			CursorData:               cursordata,
		})
	}

	fmt.Printf("Number of rows: %d\n", rowNumber)
	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred while iterating rows: %v", err)
	}

	// marshal the results into JSON format
	jsonData, err := json.Marshal(results1)
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to JSON file: %v", err)
	}
	log.Println("First data results written to plcsourcecursordatasection1.json")

	sqlScript2 := fmt.Sprintf(`
	USE [%s];
	DECLARE      @TenantId       uniqueidentifier = '%s';

	WITH T AS 
	(SELECT sq1.*
	FROM (
	VALUES 

	--Section 1
	-----------------------------------------------------------------------
	%s
	-----------------------------------------------------------------------
	--End of Section 1

	) AS sq1 (Module, Domain, PCtype, PCname, IsSchema)
	)
	select distinct '('''+pc.name+''','''+pc.Description+'''),' cursordata
	from t join 
	(
	select p.PlatformConfigurationId, g.Name, g.Description, p.ConfigurationModule, p.ConfigurationDomain, p.ConfigurationType, p.ConfigurationName, 
	p.StateOf, p.ConfigurationDescription, p.IsSchema, p.Version,i.ConfigurationInfo, p.IsOOTBEditable
	from PlatformConfiguration p 
	join PlatformConfigurationInfo i on i.PlatformConfigurationId=p.PlatformConfigurationId
	left join PlatformConfigurationGroup g on g.PlatformConfigurationGroupId = p.PlatformConfigurationGroupId
	where p.TenantId= @TenantId
	) pc on pc.ConfigurationDomain=t.Domain and pc.ConfigurationModule=t.Module and pc.ConfigurationType=t.PCtype 
			and pc.ConfigurationName=t.PCname and pc.IsSchema = t.IsSchema --and pc.name = t.GroupName
	where 1=1
	and (stateof is null or stateof <> 'Deleted' or stateof='')`, useCorrectDB, sourceTenantID, formResult)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript2))
	rows, err = db.Query(sqlScript2)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err = os.Create("plcsourcecursordatasection2.json")
	if err != nil {
		log.Fatalf("Error creating plcsourcecursordatasection2.json file: %v", err)
	}
	defer file.Close()

	var results2 []CursorData2
	var rownumbah = 0

	for rows.Next() {
		rownumbah++
		var cursorData string

		if err := rows.Scan(&cursorData); err != nil {
			log.Fatalf("Error scanning row %d: %v", rownumbah, err)
		}

		// Remove trailing commas from cursorData
		cursorData = strings.TrimSuffix(cursorData, ",")
		results2 = append(results2, CursorData2{CursorData: cursorData})
	}

	fmt.Printf("Number of rows: %d\n", rownumbah)
	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred while iterating rows: %v", err)
	}

	// marshal the results into JSON format
	jsonData, err = json.Marshal(results2)
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to JSON file: %v", err)
	}
	log.Println("First data results written to plcsourcecursordatasection2.json")

	return results1, results2

}
