package audit

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/yourname/security-orchestrator/internal/db"
)

// AuditLog represents a single audit entry.
type AuditLog struct {
	ID         int64     `json:"id"`
	TS         time.Time `json:"ts"`
	UserID     sql.NullInt64 `json:"user_id,omitempty"`
	EngagementID sql.NullInt64 `json:"engagement_id,omitempty"`
	EventType  string    `json:"event_type"`
	ObjectType sql.NullString `json:"object_type,omitempty"`
	ObjectID   sql.NullInt64 `json:"object_id,omitempty"`
	Details    json.RawMessage `json:"details,omitempty"`
}

// EnsureTable creates the audit_log table if it does not exist.
func EnsureTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS audit_log (
		id BIGSERIAL PRIMARY KEY,
		ts TIMESTAMPTZ NOT NULL DEFAULT now(),
		user_id INT REFERENCES users(id),
		engagement_id INT REFERENCES engagements(id),
		event_type VARCHAR NOT NULL,
		object_type VARCHAR,
		object_id INT,
		details JSONB
	);
	CREATE INDEX IF NOT EXISTS idx_audit_user_time ON audit_log(user_id, ts DESC);
	CREATE INDEX IF NOT EXISTS idx_audit_eng_time ON audit_log(engagement_id, ts DESC);
	`
	_, err := db.Exec(query)
	return err
}

// Insert writes a new audit entry.
func Insert(db *sql.DB, userID, engagementID *int64, eventType, objectType *string, objectID *int64, details []byte) error {
	query := `
	INSERT INTO audit_log (user_id, engagement_id, event_type, object_type, object_id, details)
	VALUES ($1,$2,$3,$4,$5,$6)
	`
	_, err := db.Exec(query,
		userID, engagementID, eventType, objectType, objectID, details)
	return err
}
