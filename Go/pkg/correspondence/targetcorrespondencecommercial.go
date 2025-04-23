package correspondence

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

func TargetCorrespondenceMigrationCommercialSQL(targetEnv string, sourceTenantID string, targetTenantId string, targetUserName string, targetDBPassword string, data []HeaderAndFooter) (err error) {
	log.Printf("Beginning Target Correspondence Commercial SQL Migration Script 1")

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
	file, err := os.Open("headerandfooter.json")
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
		log.Fatal("headerandfooter.json file is empty")
	}

	var unmarshaledData []HeaderAndFooter
	err = json.Unmarshal(byteValue, &unmarshaledData)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON := unmarshaledData[0].Header_footer_cursor
	for numbah, rows := range data {
		if numbah == 0 {
			finishedJSON = rows.Header_footer_cursor
		} else {
			finishedJSON = fmt.Sprintf("%s,%s", finishedJSON, rows.Header_footer_cursor)
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

	sqlScript := fmt.Sprintf(`
	USE [%s];
	DECLARE @TenantId uniqueidentifier = '%s';

	SET NOCOUNT ON
	SET XACT_ABORT ON;

	DECLARE @Code          nvarchar(150)                 
	,       @Description   nvarchar(250)                 
	,       @SourceTenantId uniqueidentifier              = N'%s'
	,       @CreateDate    datetime                       = getutcdate()
	,       @CreatedBy     nvarchar(450)                  = 'AdminUser'
	,       @UpdatedDttm   datetime                       = getutcdate()
	,       @UpdatedBy     nvarchar(450)                  = 'AdminUser'
	,       @newid         uniqueidentifier              
	,       @debug         bit -- set to 1 for dubugging 
	,       @commit        bit -- set to 1 for committing
	,       @err_message   nvarchar(1000)                
	--,       @ActiveFlag     bit null
	,       @CorrelationId uniqueidentifier               = NULL
	,       @SeqID         int    
	,       @ModuleId        int


	--IF N'$(varDebug)' = 1 SET @debug = 1 ELSE SET @debug = 0;
	--IF N'$(varCommit)' = 1 SET @commit = 1 ELSE SET @commit = 0;
	SET @debug = 1; --For testing 
	SET @commit = 1; --For testing 

	BEGIN TRY
		BEGIN TRANSACTION

	----------------Insert into [dbo].[ReusableContentType] and  [dbo].[ReusableContent] --------------------------------------------------
	---------------------------------------------------------------------------------------------------------------------------------------

	DECLARE @RCTContentTypeId   int             
	,       @RCTName            nvarchar(50)    
	--,       @RCTFooterDisplayId nvarchar(50)    
	--,       @RCTHeaderDisplayId nvarchar(50)    
	,       @RCTCategory        nvarchar(50)    
	,       @RCTCorrelationId   uniqueidentifier
	,       @RCOpenXml          nvarchar(MAX)  
	,       @RCStatus           nvarchar(25)    
	,       @RCDescription      nvarchar(250)    =NULL
	,       @RCDisplayID        nvarchar(50)  
	,       @RCTypeDisplayID        nvarchar(50) 
	,       @RCCorrelationId    uniqueidentifier
	,       @RCVersion          int 
	,       @contentversion    int
	,       @TargetDisplayId  nvarchar(50)

	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES 

	--Section 1: Data for ReusableContent and ReusableContentType

	-----------------------------------------------------------------------------------------------
	%s
	---------------------------------------------------------------------------------------------
	--End of Section 6: Data for ReusableContent

	) v ([Name],[ReusableContentTypeCategory],[Status],[Version],[Description],[ReusableContentTypeDisplayId],[ReusableContentDisplayId],[ModuleId],[ContentOpenXml] );

	OPEN curList
	FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId,@RCOpenXml

	WHILE @@FETCH_STATUS = 0
	BEGIN

		IF EXISTS (SELECT 1
			FROM [ReusableContentType]
			WHERE TenantId= @tenantId
				AND Name = @RCTName )
				
		IF @debug = 1 PRINT 'ReusableContentType  already exists :- '+@RCTName+''

		IF NOT EXISTS (SELECT 1
			FROM [ReusableContentType]
			WHERE TenantId= @tenantId
				AND Name = @RCTName )
				
		
	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  
			
		BEGIN

			INSERT INTO [ReusableContentType] ( [Name],   [ReusableContentTypeDisplayId], [ReusableContentTypeCategory], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [CorrelationId] ,moduleid  )
			VALUES                            ( @RCTName, @RCTypeDisplayID,   @RCTCategory,  @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantID,  @RCTCorrelationId, @ModuleId );

		--	UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'HDT' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT 'ReusableContentType inserted successfully - HEADER Name: '+@RCTName+'/'++@RCTypeDisplayID+'/'
		END


	------------------------------------------------------------------------------------------------------
	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  

	set @RCTContentTypeId = (select [ReusableContentTypeId] from [dbo].[ReusableContentType] where Name = @RCTName and [TenantId]=@tenantId)
		
	set @contentversion = isnull((select max(version) from [ReusableContent] where [TenantId]=@tenantId  AND [ReusableContentTypeId] =@RCTContentTypeId ),1)
	set @contentversion = (@contentversion + 1)
	set @TargetDisplayId = (select [ReusableContentDisplayId] from [dbo].[ReusableContent] RC join [ReusableContentType] RCT
															on RC.[ReusableContentTypeId]= RCT.[ReusableContentTypeId] Where RCT.Name = @RCTName
															AND RC.status= 'ACTIVE' and RC.tenantid = @tenantid)


	IF EXISTS (SELECT 1
			FROM [ReusableContent]
			WHERE [TenantId]= @tenantId
				AND [ReusableContentTypeId] =@RCTContentTypeId
				AND [ReusableContentDisplayId] = @TargetDisplayId
				
				)

		begin

				update [ReusableContent]
				set [ContentOpenXml] = ( select [ContentOpenXml] from [ReusableContent] where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion and status = 'Active'),
				--set [ContentOpenXml] = ( select [ContentOpenXml] from [TemplateExport2]..ReusableContent_copy1 where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion  and status = 'Active'),
				[Description]= @RCDescription
			,[UpdatedBy]= @UpdatedBy 
			,[UpdatedDate]=@UpdatedDttm
			-- , [Version]= @contentversion
				WHERE [TenantId]= @tenantId
				AND [ReusableContentTypeId] =@RCTContentTypeId
				AND [ReusableContentDisplayId] = @TargetDisplayId	
				AND status= 'ACTIVE'


		IF @debug = 1 PRINT 'ReusableContent  updated  :- '+@RCDisplayId+''
		Print '------------------------------------------------------------------'		
		end

	IF NOT  EXISTS (SELECT 1
			FROM [ReusableContent]
			WHERE [TenantId]= @tenantId
				AND [ReusableContentTypeId] =@RCTContentTypeId
				--AND [Version] = @RCVersion
				AND [ReusableContentDisplayId] = @TargetDisplayId
				)
				
		BEGIN


			INSERT INTO [dbo].[ReusableContent] ( [ContentOpenXml],               [Status],  [Version],  [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [Description],  [ReusableContentDisplayId], [ReusableContentTypeId], [CorrelationId], moduleId  )
			VALUES                   ( 
			( select [ContentOpenXml] from [ReusableContent] where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion and status = 'ACTIVE'),
			--( select [ContentOpenXml] from [TemplateExport2]..ReusableContent_copy1 where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion and status = 'ACTIVE'),
			@RCStatus, @contentversion, @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantID,  @RCDescription, @RCDisplayId,        @RCTContentTypeId,       @RCCorrelationId ,@ModuleId)

			Update [ReusableContent]
			set [Status] = 'INACTIVE'
			,updatedby = 'AdminUser'
			WHERE [TenantId]= @tenantId
			AND [ReusableContentTypeId] =@RCTContentTypeId
			and version != @contentversion

			IF @debug = 1 PRINT 'ReusableContent     inserted successfully - HEADER Name: '+@RCDisplayId+''
		END

		FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId,@RCOpenXml
	END

	CLOSE curList;
	DEALLOCATE curList;

	--------------------------------------------------------------------------------------------------------------------

	DECLARE 
		@RCTFooterDisplayId nvarchar(50)    
	,       @RCTHeaderDisplayId nvarchar(50)    


	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES 


	--Section 1: Data for ReusableContent and ReusableContentType

	-----------------------------------------------------------------------------------------------
	%s
	-----------------------------------------------------------------------------------------------
	--End of Section 9: Data for ReusableContent

	) v ([Name],[ReusableContentTypeCategory],[Status],[Version],[Description],[ReusableContentTypeDisplayId],[ReusableContentDisplayId],[ModuleID],[ContentOpenXml] );

	OPEN curList
	FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId, @RCOpenXml
	WHILE @@FETCH_STATUS = 0
	BEGIN

	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  


	IF EXISTS (SELECT 1
			FROM [ReusableContentType]
			WHERE TenantId= @tenantId
				AND Name = @RCTName )	
		Begin 
		IF (@RCTCategory = 'HEADER' )

		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'HDT' AND [TenantId]=@TenantId)
			SET @RCTHeaderDisplayId = 'HDT'+FORMAT(@SeqID, '000000000')

			Update  [ReusableContentType] 
			set [ReusableContentTypeDisplayId] = @RCTHeaderDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			,ModuleId=1
			WHERE TenantId= @tenantId
				AND Name = @RCTName 
				AND [ReusableContentTypeCategory]= @RCTCategory

			UPDATE [IdSequence] 
			SET [LastUsedNumber] = @SeqID 
			--Updatedby ='AdminUser'
			--Updateddate= getdate()
			WHERE [Prefix] = 'HDT' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContentType] updated successfully - HEADER Name: '+@RCTHeaderDisplayId+''
			Print +@RCTHeaderDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END

		IF (@RCTCategory = 'FOOTER' )

		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'FTT' AND [TenantId]=@TenantId)
			SET @RCTFooterDisplayId = 'FTT'+FORMAT(@SeqID, '000000000')
			
			Update  [ReusableContentType] 
			set [ReusableContentTypeDisplayId] = @RCTFooterDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			,ModuleId=@moduleid
			WHERE TenantId= @tenantId
				AND Name = @RCTName 
				AND [ReusableContentTypeCategory]= @RCTCategory
			
			UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID 
			--Updatedby ='AdminUser'
			--Updateddate= getdate()
			WHERE [Prefix] = 'FTT' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContentType] updated successfully - FOOTER Name: '+@RCTFooterDisplayId+''
			Print +@RCTFooterDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END
	END


	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  

	set @RCTContentTypeId = (select [ReusableContentTypeId] from [dbo].[ReusableContentType] where Name = @RCTName and [TenantId]=@tenantId)
	
	IF EXISTS (SELECT 1
			FROM [ReusableContent]
			WHERE [TenantId]= @tenantId
				)
	Begin

		IF (@RCTCategory = 'Header' )
		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'HDR' AND [TenantId]=@TenantId)
			SET @RCTHeaderDisplayId = 'HDR'+FORMAT(@SeqID, '000000000')

			Update  [ReusableContent] 
			set [ReusableContentDisplayId] = @RCTHeaderDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			WHERE TenantId= @tenantId
			AND [ReusableContentTypeId] = @RCTContentTypeId
			AND @RCTCategory = 'Header'

			UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'HDR' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContent]     updated successfully - HEADER Name: '+@RCTName+''
			Print +@RCTHeaderDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END

		IF (@RCTCategory = 'FOOTER' )
		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'FTR' AND [TenantId]=@TenantId)
			SET @RCTFooterDisplayId = 'FTR'+FORMAT(@SeqID, '000000000')

			Update  [ReusableContent] 
			set [ReusableContentDisplayId] = @RCTFooterDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			WHERE TenantId= @tenantId
			AND [ReusableContentTypeId] = @RCTContentTypeId
			AND @RCTCategory = 'Header'

			UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'FTR' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContent]     updated successfully - FOOTER Name: '+@RCTName+''
			Print +@RCTFooterDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END

	END

		FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId, @RCOpenXml
	END

	CLOSE curList;
	DEALLOCATE curList;

	--Template cache cleanup--

	delete from [TemplateCache] where tenantid = @TenantId
	IF @debug = 1
						PRINT '[TemplateCache] truncated - : ' ++cast(@TenantId as varchar(150)) 

	-----------------------------------------------------------------
		--Tenant Summary--

		print '-----------------------------------------------------------' 
		print 'Counts for the tenant: '+cast(@TenantId as varchar(150)) 
		print '-----------------------------------------------------------'

		select count(*) as [ReusableContent] from [dbo].[ReusableContent] where TenantId=@TenantId
		select count(*) as [ReusableContentType] from [dbo].[ReusableContentType] where TenantId=@TenantId

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

	END CATCH

`, useCorrectDB, targetTenantId, sourceTenantID, finishedJSON, finishedJSON)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	lol, err := db.Prepare(sqlScript)
	if err != nil {
		log.Fatalf("Error preparing Target Correspondence script: %v", err)
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

	log.Printf("Target Correspondence script executed successfully. Rows affected: %d", rowsAffected)
	return nil

}

func TargetCorrespondenceMigrationGOVSQL(targetEnv string, sourceTenantID string, targetTenantId string, targetUserName string, targetDBPassword string, data []HeaderAndFooter) (err error) {
	log.Printf("Beginning Target Correspondence Gov SQL Migration Script")

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
	file, err := os.Open("headerandfooter.json")
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
		log.Fatal("headerandfooter.json file is empty")
	}

	var unmarshaledData []HeaderAndFooter
	err = json.Unmarshal(byteValue, &unmarshaledData)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON := unmarshaledData[0].Header_footer_cursor
	for numbah, rows := range data {
		if numbah == 0 {
			finishedJSON = rows.Header_footer_cursor
		} else {
			finishedJSON = fmt.Sprintf("%s,%s", finishedJSON, rows.Header_footer_cursor)
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

	sqlScript := fmt.Sprintf(`
	USE [%s];
	DECLARE @TenantId uniqueidentifier = '%s';

	SET NOCOUNT ON
	SET XACT_ABORT ON;

	DECLARE @Code          nvarchar(150)                 
	,       @Description   nvarchar(250)                 
	,       @SourceTenantId uniqueidentifier              = N'%s'   
	,       @CreateDate    datetime                       = getutcdate()
	,       @CreatedBy     nvarchar(450)                  = 'AdminUser'
	,       @UpdatedDttm   datetime                       = getutcdate()
	,       @UpdatedBy     nvarchar(450)                  = 'AdminUser'
	,       @newid         uniqueidentifier              
	,       @debug         bit -- set to 1 for dubugging 
	,       @commit        bit -- set to 1 for committing
	,       @err_message   nvarchar(1000)                
	--,       @ActiveFlag     bit null
	,       @CorrelationId uniqueidentifier               = NULL
	,       @SeqID         int    
	,       @ModuleId        int


	--IF N'$(varDebug)' = 1 SET @debug = 1 ELSE SET @debug = 0;
	--IF N'$(varCommit)' = 1 SET @commit = 1 ELSE SET @commit = 0;
	SET @debug = 1; --For testing 
	SET @commit = 0; --For testing 

	BEGIN TRY
		BEGIN TRANSACTION

	----------------Insert into [dbo].[ReusableContentType] and  [dbo].[ReusableContent] --------------------------------------------------
	---------------------------------------------------------------------------------------------------------------------------------------

	DECLARE @RCTContentTypeId   int             
	,       @RCTName            nvarchar(50)    
	--,       @RCTFooterDisplayId nvarchar(50)    
	--,       @RCTHeaderDisplayId nvarchar(50)    
	,       @RCTCategory        nvarchar(50)    
	,       @RCTCorrelationId   uniqueidentifier
	,       @RCOpenXml          nvarchar(MAX)  
	,       @RCStatus           nvarchar(25)    
	,       @RCDescription      nvarchar(250)    =NULL
	,       @RCDisplayID        nvarchar(50)  
	,       @RCTypeDisplayID        nvarchar(50) 
	,       @RCCorrelationId    uniqueidentifier
	,       @RCVersion          int 
	,       @contentversion    int
	,       @TargetDisplayId  nvarchar(50)

	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES 

	--Section 1: Data for ReusableContent and ReusableContentType

	-----------------------------------------------------------------------------------------------
	%s
	---------------------------------------------------------------------------------------------
	--End of Section 6: Data for ReusableContent

	) v ([Name],[ReusableContentTypeCategory],[Status],[Version],[Description],[ReusableContentTypeDisplayId],[ReusableContentDisplayId],[ModuleId],[ContentOpenXml] );

	OPEN curList
	FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId,@RCOpenXml

	WHILE @@FETCH_STATUS = 0
	BEGIN

		IF EXISTS (SELECT 1
			FROM [ReusableContentType]
			WHERE TenantId= @tenantId
				AND Name = @RCTName )
				
		IF @debug = 1 PRINT 'ReusableContentType  already exists :- '+@RCTName+''

		IF NOT EXISTS (SELECT 1
			FROM [ReusableContentType]
			WHERE TenantId= @tenantId
				AND Name = @RCTName )
				
		
	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  
			
		BEGIN

			INSERT INTO [ReusableContentType] ( [Name],   [ReusableContentTypeDisplayId], [ReusableContentTypeCategory], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [CorrelationId] ,moduleid  )
			VALUES                            ( @RCTName, @RCTypeDisplayID,   @RCTCategory,  @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantID,  @RCTCorrelationId, @ModuleId );

		--	UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'HDT' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT 'ReusableContentType inserted successfully - HEADER Name: '+@RCTName+'/'++@RCTypeDisplayID+'/'
		END


	------------------------------------------------------------------------------------------------------
	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  

	set @RCTContentTypeId = (select [ReusableContentTypeId] from [dbo].[ReusableContentType] where Name = @RCTName and [TenantId]=@tenantId)
		
	set @contentversion = isnull((select max(version) from [ReusableContent] where [TenantId]=@tenantId  AND [ReusableContentTypeId] =@RCTContentTypeId ),1)
	set @contentversion = (@contentversion + 1)
	set @TargetDisplayId = (select [ReusableContentDisplayId] from [dbo].[ReusableContent] RC join [ReusableContentType] RCT
															on RC.[ReusableContentTypeId]= RCT.[ReusableContentTypeId] Where RCT.Name = @RCTName
															AND RC.status= 'ACTIVE' and RC.tenantid = @tenantid)


	IF EXISTS (SELECT 1
			FROM [ReusableContent]
			WHERE [TenantId]= @tenantId
				AND [ReusableContentTypeId] =@RCTContentTypeId
				AND [ReusableContentDisplayId] = @TargetDisplayId
				
				)

		begin

				update [ReusableContent]
				set [ContentOpenXml] = ( select [ContentOpenXml] from [ReusableContent] where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion and status = 'Active'),
				--set [ContentOpenXml] = ( select [ContentOpenXml] from [TemplateExport2]..ReusableContent_copy1 where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion  and status = 'Active'),
				[Description]= @RCDescription
			,[UpdatedBy]= @UpdatedBy 
			,[UpdatedDate]=@UpdatedDttm
			-- , [Version]= @contentversion
				WHERE [TenantId]= @tenantId
				AND [ReusableContentTypeId] =@RCTContentTypeId
				AND [ReusableContentDisplayId] = @TargetDisplayId	
				AND status= 'ACTIVE'


		IF @debug = 1 PRINT 'ReusableContent  updated  :- '+@RCDisplayId+''
		Print '------------------------------------------------------------------'		
		end

	IF NOT  EXISTS (SELECT 1
			FROM [ReusableContent]
			WHERE [TenantId]= @tenantId
				AND [ReusableContentTypeId] =@RCTContentTypeId
				--AND [Version] = @RCVersion
				AND [ReusableContentDisplayId] = @TargetDisplayId
				)
				
		BEGIN


			INSERT INTO [dbo].[ReusableContent] ( [ContentOpenXml],               [Status],  [Version],  [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [Description],  [ReusableContentDisplayId], [ReusableContentTypeId], [CorrelationId], moduleId  )
			VALUES                   ( 
			( select [ContentOpenXml] from [ReusableContent] where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion and status = 'ACTIVE'),
			--( select [ContentOpenXml] from [TemplateExport2]..ReusableContent_copy1 where [TenantId] = @SourceTenantId and [ReusableContentDisplayId]= @RCDisplayId and [Version] = @RCVersion and status = 'ACTIVE'),
			@RCStatus, @contentversion, @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantID,  @RCDescription, @RCDisplayId,        @RCTContentTypeId,       @RCCorrelationId ,@ModuleId)

			Update [ReusableContent]
			set [Status] = 'INACTIVE'
			,updatedby = 'AdminUser'
			WHERE [TenantId]= @tenantId
			AND [ReusableContentTypeId] =@RCTContentTypeId
			and version != @contentversion

			IF @debug = 1 PRINT 'ReusableContent     inserted successfully - HEADER Name: '+@RCDisplayId+''
		END

		FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId,@RCOpenXml
	END

	CLOSE curList;
	DEALLOCATE curList;

	--------------------------------------------------------------------------------------------------------------------

	DECLARE 
		@RCTFooterDisplayId nvarchar(50)    
	,       @RCTHeaderDisplayId nvarchar(50)    


	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES 


	--Section 1: Data for ReusableContent and ReusableContentType

	-----------------------------------------------------------------------------------------------
	%s
	-----------------------------------------------------------------------------------------------
	--End of Section 9: Data for ReusableContent

	) v ([Name],[ReusableContentTypeCategory],[Status],[Version],[Description],[ReusableContentTypeDisplayId],[ReusableContentDisplayId],[ModuleID],[ContentOpenXml] );

	OPEN curList
	FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId, @RCOpenXml
	WHILE @@FETCH_STATUS = 0
	BEGIN

	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  


	IF EXISTS (SELECT 1
			FROM [ReusableContentType]
			WHERE TenantId= @tenantId
				AND Name = @RCTName )	
		Begin 
		IF (@RCTCategory = 'HEADER' )

		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'HDT' AND [TenantId]=@TenantId)
			SET @RCTHeaderDisplayId = 'HDT'+FORMAT(@SeqID, '000000000')

			Update  [ReusableContentType] 
			set [ReusableContentTypeDisplayId] = @RCTHeaderDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			,ModuleId=1
			WHERE TenantId= @tenantId
				AND Name = @RCTName 
				AND [ReusableContentTypeCategory]= @RCTCategory

			UPDATE [IdSequence] 
			SET [LastUsedNumber] = @SeqID 
			--Updatedby ='AdminUser'
			--Updateddate= getdate()
			WHERE [Prefix] = 'HDT' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContentType] updated successfully - HEADER Name: '+@RCTHeaderDisplayId+''
			Print +@RCTHeaderDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END

		IF (@RCTCategory = 'FOOTER' )

		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'FTT' AND [TenantId]=@TenantId)
			SET @RCTFooterDisplayId = 'FTT'+FORMAT(@SeqID, '000000000')
			
			Update  [ReusableContentType] 
			set [ReusableContentTypeDisplayId] = @RCTFooterDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			,ModuleId=@moduleid
			WHERE TenantId= @tenantId
				AND Name = @RCTName 
				AND [ReusableContentTypeCategory]= @RCTCategory
			
			UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID 
			--Updatedby ='AdminUser'
			--Updateddate= getdate()
			WHERE [Prefix] = 'FTT' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContentType] updated successfully - FOOTER Name: '+@RCTFooterDisplayId+''
			Print +@RCTFooterDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END
	END


	--*FTT and HDT are for reusable content types
	--*FTR and HDR are for reusable contents  

	set @RCTContentTypeId = (select [ReusableContentTypeId] from [dbo].[ReusableContentType] where Name = @RCTName and [TenantId]=@tenantId)
	
	IF EXISTS (SELECT 1
			FROM [ReusableContent]
			WHERE [TenantId]= @tenantId
				)
	Begin

		IF (@RCTCategory = 'Header' )
		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'HDR' AND [TenantId]=@TenantId)
			SET @RCTHeaderDisplayId = 'HDR'+FORMAT(@SeqID, '000000000')

			Update  [ReusableContent] 
			set [ReusableContentDisplayId] = @RCTHeaderDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			WHERE TenantId= @tenantId
			AND [ReusableContentTypeId] = @RCTContentTypeId
			AND @RCTCategory = 'Header'

			UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'HDR' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContent]     updated successfully - HEADER Name: '+@RCTName+''
			Print +@RCTHeaderDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END

		IF (@RCTCategory = 'FOOTER' )
		BEGIN
			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'FTR' AND [TenantId]=@TenantId)
			SET @RCTFooterDisplayId = 'FTR'+FORMAT(@SeqID, '000000000')

			Update  [ReusableContent] 
			set [ReusableContentDisplayId] = @RCTFooterDisplayId
			,[UpdatedBy] = @UpdatedBy
			,[UpdatedDate]= @UpdatedDttm
			WHERE TenantId= @tenantId
			AND [ReusableContentTypeId] = @RCTContentTypeId
			AND @RCTCategory = 'Header'

			UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'FTR' AND [TenantId]=@TenantId

			IF @debug = 1 PRINT '[ReusableContent]     updated successfully - FOOTER Name: '+@RCTName+''
			Print +@RCTFooterDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''
		END

	END

		FETCH NEXT FROM curList
	INTO    @RCTName ,@RCTCategory  ,@RCStatus,@RCVersion,@RCDescription ,@RCTypeDisplayID , @RCDisplayID ,@ModuleId, @RCOpenXml
	END

	CLOSE curList;
	DEALLOCATE curList;

	--Template cache cleanup--

	delete from [TemplateCache] where tenantid = @TenantId
	IF @debug = 1
						PRINT '[TemplateCache] truncated - : ' ++cast(@TenantId as varchar(150)) 

	-----------------------------------------------------------------
		--Tenant Summary--

		print '-----------------------------------------------------------' 
		print 'Counts for the tenant: '+cast(@TenantId as varchar(150)) 
		print '-----------------------------------------------------------'

		select count(*) as [ReusableContent] from [dbo].[ReusableContent] where TenantId=@TenantId
		select count(*) as [ReusableContentType] from [dbo].[ReusableContentType] where TenantId=@TenantId

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

	END CATCH`, useCorrectDB, targetTenantId, sourceTenantID, finishedJSON, finishedJSON)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	lol, err := db.Prepare(sqlScript)
	if err != nil {
		log.Fatalf("Error preparing Target Correspondence script: %v", err)
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

	log.Printf("Target Correspondence script executed successfully. Rows affected: %d", rowsAffected)
	return nil

}
