package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourname/security-orchestrator/internal/logger"
)

// RegisterHandlers attaches a handful of endpoints to the router.
// In a production build you would expand this with all the real endpoints.
func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/auth/login", loginHandler).Methods("POST")

	r.HandleFunc("/engagements", listEngagementsHandler).Methods("GET")
	r.HandleFunc("/engagements", createEngagementHandler).Methods("POST")
	r.HandleFunc("/engagements/{id}", getEngagementHandler).Methods("GET")
	r.HandleFunc("/engagements/{id}", updateEngagementHandler).Methods("PUT")
	r.HandleFunc("/engagements/{id}", deleteEngagementHandler).Methods("DELETE")

	r.HandleFunc("/{any:.*}", notImplementedHandler).Methods(
		http.MethodGet, http.MethodPost, http.MethodPut,
		http.MethodDelete, http.MethodPatch, http.Options,
	)
}

// -------------------------------------------------------------------
// Handlers -----------------------------------------------------------

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if creds.Username == "harryadmin" && creds.Password == "StrongPass!2025!" {
		// In a real build you would sign a JWT with the secret from /run/secrets/jwt_secret.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": "fake-jwt-token-for-demo",
		})
		return
	}
	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}

// -------------------------------------------------------------------
// Engagement CRUD stubs ----------------------------------------------

func listEngagementsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]string{
		{"id": "placeholder-id", "name": "example"},
	})
}

func createEngagementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":   "new-id",
		"name": "created",
	})
}

func getEngagementHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":   id,
		"name": "example",
	})
}

func updateEngagementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func deleteEngagementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// -------------------------------------------------------------------
// Not‑implemented fallback -------------------------------------------

func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "not implemented in this skeleton",
	})
}
