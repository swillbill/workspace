package workflow

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

func RunTargetWorkflowSQL(targetEnv, targetTenantId string, targetUserName string, targetDBPassword string, data1 []CursorDataSection1, data2 []CursorDataSection2, data3 []CursorDataSection3) (err error) {

	log.Printf("Beginning Worlflow SQL Migration Script 1")

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
	file, err := os.Open("workflowsourcecursordatasection1.json")
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
		log.Fatal("workflowsourcecursordatasection1.json file is empty")
	}

	var unmarshaledData1 []CursorDataSection1
	err = json.Unmarshal(byteValue, &unmarshaledData1)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON1 := unmarshaledData1[0].CursorData_Section1
	for numbah, rows := range data1 {
		if numbah == 0 {
			finishedJSON1 = rows.CursorData_Section1
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON1 = fmt.Sprintf("%s,%s", finishedJSON1, rows.CursorData_Section1)
		}

	}

	// read artifacted file
	file, err = os.Open("workflowsourcecursordatasection2.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file contents
	byteValue, err = io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue) == 0 {
		log.Fatal("workflowsourcecursordatasection2.json file is empty")
	}

	var unmarshaledData2 []CursorDataSection2
	err = json.Unmarshal(byteValue, &unmarshaledData2)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON2 := unmarshaledData2[0].CursorData_Section2
	for numbah, rows := range data2 {
		if numbah == 0 {
			finishedJSON2 = rows.CursorData_Section2
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON2 = fmt.Sprintf("%s,%s", finishedJSON2, rows.CursorData_Section2)
		}

	}

	// read artifacted file
	file, err = os.Open("workflowsourcecursordatasection3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file contents
	byteValue, err = io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue) == 0 {
		log.Fatal("workflowsourcecursordatasection3.json file is empty")
	}

	var unmarshaledData3 []CursorDataSection3
	err = json.Unmarshal(byteValue, &unmarshaledData3)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON3 := unmarshaledData3[0].CursorData_Section3
	for numbah, rows := range data3 {
		if numbah == 0 {
			finishedJSON3 = rows.CursorData_Section3
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON3 = fmt.Sprintf("%s,%s", finishedJSON3, rows.CursorData_Section3)
		}

	}

	var useCorrectDB = "WorkflowEngineDB"
	// get sourceEnv and add as a prefix to the DB
	if targetEnv != "" {
		if strings.ToLower(targetEnv) == "ref" {
			useCorrectDB = targetEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "WorkflowEngineDB"
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		}
	}

	sqlScript := fmt.Sprintf(`
USE [%s];
DECLARE      @TenantId       uniqueidentifier = '%s';

SET NOCOUNT ON
SET XACT_ABORT ON;

DECLARE @Code                nvarchar(150)                 
,       @Description         nvarchar(250)                 
,       @CreateDate          datetime                       = getutcdate()
,       @CreatedBy           nvarchar(450)                  = 'AdminUser'
,       @UpdatedDttm         datetime                       = getutcdate()
,       @UpdatedBy           nvarchar(450)                  = 'AdminUser'
,       @newid               uniqueidentifier              
,       @debug               bit -- set to 1 for dubugging 
,       @commit              bit -- set to 1 for committing
,       @err_message         nvarchar(1000)                
--,       @ActiveFlag     bit null
,       @wfvCode             nvarchar(50)                  
,       @wfgCode             nvarchar(50)                  
,       @wgDescription       nvarchar(250)                 
,       @wvDescription       nvarchar(250)                 
,       @WfgShortDescription nvarchar(150)                 
,       @wfvShortDescription nvarchar(150)                 
,       @wvModule            nvarchar(100)                 
,       @workflowGroupId     int                           
,       @WorkflowVariantId   int                           
,       @PrefixId            varchar(10)                   
,       @CorrelationId  uniqueidentifier 
,       @MVModuleId        int
,       @MGModuleId        int

print '-----------------------------------------------------------'
print  'Tenant ID: ' + cast (@tenantId as nvarchar(100))
print '-----------------------------------------------------------'

--IF N'$(varDebug)' = 1 SET @debug = 1 ELSE SET @debug = 0;
--IF N'$(varCommit)' = 1 SET @commit = 1 ELSE SET @commit = 0;
SET @debug = 1; --For testing 
SET @commit = 1; --For testing 


BEGIN TRY
    BEGIN TRANSACTION
------ Section 1: Insert ------
DECLARE curList CURSOR FOR SELECT *
FROM (
VALUES %s) v ([wfvCode],[wfgCode], [wfgDescription],[wfgShortDescription],[wfvDescription],[wfvShortDescription],[Module],[prefixId],[MVmoduleId],[MGModuleId]);

OPEN curList
FETCH NEXT FROM curList
INTO @wfvCode ,@wfgCode, @wgDescription,@wfgShortDescription, @wvDescription,@wfvShortDescription, @wvModule,@prefixId, @MvModuleid, @MGModuleId

WHILE @@FETCH_STATUS = 0
BEGIN


-------Insert into [WorkflowGroups] table

	IF EXISTS (SELECT 1 
		FROM [WorkflowGroups]
		WHERE Code= @wfgCode
			AND TenantId = @TenantId
			and Moduleid = @MGModuleId
		)
		begin
			IF @debug = 1 PRINT 'WorkflowGroups already exists - Name: '+@wfgCode+' '	
        end

	IF NOT EXISTS (SELECT 1
		FROM [WorkflowGroups]
		WHERE Code= @wfgCode
			AND TenantId= @TenantId
			and Moduleid = @MGModuleId
		)
	BEGIN

	   SET @CorrelationId = newid()

		INSERT INTO [WorkflowGroups] ( [Code],   [TenantId], [CreatedDttm], [CreatedBy], [UpdatedDttm], [UpdatedBy], [LongDescription], [ShortDescription],   [PrefixId], [CorrelationId],[ModuleId] )
		VALUES                       ( @wfgCode, @TenantId,  @CreateDate,   @CreatedBy,  @UpdatedDttm,  @UpdatedBy,  @wgDescription,    @WfgShortDescription, @prefixId ,@CorrelationId, @MGModuleId );
		
		IF @debug = 1 PRINT 'WorkflowGroups inserted successfully - Name: '+@wfgCode+' '
	END
	

---insert into  table [WorkflowVariants]


    IF  EXISTS (SELECT 1
		  FROM [WorkflowVariants]
		  WHERE TenantId= @tenantId
			AND Code = @wfvCode
			and moduleid = @MvModuleid)
			
			begin 
			
			IF @debug = 1 PRINT 'Workflowvariant/code  already exists :- '+@wfvCode+''
			print '----------------------------------'

			end

 IF  EXISTS (SELECT 1 FROM [WorkflowVariants] WHERE TenantId= @tenantId AND Code = @wfvCode  )
 and  not exists (SELECT 1  FROM [WorkflowVariants] WHERE TenantId= @tenantId AND Code = @wfvCode and Moduleid = @MvModuleid )
		
		BEGIN

			update [WorkflowVariants]
			set Moduleid = @MvModuleid
			,[UpdatedDttm]=  @UpdatedDttm
			,[UpdatedBy]= @UpdatedBy
		    WHERE TenantId= @tenantId
			AND Code = @wfvCode 
			--and Moduleid = @MvModuleid

			IF @debug = 1 PRINT 'Workflowvariant/code  Moduleid is updated :- '+@wfvCode+''
	    END
	
	IF Exists (SELECT 1  FROM [WorkflowVariants] WHERE TenantId= @tenantId  AND Code = @wfvCode and Moduleid = @MvModuleid AND [CorrelationId] is NULL)
	
	BEGIN

    SET @CorrelationId = newid()
	update [WorkflowVariants]
			set [CorrelationId] = @CorrelationId
			,[UpdatedDttm]=  @UpdatedDttm
			,[UpdatedBy]= @UpdatedBy
		    WHERE TenantId= @tenantId
			AND Code = @wfvCode 
			and Moduleid = @MvModuleid

			IF @debug = 1 PRINT 'Workflowvariant CorrelationId is updated :- '+@wfvCode+''
			print '----------------------------------'
	END
	
			
	IF NOT EXISTS (SELECT 1
		FROM [WorkflowVariants]
		WHERE TenantId= @tenantId
			AND Code = @wfvCode
			and Moduleid = @MvModuleid )	
	
		BEGIN
   
   SET @CorrelationId = newid()
   SET @workflowGroupId =(SELECT [WorkflowGroupId] FROM [WorkflowGroups] WHERE code =@wfgCode AND TenantId =@TenantId)
           
		   

			INSERT INTO [WorkflowVariants] ( [Code],   [LongDescription], [TenantId], [CreatedDttm], [CreatedBy], [UpdatedDttm], [UpdatedBy], [WorkflowGroupId], [ShortDescription],   [Module] , [CorrelationId] ,[ModuleId] )
			VALUES                         ( @wfvCode, @wvDescription,    @tenantId,  @CreateDate,   @CreatedBy,  @UpdatedDttm,  @UpdatedBy,  @workflowGroupId,  @wfvShortDescription, @wvModule,@CorrelationId,@MvModuleid  )
			IF @debug = 1
			
			PRINT 'WorkflowVaraint Code inserted successfully - Name: '+@wfvCode+''
			print '----------------------------------'
		END

	FETCH NEXT FROM curList
	INTO @wfvCode ,@wfgCode, @wgDescription,@wfgShortDescription, @wvDescription,@wfvShortDescription, @wvModule,@prefixId, @MvModuleid, @MGModuleId
END

CLOSE curList;
DEALLOCATE curList;

-----------------------------------------------------------------------------------------------
--Insert into [Queue], [QueueRole] tables

/*
AMazur: 
	1. added administrator role variable (gets derived from System Settings)
	2. added insert into [QueueRole]
*/
-----------------------------------------------------------------------------------------------

DECLARE @QName             nvarchar(50) 
,       @QStatusId         int          
,       @Qcode             varchar(50)  
,       @QLongDescription  nvarchar(256)
,       @QShortDescription nvarchar(100)
,       @QStartDate        datetime      = '2023-01-01 00:00:00.000'
,       @QEndDate          datetime      = '9999-12-31 23:59:59.997'
,		@IdentityVal	   int
,       @AdminRole         nvarchar(150) 

SET @AdminRole = (SELECT isnull(max(SettingValue),'Administrator') FROM SystemSettings s
				  WHERE s.SettingKey='QueueAdminRoleName'
				  --AND s.TenantId=@TenantId !!! column does not exist yet
				 )

------ Section 2: Insert ------
DECLARE curList CURSOR FOR SELECT *
FROM (
VALUES %s) v ([QName],[QLongDescription],[QShortDescription],[Qcode]);


OPEN curList
FETCH NEXT FROM curList
INTO @QName,@QLongDescription, @QShortDescription,@Qcode;

WHILE @@FETCH_STATUS = 0
BEGIN

	SET @QStatusId = (SELECT [QueueStatusId] FROM [QueueStatus] WHERE [Code] = @Qcode)

    IF EXISTS (SELECT 1
		FROM [Queue]
		WHERE TenantId= @tenantId
			AND Name = @QName )
	
       IF @debug = 1 PRINT 'Queue Name already exists :- '+@QName+''

    IF NOT EXISTS (SELECT 1
		FROM [Queue]
		WHERE TenantId= @tenantId
			AND Name = @QName )
		BEGIN
  
			INSERT INTO [Queue] ( [Name], [QueueStatusId], [LongDescription], [ShortDescription], [StartDate], [EndDate], [CreatedDttm], [CreatedBy], [UpdatedDttm], [UpdatedBy], [TenantId] )
			VALUES              ( @QName, @QStatusId,      @QLongDescription, @QShortDescription, @QStartDate, @QEndDate, @CreateDate,   @CreatedBy,  @UpdatedDttm,  @UpdatedBy,  @TenantId  );

			SET @IdentityVal = scope_identity()

			INSERT INTO [QueueRole] ( Role,       QueueId,      CreatedDttm, CreatedBy,  UpdatedDttm,  UpdatedBy,  TenantId  )
			VALUES                  ( @AdminRole, @IdentityVal, @CreateDate, @CreatedBy, @UpdatedDttm, @UpdatedBy, @TenantId );

			IF @debug = 1 PRINT 'Queue Name inserted successfully - Name: '+@QName+''
		END
	FETCH NEXT FROM curList
	INTO @QName,@QLongDescription, @QShortDescription,@Qcode;
END

CLOSE curList;
DEALLOCATE curList;

-------------------------------------------------------------------------------------------------------
---Insert into [dbo].[WorkflowVariantObjects] table

DECLARE @VariantObject nvarchar(MAX)
,       @VOtypeOf      nvarchar(500) =''
,       @VOvariantId   bigint       
,       @VOQueueId     bigint       
,       @VOSchemaId    bigint       
,       @VOActiveFlag  bit          
,       @VODisplayName  varchar(256) 
,       @VOQname       varchar(50)  
,       @VOcode        varchar(50)  
,       @VOModuleId    int

------ Section 3: Insert ------
DECLARE curList CURSOR FOR SELECT *
FROM (
VALUES %s) v ([VariantObject],[VOQname],[VODisplayName],[VOcode],[VOActiveFlag],[VOModuleId]);

---insert into WorkflowVariantObjects table
OPEN curList
FETCH NEXT FROM curList
INTO @VariantObject,@VOQname,@VODisplayName,@VOcode,@VOActiveFlag, @VOModuleId;

WHILE @@FETCH_STATUS = 0
BEGIN

	SET @VOvariantId = (SELECT [WorkflowVariantId] FROM [WorkflowVariants] WHERE [Code] = @VOcode AND TenantId =@TenantId /*and moduleid = @VOModuleId*/)
	SET @VOQueueId = (SELECT [QueueId] FROM [Queue] WHERE [Name] = @VOQname AND TenantId =@TenantId )
	SET @VOSchemaId = (SELECT [WorkflowSchemaId] FROM [WorkflowSchemas] WHERE [WorkflowSchemaType] = @VODisplayName )
	--AND TenantId =@TenantId) --AMazur added filter on the tenant ID

    IF EXISTS (SELECT 1
		FROM [WorkflowVariantObjects]
		WHERE TenantId= @tenantId
			AND [WorkflowVariantId] = @VOvariantId
			--AND [QueueId] = @VOQueueId
			--AND [WorkflowSchemaId] = @VOSchemaId
			--and [ModuleId] = @VOModuleId
			)
	
	begin
	    update [dbo].[WorkflowVariantObjects]
		set [WorkflowVariantId] = @VOvariantId
		    ,[VariantObject]= @VariantObject
			, [QueueId] = @VOQueueId
			, [WorkflowSchemaId] = @VOSchemaId
			,[UpdatedBy] = @UpdatedBy
            ,[UpdatedDttm] = @UpdatedDttm
			,[ModuleId] = @VOModuleId
			,[CorrelationId] = @CorrelationId
        WHERE TenantId= @tenantId
			AND [WorkflowVariantId] = @VOvariantId
			--AND [QueueId] = @VOQueueId
		    --AND [WorkflowSchemaId] = @VOSchemaId
		    --and [ModuleId] = @VOModuleId

       IF @debug = 1 PRINT 'Variant Object updated  for:- '+@VOcode+','+@VOQname+','+@VODisplayName+' '

	end


    IF NOT EXISTS (SELECT 1
		FROM [dbo].[WorkflowVariantObjects]
		WHERE TenantId= @tenantId
			AND [WorkflowVariantId] = @VOvariantId
			AND [QueueId] = @VOQueueId
			AND [WorkflowSchemaId] = @VOSchemaId
			and [ModuleId] = @VOModuleId)

		BEGIN
  SET @CorrelationId = newid()

			INSERT INTO [dbo].[WorkflowVariantObjects] ( [TypeOf],  [VariantObject], [WorkflowVariantId], [TenantId], [CreatedDttm], [CreatedBy], [UpdatedDttm], [UpdatedBy], [QueueId],  [WorkflowSchemaId], [ActiveFlag] , [CorrelationId], [ModuleId] )
			VALUES                                     ( @VOtypeOf, @VariantObject,  @VOvariantId,        @tenantId,  @CreateDate,   @CreatedBy,  @UpdatedDttm,  @UpdatedBy,  @VOQueueId, @VOSchemaId, @VOActiveFlag ,@CorrelationId,  @VOModuleId );

			IF @debug = 1 PRINT 'Variant Object inserted successfully for:- '+@VOcode+','+@VOQname+','+@VODisplayName+' '
		END

	FETCH NEXT FROM curList
	INTO @VariantObject,@VOQname,@VODisplayName,@VOcode,@VOActiveFlag, @VOModuleId;
END

CLOSE curList;
DEALLOCATE curList;

/*Test*/

select count(*) as WorkflowGroups from [dbo].[WorkflowGroups] where tenantid = @TenantId  and [CreatedDttm]= @CreateDate
select count(*) as WorkflowVariants from [WorkflowVariants] where tenantid = @TenantId and [CreatedDttm]= @CreateDate
select count(*)  as WorkflowVariantObjects from [dbo].[WorkflowVariantObjects] where tenantid = @TenantId and [CreatedDttm]= @CreateDate
----------------------------------------------------------------------------------------------

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

END CATCH`, useCorrectDB, targetTenantId, finishedJSON1, finishedJSON2, finishedJSON3)

	//log.Printf("Target Active Directory SQL Migration Script: %s", sqlScript1)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	lol, err := db.Prepare(sqlScript)
	if err != nil {
		log.Fatalf("Error preparing Workflow Target SQL script: %v", err)
	}
	result, err := lol.Exec()
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	//debugtown(fmt.Sprintf("RunSQLMigration: Results: %s", result))
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error fetching rows affected: %v", err)
	}

	log.Printf("Target Workflow script executed successfully. Rows affected: %d", rowsAffected)
	return nil

}
