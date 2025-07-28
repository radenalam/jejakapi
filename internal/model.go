package internal

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// RequestLog menyimpan log request dan response
type RequestLog struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Method          string    `json:"method" gorm:"size:10;not null"`
	URL             string    `json:"url" gorm:"size:500;not null"`
	Headers         JSON      `json:"headers" gorm:"type:jsonb"`
	Body            *string   `json:"body,omitempty" gorm:"type:text"`
	StatusCode      int       `json:"status_code" gorm:"not null"`
	ResponseHeaders JSON      `json:"response_headers" gorm:"type:jsonb"`
	ResponseBody    *string   `json:"response_body,omitempty" gorm:"type:text"`
	Duration        int64     `json:"duration" gorm:"not null"` // in microseconds
	IP              string    `json:"ip" gorm:"size:45"`
	UserAgent       string    `json:"user_agent" gorm:"size:500"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// JSON adalah custom type untuk handling JSON dalam GORM
type JSON map[string]interface{}

func (j JSON) Value() interface{} {
	if j == nil {
		return nil
	}
	b, _ := json.Marshal(j)
	return string(b)
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case []byte:
		return json.Unmarshal(v, j)
	}
	return nil
}
