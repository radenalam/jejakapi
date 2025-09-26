package internal

import (
	_ "embed"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

//go:embed views/dashboard.html
var dashboardHTML []byte

//go:embed views/detail.html
var detailHTML []byte

// WebHandler menangani routing untuk web interface jejakapi
type WebHandler struct {
	collector *Collector
}

// NewWebHandler membuat instance baru dari WebHandler
func NewWebHandler(collector *Collector) *WebHandler {
	return &WebHandler{
		collector: collector,
	}
}

// SetupRoutes mensetup routing untuk jejakapi web interface
func (h *WebHandler) SetupRoutes(app *fiber.App) {
	jejakapi := app.Group("/jejakapi")

	// Web interface
	jejakapi.Get("/", h.dashboard)
	jejakapi.Get("/logs/:id", h.logDetail)

	// API endpoints
	api := jejakapi.Group("/api")
	api.Get("/logs", h.getLogs)
	api.Get("/logs/:id", h.getLogByID)
	api.Delete("/logs", h.clearLogs)
	api.Post("/toggle", h.toggleEnabled)
}

// dashboard menampilkan halaman dashboard utama
func (h *WebHandler) dashboard(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return c.Send(dashboardHTML)
}

// getLogs mengembalikan list log dalam format JSON
func (h *WebHandler) getLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	logs, total, err := h.collector.GetLogs(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch logs",
		})
	}

	return c.JSON(fiber.Map{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// getLogByID mengembalikan detail log berdasarkan ID
func (h *WebHandler) getLogByID(c *fiber.Ctx) error {
	id := c.Params("id")

	log, err := h.collector.GetLogByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Log not found",
		})
	}

	return c.JSON(log)
}

// clearLogs menghapus semua log
func (h *WebHandler) clearLogs(c *fiber.Ctx) error {
	err := h.collector.ClearLogs()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to clear logs",
		})
	}

	return c.JSON(fiber.Map{
		"message": "All logs cleared successfully",
	})
}

// toggleEnabled mengaktifkan/menonaktifkan monitoring
func (h *WebHandler) toggleEnabled(c *fiber.Ctx) error {
	// Untuk saat ini sederhana, nanti bisa disimpan di database atau config
	return c.JSON(fiber.Map{
		"message": "Monitoring toggled",
	})
}

// logDetail menampilkan detail log
func (h *WebHandler) logDetail(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return c.Send(detailHTML)
}
