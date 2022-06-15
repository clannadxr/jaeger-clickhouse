package clickhousedependencystore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/jaegertracing/jaeger-clickhouse/storage/clickhousespanstore/mocks"
)

const (
	testDependenciesTable = "test_dependencies_table"
)

func TestDependencyReader_GetDependencies(t *testing.T) {
	db, mock, err := mocks.GetDbMock()
	require.NoError(t, err, "an error was not expected when opening a stub database connection")
	defer db.Close()
	dependencyReader := NewDependencyReader(db, testDependenciesTable)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"timestamp", "call_count", "parent", "child"}).
		AddRow(time.Now(), 100, "s1", "s2")
	mock.ExpectQuery(fmt.Sprintf(`SELECT timestamp,max(call_count) as call_count,parent,child 
FROM %s 
WHERE timestamp <= ? and timestamp >= ? 
GROUP BY timestamp, parent, child`, testDependenciesTable)).WillReturnRows(rows)
	dependencies, err := dependencyReader.GetDependencies(context.Background(), now, time.Hour)
	require.NoError(t, err)
	require.Equal(t, uint64(100), dependencies[0].CallCount)
}
