package domain

import "time"

// AuditEntry records one admin-initiated action for compliance and review.
type AuditEntry struct {
	ID         string    `json:"id"`
	ActorID    string    `json:"actorId"`
	ActorEmail string    `json:"actorEmail"`
	Action     string    `json:"action"`
	TargetType string    `json:"targetType"`
	TargetID   string    `json:"targetId"`
	Meta       string    `json:"meta"`
	CreatedAt  time.Time `json:"createdAt"`
}
