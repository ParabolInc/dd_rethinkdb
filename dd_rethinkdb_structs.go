package main

type queryEngine struct {
	ClientConnections float64 `gorethink:"client_connections"`
	ClientsActive     float64 `gorethink:"clients_active"`
	QueriesPerSec     float64 `gorethink:"queries_per_sec"`
	QueriesTotal      float64 `gorethink:"queries_total"`
	ReadDocsPerSec    float64 `gorethink:"read_docs_per_sec"`
	ReadDocsTotal     float64 `gorethink:"read_docs_total"`
	WrittenDocsPerSec float64 `gorethink:"written_docs_per_sec"`
	WrittenDocsTotal  float64 `gorethink:"written_docs_total"`
}

type storageEngine struct {
	Cache struct {
		InUseBytes float64 `gorethink:"in_use_bytes"`
	}
	Disk struct {
		ReadBytesPerSec    float64 `gorethink:"read_bytes_per_sec"`
		ReadBytesTotal     float64 `gorethink:"read_bytes_total"`
		WrittenBytesPerSec float64 `gorethink:"written_bytes_per_sec"`
		WrittenBytesTotal  float64 `gorethink:"written_bytes_total"`
		SpaceUsage         struct {
			DataBytes         float64 `gorethink:"data_bytes"`
			MetadataBytes     float64 `gorethink:"metadata_bytes"`
			GarbageBytes      float64 `gorethink:"garbage_bytes"`
			PreallocatedBytes float64 `gorethink:"preallocated_bytes"`
		} `gorethink:"space_usage"`
	}
}

type Stat struct {
	ID            []string      `gorethink:"id"`
	QueryEngine   queryEngine   `gorethink:"query_engine,omitempty" `
	StorageEngine storageEngine `gorethink:"storage_engine,omitempty" `
	Server        string        `gorethink:"server,omitempty" `
	DB            string        `gorethink:"db,omitempty" `
	Table         string        `gorethink:"table,omitempty" `
	Error         string        `gorethink:"error,omitempty" `
}
