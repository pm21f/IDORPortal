package handlers

import (
        "encoding/json"
        "net/http"
        "strconv"
        "strings"
        "cyclesync/models"
)

// UserUpdateRequest represents a user update request
type UserUpdateRequest struct {
        Username string `json:"username"`
        Email    string `json:"email"`
}

// UsersHandler handles requests for all users
func UsersHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
                // Get all users
                users, err := models.GetAllUsers()
                if err != nil {
                        sendJSONResponse(w, false, "Error fetching users", nil, http.StatusInternalServerError)
                        return
                }
                sendJSONResponse(w, true, "", users, http.StatusOK)
        default:
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// UserHandler handles requests for a specific user
// VULNERABLE TO IDOR: No authorization check on user access
func UserHandler(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from path
        idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                sendJSONResponse(w, false, "Invalid user ID", nil, http.StatusBadRequest)
                return
        }

        switch r.Method {
        case http.MethodGet:
                // Get user by ID
                user, err := models.GetUserByID(id)
                if err != nil {
                        sendJSONResponse(w, false, "Error fetching user", nil, http.StatusInternalServerError)
                        return
                }
                if user == nil {
                        sendJSONResponse(w, false, "User not found", nil, http.StatusNotFound)
                        return
                }
                sendJSONResponse(w, true, "", user.ToPublic(), http.StatusOK)

        case http.MethodPut:
                // Update user
                // VULNERABLE: No check if the currently logged-in user is updating their own profile

                var req UserUpdateRequest
                err := json.NewDecoder(r.Body).Decode(&req)
                if err != nil {
                        sendJSONResponse(w, false, "Invalid request", nil, http.StatusBadRequest)
                        return
                }

                // Update user
                err = models.UpdateUser(id, req.Username, req.Email)
                if err != nil {
                        sendJSONResponse(w, false, "Error updating user", nil, http.StatusInternalServerError)
                        return
                }

                // Get updated user
                user, err := models.GetUserByID(id)
                if err != nil {
                        sendJSONResponse(w, false, "User updated but could not retrieve details", nil, http.StatusInternalServerError)
                        return
                }
                sendJSONResponse(w, true, "User updated successfully", user.ToPublic(), http.StatusOK)

        case http.MethodDelete:
                // Delete user
                // VULNERABLE: No check if the currently logged-in user is deleting their own account

                err := models.DeleteUser(id)
                if err != nil {
                        sendJSONResponse(w, false, "Error deleting user", nil, http.StatusInternalServerError)
                        return
                }
                sendJSONResponse(w, true, "User deleted successfully", nil, http.StatusOK)

        default:
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}
