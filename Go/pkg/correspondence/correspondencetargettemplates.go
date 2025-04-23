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

func TargetCorrespondenceMigrationCommercialSQLTemplate(targetEnv string, sourceTenantID string, targetTenantId string, targetUserName string, targetDBPassword string, data1 []CursorData1, data2 []CursorData2, data3 []CursorData3, data4 []CursorData4) (err error) {
	log.Printf("Beginning Target Correspondence Commercial SQL Migration Script 1")

	// Get the SQL configuration for the target environment
	sqlConfig, ok := formsmigration.GetMigrationSQLConfig(targetEnv, targetUserName, targetDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", targetEnv)
		return fmt.Errorf("invalid environment: %s", targetEnv)
	}

	// Connect to the SQL Server
	db, err := sql.Open("sqlserver", sqlConfig.ConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL Server: %v", err)
	}
	defer db.Close()

	// read artifacted file
	file1, err := os.Open("corrsourcecursordatasection1.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()

	// read file contents
	byteValue1, err := io.ReadAll(file1)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue1) == 0 {
		log.Fatal("corrsourcecursordatasection1.json file is empty")
	}

	var unmarshaledData1 []CursorData1
	err = json.Unmarshal(byteValue1, &unmarshaledData1)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON1 := unmarshaledData1[0].CursorData1
	for numbah, rows := range data1 {

		escapedCursorData1 := rows.CursorData1

		// Check for occurrences of "don't" and replace with "dont''t"
		escapedCursorData1 = strings.Replace(escapedCursorData1, "Don't", "Dont''t", -1)

		if numbah == 0 {
			finishedJSON1 = escapedCursorData1
		} else {
			finishedJSON1 = fmt.Sprintf("%s,%s", finishedJSON1, escapedCursorData1)
		}
	}

	// read artifacted file
	file2, err := os.Open("corrsourcecursordatasection2.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// read file contents
	byteValue2, err := io.ReadAll(file2)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue2) == 0 {
		log.Fatal("corrsourcecursordatasection2.json file is empty")
	}

	var unmarshaledData2 []CursorData2
	err = json.Unmarshal(byteValue2, &unmarshaledData2)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	// Process data and escape single quotes, and handle "don't" -> "dont''t"
	finishedJSON2 := unmarshaledData2[0].CursorData2
	for numbah, rows := range data2 {
		// Escape single quotes in each row's CursorData1
		escapedCursorData2 := rows.CursorData2

		// Check for occurrences of "don't" and replace with "dont''t"
		escapedCursorData2 = strings.Replace(escapedCursorData2, "Don't", "Dont''t", -1)

		if numbah == 0 {
			finishedJSON2 = escapedCursorData2
		} else {
			finishedJSON2 = fmt.Sprintf("%s,%s", finishedJSON2, escapedCursorData2)
		}
	}

	// read artifacted file
	file3, err := os.Open("corrsourcecursordatasection3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// read file contents
	byteValue3, err := io.ReadAll(file3)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue3) == 0 {
		log.Fatal("corrsourcecursordatasection3.json file is empty")
	}

	var unmarshaledData3 []CursorData3
	err = json.Unmarshal(byteValue2, &unmarshaledData3)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	// Process data and escape single quotes, and handle "don't" -> "dont''t"
	finishedJSON3 := unmarshaledData3[0].CursorData3
	for numbah, rows := range data3 {
		// Escape single quotes in each row's CursorData1
		escapedCursorData3 := rows.CursorData3

		// Check for occurrences of "don't" and replace with "dont''t"
		escapedCursorData3 = strings.Replace(escapedCursorData3, "Don't", "Dont''t", -1)

		if numbah == 0 {
			finishedJSON3 = escapedCursorData3
		} else {
			finishedJSON3 = fmt.Sprintf("%s,%s", finishedJSON3, escapedCursorData3)
		}
	}

	// read artifacted file
	file4, err := os.Open("corrsourcecursordatasection4.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// read file contents
	byteValue4, err := io.ReadAll(file4)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue4) == 0 {
		log.Fatal("corrsourcecursordatasection4.json file is empty")
	}

	var unmarshaledData4 []CursorData4
	err = json.Unmarshal(byteValue4, &unmarshaledData4)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON4 := unmarshaledData4[0].Cursor_for_Roles_section_4
	for numbah, rows := range data4 {

		escapedCursorData4 := rows.Cursor_for_Roles_section_4

		// Check for occurrences of "don't" and replace with "dont''t"
		escapedCursorData4 = strings.Replace(escapedCursorData4, "Don't", "Dont''t", -1)

		if numbah == 0 {
			finishedJSON4 = escapedCursorData4
		} else {
			finishedJSON4 = fmt.Sprintf("%s,%s", finishedJSON4, escapedCursorData4)
		}
	}
	// Determine the correct database to use
	var useCorrectDB = "NgCorrespondence"
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
	,       @ModuleId      int

	--IF N'$(varDebug)' = 1 SET @debug = 1 ELSE SET @debug = 0;
	--IF N'$(varCommit)' = 1 SET @commit = 1 ELSE SET @commit = 0;
	SET @debug = 1; --For testing 
	SET @commit = 1; --For testing 

	BEGIN TRY
		BEGIN TRANSACTION


	----Table 1. Insert into [dbo].[CorrespondenceType]--------------------------------------------------------------------------

	DECLARE @CRPTName              nvarchar(50)
	--,       @CRPTDescription nvarchar(200)
	,       @CRPTCategoryId        int
	,       @CRPTPrimaryIdType     varchar(50)
	,       @RowVersion            timestamp
	,       @CRPTDisplayId         nvarchar(50)
	,       @CRPTCertifiedMailFlag bit
	,       @CRPTContextLevel      nvarchar(MAX)
	,       @CRPTFTIFlag           bit
	,       @CRPTWorkflowSubTypeId nvarchar(50)
	,       @CRPTworkflowTypeId    nvarchar(50)
	,       @CRPTCName             nvarchar(50)
	,       @Printgroupname        nvarchar(50)
	,       @PrintgroupAssignmenttype nvarchar(50)
	,       @printgroupId int
	,       @corrtypeId   int
	,       @wfsubtype    int --nvarchar(50)
	,       @wftype       int  --nvarchar(50)

	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES 

	--Section 1: Data for [CorrespondenceType]
	-----------------------------------------------------------------------------------------------
	%s
	-----------------------------------------------------------------------------------------------
	--End of Section 1: Data for [dbo].[CorrespondenceType]

	) v ([Name],[CertifiedMailFlag],[ContextLevel],[FTIFlag],[WorkflowSubTypeId],[WorkflowTypeId],[CTCName],[PGName],[PGAssignmentType], [ModuleId],[CorrespondenceTypeDisplayId]);


	OPEN curList
	FETCH NEXT FROM curList
	INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@Printgroupname,@PrintgroupAssignmenttype,@ModuleId,@CRPTDisplayId

	WHILE @@FETCH_STATUS = 0
	BEGIN

	SET @CorrelationId = newid()
	SET @wfsubtype = (select [WorkflowVariants].[WorkflowVariantId] from [ref_WorkflowEngineDB]..[WorkflowVariants] where WorkflowVariants.code = @CRPTWorkflowSubTypeId and WorkflowVariants.tenantid= @tenantid and moduleid = @moduleid)

	--IF (@CRPTWorkflowSubTypeId is not null or  @wfsubtype is null)

	--Begin
	--PRINT 'Workflow Subtype "'+@wfsubtype+'" Does not exist for TenandID "'+ cast(@tenantId as varchar(50))+ '"!';
	--	THROW 50001, 'Workflow subtype does not exist.', 1;
	--END

	SET @wftype    = (select [WorkflowGroups].[WorkflowGroupId] from [ref_WorkflowEngineDB]..[WorkflowGroups] where [WorkflowGroups].code = @CRPTworkflowTypeId and [WorkflowGroups].tenantid= @tenantid)

	--IF (@CRPTworkflowTypeId is not null) and (@wftype is null)

	--Begin
	--PRINT 'Workflow type "'+@wftype+'" Does not exist for TenandID "'+ cast(@tenantId as varchar(50))+ '"!';
	--	THROW 50001, 'Workflow type does not exist.', 1;
	--END	

		IF EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and moduleid = @ModuleId)
		
			IF @debug = 1 PRINT 'CorrespondenceType already exists - Name: '+@CRPTName+' '	
			
		IF  EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and moduleid = @ModuleId)
		and  not exists (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and contextlevel = @CRPTContextLevel and [WorkflowSubTypeId]=@wfsubtype and [WorkflowTypeId]= @wftype  and moduleid = @ModuleId)
	
	begin
		update [CorrespondenceType]
		set [ContextLevel] = @CRPTContextLevel
		,[WorkflowSubTypeId] = @wfsubtype
		,[WorkflowTypeId] =    @wftype
		,[UpdatedBy] = @UpdatedBy
		,[UpdatedDate]= @UpdatedDttm
		,moduleid = @ModuleId
		where Name= @CRPTName 
		AND TenantId = @TenantId
		and moduleid = @ModuleId


		IF @debug = 1 PRINT 'CorrespondenceType Contectlevel/[WorkflowSubTypeId]/[WorkflowTypeId] are updated - Name: '+@CRPTName+' '	
		
	end

	IF exists (select 1 from [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and  [CorrelationId]= NULL)
	BEGIN
	
	SET @CorrelationId = newid()
	update [CorrespondenceType]
		set  [CorrelationId] = @CorrelationId
		,[UpdatedBy] = @UpdatedBy
		,[UpdatedDate]= @UpdatedDttm
			where Name= @CRPTName 
			AND TenantId = @TenantId
			and moduleid = @ModuleId

			IF @debug = 1 PRINT 'CorrespondenceType CorrelationId is updated - Name: '+@CRPTName+' '	
	END

	SET @CRPTCategoryId = (SELECT top 1 [CorrespondenceTypeCategoryId]
			FROM [dbo].[CorrespondenceTypeCategory]
			WHERE Name = @CRPTCName
				AND TenantId =@TenantId
				/*and moduleid = @ModuleId*/)
				
	---fix for correspondencetypecategoryid nulls-- 5/9/2024


	if exists(select 1 from [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and  [CorrespondenceTypeCategoryId]= NULL)
	
	
	begin
	update [CorrespondenceType]
		set  [CorrespondenceTypeCategoryId] = @CRPTCategoryId
		,[UpdatedBy] = @UpdatedBy
		,[UpdatedDate]= @UpdatedDttm
			where Name= @CRPTName 
			AND TenantId = @TenantId
			--and moduleid = @ModuleId

	IF @debug = 1 PRINT 'CorrespondenceType [CorrespondenceTypeCategoryId] is updated - Name: '+@CRPTName+' '	
		end	
	

		IF NOT EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId)

		BEGIN
		
		

			INSERT INTO [CorrespondenceType] ( [Name],    [CorrespondenceTypeCategoryId], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [CorrespondenceTypeDisplayId], [CertifiedMailFlag],    [ContextLevel],    [WorkflowSubTypeId], [WorkflowTypeId] ,  [FTIFlag]  , ModuleId  ,  [CorrelationId]  )
			VALUES                           ( @CRPTName, @CRPTCategoryId,                @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantId,  @CRPTDisplayId,                @CRPTCertifiedMailFlag,  @CRPTContextLevel, @wftype,            @wftype,            @CRPTFTIFlag , @ModuleId, @CorrelationId);

		

			IF @debug = 1 PRINT '[CorrespondenceType] inserted successfully - Name: '+@CRPTName+' '
		END

	---updated 9/26 --		
			--Insert into [dbo].[CorrespondenceTypePrintGroup]
			
			set @printgroupId = (select [PrintGroupId] from [dbo].[PrintGroup] where [Name] = @Printgroupname and tenantid = @TenantId )
			set @corrtypeId =(select [CorrespondenceTypeId] from [dbo].[CorrespondenceType] where [Name] = @CRPTName and tenantid = @TenantId)

			

			IF (@printgroupId is not null and @corrtypeId is not  null)
				BEGIN
		
		IF NOT EXISTS (SELECT 1  FROM [dbo].[CorrespondenceTypePrintGroup] where [PrintGroupId] = @printgroupId  and  [CorrespondenceTypeId] = @corrtypeId and tenantid = @TenantId)
		
		begin

			INSERT INTO [dbo].[CorrespondenceTypePrintGroup] ([CorrespondenceTypeId],[PrintGroupId],[CreatedBy],[CreatedDate],[UpdatedBy],[UpdatedDate],[TenantId],[PrintGroupAssignmentType],[CorrelationId],[ModuleId])
			VALUES (@corrtypeId,@printgroupId, @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm, @TenantId,@PrintgroupAssignmenttype,@CorrelationId, @ModuleId)

			IF @debug = 1 PRINT '[CorrespondenceTypePrintGroup] inserted successfully - Name: '+@Printgroupname+' ' 

			end				
	else
	IF @debug = 1 PRINT '[CorrespondenceTypePrintGroup] already exists - Name: '+@Printgroupname+' ' 
	end
	
	Else
	IF @debug = 1 PRINT '[CorrespondenceTypePrintGroup] are empty - Name: '+@Printgroupname+' ' 
	Print '------------------------------------------------------------------'		

		FETCH NEXT FROM curList
		INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@Printgroupname,@PrintgroupAssignmenttype,@ModuleId,@CRPTDisplayId
		
	END

	CLOSE curList;
	DEALLOCATE curList;


	-------------------------------------------------------------------------
	----------------------Insert into [dbo].[Template]--------------------
	DECLARE @TemplateId               int           
	,       @TemplateOpenXML          nvarchar(MAX)
	,       @TStatus                  nvarchar(25)  
	,       @Tversion                 int     
	,       @NewTversion                 int   
	,       @TsystemFieldDictionaryId int            =NULL
	,       @TCorrespondenceTypeId    int           
	,       @TDisplayId               nvarchar(50)  
	,       @newDisplayid             nvarchar(50) 
	,       @TDescription             varchar(50)   
	,       @CName                    nvarchar(250) 
	,       @contentname              nvarchar(250)
	,       @AssociationType          nvarchar(250)

			
	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES 


	--Section 2: Data for [Template]
	-----------------------------------------------------------------------------------------------
	%s
	-----------------------------------------------------------------------------------------------
	--End of Section 2: Data for [Template]

	) v (  [TemplateOpenXml],[TemplateDisplayId],[Status],[Version],[Description],[CName],[ModuleId],[contentname], [AssociationType]);
			
	--PRINT '--Table . [Template]--------------------------------------------------------------------------'

	OPEN curList
	FETCH NEXT FROM curList
	INTO   @TemplateOpenXML ,@TDisplayId , @TStatus ,@Tversion ,@TDescription ,@CName,@ModuleId, @contentname , @AssociationType

	WHILE @@FETCH_STATUS = 0
	BEGIN	

	set @TCorrespondenceTypeId = (select [CorrespondenceTypeId] from [dbo].[CorrespondenceType] where Name= @CName and TenantId=@TenantId )



	IF  EXISTS (SELECT 1
			FROM [Template]
			WHERE [TenantId]= @tenantId
				AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
				AND [TemplateDisplayId] = @newDisplayid
				AND [Description] = @TDescription
				And [Moduleid] =@ModuleId
				
	)

		
		IF @debug = 1 PRINT 'Template already exists Name/Version            : '+@newDisplayid+'/'+@CName+'/'+cast(@Tversion as nvarchar(16))

	IF  NOT EXISTS (SELECT 1
			FROM [dbo].[Template]
			WHERE [TenantId]= @tenantId
				AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
				AND [TemplateDisplayId] = @newDisplayid
				And [Moduleid] =@ModuleId
			
	) 
		BEGIN

		set @NewTversion = isnull((select (version)+1 from template where  [TenantId]= @tenantId AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId and status = 'Active'),1)

		
			INSERT INTO [dbo].[Template] ( [TemplateOpenXml],                    [Status], [Version], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [SystemFieldDictionaryId], [CorrespondenceTypeId], [TemplateDisplayId], [Description], [CorrelationId], moduleId )
			VALUES                       ( (select [TemplateOpenXml] from [Template] WHERE [TenantId]=@SourceTenantId  and [TemplateDisplayId] = @TDisplayId and status = 'ACTIVE' and  version= @Tversion)
											--(select [TemplateOpenXml] from templateExport2..template_copy1 WHERE [TenantId]=@SourceTenantId  and [TemplateDisplayId] = @TDisplayId and status= 'ACTIVE')
			, @TStatus, @NewTversion, @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantId,  @TsystemFieldDictionaryId, @TCorrespondenceTypeId, @TDisplayId,         @TDescription, @CorrelationId, @ModuleId  )



			IF @debug = 1 PRINT '[Template] inserted successfully [Name/Version] : '+@TDisplayId+'/'++@CName+'/'+cast(@Tversion as nvarchar(16))
			

			Update [dbo].[Template]
			set [Status] = 'INACTIVE'
			,updatedby = 'AdminUser'
			WHERE [TenantId]= @tenantId
			AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
			and version != @NewTversion

		SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'TMP' AND [TenantId]=@TenantId)
		SET @newDisplayid = 'TMP'+FORMAT(@SeqID, '000000000')

		UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'TMP' AND [TenantId]=@TenantId
		
		update template
		set TemplateDisplayId = @newDisplayid
		where  [CorrespondenceTypeId]= @TCorrespondenceTypeId
		AND [Version] = @NewTversion
		and [TenantId]= @tenantId
		IF @debug = 1 PRINT '[Template] Displayid Updated                    :'+@newDisplayid
		
	END

	--insert into [dbo].[TemplateReusableContent]

	if (@contentname is not null)

	begin

		INSERT INTO [dbo].[TemplateReusableContent] ([CreatedBy],[CreatedDate],[ReusableContentId],[TemplateId],[TenantId],[UpdatedBy],[UpdatedDate],[EffectiveEndDate],[CorrelationId],[ReusableContentAssociationType], [ModuleId])
		values (@CreatedBy, @CreateDate,
		(select  b.[ReusableContentId] from [dbo].[ReusableContentType] a join [dbo].[ReusableContent] b on a.[ReusableContentTypeId] = b.[ReusableContentTypeId] 
		where a.[Name] =@contentname  and a.tenantid =@TenantId and b.status = 'ACTIVE' )
		,(select templateid from template where 	 [TenantId]= @tenantId
				AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
				--AND [TemplateDisplayId] = @TDisplayId 
				and status = 'ACTIVE')
				
			,@TenantId,@UpdatedBy,@UpdatedDttm, NULL, null, @AssociationType , @ModuleId)
		
	print 'Template_Reusablecontent inserted               : ' +@AssociationType
	end 

	--Template cache cleanup--
	

	delete from [TemplateCache] where tenantid = @TenantId and templateid = (select templateid from template where 	 [TenantId]= @tenantId AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId 
																			and status = 'ACTIVE')
	IF @debug = 1
						PRINT '[TemplateCache] truncated for                   : ' +@CName 
	Print '------------------------------------------------------------------'


		FETCH NEXT FROM curList
	INTO   @TemplateOpenXML ,@TDisplayId , @TStatus ,@Tversion ,@TDescription ,@CName,@ModuleId, @contentname , @AssociationType
	END

	CLOSE curList;
	DEALLOCATE curList;



	

	--------------updating  [CorrespondenceTypeDisplayId]----------------------------
				

	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES 


	----Section 3: Data for [CorrespondenceType]

	-------------------------------------------------------------------------------------------------
	%s
	-------------------------------------------------------------------------------------------------
	----End of Section 3: Data for [dbo].[CorrespondenceType]

	) v ([Name],[CertifiedMailFlag],[ContextLevel],[FTIFlag],[WorkflowSubTypeId],[WorkflowTypeId],[CTCName],[ModuleId],[CorrespondenceTypeDisplayId]);


	OPEN curList
	FETCH NEXT FROM curList
	INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@ModuleId,@CRPTDisplayId 

	WHILE @@FETCH_STATUS = 0
	BEGIN
		

		IF  EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND [CorrespondenceTypeDisplayId]=@CRPTDisplayId 
																			AND TenantId = @TenantId)

		BEGIN

			SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'COR' AND [TenantId]=@TenantId)
			SET @CRPTDisplayId = 'COR'+FORMAT(@SeqID, '000000000')
			
		Update [dbo].[CorrespondenceType]
		set [CorrespondenceTypeDisplayId] = @CRPTDisplayId
		where TenantId = @TenantId
		AND Name= @CRPTName
		
		
		UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'COR' AND [TenantId]=@TenantId

		IF @debug = 1 PRINT '[CorrespondenceType] updated successfully - Name: '+@CRPTName+' '
		Print +@CRPTDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''

		END

		FETCH NEXT FROM curList
		INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@ModuleId,@CRPTDisplayId 
		
	END

	CLOSE curList;
	DEALLOCATE curList;
	Print '------------------------------------------------------------------'
	--------------------------------------------------------------------------------------------

	-------------Inserting Roles to [CorrespondenceTypeUserRole]  ------------------------


	DECLARE @corresTypeName nvarchar(50)
	DECLARE @corresTypeId int
	DECLARE @Userroleid uniqueidentifier
	DECLARE @rolename varchar(250)
	DECLARE @CanAddFlag bit
	DECLARE @TenantId2 uniqueidentifier = '00000000-0000-0000-0000-000000000000'

	DECLARE curList CURSOR FOR SELECT *
	FROM (
	VALUES

	-------------------------------------------------------------------------------------------------
	----Section 4: Data for [CorrespondenceType]
	-------------------------------------------------------------------------------------------------
	%s
	-------------------------------------------------------------------------------------------------
	----End of Section 4: Data for [CorrespondenceType]
	-------------------------------------------------------------------------------------------------
	) v ([rolename],[CorresTypeName],[ModuleId],[CanAddFlag]);


	OPEN curList
	FETCH NEXT FROM curList
	INTO @rolename, @corresTypeName, @moduleId,@CanAddFlag

	WHILE @@FETCH_STATUS = 0
	BEGIN

		SET @corresTypeId = (SELECT [CorrespondenceTypeId] FROM [CorrespondenceType] WHERE tenantid = @tenantid AND name = @corresTypeName)
		SET @Userroleid = (SELECT [RoleId] FROM ref_NgRoleManagement..role WHERE tenantid IN( @tenantid ,@tenantid2) AND name = @rolename
							and moduleId = @ModuleId /*and name like 'STL'*/ )
		

		
		IF EXISTS (
							SELECT 1
								FROM [CorrespondenceTypeUserRole]
								WHERE [CorrespondenceTypeId] = @corresTypeId
									AND [CanAddFlag] = @CanAddFlag
									and moduleid = @moduleId
									AND [UserRoleId] = (
															SELECT [RoleId] FROM ref_NgRoleManagement..role 
															WHERE tenantid IN( @tenantid ,@tenantid2) AND name =@rolename /*AND name LIKE 'STL%'*/
															AND roleid IN (SELECT DISTINCT UserRoleId FROM ref_NgCorrespondence..[CorrespondenceTypeUserRole] WHERE CorrespondenceTypeId =@corresTypeId)
													)
													AND [CorrelationId] is null
													
						)
		BEGIN
		set  @CorrelationId = newid()

		update [CorrespondenceTypeUserRole]
		set [CorrelationId] = @CorrelationId
		,updatedby  = @UpdatedBy
		,updateddate = @UpdatedDttm
		WHERE [CorrespondenceTypeId] = @corresTypeId
									AND [CanAddFlag] = @CanAddFlag
									and moduleid = @moduleId
									AND [UserRoleId] = (
															SELECT [RoleId] FROM ref_NgRoleManagement..role 
															WHERE tenantid IN( @tenantid ,@tenantid2) AND name =@rolename /*AND name LIKE 'STL%'*/
															AND roleid IN (SELECT DISTINCT UserRoleId FROM ref_NgCorrespondence..[CorrespondenceTypeUserRole] WHERE CorrespondenceTypeId =@corresTypeId)
													)
													
		and [CorrelationId] is null
		
		PRINT '[CorrespondenceTypeUserRole] CorrelationID Updated - Role/Correspondence Type: '+@rolename+'/'+@corresTypeName+'' 

		END

		IF @Userroleid IS NULL
			PRINT '[Role] does not exists: [RoleName] = '+@rolename


		IF @Userroleid IS NOT NULL
		BEGIN 

						
			IF NOT EXISTS (
							SELECT 1
								FROM [CorrespondenceTypeUserRole]
								WHERE [CorrespondenceTypeId] = @corresTypeId
									AND [CanAddFlag] = @CanAddFlag
									and moduleid = @moduleId
									AND [UserRoleId] = (
															SELECT [RoleId] FROM ref_NgRoleManagement..role 
															WHERE tenantid IN( @tenantid ,@tenantid2) AND name =@rolename /*AND name LIKE 'STL%'*/
															AND roleid IN (SELECT DISTINCT UserRoleId FROM ref_NgCorrespondence..[CorrespondenceTypeUserRole] WHERE CorrespondenceTypeId =@corresTypeId)
													)
						)

				
				BEGIN

				set  @CorrelationId = newid()
					INSERT INTO [CorrespondenceTypeUserRole] ( [CorrespondenceTypeId], [UserRoleId], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [CanAddFlag], [CorrelationId], [ModuleId] )
					VALUES                                   ( @corresTypeId,          @Userroleid,  @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantId,  @CanAddFlag,  @CorrelationId , @ModuleId )

					IF @debug = 1
						PRINT '[CorrespondenceTypeUserRole] inserted successfully - Role/Correspondence Type: '+@rolename+'/'+@corresTypeName+''
				END
			ELSE 
				PRINT '[CorrespondenceTypeUserRole] already exists - Role/Correspondence Type: '+@rolename+'/'+@corresTypeName+''
		END

	

		FETCH NEXT FROM curList
		INTO @rolename, @corresTypeName, @moduleId,@CanAddFlag

	END

	CLOSE curList;
	DEALLOCATE curList;



	-----------------------------------------------------------------
		--Tenant Summary--

		print '-----------------------------------------------------------' 
		print 'Counts for the tenant: '+cast(@TenantId as varchar(150)) 
		print '-----------------------------------------------------------'

		SELECT count(*) as [CorrespondenceType] FROM [dbo].[CorrespondenceType] where TenantId=@TenantId 	
		SELECT count(*) as CorrespondenceTypeCategory FROM [dbo].[CorrespondenceTypeCategory] where TenantId=@TenantId
		select count(*) as Template from [dbo].[Template]
		select count(*) as [ReusableContent] from [dbo].[ReusableContent]
		select count(*) as [ReusableContentType] from [dbo].[ReusableContentType]	

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

	END CATCH`, useCorrectDB, targetTenantId, sourceTenantID, finishedJSON1, finishedJSON2, finishedJSON3, finishedJSON4)

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

	log.Printf("target correspondence template script executed successfully. rows affected: %d", rowsAffected)
	return nil

}

func TargetCorrespondenceTemplateMigrationGOVSQL(targetEnv string, sourceTenantID string, targetTenantId string, targetUserName string, targetDBPassword string, data1 []CursorData1, data2 []CursorData2, data3 []CursorData3, data4 []CursorData4) (err error) {
	log.Printf("Beginning Target Correspondence Commercial SQL Migration Script 1")

	// Get the SQL configuration for the target environment
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
	file1, err := os.Open("corrsourcecursordatasection1.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()

	// read file contents
	byteValue1, err := io.ReadAll(file1)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue1) == 0 {
		log.Fatal("corrsourcecursordatasection1.json file is empty")
	}

	var unmarshaledData1 []CursorData1
	err = json.Unmarshal(byteValue1, &unmarshaledData1)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON1 := unmarshaledData1[0].CursorData1
	for numbah, rows := range data1 {
		if numbah == 0 {
			finishedJSON1 = rows.CursorData1
		} else {
			finishedJSON1 = fmt.Sprintf("%s,%s", finishedJSON1, rows.CursorData1)
		}
	}

	// read artifacted file
	file2, err := os.Open("corrsourcecursordatasection2.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// read file contents
	byteValue2, err := io.ReadAll(file2)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue2) == 0 {
		log.Fatal("corrsourcecursordatasection2.json file is empty")
	}

	var unmarshaledData2 []CursorData2
	err = json.Unmarshal(byteValue2, &unmarshaledData2)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON2 := unmarshaledData2[0].CursorData2
	for numbah, rows := range data2 {
		if numbah == 0 {
			finishedJSON2 = rows.CursorData2
		} else {
			finishedJSON2 = fmt.Sprintf("%s,%s", finishedJSON2, rows.CursorData2)
		}
	}

	// read artifacted file
	file3, err := os.Open("corrsourcecursordatasection3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// read file contents
	byteValue3, err := io.ReadAll(file3)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue3) == 0 {
		log.Fatal("corrsourcecursordatasection3.json file is empty")
	}

	var unmarshaledData3 []CursorData3
	err = json.Unmarshal(byteValue2, &unmarshaledData3)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON3 := unmarshaledData3[0].CursorData3
	for numbah, rows := range data3 {
		if numbah == 0 {
			finishedJSON2 = rows.CursorData3
		} else {
			finishedJSON2 = fmt.Sprintf("%s,%s", finishedJSON2, rows.CursorData3)
		}
	}

	// read artifacted file
	file4, err := os.Open("corrsourcecursordatasection4.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	// read file contents
	byteValue4, err := io.ReadAll(file4)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	if len(byteValue4) == 0 {
		log.Fatal("corrsourcecursordatasection4.json file is empty")
	}

	var unmarshaledData4 []CursorData4
	err = json.Unmarshal(byteValue4, &unmarshaledData4)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON4 := unmarshaledData4[0].Cursor_for_Roles_section_4
	for numbah, rows := range data4 {
		if numbah == 0 {
			finishedJSON2 = rows.Cursor_for_Roles_section_4
		} else {
			finishedJSON2 = fmt.Sprintf("%s,%s", finishedJSON2, rows.Cursor_for_Roles_section_4)
		}
	}

	// Determine the correct database to use
	var useCorrectDB = "NgCorrespondence"
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
,       @ModuleId      int

--IF N'$(varDebug)' = 1 SET @debug = 1 ELSE SET @debug = 0;
--IF N'$(varCommit)' = 1 SET @commit = 1 ELSE SET @commit = 0;
SET @debug = 1; --For testing 
SET @commit = 1; --For testing 

BEGIN TRY
    BEGIN TRANSACTION


----Table 1. Insert into [dbo].[CorrespondenceType]--------------------------------------------------------------------------

DECLARE @CRPTName              nvarchar(50)
--,       @CRPTDescription nvarchar(200)
,       @CRPTCategoryId        int
,       @CRPTPrimaryIdType     varchar(50)
,       @RowVersion            timestamp
,       @CRPTDisplayId         nvarchar(50)
,       @CRPTCertifiedMailFlag bit
,       @CRPTContextLevel      nvarchar(MAX)
,       @CRPTFTIFlag           bit
,       @CRPTWorkflowSubTypeId nvarchar(50)
,       @CRPTworkflowTypeId    nvarchar(50)
,       @CRPTCName             nvarchar(50)
,       @Printgroupname        nvarchar(50)
,       @PrintgroupAssignmenttype nvarchar(50)
,       @printgroupId int
,       @corrtypeId   int
,       @wfsubtype    int --nvarchar(50)
,       @wftype       int  --nvarchar(50)

DECLARE curList CURSOR FOR SELECT *
FROM (
VALUES 

--Section 1: Data for [CorrespondenceType]
-----------------------------------------------------------------------------------------------
%s
-----------------------------------------------------------------------------------------------
--End of Section 1: Data for [dbo].[CorrespondenceType]

) v ([Name],[CertifiedMailFlag],[ContextLevel],[FTIFlag],[WorkflowSubTypeId],[WorkflowTypeId],[CTCName],[PGName],[PGAssignmentType], [ModuleId],[CorrespondenceTypeDisplayId]);


OPEN curList
FETCH NEXT FROM curList
INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@Printgroupname,@PrintgroupAssignmenttype,@ModuleId,@CRPTDisplayId

WHILE @@FETCH_STATUS = 0
BEGIN

SET @CorrelationId = newid()
SET @wfsubtype = (select [WorkflowVariants].[WorkflowVariantId] from [ref_WorkflowEngineDB]..[WorkflowVariants] where WorkflowVariants.code = @CRPTWorkflowSubTypeId and WorkflowVariants.tenantid= @tenantid and moduleid = @moduleid)

--IF (@CRPTWorkflowSubTypeId is not null or  @wfsubtype is null)

--Begin
--PRINT 'Workflow Subtype "'+@wfsubtype+'" Does not exist for TenandID "'+ cast(@tenantId as varchar(50))+ '"!';
--	THROW 50001, 'Workflow subtype does not exist.', 1;
--END

SET @wftype    = (select [WorkflowGroups].[WorkflowGroupId] from [ref_WorkflowEngineDB]..[WorkflowGroups] where [WorkflowGroups].code = @CRPTworkflowTypeId and [WorkflowGroups].tenantid= @tenantid)

--IF (@CRPTworkflowTypeId is not null) and (@wftype is null)

--Begin
--PRINT 'Workflow type "'+@wftype+'" Does not exist for TenandID "'+ cast(@tenantId as varchar(50))+ '"!';
--	THROW 50001, 'Workflow type does not exist.', 1;
--END	

	IF EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and moduleid = @ModuleId)
	
		IF @debug = 1 PRINT 'CorrespondenceType already exists - Name: '+@CRPTName+' '	
		
	IF  EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and moduleid = @ModuleId)
	and  not exists (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and contextlevel = @CRPTContextLevel and [WorkflowSubTypeId]=@wfsubtype and [WorkflowTypeId]= @wftype  and moduleid = @ModuleId)
   
   begin
	   update [CorrespondenceType]
	   set [ContextLevel] = @CRPTContextLevel
	   ,[WorkflowSubTypeId] = @wfsubtype
	   ,[WorkflowTypeId] =    @wftype
	   ,[UpdatedBy] = @UpdatedBy
	   ,[UpdatedDate]= @UpdatedDttm
	   ,moduleid = @ModuleId
	   where Name= @CRPTName 
	   AND TenantId = @TenantId
	   and moduleid = @ModuleId


	   IF @debug = 1 PRINT 'CorrespondenceType Contectlevel/[WorkflowSubTypeId]/[WorkflowTypeId] are updated - Name: '+@CRPTName+' '	
	   
   end

   IF exists (select 1 from [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and  [CorrelationId]= NULL)
   BEGIN
   
   SET @CorrelationId = newid()
   update [CorrespondenceType]
	   set  [CorrelationId] = @CorrelationId
	   ,[UpdatedBy] = @UpdatedBy
	   ,[UpdatedDate]= @UpdatedDttm
	    where Name= @CRPTName 
	    AND TenantId = @TenantId
		and moduleid = @ModuleId

		IF @debug = 1 PRINT 'CorrespondenceType CorrelationId is updated - Name: '+@CRPTName+' '	
   END

   SET @CRPTCategoryId = (SELECT top 1 [CorrespondenceTypeCategoryId]
		FROM [dbo].[CorrespondenceTypeCategory]
		WHERE Name = @CRPTCName
			AND TenantId =@TenantId
			/*and moduleid = @ModuleId*/)
			  
---fix for correspondencetypecategoryid nulls-- 5/9/2024


if exists(select 1 from [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId and  [CorrespondenceTypeCategoryId]= NULL)
 
  
begin
   update [CorrespondenceType]
	   set  [CorrespondenceTypeCategoryId] = @CRPTCategoryId
	   ,[UpdatedBy] = @UpdatedBy
	   ,[UpdatedDate]= @UpdatedDttm
	    where Name= @CRPTName 
	    AND TenantId = @TenantId
		--and moduleid = @ModuleId

IF @debug = 1 PRINT 'CorrespondenceType [CorrespondenceTypeCategoryId] is updated - Name: '+@CRPTName+' '	
	end	
   

    IF NOT EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND TenantId = @TenantId)

	BEGIN
	
	

		INSERT INTO [CorrespondenceType] ( [Name],    [CorrespondenceTypeCategoryId], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [CorrespondenceTypeDisplayId], [CertifiedMailFlag],    [ContextLevel],    [WorkflowSubTypeId], [WorkflowTypeId] ,  [FTIFlag]  , ModuleId  ,  [CorrelationId]  )
		VALUES                           ( @CRPTName, @CRPTCategoryId,                @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantId,  @CRPTDisplayId,                @CRPTCertifiedMailFlag,  @CRPTContextLevel, @wftype,            @wftype,            @CRPTFTIFlag , @ModuleId, @CorrelationId);

	

		IF @debug = 1 PRINT '[CorrespondenceType] inserted successfully - Name: '+@CRPTName+' '
	END

---updated 9/26 --		
		--Insert into [dbo].[CorrespondenceTypePrintGroup]
		
		set @printgroupId = (select [PrintGroupId] from [dbo].[PrintGroup] where [Name] = @Printgroupname and tenantid = @TenantId )
		set @corrtypeId =(select [CorrespondenceTypeId] from [dbo].[CorrespondenceType] where [Name] = @CRPTName and tenantid = @TenantId)

		

        IF (@printgroupId is not null and @corrtypeId is not  null)
			BEGIN
	
	   IF NOT EXISTS (SELECT 1  FROM [dbo].[CorrespondenceTypePrintGroup] where [PrintGroupId] = @printgroupId  and  [CorrespondenceTypeId] = @corrtypeId and tenantid = @TenantId)
	  
	  begin

		INSERT INTO [dbo].[CorrespondenceTypePrintGroup] ([CorrespondenceTypeId],[PrintGroupId],[CreatedBy],[CreatedDate],[UpdatedBy],[UpdatedDate],[TenantId],[PrintGroupAssignmentType],[CorrelationId],[ModuleId])
		VALUES (@corrtypeId,@printgroupId, @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm, @TenantId,@PrintgroupAssignmenttype,@CorrelationId, @ModuleId)

		IF @debug = 1 PRINT '[CorrespondenceTypePrintGroup] inserted successfully - Name: '+@Printgroupname+' ' 

		end				
  else
  IF @debug = 1 PRINT '[CorrespondenceTypePrintGroup] already exists - Name: '+@Printgroupname+' ' 
  end
  
  Else
  IF @debug = 1 PRINT '[CorrespondenceTypePrintGroup] are empty - Name: '+@Printgroupname+' ' 
Print '------------------------------------------------------------------'		

	FETCH NEXT FROM curList
	INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@Printgroupname,@PrintgroupAssignmenttype,@ModuleId,@CRPTDisplayId
	
END

CLOSE curList;
DEALLOCATE curList;


-------------------------------------------------------------------------
----------------------Insert into [dbo].[Template]--------------------
DECLARE @TemplateId               int           
,       @TemplateOpenXML          nvarchar(MAX)
,       @TStatus                  nvarchar(25)  
,       @Tversion                 int     
,       @NewTversion                 int   
,       @TsystemFieldDictionaryId int            =NULL
,       @TCorrespondenceTypeId    int           
,       @TDisplayId               nvarchar(50)  
,       @newDisplayid             nvarchar(50) 
,       @TDescription             varchar(50)   
,       @CName                    nvarchar(250) 
,       @contentname              nvarchar(250)
,       @AssociationType          nvarchar(250)

		   
DECLARE curList CURSOR FOR SELECT *
FROM (
VALUES 


--Section 2: Data for [Template]
-----------------------------------------------------------------------------------------------
%s
-----------------------------------------------------------------------------------------------
--End of Section 2: Data for [Template]

) v (  [TemplateOpenXml],[TemplateDisplayId],[Status],[Version],[Description],[CName],[ModuleId],[contentname], [AssociationType]);
		
--PRINT '--Table . [Template]--------------------------------------------------------------------------'

OPEN curList
FETCH NEXT FROM curList
INTO   @TemplateOpenXML ,@TDisplayId , @TStatus ,@Tversion ,@TDescription ,@CName,@ModuleId, @contentname , @AssociationType

WHILE @@FETCH_STATUS = 0
BEGIN	

set @TCorrespondenceTypeId = (select [CorrespondenceTypeId] from [dbo].[CorrespondenceType] where Name= @CName and TenantId=@TenantId )



IF  EXISTS (SELECT 1
		FROM [Template]
		WHERE [TenantId]= @tenantId
			AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
			AND [TemplateDisplayId] = @newDisplayid
			AND [Description] = @TDescription
			And [Moduleid] =@ModuleId
			
)

	
       IF @debug = 1 PRINT 'Template already exists Name/Version            : '+@newDisplayid+'/'++@CName+'/'+cast(@Tversion as nvarchar(16))

IF  NOT EXISTS (SELECT 1
		FROM [dbo].[Template]
		WHERE [TenantId]= @tenantId
			AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
			AND [TemplateDisplayId] = @newDisplayid
			And [Moduleid] =@ModuleId
		
) 
	BEGIN

       set @NewTversion = isnull((select (version)+1 from template where  [TenantId]= @tenantId AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId and status = 'Active'),1)

	
		INSERT INTO [dbo].[Template] ( [TemplateOpenXml],                    [Status], [Version], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [SystemFieldDictionaryId], [CorrespondenceTypeId], [TemplateDisplayId], [Description], [CorrelationId], moduleId )
		VALUES                       ( --(select [TemplateOpenXml] from [Template] WHERE [TenantId]=@SourceTenantId  and [TemplateDisplayId] = @TDisplayId and status = 'ACTIVE' and  version= @Tversion)
		                                (select [TemplateOpenXml] from templateExport2..template_copy1 WHERE [TenantId]=@SourceTenantId  and [TemplateDisplayId] = @TDisplayId and status= 'ACTIVE')
		, @TStatus, @NewTversion, @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantId,  @TsystemFieldDictionaryId, @TCorrespondenceTypeId, @TDisplayId,         @TDescription, @CorrelationId, @ModuleId  )



		IF @debug = 1 PRINT '[Template] inserted successfully [Name/Version] : '+@TDisplayId+'/'++@CName+'/'+cast(@Tversion as nvarchar(16))
		

	    Update [dbo].[Template]
		set [Status] = 'INACTIVE'
		,updatedby = 'AdminUser'
		WHERE [TenantId]= @tenantId
		AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
		and version != @NewTversion

	SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'TMP' AND [TenantId]=@TenantId)
    SET @newDisplayid = 'TMP'+FORMAT(@SeqID, '000000000')

	UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'TMP' AND [TenantId]=@TenantId
	
	update template
	set TemplateDisplayId = @newDisplayid
	where  [CorrespondenceTypeId]= @TCorrespondenceTypeId
	AND [Version] = @NewTversion
	and [TenantId]= @tenantId
	IF @debug = 1 PRINT '[Template] Displayid Updated                    :'+@newDisplayid
	
END

--insert into [dbo].[TemplateReusableContent]

if (@contentname is not null)

begin

	INSERT INTO [dbo].[TemplateReusableContent] ([CreatedBy],[CreatedDate],[ReusableContentId],[TemplateId],[TenantId],[UpdatedBy],[UpdatedDate],[EffectiveEndDate],[CorrelationId],[ReusableContentAssociationType], [ModuleId])
	values (@CreatedBy, @CreateDate,
	(select  b.[ReusableContentId] from [dbo].[ReusableContentType] a join [dbo].[ReusableContent] b on a.[ReusableContentTypeId] = b.[ReusableContentTypeId] 
       where a.[Name] =@contentname  and a.tenantid =@TenantId and b.status = 'ACTIVE' )
	   ,(select templateid from template where 	 [TenantId]= @tenantId
			AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId
			--AND [TemplateDisplayId] = @TDisplayId 
			and status = 'ACTIVE')
			
		,@TenantId,@UpdatedBy,@UpdatedDttm, NULL, null, @AssociationType , @ModuleId)
	
 print 'Template_Reusablecontent inserted               : ' +@AssociationType
 end 

  --Template cache cleanup--
 

  delete from [TemplateCache] where tenantid = @TenantId and templateid = (select templateid from template where 	 [TenantId]= @tenantId AND  [CorrespondenceTypeId]= @TCorrespondenceTypeId 
                                                                           and status = 'ACTIVE')
  IF @debug = 1
					PRINT '[TemplateCache] truncated for                   : ' +@CName 
Print '------------------------------------------------------------------'


	FETCH NEXT FROM curList
INTO   @TemplateOpenXML ,@TDisplayId , @TStatus ,@Tversion ,@TDescription ,@CName,@ModuleId, @contentname , @AssociationType
END

CLOSE curList;
DEALLOCATE curList;



 

--------------updating  [CorrespondenceTypeDisplayId]----------------------------
             

DECLARE curList CURSOR FOR SELECT *
FROM (
VALUES 


----Section 3: Data for [CorrespondenceType]

-------------------------------------------------------------------------------------------------
%s
-------------------------------------------------------------------------------------------------
----End of Section 3: Data for [dbo].[CorrespondenceType]

) v ([Name],[CertifiedMailFlag],[ContextLevel],[FTIFlag],[WorkflowSubTypeId],[WorkflowTypeId],[CTCName],[ModuleId],[CorrespondenceTypeDisplayId]);


OPEN curList
FETCH NEXT FROM curList
INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@ModuleId,@CRPTDisplayId 

WHILE @@FETCH_STATUS = 0
BEGIN
     

    IF  EXISTS (SELECT 1  FROM [CorrespondenceType] WHERE Name= @CRPTName AND [CorrespondenceTypeDisplayId]=@CRPTDisplayId 
	                                                                      AND TenantId = @TenantId)

	BEGIN

		SET @SeqID = (SELECT [LastUsedNumber]+1 FROM [dbo].[IdSequence] WHERE [Prefix] = 'COR' AND [TenantId]=@TenantId)
		SET @CRPTDisplayId = 'COR'+FORMAT(@SeqID, '000000000')
		
       Update [dbo].[CorrespondenceType]
	   set [CorrespondenceTypeDisplayId] = @CRPTDisplayId
	   where TenantId = @TenantId
	   AND Name= @CRPTName
	 
	
	   UPDATE [IdSequence] SET [LastUsedNumber] = @SeqID WHERE [Prefix] = 'COR' AND [TenantId]=@TenantId

	   IF @debug = 1 PRINT '[CorrespondenceType] updated successfully - Name: '+@CRPTName+' '
	   Print +@CRPTDisplayId+ '/' + cast(@SeqID as nvarchar(50)) + ''

	END

	FETCH NEXT FROM curList
	INTO @CRPTName,@CRPTCertifiedMailFlag,@CRPTContextLevel,@CRPTFTIFlag,@CRPTWorkflowSubTypeId,@CRPTworkflowTypeId,@CRPTCName,@ModuleId,@CRPTDisplayId 
	
END

CLOSE curList;
DEALLOCATE curList;
Print '------------------------------------------------------------------'
--------------------------------------------------------------------------------------------

-------------Inserting Roles to [CorrespondenceTypeUserRole]  ------------------------


DECLARE @corresTypeName nvarchar(50)
DECLARE @corresTypeId int
DECLARE @Userroleid uniqueidentifier
DECLARE @rolename varchar(250)
DECLARE @CanAddFlag bit
DECLARE @TenantId2 uniqueidentifier = '00000000-0000-0000-0000-000000000000'

DECLARE curList CURSOR FOR SELECT *
FROM (
VALUES

-------------------------------------------------------------------------------------------------
----Section 4: Data for [CorrespondenceType]
-------------------------------------------------------------------------------------------------
%s
-------------------------------------------------------------------------------------------------
----End of Section 4: Data for [CorrespondenceType]
-------------------------------------------------------------------------------------------------
) v ([rolename],[CorresTypeName],[ModuleId],[CanAddFlag]);


OPEN curList
FETCH NEXT FROM curList
INTO @rolename, @corresTypeName, @moduleId,@CanAddFlag

WHILE @@FETCH_STATUS = 0
BEGIN

	SET @corresTypeId = (SELECT [CorrespondenceTypeId] FROM [CorrespondenceType] WHERE tenantid = @tenantid AND name = @corresTypeName)
	SET @Userroleid = (SELECT [RoleId] FROM ref_NgRoleManagement..role WHERE tenantid IN( @tenantid ,@tenantid2) AND name = @rolename
	                    and moduleId = @ModuleId /*and name like 'STL%'*/ )
	

	
	IF EXISTS (
						SELECT 1
							FROM [CorrespondenceTypeUserRole]
							WHERE [CorrespondenceTypeId] = @corresTypeId
							    AND [CanAddFlag] = @CanAddFlag
								and moduleid = @moduleId
								AND [UserRoleId] = (
														SELECT [RoleId] FROM ref_NgRoleManagement..role 
														WHERE tenantid IN( @tenantid ,@tenantid2) AND name =@rolename /*AND name LIKE 'STL%'*/
														AND roleid IN (SELECT DISTINCT UserRoleId FROM ref_NgCorrespondence..[CorrespondenceTypeUserRole] WHERE CorrespondenceTypeId =@corresTypeId)
												   )
												   AND [CorrelationId] is null
												  
					  )
	BEGIN
	set  @CorrelationId = newid()

	update [CorrespondenceTypeUserRole]
	set [CorrelationId] = @CorrelationId
	,updatedby  = @UpdatedBy
	,updateddate = @UpdatedDttm
	WHERE [CorrespondenceTypeId] = @corresTypeId
							    AND [CanAddFlag] = @CanAddFlag
								and moduleid = @moduleId
								AND [UserRoleId] = (
														SELECT [RoleId] FROM ref_NgRoleManagement..role 
														WHERE tenantid IN( @tenantid ,@tenantid2) AND name =@rolename /*AND name LIKE 'STL%'*/
														AND roleid IN (SELECT DISTINCT UserRoleId FROM ref_NgCorrespondence..[CorrespondenceTypeUserRole] WHERE CorrespondenceTypeId =@corresTypeId)
												   )
												 
	and [CorrelationId] is null
	
	PRINT '[CorrespondenceTypeUserRole] CorrelationID Updated - Role/Correspondence Type: '+@rolename+'/'+@corresTypeName+'' 

	END

	IF @Userroleid IS NULL
		PRINT '[Role] does not exists: [RoleName] = '+@rolename


	IF @Userroleid IS NOT NULL
	BEGIN 

					  
		IF NOT EXISTS (
						SELECT 1
							FROM [CorrespondenceTypeUserRole]
							WHERE [CorrespondenceTypeId] = @corresTypeId
							    AND [CanAddFlag] = @CanAddFlag
								and moduleid = @moduleId
								AND [UserRoleId] = (
														SELECT [RoleId] FROM ref_NgRoleManagement..role 
														WHERE tenantid IN( @tenantid ,@tenantid2) AND name =@rolename /*AND name LIKE 'STL%'*/
														AND roleid IN (SELECT DISTINCT UserRoleId FROM ref_NgCorrespondence..[CorrespondenceTypeUserRole] WHERE CorrespondenceTypeId =@corresTypeId)
												   )
					  )

			
			BEGIN

			set  @CorrelationId = newid()
				INSERT INTO [CorrespondenceTypeUserRole] ( [CorrespondenceTypeId], [UserRoleId], [CreatedBy], [CreatedDate], [UpdatedBy], [UpdatedDate], [TenantId], [CanAddFlag], [CorrelationId], [ModuleId] )
				VALUES                                   ( @corresTypeId,          @Userroleid,  @CreatedBy,  @CreateDate,   @UpdatedBy,  @UpdatedDttm,  @TenantId,  @CanAddFlag,  @CorrelationId , @ModuleId )

				IF @debug = 1
					PRINT '[CorrespondenceTypeUserRole] inserted successfully - Role/Correspondence Type: '+@rolename+'/'+@corresTypeName+''
			END
		ELSE 
			PRINT '[CorrespondenceTypeUserRole] already exists - Role/Correspondence Type: '+@rolename+'/'+@corresTypeName+''
	END

 

	FETCH NEXT FROM curList
	INTO @rolename, @corresTypeName, @moduleId,@CanAddFlag

END

CLOSE curList;
DEALLOCATE curList;



-----------------------------------------------------------------
	--Tenant Summary--

	print '-----------------------------------------------------------' 
	print 'Counts for the tenant: '+cast(@TenantId as varchar(150)) 
	print '-----------------------------------------------------------'

	SELECT count(*) as [CorrespondenceType] FROM [dbo].[CorrespondenceType] where TenantId=@TenantId 	
	SELECT count(*) as CorrespondenceTypeCategory FROM [dbo].[CorrespondenceTypeCategory] where TenantId=@TenantId
	select count(*) as Template from [dbo].[Template]
	select count(*) as [ReusableContent] from [dbo].[ReusableContent]
	select count(*) as [ReusableContentType] from [dbo].[ReusableContentType]	

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

END CATCH`, useCorrectDB, targetTenantId, sourceTenantID, finishedJSON1, finishedJSON2, finishedJSON3, finishedJSON4)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	lol, err := db.Prepare(sqlScript)
	if err != nil {
		log.Fatalf("Error preparing target correspondence gov script: %v", err)
	}
	result, err := lol.Exec()
	// result, err := db.Exec(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	// debugtown(fmt.Sprintf("RunSQLMigration: Results: %s", result))
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error fetching rows affected: %v", err)
	}

	log.Printf("target correspondence template gov script executed successfully. rows affected: %d", rowsAffected)
	return nil

}
