// kaiops-azure-demo: a tiny stateless Go HTTP API used to exercise the
// KaiOps deployment pipeline (GitHub CI -> Artifact Registry -> ArgoCD -> AKS).
// It exposes /, /healthz, and /api/info so the app is observable.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type info struct {
	App       string `json:"app"`
	Version   string `json:"version"`
	Pod       string `json:"pod"`
	Host      string `json:"host"`
	Timestamp string `json:"timestamp"`
}

var (
	appName  = getenv("APP_NAME", "kaiops-azure-demo")
	version  = getenv("APP_VERSION", "v1.0.0")
	hostname, _ = os.Hostname()
)

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func writeJSON(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"service": appName, "status": "ok"})
}

func healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy", "app": appName})
}

func apiInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, info{
		App:       appName,
		Version:   version,
		Pod:       getenv("POD_NAME", ""),
		Host:      hostname,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/healthz", healthz).Methods("GET")
	r.HandleFunc("/api/info", apiInfo).Methods("GET")
	log.Printf("%s listening on :8080", appName)
	log.Fatal(http.ListenAndServe(":8080", r))
}
