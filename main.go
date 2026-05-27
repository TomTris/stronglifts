package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"stronglifts/db"
	"stronglifts/handlers"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	loadEnv(".env")

	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "stronglifts.db")

	if os.Getenv("GOOGLE_CLIENT_ID") == "" {
		log.Println("WARNING: GOOGLE_CLIENT_ID not set. Google login will not work.")
		log.Println("Set it in .env or as environment variable.")
	}

	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	h := handlers.New(database)
	mux := http.NewServeMux()

	// API routes
	h.RegisterRoutes(mux)

	// Static files from embedded frontend
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to load frontend: %v", err)
	}
	fileServer := http.FileServer(http.FS(distFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// For SPA: serve index.html for non-file routes
		path := r.URL.Path
		if path != "/" {
			// Try to serve the file
			if _, err := fs.Stat(distFS, path[1:]); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		// Serve index.html for SPA routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	fmt.Printf("StrongLifts 5x5 running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// loadEnv reads a .env file and sets variables that are not already in the environment.
func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // .env is optional
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove surrounding quotes
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		// Don't override existing env vars
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
