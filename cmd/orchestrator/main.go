package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourname/security-orchestrator/internal/app"
)

func main() {
	// Simple sub‑command dispatcher.
	if len(os.Args) > 1 && os.Args[1] == "api" {
		app.RunAPI()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		app.RunCLI(os.Args[2:])
		return
	}
	// Default REPL – useful while developing.
	fmt.Println("=== Orchestrator CLI (demo) ===")
	fmt.Println("Available commands: init, admin, passive, active, triage, manual, report, stop")
	fmt.Println("Type 'help <cmd>' for details, or 'exit' to quit.")
	app.RunREPL()
	// Wait for SIGINT/SIGTERM to shut down cleanly.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\nShutting down…")
}
