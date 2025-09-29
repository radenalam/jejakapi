# JejakAPI

JejakAPI adalah tools monitoring request dan response seperti Laravel Telescope untuk aplikasi Go dengan Fiber framework.

## Features

- 📊 Real-time monitoring request dan response
- 🔍 Web interface untuk melihat history logs
- 📝 Detail informasi request/response headers dan body
- ⏱️ Durasi response time
- 🎯 Filter berdasarkan method HTTP
- 🗑️ Clear logs functionality
- 🔄 Auto-refresh logs


## Installation

Install JejakAPI dengan perintah berikut:

```bash
go get github.com/radenalam/jejakapi
```

Lalu integrasikan dengan aplikasi Fiber Anda sesuai contoh di bawah.

## Usage

### Basic Setup

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/radenalam/jejakapi"
    "gorm.io/gorm"
)

func main() {
    app := fiber.New()

    // Setup database connection
    db := setupDatabase() // your database setup

    // Initialize JejakAPI
    jejakAPI := jejakapi.New(db)

    // Add middleware untuk monitoring (sebelum routes lain)
    app.Use(jejakAPI.Middleware())

    // Setup routes untuk web interface
    jejakAPI.SetupRoutes(app)

    // Your application routes
    app.Get("/api/users", getUsersHandler)
    app.Post("/api/users", createUserHandler)

    app.Listen(":3000")
}
```

### Accessing Web Interface

Setelah setup, Anda bisa mengakses web interface di:

```
http://localhost:3000/jejakapi
```

### API Endpoints

JejakAPI menyediakan beberapa API endpoints:

- `GET /jejakapi/api/logs` - Ambil list logs dengan pagination
- `GET /jejakapi/api/logs/:id` - Detail log berdasarkan ID
- `DELETE /jejakapi/api/logs` - Hapus semua logs
- `POST /jejakapi/api/toggle` - Toggle monitoring on/off

## Web Interface Features

### Dashboard

- List semua request logs dengan informasi:
  - HTTP Method (GET, POST, PUT, DELETE, etc.)
  - URL endpoint
  - Status code dengan color coding
  - Response time
  - Timestamp

### Log Details

Klik pada log entry untuk melihat detail:

- Request headers
- Request body (jika ada)
- Response headers
- Response body (jika ada)
- Client information (IP, User Agent)

### Controls

- **Toggle Monitoring**: Enable/disable logging
- **Clear All Logs**: Hapus semua log entries
- **Refresh**: Manual refresh logs
- **Auto-refresh**: Otomatis refresh setiap 10 detik

## Database Schema

JejakAPI menggunakan table `request_logs` dengan structure:

```sql
CREATE TABLE request_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    method VARCHAR(10) NOT NULL,
    url VARCHAR(500) NOT NULL,
    headers JSONB,
    body TEXT,
    status_code INTEGER NOT NULL,
    response_headers JSONB,
    response_body TEXT,
    duration BIGINT NOT NULL, -- in microseconds
    ip VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```


## Performance Considerations

- Logs disimpan secara asynchronous untuk tidak mempengaruhi response time
- Gunakan index database untuk query yang lebih cepat
- Pertimbangkan untuk membersihkan logs lama secara berkala
- Untuk production, batasi ukuran response body yang disimpan

## Todo / Future Enhancements

- [ ] Filter berdasarkan status code
- [ ] Search functionality
- [ ] Export logs (JSON, CSV)
- [ ] Retention policy untuk auto-cleanup logs lama
- [ ] Alerting untuk error rates
- [ ] Performance metrics dan charts
- [ ] Authentication untuk web interface
- [ ] Configuration file support
