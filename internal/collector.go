package internal

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Collector bertanggung jawab untuk mengumpulkan dan menyimpan log request
type Collector struct {
	db           *gorm.DB
	enabled      bool
	excludePaths []string
}

// NewCollector membuat instance baru dari Collector dengan SQL query monitoring
func NewCollector(db *gorm.DB) *Collector {
	collector := &Collector{
		db:      db,
		enabled: true,
		excludePaths: []string{
			"/jejakapi",
			"/favicon.ico",
			"/health",
		},
	}

	// Setup SQL logger otomatis
	collector.setupDefaultSQLLogger()

	collector.AutoMigrate()

	return collector
}

// setupDefaultSQLLogger mengatur default SQL logger untuk monitoring
func (c *Collector) setupDefaultSQLLogger() {
	// Buat default logger untuk SQL monitoring
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	// Setup custom SQL logger
	c.SetupSQLLogger(newLogger)
}

// SetupSQLLogger mengatur custom SQL logger untuk GORM
func (c *Collector) SetupSQLLogger(baseLogger logger.Interface) {
	customLogger := NewCustomSQLLogger(baseLogger)
	c.db.Logger = customLogger
}

// AutoMigrate membuat table jika belum ada
func (c *Collector) AutoMigrate() error {
	return c.db.AutoMigrate(&RequestLog{})
}

// Middleware adalah Fiber middleware untuk jejakapi
func (c *Collector) Middleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !c.enabled {
			return ctx.Next()
		}

		// Skip monitoring untuk exclude paths
		if c.shouldExcludePath(ctx.Path()) {
			return ctx.Next()
		}

		start := time.Now()

		// Siapkan slice untuk menampung SQL queries
		var sqlQueries []SQLQuery

		// Tambahkan SQL queries ke context
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), SQLQueriesKey, &sqlQueries))

		// Capture request body
		var requestBody *string
		if ctx.Body() != nil && len(ctx.Body()) > 0 {
			body := string(ctx.Body())
			requestBody = &body
		}

		// Capture request headers
		requestHeaders := make(map[string]interface{})
		ctx.GetReqHeaders()
		for key, values := range ctx.GetReqHeaders() {
			if len(values) > 0 {
				requestHeaders[key] = values[0]
			}
		}

		// Continue to next handler
		err := ctx.Next()

		// Calculate duration
		duration := time.Since(start).Microseconds()

		// Capture response
		var responseBody *string
		if len(ctx.Response().Body()) > 0 {
			body := string(ctx.Response().Body())
			responseBody = &body
		}

		// Capture response headers
		responseHeaders := make(map[string]interface{})
		ctx.Response().Header.VisitAll(func(key, value []byte) {
			responseHeaders[string(key)] = string(value)
		})

		// Create log entry
		var sqlQueriesJSON JSON
		if len(sqlQueries) > 0 {
			sqlQueriesMap := make(map[string]interface{})
			sqlQueriesMap["queries"] = sqlQueries
			sqlQueriesJSON = JSON(sqlQueriesMap)
		}

		logEntry := RequestLog{
			ID:              uuid.New(),
			Method:          ctx.Method(),
			URL:             ctx.OriginalURL(),
			Headers:         JSON(requestHeaders),
			Body:            requestBody,
			StatusCode:      ctx.Response().StatusCode(),
			ResponseHeaders: JSON(responseHeaders),
			ResponseBody:    responseBody,
			SQLQueries:      sqlQueriesJSON,
			Duration:        duration,
			IP:              ctx.IP(),
			UserAgent:       ctx.Get("User-Agent"),
			CreatedAt:       time.Now(),
		}

		// Save to database (async)
		go c.saveLog(logEntry)

		return err
	}
}

// saveLog menyimpan log ke database
func (c *Collector) saveLog(log RequestLog) {
	c.db.Create(&log)
}

// GetLogs mengambil list log dengan pagination
func (c *Collector) GetLogs(page, limit int) ([]RequestLog, int64, error) {
	var logs []RequestLog
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := c.db.Model(&RequestLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get logs with pagination - hanya mengambil field yang diperlukan untuk list
	if err := c.db.Select("id, method, url, status_code, duration, created_at").Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLogByID mengambil log berdasarkan ID
func (c *Collector) GetLogByID(id string) (*RequestLog, error) {
	var log RequestLog
	if err := c.db.Where("id = ?", id).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// ClearLogs menghapus semua log
func (c *Collector) ClearLogs() error {
	return c.db.Where("1 = 1").Delete(&RequestLog{}).Error
}

// SetEnabled mengaktifkan/menonaktifkan collector
func (c *Collector) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// shouldExcludePath memeriksa apakah path harus diexclude dari monitoring
func (c *Collector) shouldExcludePath(path string) bool {
	for _, excludePath := range c.excludePaths {
		if strings.HasPrefix(path, excludePath) {
			return true
		}
	}
	return false
}

// AddExcludePath menambahkan path ke exclude list
func (c *Collector) AddExcludePath(path string) {
	c.excludePaths = append(c.excludePaths, path)
}

// RemoveExcludePath menghapus path dari exclude list
func (c *Collector) RemoveExcludePath(path string) {
	for i, excludePath := range c.excludePaths {
		if excludePath == path {
			c.excludePaths = append(c.excludePaths[:i], c.excludePaths[i+1:]...)
			break
		}
	}
}

// SetExcludePaths mengatur semua exclude paths
func (c *Collector) SetExcludePaths(paths []string) {
	c.excludePaths = paths
}

// GetExcludePaths mengambil semua exclude paths
func (c *Collector) GetExcludePaths() []string {
	return c.excludePaths
}
