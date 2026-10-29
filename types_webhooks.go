package lyngbak

import (
	"encoding/json"
	"time"
)

type Webhook struct {
	UUID              string    `json:"uuid"`
	Name              string    `json:"name"`
	URL               string    `json:"url"`
	Enabled           bool      `json:"enabled"`
	EventFilter       []string  `json:"event_filter"`
	AppUUIDs          []string  `json:"app_uuids"`
	IncludeReviewApps bool      `json:"include_review_apps"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type WebhookCreated struct {
	Webhook
	Secret string `json:"secret"`
}

type WebhookDelivery struct {
	UUID          string          `json:"uuid"`
	EventID       string          `json:"event_id"`
	EventType     string          `json:"event_type"`
	Status        string          `json:"status"`
	ResponseCode  *int            `json:"response_code"`
	Message       *string         `json:"message"`
	Attempts      int             `json:"attempts"`
	LastAttemptAt *time.Time      `json:"last_attempt_at"`
	CreatedAt     time.Time       `json:"created_at"`
	Payload       json.RawMessage `json:"payload"`
}

type CreateWebhookRequest struct {
	Name              string   `json:"name"`
	URL               string   `json:"url"`
	EventFilter       []string `json:"event_filter"`
	AppUUIDs          []string `json:"app_uuids,omitempty"`
	Enabled           *bool    `json:"enabled,omitempty"`
	IncludeReviewApps *bool    `json:"include_review_apps,omitempty"`
}

type UpdateWebhookRequest struct {
	Name              *string   `json:"name,omitempty"`
	URL               *string   `json:"url,omitempty"`
	EventFilter       *[]string `json:"event_filter,omitempty"`
	AppUUIDs          *[]string `json:"app_uuids,omitempty"`
	Enabled           *bool     `json:"enabled,omitempty"`
	IncludeReviewApps *bool     `json:"include_review_apps,omitempty"`
}
