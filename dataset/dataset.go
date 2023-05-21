// Dataset contains logic for creating a collection of data.
package dataset

import (
    "fmt"
    "sync"
    "io"
    "encoding/csv"
)

// Dataset contains all logic for managing data in Spectacle.
type Dataset struct {
    // Id of the dataset.
    Id uint64 `json:"datasetId"`

    // DisplayName of Datast.
    DisplayName string `json:"displayName"`

    // Number of Records (Rows) in the dataset.
    NumRecords uint64 `json:"numRecords"`

    // HasHeaders means the data has named columns. We should
    // read the first column and use that row as the headers.
    // Default: true 
    HasHeaders bool    `json:"hasHeaders"`    
    headersSet bool
    headers []string

    // Data
    // this will also be stored in the database
    data [][]string

    // Maximum number of threads to run when importing data.
    // Default: 1000
    maxThreads int
    twg sync.WaitGroup

    // Lock
    mu sync.RWMutex
}

// Option for creating new Datasets.
type Option func(*Dataset)

// New creats a new Dataset with the given options.
func New(opts ...Option) (*Dataset, error) {
    ds, err := Default()
    if err != nil {
        return nil, err
    }
    for _, o := range opts {
        o(ds)
    }
    return ds, nil
}

// Returns an Option with the Id set.
func WithId(id uint64) Option {
    return func(ds *Dataset) {
        ds.Id = id
    }
}

// Returns an Option with the DisplayName set.
func WithDisplayName(displayName string) Option {
    return func(ds *Dataset) {
        ds.DisplayName = displayName
    }
}

// Returns an Option with HasHeaders set.
func WithHasHeaders(hasHeaders bool) Option {
    return func(ds *Dataset) {
        ds.HasHeaders = hasHeaders
    } 
}

// Returns an Option with maxImportThreads set.
func WithMaxImportThreads(n int) Option {
    return func(ds *Dataset) {
        ds.maxThreads = n
    }
}

// NewWithId builds a Dataset with the provided Id.
func NewWithId(id uint64) (*Dataset, error) {
    ds, _ := Default()
    ds.Id = id
    return ds, nil
}

// Default returns an initialized Dataset.
func Default() (*Dataset, error) {
    return &Dataset{
        maxThreads: 100,
        HasHeaders: true,
        headersSet: false,
        headers: make([]string, 0),
    }, nil
}

// Shallow copy of the dataset. This can be used
// for thread-safe printing of the Dataset.
func (d *Dataset) Copy() *Dataset {
    d.mu.RLock()
    defer d.mu.RUnlock()

    return &Dataset{
        Id: d.Id,
        DisplayName: d.DisplayName,
        NumRecords: d.NumRecords,
    }
}

// Equal returns a bool, string of if two datasets are equal and a diff.
func (d *Dataset) Equal(other *Dataset) (bool, string)  {
    d.mu.RLock()
    defer d.mu.RUnlock()
    other.mu.RLock()
    defer other.mu.RUnlock()

    if other == nil {
        return false, "other is nil"
    }
    if d.Id != other.Id {
        return false, fmt.Sprintf("Id: %d vs %d", d.Id, other.Id)
    }
    if d.DisplayName != other.DisplayName {
        return false, fmt.Sprintf("DisplayName: %s vs %s", d.DisplayName, other.DisplayName)
    }
    if d.HasHeaders != other.HasHeaders {
        return false, fmt.Sprintf("HasHeaders: %v vs %v", d.HasHeaders, other.HasHeaders)
    }
    if d.NumRecords != other.NumRecords {
        return false, fmt.Sprintf("NumRecords: %d vs %d", d.NumRecords, other.NumRecords)
    }
    return true, ""
}

func (d *Dataset) processHeader(reader *csv.Reader) error {
    d.mu.Lock()
    defer d.mu.Unlock()

    // We already have process the headers or this
    // dataset doesn't have headers
    if d.headersSet || !d.HasHeaders {
        return nil
    }

    record, err := reader.Read()
    if err != nil {
        return err
    }
    for _, r := range record {
        d.headers = append(d.headers, r)
    }
    d.headersSet = true
    return nil
}

// TODO make this do something
func (d *Dataset) addRecord(record []string) error {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.NumRecords++
    return nil
}

func (d *Dataset) Upload(r io.Reader) error {
    reader := csv.NewReader(r)
    d.processHeader(reader)
    wg := sync.WaitGroup{}
    for {
        record, err := reader.Read()
        // Unexpected Error
        if err != nil && err != io.EOF {
            return err
        }
        // EOF, wait until all threads have finished
        if err == io.EOF {
            break
        }

        // Synchronously read this record
        wg.Add(1)
        go func(){
            d.addRecord(record)
            wg.Done()
        }()
    }
    wg.Wait()
    return nil
}

// SetUntitledDisplayName sets the DisplayName to untitled-{n} if displayName is empty.
// Increments counter if the Dataset is untitled.
func (d *Dataset) SetUntitledDisplayName(u int) int {
    if d.DisplayName == "" {
        d.DisplayName = fmt.Sprintf("untitled-%d", u)
        return u + 1
    }
    return u
}