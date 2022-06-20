package clickhousedependencystore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/dependencystore"

	"github.com/jaegertracing/jaeger-clickhouse/storage/clickhousespanstore"
)

type DependencyReader struct {
	db                *sql.DB
	dependenciesTable clickhousespanstore.TableName
}

var _ dependencystore.Reader = (*DependencyReader)(nil)

func NewDependencyReader(db *sql.DB, dependenciesTable clickhousespanstore.TableName) *DependencyReader {
	return &DependencyReader{
		db:                db,
		dependenciesTable: dependenciesTable,
	}
}

func (d *DependencyReader) GetDependencies(ctx context.Context, endTs time.Time, lookback time.Duration) ([]model.DependencyLink, error) {
	query := fmt.Sprintf(`SELECT timestamp,argMax(call_count, time) as call_count,parent,child 
FROM %s 
WHERE timestamp <= ? and timestamp >= ? 
GROUP BY timestamp, parent, child`, d.dependenciesTable)
	rows, err := d.db.QueryContext(ctx, query, endTs, endTs.Add(-lookback))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dependencyLinks []model.DependencyLink
	for rows.Next() {
		var ts time.Time
		var depLink = model.DependencyLink{}
		err = rows.Scan(&ts, &depLink.CallCount, &depLink.Parent, &depLink.Child)
		if err != nil {
			return nil, err
		}
		dependencyLinks = append(dependencyLinks, depLink)
	}
	return dependencyLinks, nil
}
