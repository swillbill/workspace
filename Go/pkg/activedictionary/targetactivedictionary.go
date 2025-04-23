package activedictionary

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/formsmigration"
)

func RunTargetActiveDirectorySQL1(targetEnv, targetTenantId string, targetUserName string, targetDBPassword string, data []SystemFieldDictionary) (err error) {
	log.Printf("Beginning Target Active Directory SQL Migration Script 1")

	sqlConfig, ok := formsmigration.GetMigrationSQLConfig(targetEnv, targetUserName, targetDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", targetEnv)
		return fmt.Errorf("invalid environment: %s", targetEnv)
	}

	db, err := sql.Open("sqlserver", sqlConfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	// read artifacted file
	file, err := os.Open("ADsourcecursordata1.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue) == 0 {
		log.Fatal("ADsourcecursordata1.json file is empty")
	}

	var unmarshaledData []SystemFieldDictionary
	err = json.Unmarshal(byteValue, &unmarshaledData)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON := unmarshaledData[0].CursorData
	for numbah, rows := range data {
		if numbah == 0 {
			finishedJSON = rows.CursorData
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON = fmt.Sprintf("%s,%s", finishedJSON, rows.CursorData)
		}

	}

	// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Finished JSON: %s", finishedJSON))
	var useCorrectDB = "NgCorrespondence"
	// get sourceEnv and add as a prefix to the DB
	if targetEnv != "" {
		if strings.ToLower(targetEnv) == "ref" {
			useCorrectDB = targetEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		}
	}

	sqlScript1 := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';

		SET NOCOUNT ON
		SET XACT_ABORT ON;

		DECLARE @Code          nvarchar(150)                 
		,       @Description   nvarchar(250)                 
		,       @CreateDate    datetime                       = getutcdate()
		,       @CreatedBy     nvarchar(450)                  = 'AdminUser'
		,       @UpdatedDttm   datetime                       = getutcdate()
		,       @UpdatedBy     nvarchar(450)                  = 'AdminUser'
		,       @newid         uniqueidentifier              
		,       @debug         bit -- set to 1 for dubugging 
		,       @commit        bit -- set to 1 for committing
		,       @err_message   nvarchar(1000)                
		,       @CorrelationId uniqueidentifier               = NULL
		,       @SeqID         int                           
		,       @ModuleID      int                           
		--,       @ActiveFlag     bit null

		SET @debug = 1; --For testing 
		SET @commit = 1; --For testing 

		BEGIN TRY
			BEGIN TRANSACTION

		PRINT '-----------------------------------------------------------'
		PRINT 'Tenant: '+cast(@TenantId AS varchar(150))
		PRINT '-----------------------------------------------------------'

		---1. Insert into [SystemFieldDictionary]-------------------------
		----------------------------------------------------------------------------------
		DECLARE @Status                        nvarchar(25)    
		,       @LogixInputSchemaJson          nvarchar(MAX)   
		,       @FixedId                       int             
		,       @Version                       int              = 1
		,       @SystemFieldDictionaryUniqueId uniqueidentifier
		,       @IsOOTB                        bit              = 1
		,       @IsOOTBEditable                bit              = 0
		,       @IsDeleted                     int              = 0
		,       @SDescription                  nvarchar(250)    = NULL
		,@Newversion   int 


		DECLARE curListDict CURSOR FOR SELECT *
		FROM (
		VALUES 

		--Section 1: Data for SystemFieldDictionary
		/*

		
		SELECT 
		[FixedId]                                                                 
		,      [Status]                                                                  
		,     [Version]    
		,      [LogixInputSchemaJson]                                                    
		,      SystemFieldDictionaryUniqueId                                             
		,      SystemFieldDictionaryId                                                   
		,      ModuleId                                                                  
		,      '('''+isnull(cast([FixedId] AS nvarchar(10)),'') +''','''+ isnull([Status],'') +''',
		'''+ isnull(cast([Version] AS nvarchar(10)),'') +''',
		''' + isnull(cast([ModuleId] AS nvarchar(10)),'') 
				+''','''+ isnull(replace([LogixInputSchemaJson],'''',''''''),'')
				+''','''+ isnull(cast([SystemFieldDictionaryUniqueId] AS nvarchar(50)),'')+'''),' cursordata
		FROM [SystemFieldDictionary]
		WHERE [TenantId] = 'aa106e70-19f7-4269-9eeb-fd6e9bfdc88b'
			AND status IN ( 'ACTIVE') 
			--AND status IN ( 'ACTIVE', 'DRAFT')
		
		AMazur: some JSON elements have single quotes, might need to rework using the approach with:  SET QUOTED_IDENTIFIER OFF; "It's" SET QUOTED_IDENTIFIER ON;
		
		select max(version) from SystemFieldDictionary where tenantid = 'dbe67321-c27a-4035-8266-5979a0b6b730' and status = 'Active'

		*/

		-----------------------------------------------------------------------------------------------
		%s
		-----------------------------------------------------------------------------------------------

		--End of Section 1: Data for SystemFieldDictionary
		) v ([FixedId],[Status],[Version],[ModuleId],[LogixInputSchemaJson],[SystemFieldDictionaryUniqueId])

		PRINT '--Table 1. [SystemFieldDictionary]--------------------------------------------------------------------------'

		UPDATE SystemFieldDictionary
		SET status =    'INACTIVE'
		,   updatedby = 'AdminUser'
		WHERE TenantId = @TenantId
			AND Status = 'ACTIVE'

		OPEN curListDict
			FETCH NEXT FROM curListDict
			INTO @FixedId,@Status,@version,@ModuleId, @LogixInputSchemaJson, @SystemFieldDictionaryUniqueId

		WHILE @@FETCH_STATUS = 0
		BEGIN

			IF EXISTS (SELECT 1  FROM SystemFieldDictionary WHERE SystemFieldDictionaryUniqueId = @SystemFieldDictionaryUniqueId and TenantId= @TenantId and version= @version and ModuleId=@ModuleId)
				IF @debug = 1 PRINT 'System Field Dictionary already exists: SystemFieldDictionaryUniqueId = ' + cast(@SystemFieldDictionaryUniqueId as nvarchar(100)) + ', Version = '+cast(@version as nvarchar(100)) +', ModuleId = '+cast(@ModuleId as nvarchar(100))

		SET @Newversion = (select max(version) from SystemFieldDictionary  WHERE SystemFieldDictionaryUniqueId = @SystemFieldDictionaryUniqueId and TenantId= @TenantId and ModuleId=@ModuleId)+1
		select @Newversion

			IF NOT EXISTS (SELECT 1  FROM SystemFieldDictionary WHERE SystemFieldDictionaryUniqueId = @SystemFieldDictionaryUniqueId and TenantId= @TenantId and (version = @Newversion  and version !=@Version)  and ModuleId=@ModuleId)

			BEGIN

				INSERT INTO [SystemFieldDictionary] ( [FixedId], [Status], [LogixInputSchemaJson], [Version], [CreatedDate], [CreatedBy], [UpdatedDate], [UpdatedBy], [IsOOTB], [IsOOTBEditable], [SystemFieldDictionaryUniqueId], [TenantId], [ModuleId] )
				VALUES                              ( @FixedId,  @Status,  @LogixInputSchemaJson,  @Newversion,  @CreateDate,   @CreatedBy,  @UpdatedDttm,  @UpdatedBy,  @IsOOTB,  @IsOOTBEditable,  @SystemFieldDictionaryUniqueId,  @TenantId, @ModuleId  );

				IF @debug = 1 PRINT 'System Field Dictionary inserted successfully: SystemFieldDictionaryUniqueId = ' + cast(@SystemFieldDictionaryUniqueId as nvarchar(100)) + ', Version = '+cast(@Newversion as nvarchar(100)) +', ModuleId = '+cast(@ModuleId as nvarchar(100))
			
			END

			FETCH NEXT FROM curListDict
					INTO @FixedId,@Status,@version,@ModuleId, @LogixInputSchemaJson, @SystemFieldDictionaryUniqueId

		END 

		CLOSE curListDict;
		DEALLOCATE curListDict;`, useCorrectDB, targetTenantId, finishedJSON)

	//log.Printf("Target Active Directory SQL Migration Script: %s", sqlScript1)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript1))
	lol, err := db.Prepare(sqlScript1)
	if err != nil {
		log.Fatalf("Error preparing AD SQL script 1: %v", err)
	}
	result, err := lol.Exec()
	//result, err := db.Exec(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	//debugtown(fmt.Sprintf("RunSQLMigration: Results: %s", result))
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error fetching rows affected: %v", err)
	}

	log.Printf("AD SQL script 1 executed successfully. Rows affected: %d", rowsAffected)
	return nil

}

func RunTargetActiveDirectorySQL2(targetEnv, targetTenantId string, targetUserName string, targetDBPassword string, data []SystemFieldDictionaryOOTBSection) (err error) {

	log.Printf("Beginning Target Active Directory SQL Migration Script 2")

	sqlConfig, ok := formsmigration.GetMigrationSQLConfig(targetEnv, targetUserName, targetDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", targetEnv)
		return fmt.Errorf("invalid environment: %s", targetEnv)
	}

	db, err := sql.Open("sqlserver", sqlConfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	// read artifacted file
	file, err := os.Open("ADsourcecursordata2.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue) == 0 {
		log.Fatal("ADsourcecursordata2.json file is empty")
	}

	var unmarshaledData []SystemFieldDictionaryOOTBSection
	err = json.Unmarshal(byteValue, &unmarshaledData)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON := unmarshaledData[0].CursorData
	for numbah, rows := range data {
		if numbah == 0 {
			finishedJSON = rows.CursorData
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL2: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON = fmt.Sprintf("%s,%s", finishedJSON, rows.CursorData)
		}

	}

	// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL2: Finished JSON: %s", finishedJSON))
	var useCorrectDB = "NgCorrespondence"
	// get sourceEnv and add as a prefix to the DB
	if targetEnv != "" {
		if strings.ToLower(targetEnv) == "ref" {
			useCorrectDB = targetEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		}
	}

	sqlScript2 := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';


		DECLARE @SectionJson                  nvarchar(max)
		,       @ContextName                  nvarchar(60)
		,       @SectionName                  nvarchar(60)
		,       @SystemFieldDictionaryId      int
		,       @DContextname                 nvarchar(60)
		,       @DSectionJson                 nvarchar(max)
		
		DECLARE curListSec CURSOR FOR SELECT *
		FROM (
		VALUES
		-- data for SystemFieldDictionaryOOTBSection
		/*
		SELECT o.[SectionJson]                                                                                                                                                           
		,      o.[ContextName]                                                                                                                                                           
		,      o.[SectionName]                                                                                                                                                           
		,      o.[ModuleId]                                                                                                                                                              
		,      '('''+isnull(o.[SectionJson],'') +''','''+ isnull(o.[ContextName],'') +''','''+ isnull(o.[SectionName],'')+''','''+ isnull(cast(o.[ModuleId] AS nvarchar(10)),'') + '''),'
		cursordata
		FROM [SystemFieldDictionaryOOTBSection] o
		JOIN [SystemFieldDictionary]            d ON d.[SystemFieldDictionaryId] = o.[SystemFieldDictionaryId]
		WHERE o.[TenantId] = 'aa106e70-19f7-4269-9eeb-fd6e9bfdc88b'
			AND d.[Status] = 'ACTIVE'
		*/
		
		--Section 2: Data for SystemFieldDictionaryOOTBSection
		-----------------------------------------------------------------------------------------------
		%s
		-----------------------------------------------------------------------------------------------
		--Section 2: End of data for SystemFieldDictionaryOOTBSection
		
		) v ([SectionJson], [ContextName], [SectionName], [ModuleId])
		
		OPEN curListSec
			FETCH NEXT FROM curListSec
			INTO @SectionJson, @ContextName, @SectionName, @ModuleId;
		
		PRINT '--Table 2. [SystemFieldDictionaryOOTBSection]--------------------------------------------------------------------------'
		
		WHILE @@FETCH_STATUS = 0
		BEGIN
		
			DECLARE @SystemFieldDictId int;
			DECLARE @cnt int;
			SET @SystemFieldDictId = (SELECT SystemFieldDictionaryId FROM SystemFieldDictionary WHERE TenantId = @TenantId and Status = 'ACTIVE' and ModuleId=@ModuleID)
		
			IF(@SystemFieldDictId > 0)
			BEGIN		
				SELECT @cnt = count(*) FROM [SystemFieldDictionaryOOTBSection] WHERE TenantId= @TenantId AND SectionName = @SectionName AND SystemFieldDictionaryId = @SystemFieldDictId;
				IF (@cnt > 0)
				BEGIN
					IF @debug = 1 PRINT 'System Field Dictionary OOTBSection already exists for ' + @SectionName + ' SystemFieldDictionaryId: '+cast(@SystemFieldDictId as nvarchar(100)) + ' //ACTIVE'
				END
				ELSE
				BEGIN
					INSERT INTO [SystemFieldDictionaryOOTBSection] ( [SectionJson], [ContextName], [SectionName], [SystemFieldDictionaryId], [CreatedDate], [CreatedBy], [UpdatedDate], [UpdatedBy], [TenantId], ModuleId  )
					VALUES                                         ( @SectionJson,  @ContextName,  @SectionName,  @SystemFieldDictId,        @CreateDate,   @CreatedBy,  @UpdatedDttm,  @updatedby,  @TenantId,  @ModuleId );
		
				IF @debug = 1 PRINT 'System Field Dictionary OOTBSection inserted for ' + @SectionName + ' SystemFieldDictionaryId: '+cast(@SystemFieldDictId as nvarchar(100)) + ' //ACTIVE'
				END
		
			END
		
			SET @SystemFieldDictId = (SELECT SystemFieldDictionaryId FROM SystemFieldDictionary WHERE TenantId = @TenantId AND Status = 'DRAFT' AND ModuleId=@ModuleID)
		
			IF(@SystemFieldDictId > 0)
			BEGIN
				SELECT @cnt = count(*) FROM [SystemFieldDictionaryOOTBSection] WHERE TenantId= @TenantId AND SectionName = @SectionName AND SystemFieldDictionaryId = @SystemFieldDictId;
				IF (@cnt > 0)
				BEGIN
					IF @debug = 1 PRINT 'System Field Dictionary OOTBsection already exists for ' + @SectionName + ' SystemFieldDictionaryId: '+cast(@SystemFieldDictId as nvarchar(100)) + ' //DRAFT'
				END
				ELSE
				BEGIN
					INSERT INTO [SystemFieldDictionaryOOTBSection] ( [SectionJson], [ContextName], [SectionName], [SystemFieldDictionaryId], [CreatedDate], [CreatedBy], [UpdatedDate], [UpdatedBy], [TenantId], ModuleId  )
					VALUES                                         ( @SectionJson,  @ContextName,  @SectionName,  @SystemFieldDictId,        @CreateDate,   @CreatedBy,  @UpdatedDttm,  @updatedby,  @TenantId,  @ModuleId );
		
					IF @debug = 1 PRINT 'System Field Dictionary OOTBSection inserted for ' + @SectionName + ' SystemFieldDictionaryId: '+cast(@SystemFieldDictId as nvarchar(100)) + ' //DRAFT'
				END
		 
			END
		
			FETCH NEXT FROM curListSec
			INTO @SectionJson, @ContextName, @SectionName, @ModuleId;
		
		END
		
		CLOSE curListSec;
		DEALLOCATE curListSec;`, useCorrectDB, targetTenantId, finishedJSON)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript2))
	lol, err := db.Prepare(sqlScript2)
	if err != nil {
		log.Fatalf("Error preparing AD SQL script 2: %v", err)
	}
	result, err := lol.Exec()
	//result, err := db.Exec(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	//debugtown(fmt.Sprintf("RunSQLMigration: Results: %s", result))
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error fetching rows affected: %v", err)
	}

	log.Printf("AD SQL script 2 executed successfully. Rows affected: %d", rowsAffected)

	return nil
}

func RunTargetActiveDirectorySQL3(targetEnv, targetTenantId string, targetUserName string, targetDBPassword string, data []SystemFieldDictionarySection) (err error) {

	log.Printf("Beginning Target Active Directory SQL Migration Script 3")

	sqlConfig, ok := formsmigration.GetMigrationSQLConfig(targetEnv, targetUserName, targetDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", targetEnv)
		return fmt.Errorf("invalid environment: %s", targetEnv)
	}

	db, err := sql.Open("sqlserver", sqlConfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	// read artifacted file
	file, err := os.Open("ADsourcecursordata3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue) == 0 {
		log.Fatal("ADsourcecursordata3.json file is empty")
	}

	var unmarshaledData []SystemFieldDictionarySection
	err = json.Unmarshal(byteValue, &unmarshaledData)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON := unmarshaledData[0].CursorData
	for numbah, rows := range data {
		if numbah == 0 {
			finishedJSON = rows.CursorData
		} else {
			finishedJSON = fmt.Sprintf("%s,%s", finishedJSON, rows.CursorData)
		}

	}

	var useCorrectDB = "NgCorrespondence"
	// get sourceEnv and add as a prefix to the DB
	if targetEnv != "" {
		if strings.ToLower(targetEnv) == "ref" {
			useCorrectDB = targetEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		}
	}

	sqlScript3 := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';
		DECLARE @DSSectionJson                  nvarchar(max)
		,       @DSContextName                  nvarchar(60)
		,       @DSSystemFieldDictionaryId      int

		DECLARE curListSec2 CURSOR FOR SELECT *
		FROM (
		VALUES

		/*
		SELECT s.[SectionJson]                                                                                                                          
		,      s.[ContextName]                                                                                                                          
		,      d. [SystemFieldDictionaryId]                                                                                                             
		,      s. ModuleId                                                                                                                              
		,      '('''+isnull(s.[SectionJson],'') +''','''+ isnull(s.[ContextName],'')+''',''' + isnull(cast(s.[ModuleId] AS nvarchar(10)),'') +'' +'''),' cursordata

		FROM [SystemFieldDictionarySection] s
		JOIN [SystemFieldDictionary]        d ON d.[SystemFieldDictionaryId] = s.[SystemFieldDictionaryId]
		WHERE d.[TenantId] = 'aa106e70-19f7-4269-9eeb-fd6e9bfdc88b'
			AND d.[Status] = 'ACTIVE'
		*/

		--Section 3: Data for SystemFieldDictionarySection
		-----------------------------------------------------------------------------------------------
		%s
		----------------------------------------------------------------------------------------------
		--Section 3:End of data for SystemFieldDictionarySection
		) v ([SectionJson],[ContextName], [ModuleID])

		OPEN curListSec2
		FETCH NEXT FROM curListSec2
		INTO @DSSectionJson, @DSContextName , @ModuleId;

		PRINT '--Table 3. [SystemFieldDictionarySection]--------------------------------------------------------------------------'

		WHILE @@FETCH_STATUS = 0
		BEGIN

			SET @SystemFieldDictId = (SELECT SystemFieldDictionaryId FROM SystemFieldDictionary WHERE TenantId = @TenantId and Status = 'ACTIVE' and ModuleId=@ModuleId)

				--IF not exists(select 1 from [SystemFieldDictionarySection] where [TenantId]= @TenantId and [SystemFieldDictionaryId] = 
				--(select [SystemFieldDictionaryId] from [SystemFieldDictionary] where [TenantId]= @TenantId and status = 'ACTIVE'))

			BEGIN
				INSERT INTO [SystemFieldDictionarySection] ( [SectionJson],  [ContextName],  [SystemFieldDictionaryId], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [CorrelationId], ModuleId  )
				VALUES                                     ( @DSSectionJson, @DSContextName, @SystemFieldDictId,        @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantId,  @CorrelationId,  @ModuleId )

				IF @debug = 1 PRINT 'System Field Dictionary section inserted: SectionName = '+@DSContextName+ ', ModuleId = '+cast(@ModuleId as nvarchar(100))
			END

			FETCH NEXT FROM curListSec2
			INTO @DSSectionJson, @DSContextName , @ModuleId;

		END
		CLOSE curListSec2;
		DEALLOCATE curListSec2;`, useCorrectDB, targetTenantId, finishedJSON)

	//log.Printf("Target Active Directory SQL Migration Script: %s", sqlScript3)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript3))
	lol, err := db.Prepare(sqlScript3)
	if err != nil {
		log.Fatalf("Error preparing AD SQL script 3: %v", err)
	}
	result, err := lol.Exec()
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	// debugtown(fmt.Sprintf("RunSQLMigration: Results: %s", result))
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error fetching rows affected: %v", err)
	}

	log.Printf("AD SQL script 3 executed successfully. Rows affected: %d", rowsAffected)
	return nil

}

func RunTargetActiveDirectorySQL4(targetEnv, targetTenantId string, targetUserName string, targetDBPassword string, data []SystemFieldDictionaryOOTBField) (err error) {

	log.Printf("Beginning Target Active Directory SQL Migration Script 4")

	sqlConfig, ok := formsmigration.GetMigrationSQLConfig(targetEnv, targetUserName, targetDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", targetEnv)
		return fmt.Errorf("invalid environment: %s", targetEnv)
	}

	db, err := sql.Open("sqlserver", sqlConfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	// read artifacted file
	file, err := os.Open("ADsourcecursordata3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue) == 0 {
		log.Fatal("ADsourcecursordata3.json file is empty")
	}

	var unmarshaledData []SystemFieldDictionaryOOTBField
	err = json.Unmarshal(byteValue, &unmarshaledData)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON := unmarshaledData[0].CursorData
	for numbah, rows := range data {
		if numbah == 0 {
			finishedJSON = rows.CursorData
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL3: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON = fmt.Sprintf("%s,%s", finishedJSON, rows.CursorData)
		}

	}

	// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL3: Finished JSON: %s", finishedJSON))

	var useCorrectDB = "NgCorrespondence"
	// get sourceEnv and add as a prefix to the DB
	if targetEnv != "" {
		if strings.ToLower(targetEnv) == "ref" {
			useCorrectDB = targetEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgCorrespondence"
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		}
	}

	sqlScript4 := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';
		DECLARE @FieldJson            nvarchar(MAX)
		,       @FieldSectionName     nvarchar(60)
		,       @FieldName            nvarchar(60)
		,       @SystemFieldSectionId int

		DECLARE curListField CURSOR FOR SELECT *
		FROM (
		VALUES

		-- data for SystemFieldDictionaryOOTBField
		/*
		SELECT f.[FieldJson]                                                                                                                                                         
		,      f.[SectionName]                                                                                                                                                       
		,      f.[FieldName]                                                                                                                                                         
		,      f.ModuleId                                                                                                                                                            
		,      '('''+isnull(f.[FieldJson],'') +''','''+ isnull(f.[SectionName],'') +''','''+ isnull(f.[FieldName],'') +''',''' + isnull(cast(f.[ModuleId] AS nvarchar(10)),'')+'''),' cursordata
		FROM [SystemFieldDictionaryOOTBField]   f 
		JOIN [SystemFieldDictionaryOOTBSection] os ON os.[SystemFieldDictionaryOOTBSectionId] = f.[SystemFieldDictionaryOOTBSectionId]
		JOIN [SystemFieldDictionary]            s  ON s.[SystemFieldDictionaryId] = os.[SystemFieldDictionaryId]
		WHERE f.[TenantId] = 'aa106e70-19f7-4269-9eeb-fd6e9bfdc88b'
			AND s.[Status] = 'ACTIVE'
		*/

		--Section 4: Data for SystemFieldDictionaryOOTBField
		-----------------------------------------------------------------------------------------------
		%s
		-----------------------------------------------------------------------------------------------
		--Section 4:End of data for SystemFieldDictionaryOOTBField
		) v ([FieldJson], [SectionName], [FieldName], [ModuleId])

		OPEN curListField
		FETCH NEXT FROM curListField
		INTO
		@FieldJson, @FieldSectionName, @FieldName, @ModuleId;

		PRINT '--Table 4. [SystemFieldDictionaryOOTBField]--------------------------------------------------------------------------'

		WHILE @@FETCH_STATUS = 0
		BEGIN
			DECLARE @SystemFieldSectId int = 0;	

			SET @SystemFieldSectId = (SELECT SystemFieldDictionaryOOTBSectionId FROM SystemFieldDictionaryOOTBSection ss	
												INNER JOIN SystemFieldDictionary sd ON ss.SystemFieldDictionaryId = sd.SystemFieldDictionaryId and ss.ModuleId=sd.ModuleId and sd.Status = 'ACTIVE'
												WHERE ss.SectionName = @FieldSectionName and ss.TenantId = @TenantId and ss.ModuleId = @ModuleID)

			IF EXISTS (SELECT 1 FROM [SystemFieldDictionaryOOTBField] WHERE TenantId= @TenantId AND FieldName = @FieldName and [SystemFieldDictionaryOOTBSectionId] = @SystemFieldSectId)

			--[SystemFieldDictionaryOOTBSectionId]=
			--(SELECT o.[SystemFieldDictionaryOOTBSectionId]
			--FROM [SystemFieldDictionaryOOTBSection] o
			--join [SystemFieldDictionaryOOTBField] f
			--on o.[SystemFieldDictionaryOOTBSectionId] = f.[SystemFieldDictionaryOOTBSectionId]
			--join [SystemFieldDictionary] d
			--on d.[SystemFieldDictionaryId] = o.[SystemFieldDictionaryId]
			--WHERE o.[TenantId] = @TenantId
			--and d.[Status] = 'ACTIVE') )

			BEGIN
				IF @debug = 1 PRINT 'System Field Dictionary field already exists: FieldName = ' + @FieldName + ', ModuleId = '+cast(@ModuleId as nvarchar(100))
			END

			IF NOT EXISTS (SELECT 1 FROM [SystemFieldDictionaryOOTBField] WHERE TenantId= @TenantId AND FieldName = @FieldName and [SystemFieldDictionaryOOTBSectionId] = @SystemFieldSectId)
			BEGIN
				
				SET @SystemFieldDictId = (SELECT SystemFieldDictionaryId FROM SystemFieldDictionary WHERE TenantId = @TenantId and Status = 'ACTIVE' and ModuleId=@ModuleId)

				IF(@SystemFieldSectId > 0 AND @SystemFieldDictId > 0)
				BEGIN

					INSERT INTO [SystemFieldDictionaryOOTBField] ( [FieldJson], [FieldName], [SectionName],     [SystemFieldDictionaryOOTBSectionId], [CreatedDate], [CreatedBy], [UpdatedDate], [UpdatedBy], [TenantId], [ModuleId] )
					VALUES                                       ( @FieldJson,  @FieldName,  @FieldSectionName, @SystemFieldSectId,                   @CreateDate,   @CreatedBy,  @UpdatedDttm,  @UpdatedBy,  @TenantId,  @ModuleId  );

					IF @debug = 1 PRINT 'System Field Dictionary OOTB field inserted: FieldName = ' + @FieldName + ', ModuleId = '+cast(@ModuleId as nvarchar(100)) + ' //ACTIVE'

					SET @SystemFieldSectId = 0;
				END 


				SET @SystemFieldSectId = (SELECT SystemFieldDictionaryOOTBSectionId FROM SystemFieldDictionaryOOTBSection ss	
											INNER JOIN SystemFieldDictionary sd ON ss.SystemFieldDictionaryId = sd.SystemFieldDictionaryId and ss.ModuleId=sd.ModuleId AND sd.Status = 'DRAFT'
											WHERE ss.SectionName = @FieldSectionName and ss.TenantId = @TenantId and ss.ModuleId = @ModuleID)

				IF(@SystemFieldSectId > 0 AND @SystemFieldDictId > 0)
				BEGIN
					INSERT INTO [SystemFieldDictionaryOOTBField] ( [FieldJson], [FieldName], [SectionName],     [SystemFieldDictionaryOOTBSectionId], [CreatedDate], [CreatedBy], [UpdatedDate], [UpdatedBy], [TenantId], [ModuleId] )
					VALUES                                       ( @FieldJson,  @FieldName,  @FieldSectionName, @SystemFieldSectId,                   @CreateDate,   @CreatedBy,  @UpdatedDttm,  @UpdatedBy,  @TenantId , @ModuleId );
				
					IF @debug = 1 PRINT 'System Field Dictionary OOTB field inserted: FieldName = ' + @FieldName  + ', ModuleId = '+cast(@ModuleId as nvarchar(100)) + ' //DRAFT'
				END 

			END

			FETCH NEXT FROM curListField
			into @FieldJson, @FieldSectionName, @FieldName, @ModuleId;

		END
		CLOSE curListField;
		DEALLOCATE curListField;

			DELETE FROM [TemplateCache]
			WHERE tenantid = @TenantId
			IF @debug = 1 PRINT '[TemplateCache] truncated for - : '+cast(@TenantId AS varchar(150))

		-----------------------------------------------------------------
			--Tenant Summary--

			print '-----------------------------------------------------------' 
			print 'Counts for the tenant: '+cast(@TenantId as varchar(150)) 
			print '-----------------------------------------------------------'

			SELECT count(*) as [SystemFieldDictionary] FROM [SystemFieldDictionary] where TenantId=@TenantId	
			select count(*) as [SystemFieldDictionaryOOTBField] from [SystemFieldDictionaryOOTBField] where TenantId=@TenantId
			SELECT count(*) as [SystemFieldDictionaryOOTBSection] FROM [SystemFieldDictionaryOOTBSection] where TenantId=@TenantId
			select count(*) as [SystemFieldDictionarySection] from [SystemFieldDictionarySection]

			IF @commit =1 
				BEGIN
					COMMIT TRANSACTION
					print '-----------------------------------------------------------'
					print  'Committed!'
				END
			ELSE
				BEGIN
					print '-----------------------------------------------------------'
					print  'Rolled back!'
					ROLLBACK TRANSACTION
				END

		END TRY

		BEGIN CATCH
			THROW;

			WHILE @@TRANCOUNT > 0
			BEGIN
				ROLLBACK TRANSACTION;
			END

		END CATCH`, useCorrectDB, targetTenantId, finishedJSON)
	//log.Printf("Target Active Directory SQL Migration Script: %s", sqlScript4)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript4))
	lol, err := db.Prepare(sqlScript4)
	if err != nil {
		log.Fatalf("Error preparing AD SQL script 4: %v", err)
	}
	result, err := lol.Exec()
	//result, err := db.Exec(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	//debugtown(fmt.Sprintf("RunSQLMigration: Results: %s", result))
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error fetching rows affected: %v", err)
	}

	log.Printf("AD SQL script 4 executed successfully. Rows affected: %d", rowsAffected)
	return nil

}
