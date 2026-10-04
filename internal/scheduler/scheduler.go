package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yourname/security-orchestrator/internal/logger"
)

// Scheduler wraps a cron.Cron and knows how to invoke the orchestrator’s
// internal job functions (e.g., periodic passive recon).
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a scheduler that will run jobs in the given ctx.
func NewScheduler(ctx context.Context) *Scheduler {
	c := cron.New(cron.WithSeconds()) // allow second‑resolution if you like
	// Example job: run passive recon every 6 hours.
	_, err := c.AddFunc("0 */6 * * *", func() {
		logger.Infof("Scheduler: triggering passive recon job")
	})
	if err != nil {
		panic(err)
	}
	return &Scheduler{cron: c}
}

// Start begins the scheduler; it runs until ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	s.cron.Start()
	<-ctx.Done()
	s.cron.Stop()
}
