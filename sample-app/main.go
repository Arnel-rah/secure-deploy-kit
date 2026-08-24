package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"secure-deploy-kit/sample-app/internal/oidcverify"
)

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(healthResponse{
			Status: "UP",
			Time:   time.Now().UTC().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secure-deploy-kit demo app is running\n"))
	})

	issuerURL := os.Getenv("OIDC_ISSUER_URL")
	if issuerURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		verifier, err := oidcverify.NewVerifier(ctx, issuerURL)
		if err != nil {
			log.Fatalf("failed to initialize OIDC verifier for issuer %q: %v", issuerURL, err)
		}

		profileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, _ := oidcverify.ClaimsFromContext(r.Context())
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"subject": claims.Subject,
				"roles":   claims.Roles,
			})
		})

		mux.Handle("/api/profile", verifier.Middleware(profileHandler))
		log.Printf("OIDC enforcement active, issuer: %s", issuerURL)
	} else {
		log.Printf("OIDC_ISSUER_URL not set — /api/profile is disabled, running without auth enforcement")
		mux.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"error":"OIDC not configured on this deployment"}`, http.StatusNotImplemented)
		})
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
