package formsmigration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

func GetMigrationSQLConfig(targetEnv string, targetUserName string, targetDBPassword string) (Config, error) {

	var sqlconfigs = map[string]Config{
		"dev": {
			ConnectionString: fmt.Sprintf("Server=tcp:dev-revx-managed-sql-instance.public.5f8be3c9b863.database.windows.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"qa": {
			ConnectionString: fmt.Sprintf("Server=tcp:qa-revx-managed-sql-instance.public.e5d2b6f8aec1.database.windows.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"ref": {
			ConnectionString: fmt.Sprintf("Server=tcp:rfsb-revx-managed-sql-instance.public.96a031c3633f.database.windows.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "ref_NgPlatform",
		},
		"stg": {
			ConnectionString: fmt.Sprintf("Server=tcp:revx-va-stg-smi.public.76b6315c2934.database.usgovcloudapi.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"stl": {
			ConnectionString: fmt.Sprintf("Server=tcp:va-01-stl-smi.public.f190e0b62306.database.usgovcloudapi.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"stp": {
			ConnectionString: fmt.Sprintf("Server=tcp:va-01-stp-smi.public.0ee6689d0f92.database.usgovcloudapi.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"devgov": {
			ConnectionString: fmt.Sprintf("Server=tcp:va-02-dev-smi.public.cea638793093.database.usgovcloudapi.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"qagov": {
			ConnectionString: fmt.Sprintf("Server=tcp:va-qal-smi.public.d88afbb0de5d.database.usgovcloudapi.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"refgov": {
			ConnectionString: fmt.Sprintf("Server=tcp:va-ref-smi.886a9d260a75.database.usgovcloudapi.net,1433;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
		"demogov": {
			ConnectionString: fmt.Sprintf("Server=tcp:va-demo-smi.public.6498609a439c.database.usgovcloudapi.net,3342;Persist Security Info=False;User ID=%s;Password=%s;MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;", targetUserName, targetDBPassword),
			Database:         "NgPlatform",
		},
	}

	// Find the sqlconfig environment
	sqlconfig, ok := sqlconfigs[strings.ToLower(targetEnv)]
	if !ok {
		return Config{}, fmt.Errorf("environment %s not found", targetEnv)
	}
	return sqlconfig, nil
}

func RunMigrationScript(targetEnv, targetTenantId string, targetUserName string, targetDBPassword string, data []LayoutSelectionConfigData) (err error) {

	log.Printf("Beginning SQL Migration Script")

	sqlConfig, ok := GetMigrationSQLConfig(targetEnv, targetUserName, targetDBPassword)
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
	file, err := os.Open("cursordata.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, err := os.Stat("cursordata.json")
	if err != nil {
		log.Fatal(err)
	}
	if content.Size() == 0 {
		log.Fatal("Cursor Data file is empty")
	}

	finishedJSON := ""
	for numbah, rows := range data {
		if numbah == 0 {
			finishedJSON = rows.CursorData
		} else {
			finishedJSON = fmt.Sprintf("%s,%s", finishedJSON, rows.CursorData)
		}

	}

	sqlScript := fmt.Sprintf(`
        USE [%s];
        SET NOCOUNT ON;

        DECLARE 
            @TenantId varchar(100) = N'%s',
            @LayoutConfigurationId varchar(100),
            @CreateDate datetime = GETUTCDATE(),
            @CreatedBy nvarchar(450) = 'AdminUser',
            @ModifiedDate datetime = GETUTCDATE(),
            @ModifiedBy nvarchar(450) = 'AdminUser',
            @debug bit = 1, 
            @commit bit = 1;

        BEGIN TRY
            BEGIN TRANSACTION;

            PRINT '-----------------------------------------------------------';
            PRINT 'TenantID: '+@TenantId;
            PRINT '-----------------------------------------------------------';

            DECLARE 
                @LayoutType nvarchar(500),
                @Context nvarchar(500),
                @Layout nvarchar(MAX),
                @Version int,
                @IsDeleted bit,
                @ConfigurationId nvarchar(450),
                @ConfigurationVersion int,
                @MaxVersion int,
                @MaxConfigurationVersion int;

            DECLARE curList CURSOR FOR 
            SELECT * 
            FROM (VALUES %s) v (LayoutType, Context, Layout, Version, IsDeleted, ConfigurationId, ConfigurationVersion);

            OPEN curList;
            FETCH NEXT FROM curList INTO @LayoutType, @Context, @Layout, @Version, @IsDeleted, @ConfigurationId, @ConfigurationVersion;

            WHILE @@FETCH_STATUS = 0
            BEGIN
                SET @LayoutConfigurationId = NEWID();

                IF EXISTS (SELECT 1 FROM LayoutConfiguration WHERE ConfigurationId = @ConfigurationId AND TenantId = @TenantId AND Context = @Context)
                BEGIN
                    SET @MaxVersion = (SELECT MAX(Version) + 1 FROM LayoutConfiguration WHERE ConfigurationId = @ConfigurationId AND TenantId = @TenantId AND Context = @Context);
                    SET @MaxConfigurationVersion = (SELECT MAX(ConfigurationVersion) + 1 FROM LayoutConfiguration WHERE ConfigurationId = @ConfigurationId AND TenantId = @TenantId AND Context = @Context);

                    IF @debug = 1 PRINT 'Layout exists. Adding new version.';
                END
                ELSE
                BEGIN
                    SET @MaxVersion = @Version;
                    SET @MaxConfigurationVersion = @ConfigurationVersion;

                    IF @debug = 1 PRINT 'Layout does not exist. Inserting.';
                END

                INSERT INTO LayoutConfiguration (
                    LayoutConfigurationId, TenantId, LayoutType, Context, Layout, Version, IsDeleted, 
                    CreateDate, CreatedBy, ModifiedDate, ModifiedBy, ConfigurationId, ConfigurationVersion, CorrelationId
                )
                VALUES (
                    @LayoutConfigurationId, @TenantId, @LayoutType, @Context, @Layout, @MaxVersion, @IsDeleted, 
                    @CreateDate, @CreatedBy, @ModifiedDate, @ModifiedBy, @ConfigurationId, @MaxConfigurationVersion, NEWID()
                );

                IF @debug = 1 PRINT 'Inserted record, LayoutConfigurationId: ' + @LayoutConfigurationId + ' ConfigurationVersion: ' + CAST(@MaxConfigurationVersion AS varchar(10)) + ' Version: ' + CAST(@MaxVersion AS varchar(10)) + ' Context: ' + @Context;

                FETCH NEXT FROM curList INTO @LayoutType, @Context, @Layout, @Version, @IsDeleted, @ConfigurationId, @ConfigurationVersion;
            END

            CLOSE curList;
            DEALLOCATE curList;

            IF @commit = 1 
            BEGIN
                COMMIT TRANSACTION;
                PRINT '-----------------------------------------------------------';
                PRINT 'Committed!';
            END
            ELSE
            BEGIN
                PRINT '-----------------------------------------------------------';
                PRINT 'Rolled back!';
                ROLLBACK TRANSACTION;
            END

        END TRY
        BEGIN CATCH
            THROW;

            WHILE @@TRANCOUNT > 0
            BEGIN
                ROLLBACK TRANSACTION;
            END
        END CATCH
        `, sqlConfig.Database, targetTenantId, finishedJSON)

	//fmt.Printf("SQL Script: %s\n", sqlScript)

	//debugtown(fmt.Sprintf("runSQLMigration: SQL Script: %s", sqlScript))
	lol, err := db.Prepare(sqlScript)
	if err != nil {
		log.Fatalf("Error preparing SQL script: %v", err)
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

	log.Printf("SQL script executed successfully. Rows affected: %d", rowsAffected)
	return nil
}
