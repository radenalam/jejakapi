package internal

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

// SQLQuery menyimpan informasi tentang SQL query yang dijalankan
type SQLQuery struct {
	SQL      string        `json:"sql"`
	Duration time.Duration `json:"duration"`
	Rows     int64         `json:"rows_affected"`
	Error    string        `json:"error,omitempty"`
}

// ContextKey untuk menyimpan SQL queries dalam context
type ContextKey string

const SQLQueriesKey ContextKey = "sql_queries"

// CustomSQLLogger adalah custom logger untuk GORM yang menangkap SQL queries
type CustomSQLLogger struct {
	logger.Interface
}

// NewCustomSQLLogger membuat instance baru dari CustomSQLLogger
func NewCustomSQLLogger(baseLogger logger.Interface) *CustomSQLLogger {
	return &CustomSQLLogger{
		Interface: baseLogger,
	}
}

// Trace menangkap SQL query dan menyimpannya dalam context
func (l *CustomSQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// Panggil logger asli terlebih dahulu
	l.Interface.Trace(ctx, begin, fc, err)

	// Tangkap query info
	sql, rows := fc()
	duration := time.Since(begin)

	query := SQLQuery{
		SQL:      sql,
		Duration: duration,
		Rows:     rows,
	}

	if err != nil {
		query.Error = err.Error()
	}

	// Simpan ke context jika ada
	if queries, ok := ctx.Value(SQLQueriesKey).(*[]SQLQuery); ok {
		*queries = append(*queries, query)
	}
}
