package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-kusto-go/kusto"
	"github.com/Azure/azure-kusto-go/kusto/ingest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	kustoTable          = "StormEvents"
	kustoMappingRefName = "StormEvents_CSV_Mapping"

	createTableCommand = ".create table StormEvents (StartTime: datetime, EndTime: datetime, EpisodeId: int, EventId: int, State: string, EventType: string, InjuriesDirect: int, InjuriesIndirect: int, DeathsDirect: int, DeathsIndirect: int, DamageProperty: int, DamageCrops: int, Source: string, BeginLocation: string, EndLocation: string, BeginLat: real, BeginLon: real, EndLat: real, EndLon: real, EpisodeNarrative: string, EventNarrative: string, StormSummary: dynamic)"

	createMappingCommand = `.create table StormEvents ingestion csv mapping 'StormEvents_CSV_Mapping' '[{"Name":"StartTime","datatype":"datetime","Ordinal":0}, {"Name":"EndTime","datatype":"datetime","Ordinal":1},{"Name":"EpisodeId","datatype":"int","Ordinal":2},{"Name":"EventId","datatype":"int","Ordinal":3},{"Name":"State","datatype":"string","Ordinal":4},{"Name":"EventType","datatype":"string","Ordinal":5},{"Name":"InjuriesDirect","datatype":"int","Ordinal":6},{"Name":"InjuriesIndirect","datatype":"int","Ordinal":7},{"Name":"DeathsDirect","datatype":"int","Ordinal":8},{"Name":"DeathsIndirect","datatype":"int","Ordinal":9},{"Name":"DamageProperty","datatype":"int","Ordinal":10},{"Name":"DamageCrops","datatype":"int","Ordinal":11},{"Name":"Source","datatype":"string","Ordinal":12},{"Name":"BeginLocation","datatype":"string","Ordinal":13},{"Name":"EndLocation","datatype":"string","Ordinal":14},{"Name":"BeginLat","datatype":"real","Ordinal":16},{"Name":"BeginLon","datatype":"real","Ordinal":17},{"Name":"EndLat","datatype":"real","Ordinal":18},{"Name":"EndLon","datatype":"real","Ordinal":19},{"Name":"EpisodeNarrative","datatype":"string","Ordinal":20},{"Name":"EventNarrative","datatype":"string","Ordinal":21},{"Name":"StormSummary","datatype":"dynamic","Ordinal":22}]'`

	blobStorePathFormat  = "https://%s.blob.core.windows.net/%s/%s%s"
	blobStoreAccountName = "kustosamplefiles"
	blobStoreContainer   = "samplefiles"
	blobStoreToken       = "?st=2018-08-31T22%3A02%3A25Z&se=2020-09-01T22%3A02%3A00Z&sp=r&sv=2018-03-28&sr=b&sig=LQIbomcKI8Ooz425hWtjeq6d61uEaq21UVX7YrM61N4%3D"
	blobStoreFileName    = "StormEvents.csv"

	dropTableQ = ".drop table StormEvents"

	defaultDBName = "testkustodb"
)

var kc *kusto.Client

func init() {
	kc = connect()
}

func main() {

	kustoDB := os.Getenv("KUSTO_DB")
	if kustoDB == "" {
		log.Fatal("missing KUSTO_DB environment variable")
	}

	log.Println("Using database -", kustoDB)

	dropTable(kc, kustoDB)
	createTable(kc, kustoDB)
	createMapping(kc, kustoDB)
	ingestFile(kc, blobStoreAccountName, blobStoreContainer, blobStoreToken, blobStoreFileName, kustoMappingRefName, kustoDB, kustoTable)
}

func dropTable(kc *kusto.Client, kustoDB string) {
	_, err := kc.Mgmt(context.Background(), kustoDB, kusto.NewStmt(dropTableQ))
	if err != nil {
		log.Println("Failed to drop StormEvents table. Maybe it does not exist?")
		return
	}
	log.Println("Table StormEvents dropped from DB", kustoDB)
}

func createTable(kc *kusto.Client, kustoDB string) {
	_, err := kc.Mgmt(context.Background(), kustoDB, kusto.NewStmt(createTableCommand))
	if err != nil {
		log.Fatal("failed to create table", err)
	}

	log.Printf("Table %s created in DB %s\n", kustoTable, kustoDB)
}

func createMapping(kc *kusto.Client, kustoDB string) {
	_, err := kc.Mgmt(context.Background(), kustoDB, kusto.NewStmt(createMappingCommand))
	if err != nil {
		log.Fatal("failed to create mapping - ", err)
	}
	log.Printf("Mapping %s created\n", kustoMappingRefName)
}

func ingestFile(kc *kusto.Client, blobStoreAccountName, blobStoreContainer, blobStoreToken, blobStoreFileName, kustoMappingRefName, kustoDB, kustoTable string) {
	kIngest, err := ingest.New(kc, kustoDB, kustoTable)
	if err != nil {
		log.Fatal("failed to create ingestion client", err)
	}
	blobStorePath := fmt.Sprintf(blobStorePathFormat, blobStoreAccountName, blobStoreContainer, blobStoreFileName, blobStoreToken)
	err = kIngest.FromFile(context.Background(), blobStorePath, ingest.FileFormat(ingest.CSV), ingest.IngestionMappingRef(kustoMappingRefName, ingest.CSV))

	if err != nil {
		log.Fatal("failed to ingest file", err)
	}
	log.Println("Ingested file from -", blobStorePath)
}

func connect() *kusto.Client {
	clientID := os.Getenv("AZURE_SP_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_SP_CLIENT_SECRET")
	tenantID := os.Getenv("AZURE_SP_TENANT_ID")
	kustoEndpoint := os.Getenv("KUSTO_ENDPOINT")

	if clientID == "" || clientSecret == "" || tenantID == "" || kustoEndpoint == "" {
		log.Fatal("missing one or more environment variables: AZURE_SP_CLIENT_ID, AZURE_SP_CLIENT_SECRET, AZURE_SP_TENANT_ID, KUSTO_ENDPOINT")
	}

	auth := kusto.Authorization{Config: auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)}
	dataxClient, err := kusto.New(kustoEndpoint, auth)
	if err != nil {
		log.Fatal("failed to create kusto client", err)
	}

	log.Println("Connected to Azure Data Explorer")
	return dataxClient
}
