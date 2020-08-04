---
page_type: sample
languages:
- go
products:
- azure
description: "Sample to demonstrate how to ingest (load) data into Azure Data Explorer using the Go SDK"
urlFragment: cognitive-services-personalizer-samples
---

# Ingest data using the Azure Data Explorer Go SDK 

Azure Data Explorer is a fast and highly scalable data exploration service for log and telemetry data. It provides a [Go client library](https://docs.microsoft.com/azure/data-explorer/kusto/api/golang/kusto-golang-client-library?WT.mc_id=adxgo-docs-abhishgu) for interacting with the Azure Data Explorer service. You can use the [Go SDK](https://github.com/Azure/azure-kusto-go) to query, control, and ingest into Azure Data Explorer clusters

You will start by using the sample program to create a table and data mapping in a test cluster. You then queue an ingestion to the cluster using the Go SDK and validate the results.

## Getting Started

### Prerequisites

* If you don't have an Azure subscription, create a [free Azure account](https://azure.microsoft.com/free/?WT.mc_id=adxgo-docs-abhishgu) before you begin.

* You need [Go](https://golang.org/) installed on your computer. The [Go SDK requires Go 1.13](https://docs.microsoft.com/azure/data-explorer/kusto/api/golang/kusto-golang-client-library?WT.mc_id=adxgo-docs-abhishgu#minimum-requirements) as a minimum version

* [A test cluster and database](create-cluster-database-portal.md)

* You need [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) installed on your computer


### Installation

To install using [Go modules](https://github.com/golang/go/wiki/Modules), create a module (if needed) and fetch the Azure Data Explorer package. For example:


```shell
go mod init foo.com/bar
go get github.com/Azure/azure-kusto-go/kusto
```

You should see the package dependency being added to the `go.mod` file

### Azure Data Explorer resources

[Create an Azure Data Explorer cluster and database](https://docs.microsoft.com/azure/data-explorer/net-sdk-ingest-data?WT.mc_id=adxgo-docs-abhishgu#create-a-table-on-your-test-cluster) before executing the sample program in this repository


## Run the application

Clone the sample code:

```shell
git clone https://github.com/Azure-Samples/azure-data-explorer-go-ingest.git
cd azure-data-explorer-go-ingest
```

You will need to create Service Principal for authentication. For example, you can do it with Azure CLI using the [az ad sp create-for-rbac](https://docs.microsoft.com/cli/azure/ad/sp?view=azure-cli-latest&WT.mc_id=adxgo-docs-abhishgu#az-ad-sp-create-for-rbac) command

Set Service Principal information along with cluster endpoint in the form of environment variables that will be used by the program:

```shell
export AZURE_SP_CLIENT_ID="<replace with appID>"
export AZURE_SP_CLIENT_SECRET="<replace with password>"
export AZURE_SP_TENANT_ID="<replace with tenant>"
export KUSTO_ENDPOINT="https://<cluster name>.<azure region>.kusto.windows.net"
export KUSTO_DB="name of the database"
```

Finally, to run the program:

```shell
go run main.go
```

You will see an output similar to the following:

```shell
Connected to Azure Data Explorer
Using database - testkustodb
Failed to drop StormEvents table. Maybe it does not exist?
Table StormEvents created in DB testkustodb
Mapping StormEvents_CSV_Mapping created
Ingested file from - https://kustosamplefiles.blob.core.windows.net/samplefiles/StormEvents.csv?st=2018-08-31T22%3A02%3A25Z&se=2020-09-01T22%3A02%3A00Z&sp=r&sv=2018-03-28&sr=b&sig=LQIbomcKI8Ooz425hWtjeq6d61uEaq21UVX7YrM61N4%3D
```

## Resources

- [Azure Data Explorer](https://docs.microsoft.com/azure/data-explorer/?WT.mc_id=adxgo-docs-abhishgu)
- [Azure Data Explorer Go SDK](https://docs.microsoft.com/azure/data-explorer/kusto/api/golang/kusto-golang-client-library?WT.mc_id=adxgo-docs-abhishgu)
- [Azure Data Explorer Go SDK on GitHub](https://github.com/Azure/azure-kusto-go/)
- [Azure Data Explorer Go SDK reference](https://godoc.org/github.com/Azure/azure-kusto-go/kusto)
