package main

type queryEngine struct {
	ClientConnections float64 `gorethink:"client_connections"      exporter:"cluster,server"`
	ClientsActive     float64 `gorethink:"clients_active"          exporter:"cluster,server"`
	QueriesPerSec     float64 `gorethink:"queries_per_sec"         exporter:"cluster,server"`
	QueriesTotal      float64 `gorethink:"queries_total"           exporter:"server"`
	ReadDocsPerSec    float64 `gorethink:"read_docs_per_sec"       exporter:"all"`
	ReadDocsTotal     float64 `gorethink:"read_docs_total"         exporter:"server,table_server"`
	WrittenDocsPerSec float64 `gorethink:"written_docs_per_sec"    exporter:"all"`
	WrittenDocsTotal  float64 `gorethink:"written_docs_total"      exporter:"server,table"`
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

	Server string `gorethink:"server,omitempty" `
	DB     string `gorethink:"db,omitempty" `
	Table  string `gorethink:"table,omitempty" `

	Error string `gorethink:"error,omitempty" `
}
