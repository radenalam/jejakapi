package jejakapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radenalam/jejakapi/internal"
	"gorm.io/gorm"
)

// JejakAPI adalah struct utama untuk tools monitoring
type JejakAPI struct {
	collector  *internal.Collector
	webHandler *internal.WebHandler
}

// New membuat instance baru dari JejakAPI
func New(db *gorm.DB) *JejakAPI {
	// Auto migrate table
	db.AutoMigrate(&internal.RequestLog{})

	collector := internal.NewCollector(db)
	webHandler := internal.NewWebHandler(collector)

	return &JejakAPI{
		collector:  collector,
		webHandler: webHandler,
	}
}

// Middleware mengembalikan Fiber middleware untuk monitoring
func (j *JejakAPI) Middleware() fiber.Handler {
	return j.collector.Middleware()
}

// SetupRoutes mensetup routing untuk web interface
func (j *JejakAPI) SetupRoutes(app *fiber.App) {
	j.webHandler.SetupRoutes(app)
}

// GetCollector mengembalikan instance collector
func (j *JejakAPI) GetCollector() *internal.Collector {
	return j.collector
}
