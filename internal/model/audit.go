package model

import "time"

type AuditEvent struct {
	ID        int64     `json:"id"`
	Subject   string    `json:"subject"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}
