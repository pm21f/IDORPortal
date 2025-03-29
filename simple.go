package main

import (
        "fmt"
        "log"
        "net/http"
        "os"
        "path/filepath"
        "time"
        "cyclesync/models"
        "cyclesync/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                start := time.Now()
                log.Printf("Started %s %s", r.Method, r.URL.Path)
                next.ServeHTTP(w, r)
                log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
        })
}

func main() {
        log.Println("Starting CycleSync application...")
        
        // Initialize database connection
        log.Println("Initializing database connection...")
        err := models.InitDB()
        if err != nil {
                log.Fatalf("Failed to connect to database: %v", err)
        }
        defer models.CloseDB()
        log.Println("Database connection successful")
        
        // Create tables if they don't exist
        log.Println("Creating database tables if needed...")
        err = models.CreateTables()
        if err != nil {
                log.Fatalf("Failed to create tables: %v", err)
        }
        log.Println("Database tables checked/created successfully")

        // Static file server
        log.Println("Setting up static file server...")
        fs := http.FileServer(http.Dir("static"))
        http.Handle("/static/", http.StripPrefix("/static/", fs))

        // Main routes
        log.Println("Registering route handlers...")
        http.HandleFunc("/", handlers.IndexHandler)
        http.HandleFunc("/login", handlers.LoginPageHandler)
        http.HandleFunc("/signup", handlers.SignupPageHandler)
        http.HandleFunc("/dashboard", handlers.DashboardHandler)
        http.HandleFunc("/profile", handlers.ProfileHandler)

        // API routes
        http.HandleFunc("/api/login", handlers.LoginHandler)
        http.HandleFunc("/api/signup", handlers.SignupHandler)
        http.HandleFunc("/api/logout", handlers.LogoutHandler)
        http.HandleFunc("/api/users", handlers.UsersHandler)
        http.HandleFunc("/api/user/", handlers.UserHandler) // Vulnerable to IDOR
        http.HandleFunc("/api/posts", handlers.PostsHandler)
        http.HandleFunc("/api/post/", handlers.PostHandler) // Vulnerable to IDOR

        // Serve on port 5000
        log.Println("Server starting on http://0.0.0.0:5000")
        log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}