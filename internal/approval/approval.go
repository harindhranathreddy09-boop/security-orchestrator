package approval

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yourname/security-orchestrator/internal/db"
)

// Approval represents a pending or decided request.
type Approval struct {
	ID          int64     `json:"id"`
	EngagementID int64    `json:"engagement_id"`
	RequesterID int64    `json:"requester_id"`
	ApproverID  sql.NullInt64 `json:"approver_id,omitempty"`
	ActionType  string   `json:"action_type"` // e.g., "run_tool", "exploit", "manual_send"
	CommandHash []byte   `json:"command_hash"` // SHA‑256 of the sanitized argv
	RequestedAt time.Time `json:"requested_at"`
	DecidedAt   sql.NullTime `json:"decided_at,omitempty"`
	Decision    sql.NullString `json:"decision,omitempty"` // "allowed" or "denied"
	Comment     sql.NullString `json:"comment,omitempty"`
	JWTToken    string   `json:"jwt_token,omitempty"` // signed token handed to the caller
}

// NewApprovalTable ensures the table exists (run once on startup).
func NewApprovalTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS approvals (
		id SERIAL PRIMARY KEY,
		engagement_id INT NOT NULL REFERENCES engagements(id) ON DELETE CASCADE,
		requester_id INT NOT NULL REFERENCES users(id),
		approver_id INT REFERENCES users(id),
		action_type VARCHAR NOT NULL,
		command_hash BYTEA NOT NULL,
		requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		decided_at TIMESTAMPTZ,
		decision VARCHAR CHECK (decision IN ('allowed','denied') OR decision IS NULL),
		comment TEXT,
		jwt_token TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_approvals_engagement ON approvals(engagement_id);
	CREATE INDEX IF NOT EXISTS idx_approvals_requested ON approvals(requested_at);
	`
	_, err := db.Exec(query)
	return err
}

// CreateApproval inserts a pending request and returns the JWT that should be
// handed to the caller.
func CreateApproval(db *sql.DB, engagementID, requesterID int64, actionType string, commandHash []byte, ttlSeconds int) (*Approval, error) {
	// Build a simple JWT (HS256) – in a real build you would use the secret from /run/secrets/jwt_secret.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"engagement_id": engagementID,
		"action_type":   actionType,
		"command_hash":  fmt.Sprintf("%x", commandHash),
		"exp":           time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix(),
	})
	// secret is read from the environment or a file – here we use a hard‑coded fallback.
	secret := []byte("fake-jwt-secret-for-demo")
	signed, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	res := Approval{
		EngagementID: engagementID,
		RequesterID:  requesterID,
		ActionType:   actionType,
		CommandHash:  commandHash,
		RequestedAt:  time.Now(),
		JWTToken:     signed,
	}
	query := `
	INSERT INTO approvals (engagement_id, requester_id, action_type, command_hash, requested_at, jwt_token)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id
	`
	err = db.QueryRow(query,
		res.EngagementID,
		res.RequesterID,
		res.ActionType,
		res.CommandHash,
		res.RequestedAt,
		res.JWTToken).Scan(&res.ID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// DecideApproval marks an approval as allowed/denied.
func DecideApproval(db *sql.DB, approvalID int64, decision string, comment string, approverID int64) error {
	query := `
	UPDATE approvals
	SET decided_at = now(),
	    decision = $2,
	    comment = $3,
	    approver_id = $4
	WHERE id = $1 AND decided_at IS NULL
	`
	_, err := db.Exec(query, approvalID, decision, comment, approverID)
	return err
}

// GetPending returns all approvals that are still waiting for a decision.
func GetPending(db *sql.DB) ([]Approval, error) {
	rows, err := db.Query(`
	SELECT id, engagement_id, requester_id, action_type, command_hash,
	       requested_at, decided_at, decision, comment, jwt_token
	FROM approvals
	WHERE decided_at IS NULL
	ORDER BY requested_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Approval
	for rows.Next() {
		var a Approval
		if err := rows.Scan(&a.ID, &a.EngagementID, &a.RequesterID,
			&a.ActionType, &a.CommandHash, &a.RequestedAt,
			&a.DecidedAt, &a.Decision, &a.Comment, &a.JWTToken); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}
