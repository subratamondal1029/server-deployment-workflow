package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// middlewares
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", ORIGIN)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode != 0 {
		return
	}

	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(body)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWriter := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(responseWriter, r)

		statusCode := responseWriter.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		fmt.Printf("%s %s %d\n", r.Method, r.URL.Path, statusCode)
	})
}

// router handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Hello, World!")
}

var PORT string = os.Getenv("PORT")
var ORIGIN string = os.Getenv("ORIGIN")

func main() {

	if PORT == "" {
		log.Fatal("PORT is not set")
	}

	if ORIGIN == "" {
		log.Fatal("ORIGIN is not set")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)

	router := cors(mux)
	router = logging(router)

	fmt.Printf("Starting server on port %s\n", PORT)
	if err := http.ListenAndServe(":"+PORT, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
