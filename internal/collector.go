package internal

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Collector bertanggung jawab untuk mengumpulkan dan menyimpan log request
type Collector struct {
	db      *gorm.DB
	enabled bool
}

// NewCollector membuat instance baru dari Collector
func NewCollector(db *gorm.DB) *Collector {
	collector := &Collector{
		db:      db,
		enabled: true,
	}

	collector.AutoMigrate()

	return collector
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

		// Skip monitoring untuk routes JejakAPI sendiri
		if strings.HasPrefix(ctx.Path(), "/jejakapi") {
			return ctx.Next()
		}

		start := time.Now()

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
		logEntry := RequestLog{
			ID:              uuid.New(),
			Method:          ctx.Method(),
			URL:             ctx.OriginalURL(),
			Headers:         JSON(requestHeaders),
			Body:            requestBody,
			StatusCode:      ctx.Response().StatusCode(),
			ResponseHeaders: JSON(responseHeaders),
			ResponseBody:    responseBody,
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
