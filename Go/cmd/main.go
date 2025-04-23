package main

import (
	"flag"
	"fmt"
	"log"

	// import code packages
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/activedictionary"
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/channelmanagement"
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/correspondence"
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/formsmigration"
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/platformconfiguration"
	"github.com/revenue-solutions-inc/DevOps-CICD/Scripts/Go/pkg/workflow"
)

func main() {
	sourceEnv := flag.String("sourceEnv", "", "Source environment")
	targetEnv := flag.String("targetEnv", "", "Target environment")
	sourceTenantID := flag.String("sourceTenantID", "", "Source tenant ID")
	targetTenantID := flag.String("targetTenantID", "", "Target tenant ID")
	fileName := flag.String("file", "NameList.txt", "File name containing form group names")
	sourceDB := flag.String("sourceDB", "", "Source SQL DataBase")
	sourceUserName := flag.String("sourceUserName", "", "Source DB SQL Username")
	targetUserName := flag.String("targetUserName", "", "Target DB SQL Username")
	sourceDBPassword := flag.String("sourceDBPassword", "", "Target SQL Database Password")
	targetDBPassword := flag.String("targetDBPassword", "", "Target SQL Database Password")
	workFlow := flag.String("workflow", "nothing", "Which workflow to run, EX: mongo, sql, activeDictionary")
	module := flag.String("module", "", "Channel Management Module")
	channelId := flag.String("channelID", "", "Channel Management Channel ID")
	channelName := flag.String("channelName", "", "Channel Management Channel Name")
	flag.Parse()

	var source = false
	if *sourceEnv != "" {
		source = true
	}

	if source {
		switch *workFlow {
		case "mongo":
			fmt.Println("Starting Mongo Source Find and Export!")
			formsmigration.DoSourceMongoStuff(sourceEnv, sourceTenantID, fileName, sourceDB, sourceUserName, sourceDBPassword)
		case "sql":
			fmt.Println("Executing SQL Selection Query!")
			formsmigration.DoSourceSQLStuff(sourceEnv, sourceTenantID, fileName, sourceDB, sourceUserName, sourceDBPassword)
		case "activedictionary":
			fmt.Println("Executing Source Active Directory SQL Selection Query!")
			activedictionary.DoSourceADSQLStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword)
		case "workFlow":
			fmt.Println("Executing Source Workflow SQL Selection Query!")
			workflow.DoSourceWorkflowSQLStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword)
		case "plc":
			fmt.Println("Executing Source PLC SQL Selection Query!")
			platformconfiguration.DoSourcePLCStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword)
		case "corr":
			fmt.Println("Executing Source Correspondence SQL Selection Query!")
			correspondence.DoSourceCorrespondenceStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword)
		case "ct":
			fmt.Println("Executing Source Correspondence Template SQL Selection Query!")
			correspondence.DoSourceCorrespondenceTemplateStuff(sourceEnv, sourceTenantID, sourceUserName, sourceDBPassword)
		case "channel":
			fmt.Println("Executing Source Channel Management Migration!")
			channnelmanagement.DoChannelManagementSourceStuff(sourceEnv, module, channelId, channelName)
		case "nothing":
			log.Fatal("No workflow specified.")
		default:
			log.Fatal("What even is this?")
		}
	} else {
		switch *workFlow {
		case "mongo":
			fmt.Println("Starting Mongo Target Document Build!")
			formsmigration.DoTargetMongoStuff(targetEnv, targetTenantID, fileName, targetUserName, targetDBPassword, sourceTenantID)
		case "sql":
			fmt.Println("Executing Target SQL Query!")
			formsmigration.DoTargetSQLStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword)
		case "activedictionary":
			fmt.Println("Executing Target Active Directory SQL Selection Query!")
			activedictionary.DoTargetADSQLStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword)
		case "workFlow":
			fmt.Println("Executing Target Workflow SQL Selection Query!")
			workflow.DoTargetWorkflowSQLStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword)
		case "plc":
			fmt.Println("Executing Target PLC SQL Selection Query!")
			platformconfiguration.DoTargetPLCStuff(targetEnv, targetTenantID, targetUserName, targetDBPassword)
		case "corr":
			fmt.Println("Executing Target Correspondence SQL Selection Query!")
			correspondence.DoTargetCorrespondenceCommercialStuff(targetEnv, sourceEnv, sourceTenantID, targetTenantID, targetUserName, targetDBPassword)
		case "ct":
			fmt.Println("Executing Target Correspondence Template SQL Selection Query!")
			correspondence.DoTargetCorrespondenceTemplateCommercialStuff(targetEnv, sourceEnv, sourceTenantID, targetTenantID, targetUserName, targetDBPassword)
		// case "channel":
		// 	fmt.Println("Executing Target Channel Management Migration!")
		// 	channnelmanagement.DoChannelManagementTargetMigration(targetEnv, module, channelID, channelName)
		case "nothing":
			log.Fatal("No workflow specified.")
		default:
			log.Fatal("What even is this?")
		}
	}
}
