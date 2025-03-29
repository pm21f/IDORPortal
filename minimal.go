package main

import (
        "fmt"
        "log"
        "net/http"
)

// Create a minimal version for demonstration
func main() {
        // Simple handler function
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, `
                        <!DOCTYPE html>
                        <html>
                        <head>
                                <title>CycleSync - Simple Demo</title>
                                <style>
                                        body {
                                                font-family: Arial, sans-serif;
                                                margin: 0;
                                                padding: 50px;
                                                background-color: #f4f4f4;
                                                text-align: center;
                                        }
                                        .container {
                                                max-width: 800px;
                                                margin: 0 auto;
                                                padding: 20px;
                                                background-color: white;
                                                border-radius: 8px;
                                                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                                        }
                                        h1 {
                                                color: #50b3a2;
                                        }
                                        .card {
                                                border: 1px solid #ddd;
                                                padding: 15px;
                                                margin: 20px 0;
                                                border-radius: 4px;
                                                text-align: left;
                                        }
                                        .card h2 {
                                                margin-top: 0;
                                                color: #e8491d;
                                        }
                                        .footer {
                                                margin-top: 30px;
                                                font-size: 0.8em;
                                                color: #666;
                                        }
                                        a {
                                                color: #50b3a2;
                                                text-decoration: none;
                                        }
                                        a:hover {
                                                text-decoration: underline;
                                        }
                                </style>
                        </head>
                        <body>
                                <div class="container">
                                        <h1>CycleSync - IDOR Demo</h1>
                                        <p>This is a simplified demonstration of a web application with IDOR vulnerabilities.</p>
                                        
                                        <div class="card">
                                                <h2>User Endpoints</h2>
                                                <p>All users can be viewed at: <a href="/users">/users</a></p>
                                                <p>The IDOR vulnerability allows accessing any user profile without authentication at: <a href="/user/1">/user/1</a>, <a href="/user/2">/user/2</a>, etc.</p>
                                        </div>
                                        
                                        <div class="card">
                                                <h2>Post Endpoints</h2>
                                                <p>All posts can be viewed at: <a href="/posts">/posts</a></p>
                                                <p>The IDOR vulnerability allows accessing any post without authorization at: <a href="/post/1">/post/1</a>, <a href="/post/2">/post/2</a>, etc.</p>
                                        </div>

                                        <div class="footer">
                                                CycleSync Demo &copy; 2025
                                        </div>
                                </div>
                        </body>
                        </html>
                `)
        })

        // Users endpoint
        http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, `
                        <!DOCTYPE html>
                        <html>
                        <head>
                                <title>CycleSync - Users</title>
                                <style>
                                        body {
                                                font-family: Arial, sans-serif;
                                                margin: 0;
                                                padding: 50px;
                                                background-color: #f4f4f4;
                                        }
                                        .container {
                                                max-width: 800px;
                                                margin: 0 auto;
                                                padding: 20px;
                                                background-color: white;
                                                border-radius: 8px;
                                                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                                        }
                                        h1 {
                                                color: #50b3a2;
                                        }
                                        .user {
                                                border: 1px solid #ddd;
                                                padding: 15px;
                                                margin: 10px 0;
                                                border-radius: 4px;
                                        }
                                        .user h2 {
                                                margin-top: 0;
                                                color: #e8491d;
                                        }
                                        .footer {
                                                margin-top: 30px;
                                                font-size: 0.8em;
                                                color: #666;
                                                text-align: center;
                                        }
                                        a {
                                                color: #50b3a2;
                                                text-decoration: none;
                                        }
                                        a:hover {
                                                text-decoration: underline;
                                        }
                                        .nav {
                                                margin-bottom: 20px;
                                        }
                                </style>
                        </head>
                        <body>
                                <div class="container">
                                        <div class="nav">
                                                <a href="/">Home</a>
                                        </div>

                                        <h1>All Users</h1>
                                        
                                        <div class="user">
                                                <h2>Admin</h2>
                                                <p><strong>ID:</strong> 1</p>
                                                <p><strong>Email:</strong> admin@example.com</p>
                                                <p><strong>Joined:</strong> Yesterday</p>
                                                <p><a href="/user/1">View Profile</a></p>
                                        </div>
                                        
                                        <div class="user">
                                                <h2>User1</h2>
                                                <p><strong>ID:</strong> 2</p>
                                                <p><strong>Email:</strong> user1@example.com</p>
                                                <p><strong>Joined:</strong> Today</p>
                                                <p><a href="/user/2">View Profile</a></p>
                                        </div>

                                        <div class="footer">
                                                CycleSync Demo &copy; 2025
                                        </div>
                                </div>
                        </body>
                        </html>
                `)
        })

        // Individual user profile - IDOR Vulnerable
        http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
                // Extract user ID from URL
                userID := r.URL.Path[6:]
                
                // In a secure application, we would check if the logged-in user
                // has permission to view this profile
                
                var username, email string
                if userID == "1" {
                        username = "Admin"
                        email = "admin@example.com"
                } else if userID == "2" {
                        username = "User1"
                        email = "user1@example.com"
                } else {
                        // User not found
                        http.Error(w, "User not found", http.StatusNotFound)
                        return
                }
                
                fmt.Fprintf(w, `
                        <!DOCTYPE html>
                        <html>
                        <head>
                                <title>CycleSync - User Profile</title>
                                <style>
                                        body {
                                                font-family: Arial, sans-serif;
                                                margin: 0;
                                                padding: 50px;
                                                background-color: #f4f4f4;
                                        }
                                        .container {
                                                max-width: 800px;
                                                margin: 0 auto;
                                                padding: 20px;
                                                background-color: white;
                                                border-radius: 8px;
                                                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                                        }
                                        h1 {
                                                color: #50b3a2;
                                        }
                                        .profile {
                                                padding: 15px;
                                                margin: 10px 0;
                                                background-color: #f9f9f9;
                                                border-radius: 4px;
                                        }
                                        .actions {
                                                margin-top: 20px;
                                                padding-top: 10px;
                                                border-top: 1px solid #ddd;
                                        }
                                        .btn {
                                                display: inline-block;
                                                background-color: #e8491d;
                                                color: white;
                                                padding: 8px 16px;
                                                text-decoration: none;
                                                border-radius: 4px;
                                        }
                                        .btn:hover {
                                                background-color: #d43d1a;
                                        }
                                        .footer {
                                                margin-top: 30px;
                                                font-size: 0.8em;
                                                color: #666;
                                                text-align: center;
                                        }
                                        a {
                                                color: #50b3a2;
                                                text-decoration: none;
                                        }
                                        a:hover {
                                                text-decoration: underline;
                                        }
                                        .nav {
                                                margin-bottom: 20px;
                                        }
                                </style>
                        </head>
                        <body>
                                <div class="container">
                                        <div class="nav">
                                                <a href="/">Home</a> | 
                                                <a href="/users">All Users</a>
                                        </div>

                                        <h1>User Profile</h1>
                                        
                                        <div class="profile">
                                                <h2>%s</h2>
                                                <p><strong>ID:</strong> %s</p>
                                                <p><strong>Email:</strong> %s</p>
                                                
                                                <div class="actions">
                                                        <!-- IDOR Vulnerability: Anyone can access these actions -->
                                                        <a href="#" class="btn" onclick="alert('IDOR Vulnerability: Anyone can edit this profile!')">Edit Profile</a>
                                                        <a href="#" class="btn" onclick="alert('IDOR Vulnerability: Anyone can delete this profile!')">Delete Profile</a>
                                                </div>
                                        </div>

                                        <div class="footer">
                                                CycleSync Demo &copy; 2025
                                        </div>
                                </div>
                        </body>
                        </html>
                `, username, userID, email)
        })

        // Posts endpoint
        http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, `
                        <!DOCTYPE html>
                        <html>
                        <head>
                                <title>CycleSync - Posts</title>
                                <style>
                                        body {
                                                font-family: Arial, sans-serif;
                                                margin: 0;
                                                padding: 50px;
                                                background-color: #f4f4f4;
                                        }
                                        .container {
                                                max-width: 800px;
                                                margin: 0 auto;
                                                padding: 20px;
                                                background-color: white;
                                                border-radius: 8px;
                                                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                                        }
                                        h1 {
                                                color: #50b3a2;
                                        }
                                        .post {
                                                border: 1px solid #ddd;
                                                padding: 15px;
                                                margin: 10px 0;
                                                border-radius: 4px;
                                        }
                                        .post h2 {
                                                margin-top: 0;
                                                color: #e8491d;
                                        }
                                        .post .meta {
                                                color: #666;
                                                font-size: 0.8em;
                                                margin-bottom: 10px;
                                        }
                                        .footer {
                                                margin-top: 30px;
                                                font-size: 0.8em;
                                                color: #666;
                                                text-align: center;
                                        }
                                        a {
                                                color: #50b3a2;
                                                text-decoration: none;
                                        }
                                        a:hover {
                                                text-decoration: underline;
                                        }
                                        .nav {
                                                margin-bottom: 20px;
                                        }
                                </style>
                        </head>
                        <body>
                                <div class="container">
                                        <div class="nav">
                                                <a href="/">Home</a>
                                        </div>

                                        <h1>All Posts</h1>
                                        
                                        <div class="post">
                                                <h2>Admin Post</h2>
                                                <div class="meta">Posted by <a href="/user/1">Admin</a></div>
                                                <p>This is a post by admin with some content.</p>
                                                <p><a href="/post/1">View Full Post</a></p>
                                        </div>
                                        
                                        <div class="post">
                                                <h2>User Post</h2>
                                                <div class="meta">Posted by <a href="/user/2">User1</a></div>
                                                <p>This is a post by user1 with some different content.</p>
                                                <p><a href="/post/2">View Full Post</a></p>
                                        </div>

                                        <div class="footer">
                                                CycleSync Demo &copy; 2025
                                        </div>
                                </div>
                        </body>
                        </html>
                `)
        })

        // Individual post - IDOR Vulnerable
        http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
                // Extract post ID from URL
                postID := r.URL.Path[6:]
                
                // In a secure application, we would check if the post is public
                // or if the logged-in user has permission to view it
                
                var title, content, authorID, authorName string
                if postID == "1" {
                        title = "Admin Post"
                        content = "This is a detailed post by admin with some content. In a real application, this could contain sensitive information that should only be accessible to authorized users."
                        authorID = "1"
                        authorName = "Admin"
                } else if postID == "2" {
                        title = "User Post"
                        content = "This is a detailed post by user1 with some different content. This post might contain personal information that shouldn't be accessible to just anyone."
                        authorID = "2"
                        authorName = "User1"
                } else {
                        // Post not found
                        http.Error(w, "Post not found", http.StatusNotFound)
                        return
                }
                
                fmt.Fprintf(w, `
                        <!DOCTYPE html>
                        <html>
                        <head>
                                <title>CycleSync - %s</title>
                                <style>
                                        body {
                                                font-family: Arial, sans-serif;
                                                margin: 0;
                                                padding: 50px;
                                                background-color: #f4f4f4;
                                        }
                                        .container {
                                                max-width: 800px;
                                                margin: 0 auto;
                                                padding: 20px;
                                                background-color: white;
                                                border-radius: 8px;
                                                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                                        }
                                        h1 {
                                                color: #50b3a2;
                                        }
                                        .post {
                                                padding: 15px;
                                                margin: 10px 0;
                                                background-color: #f9f9f9;
                                                border-radius: 4px;
                                        }
                                        .post .meta {
                                                color: #666;
                                                font-size: 0.8em;
                                                margin-bottom: 10px;
                                        }
                                        .post .content {
                                                line-height: 1.6;
                                        }
                                        .actions {
                                                margin-top: 20px;
                                                padding-top: 10px;
                                                border-top: 1px solid #ddd;
                                        }
                                        .btn {
                                                display: inline-block;
                                                background-color: #e8491d;
                                                color: white;
                                                padding: 8px 16px;
                                                text-decoration: none;
                                                border-radius: 4px;
                                        }
                                        .btn:hover {
                                                background-color: #d43d1a;
                                        }
                                        .footer {
                                                margin-top: 30px;
                                                font-size: 0.8em;
                                                color: #666;
                                                text-align: center;
                                        }
                                        a {
                                                color: #50b3a2;
                                                text-decoration: none;
                                        }
                                        a:hover {
                                                text-decoration: underline;
                                        }
                                        .nav {
                                                margin-bottom: 20px;
                                        }
                                </style>
                        </head>
                        <body>
                                <div class="container">
                                        <div class="nav">
                                                <a href="/">Home</a> | 
                                                <a href="/posts">All Posts</a>
                                        </div>

                                        <h1>Post Details</h1>
                                        
                                        <div class="post">
                                                <h2>%s</h2>
                                                <div class="meta">Posted by <a href="/user/%s">%s</a></div>
                                                <div class="content">
                                                        <p>%s</p>
                                                </div>
                                                
                                                <div class="actions">
                                                        <!-- IDOR Vulnerability: Anyone can access these actions -->
                                                        <a href="#" class="btn" onclick="alert('IDOR Vulnerability: Anyone can edit this post!')">Edit Post</a>
                                                        <a href="#" class="btn" onclick="alert('IDOR Vulnerability: Anyone can delete this post!')">Delete Post</a>
                                                </div>
                                        </div>

                                        <div class="footer">
                                                CycleSync Demo &copy; 2025
                                        </div>
                                </div>
                        </body>
                        </html>
                `, title, title, authorID, authorName, content)
        })

        // Start server
        log.Println("Starting CycleSync application on http://0.0.0.0:5000")
        log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}