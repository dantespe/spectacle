# Spectacle
Good Performance, Great Charts 

## Documentation 

To view the documentation for this repository, you can install `godoc`:
```
$ go install golang.org/x/tools/cmd/godoc@latest
```

Then you can start a localhost server and navigate to [localhost:6080](http://localhost:6080/pkg/github.com/dantespe/spectacle/):
```
$ godoc --http=:6080
```

## REST Reference

### Overview
| Endpoint                                 | Description                                       | Method |
| ---------------------------------------- | ------------------------------------------------- | ------ |
| [`/status`](#status)                     | The status of the server.                         | `GET`  |
| [`/datasets`](#list-datasets)            | Returns all datasets.                             | `GET`  |
| [`/datasets/<datasetId>`](#get-dataset)  | Returns a single dataset.                         | `GET`  |
| [`/datasets/<datasetId>/headers`](#get-headers)  | Returns a single dataset.                         | `GET`  |
| [`/dataset`](#create-dataset)            | Creates a new dataset                             | `POST` |
| [`/dataset/<datasetId>/upload`](#upload) | Uploads a new file to the dataset with datasetId. | `POST` |


#### [Status](#status)

The status of the server.

Example:
```
curl localhost:8080/rest/status
{
   "code" : 200,
   "numDatasets" : 0,
   "numRecords" : 0,
   "status" : "HEALTHY"
}
```

#### [List Datasets](#list-datasets)

Returns a list of datasets.

Example:
```
curl localhost:8080/rest/datasets
{
   "code" : 200,
   "results" : [
      {
         "datasetId" : 1,
         "displayName" : "large-dataset",
         "headersSet" : true,
         "numRecords" : 1048577
      },
      {
         "datasetId" : 2,
         "displayName" : "teams",
         "headersSet" : true,
         "numRecords" : 31
      },
      {
         "datasetId" : 3,
         "displayName" : "ranking",
         "headersSet" : true,
         "numRecords" : 210343
      },
   ],
   "totalDatasets" : 3
}
```

#### [Get Dataset](#get-dataset)

Returns the dataset with given dataset id.

Example:
```
curl localhost:8080/rest/dataset/1
{
   "code" : 200,
   "dataset" : {
      "datasetId" : 1,
      "displayName" : "teams",
      "headersSet" : true,
      "numRecords" : 31
   }
}
```

#### [Get Headers](#get-headers)

Returns the headers with given dataset id.

Example:
```
curl localhost:8080/rest/dataset/1/headers
{
   "code" : 200,
   "results" : [
      {
         "displayName" : "LEAGUE_ID",
         "headerId" : 1
      },
      {
         "displayName" : "TEAM_ID",
         "headerId" : 2
      },
      {
         "displayName" : "MIN_YEAR",
         "headerId" : 3
      },
      {
         "displayName" : "MAX_YEAR",
         "headerId" : 4
      },
    ...
      {
         "displayName" : "HEADCOACH",
         "headerId" : 13
      },
      {
         "displayName" : "DLEAGUEAFFILIATION",
         "headerId" : 14
      }
   ]
}
```

#### [Create Dataset](#create-dataset)

Creates an empty dataset.

**Options:**

* `displayName`: the name of the dataset. If unset will be "untitled-{datasetId}"

Examples:
```
curl -X POST localhost:8080/rest/dataset
{
   "code" : 201,
   "datasetId" : 8,
   "datasetUrl" : "/dataset/8",
   "displayName" : "untitled-8"
}

curl -X POST -d '{"displayName": "secret-dataset"}' -H "Content-Type: application/json" localhost:8080/rest/dataset
{
   "code" : 201,
   "datasetId" : 9,
   "datasetUrl" : "/dataset/9",
   "displayName" : "secret-dataset"
}
```

#### [Upload](#upload)

Upload a CSV into a dataset. 

Example:

This uploads the file `./data/top_1000.csv` into dataset with `id=9`
```
curl -X POST -F "file=@./data/top_1000.csv"  -H "Content-Type: multipart/form-data" localhost:8080/rest/dataset/9/upload
{
   "code" : 200,
   "operation" : "/operation/8"
}
```
