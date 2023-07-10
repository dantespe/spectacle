CREATE TABLE IF NOT EXISTS Datasets(
    DatasetId SERIAL,
    DisplayName TEXT NOT NULL,
    HeadersSet INT NOT NULL,
    NumRecords INTEGER, 
    PRIMARY KEY (DatasetId)
);

CREATE TABLE IF NOT EXISTS Headers(
    HeaderId SERIAL,
    DatasetId INTEGER REFERENCES Datasets(DatasetId),
    ColumnIndex INTEGER,
    DisplayName TEXT,
    ValueType TEXT NOT NULL,
    PRIMARY KEY (HeaderId)
);

CREATE INDEX IF NOT EXISTS idx_datasetid_headers ON Headers(DatasetId);

CREATE TABLE IF NOT EXISTS Operations (
    OperationId SERIAL,
    OperationStatus TEXT NOT NULL,
    ErrorMessage TEXT,
    CreationTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FinishTime TIMESTAMP,
    PRIMARY KEY (OperationId)
);

CREATE TABLE IF NOT EXISTS Records (
    RecordId SERIAL,
    OperationId INTEGER REFERENCES Operations(OperationId),
    DatasetId INTEGER REFERENCES Datasets(DatasetId),
    DatasetIndex INTEGER,
    PRIMARY KEY (RecordId)
);

CREATE INDEX IF NOT EXISTS idx_datasetid_records ON Records(DatasetId);

CREATE TABLE IF NOT EXISTS RecordsProcessed (
    RecordId INTEGER REFERENCES Records(RecordId),
    DatasetId INTEGER REFERENCES Datasets(DatasetId),
    UNIQUE (RecordId, DatasetId)
);

CREATE TABLE IF NOT EXISTS Cells (
    CellId SERIAL,
    RecordId INTEGER REFERENCES Records(RecordId),
    HeaderId INTEGER REFERENCES Headers(HeaderId),
    OperationId INTEGER REFERENCES Operations(OperationId),
    RawValue TEXT,
    IntValue INTEGER,
    FloatValue FLOAT,
    UNIQUE (HeaderId, RecordId),
    PRIMARY KEY (CellId)
);

CREATE INDEX IF NOT EXISTS idx_datasetid_cells ON Cells(RecordId);

