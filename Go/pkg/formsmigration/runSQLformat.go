package formsmigration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSelectionSQLConfig(sourceEnv string, sourceUserName string, sourceDBPassword string) (Config, error) {

	var sqlconfigs = map[string]Config{
		"dev": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"qa": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"ref": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "ref_NgPlatform",
		},
		"stg": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"stl": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"stp": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"devgov": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"qagov": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"refgov": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
		"demogov": {
			ConnectionString: fmt.Sprintf("<redacted>", sourceUserName, sourceDBPassword),
			Database:         "NgPlatform",
		},
	}

	// Find the sqlconfig environment
	sqlconfig, ok := sqlconfigs[strings.ToLower(sourceEnv)]
	if !ok {
		return Config{}, fmt.Errorf("environment %s not found", sourceEnv)
	}
	return sqlconfig, nil
}

type LayoutSelectionConfigData struct {
	LayoutConfigurationId string `json:"layout_configuration_id"`
	ConfigurationId       string `json:"configuration_id"`
	CursorData            string `json:"cursordata"`
	Version               string `json:"version"`
	SortId                int    `json:"sort_id"`
}

type LayoutMigrationConfigData struct {
	LayoutType           string          `json:"layout_type"`
	Context              string          `json:"context"`
	Layout               string          `json:"layout"`
	Version              json.RawMessage `json:"version"`
	IsDeleted            bool            `json:"is_deleted"`
	ConfigurationId      string          `json:"configuration_id"`
	ConfigurationVersion int             `json:"configuration_version"`
}

// Get ID values from existing mongo dataset without querying the database
func GetIdValuesFromMongoDataSet(documents map[string]bson.M) []string {
	UniqueVersionedResults := make(map[string]bson.M)
	for _, doc := range documents {
		formGroupName := doc["formGroupName"].(string)
		incomingversion, err := strconv.Atoi(doc["version"].(string))
		if err != nil {
			log.Fatalf("Document missing version field: %+v", doc)
		}

		Debugtown(fmt.Sprintf("GetIdValuesFromMongoDataSet: Incoming Version: %d, working on %s, ID: %s", incomingversion, formGroupName, doc["_id"].(string)))

		if _, exists := UniqueVersionedResults[formGroupName]; !exists {
			Debugtown(fmt.Sprintf("GetIdValuesFromMongoDataSet: %s didn't exist. Adding to map!", formGroupName))
			UniqueVersionedResults[formGroupName] = doc
		}
		existingVersion, err := strconv.Atoi(UniqueVersionedResults[formGroupName]["version"].(string))
		if err != nil {
			log.Fatalf("Document missing version field: %+v", doc)
		}

		Debugtown(fmt.Sprintf("GetIdValuesFromMongoDataSet: existing Version %d, ID: %s....incoming is %d. Adding to map!", existingVersion, UniqueVersionedResults[formGroupName]["_id"].(string), incomingversion))
		if existingVersion < incomingversion {
			UniqueVersionedResults[formGroupName] = doc
		}
	}
	finalIDs := []string{}
	Debugtown(fmt.Sprintf("GetIdValuesFromMongoDataSet: final IDs is: %+v", finalIDs))
	for _, doc := range UniqueVersionedResults {
		finalIDs = append(finalIDs, doc["_id"].(string))
	}
	Debugtown(fmt.Sprintf("GetIdValuesFromMongoDataSet: final IDs is: %+v", finalIDs))
	return finalIDs
}

func GetIdValuesFromMongo(client *mongo.Client, config Config, formGroups []string, tenantID string) []string {
	log.Printf("Starting SQL Migrations for form groups: %v", formGroups)
	collection := client.Database(config.Database).Collection(config.Collection)

	// Query the source collection to get the _id values
	filter := bson.M{
		"formGroupName": bson.M{"$in": formGroups},
		"tenantID":      tenantID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error querying source collection: %v", err)
	}
	defer cursor.Close(context.Background())
	var documents []bson.M
	for cursor.Next(context.Background()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			log.Fatalf("Error decoding source document: %v", err)
		}
		documents = append(documents, document)

		if err := cursor.Err(); err != nil {
			log.Fatalf("Error occurred while iterating source cursor: %v", err)
		}
	}

	UniqueVersionedResults := make(map[string]bson.M)
	for _, doc := range documents {
		formGroupName := doc["formGroupName"].(string)
		incomingversion, err := strconv.Atoi(doc["version"].(string))
		if err != nil {
			log.Fatalf("Document missing version field: %+v", doc)
		}

		Debugtown(fmt.Sprintf("GetIDValuesFromMongo: Incoming Version: %d, working on %s, ID: %s", incomingversion, formGroupName, doc["_id"].(string)))

		if _, exists := UniqueVersionedResults[formGroupName]; !exists {
			Debugtown(fmt.Sprintf("GetIdValuesFromMongo: %s didn't exist. Adding to map!", formGroupName))
			UniqueVersionedResults[formGroupName] = doc
		}
		existingVersion, err := strconv.Atoi(UniqueVersionedResults[formGroupName]["version"].(string))
		if err != nil {
			log.Fatalf("Document missing version field: %+v", doc)
		}

		Debugtown(fmt.Sprintf("GetIdValuesFromMongo: existing Version %d, ID: %s....incoming is %d. Adding to map!", existingVersion, UniqueVersionedResults[formGroupName]["_id"].(string), incomingversion))
		if existingVersion < incomingversion {
			UniqueVersionedResults[formGroupName] = doc
		}
	}
	finalIDs := []string{}
	Debugtown(fmt.Sprintf("GetIdValuesFromMongo: final IDs is: %+v", finalIDs))
	for _, doc := range UniqueVersionedResults {
		finalIDs = append(finalIDs, doc["_id"].(string))
	}
	Debugtown(fmt.Sprintf("GetIdValuesFromMongo: final IDs is: %+v", finalIDs))
	return finalIDs
}

func RunSelectionScript(mongojsonstruct []MongoJsonWithID, tenantID string, sourceEnv string, sourceDB string, sourceUserName string, sourceDBPassword string) []LayoutSelectionConfigData {
	// getconfiguration based on the environment, source username, and source password	flags
	sqlconfig, ok := GetSelectionSQLConfig(sourceEnv, sourceUserName, sourceDBPassword)
	if ok != nil {
		log.Printf("Invalid environment: %s", sqlconfig)
		return nil
	}

	// connection string from the config for SQL Server connection
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

	// build SQL script with the _id values
	idValues := make([]string, len(mongojsonstruct))
	for i, docs := range mongojsonstruct {
		Debugging(fmt.Sprintf("RunSelectionScript: SourceID: %s, TargetID: %s", docs.SourceID, docs.TargetID))
		idValues[i] = fmt.Sprintf("('%s', '%s')", docs.SourceID, docs.TargetID)
	}
	idValuesStr := strings.Join(idValues, ",")

	sqlScript := fmt.Sprintf(`USE [%s];
	;WITH 
	CTE_IdTable AS
	(
		SELECT DISTINCT 
		SourceID, TargetID
		FROM
		(
			VALUES
			%s
		) AS IdTable (SourceID, TargetID)
	),
	CTE_MyTable_Sorted AS
	(
		SELECT p.LayoutConfigurationId, p.version, p.ConfigurationId
		, '('''+ISNULL(p.LayoutType,'') +''','''+ ISNULL(p.Context,'') +''','''+ ISNULL(REPLACE(p.Layout,'''',''''''),'') +''','''+  ISNULL(CAST(Version AS nvarchar(16)),'') +''','''+ ISNULL(CAST(p.IsDeleted AS nvarchar(16)),'')+''','''+ ISNULL(CAST(CTE_IdTable.TargetID AS nvarchar(100)),'') +''',''' +ISNULL(CAST(ConfigurationVersion AS nvarchar(30)),'')+''')' AS cursordata
		, ROW_NUMBER() OVER (PARTITION BY p.ConfigurationId, p.context ORDER BY p.version DESC) AS SortId
		FROM LayoutConfiguration p
		JOIN CTE_IdTable ON CTE_IdTable.SourceID = p.ConfigurationId
		WHERE p.TenantId = '%s'
	)
	SELECT LayoutConfigurationId, ConfigurationId, cursordata, version, SortId
	FROM CTE_MyTable_Sorted
	WHERE SortId = 1`, sqlconfig.Database, idValuesStr, tenantID)

	//fmt.Println("SQL Script: ", sqlScript)

	// execute the SQL script
	Debugtown(fmt.Sprintf("runSQLFormat: SQL Script: %s", sqlScript))
	rows, err := db.Query(sqlScript)
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}
	defer rows.Close()

	// open a file to write the JSON data
	file, err := os.Create("cursordata.json")
	if err != nil {
		log.Fatalf("Error creating cursordata.json file: %v", err)
	}
	defer file.Close()

	// create a slice to hold the results
	var results []LayoutSelectionConfigData
	var rownumbah int = 0
	// process the results
	for rows.Next() {
		rownumbah++
		var layoutConfigID, configID, cursorData, version string
		var sortID int
		if err := rows.Scan(&layoutConfigID, &configID, &cursorData, &version, &sortID); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		results = append(results, LayoutSelectionConfigData{
			LayoutConfigurationId: layoutConfigID,
			ConfigurationId:       configID,
			CursorData:            cursorData,
			Version:               version,
			SortId:                sortID,
		})
	}
	fmt.Printf("Number of rows: %d\n", rownumbah)
	if err := rows.Err(); err != nil {
		log.Fatalf("Error occurred while iterating rows: %v", err)
	}

	// marshal the results into JSON format
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing data to JSON file: %v", err)
	}

	log.Println("SQL Selection Script Successfully Run and Data Written to cursordata.json")
	return results
}
