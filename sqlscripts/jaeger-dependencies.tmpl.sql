CREATE TABLE IF NOT EXISTS {{.DependenciesTable}}
{{if .Replication}}ON CLUSTER '{cluster}'{{end}}
(
    {{if .Multitenant -}}
    tenant    LowCardinality(String) CODEC (ZSTD(1)),
    {{- end -}}
    timestamp DateTime CODEC (Delta, ZSTD(1)),
    parent   String CODEC (ZSTD(1)),
    child   String CODEC (ZSTD(1)),
    call_count   UInt64 CODEC (ZSTD(1)),
    server_duration_p50 Float64 CODEC (ZSTD(1)),
    server_duration_p90 Float64 CODEC (ZSTD(1)),
    server_duration_p99 Float64 CODEC (ZSTD(1)),
    client_duration_p50 Float64 CODEC (ZSTD(1)),
    client_duration_p90 Float64 CODEC (ZSTD(1)),
    client_duration_p99 Float64 CODEC (ZSTD(1)),
    server_success_rate Float64 CODEC (ZSTD(1)),
    client_success_rate Float64 CODEC (ZSTD(1)),
    time DateTime CODEC (Delta, ZSTD(1))
) ENGINE {{if .Replication}}ReplicatedReplacingMergeTree(time){{else}}ReplacingMergeTree(time){{end}}
    {{.TTLDependencies}}
    PARTITION BY (
        {{if .Multitenant -}}
        tenant,
        {{- end -}}
        toDate(timestamp)
    )
    ORDER BY (timestamp,parent,child)
    SETTINGS index_granularity = 1024
