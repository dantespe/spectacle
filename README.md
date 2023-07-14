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
| Endpoint                                             | Description                                       | Method |
| ---------------------------------------------------- | ------------------------------------------------- | ------ |
| [`/rest/status`](#status)                            | The status of the server.                         | `GET`  |
| [`/rest/datasets`](#list-datasets)                   | Returns all datasets.                             | `GET`  |
| [`/rest/datasets/<datasetId>`](#get-dataset)         | Returns a single dataset.                         | `GET`  |
| [`/rest/datasets/<datasetId>/headers`](#get-headers) | Returns headers for a dataset.                    | `GET`  |
| [`/rest/data/<datasetId>`](#data-api)                | Returns data from a dataset.                      | `GET`  |
| [`/rest/dataset`](#create-dataset)                   | Creates a new dataset                             | `POST` |
| [`/rest/dataset/<datasetId>/upload`](#upload)        | Uploads a new file to the dataset with datasetId. | `POST` |


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

### [Data API](#data-api)

Returns the raw data from the dataset.

**Options:**

* `headers`: a comma-seperated list of header ids. Defaults to all headers in the dataset.

* `recordid`: the recordid that was last seen. Default is 0.

Example:
```
curl localhost:8080/rest/data/1
{
   "code" : 200,
   "results" : [
      {
         "displayName" : "ABBREVIATION",
         "headerId" : 5,
         "rows" : [
            "ATL",
            "BOS",
            "NOP",
            "CHI",
            "DAL",
            "DEN",
            "HOU",
            "LAC",
            "LAL",
            "MIA",
            "MIL",
            "MIN",
            "BKN",
            "NYK",
            "ORL",
            "IND",
            "PHI",
            "PHX",
            "POR",
            "SAC",
            "SAS",
            "OKC",
            "TOR",
            "UTA",
            "MEM",
            "WAS",
            "DET",
            "CHA",
            "CLE",
            "GSW"
         ]
      },
      {
         "displayName" : "NICKNAME",
         "headerId" : 6,
         "rows" : [
            "Hawks",
            "Celtics",
            "Pelicans",
            "Bulls",
            "Mavericks",
            "Nuggets",
            "Rockets",
            "Clippers",
            "Lakers",
            "Heat",
            "Bucks",
            "Timberwolves",
            "Nets",
            "Knicks",
            "Magic",
            "Pacers",
            "76ers",
            "Suns",
            "Trail Blazers",
            "Kings",
            "Spurs",
            "Thunder",
            "Raptors",
            "Jazz",
            "Grizzlies",
            "Wizards",
            "Pistons",
            "Hornets",
            "Cavaliers",
            "Warriors"
         ]
      },
   ]
}
```