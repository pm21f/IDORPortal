package handlers

import (
        "encoding/json"
        "net/http"
        "strconv"
        "strings"
        "cyclesync/models"
)

// PostRequest represents a post create/update request
type PostRequest struct {
        Title   string `json:"title"`
        Content string `json:"content"`
}

// PostsHandler handles requests for all posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
                // Get all posts or filter by user ID
                userIDStr := r.URL.Query().Get("user_id")
                if userIDStr != "" {
                        userID, err := strconv.Atoi(userIDStr)
                        if err != nil {
                                sendJSONResponse(w, false, "Invalid user ID", nil, http.StatusBadRequest)
                                return
                        }

                        posts, err := models.GetPostsByUserID(userID)
                        if err != nil {
                                sendJSONResponse(w, false, "Error fetching posts", nil, http.StatusInternalServerError)
                                return
                        }
                        sendJSONResponse(w, true, "", posts, http.StatusOK)
                        return
                }

                // Get all posts
                posts, err := models.GetAllPosts()
                if err != nil {
                        sendJSONResponse(w, false, "Error fetching posts", nil, http.StatusInternalServerError)
                        return
                }
                sendJSONResponse(w, true, "", posts, http.StatusOK)

        case http.MethodPost:
                // Create a new post
                session, ok := getSession(r)
                if !ok {
                        sendJSONResponse(w, false, "Not logged in", nil, http.StatusUnauthorized)
                        return
                }

                var req PostRequest
                err := json.NewDecoder(r.Body).Decode(&req)
                if err != nil {
                        sendJSONResponse(w, false, "Invalid request", nil, http.StatusBadRequest)
                        return
                }

                postID, err := models.CreatePost(session.UserID, req.Title, req.Content)
                if err != nil {
                        sendJSONResponse(w, false, "Error creating post", nil, http.StatusInternalServerError)
                        return
                }

                post, err := models.GetPostByID(postID)
                if err != nil {
                        sendJSONResponse(w, false, "Post created but could not retrieve details", nil, http.StatusInternalServerError)
                        return
                }
                sendJSONResponse(w, true, "Post created successfully", post, http.StatusCreated)

        default:
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// PostHandler handles requests for a specific post
// VULNERABLE TO IDOR: No authorization check for viewing/modifying posts
func PostHandler(w http.ResponseWriter, r *http.Request) {
        // Extract post ID from path
        idStr := strings.TrimPrefix(r.URL.Path, "/api/post/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                sendJSONResponse(w, false, "Invalid post ID", nil, http.StatusBadRequest)
                return
        }

        switch r.Method {
        case http.MethodGet:
                // Get post by ID
                post, err := models.GetPostByID(id)
                if err != nil {
                        sendJSONResponse(w, false, "Error fetching post", nil, http.StatusInternalServerError)
                        return
                }
                if post == nil {
                        sendJSONResponse(w, false, "Post not found", nil, http.StatusNotFound)
                        return
                }
                sendJSONResponse(w, true, "", post, http.StatusOK)

        case http.MethodPut:
                // Update post
                // VULNERABLE: No check if the currently logged-in user is the owner of the post

                var req PostRequest
                err := json.NewDecoder(r.Body).Decode(&req)
                if err != nil {
                        sendJSONResponse(w, false, "Invalid request", nil, http.StatusBadRequest)
                        return
                }

                // Update post
                err = models.UpdatePost(id, req.Title, req.Content)
                if err != nil {
                        sendJSONResponse(w, false, "Error updating post", nil, http.StatusInternalServerError)
                        return
                }

                // Get updated post
                post, err := models.GetPostByID(id)
                if err != nil {
                        sendJSONResponse(w, false, "Post updated but could not retrieve details", nil, http.StatusInternalServerError)
                        return
                }
                sendJSONResponse(w, true, "Post updated successfully", post, http.StatusOK)

        case http.MethodDelete:
                // Delete post
                // VULNERABLE: No check if the currently logged-in user is the owner of the post

                err := models.DeletePost(id)
                if err != nil {
                        sendJSONResponse(w, false, "Error deleting post", nil, http.StatusInternalServerError)
                        return
                }
                sendJSONResponse(w, true, "Post deleted successfully", nil, http.StatusOK)

        default:
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}
