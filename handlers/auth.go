package handlers

import (
        "encoding/json"
        "html/template"
        "net/http"
        "time"
        "cyclesync/models"
)

// Session represents a user session
type Session struct {
        UserID    int
        Username  string
        ExpiresAt time.Time
}

// Sessions store (in-memory for simplicity)
// In a production environment, use a more secure session store
var sessions = make(map[string]Session)

// Response represents a JSON response
type Response struct {
        Success bool        `json:"success"`
        Message string      `json:"message,omitempty"`
        Data    interface{} `json:"data,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
}

// SignupRequest represents a signup request
type SignupRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
}

// IndexHandler handles the root path
func IndexHandler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
                http.NotFound(w, r)
                return
        }

        tmpl, err := template.ParseFiles("templates/index.html")
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        tmpl.Execute(w, nil)
}

// LoginPageHandler renders the login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/login.html")
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        tmpl.Execute(w, nil)
}

// SignupPageHandler renders the signup page
func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
        tmpl, err := template.ParseFiles("templates/signup.html")
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        tmpl.Execute(w, nil)
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        var req LoginRequest
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
                sendJSONResponse(w, false, "Invalid request", nil, http.StatusBadRequest)
                return
        }

        // Get user by username
        user, err := models.GetUserByUsername(req.Username)
        if err != nil {
                sendJSONResponse(w, false, "Internal server error", nil, http.StatusInternalServerError)
                return
        }

        if user == nil || !user.VerifyPassword(req.Password) {
                sendJSONResponse(w, false, "Invalid username or password", nil, http.StatusUnauthorized)
                return
        }

        // Create session
        sessionID := generateSessionID()
        sessions[sessionID] = Session{
                UserID:    user.ID,
                Username:  user.Username,
                ExpiresAt: time.Now().Add(24 * time.Hour),
        }

        // Set session cookie
        http.SetCookie(w, &http.Cookie{
                Name:     "session",
                Value:    sessionID,
                Path:     "/",
                HttpOnly: true,
                MaxAge:   86400, // 24 hours
        })

        sendJSONResponse(w, true, "Login successful", user.ToPublic(), http.StatusOK)
}

// SignupHandler handles user registration
func SignupHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        var req SignupRequest
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
                sendJSONResponse(w, false, "Invalid request", nil, http.StatusBadRequest)
                return
        }

        // Check if username already exists
        existingUser, err := models.GetUserByUsername(req.Username)
        if err != nil {
                sendJSONResponse(w, false, "Internal server error", nil, http.StatusInternalServerError)
                return
        }
        if existingUser != nil {
                sendJSONResponse(w, false, "Username already taken", nil, http.StatusConflict)
                return
        }

        // Check if email already exists
        existingUser, err = models.GetUserByEmail(req.Email)
        if err != nil {
                sendJSONResponse(w, false, "Internal server error", nil, http.StatusInternalServerError)
                return
        }
        if existingUser != nil {
                sendJSONResponse(w, false, "Email already in use", nil, http.StatusConflict)
                return
        }

        // Create new user
        userID, err := models.CreateUser(req.Username, req.Email, req.Password)
        if err != nil {
                sendJSONResponse(w, false, "Failed to create user", nil, http.StatusInternalServerError)
                return
        }

        // Get the newly created user
        user, err := models.GetUserByID(userID)
        if err != nil {
                sendJSONResponse(w, false, "User created but could not retrieve details", nil, http.StatusInternalServerError)
                return
        }

        // Create session
        sessionID := generateSessionID()
        sessions[sessionID] = Session{
                UserID:    user.ID,
                Username:  user.Username,
                ExpiresAt: time.Now().Add(24 * time.Hour),
        }

        // Set session cookie
        http.SetCookie(w, &http.Cookie{
                Name:     "session",
                Value:    sessionID,
                Path:     "/",
                HttpOnly: true,
                MaxAge:   86400, // 24 hours
        })

        sendJSONResponse(w, true, "Signup successful", user.ToPublic(), http.StatusCreated)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Get session cookie
        cookie, err := r.Cookie("session")
        if err != nil {
                sendJSONResponse(w, false, "Not logged in", nil, http.StatusUnauthorized)
                return
        }

        // Delete session
        delete(sessions, cookie.Value)

        // Clear session cookie
        http.SetCookie(w, &http.Cookie{
                Name:     "session",
                Value:    "",
                Path:     "/",
                HttpOnly: true,
                MaxAge:   -1, // Delete cookie
        })

        sendJSONResponse(w, true, "Logout successful", nil, http.StatusOK)
}

// DashboardHandler renders the dashboard page
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
        // Check if user is logged in
        session, ok := getSession(r)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

        tmpl, err := template.ParseFiles("templates/dashboard.html")
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        tmpl.Execute(w, session)
}

// ProfileHandler renders the profile page
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
        // Check if user is logged in
        session, ok := getSession(r)
        if !ok {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
        }

        tmpl, err := template.ParseFiles("templates/profile.html")
        if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }
        tmpl.Execute(w, session)
}

// Helper functions

// sendJSONResponse sends a JSON response
func sendJSONResponse(w http.ResponseWriter, success bool, message string, data interface{}, statusCode int) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(Response{
                Success: success,
                Message: message,
                Data:    data,
        })
}

// getSession retrieves the current session from a request
func getSession(r *http.Request) (Session, bool) {
        cookie, err := r.Cookie("session")
        if err != nil {
                return Session{}, false
        }

        session, ok := sessions[cookie.Value]
        if !ok || time.Now().After(session.ExpiresAt) {
                return Session{}, false
        }

        return session, true
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
        return time.Now().Format("20060102150405") + ":" + string(randStringBytes(16))
}

// randStringBytes generates a random string of length n
func randStringBytes(n int) string {
        const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
        b := make([]byte, n)
        for i := range b {
                b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
                time.Sleep(1 * time.Nanosecond) // Ensure uniqueness
        }
        return string(b)
}
