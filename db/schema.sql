CREATE TABLE IF NOT EXISTS Datasets(
    DatasetId SERIAL,
    DisplayName TEXT NOT NULL,
    HeadersSet INT NOT NULL,
    PRIMARY KEY (DatasetId)
);

CREATE TABLE IF NOT EXISTS Headers(
    HeaderId SERIAL,
    DatasetId INTEGER REFERENCES Datasets(DatasetId),
    ColumnIndex SERIAL,
    DisplayName TEXT,
    ValueType TEXT NOT NULL,
    PRIMARY KEY (HeaderId)
);

CREATE TABLE IF NOT EXISTS Records (
    RecordId SERIAL,
    DatasetId INTEGER REFERENCES Datasets(DatasetId),
    DatasetIndex INTEGER, -- UNUSED
    PRIMARY KEY (RecordId)
);

CREATE TABLE IF NOT EXISTS Operations (
    OperationId SERIAL,
    OperationStatus TEXT NOT NULL,
    ErrorMessage TEXT,
    CreationTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (OperationId)
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

