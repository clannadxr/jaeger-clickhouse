CREATE TABLE IF NOT EXISTS {{.DependenciesTable}}
{{if .Replication}}ON CLUSTER '{cluster}'{{end}}
(
    {{if .Multitenant -}}
    tenant    LowCardinality(String) CODEC (ZSTD(1)),
    {{- end -}}
    timestamp DateTime CODEC (Delta, ZSTD(1)),
    parent   String CODEC (ZSTD(1)),
    child   String CODEC (ZSTD(1)),
    call_count   UInt64 CODEC (ZSTD(1))
) ENGINE {{if .Replication}}ReplicatedReplacingMergeTree(call_count){{else}}ReplacingMergeTree(call_count){{end}}
    {{.TTLDependencies}}
    PARTITION BY (
        {{if .Multitenant -}}
        tenant,
        {{- end -}}
        toDate(timestamp)
    )
    ORDER BY (timestamp,parent,child)
    SETTINGS index_granularity = 1024
