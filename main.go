package main

import (
	"fmt"
	"gamecraft-backend/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowed := map[string]bool{
			"http://localhost:3000":     true,
			"http://172.18.126.70:3000": true,
		}

		if allowed[origin] {
			// IMPORTANT: set exact origin (not "*") when credentials are used
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// Optional but good: tell caches that response varies by Origin
			w.Header().Set("Vary", "Origin")
		}

		// Allow needed methods and headers for preflight and actual requests
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		// If your client sends custom headers, add them above

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	apiMux := http.NewServeMux()
	getMux := http.NewServeMux()

	routes.RegisterRouter(apiMux)
	routes.RegisterRouterGet(getMux)

	// Combine routes under prefixes
	http.Handle("/api/", enableCORS(http.StripPrefix("/api", apiMux)))
	http.Handle("/getapis/", enableCORS(http.StripPrefix("/getapis", getMux)))

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("❌ DATABASE_URL environment variable is not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("❓ Unhandled request: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	fmt.Println("🚀 Server running at http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
