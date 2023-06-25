CREATE TABLE IF NOT EXISTS Datasets (
    DatasetId INTEGER PRIMARY KEY,
    DisplayName TEXT NOT NULL,
    HeadersSet INT NOT NULL
);

CREATE TABLE IF NOT EXISTS Headers (
    HeaderId INTEGER PRIMARY KEY,
    DatasetId INTEGER NOT NULL,
    ColumnIndex INTEGER, -- UNUSED
    DisplayName TEXT,
    ValueType TEXT NOT NULL,
    UNIQUE (DatasetId, ColumnIndex),
    FOREIGN KEY(DatasetId) REFERENCES Datasets(DatasetId)
);

CREATE TABLE IF NOT EXISTS Records (
    RecordId INTEGER PRIMARY KEY,
    DatasetId INTEGER NOT NULL,
    DatasetIndex INTEGER, -- UNUSED
    UNIQUE (DatasetId, DatasetIndex),
    FOREIGN KEY(DatasetId) REFERENCES Datasets(DatasetId)
);

CREATE TABLE IF NOT EXISTS Operations (
    OperationId INTEGER PRIMARY KEY,
    OperationStatus TEXT NOT NULL,
    ErrorMessage TEXT,
    CreationTime DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Cells (
    CellId INTEGER PRIMARY KEY,
    RecordId INTEGER NOT NULL,
    HeaderId INTEGER NOT NULL,
    OperationId INTEGER NOT NULL,
    RawValue TEXT,
    IntValue INTEGER,
    FloatValue FLOAT,
    UNIQUE (HeaderId, RecordId),
    FOREIGN KEY(RecordId) REFERENCES Records(RecordId),
    FOREIGN KEY(HeaderId) REFERENCES Headers(HeaderId),
    FOREIGN KEY(OperationId) REFERENCES Operation(OperationId)
);

