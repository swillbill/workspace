package formsmigration

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type sourceDocuments struct {
// 	jsontown []bson.M `json:"jsonTown"`
// }

type Config struct {
	ConnectionString string
	Database         string
	Collection       string
}

type Ids struct {
	Ids []string `json:"ids"`
}

type Jsonids struct {
	SourceID Ids `json:"sourceID"`
	TargetID Ids `json:"targetID"`
}

type MongoJson struct {
	MongoJson map[string]bson.M `json:"mongoJson"`
}

type MongoJsonWithID struct {
	MongoJson     bson.M `json:"mongoJson"`
	SourceID      string `json:"sourceID"`
	TargetID      string `json:"targetID"`
	FormGroupName string `json:"formGroupName"`
}

var configs = map[string]Config{
	"dev": {
		ConnectionString: os.Getenv("MONGO_CONN_DEV"),
		Database:         "NgForms",
		Collection:       "FormGroups",
	},
	"ref": {
		ConnectionString: os.Getenv("MONGO_CONN_REF"),
		Database:         "ref_NgForms",
		Collection:       "FormGroups",
	},
	"stg": {
		ConnectionString: os.Getenv("MONGO_CONN_STG"),
		Database:         "stg_NgForms",
		Collection:       "FormGroups",
	},
	"qa": {
		ConnectionString: os.Getenv("MONGO_CONN_QA"),
		Database:         "NgForms",
		Collection:       "FormGroups",
	},
	"stp": {
		ConnectionString: os.Getenv("MONGO_CONN_STP"),
		Database:         "stp_NgForms",
		Collection:       "FormGroups",
	},
	"stl": {
		ConnectionString: os.Getenv("MONGO_CONN_STL"),
		Database:         "stl_NgForms",
		Collection:       "FormGroups",
	},
	"devgov": {
		ConnectionString: os.Getenv("MONGO_CONN_DEVGOV"),
		Database:         "NgForms",
		Collection:       "FormGroups",
	},
	"qagov": {
		ConnectionString: os.Getenv("MONGO_CONN_QAGOV"),
		Database:         "NgForms",
		Collection:       "FormGroups",
	},
	"refgov": {
		ConnectionString: os.Getenv("MONGO_CONN_REFGOV"),
		Database:         "NgForms",
		Collection:       "FormGroups",
	},
	"demogov": {
		ConnectionString: os.Getenv("MONGO_CONN_DEMOGOV"),
		Database:         "NgForms",
		Collection:       "FormGroups",
	},
}

var debug *bool

// gets env
func GetConfig(env string) Config {
	if config, ok := configs[strings.ToLower(env)]; ok {
		return config
	}
	log.Fatalf("Unknown environment, %s, Beginning Script ", env)
	return Config{}
}

func ConnectToMongoDB(env string, config Config) (*mongo.Client, string) {
	clientOptions := options.Client().ApplyURI(config.ConnectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Printf("Connected to %s database in %s environment.\n", config.Database, env)
	return client, config.Database
}

// reads from Namelist.txt file
func ReadFormGroupNames(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
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
		return nil, fmt.Errorf("failed to read lines from file: %w", err)
	}

	return formGroups, nil
}

// validate form group names in Mongo
func ValidateFormGroupNames(client *mongo.Client, database string, formGroups []string) error {
	log.Printf("Validating form group names against MongoDB collection")
	collection := client.Database(database).Collection("FormGroups")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error fetching documents from MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	// Return error if no docs are found
	if !cursor.Next(ctx) {
		log.Println("Cursor is empty, no documents found in the collection.")
		return fmt.Errorf("no documents found in the MongoDB collection")
	}

	formGroupsSet := make(map[string]bool)
	for _, formGroupName := range formGroups {
		formGroupsSet[formGroupName] = false // Default to not matched
	}

	log.Println("Starting iteration...")

	// Iterate over Mongo docs
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return fmt.Errorf("error decoding document: %w", err)
		}

		// Get formGroupName
		dbFormGroupName, ok := result["formGroupName"].(string)
		if !ok {
			log.Println("Warning: 'formGroupName' field is missing or not a string in the document.")
			continue
		}

		// Check if dbFormGroupName exists in formGroupsSet and hasn't been matched yet
		if matched, exists := formGroupsSet[dbFormGroupName]; exists && !matched {
			log.Printf("Form name '%s' matches from list\n", dbFormGroupName)
			formGroupsSet[dbFormGroupName] = true
		}
	}

	// Log unmatched form group names
	for formGroupName, matched := range formGroupsSet {
		if !matched {
			log.Printf("Form name '%s' from the list was not found in the database.\n", formGroupName)
		}
	}

	log.Println("Finished iterating over MongoDB collection.")
	return nil
}

func exportAllDocumentsAsJSON(client *mongo.Client, config Config, formGroups []string, sourceTenantID string) ([]MongoJsonWithID, error) {
	collection := client.Database(config.Database).Collection(config.Collection)
	ValidateFormGroupNames(client, config.Database, formGroups)
	filter := bson.M{"formGroupName": bson.M{"$in": formGroups}, "tenantID": sourceTenantID}
	Debugging(fmt.Sprintf("ExportDocumentsAsJSON: Filter: %+v", filter))

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return []MongoJsonWithID{}, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(context.Background())

	allMongoWithIDs := []MongoJsonWithID{}
	var documents = make(map[string]bson.M)

	// Iterate though docs and find highest version
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode document: %v", err)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to find documents: %w", err)
		}
		if documents[result["formGroupName"].(string)]["updateDate"] != nil { // if existing Result has a updatedDate
			Debugtown(fmt.Sprintf("existing result has an updatedDate: %s", result["updatedDate"]))
			if result["updatedDate"] != nil { //Incoming result has a updatedDate
				Debugtown(fmt.Sprintf("incoming result has an updatedDate: %s", documents[result["formGroupName"].(string)]))
				if documents[result["formGroupName"].(string)]["updatedDate"].(primitive.DateTime) < result["updatedDate"].(primitive.DateTime) { //if incoming result is newer than existing
					Debugtown(fmt.Sprintf("Incoming result is newer than existing: old: %s, new: %s", documents[result["formGroupName"].(string)]["updatedDate"].(primitive.DateTime).Time().String(), result["updatedDate"].(primitive.DateTime).Time().String()))
					documents[result["formGroupName"].(string)] = result
				} else {
					Debugtown(fmt.Sprintf("Incoming result is older than existing: old: %s, new: %s", documents[result["formGroupName"].(string)]["updatedDate"].(primitive.DateTime).Time().String(), result["updatedDate"]))
				}
			} else { //if existing result does not have an updatedDate
				Debugtown("Incoming result does not have an update date. Ignoring.")
			}
		} else { //existing doesn't have one.
			documents[result["formGroupName"].(string)] = result
		}
		Debugging(fmt.Sprintf("ExportDocumentsAsJSON: Cursor.next: Working on: %s", result["_id"]))
	}

	for _, doc := range documents {
		allMongoWithIDs = append(allMongoWithIDs, MongoJsonWithID{MongoJson: doc, SourceID: doc["_id"].(string), TargetID: "", FormGroupName: doc["formGroupName"].(string)})
	}

	file, err := os.OpenFile("source_documents.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(allMongoWithIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to write to JSON file: %w", err)
	}

	log.Printf("Processed %d documents.", len(allMongoWithIDs))

	if len(allMongoWithIDs) == 0 {
		return nil, fmt.Errorf("no documents were processed, ensure the filter is correct")
	}

	return allMongoWithIDs, nil

}

// Gets docs from source DB and saves as JSON
func ExportDocumentsAsJSON(client *mongo.Client, config Config, formGroups []string, sourceTenantID string) (map[string]bson.M, error) {
	collection := client.Database(config.Database).Collection(config.Collection)
	filter := bson.M{"formGroupName": bson.M{"$in": formGroups}, "tenantID": sourceTenantID}
	Debugging(fmt.Sprintf("ExportDocumentsAsJSON: Filter: %+v", filter))

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(context.Background())

	if len(formGroups) > 1000 {
		log.Printf("File contains too many form groups (%d). Consider splitting the file.", len(formGroups))
	}

	var documents []bson.M

	// Iterate though docs and find highest version
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode document: %v", err)
			continue
		}
		Debugging(fmt.Sprintf("ExportDocumentsAsJSON: Cursor.next: Working on: %s", result["_id"]))
		documents = append(documents, result)
	}

	// Iterate through the docs and find the newest updated date EX // "updatedDate" : {"$date" : 1718040248255}
	groupedResultsUpdatedDate := make(map[string]bson.M)

	for _, doc := range documents {
		formGroupName := doc["formGroupName"].(string)

		// Get the updatedDate as a primitive.DateTime
		updatedDate, ok := doc["updatedDate"].(primitive.DateTime)
		if !ok {
			Debugtown(fmt.Sprintf("Document missing updatedDate field: %s", doc["formGroupName"]))
		}

		latestUpdateDate := updatedDate.Time().UnixMilli()

		// Check if the formGroupName already exists in the map
		if existingDoc, exists := groupedResultsUpdatedDate[formGroupName]; exists {
			existingUpdateDateBSON, ok := existingDoc["updatedDate"].(primitive.DateTime)
			if !ok {
				log.Fatalf("Document has invalid updatedDate field: %+v", existingDoc)
			}

			existingUpdateDate := existingUpdateDateBSON.Time().UnixMilli()
			log.Printf("Existing formGroupName: %s, existingUpdateDate: %d", formGroupName, existingUpdateDate)

			// Only update if the existing date is older
			if existingUpdateDate < latestUpdateDate {
				Debugging(fmt.Sprintf("ExportDocumentsAsJson: Timestamp Check: Replacing %s with %s", existingDoc["_id"], doc["_id"]))
				groupedResultsUpdatedDate[formGroupName] = doc
				//log.Printf("Updated formGroupName: %s with new latestUpdateDate: %d", formGroupName, latestUpdateDate)
				//log.Printf("Taking form with latest updatedDate: %v", updatedDate.Time())
			} else {
				Debugtown(fmt.Sprintf("Keeping existing formGroupName: %s with existingUpdateDate: %d", formGroupName, existingUpdateDate))
			}
		} else {
			// Add new formGroupName if it doesn't exist
			groupedResultsUpdatedDate[formGroupName] = doc
			log.Printf("Added new formGroupName: %s", formGroupName)
			log.Printf("Taking form with updatedDate: %v", updatedDate.Time())
		}
	}

	// After processing all documents, groupedResultsUpdatedDate will have the latest form for each formGroupName

	// Open the file for writing
	file, err := os.OpenFile("source_documents.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(MongoJson{MongoJson: groupedResultsUpdatedDate})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to write to JSON file: %w", err)
	}

	log.Printf("Processed %d documents.", len(groupedResultsUpdatedDate))

	if len(groupedResultsUpdatedDate) == 0 {
		return nil, fmt.Errorf("no documents were processed, ensure the filter is correct")
	}

	log.Printf("JSON export completed successfully.")
	return groupedResultsUpdatedDate, nil
}

// updates documents with _id from targetDB and saves as JSON
func UpdateDocumentsWithTargetIDs(targetClient *mongo.Client, targetConfig Config, targetTenantID string) ([]MongoJsonWithID, error) {
	targetCollection := targetClient.Database(targetConfig.Database).Collection(targetConfig.Collection)

	updatedResults := make(map[string]bson.M)

	file, err := os.Open("source_documents.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	fileResults, err := os.ReadFile("source_documents.json")
	if err != nil {
		log.Fatal("Couldn't open source_documents.json!")
	}

	var sourceResultsJson []MongoJsonWithID
	var newResultsJson []MongoJsonWithID
	err = json.Unmarshal(fileResults, &sourceResultsJson)
	if err != nil {
		log.Fatalf("Failed to unmarshal source documents: %v", err)
	}

	Debugtown(fmt.Sprintf("UpdateDocumentsWithTargetIDs: sourceResultsJson: %+v", sourceResultsJson))

	for _, sourceDoc := range sourceResultsJson {
		filter := bson.M{"formGroupName": sourceDoc.FormGroupName, "tenantID": targetTenantID}
		var targetDocData bson.M
		err := targetCollection.FindOne(context.Background(), filter).Decode(&targetDocData)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No matching document found in target for formGroupName: %s\n", sourceDoc.FormGroupName)
			sourceDoc.MongoJson["_id"] = primitive.NewObjectID().Hex()
			Debugging(fmt.Sprintf("Setting %s to ID %s.", sourceDoc.FormGroupName, sourceDoc.MongoJson["_id"]))
			sourceDoc.MongoJson["tenantID"] = targetTenantID
			updatedResults[sourceDoc.FormGroupName] = sourceDoc.MongoJson
			sourceDoc.TargetID = sourceDoc.MongoJson["_id"].(string)
			newResultsJson = append(newResultsJson, sourceDoc)
		} else {
			existingID, err := getIDAsString(targetDocData)
			if err != nil {
				return nil, fmt.Errorf("failed to find document for formGroupName %s: %w", sourceDoc.FormGroupName, err)
			}
			Debugtown(fmt.Sprintf("Updating formGroupName %s: setting _id to %s\n", sourceDoc.FormGroupName, existingID))
			sourceDoc.MongoJson["_id"] = existingID
			Debugging(fmt.Sprintf("Setting %s to ID %s.", sourceDoc.FormGroupName, existingID))
			sourceDoc.MongoJson["tenantID"] = targetTenantID
			updatedResults[sourceDoc.FormGroupName] = sourceDoc.MongoJson
			sourceDoc.TargetID = sourceDoc.MongoJson["_id"].(string)
			newResultsJson = append(newResultsJson, sourceDoc)
		}

	}

	//updatedResultsStruct := MongoJson{MongoJson: updatedResults}
	//fmt.Printf("Updated documents: %+v\n", updatedResultsStruct)
	for _, v := range newResultsJson {
		Debugging(fmt.Sprintf("the new target ID is: %s", v.TargetID))
	}

	jsonData, err := json.Marshal(newResultsJson)
	//fmt.Printf("jsonData: %+v\n", jsonData)
	if err != nil {
		log.Fatalf("failed to marshal data to JSON: %v", err)
	}
	if err := os.WriteFile("target_documents.json", jsonData, 0644); err != nil {
		log.Fatalf("failed to write JSON file: %v", err)
	}
	//if err := saveToJSONFile("target_documents.json", updatedResultsStruct); err != nil {
	//	return nil, err
	//}

	fmt.Println("Updated documents have been saved to target_documents.json.")
	return newResultsJson, nil
}

// helper to convert _id from various types to string from objectID
func getIDAsString(doc bson.M) (string, error) {
	if id, ok := doc["_id"].(primitive.ObjectID); ok {
		return id.Hex(), nil
	}
	if idStr, ok := doc["_id"].(string); ok {
		return idStr, nil
	}
	return "", fmt.Errorf("unsupported _id type")
}

// deletes docs from the targetDB
func DeleteExistingDocuments(targetClient *mongo.Client, targetConfig Config, documents []MongoJsonWithID) error {
	collection := targetClient.Database(targetConfig.Database).Collection(targetConfig.Collection)
	for _, doc := range documents {
		if id, ok := doc.MongoJson["_id"].(string); ok {
			if err := deleteDocumentByID(collection, id); err != nil {
				return err
			}
		} else if id, ok := doc.MongoJson["_id"].(primitive.ObjectID); ok {
			if err := deleteDocumentByID(collection, id.Hex()); err != nil {
				return err
			}
		}
	}
	return nil
}

// deletes a document from the collection by _id
func deleteDocumentByID(collection *mongo.Collection, id string) error {
	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete document with _id %s: %w", id, err)
	}
	fmt.Printf("Deleted document with _id: %s\n", id)
	return nil
}

func InsertDocumentsIntoTarget(targetClient *mongo.Client, targetConfig Config, documents []MongoJsonWithID) error {
	collection := targetClient.Database(targetConfig.Database).Collection(targetConfig.Collection)

	moduleIDUpdate := UpdateDoc(documents)

	for _, doc := range moduleIDUpdate {
		if _, err := collection.InsertOne(context.Background(), doc.MongoJson); err != nil {
			return fmt.Errorf("failed to insert document with _id %v: %w", doc.MongoJson["_id"], err)
		}
		fmt.Printf("Inserted document with _id: %v\n", doc.MongoJson["_id"])
	}
	return nil
}

// converting moduleID value
func UpdateDoc(jsonData []MongoJsonWithID) []MongoJsonWithID {
	// if float64
	for _, doc := range jsonData {
		switch doc.MongoJson["moduleId"].(type) {
		case float64:
			// Convert float64 to int
			doc.MongoJson["moduleId"] = int(doc.MongoJson["moduleId"].(float64))
		case string:
			// Try to convert string to int
			doc.MongoJson["moduleId"], _ = strconv.Atoi(doc.MongoJson["moduleId"].(string))
		}
	}
	return jsonData
}
