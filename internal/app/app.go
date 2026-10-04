package app

import (
	"log"
	"net/http"
	"time"

	"github.com/yourname/security-orchestrator/internal/api"
	"github.com/yourname/security-orchestrator/internal/logger"
)

// RunAPI starts the HTTP server.
func RunAPI() {
	r := mux.NewRouter()
	api.RegisterHandlers(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Println("API listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

// RunCLI forwards CLI sub‑commands to the appropriate internal packages.
func RunCLI(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: orchestrator cli <command> [args…]")
		return
	}
	switch args[0] {
	case "admin":
		adminCmd(args[1:])
	case "init":
		initCmd(args[1:])
	case "passive":
		passiveCmd(args[1:])
	case "active":
		activeCmd(args[1:])
	case "triage":
		triageCmd(args[1:])
	case "manual":
		manualCmd(args[1:])
	case "report":
		reportCmd(args[1:])
	case "stop":
		stopCmd()
	default:
		fmt.Printf("Unknown CLI command: %s\n", args[0])
	}
}

// Stub implementations – replace with real calls to internal packages.
func adminCmd(args []string) {
	fmt.Println("[admin] stub – would create/delete admin users")
}
func initCmd(args []string) {
	fmt.Println("[init] stub – would create a new engagement auth.yml")
}
func passiveCmd(args []string) {
	fmt.Println("[passive] stub – would run subfinder/amass")
}
func activeCmd(args []string) {
	fmt.Println("[active] stub – would run httpx/nmap/nuclei")
}
func triageCmd(args []string) {
	fmt.Println("[triage] stub – would call LLM to turn raw output into findings")
}
func manualCmd(args []string) {
	fmt.Println("[manual] stub – would show safe payload prompt and send request")
}
func reportCmd(args []string) {
	fmt.Println("[report] stub – would generate HTML/MD/PDF report")
}
func stopCmd() {
	fmt.Println("[stop] stub – would touch STOP file to halt all side‑cars")
}

// RunREPL provides a tiny interactive loop for demo purposes.
func RunREPL() {
	for {
		fmt.Print("orchestrator> ")
		var cmd string
		if _, err := fmt.Scanln(&cmd); err != nil {
			return
		}
		if cmd == "exit" || cmd == "quit" {
			fmt.Println("Goodbye!")
			return
		}
		fmt.Printf("You entered: %s\n", cmd)
	}
}
