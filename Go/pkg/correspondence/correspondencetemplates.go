package correspondence

import (
	//"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	//this import grabs the connection strings and env vars
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/formsmigration"
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/platformconfiguration"
)

type CursorData1 struct {
	Name        string `json:"name"`
	CursorData1 string `json:"cursor_for_section_1"`
}

type CursorData2 struct {
	TemplateDisplayId string         `json:"template_display_id"`
	Status            string         `json:"status"`
	Name              string         `json:"name"`
	Association_Type  sql.NullString `json:"association_type"`
	CursorData2       string         `json:"cursor_for_section_2"`
}

type CursorData3 struct {
	Name        string `json:"name"`
	CursorData3 string `json:"cursor_for_section_3"`
}

type CursorData4 struct {
	TenantId                   string `json:"tenant_id"`
	RoleName                   string `json:"role_name"`
	CorrName                   string `json:"corr_name"`
	CanAddFlag                 bool   `json:"can_add_flag"`
	Cursor_for_Roles_section_4 string `json:"cursor_for_roles_section_4"`
}

func CorrespondenceTemplate1(sourceEnv, sourceUserName, sourceDBPassword, tenantID, formGroups string) []CursorData1 {

	log.Printf("Getting Cursor Data Section 1 for Correspondence Source Selection")

	// Get form group names from the file
	formGroups, _ = platformconfiguration.ReadFormGroupNames("corrbundleNameList.txt")
	fmt.Printf("Correspondence Bundle List of Names: %v\n", formGroups)

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

	sqlScript1 := fmt.Sprintf(
		`USE [%s];
	DECLARE      @TenantId       uniqueidentifier = '%s';
	DECLARE        @CorrtypeList   TABLE ( Corrname    nvarchar(200)
	,                                    Templatedisplayid    nvarchar(50)                     );

	INSERT INTO @CorrtypeList ( Corrname )
	values
	-----------------------
	%s
	------------------------------------------------------------------------------------------
	--section 1 

	;WITH TL AS ( SELECT * FROM @CorrtypeList)
	SELECT ct.[Name]
	,      '('''+isnull(ct.[Name],'') +''','''
	+ isnull(cast(ct.[CertifiedMailFlag] as nvarchar(10)),'') +''','''
	+ isnull(ct.[ContextLevel],'') +''','''
	+ isnull(cast(ct.[FTIFlag] as nvarchar(10)),'') +''','''
	+isnull((select distinct wv.code -- workflowsubtype
		from     [ref_WorkflowEngineDB]..[WorkflowVariants] wv 
		join     [ref_NgCorrespondence]..[CorrespondenceType] ct
		on       wv.[WorkflowVariantId]= ct.WorkflowSubTypeId
		where    wv.[WorkflowVariantId] = (select  WorkflowSubTypeId from [CorrespondenceType]  where  name= tl.Corrname and tenantid = @TenantID  )
		and      wv.tenantid = @TenantID ),'') +''',''' 

	--+ isnull(cast(ct.[WorkflowSubTypeId] as nvarchar(25)),'') +''',''' 
	--+ isnull(cast(ct.[WorkflowTypeId] as nvarchar(25)),'') +''','''
	+isnull((select distinct wg.code -- workflowtype
		from     [ref_WorkflowEngineDB]..[WorkflowGroups] wg
		join     [ref_NgCorrespondence]..[CorrespondenceType] ct
		on       wg.[WorkflowGroupId]= ct.[WorkflowTypeId]
		where    wg.[WorkflowGroupId] = (select [WorkflowTypeId]  from [CorrespondenceType]  where  name= tl.Corrname and tenantid = @TenantID  )
		and      wg.tenantid = @TenantID ),'') +''',''' 

	+ isnull(ctc.Name,'')+''',''' 
	+ isnull(cast(pg.[Name]as nvarchar(50)),'') +''',''' 
	+ isnull(cpg.[PrintGroupAssignmentType],'') +''',''' 
	+isnull(cast(ct.[ModuleId] as nvarchar(10)),'') +''','''
	+ isnull(cast(ct.[CorrespondenceTypeDisplayId] as nvarchar(50)),'')

	+'''),'  as cursor_for_section1 

	FROM [dbo].[CorrespondenceType]         ct 
	left JOIN [dbo].[CorrespondenceTypeCategory] ctc ON ct.[CorrespondenceTypeCategoryId]=ctc.[CorrespondenceTypeCategoryId]
	left join [dbo].[CorrespondenceTypePrintGroup] cpg on cpg.[CorrespondenceTypeId] = ct.[CorrespondenceTypeId]
	left join [dbo].[PrintGroup] pg on cpg.[PrintGroupId] = pg.[PrintGroupId]
	Join  tl                               on tl.Corrname = ct.Name
	WHERE ctc. TenantId= @TenantID`, useCorrectDB, tenantID, formGroups)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript1))
	rows, err := db.Query(sqlScript1)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("corrsourcecursordatasection1.json")
	if err != nil {
		log.Fatalf("Error creating corrsourcecursordatasection1.json file: %v", err)
	}
	defer file.Close()

	var results []CursorData1
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var (
			name                 string
			cursor_for_section_1 string
		)

		if err := rows.Scan(&name, &cursor_for_section_1); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursor_for_section_1 = strings.TrimSuffix(cursor_for_section_1, ",")

		results = append(results, CursorData1{
			Name:        name,
			CursorData1: cursor_for_section_1,
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
	log.Println("First data results written to corrsourcecursordatasection1.json")

	return results

}

func CorrespondenceTemplate2(sourceEnv, sourceUserName, sourceDBPassword, tenantID, formGroups string) []CursorData2 {

	log.Printf("Getting Cursor Data Section 2 for Correspondence Source Selection")

	// Get form group names from the file
	formGroups, _ = platformconfiguration.ReadFormGroupNames("corrbundleNameList.txt")
	fmt.Printf("Correspondence Bundle List of Names: %v\n", formGroups)

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

	sqlScript1 := fmt.Sprintf(
		`USE [%s];
	DECLARE      @TenantId       uniqueidentifier = '%s';
	DECLARE        @CorrtypeList   TABLE ( Corrname    nvarchar(200)
	,                                    Templatedisplayid    nvarchar(50)                     );

	INSERT INTO @CorrtypeList ( Corrname )
	values
	-----------------------
	--Section 1
	--> Copy SQL1 from Inventory
	-----------------------
	%s
	---section 2
	;WITH TL AS ( SELECT * FROM @CorrtypeList)
	SELECT distinct
		t.TemplateDisplayId
	,      t.[Status]
	--,      t.[Version]
	--,      t.[Description]
	,    ct.Name,
	--,     (select  a.name from [dbo].[ReusableContentType] a join [dbo].[ReusableContent] b on a.[ReusableContentTypeId] = b.[ReusableContentTypeId] 
	--      where b.[ReusableContentId] = trc.[ReusableContentId]  and b.tenantid =@TenantId   ) as Header_Associated
		(select  trc.ReusableContentAssociationType from [dbo].[ReusableContentType] a join [dbo].[ReusableContent] b on a.[ReusableContentTypeId] = b.[ReusableContentTypeId] 
		where b.[ReusableContentId] = trc.[ReusableContentId]  and b.tenantid =@TenantId   ) as Association_Type,

	case when (select  a.name from [dbo].[ReusableContentType] a left join [dbo].[ReusableContent] b on a.[ReusableContentTypeId] = b.[ReusableContentTypeId] 
		where b.[ReusableContentId] = trc.[ReusableContentId]  and b.tenantid =@TenantId   ) is not null
	then 
	'('''''+','''
			+ isnull(t.[TemplateDisplayId],'') +''',''' 
			+ isnull(t.[Status],'') +''','''
			+ isnull(cast(t.[Version] as nvarchar(10)),'') +''','''
			+ isnull(t.[Description],'') +''',''' 
			+ isnull(ct.Name,'') +''','''
			+isnull(cast(t.[ModuleId] as nvarchar(10)),'') +''','''
			+ (select  a.name from [dbo].[ReusableContentType] a left join [dbo].[ReusableContent] b on a.[ReusableContentTypeId] = b.[ReusableContentTypeId] 
		where b.[ReusableContentId] = trc.[ReusableContentId]  and b.tenantid =@TenantId   ) +''',''' 
		+(select  trc.ReusableContentAssociationType from [dbo].[ReusableContentType] a join [dbo].[ReusableContent] b on a.[ReusableContentTypeId] = b.[ReusableContentTypeId] 
		where b.[ReusableContentId] = trc.[ReusableContentId]  and b.tenantid =@TenantId   )
			+'''),'
	else 
	'('''''+','''
			+ isnull(t.[TemplateDisplayId],'') +''',''' 
			+ isnull(t.[Status],'') +''',''' 
			+ isnull(cast(t.[Version] as nvarchar(10)),'') +''','''
			+ isnull(t.[Description],'') +''',''' 
			+ isnull(ct.Name,'') +''','''
			+isnull(cast(t.[ModuleId] as nvarchar(10)),'') +''','
			+ 'null' +','
			+ 'null' 
			+'),'
	end as cursor_for_section2
	FROM Template             t
	JOIN [CorrespondenceType] ct ON ct.CorrespondenceTypeId = t.CorrespondenceTypeId
	left join [dbo].[TemplateReusableContent] trc on trc.[TemplateId] = t.[TemplateId]
	join tl           on tl.Corrname =  ct.Name
	where ct.name = tl.Corrname
	and t.status in('Active' )
	and t.tenantid = @TenantID
	--and t.[TemplateDisplayId] = @tempId
	--IN ('TMP000000682')`, useCorrectDB, tenantID, formGroups)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript1))
	rows, err := db.Query(sqlScript1)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("corrsourcecursordatasection2.json")
	if err != nil {
		log.Fatalf("Error creating corrsourcecursordatasection2.json file: %v", err)
	}
	defer file.Close()

	var results []CursorData2
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var (
			template_display_id  string
			status               string
			name                 string
			association_type     sql.NullString
			cursor_for_section_2 string
		)

		if err := rows.Scan(&template_display_id, &status, &name, &association_type, &cursor_for_section_2); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursor_for_section_2 = strings.TrimSuffix(cursor_for_section_2, ",")

		results = append(results, CursorData2{
			TemplateDisplayId: template_display_id,
			Status:            status,
			Name:              name,
			Association_Type:  association_type,
			CursorData2:       cursor_for_section_2,
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
	log.Println("First data results written to corrsourcecursordatasection2.json")

	return results

}

func CorrespondenceTemplate3(sourceEnv, sourceUserName, sourceDBPassword, tenantID, formGroups string) []CursorData3 {

	log.Printf("Getting Cursor Data Section 3 for Correspondence Source Selection")

	// Get form group names from the file
	formGroups, _ = platformconfiguration.ReadFormGroupNames("corrbundleNameList.txt")
	fmt.Printf("Correspondence Bundle List of Names: %v\n", formGroups)

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

	sqlScript1 := fmt.Sprintf(
		`USE [%s];
	DECLARE      @TenantId       uniqueidentifier = '%s';
	DECLARE        @CorrtypeList   TABLE ( Corrname    nvarchar(200)
	,                                    Templatedisplayid    nvarchar(50)                     );

	INSERT INTO @CorrtypeList ( Corrname )
	values
-----------------------
%s
-----------------------

	;WITH TL AS ( SELECT * FROM @CorrtypeList)
	SELECT ct.[Name]
	,      '('''+isnull(ct.[Name],'') +''','''
	+ isnull(cast(ct.[CertifiedMailFlag] as nvarchar(10)),'') +''','''
	+ isnull(ct.[ContextLevel],'') +''','''
	+ isnull(cast(ct.[FTIFlag] as nvarchar(10)),'') +''','''
	+ isnull(cast(ct.[WorkflowSubTypeId] as nvarchar(10)),'') +''',''' 
	+ isnull(cast(ct.[WorkflowTypeId] as nvarchar(10)),'') +''','''
	+ isnull(ctc.Name,'')+''','''
	+isnull(cast(ct.[ModuleId] as nvarchar(10)),'') +''','''
	+ isnull(cast(ct.[CorrespondenceTypeDisplayId] as nvarchar(50)),'') +'''),'  as cursor_for_section3 

	FROM [dbo].[CorrespondenceType]         ct 
	left JOIN [dbo].[CorrespondenceTypeCategory] ctc ON ct.[CorrespondenceTypeCategoryId]=ctc.[CorrespondenceTypeCategoryId]
	left join [dbo].[CorrespondenceTypePrintGroup] cpg on cpg.[CorrespondenceTypeId] = ct.[CorrespondenceTypeId]
	left join [dbo].[PrintGroup] pg on cpg.[PrintGroupId] = pg.[PrintGroupId]
	Join  tl                               on tl.Corrname = ct.Name
	WHERE ctc. TenantId= @TenantID`, useCorrectDB, tenantID, formGroups)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript1))
	rows, err := db.Query(sqlScript1)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("corrsourcecursordatasection3.json")
	if err != nil {
		log.Fatalf("Error creating corrsourcecursordatasection3.json file: %v", err)
	}
	defer file.Close()

	var results []CursorData3
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var (
			name                 string
			cursor_for_section_3 string
		)

		if err := rows.Scan(&name, &cursor_for_section_3); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursor_for_section_3 = strings.TrimSuffix(cursor_for_section_3, ",")

		results = append(results, CursorData3{
			Name:        name,
			CursorData3: cursor_for_section_3,
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
	log.Println("First data results written to corrsourcecursordatasection3.json")

	return results

}

func CorrespondenceTemplate4(sourceEnv, sourceUserName, sourceDBPassword, tenantID, formGroups string) []CursorData4 {

	log.Printf("Getting Cursor Data Section 4 for Correspondence Source Selection")

	// Get form group names from the file
	formGroups, _ = platformconfiguration.ReadFormGroupNames("corrbundleNameList.txt")
	fmt.Printf("Correspondence Bundle List of Names: %v\n", formGroups)

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

	var useCorrectDB = "NgRoleManagement"
	// get sourceEnv and add as a prefix to the DB
	if sourceEnv != "" {
		if strings.ToLower(sourceEnv) == "ref" {
			useCorrectDB = sourceEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB1: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgRoleManagement"
			fmt.Printf("Use Correct DB2: %s\n", useCorrectDB)
		}
	}

	sqlScript1 := fmt.Sprintf(
		`USE [%s];
	DECLARE      @TenantId       uniqueidentifier = '%s';
	DECLARE        @CorrtypeList   TABLE ( Corrname    nvarchar(200)
	,                                    Templatedisplayid    nvarchar(50)                     );

	INSERT INTO @CorrtypeList ( Corrname )
	values
	-----------------------
	%s
	-----------------------
	;WITH TL AS ( SELECT * FROM @CorrtypeList)
	SELECT a.TenantId                                                                                                      
	,      a.name           roleName
	,      c.CorrName                                                                                                      
	,      b.CanAddFlag                                                                                                    
	,      '('''+isnull(a.[Name],'') +''','''
				+ isnull(c.CorrName,'') +''','''
				+isnull(cast(b.[ModuleId] as nvarchar(10)),'') +''','''
				+ isnull(cast(b.CanAddFlag AS varchar(10)),'') +'''),' AS cursor_for_Roles_section_4

	FROM ref_NgRoleManagement..role                         a
	JOIN ref_NgCorrespondence..[CorrespondenceTypeUserRole] b ON b.UserRoleId = a.RoleId
	JOIN ref_NgCorrespondence..[CorrespondenceType]         t ON t.[CorrespondenceTypeId] = b.CorrespondenceTypeId
	JOIN @CorrtypeList                                         c ON c.Corrname = t.Name
	WHERE a.tenantid IN( @TenantID,'00000000-0000-0000-0000-000000000000')
		--AND a.name LIKE 'STL%'
		AND t.tenantid = @TenantID`, useCorrectDB, tenantID, formGroups)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript1))
	rows, err := db.Query(sqlScript1)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	//open a file to write the JSON data
	file, err := os.Create("corrsourcecursordatasection4.json")
	if err != nil {
		log.Fatalf("Error creating corrsourcecursordatasection4.json file: %v", err)
	}
	defer file.Close()

	var results []CursorData4
	var rowNumber = 0

	for rows.Next() {
		rowNumber++
		var (
			TenantId                   string
			roleName                   string
			CorrName                   string
			CanAddFlag                 bool
			cursor_for_Roles_section_4 string
		)

		if err := rows.Scan(&TenantId, &roleName, &CorrName, &CanAddFlag, &cursor_for_Roles_section_4); err != nil {
			log.Fatalf("Error scanning row %d: %v", rowNumber, err)
		}

		// Remove trailing commas from cursorData
		cursor_for_Roles_section_4 = strings.TrimSuffix(cursor_for_Roles_section_4, ",")

		results = append(results, CursorData4{
			TenantId:                   TenantId,
			RoleName:                   roleName,
			CorrName:                   CorrName,
			CanAddFlag:                 CanAddFlag,
			Cursor_for_Roles_section_4: cursor_for_Roles_section_4,
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
	log.Println("First data results written to corrsourcecursordatasection4.json")

	return results

}
