package internal

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>JejakAPI - Request Monitoring</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: #2c3e50;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .controls {
            padding: 20px;
            border-bottom: 1px solid #eee;
            display: flex;
            gap: 10px;
            align-items: center;
        }
        .btn {
            background: #3498db;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
        }
        .btn:hover {
            background: #2980b9;
        }
        .btn.danger {
            background: #e74c3c;
        }
        .btn.danger:hover {
            background: #c0392b;
        }
        .logs-container {
            padding: 20px;
        }
        .log-entry {
            border: 1px solid #ddd;
            margin-bottom: 10px;
            border-radius: 4px;
            overflow: hidden;
        }
        .log-header {
            background: #f8f9fa;
            padding: 10px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            cursor: pointer;
        }
        .log-header:hover {
            background: #e9ecef;
        }
        .method {
            font-weight: bold;
            padding: 2px 8px;
            border-radius: 3px;
            color: white;
            font-size: 12px;
        }
        .method.GET { background: #28a745; }
        .method.POST { background: #007bff; }
        .method.PUT { background: #ffc107; color: #000; }
        .method.DELETE { background: #dc3545; }
        .method.PATCH { background: #6f42c1; }
        .status {
            font-weight: bold;
            padding: 2px 8px;
            border-radius: 3px;
            color: white;
            font-size: 12px;
        }
        .status.success { background: #28a745; }
        .status.error { background: #dc3545; }
        .status.warning { background: #ffc107; color: #000; }
        .log-details {
            padding: 15px;
            background: #f8f9fa;
            display: none;
        }
        .detail-section {
            margin-bottom: 15px;
        }
        .detail-title {
            font-weight: bold;
            margin-bottom: 5px;
            color: #495057;
        }
        .detail-content {
            background: white;
            padding: 10px;
            border-radius: 3px;
            border: 1px solid #dee2e6;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            max-height: 200px;
            overflow-y: auto;
        }
        .pagination {
            text-align: center;
            padding: 20px;
        }
        .pagination a {
            display: inline-block;
            padding: 8px 12px;
            margin: 0 5px;
            background: #f8f9fa;
            border: 1px solid #dee2e6;
            border-radius: 4px;
            text-decoration: none;
            color: #495057;
        }
        .pagination a:hover {
            background: #e9ecef;
        }
        .pagination a.active {
            background: #007bff;
            color: white;
            border-color: #007bff;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 JejakAPI</h1>
            <p>Request & Response Monitoring Tool</p>
        </div>
        
        <div class="controls">
            <button class="btn" onclick="toggleEnabled()">Toggle Monitoring</button>
            <button class="btn danger" onclick="clearLogs()">Clear All Logs</button>
            <button class="btn" onclick="refreshLogs()">Refresh</button>
        </div>
        
        <div class="logs-container" id="logs-container">
            <div style="text-align: center; padding: 40px;">
                <p>Loading logs...</p>
            </div>
        </div>
    </div>

    <script>
        let currentPage = 1;
        const logsPerPage = 20;

        function formatTimestamp(timestamp) {
            return new Date(timestamp).toLocaleString();
        }

        function formatDuration(microseconds) {
            if (microseconds < 1000) {
                return microseconds + 'μs';
            } else if (microseconds < 1000000) {
                return Math.round(microseconds / 1000) + 'ms';
            } else {
                return Math.round(microseconds / 1000000) + 's';
            }
        }

        function getStatusClass(status) {
            if (status >= 200 && status < 300) return 'success';
            if (status >= 400) return 'error';
            if (status >= 300) return 'warning';
            return '';
        }

        function formatJSON(str) {
            try {
                return JSON.stringify(JSON.parse(str), null, 2);
            } catch (e) {
                return str;
            }
        }

        function toggleLogDetails(id) {
            const details = document.getElementById('details-' + id);
            if (details.style.display === 'none' || !details.style.display) {
                details.style.display = 'block';
            } else {
                details.style.display = 'none';
            }
        }

        async function loadLogs(page = 1) {
            try {
                const response = await fetch('/jejakapi/api/logs?page=' + page + '&limit=' + logsPerPage);
                const data = await response.json();
                
                let html = '';
                
                if (data.logs && data.logs.length > 0) {
                    data.logs.forEach(log => {
                        const statusClass = getStatusClass(log.status_code);
                        html += '<div class="log-entry">';
                        html += '<div class="log-header" onclick="toggleLogDetails(\'' + log.id + '\')">';
                        html += '<div>';
                        html += '<span class="method ' + log.method + '">' + log.method + '</span>';
                        html += '<span style="margin-left: 10px;">' + log.url + '</span>';
                        html += '</div>';
                        html += '<div>';
                        html += '<span class="status ' + statusClass + '">' + log.status_code + '</span>';
                        html += '<span style="margin-left: 10px;">' + formatDuration(log.duration) + '</span>';
                        html += '<span style="margin-left: 10px; color: #6c757d;">' + formatTimestamp(log.created_at) + '</span>';
                        html += '</div>';
                        html += '</div>';
                        html += '<div class="log-details" id="details-' + log.id + '">';
                        html += '<div class="detail-section">';
                        html += '<div class="detail-title">Request Headers</div>';
                        html += '<div class="detail-content">' + formatJSON(JSON.stringify(log.headers)) + '</div>';
                        html += '</div>';
                        if (log.body) {
                            html += '<div class="detail-section">';
                            html += '<div class="detail-title">Request Body</div>';
                            html += '<div class="detail-content">' + formatJSON(log.body) + '</div>';
                            html += '</div>';
                        }
                        html += '<div class="detail-section">';
                        html += '<div class="detail-title">Response Headers</div>';
                        html += '<div class="detail-content">' + formatJSON(JSON.stringify(log.response_headers)) + '</div>';
                        html += '</div>';
                        if (log.response_body) {
                            html += '<div class="detail-section">';
                            html += '<div class="detail-title">Response Body</div>';
                            html += '<div class="detail-content">' + formatJSON(log.response_body) + '</div>';
                            html += '</div>';
                        }
                        html += '<div class="detail-section">';
                        html += '<div class="detail-title">Client Info</div>';
                        html += '<div class="detail-content">IP: ' + log.ip + '<br>User Agent: ' + log.user_agent + '</div>';
                        html += '</div>';
                        html += '</div>';
                        html += '</div>';
                    });
                    
                    // Add pagination
                    const totalPages = Math.ceil(data.total / logsPerPage);
                    if (totalPages > 1) {
                        html += '<div class="pagination">';
                        for (let i = 1; i <= totalPages; i++) {
                            const activeClass = i === page ? 'active' : '';
                            html += '<a href="#" class="' + activeClass + '" onclick="loadLogs(' + i + '); return false;">' + i + '</a>';
                        }
                        html += '</div>';
                    }
                } else {
                    html = '<div style="text-align: center; padding: 40px;"><p>No logs found</p></div>';
                }
                
                document.getElementById('logs-container').innerHTML = html;
                currentPage = page;
            } catch (error) {
                console.error('Error loading logs:', error);
                document.getElementById('logs-container').innerHTML = '<div style="text-align: center; padding: 40px;"><p>Error loading logs</p></div>';
            }
        }

        async function clearLogs() {
            if (confirm('Are you sure you want to clear all logs?')) {
                try {
                    await fetch('/jejakapi/api/logs', { method: 'DELETE' });
                    loadLogs(1);
                } catch (error) {
                    console.error('Error clearing logs:', error);
                }
            }
        }

        async function toggleEnabled() {
            try {
                await fetch('/jejakapi/api/toggle', { method: 'POST' });
                alert('Monitoring toggled');
            } catch (error) {
                console.error('Error toggling monitoring:', error);
            }
        }

        function refreshLogs() {
            loadLogs(currentPage);
        }

        // Load logs on page load
        document.addEventListener('DOMContentLoaded', function() {
            loadLogs();
            
            // Auto refresh every 10 seconds
            setInterval(function() {
                if (currentPage === 1) {
                    loadLogs(1);
                }
            }, 10000);
        });
    </script>
</body>
</html>`
	c.Set("Content-Type", "text/html")
	return c.SendString(html)
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

// logDetail menampilkan detail log (untuk implementasi detail view nantinya)
func (h *WebHandler) logDetail(c *fiber.Ctx) error {
	return c.Redirect("/jejakapi/")
}
