package platformconfiguration

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

func RunTargetPLCSQL(targetEnv string, targetTenantId string, targetUserName string, targetDBPassword string, data1 []CursorData1, data2 []CursorData2) (err error) {
	log.Printf("Beginning Target PLC SQL Migration Script ")

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
	file, err := os.Open("plcsourcecursordatasection1.json")
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
		log.Fatal("plcsourcecursordatasection1.json file is empty")
	}

	var unmarshaledData1 []CursorData1
	err = json.Unmarshal(byteValue, &unmarshaledData1)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON1 := unmarshaledData1[0].CursorData
	for numbah, rows := range data1 {
		if numbah == 0 {
			finishedJSON1 = rows.CursorData
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON1 = fmt.Sprintf("%s,%s", finishedJSON1, rows.CursorData)
		}

	}

	// read artifacted file
	file, err = os.Open("plcsourcecursordatasection2.json")
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
		log.Fatal("plcsourcecursordatasection2.json file is empty")
	}

	var unmarshaledData2 []CursorData2
	err = json.Unmarshal(byteValue, &unmarshaledData2)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %v", err)
	}

	finishedJSON2 := unmarshaledData2[0].CursorData
	for numbah, rows := range data2 {
		if numbah == 0 {
			finishedJSON2 = rows.CursorData
		} else {
			// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Adding %s to finishedarrays.\n", rows.CursorData))
			finishedJSON2 = fmt.Sprintf("%s,%s", finishedJSON2, rows.CursorData)
		}

	}

	// formsmigration.Debugtown(fmt.Sprintf("RunTargetActiveDirectorySQL1: Finished JSON: %s", finishedJSON))
	var useCorrectDB = "NgPlatform"
	// get sourceEnv and add as a prefix to the DB
	if targetEnv != "" {
		if strings.ToLower(targetEnv) == "ref" {
			useCorrectDB = targetEnv + "_" + useCorrectDB
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		} else {
			useCorrectDB = "NgPlatform"
			fmt.Printf("Use Correct DB: %s\n", useCorrectDB)
		}
	}

	sqlScript := fmt.Sprintf(
		`USE [%s];
		DECLARE      @TenantId       uniqueidentifier = '%s';
		SET NOCOUNT ON
		SET XACT_ABORT ON;
		DECLARE @Name           nvarchar(150)   
		,       @Description    nvarchar(250)   
		,       @CreateDate     datetime         = getutcdate()
		,       @CreatedBy      nvarchar(450)    = 'AdminUser'
		,       @ModifiedDate   datetime         = getutcdate()
		,       @ModifiedBy     nvarchar(450)    = 'AdminUser'
		,       @IsOOTB         bit              = 1
		,       @IsOOTBEditable bit             
		,       @IsDeleted      int              = 0
		,       @newid			uniqueidentifier
		,		@debug			bit				 -- set to 1 for dubugging
		,		@commit			bit				 -- set to 1 for committing
		,		@err_message	nvarchar(1000)
		,       @CorrelationId  uniqueidentifier  
		--IF N'$(varDebug)' = 1 SET @debug = 1 ELSE SET @debug = 0;
		--IF N'$(varCommit)' = 1 SET @commit = 1 ELSE SET @commit = 0;
		SET @debug = 1; --For testing 
		SET @commit = 1; --For testing 
		BEGIN TRY
			BEGIN TRANSACTION
		DECLARE curList CURSOR FOR SELECT *
		FROM (
		VALUES
		--Section 1: Data for PlatformConfigurationGroup
		-----------------------------------------------------------------------------------------------
		%s
		-----------------------------------------------------------------------------------------------
		--End of Section 1: Data for PlatformConfigurationGroup
		) v ([Name],[Description]);
		OPEN curList
		FETCH NEXT FROM curList
		INTO @Name,@Description
		WHILE @@FETCH_STATUS = 0
		BEGIN
			SET @CorrelationId = newid()
			IF EXISTS (SELECT 1
				FROM [PlatformConfigurationGroup]
				WHERE Name=@Name
					AND TenantId = @TenantId
				)
				IF @debug = 1 PRINT 'Group  already exists :- '+@Name+''
			IF NOT EXISTS (SELECT 1
				FROM [PlatformConfigurationGroup]
				WHERE Name=@Name
					AND TenantId= @TenantId
				)
			BEGIN
				INSERT INTO [PlatformConfigurationGroup] ( [Name], [Description], [TenantId], [CreateDate], [CreatedBy], [ModifiedDate], [ModifiedBy], [IsDeleted], [IsOOTB] ,[CorrelationId])
				VALUES                                   ( @Name,  @Description,  @TenantId,  @CreateDate,  @CreatedBy,  @ModifiedDate,  @ModifiedBy,  @IsDeleted,  @IsOOTB , @CorrelationId );
				IF @debug = 1 PRINT 'Group inserted successfully - Name: '+@Name+''
			END
			FETCH NEXT FROM curList
			INTO @Name,@Description
		END
		CLOSE curList;
		DEALLOCATE curList;
		---insert into platormconfiguration table
		DECLARE @GroupName                    nvarchar(150)
		,       @PlatformConfigurationId      uniqueidentifier
		,       @ConfigurationModule          int
		,       @ConfigurationDomain          int
		,       @ConfigurationType            nvarchar(50)
		,       @ConfigurationName            nvarchar(250)
		,       @StateOf                      nvarchar(50)
		,       @ConfigurationDescription     nvarchar(150)
		,       @IsSchema                     bit
		,       @PlatformConfigurationGroupId int
		,       @Version                      int
		,       @ConfigurationInfo            nvarchar(MAX)
		,       @rows                         int
		,       @displayname                  nvarchar(250)      
		DECLARE curList CURSOR FOR SELECT *
		FROM (
		VALUES
		--Section 2: Data for PlatformConfiguration
		-----------------------------------------------------------------------------------------------
		%s
		-----------------------------------------------------------------------------------------------

		--End of Section 2: Data for PlatformConfiguration
		) v ( PlatformConfigurationId, Groupname,ConfigurationModule, ConfigurationDomain, ConfigurationType, ConfigurationName,  StateOf, ConfigurationDescription, Displayname, IsSchema, Version,
		PlatformConfigurationInfo, IsOOTBEditable )
		OPEN curList
		FETCH NEXT FROM curList
		INTO @PlatformConfigurationId, @Groupname,@configurationModule,@ConfigurationDomain,@ConfigurationType,@ConfigurationName,@StateOf,@ConfigurationDescription,@displayname , @IsSchema,@Version,@ConfigurationInfo,@IsOOTBEditable;
		WHILE @@FETCH_STATUS = 0
		BEGIN
			SET @newid = newid()
			SET @CorrelationId = newid()
			IF @ConfigurationInfo NOT IN ('null','') AND ISJSON(@ConfigurationInfo)=0
			BEGIN
				SET @err_message = 'Configuration type/name: '+@ConfigurationType+':'+@ConfigurationName+' has invalid JSON.'
				RAISERROR(@err_message, 16, 1);
			END
			IF replace(replace(@configurationinfo,'"',''),' ','') like '%!A(MISSING)ttributeValue:null%!'
			BEGIN
				SET @err_message = 'Configuration type/name: '+@ConfigurationType+':'+@ConfigurationName+' has "AttributeValue:null" pattern in its JSON'
				RAISERROR(@err_message, 16, 1);
			END
			IF replace(replace(@configurationinfo,'"',''),' ','') like '%!C(MISSING)onfigurationSection:null%!'
			BEGIN
				SET @err_message = 'Configuration type/name: '+@ConfigurationType+':'+@ConfigurationName+' has "ConfigurationSection:null" pattern in its JSON'
				RAISERROR(@err_message, 16, 1);
			END
		--------------------------------
		-- comment out if we don't want to update
			IF EXISTS (SELECT 1
				FROM [PlatformConfiguration]
				WHERE TenantId= @tenantId
					AND configurationModule =@ConfigurationModule
					AND ConfigurationDomain =@ConfigurationDomain
					AND ConfigurationType =@ConfigurationType
					AND ConfigurationName =@ConfigurationName
					AND IsSchema =@IsSchema
					AND isnull(StateOf,'') <> 'Deleted'
				)
			BEGIN
				
				UPDATE [PlatformConfiguration]
				SET StateOf =                  @StateOf
				,   ConfigurationDescription = @ConfigurationDescription
				,   IsSchema =                 @IsSchema
				,   IsOOTB =                   @IsOOTB
				,   ModifiedBy =               @ModifiedBy
				,   ModifiedDate =             @ModifiedDate
				,   displayname  =             @displayname
				,   PlatformConfigurationGroupId = (SELECT PlatformConfigurationGroupId
													FROM PlatformConfigurationGroup
													WHERE Name = @GroupName
													AND TenantId= @tenantId )
				WHERE TenantId= @tenantId
					AND configurationModule =@ConfigurationModule
					AND ConfigurationDomain =@ConfigurationDomain
					AND ConfigurationType =@ConfigurationType
					AND ConfigurationName =@ConfigurationName
					AND IsSchema =@IsSchema
					AND isnull(StateOf,'') <> 'Deleted'
				IF @debug = 1 PRINT 'Configuration already exists for type/name: '+@ConfigurationType+':'+@ConfigurationName+'. Updating.'
				IF NOT EXISTS (SELECT 1
				FROM [PlatformConfigurationInfo]
				WHERE [PlatformConfigurationId] IN
					(SELECT [PlatformConfigurationId]
					FROM [PlatformConfiguration]
					WHERE TenantId= @tenantId
						AND configurationModule =@ConfigurationModule
						AND ConfigurationDomain =@ConfigurationDomain
						AND ConfigurationType =@ConfigurationType
						AND ConfigurationName =@ConfigurationName
						AND IsSchema =@IsSchema)
				)
				BEGIN
					INSERT INTO PlatformConfigurationInfo
					SELECT PlatformConfigurationId,'null'
					FROM PlatformConfiguration p
					WHERE TenantId= @tenantId
						AND configurationModule =@ConfigurationModule
						AND ConfigurationDomain =@ConfigurationDomain
						AND ConfigurationType =@ConfigurationType
						AND ConfigurationName =@ConfigurationName
						AND IsSchema =@IsSchema
					IF @debug = 1 PRINT 'Configuration placeholder inserted successfully: '+@ConfigurationType+':'+@ConfigurationName
				END
			END
		--------------------------------*/
			--adding PlatformConfiguration record if it does not exist
			IF NOT EXISTS (SELECT 1
				FROM [PlatformConfiguration]
				WHERE TenantId= @tenantId
					AND configurationModule =@ConfigurationModule
					AND ConfigurationDomain =@ConfigurationDomain
					AND ConfigurationType =@ConfigurationType
					AND ConfigurationName =@ConfigurationName
					AND IsSchema =@IsSchema
					AND isnull(StateOf,'') <> 'Deleted'
				)
				BEGIN
				INSERT INTO [PlatformConfiguration] ( [PlatformConfigurationId], [ConfigurationModule], [TenantId], [ConfigurationDomain], [ConfigurationType], [ConfigurationName], [CreateDate], [CreatedBy], [ModifiedDate], [ModifiedBy], [StateOf], [ConfigurationDescription], [IsSchema], [PlatformConfigurationGroupId], [Version], [IsOOTB], [IsOOTBEditable], [CorrelationId], displayname)
				VALUES                              ( @newid,                    @configurationModule,  @TenantId,  @ConfigurationDomain,  @ConfigurationType,  @ConfigurationName,  @CreateDate,  @CreatedBy,  @ModifiedDate,  @ModifiedBy,  @StateOf,  @ConfigurationDescription,  @IsSchema,  (SELECT PlatformConfigurationGroupId
																																															FROM PlatformConfigurationGroup
																																															WHERE Name = @GroupName
																																																AND TenantId= @tenantId ),                                                                                                                                                                                                                                                                                   
													@Version,  @IsOOTB,  @IsOOTBEditable , @CorrelationId , @displayname)
				IF @debug = 1 PRINT 'Configuration inserted successfully: '+@ConfigurationType+':'+@ConfigurationName
				END
			-- check Info existence, if there is a record - update, else - insert
			IF EXISTS (SELECT 1
				FROM [PlatformConfigurationInfo]
				WHERE [PlatformConfigurationId] IN
					(SELECT [PlatformConfigurationId]
					FROM [PlatformConfiguration]
					WHERE TenantId= @tenantId
						AND configurationModule =@ConfigurationModule
						AND ConfigurationDomain =@ConfigurationDomain
						AND ConfigurationType =@ConfigurationType
						AND ConfigurationName =@ConfigurationName
						AND IsSchema =@IsSchema
						AND isnull(StateOf,'') <> 'Deleted')
				)  
				BEGIN
					DECLARE @platformConfigurationInfoUpdate nvarchar(MAX)
		--------------------------------
		-- comment out if we don't want to update
					SET @platformConfigurationInfoUpdate = (SELECT TOP 1 PlatformConfigurationId
															FROM [PlatformConfigurationInfo]
															WHERE [PlatformConfigurationId] IN 
																(SELECT [PlatformConfigurationId]
																FROM [PlatformConfiguration]
																WHERE TenantId= @tenantId
																	AND configurationModule =@ConfigurationModule
																	AND ConfigurationDomain =@ConfigurationDomain
																	AND ConfigurationType =@ConfigurationType
																	AND ConfigurationName =@ConfigurationName
																	AND IsSchema =@IsSchema
																	AND isnull(StateOf,'') <> 'Deleted'))
					UPDATE [PlatformConfigurationInfo]
					SET [ConfigurationInfo] = replace(@ConfigurationInfo,@PlatformConfigurationId, @platformConfigurationInfoUpdate)
					WHERE [PlatformConfigurationId] = @platformConfigurationInfoUpdate
					AND [ConfigurationInfo] <> replace(@ConfigurationInfo,@PlatformConfigurationId, @platformConfigurationInfoUpdate)
					set @rows = @@ROWCOUNT	
					IF @rows > 0
						UPDATE PlatformConfiguration
						SET ModifiedBy =   @ModifiedBy
						,   ModifiedDate = @ModifiedDate
						WHERE PlatformConfigurationId=@platformConfigurationInfoUpdate
				
					IF @debug = 1 
					BEGIN
							IF @rows > 0
							PRINT 'ConfigurationInfo updated successfully for type/name: '+@ConfigurationType+':'+@ConfigurationName
						ELSE
							PRINT 'ConfigurationInfo identical for type/name: '+@ConfigurationType+':'+@ConfigurationName
					END
		--------------------------------*/
				END
			ELSE
				IF 	(@ConfigurationInfo <> 'null') 
				BEGIN
				
				IF @debug = 1 
					PRINT 'Inserting --: '+cast (@PlatformConfigurationId as varchar(150))
				INSERT INTO [PlatformConfigurationInfo] ( [PlatformConfigurationId], [ConfigurationInfo] )
				VALUES                                  ( @newid,                    replace(@ConfigurationInfo,@PlatformConfigurationId, @newid))
				
					UPDATE PlatformConfiguration
					SET ModifiedBy =   @ModifiedBy
					,   ModifiedDate = @ModifiedDate
					WHERE PlatformConfigurationId=@newid
				IF @debug = 1 
					PRINT 'ConfigurationInfo inserted successfully: '+cast (@newid as varchar(150))
				END
			FETCH NEXT FROM curList
		INTO @PlatformConfigurationId, @Groupname,@configurationModule,@ConfigurationDomain,@ConfigurationType,@ConfigurationName,@StateOf,@ConfigurationDescription,@displayname , @IsSchema,@Version,@ConfigurationInfo,@IsOOTBEditable;
			END
			CLOSE curList;
			DEALLOCATE curList;
			-----------------------------------------------------------------
			--Tenant Summary--
			declare @cnt nvarchar(100)
			print '-----------------------------------------------------------' 
			print 'Counts for the tenant: '+cast(@TenantId as varchar(150)) 
			print '-----------------------------------------------------------'
			SELECT @cnt=count(*) FROM [PlatformConfigurationGroup] where TenantId=@TenantId
			print 'PlatformConfigurationGroup: '+@cnt 
			SELECT @cnt=count(*) FROM [PlatformConfiguration] where TenantId=@TenantId
			print 'PlatformConfiguration: '+@cnt 
			SELECT @cnt=count(*) FROM [PlatformConfigurationInfo] where PlatformConfigurationId in (SELECT PlatformConfigurationId FROM [PlatformConfiguration] where TenantId=@TenantId)
			print 'PlatformConfigurationInfo: '+@cnt 
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
		END CATCH`, useCorrectDB, targetTenantId, finishedJSON2, finishedJSON1)

	formsmigration.Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	lol, err := db.Prepare(sqlScript)
	if err != nil {
		log.Fatalf("Error preparing PLC SQL script: %v", err)
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

	log.Printf("PLC script executed successfully. Rows affected: %d", rowsAffected)
	return nil

}
