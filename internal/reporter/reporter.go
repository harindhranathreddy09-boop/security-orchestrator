package reporter

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/yourname/security-orchestrator/internal/db"
)

// Report holds the data needed to render a template.
type Report struct {
	EngagementID int64
	GeneratedAt  string
	Findings     []Finding
}

// Finding is a tiny representation of a finding.
type Finding struct {
	Title       string
	Description string
	CVSS        string
	Severity    string
	Steps       []string
}

// NewReport builds a Report from the DB (stubbed).
func NewReport(db *db.DB, engagementID int64) (*Report, error) {
	// In a real build we would query findings, etc.
	// Here we return a dummy report.
	return &Report{
		EngagementID: engagementID,
		GeneratedAt:  time.Now().UTC().Format(time.RFC3339),
		Findings: []Finding{
			{
				Title:       "Example Finding",
				Description: "This is a placeholder finding.",
				CVSS:        "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:L/A:N",
				Severity:    "medium",
				Steps:       []string{"Step 1: do X", "Step 2: observe Y"},
			},
		},
	}, nil
}

// RenderHTML executes the HTML template and returns the byte slice.
func RenderHTML(tmplPath string, r *Report) ([]byte, error) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r); err != nil {
		return nil, err
	}
	return ioutil.ReadAll(&buf)
}

// WriteReport writes the rendered content to disk.
func WriteReport(dir string, filename string, data []byte) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, filename), data, 0o644)
}
